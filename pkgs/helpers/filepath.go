package helpers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func EnsurePathExist(p string) error {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err := os.MkdirAll(p, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func CopyDir(srcDir, destDir string) error {
	info, err := os.Stat(srcDir)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("source is not a directory: %s", srcDir)
	}

	if err := os.MkdirAll(destDir, info.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	destFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	return destFile.Sync()
}
