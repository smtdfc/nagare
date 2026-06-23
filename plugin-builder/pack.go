package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func createPluginPackage(binaryPath, metadataPath, outputName string) error {
	archive, err := os.Create(outputName)
	if err != nil {
		return err
	}
	defer archive.Close()

	writer := zip.NewWriter(archive)
	defer writer.Close()

	files := []string{binaryPath, metadataPath}

	for _, file := range files {
		if err := addFileToZip(writer, file); err != nil {
			return err
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
	header.Name = filepath.Base(filename)
	header.Method = zip.Deflate

	writerEntry, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writerEntry, fileToZip)
	return err
}
