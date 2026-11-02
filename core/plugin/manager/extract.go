package manager

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extractFile(fpath string, file *zip.File) error {
	if file.UncompressedSize64 > uint64(maxPluginFileSize) {
		return fmt.Errorf("file %s exceeds maximum allowed size (%d bytes)", file.Name, maxPluginFileSize)
	}

	outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {

		}
	}(outFile)

	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer func(rc io.ReadCloser) {
		err := rc.Close()
		if err != nil {

		}
	}(rc)

	limitedReader := io.LimitReader(rc, maxPluginFileSize+1)
	written, err := io.Copy(outFile, limitedReader)
	if err != nil {
		return err
	}
	if written > maxPluginFileSize {
		return fmt.Errorf("file %s exceeded maximum decompression limit", file.Name)
	}

	mode := file.Mode() & 0755
	if mode&0111 == 0 {
		mode = 0644
	} else {
		mode = 0755
	}

	return os.Chmod(fpath, mode)
}

func unpackPlugin(archivePath, destDir string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open plugin archive: %w", err)
	}
	defer func(reader *zip.ReadCloser) {
		_ = reader.Close()
	}(reader)

	absDestDir, err := filepath.Abs(destDir)
	if err != nil {
		return fmt.Errorf("failed to resolve destination directory: %w", err)
	}
	absDestDir = filepath.Clean(absDestDir)
	if err := os.MkdirAll(absDestDir, 0775); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	var totalUncompressedBytes int64

	for _, file := range reader.File {
		// Reject symlinks to prevent symlink traversal
		if file.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("illegal symlink detected in archive: %s", file.Name)
		}

		if filepath.IsAbs(file.Name) {
			return fmt.Errorf("illegal absolute file path detected in archive: %s", file.Name)
		}

		normalizedName := strings.ReplaceAll(file.Name, "\\", "/")
		if strings.Contains(normalizedName, "../") || normalizedName == ".." || strings.HasPrefix(normalizedName, "/") {
			return fmt.Errorf("illegal file path detected in archive: %s", file.Name)
		}

		fpath := filepath.Clean(filepath.Join(absDestDir, file.Name))
		destPrefix := absDestDir + string(filepath.Separator)
		if fpath != absDestDir && !strings.HasPrefix(fpath, destPrefix) {
			return fmt.Errorf("illegal file path detected in archive: %s", file.Name)
		}

		dirPath := filepath.Dir(fpath)
		if evalDir, err := filepath.EvalSymlinks(dirPath); err == nil {
			if evalDest, err := filepath.EvalSymlinks(absDestDir); err == nil {
				if evalDir != evalDest && !strings.HasPrefix(evalDir, evalDest+string(filepath.Separator)) {
					return fmt.Errorf("illegal directory traversal detected in archive: %s", file.Name)
				}
			}
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, 0775); err != nil {
				return err
			}
			continue
		}

		totalUncompressedBytes += int64(file.UncompressedSize64)
		if totalUncompressedBytes > maxPluginTotalSize {
			return fmt.Errorf("archive total uncompressed size exceeds limit (%d bytes)", maxPluginTotalSize)
		}

		if err := os.MkdirAll(dirPath, 0775); err != nil {
			return err
		}

		if err := extractFile(fpath, file); err != nil {
			return err
		}
	}

	return nil
}
