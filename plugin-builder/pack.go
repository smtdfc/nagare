package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func createPluginPackage(binaries map[string]string, metadataPath, outputName string) error {
	archive, err := os.Create(outputName)
	if err != nil {
		return err
	}
	defer archive.Close()

	writer := zip.NewWriter(archive)
	defer writer.Close()
	files := []string{metadataPath}
	for _, output := range binaries {
		files = append(files, output)
	}

	for _, file := range files {
		fmt.Println("Compressing ", file)
		if err := addFileToZip(writer, file); err != nil {
			return fmt.Errorf("failed to add file %s to zip: %w", file, err)
		}
	}
	return nil
}

func addFileToZip(writer *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	if filepath.Base(filename) == "metadata.json" {
		header.Name = "metadata.json"
	} else {
		header.Name = filename
	}

	header.Method = zip.Deflate

	writerEntry, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writerEntry, fileToZip)
	return err
}
