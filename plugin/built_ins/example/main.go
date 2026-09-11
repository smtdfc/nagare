package main

import "github.com/smtdfc/nagare/plugin/client"

func main() {
	pluginClient := client.NewPlugin()
	pluginClient.Start()
}
