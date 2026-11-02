//go:build windows

package external_module

import "os"

func terminate(p *os.Process) error { return p.Kill() }
func alive(p *os.Process) bool      { return true }
