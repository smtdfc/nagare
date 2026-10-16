use crate::client::PluginClient;
use crate::dtos::{
    PrepareChatSessionEventPayload, PrepareChatSessionFailedEventPayload,
    PrepareChatSessionSuccessEventPayload, ResetChatChannelEventPayload,
    ResetChatChannelFailedEventPayload, ResetChatChannelSuccessEventPayload,
    SendChatMessageEventPayload, SendChatMessageFailedEventPayload,
    SendChatMessageSuccessEventPayload, PREPARE_CHAT_SESSION_EVENT,
    PREPARE_CHAT_SESSION_FAILED_EVENT, PREPARE_CHAT_SESSION_SUCCESS_EVENT,
    RESET_CHAT_CHANNEL_EVENT, RESET_CHAT_CHANNEL_FAILED_EVENT,
    RESET_CHAT_CHANNEL_SUCCESS_EVENT, SEND_CHAT_MESSAGE_EVENT,
    SEND_CHAT_MESSAGE_FAILED_EVENT, SEND_CHAT_MESSAGE_SUCCESS_EVENT,
};
use crate::errors::PluginError;

impl PluginClient {
    /// Mirror of Go `PrepareChatSession(ctx, channelID)`.
    pub async fn prepare_chat_session(
        &self,
        channel_id: &str,
    ) -> Result<String, PluginError> {
        let payload = PrepareChatSessionEventPayload {
            channel_id: channel_id.to_string(),
        };
        let resp = self
            .send_and_wait::<
                PrepareChatSessionSuccessEventPayload,
                PrepareChatSessionFailedEventPayload,
            >(
                PREPARE_CHAT_SESSION_EVENT,
                &payload,
                PREPARE_CHAT_SESSION_FAILED_EVENT,
                PREPARE_CHAT_SESSION_SUCCESS_EVENT,
            )
            .await
            .map_err(|e| {
                tracing::error!(error = %e, "failed to prepare chat session");
                e
            })?;

        if !resp.is_success {
            let cause = resp
                .error
                .map(|e| e.cause)
                .unwrap_or_else(|| "unknown".to_string());
            tracing::error!(cause = %cause, "failed to prepare chat session");
            return Err(PluginError::PrepareChatSessionFailed(cause));
        }
        Ok(resp.data.map(|d| d.session_id).unwrap_or_default())
    }

    /// Mirror of Go `SendChatMessage(ctx, sessionID, text)`.
    pub async fn send_chat_message(
        &self,
        session_id: &str,
        text: &str,
    ) -> Result<(), PluginError> {
        let payload = SendChatMessageEventPayload {
            session_id: session_id.to_string(),
            text: text.to_string(),
        };
        let resp = self
            .send_and_wait::<
                SendChatMessageSuccessEventPayload,
                SendChatMessageFailedEventPayload,
            >(
                SEND_CHAT_MESSAGE_EVENT,
                &payload,
                SEND_CHAT_MESSAGE_FAILED_EVENT,
                SEND_CHAT_MESSAGE_SUCCESS_EVENT,
            )
            .await
            .map_err(|e| {
                tracing::error!(error = %e, "failed to send chat session messages");
                e
            })?;

        if !resp.is_success {
            let cause = resp
                .error
                .map(|e| e.cause)
                .unwrap_or_else(|| "unknown".to_string());
            tracing::error!(cause = %cause, "failed to send chat session messages");
            return Err(PluginError::Other(cause));
        }
        Ok(())
    }

    /// Mirror of Go `ResetChatChannel(ctx, channelID)`.
    pub async fn reset_chat_channel(&self, channel_id: &str) -> Result<(), PluginError> {
        let payload = ResetChatChannelEventPayload {
            channel_id: channel_id.to_string(),
        };
        let resp = self
            .send_and_wait::<
                ResetChatChannelSuccessEventPayload,
                ResetChatChannelFailedEventPayload,
            >(
                RESET_CHAT_CHANNEL_EVENT,
                &payload,
                RESET_CHAT_CHANNEL_FAILED_EVENT,
                RESET_CHAT_CHANNEL_SUCCESS_EVENT,
            )
            .await
            .map_err(|e| {
                tracing::error!(error = %e, "failed to reset chat channel");
                e
            })?;

        if !resp.is_success {
            let cause = resp
                .error
                .map(|e| e.cause)
                .unwrap_or_else(|| "unknown".to_string());
            tracing::error!(cause = %cause, "failed to reset chat channel");
            return Err(PluginError::Other(cause));
        }
        Ok(())
    }
}
