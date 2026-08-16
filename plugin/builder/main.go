package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/smtdfc/nagare/shared/plugin"
)

// Helper function to actually copy files (avoiding symlinks which break when packed)
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

// Helper function to compress a directory into a .nagare_plugin file (ZIP format)
func zipDirectory(sourceDir, zipFile string) error {
	zipF, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer zipF.Close()

	archive := zip.NewWriter(zipF)
	defer archive.Close()

	info, err := os.Stat(sourceDir)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = sourceDir
	} else {
		baseDir = filepath.Dir(sourceDir)
	}

	return filepath.Walk(sourceDir, func(pathItem string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create header for file inside zip
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Get relative path compared to source directory
		relPath, err := filepath.Rel(baseDir, pathItem)
		if err != nil {
			return err
		}

		// Skip root directory itself
		if relPath == "." {
			return nil
		}

		header.Name = filepath.ToSlash(relPath)

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(pathItem)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	metadataFile := path.Join(cwd, "metadata.json")
	metadataContent, err := os.ReadFile(metadataFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to read metadata.json: %v", err))
	}

	var metadata plugin.PluginMetadata
	if err := json.Unmarshal(metadataContent, &metadata); err != nil {
		fmt.Println("Failed to unmarshal metadata:", err)
		panic(err)
	}

	buildDir := path.Join(cwd, "build")
	os.RemoveAll(buildDir)

	if err := os.MkdirAll(buildDir, 0755); err != nil {
		panic(err)
	}

	binsDir := path.Join(buildDir, "bins")
	if err := os.MkdirAll(binsDir, 0755); err != nil {
		panic(err)
	}

	var bins []plugin.PluginBin
	currentOS := "linux"

	for _, arch := range metadata.Architectures {
		fmt.Println("Building for architecture:", arch)

		binFileName := fmt.Sprintf("%s_%s", metadata.ID, arch)
		binFullPath := path.Join(binsDir, binFileName)
		binRelPath := path.Join("bins", binFileName)

		cmd := exec.Command("go", "build", "-v", "-x", "-o", binFullPath, ".")
		cmd.Env = append(os.Environ(), "GOOS="+currentOS, "GOARCH="+arch)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			panic(fmt.Sprintf("Failed to build for architecture %s: %v", arch, err))
		}

		bins = append(bins, plugin.PluginBin{
			Architecture: arch,
			Path:         binRelPath,
		})
	}

	metadata.Bins = bins

	pkgDir := path.Join(buildDir, "pkg")
	pkgBinsDir := path.Join(pkgDir, "bins")
	if err := os.MkdirAll(pkgBinsDir, 0755); err != nil {
		panic(err)
	}

	for _, bin := range bins {
		srcFile := path.Join(cwd, "build", bin.Path)
		dstFile := path.Join(pkgDir, bin.Path)
		if err := copyFile(srcFile, dstFile); err != nil {
			panic(fmt.Sprintf("Failed to copy binary: %v", err))
		}
	}

	updatedMetadataJSON, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		panic(err)
	}

	pkgMetadataPath := path.Join(pkgDir, "metadata.json")
	if err := os.WriteFile(pkgMetadataPath, updatedMetadataJSON, 0644); err != nil {
		panic(err)
	}

	outputPluginFile := path.Join(cwd, fmt.Sprintf("%s.nagare_plugin", metadata.ID))
	fmt.Println("Packaging plugin into:", outputPluginFile)

	if err := zipDirectory(pkgDir, outputPluginFile); err != nil {
		panic(fmt.Sprintf("Failed to compress plugin file: %v", err))
	}

	fmt.Println("Plugin successfully packaged!")
}
