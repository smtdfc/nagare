package main

import (
	"github.com/smtdfc/nagare/plugin/client"
)

func main() {
	client := client.NewPluginClient("example")
	client.Start()
}
