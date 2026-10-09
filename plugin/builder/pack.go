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

func prepareDirectories(cwd, packageName string) (pkgDir, binFile, sigFile string, err error) {
	pkgDir = filepath.Join(cwd, "pkg")
	binDir := filepath.Join(pkgDir, "bin")

	binFileName := packageName
	if runtime.GOOS == "windows" {
		binFileName += ".exe"
	}
	binFile = filepath.Join(binDir, binFileName)
	sigFile = filepath.Join(pkgDir, packageName+".sig")

	if err = os.MkdirAll(pkgDir, 0775); err != nil {
		return
	}
	if err = os.MkdirAll(binDir, 0775); err != nil {
		return
	}
	return
}

func generateFileSignature(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

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
	defer func(file *os.File) {
		_ = file.Close()
	}(outFile)

	zipWriter := zip.NewWriter(outFile)
	defer func(zipWriter *zip.Writer) {
		_ = zipWriter.Close()
	}(zipWriter)

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
		defer func(file *os.File) {
			_ = file.Close()
		}(file)

		_, err = io.Copy(fileWriter, file)
		return err
	})
}
