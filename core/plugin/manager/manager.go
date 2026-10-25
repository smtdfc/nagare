package manager

import (
	"archive/zip"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/shirou/gopsutil/v4/process"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/logger"
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
	"github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/pkgs/helpers"
	"github.com/smtdfc/nagare/pkgs/paths"
	"github.com/smtdfc/nagare/plugin/metadata"
)

const (
	maxPluginFileSize  int64 = 250 * 1024 * 1024 // 250MB per file limit
	maxPluginTotalSize int64 = 500 * 1024 * 1024 // 500MB total uncompressed limit
)

func computeFileChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to hash file %s: %w", filePath, err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func verifyFileSignature(binFile, sigFile string) error {
	sigData, err := os.ReadFile(sigFile)
	if err != nil {
		return fmt.Errorf("failed to read signature file %s: %w", sigFile, err)
	}

	expectedSig := strings.TrimSpace(string(sigData))
	if expectedSig == "" {
		return errors.New("signature file is empty")
	}

	file, err := os.Open(binFile)
	if err != nil {
		return fmt.Errorf("failed to open binary file %s: %w", binFile, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("failed to compute binary hash: %w", err)
	}

	actualSig := hex.EncodeToString(hasher.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(strings.ToLower(expectedSig)), []byte(strings.ToLower(actualSig))) != 1 {
		return fmt.Errorf("signature verification failed: expected %s, got %s", expectedSig, actualSig)
	}

	return nil
}

func (p *PluginManager) loadTrustedPublicKeys() ([]ed25519.PublicKey, error) {
	entries, err := os.ReadDir(paths.TrustedDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read trusted key directory %s: %w", paths.TrustedDir, err)
	}

	keys := make([]ed25519.PublicKey, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		keyPath := filepath.Join(paths.TrustedDir, entry.Name())
		raw, err := os.ReadFile(keyPath)
		if err != nil {
			p.logger.Warn("Failed to read trusted plugin key", "file", keyPath, "error", err)
			continue
		}

		key, err := metadata.ParsePublicKey(string(raw))
		if err != nil {
			p.logger.Warn("Skipping invalid trusted plugin key", "file", keyPath, "error", err)
			continue
		}
		keys = append(keys, key)
	}

	return keys, nil
}

func (p *PluginManager) verifyArchiveSignature(archivePath string, detachedSignaturePath ...string) (string, error) {
	signaturePath := archivePath + ".sig"
	if len(detachedSignaturePath) > 0 && strings.TrimSpace(detachedSignaturePath[0]) != "" {
		signaturePath = detachedSignaturePath[0]
	}

	rawSignature, err := os.ReadFile(signaturePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			p.logger.Warn("Plugin archive is not signed", "plugin", archivePath)
			return plugin.SignatureStatusUnsigned, nil
		}

		return "", fmt.Errorf("failed to read signature file %s: %w", signaturePath, err)
	}

	signature, err := metadata.ParseSignatureFile(rawSignature)
	if err != nil {
		return "", fmt.Errorf("failed to parse signature file %s: %w", signaturePath, err)
	}

	trustedKeys, err := p.loadTrustedPublicKeys()
	if err != nil {
		return "", err
	}
	if len(trustedKeys) == 0 {
		p.logger.Warn(
			"Plugin archive is signed but no trusted public key is available; signature was not verified",
			"plugin", archivePath,
			"trustedDir", paths.TrustedDir,
		)
		return plugin.SignatureStatusUnverified, nil
	}

	for _, publicKey := range trustedKeys {
		err := metadata.VerifyArchive(archivePath, signature, publicKey)
		if err == nil {
			return plugin.SignatureStatusVerified, nil
		}
		if !errors.Is(err, metadata.ErrInvalidSignature) {
			return "", err
		}
	}

	return "", metadata.ErrInvalidSignature
}

func (p *PluginManager) verifyPluginChecksum(plg *plugin.Plugin) error {
	expected := strings.TrimSpace(plg.BinChecksum)
	if expected == "" {
		p.logger.Warn(
			"Plugin checksum is missing, skipping binary integrity check",
			"packageName", plg.PackageName,
		)
		return nil
	}

	actual, err := computeFileChecksum(plg.Bin)
	if err != nil {
		p.logger.Error("Failed to compute plugin binary checksum", "error", err, "packageName", plg.PackageName)
		return custom_errors.ErrPluginChecksumMismatch
	}

	if subtle.ConstantTimeCompare([]byte(strings.ToLower(expected)), []byte(actual)) != 1 {
		p.logger.Error(
			"Plugin binary checksum mismatch",
			"packageName", plg.PackageName,
			"expected", expected,
			"actual", actual,
		)
		return custom_errors.ErrPluginChecksumMismatch
	}

	return nil
}

func (p *PluginManager) ensureChecksum(ctx context.Context, plg *plugin.Plugin) {
	if strings.TrimSpace(plg.BinChecksum) != "" {
		return
	}

	checksum, err := computeFileChecksum(plg.Bin)
	if err != nil {
		p.logger.Warn("Failed to backfill plugin checksum", "error", err, "packageName", plg.PackageName)
		return
	}

	plg.BinChecksum = checksum
	if err := p.pluginRepo.Update(ctx, p.pluginMapper.ToEntity(plg)); err != nil {
		p.logger.Warn("Failed to persist plugin checksum", "error", err, "packageName", plg.PackageName)
	}
}

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

func loadMetadata(metadataFile string) (*metadata.PluginMetadata, error) {
	_, err := os.Stat(metadataFile)
	if errors.Is(err, os.ErrNotExist) {
		return nil, errors.New("metadata file not exist")
	}

	metadataContent, err := os.ReadFile(metadataFile)
	if err != nil {
		return nil, err
	}

	pluginMetadata, err := helpers.UnmarshalJson[metadata.PluginMetadata](string(metadataContent))
	if err != nil {
		return nil, err
	}

	return pluginMetadata, nil
}

type PluginManager struct {
	mu           sync.RWMutex
	pluginRepo   *repositories.PluginRepository
	pluginMapper *mappers.PluginMapper
	connectCodes map[string]string
	socketPath   string
	logger       *logger.BaseLogger
}

func (p *PluginManager) getPluginByID(ctx context.Context, id string) (*plugin.Plugin, error) {
	pluginEntity, err := p.pluginRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if pluginEntity == nil {
		return nil, custom_errors.ErrPluginNotFound
	}

	return p.pluginMapper.ToDomain(pluginEntity), nil
}

func (p *PluginManager) GetListPlugin(ctx context.Context) ([]*plugin.Plugin, error) {
	ents, err := p.pluginRepo.FindAll(ctx)
	if err != nil {
		return nil, custom_errors.ErrGetListPluginFailed
	}

	return p.pluginMapper.ToDomains(ents), nil
}

func (p *PluginManager) StopPlugin(plugin *plugin.Plugin) error {
	pidPath := plugin.Bin + ".pid"

	data, err := os.ReadFile(pidPath)
	if err != nil {
		if os.IsNotExist(err) {
			p.logger.Warn("PID file not found, skipping plugin stop", "packageName", plugin.PackageName)
			return nil
		}
		p.logger.Error("Failed to read PID file", "error", err, "packageName", plugin.PackageName)
		return err
	}

	var pid int
	if _, err := fmt.Sscanf(string(data), "%d", &pid); err != nil {
		p.logger.Error("Invalid PID format in file", "error", err, "packageName", plugin.PackageName)
		_ = os.Remove(pidPath)
		return err
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		p.logger.Error("Failed to find process", "pid", pid, "error", err, "packageName", plugin.PackageName)
		_ = os.Remove(pidPath)
		return err
	}

	if err := proc.Kill(); err != nil {
		p.logger.Error("Failed to kill process", "pid", pid, "error", err, "packageName", plugin.PackageName)
	}

	if err := os.Remove(pidPath); err != nil {
		p.logger.Warn("Failed to remove PID file", "error", err, "packageName", plugin.PackageName)
	}

	p.mu.Lock()
	delete(p.connectCodes, plugin.PackageName)
	p.mu.Unlock()
	return nil
}

func (p *PluginManager) StartPlugin(plugin *plugin.Plugin) error {
	if err := p.verifyPluginChecksum(plugin); err != nil {
		return err
	}

	connectCode := uuid.New().String()
	p.mu.Lock()
	p.connectCodes[plugin.PackageName] = connectCode
	socketPath := p.socketPath
	p.mu.Unlock()

	cmd := exec.Command(plugin.Bin)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if socketPath == "" {
		socketPath = paths.PluginSocketPath
	}
	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf("NAGARE_PLUGIN_CONNECT_CODE=%s", connectCode),
		fmt.Sprintf("NAGARE_PLUGIN_SOCKET_PATH=%s", socketPath),
		fmt.Sprintf("NAGARE_PLUGIN_LOG_FILE=%s", filepath.Join(paths.PluginLogDir, plugin.PackageName+".log")),
		fmt.Sprintf("NAGARE_PLUGIN_CONFIG_DIR=%s", filepath.Join(paths.PluginConfigDir, plugin.PackageName)),
	)

	err := cmd.Start()
	if err != nil {
		p.logger.Error("Start plugin failed", "error", err, "packageName", plugin.PackageName)
		return custom_errors.ErrStartPluginFailed
	}

	pidPath := plugin.Bin + ".pid"
	pidStr := fmt.Sprintf("%d", cmd.Process.Pid)
	if writeErr := os.WriteFile(pidPath, []byte(pidStr), 0644); writeErr != nil {
		p.logger.Error("Failed to write plugin PID file", "error", writeErr, "packageName", plugin.PackageName)
	}

	return nil
}

func (p *PluginManager) Install(ctx context.Context, pluginPath string) (*plugin.Plugin, error) {
	return p.install(ctx, pluginPath, "")
}

// InstallWithSignature installs an archive with a detached signature that is
// not stored next to the archive, as happens when a web client uploads the
// .sig file separately.
func (p *PluginManager) InstallWithSignature(ctx context.Context, pluginPath, signaturePath string) (*plugin.Plugin, error) {
	return p.install(ctx, pluginPath, signaturePath)
}

func (p *PluginManager) install(ctx context.Context, pluginPath, signaturePath string) (*plugin.Plugin, error) {
	_, err := os.Stat(pluginPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, custom_errors.ErrPluginNotFound
		}

		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrInstallPluginFailed
	}

	signatureStatus, err := p.verifyArchiveSignature(pluginPath, signaturePath)
	if err != nil {
		p.logger.Error("Install plugin failed: archive signature verification failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrPluginSignatureInvalid
	}

	id := uuid.New().String()
	tempDir := filepath.Join(paths.TempDir, id)
	metadataFile := filepath.Join(tempDir, "metadata.json")
	err = os.MkdirAll(tempDir, 0775)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	err = unpackPlugin(
		pluginPath,
		tempDir,
	)
	if err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrUnpackPluginFailed
	}

	pluginMetadata, err := loadMetadata(metadataFile)
	if err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrLoadPluginMetadataFailed
	}

	err = pluginMetadata.Validate()
	if err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrPluginMetadataInvalid
	}

	binFile, isExist := pluginMetadata.Bin[runtime.GOOS]
	if !isExist {
		return nil, custom_errors.ErrPluginBinaryMissing
	}

	cleanRelBin := filepath.Clean(binFile)
	normalizedBin := strings.ReplaceAll(cleanRelBin, "\\", "/")
	if filepath.IsAbs(cleanRelBin) || strings.HasPrefix(normalizedBin, "../") || normalizedBin == ".." || strings.Contains(normalizedBin, "/../") {
		p.logger.Error("Install plugin failed: binary path attempts traversal", "bin", binFile)
		return nil, custom_errors.ErrPluginMetadataInvalid
	}

	tempBinFile := filepath.Join(tempDir, cleanRelBin)
	tempSigFile := filepath.Join(tempDir, pluginMetadata.PackageName+".sig")
	if err := verifyFileSignature(tempBinFile, tempSigFile); err != nil {
		p.logger.Error("Install plugin failed: signature verification failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrPluginSignatureInvalid
	}

	binChecksum, err := computeFileChecksum(tempBinFile)
	if err != nil {
		p.logger.Error("Install plugin failed: compute binary checksum", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrInstallPluginFailed
	}

	pluginDir := filepath.Join(paths.PluginDir, pluginMetadata.PackageName)
	cleanPluginDir := filepath.Clean(pluginDir)
	destBinFile := filepath.Clean(filepath.Join(cleanPluginDir, cleanRelBin))
	if !strings.HasPrefix(destBinFile, cleanPluginDir+string(filepath.Separator)) {
		p.logger.Error("Install plugin failed: binary path escapes plugin directory", "destBinFile", destBinFile)
		return nil, custom_errors.ErrInstallPluginFailed
	}

	err = helpers.CopyDir(tempDir, pluginDir)
	if err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrInstallPluginFailed
	}

	if err := os.Chmod(destBinFile, 0755); err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrInstallPluginFailed
	}

	destChecksum, err := computeFileChecksum(destBinFile)
	if err != nil {
		p.logger.Error("Install plugin failed: compute installed binary checksum", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrInstallPluginFailed
	}
	if subtle.ConstantTimeCompare([]byte(binChecksum), []byte(destChecksum)) != 1 {
		p.logger.Error(
			"Install plugin failed: installed binary checksum mismatch",
			"plugin", pluginPath,
			"expected", binChecksum,
			"actual", destChecksum,
		)
		return nil, custom_errors.ErrPluginChecksumMismatch
	}

	newPlugin := plugin.Plugin{
		PackageName:     pluginMetadata.PackageName,
		Name:            pluginMetadata.Name,
		Author:          pluginMetadata.Author,
		Features:        plugin.ParseFeatureString(pluginMetadata.GetFeatures()),
		Version:         pluginMetadata.Version,
		Bin:             destBinFile,
		BinChecksum:     binChecksum,
		SignatureStatus: signatureStatus,
		IsActive:        true,
	}

	pluginEntity, err := p.pluginRepo.CreateOrUpdate(ctx, p.pluginMapper.ToEntity(&newPlugin))
	if err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrInstallPluginFailed
	}

	err = p.StartPlugin(&newPlugin)
	if err != nil {
		p.logger.Error("Install plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrStartPluginFailed
	}

	return p.pluginMapper.ToDomain(pluginEntity), nil

}

// Update replaces the plugin files and metadata in place while preserving the
// existing plugin ID. This keeps linked sessions (owner_id = plugin ID) intact
// instead of orphaning them like uninstall + reinstall does.
func (p *PluginManager) Update(ctx context.Context, pluginPath string) (*plugin.Plugin, error) {
	return p.update(ctx, pluginPath, "")
}

// UpdateWithSignature updates an archive with a detached signature that is
// not stored next to the archive.
func (p *PluginManager) UpdateWithSignature(ctx context.Context, pluginPath, signaturePath string) (*plugin.Plugin, error) {
	return p.update(ctx, pluginPath, signaturePath)
}

func (p *PluginManager) update(ctx context.Context, pluginPath, signaturePath string) (*plugin.Plugin, error) {
	_, err := os.Stat(pluginPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, custom_errors.ErrPluginNotFound
		}
		p.logger.Error("Update plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrPluginUpdateFailed
	}

	signatureStatus, err := p.verifyArchiveSignature(pluginPath, signaturePath)
	if err != nil {
		p.logger.Error("Update plugin failed: archive signature verification failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrPluginSignatureInvalid
	}

	id := uuid.New().String()
	tempDir := filepath.Join(paths.TempDir, id)
	metadataFile := filepath.Join(tempDir, "metadata.json")
	err = os.MkdirAll(tempDir, 0775)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	err = unpackPlugin(pluginPath, tempDir)
	if err != nil {
		p.logger.Error("Update plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrUnpackPluginFailed
	}

	pluginMetadata, err := loadMetadata(metadataFile)
	if err != nil {
		p.logger.Error("Update plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrLoadPluginMetadataFailed
	}

	err = pluginMetadata.Validate()
	if err != nil {
		p.logger.Error("Update plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrPluginMetadataInvalid
	}

	binFile, isExist := pluginMetadata.Bin[runtime.GOOS]
	if !isExist {
		return nil, custom_errors.ErrPluginBinaryMissing
	}

	cleanRelBin := filepath.Clean(binFile)
	normalizedBin := strings.ReplaceAll(cleanRelBin, "\\", "/")
	if filepath.IsAbs(cleanRelBin) || strings.HasPrefix(normalizedBin, "../") || normalizedBin == ".." || strings.Contains(normalizedBin, "/../") {
		p.logger.Error("Update plugin failed: binary path attempts traversal", "bin", binFile)
		return nil, custom_errors.ErrPluginMetadataInvalid
	}

	tempBinFile := filepath.Join(tempDir, cleanRelBin)
	tempSigFile := filepath.Join(tempDir, pluginMetadata.PackageName+".sig")
	if err := verifyFileSignature(tempBinFile, tempSigFile); err != nil {
		p.logger.Error("Update plugin failed: signature verification failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrPluginSignatureInvalid
	}

	binChecksum, err := computeFileChecksum(tempBinFile)
	if err != nil {
		p.logger.Error("Update plugin failed: compute binary checksum", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrPluginUpdateFailed
	}

	// Locate the installed plugin by package name so we can preserve its ID.
	existing, err := p.pluginRepo.FindByPackageName(ctx, pluginMetadata.PackageName)
	if err != nil {
		return nil, custom_errors.ErrPluginUpdateFailed
	}
	if existing == nil {
		return nil, custom_errors.ErrPluginNotFound
	}

	existingDomain := p.pluginMapper.ToDomain(existing)

	// Stop the running plugin before replacing its files.
	if err := p.StopPlugin(existingDomain); err != nil {
		p.logger.Error("Update plugin failed: stop plugin", "error", err, "packageName", pluginMetadata.PackageName)
		// Continue anyway so files can still be replaced.
	}

	pluginDir := filepath.Join(paths.PluginDir, pluginMetadata.PackageName)
	cleanPluginDir := filepath.Clean(pluginDir)

	// Remove old files then copy the new version into the same install location.
	if err := os.RemoveAll(cleanPluginDir); err != nil {
		p.logger.Error("Update plugin failed: remove old plugin dir", "error", err, "packageName", pluginMetadata.PackageName)
		return nil, custom_errors.ErrPluginUpdateFailed
	}
	if err := os.MkdirAll(cleanPluginDir, 0775); err != nil {
		return nil, custom_errors.ErrPluginUpdateFailed
	}
	if err := helpers.CopyDir(tempDir, cleanPluginDir); err != nil {
		p.logger.Error("Update plugin failed: copy files", "error", err, "packageName", pluginMetadata.PackageName)
		return nil, custom_errors.ErrPluginUpdateFailed
	}

	destBinFile := filepath.Clean(filepath.Join(cleanPluginDir, cleanRelBin))
	if err := os.Chmod(destBinFile, 0755); err != nil {
		p.logger.Error("Update plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrPluginUpdateFailed
	}

	destChecksum, err := computeFileChecksum(destBinFile)
	if err != nil {
		p.logger.Error("Update plugin failed: compute installed binary checksum", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrPluginUpdateFailed
	}
	if subtle.ConstantTimeCompare([]byte(binChecksum), []byte(destChecksum)) != 1 {
		p.logger.Error(
			"Update plugin failed: installed binary checksum mismatch",
			"plugin", pluginPath,
			"expected", binChecksum,
			"actual", destChecksum,
		)
		return nil, custom_errors.ErrPluginChecksumMismatch
	}

	// Update DB record while preserving the existing ID (sessions stay linked).
	existing.Name = pluginMetadata.Name
	existing.Author = pluginMetadata.Author
	existing.Features = pluginMetadata.GetFeatures()
	existing.Version = pluginMetadata.Version
	existing.Bin = destBinFile
	existing.BinChecksum = binChecksum
	existing.SignatureStatus = signatureStatus
	existing.IsActive = true

	if err := p.pluginRepo.Update(ctx, existing); err != nil {
		p.logger.Error("Update plugin failed", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrPluginUpdateFailed
	}

	updatedDomain := p.pluginMapper.ToDomain(existing)

	if err := p.StartPlugin(updatedDomain); err != nil {
		p.logger.Error("Update plugin failed: start plugin", "error", err, "plugin", pluginPath)
		return nil, custom_errors.ErrStartPluginFailed
	}

	return updatedDomain, nil
}

func (p *PluginManager) StartAllPlugin(ctx context.Context) error {
	activePlugins, err := p.pluginRepo.FindActive(ctx)
	if err != nil {
		p.logger.Error("Start plugin failed", "error", err)
		return custom_errors.ErrStartPluginFailed
	}

	for _, plg := range p.pluginMapper.ToDomains(activePlugins) {
		p.ensureChecksum(ctx, plg)
		err := p.StartPlugin(plg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PluginManager) SetPluginSocketPath(path string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.socketPath = path
}

func (p *PluginManager) GetPluginSocketPath() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.socketPath
}

func (p *PluginManager) SetConnectCode(packageName string, code string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.connectCodes[packageName] = code
}

func (p *PluginManager) StopAllPlugin(ctx context.Context) error {
	activePlugins, err := p.pluginRepo.FindActive(ctx)
	if err != nil {
		p.logger.Error("Stop plugin failed", "error", err)
		return custom_errors.ErrStartPluginFailed
	}

	for _, plg := range p.pluginMapper.ToDomains(activePlugins) {
		err := p.StopPlugin(plg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PluginManager) ValidConnect(cxt context.Context, packageName string, connectCode string) (*plugin.Plugin, error) {
	pluginEntity, err := p.pluginRepo.FindByPackageName(cxt, packageName)
	if err != nil {
		return nil, custom_errors.ErrCheckPluginConnectionFailed
	}

	if pluginEntity == nil {
		return nil, custom_errors.ErrPluginNotFound
	}

	p.mu.RLock()
	code, ok := p.connectCodes[packageName]
	p.mu.RUnlock()
	if !ok {
		return nil, custom_errors.ErrPluginConnectionInvalid
	}

	if subtle.ConstantTimeCompare([]byte(code), []byte(connectCode)) != 1 {
		return nil, custom_errors.ErrPluginConnectionInvalid
	}

	return p.pluginMapper.ToDomain(pluginEntity), nil
}

func (p *PluginManager) Uninstall(ctx context.Context, id string) error {
	plg, err := p.getPluginByID(ctx, id)
	if err != nil {
		if errors.Is(err, custom_errors.ErrPluginNotFound) {
			return err
		}

		return custom_errors.ErrUninstallPluginFailed
	}

	p.logger.Info("Uninstalling plugin", "packageName", plg.PackageName, "name", plg.Name, "version", plg.Version)

	p.logger.Info("Stopping plugin before uninstall", "packageName", plg.PackageName)
	err = p.StopPlugin(plg)
	if err != nil {
		p.logger.Logger.Error("Failed to stop plugin before uninstall", "error", err, "packageName", plg.PackageName)
	} else {
		p.logger.Info("Plugin stopped successfully", "packageName", plg.PackageName)
	}

	err = p.pluginRepo.DeleteById(ctx, plg.ID.String())
	if err != nil {
		return custom_errors.ErrUninstallPluginFailed
	}

	pluginDir := filepath.Join(paths.PluginDir, plg.PackageName)
	err = os.RemoveAll(pluginDir)
	if err != nil {
		p.logger.Error("Failed to remove plugin directory", "error", err, "packageName", plg.PackageName)
	}

	p.logger.Info("Uninstall plugin completed", "packageName", plg.PackageName, "name", plg.Name, "version", plg.Version)
	return nil
}

func (p *PluginManager) Activate(ctx context.Context, id string) error {
	plg, err := p.getPluginByID(ctx, id)
	if err != nil {
		if errors.Is(err, custom_errors.ErrPluginNotFound) {
			return err
		}

		return custom_errors.ErrActivatePluginFailed
	}

	p.logger.Info("Activating plugin", "packageName", plg.PackageName, "name", plg.Name, "version", plg.Version)

	if plg.IsActive {
		return custom_errors.ErrPluginAlreadyActive
	}

	p.ensureChecksum(ctx, plg)

	plg.IsActive = true
	err = p.pluginRepo.Update(ctx, p.pluginMapper.ToEntity(plg))
	if err != nil {
		return custom_errors.ErrActivatePluginFailed
	}

	p.logger.Info("Plugin activated successfully", "packageName", plg.PackageName, "name", plg.Name, "version", plg.Version)
	return p.StartPlugin(plg)
}

func (p *PluginManager) Deactivate(ctx context.Context, id string) error {
	plg, err := p.getPluginByID(ctx, id)
	if err != nil {
		if errors.Is(err, custom_errors.ErrPluginNotFound) {
			return err
		}

		return custom_errors.ErrDeactivatePluginFailed
	}

	p.logger.Info("Deactivating plugin", "packageName", plg.PackageName, "name", plg.Name, "version", plg.Version)
	if !plg.IsActive {
		return custom_errors.ErrPluginNotActive
	}

	err = p.StopPlugin(plg)
	if err != nil {
		return custom_errors.ErrDeactivatePluginFailed
	}

	plg.IsActive = false
	err = p.pluginRepo.Update(ctx, p.pluginMapper.ToEntity(plg))
	if err != nil {
		return custom_errors.ErrDeactivatePluginFailed
	}
	p.logger.Info("Plugin deactivated successfully", "packageName", plg.PackageName, "name", plg.Name, "version", plg.Version)
	return nil
}

func (p *PluginManager) GetPluginStatus(ctx context.Context, id string) (*plugin.Status, error) {
	pluginDomain, err := p.getPluginByID(ctx, id)
	if err != nil {
		if errors.Is(err, custom_errors.ErrPluginNotFound) {
			return nil, err
		}

		return nil, custom_errors.ErrGetPluginStatusFailed
	}

	pidPath := pluginDomain.Bin + ".pid"
	data, err := os.ReadFile(pidPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &plugin.Status{
				PID:         "",
				PackageName: pluginDomain.PackageName,
				Name:        pluginDomain.Name,
				Version:     pluginDomain.Version,
			}, nil
		}
		p.logger.Error("Failed to read PID file", "error", err, "packageName", pluginDomain.PackageName)
		return nil, custom_errors.ErrGetPluginStatusFailed
	}

	var pid string
	if _, err := fmt.Sscanf(string(data), "%s", &pid); err != nil {
		p.logger.Error("Invalid PID format in file", "error", err, "packageName", pluginDomain.PackageName)
		return nil, custom_errors.ErrGetPluginStatusFailed
	}

	pidValue, err := strconv.ParseInt(pid, 10, 32)
	if err != nil {
		p.logger.Error("Invalid PID format in file", "error", err, "packageName", pluginDomain.PackageName)
		return nil, custom_errors.ErrGetPluginStatusFailed
	}

	proc, err := process.NewProcess(int32(pidValue))
	if err != nil {
		p.logger.Error("Failed to get process", "error", err, "packageName", pluginDomain.PackageName)
		return nil, custom_errors.ErrGetPluginStatusFailed
	}
	cpuPercent, _ := proc.CPUPercent()
	memInfo, _ := proc.MemoryInfo()

	return &plugin.Status{
		PID:         pid,
		PackageName: pluginDomain.PackageName,
		Name:        pluginDomain.Name,
		Version:     pluginDomain.Version,
		CPUPercent:  cpuPercent,
		MemoryUsage: float64(memInfo.RSS) / (1024 * 1024),
	}, nil
}

// @Injectable
func NewPluginManager(
	pluginRepo *repositories.PluginRepository,
	pluginMapper *mappers.PluginMapper,
	logger *logger.BaseLogger,
) *PluginManager {
	return &PluginManager{
		pluginRepo:   pluginRepo,
		pluginMapper: pluginMapper,
		logger:       logger.With("module", "plugin-manager"),
		connectCodes: make(map[string]string),
		socketPath:   paths.PluginSocketPath,
	}
}
