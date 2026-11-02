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

	h.mu.RLock()
	if _, ok := h.connections[data.ModuleID]; ok {
		h.mu.RUnlock()
		_ = pConn.Send(HandshakeFailedEvent, &HandshakeFailedEventPayload{
			Cause: custom_errors.ErrModuleConnectionBlocked.Error(),
		}, payload.RequestID)
		return
	}
	h.mu.RUnlock()

	isValid := h.externalModuleMgr.IsValidConnect(data.ModuleID, data.ConnectCode)
	if !isValid {
		_ = pConn.Send(HandshakeFailedEvent, &HandshakeFailedEventPayload{
			Cause: custom_errors.ErrInvalidModuleConnection.Error(),
		}, payload.RequestID)
		return
	}

	h.externalModuleMgr.RevokeConnectCode(data.ConnectCode)
	h.mu.Lock()
	h.connections[data.ModuleID] = pConn
	h.mu.Unlock()
	pConn.moduleID = data.ModuleID
	_ = pConn.Send(HandshakeSuccessEvent, &HandshakeSuccessEventPayload{}, payload.RequestID)
}
