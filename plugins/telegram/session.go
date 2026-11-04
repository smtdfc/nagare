package main

import (
	"time"
)

const idleTimeout = 2 * time.Hour

func (tp *TelegramPlugin) hasChannel(channelID string) bool {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	_, ok := tp.channels[channelID]
	return ok
}

func (tp *TelegramPlugin) setChannel(channelID string) {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	tp.channels[channelID] = struct{}{}
}

func (tp *TelegramPlugin) clearChannel(channelID string) {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	delete(tp.channels, channelID)
}

func (tp *TelegramPlugin) releaseChannel(channelID string) {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	if _, active := tp.channels[channelID]; !active {
		return
	}

	delete(tp.channels, channelID)
	delete(tp.idleTimers, channelID)
	delete(tp.queues, channelID)
	delete(tp.isProcessing, channelID)
	for requestID, approval := range tp.approvals {
		if approval.channelID == channelID {
			delete(tp.approvals, requestID)
		}
	}
	tp.logger().Info("release chat channel due to 2 hours of inactivity.", "channelID", channelID)
}

func (tp *TelegramPlugin) resetIdleTimer(channelID string) {
	if t, ok := tp.idleTimers[channelID]; ok {
		t.Stop()
	}
	tp.idleTimers[channelID] = time.AfterFunc(idleTimeout, func() {
		tp.releaseChannel(channelID)
	})
}

func (tp *TelegramPlugin) finishProcessing(channelID string) {
	tp.mu.Lock()
	tp.isProcessing[channelID] = false
	tp.mu.Unlock()

	tp.clearBuffer(channelID)
}

func (tp *TelegramPlugin) setPendingApproval(requestID, channelID string, messageID int) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	if tp.approvals == nil {
		tp.approvals = make(map[string]pendingToolApproval)
	}
	tp.approvals[requestID] = pendingToolApproval{channelID: channelID, messageID: messageID}
}

func (tp *TelegramPlugin) getPendingApproval(requestID, channelID string) (pendingToolApproval, bool) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	approval, ok := tp.approvals[requestID]
	return approval, ok && approval.channelID == channelID
}

func (tp *TelegramPlugin) clearPendingApproval(requestID string) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	delete(tp.approvals, requestID)
}
