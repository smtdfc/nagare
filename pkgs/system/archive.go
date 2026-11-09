package system

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MaxArchiveExtractBytes caps the total uncompressed size of an extraction to
// protect against archive bombs. The agent must never write unbounded output
// to the host filesystem.
const MaxArchiveExtractBytes = 2 << 30 // 2 GiB

// ArchiveFormat enumerates the archive formats supported by the helper.
type ArchiveFormat string

const (
	ArchiveFormatZip     ArchiveFormat = "zip"
	ArchiveFormatTar     ArchiveFormat = "tar"
	ArchiveFormatTarGz   ArchiveFormat = "tar.gz"
	ArchiveFormatUnknown ArchiveFormat = "unknown"
)

// DetectArchiveFormat infers the archive format from a file extension.
func DetectArchiveFormat(path string) ArchiveFormat {
	lower := strings.ToLower(path)
	switch {
	case strings.HasSuffix(lower, ".zip"):
		return ArchiveFormatZip
	case strings.HasSuffix(lower, ".tar.gz"), strings.HasSuffix(lower, ".tgz"):
		return ArchiveFormatTarGz
	case strings.HasSuffix(lower, ".tar"):
		return ArchiveFormatTar
	default:
		return ArchiveFormatUnknown
	}
}

// ArchiveCreate packs the given source paths (files or directories) into the
// archive at dest. The format is detected from the dest extension.
func ArchiveCreate(sources []string, dest string) error {
	format := DetectArchiveFormat(dest)
	if format == ArchiveFormatUnknown {
		return fmt.Errorf("unsupported archive format: '%s' (supported: .zip, .tar, .tar.gz, .tgz)", filepath.Ext(dest))
	}

	var entries []archiveEntry
	for _, source := range sources {
		resolved, err := ResolveFilePath(source)
		if err != nil {
			return err
		}
		info, err := os.Stat(resolved)
		if err != nil {
			return fmt.Errorf("source does not exist: %w", err)
		}

		if info.IsDir() {
			// Keep the directory's own name as the archive root, matching
			// how "zip -r" and "tar -cf" behave on the CLI.
			if err := collectDir(resolved, filepath.Base(resolved), &entries); err != nil {
				return err
			}
		} else {
			entries = append(entries, archiveEntry{path: resolved, name: filepath.Base(resolved), size: info.Size()})
		}
	}

	if len(entries) == 0 {
		return fmt.Errorf("no files found in the provided sources")
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	if _, err := os.Stat(dest); err == nil {
		return fmt.Errorf("destination already exists: %s", dest)
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create archive: %w", err)
	}
	defer out.Close()

	var writer io.Writer = out
	if format == ArchiveFormatTarGz {
		gz := gzip.NewWriter(out)
		defer gz.Close()
		writer = gz
	}

	switch format {
	case ArchiveFormatZip:
		return writeZip(writer, entries)
	default:
		return writeTar(writer, entries)
	}
}

// ArchiveExtract unpacks the archive at path into dest, creating dest when
// missing. Entry paths are validated against zip-slip attacks and the total
// uncompressed size is capped. It returns the destination directory contents
// were unpacked into.
func ArchiveExtract(path string, dest string) (string, error) {
	format := DetectArchiveFormat(path)
	if format == ArchiveFormatUnknown {
		return "", fmt.Errorf("unsupported archive format: '%s'", filepath.Ext(path))
	}

	if dest == "" {
		dest = filepath.Dir(path)
	}
	dest = filepath.Clean(dest)

	if err := os.MkdirAll(dest, 0o755); err != nil {
		return "", fmt.Errorf("failed to create destination directory: %w", err)
	}

	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open archive: %w", err)
	}
	defer f.Close()

	var total int64
	switch format {
	case ArchiveFormatZip:
		sizeRaw, err := fStat(f)
		if err != nil {
			return "", fmt.Errorf("failed to stat archive: %w", err)
		}
		zr, err := zip.NewReader(f, sizeRaw)
		if err != nil {
			return "", fmt.Errorf("failed to read zip archive: %w", err)
		}
		for _, entry := range zr.File {
			if strings.HasSuffix(entry.Name, "/") {
				continue // directory entry, created lazily by file entries
			}
			rc, err := entry.Open()
			if err != nil {
				return "", fmt.Errorf("failed to open zip entry %s: %w", entry.Name, err)
			}
			_, err = extractEntry(rc, entry.Name, dest, &total)
			rc.Close()
			if err != nil {
				return "", err
			}
		}
	default:
		var tr *tar.Reader
		if format == ArchiveFormatTarGz {
			gz, err := gzip.NewReader(f)
			if err != nil {
				return "", fmt.Errorf("failed to read gzip archive: %w", err)
			}
			defer gz.Close()
			tr = tar.NewReader(gz)
		} else {
			tr = tar.NewReader(f)
		}

		for {
			header, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return "", fmt.Errorf("failed to read tar entry: %w", err)
			}
			if header.Typeflag == tar.TypeDir {
				dir, err := safeJoin(dest, header.Name)
				if err != nil {
					return "", err
				}
				if err := os.MkdirAll(dir, os.FileMode(header.Mode)); err != nil {
					return "", fmt.Errorf("failed to create directory %s: %w", header.Name, err)
				}
				continue
			}
			if header.Typeflag != tar.TypeReg {
				continue
			}
			if _, err := extractEntry(tr, header.Name, dest, &total); err != nil {
				return "", err
			}
		}
	}

	return dest, nil
}

type archiveEntry struct {
	path string
	name string
	size int64
}

func fStat(f *os.File) (int64, error) {
	info, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// safeJoin validates that a relative entry name stays inside dest.
func safeJoin(dest string, name string) (string, error) {
	if strings.ContainsRune(name, 0) {
		return "", fmt.Errorf("invalid entry name: %s", name)
	}
	cleaned := filepath.Clean(name)
	if filepath.IsAbs(cleaned) || strings.HasPrefix(cleaned, "..") {
		return "", fmt.Errorf("entry '%s' escapes the destination directory (zip-slip blocked)", name)
	}
	return filepath.Join(dest, cleaned), nil
}

// extractEntry streams one entry into dest under the given name.
func extractEntry(reader io.Reader, name string, dest string, total *int64) (int64, error) {
	target, err := safeJoin(dest, name)
	if err != nil {
		return 0, err
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return 0, fmt.Errorf("failed to create parent directory for %s: %w", name, err)
	}

	out, err := os.Create(target)
	if err != nil {
		return 0, fmt.Errorf("failed to create file for entry %s: %w", name, err)
	}
	defer out.Close()

	var written int64
	buf := make([]byte, 32*1024)
	for {
		n, readErr := reader.Read(buf)
		if n > 0 {
			*total += int64(n)
			if *total > MaxArchiveExtractBytes {
				return written, fmt.Errorf("archive exceeds the %d byte extraction limit (possible archive bomb)", MaxArchiveExtractBytes)
			}
			if _, writeErr := out.Write(buf[:n]); writeErr != nil {
				return written, fmt.Errorf("failed to write entry %s: %w", name, writeErr)
			}
			written += int64(n)
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return written, fmt.Errorf("failed to read entry %s: %w", name, readErr)
		}
	}

	return written, nil
}

// collectDir recursively appends all regular files under current to entries,
// naming them below the prefix that represents the packed root directory.
func collectDir(current string, prefix string, entries *[]archiveEntry) error {
	dirEntries, err := os.ReadDir(current)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", current, err)
	}

	for _, dirEntry := range dirEntries {
		path := filepath.Join(current, dirEntry.Name())
		name := filepath.Join(prefix, dirEntry.Name())
		if dirEntry.IsDir() {
			if err := collectDir(path, name, entries); err != nil {
				return err
			}
			continue
		}
		info, err := dirEntry.Info()
		if err != nil {
			continue
		}
		if info.IsDir() || !dirEntry.Type().IsRegular() {
			continue
		}
		*entries = append(*entries, archiveEntry{path: path, name: name, size: info.Size()})
	}

	return nil
}

// writeZip writes a ZIP archive containing the entries into writer.
func writeZip(writer io.Writer, entries []archiveEntry) error {
	zw := zip.NewWriter(writer)
	defer zw.Close()

	seen := map[string]int{}
	for _, entry := range entries {
		name := uniqueName(seen, entry.name)
		header := &zip.FileHeader{Name: name, Method: zip.Deflate}
		mtime := mtimeOf(entry.path)
		if !mtime.IsZero() {
			header.SetModTime(mtime)
		}

		zf, err := zw.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create zip entry %s: %w", name, err)
		}
		if err := copyFileInto(zf, entry.path); err != nil {
			return err
		}
	}

	return zw.Close()
}

// writeTar writes a (possibly gzip-compressed) TAR archive of the entries.
func writeTar(writer io.Writer, entries []archiveEntry) error {
	tw := tar.NewWriter(writer)
	defer tw.Close()

	seen := map[string]int{}
	for _, entry := range entries {
		name := uniqueName(seen, entry.name)
		header := &tar.Header{Name: name, Size: entry.size, Mode: 0o644}
		mtime := mtimeOf(entry.path)
		if !mtime.IsZero() {
			header.ModTime = mtime
		}
		if err := tw.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write tar header %s: %w", name, err)
		}
		if err := copyFileInto(tw, entry.path); err != nil {
			return err
		}
	}

	return tw.Close()
}

// copyFileInto streams the file contents into the archive writer.
func copyFileInto(dst io.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer f.Close()

	if _, err := io.Copy(dst, f); err != nil {
		return fmt.Errorf("failed to add %s to archive: %w", path, err)
	}
	return nil
}

// uniqueName deduplicates entry names by appending a numeric suffix.
func uniqueName(seen map[string]int, name string) string {
	seen[name]++
	if seen[name] > 1 {
		ext := filepath.Ext(name)
		return fmt.Sprintf("%s_%d%s", strings.TrimSuffix(name, ext), seen[name], ext)
	}
	return name
}

// mtimeOf resolves the file modification time, ignoring stat errors.
func mtimeOf(path string) (mtime time.Time) {
	info, err := os.Stat(path)
	if err == nil {
		mtime = info.ModTime()
	}
	return mtime
}

// ArchiveSourceFiles returns every regular file under source: the file
// itself when source is a file, otherwise all files in its directory tree.
func ArchiveSourceFiles(source string) ([]string, error) {
	resolved, err := ResolveFilePath(source)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(resolved)
	if err != nil {
		return nil, fmt.Errorf("source does not exist: %w", err)
	}
	if !info.IsDir() {
		return []string{resolved}, nil
	}

	var files []string
	walkErr := filepath.WalkDir(resolved, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}
		if entry.IsDir() || !entry.Type().IsRegular() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}

	return files, nil
}
