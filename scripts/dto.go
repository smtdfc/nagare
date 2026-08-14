package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/coder/guts"
	"github.com/coder/guts/config"
)

func GenerateTS(pkg, outputPath string) error {
	golang, _ := guts.NewGolangParser()
	golang.PreserveComments()

	err := golang.IncludeGenerate(pkg)
	if err != nil {
		return err
	}

	ts, err := golang.ToTypescript()
	if err != nil {
		return err
	}

	ts.ApplyMutations(
		config.ExportTypes,
	)
	output, err := ts.Serialize()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, []byte(output), 0644)
}

func main() {
	var entries = map[string]string{
		"github.com/smtdfc/nagare/shared/dto":      "./ui/packages/dto/src/api.ts",
		"github.com/smtdfc/nagare/shared/messages": "./ui/packages/dto/src/messages.ts",
	}

	for entry, output := range entries {
		fmt.Printf("Generating %s -> %s ... ", entry, output)
		err := GenerateTS(entry, output)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
		fmt.Printf("OK \n")
	}
}
