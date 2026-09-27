package client

import (
	"context"
	"errors"

	plugin_dtos "github.com/smtdfc/nagare/dtos/plugin"
)

func (p *PluginClient) Handshake(ctx context.Context) error {
	payload := plugin_dtos.HandshakeEventPayload{
		PackageName: p.Metadata.PackageName,
		ConnectCode: p.ConnectConfig.ConnectCode,
	}

	resp, err := sendAndWait[plugin_dtos.HandshakeSuccessEventPayload, plugin_dtos.HandshakeFailedEventPayload](
		p,
		ctx,
		plugin_dtos.HandshakeEvent,
		payload,
		plugin_dtos.HandshakeFailedEvent,
		plugin_dtos.HandshakeSuccessEvent,
	)
	if err != nil {
		p.Logger.Error("failed to perform handshake", "error", err)
		return err
	}

	if !resp.IsSuccess {
		p.Logger.Error("handshake failed", "error", resp.Error)
		return errors.New(resp.Error.Cause)
	}

	p.Logger.Info("handshake successful")

	return nil
}
