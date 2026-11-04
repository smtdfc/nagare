package sdk

import (
	"context"
	"errors"

	contracts "github.com/smtdfc/nagare/pkgs/plugin/contracts"
)

func (p *PluginClient) Handshake(ctx context.Context) error {
	payload := contracts.HandshakeEventPayload{
		PackageName: p.Metadata.PackageName,
		ConnectCode: p.ConnectConfig.ConnectCode,
	}

	resp, err := sendAndWait[contracts.HandshakeSuccessEventPayload, contracts.HandshakeFailedEventPayload](
		p,
		ctx,
		contracts.HandshakeEvent,
		payload,
		contracts.HandshakeFailedEvent,
		contracts.HandshakeSuccessEvent,
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
