package main

import (
	"strconv"
	"time"
)

const idleTimeout = 2 * time.Hour

func (tp *TelegramPlugin) getSessionID(chatID string) (string, bool) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	sessionID, ok := tp.channels[chatID]
	return sessionID, ok
}

func (tp *TelegramPlugin) setSessionID(chatID, sessionID string) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	tp.channels[chatID] = sessionID
}

func (tp *TelegramPlugin) getChatIDBySessionID(sessionID string) (string, int64, bool) {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	for chatID, sID := range tp.channels {
		if sID == sessionID {
			chatIntID, err := strconv.ParseInt(chatID, 10, 64)
			if err != nil {
				return chatID, 0, false
			}
			return chatID, chatIntID, true
		}
	}
	return "", 0, false
}

func (tp *TelegramPlugin) releaseSession(chatID string) {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	if sessionID, active := tp.channels[chatID]; active {
		delete(tp.channels, chatID)
		delete(tp.idleTimers, chatID)
		delete(tp.queues, chatID)
		delete(tp.isProcessing, chatID)
		for requestID, approval := range tp.approvals {
			if approval.chatID == chatID {
				delete(tp.approvals, requestID)
			}
		}
		tp.pluginClient.Logger.Info("release chat session due to 2 hours of inactivity.", "chatID", chatID, "sessionID", sessionID)
	}
}

func (tp *TelegramPlugin) resetIdleTimer(chatID string) {
	if t, ok := tp.idleTimers[chatID]; ok {
		t.Stop()
	}
	tp.idleTimers[chatID] = time.AfterFunc(idleTimeout, func() {
		tp.releaseSession(chatID)
	})
}

func (tp *TelegramPlugin) finishProcessing(chatID, sessionID string) {
	tp.mu.Lock()
	tp.isProcessing[chatID] = false
	tp.mu.Unlock()

	tp.clearBuffer(sessionID)
}

func (tp *TelegramPlugin) setPendingApproval(requestID, chatID, sessionID string, messageID int) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	if tp.approvals == nil {
		tp.approvals = make(map[string]pendingToolApproval)
	}
	tp.approvals[requestID] = pendingToolApproval{chatID: chatID, sessionID: sessionID, messageID: messageID}
}

func (tp *TelegramPlugin) getPendingApproval(requestID, chatID string) (pendingToolApproval, bool) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	approval, ok := tp.approvals[requestID]
	return approval, ok && approval.chatID == chatID
}

func (tp *TelegramPlugin) clearPendingApproval(requestID string) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	delete(tp.approvals, requestID)
}
