package plugin

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"github.com/smtdfc/nagare/plugin-sdk/shared"
)

type Plugin struct {
	ID            string
	ConnectionStr string
	conn          net.Conn
}

func NewPlugin() *Plugin {
	cs := os.Getenv(shared.ADDR_ENV)
	id := os.Getenv(shared.PLUGIN_ID_ENV)
	if cs == "" {
		log.Fatal("Nagare not running ")
		return nil
	}

	return &Plugin{ID: id, ConnectionStr: cs}
}

func (p *Plugin) Connect() error {
	conn, err := Connect(p.ConnectionStr)
	if err != nil {
		return err
	}

	p.conn = conn
	return p.Send(shared.REGISTER_PLUGIN_REQUEST, p.ID)
}

func (p *Plugin) Send(method string, payload interface{}) error {
	data, _ := json.Marshal(payload)
	msg := shared.Message{
		PluginID: p.ID,
		Kind:     method,
		Payload:  data,
	}
	return json.NewEncoder(p.conn).Encode(msg)
}

func (p *Plugin) Listen(handler func(msg shared.Message)) {
	decoder := json.NewDecoder(p.conn)
	for {
		var msg shared.Message
		if err := decoder.Decode(&msg); err != nil {
			break
		}
		handler(msg)
	}
}
