//go:build !dix
// +build !dix

package main

import (
	"github.com/smtdfc/nagare/server/.dix/generated"
)

func main() {
	generated.Root()
}
