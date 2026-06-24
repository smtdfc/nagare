package plugin

import (
	"encoding/base64"
	"net"
)

func Connect(connectionStr string) (net.Conn, error) {
	addr, err := base64.StdEncoding.DecodeString(connectionStr)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", string(addr))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
