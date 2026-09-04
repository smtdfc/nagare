//go:build !dix
// +build !dix

//go:generate dix wire
package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"

	"github.com/smtdfc/nagare/gateway/.dix/generated"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Gateway panicked: %v\n", r)
			os.Exit(1)
		}
	}()

	_, err := generated.Root()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Gateway exited with error: %v\n", err)
		os.Exit(1)
	}
}
