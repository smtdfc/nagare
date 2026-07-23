//go:build !dix
// +build !dix

package main

import (
	providers "github.com/smtdfc/nagare/core/providers"
	"github.com/smtdfc/nagare/server/.dix/generated"
)

func main() {
	err := providers.Init()
	if err != nil {
		panic(err)
	}

	generated.Root()
}
