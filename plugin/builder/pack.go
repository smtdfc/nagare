package builder

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

func prepareDirectories(cwd, pluginID string) (pkgDir, binFile, sigFile string, err error) {
	pkgDir = filepath.Join(cwd, "pkg")
	binDir := filepath.Join(pkgDir, "bin")

	binFileName := pluginID
	if runtime.GOOS == "windows" {
		binFileName += ".exe"
	}
	binFile = filepath.Join(binDir, binFileName)
	sigFile = filepath.Join(pkgDir, pluginID+".sig")

	if err = os.MkdirAll(pkgDir, 0775); err != nil {
		return
	}
	if err = os.MkdirAll(binDir, 0775); err != nil {
		return
	}
	return
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

func generateFileSignature(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func signBinary(binFile, sigFile string) error {
	signature, err := generateFileSignature(binFile)
	if err != nil {
		return err
	}
	return os.WriteFile(sigFile, []byte(signature), 0664)
}

func packPlugin(pkgDir, outputFilePath string) error {
	outFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	return filepath.Walk(pkgDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == pkgDir {
			return nil
		}

		relPath, err := filepath.Rel(pkgDir, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)

		if info.IsDir() {
			_, err := zipWriter.Create(relPath + "/")
			return err
		}

		fileWriter, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(fileWriter, file)
		return err
	})
}
