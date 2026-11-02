package external_module

import "github.com/smtdfc/nagare/core/custom_errors"

var (
	HandshakeEvent        Event = "external_module:handshake"
	HandshakeSuccessEvent Event = "external_module:handshake:success"
	HandshakeFailedEvent  Event = "external_module:handshake:failed"
)

type HandshakeEventPayload struct {
	ModuleID    string `json:"moduleID"`
	ConnectCode string `json:"connectCode"`
}

type HandshakeFailedEventPayload struct {
	Cause string `json:"cause"`
}

type HandshakeSuccessEventPayload struct{}

func (h *ExternalModuleHost) onHandshake(pConn *Connection, payload *Payload[any]) {
	data, err := getPayloadData[HandshakeEventPayload](payload)
	if err != nil {
		_ = pConn.Send(HandshakeFailedEvent, &HandshakeFailedEventPayload{
			Cause: custom_errors.ErrInvalidPayload.Error(),
		}, payload.RequestID)
		return
	}

	h.mu.Lock()
	if _, ok := h.connections[data.ModuleID]; ok {
		h.mu.Unlock()
		_ = pConn.Send(HandshakeFailedEvent, &HandshakeFailedEventPayload{
			Cause: custom_errors.ErrModuleConnectionBlocked.Error(),
		}, payload.RequestID)
		return
	}

	isValid := h.externalModuleMgr.IsValidConnect(data.ModuleID, data.ConnectCode)
	if !isValid {
		h.mu.Unlock()
		_ = pConn.Send(HandshakeFailedEvent, &HandshakeFailedEventPayload{
			Cause: custom_errors.ErrInvalidModuleConnection.Error(),
		}, payload.RequestID)
		return
	}

	h.externalModuleMgr.RevokeConnectCode(data.ConnectCode)
	h.connections[data.ModuleID] = pConn
	pConn.moduleID = data.ModuleID
	h.mu.Unlock()

	if err := pConn.Send(HandshakeSuccessEvent, &HandshakeSuccessEventPayload{}, payload.RequestID); err != nil {
		// Rollback
		h.unregister(pConn)
	}
}
