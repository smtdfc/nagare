package main

import (
	"github.com/smtdfc/nagare/plugin-sdk/plugin"
	"github.com/smtdfc/nagare/plugin-sdk/shared"
)

func main() {
	plg := plugin.NewPlugin()

	plg.Connect()
	plg.Listen(func(msg shared.Message) {
		switch msg.Kind {
		case shared.REGISTER_PLUGIN_SUCCESS:

		}
	})
}
