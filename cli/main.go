package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/smtdfc/nagare/cli/cmd"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file. Skipping")
	}

	cmd.Execute()
}
