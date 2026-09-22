package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/coder/guts"
	"github.com/coder/guts/config"
)

func main() {
	cmd, _ := os.Getwd()

	dtoMapping := map[string]string{
		filepath.Join(cmd, "dtos", "rest"):       filepath.Join(cmd, "client", "dtos", "src", "rest.generated.ts"),
		filepath.Join(cmd, "dtos", "websocket"):  filepath.Join(cmd, "client", "dtos", "src", "websocket.generated.ts"),
		filepath.Join(cmd, "shared", "messages"): filepath.Join(cmd, "client", "messages", "src", "messages.generated.ts"),
	}

	for goPath, tsFile := range dtoMapping {
		golang, err := guts.NewGolangParser()
		if err != nil {
			fmt.Printf("Failed to initialize parser for %s: %v\n", goPath, err)
			continue
		}

		_ = golang.IncludeCustom(map[string]string{
			"chan github.com/smtdfc/nagare/shared/messages.Message": "any",
		})
		golang.PreserveComments()

		if err := golang.IncludeGenerate(goPath); err != nil {
			fmt.Printf("Failed to include generate path %s: %v\n", goPath, err)
			continue
		}

		ts, err := golang.ToTypescript()
		if err != nil {
			fmt.Printf("Failed to convert TypeScript from %s: %v\n", goPath, err)
			continue
		}

		ts.ApplyMutations(
			config.ExportTypes,
			config.InterfaceToType,
		)

		output, err := ts.Serialize()
		if err != nil {
			fmt.Printf("Failed to serialize data for %s: %v\n", tsFile, err)
			continue
		}

		if filepath.Base(goPath) == "messages" {
			messageUnion := `

// Union type for Message
export type Message = 
    | AgentCompletedMessage
    | AgentErrorMessage
    | AgentStartedMessage
    | ReasoningMessage
    | ResponseCompletedMessage
    | ResponseFailedMessage
    | ResponseStartedMessage
    | TextMessage
    | ToolCallMessage
    | ToolResultMessage;
`
			output += messageUnion
		}

		if err := os.WriteFile(tsFile, []byte(output), 0644); err != nil {
			fmt.Printf("Failed to write file %s: %v\n", tsFile, err)
			continue
		}

		fmt.Printf("Successfully generated: %s -> %s\n", goPath, tsFile)
	}
}
