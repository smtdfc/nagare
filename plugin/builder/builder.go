package builder

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func Build(privateKeyPath string) {
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
	signaturePath := outputPackage + ".sig"

	printStep(fmt.Sprintf("Packing into archive (%s)", outputPackage))
	if err := packPlugin(pkgDir, outputPackage); err != nil {
		printError("%v", err)
		return
	}
	printSuccess()

	// Drop any signature left over from a previous build so an unsigned
	// archive is never shipped next to a stale signature.
	if err := os.Remove(signaturePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		printError("%v", err)
		return
	}

	if strings.TrimSpace(privateKeyPath) == "" {
		fmt.Printf("%s[!]%s Archive is not signed. Pass --sign <private-key> to sign it before distributing.\n", ColorYellow, ColorReset)
	} else {
		printStep("Signing archive (Ed25519)")
		if err := signPluginArchive(outputPackage, privateKeyPath); err != nil {
			printError("%v", err)
			_ = os.Remove(outputPackage)
			return
		}
		printSuccess()
		fmt.Printf("%s[+]%s Signature written to %s\n", ColorGreen, ColorReset, signaturePath)
	}

	fmt.Printf("\n%s✨ All done! Plugin packed nicely at: %s%s\n", ColorGreen, outputPackage, ColorReset)
}
