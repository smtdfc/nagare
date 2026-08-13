//go:build !dix
// +build !dix

package main

import (
	"github.com/smtdfc/nagare/core/global"
	"github.com/smtdfc/nagare/server/.dix/generated"
)

func main() {
	err := global.Init()
	if err != nil {
		panic(err)
	}

	generated.Root()
}
