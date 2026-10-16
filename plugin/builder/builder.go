package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func Build() {
	fmt.Printf("%s=== Nagare Plugin Builder ===%s\n\n", ColorBold, ColorReset)

	cwd, err := os.Getwd()
	if err != nil {
		printError("%v", err)
		return
	}

	printStep("Checking metadata")
	pluginMetadata, err := loadMetadata(filepath.Join(cwd, "metadata.json"))
	if err != nil {
		printError("%v", err)
		return
	}
	printSuccess()

	printStep("Validating metadata")
	if err := validateMetadata(pluginMetadata); err != nil {
		printError("%v", err)
		return
	}
	printSuccess()

	printStep("Generating package structure")
	pkgDir, binFile, sigFile, err := prepareDirectories(cwd, pluginMetadata.PackageName)
	if err != nil {
		printError("%v", err)
		return
	}
	printSuccess()

	fmt.Printf("%s[~]%s Detecting source language ... ", ColorYellow, ColorReset)
	lang, err := detectSourceLanguage(cwd)
	if err != nil {
		printError("%v", err)
		return
	}
	fmt.Printf("(%s)\n", lang)

	fmt.Printf("%s[~]%s Building %s source for target [%s/%s]...\n", ColorYellow, ColorReset, lang, runtime.GOOS, runtime.GOARCH)
	if err := compileSource(cwd, binFile, lang); err != nil {
		printError("Compilation fatal abort: %v", err)
		return
	}
	fmt.Printf("%s[+]%s Build completed successfully!\n", ColorGreen, ColorReset)

	printStep("Signing binary (SHA-256)")
	if err := signBinary(binFile, sigFile); err != nil {
		printError("%v", err)
		return
	}
	printSuccess()

	printStep("Generating metadata package")
	if err := savePackageMetadata(pkgDir, binFile, *pluginMetadata); err != nil {
		printError("%v", err)
		return
	}
	printSuccess()

	outputPackage := pluginMetadata.PackageName + ".nagare_plugin"
	printStep(fmt.Sprintf("Packing into archive (%s)", outputPackage))
	if err := packPlugin(pkgDir, outputPackage); err != nil {
		printError("%v", err)
		return
	}
	printSuccess()

	fmt.Printf("\n%s✨ All done! Plugin packed nicely at: %s%s\n", ColorGreen, outputPackage, ColorReset)
}
