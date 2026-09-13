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
	pkgDir, binFile, sigFile, err := prepareDirectories(cwd, pluginMetadata.ID)
	if err != nil {
		printError("%v", err)
		return
	}
	printSuccess()

	fmt.Printf("%s[~]%s Building Go source for target [%s/%s]...\n", ColorYellow, ColorReset, runtime.GOOS, runtime.GOARCH)
	if err := compileSource(binFile); err != nil {
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

	outputPackage := pluginMetadata.ID + ".nagare_plugin"
	printStep(fmt.Sprintf("Packing into archive (%s)", outputPackage))
	if err := packPlugin(pkgDir, outputPackage); err != nil {
		printError("%v", err)
		return
	}
	printSuccess()

	fmt.Printf("\n%s✨ All done! Plugin packed nicely at: %s%s\n", ColorGreen, outputPackage, ColorReset)
}
