package cmd

import (
	"github.com/spf13/cobra"
)

var port int

var rootCmd = &cobra.Command{
	Use: "nagare",
}

func init() {

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
