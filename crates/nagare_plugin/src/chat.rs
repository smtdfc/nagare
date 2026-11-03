use crate::client::PluginClient;
use crate::dtos::{
    PREPARE_CHAT_CHANNEL_EVENT, PREPARE_CHAT_CHANNEL_FAILED_EVENT,
    PREPARE_CHAT_CHANNEL_SUCCESS_EVENT, PrepareChatChannelEventPayload,
    PrepareChatChannelFailedEventPayload, PrepareChatChannelSuccessEventPayload,
    RESET_CHAT_CHANNEL_EVENT, RESET_CHAT_CHANNEL_FAILED_EVENT, RESET_CHAT_CHANNEL_SUCCESS_EVENT,
    ResetChatChannelEventPayload, ResetChatChannelFailedEventPayload,
    ResetChatChannelSuccessEventPayload, SEND_CHAT_MESSAGE_EVENT, SEND_CHAT_MESSAGE_FAILED_EVENT,
    SEND_CHAT_MESSAGE_SUCCESS_EVENT, SendChatMessageEventPayload,
    SendChatMessageFailedEventPayload, SendChatMessageSuccessEventPayload,
};
use crate::errors::PluginError;

impl PluginClient {
    /// Mirror of Go `PrepareChatChannel(ctx, channelID)`.
    pub async fn prepare_chat_channel(&self, channel_id: &str) -> Result<(), PluginError> {
        let payload = PrepareChatChannelEventPayload {
            channel_id: channel_id.to_string(),
        };
        let resp = self
            .send_and_wait::<
                PrepareChatChannelSuccessEventPayload,
                PrepareChatChannelFailedEventPayload,
            >(
                PREPARE_CHAT_CHANNEL_EVENT,
                &payload,
                PREPARE_CHAT_CHANNEL_FAILED_EVENT,
                PREPARE_CHAT_CHANNEL_SUCCESS_EVENT,
            )
            .await
            .map_err(|e| {
                tracing::error!(error = %e, "failed to prepare chat channel");
                e
            })?;

        if !resp.is_success {
            let cause = resp
                .error
                .map(|e| e.cause)
                .unwrap_or_else(|| "unknown".to_string());
            tracing::error!(cause = %cause, "failed to prepare chat channel");
            return Err(PluginError::PrepareChatChannelFailed(cause));
        }
        Ok(())
    }

    /// Mirror of Go `SendChatMessage(ctx, channelID, text)`.
    pub async fn send_chat_message(&self, channel_id: &str, text: &str) -> Result<(), PluginError> {
        let payload = SendChatMessageEventPayload {
            channel_id: channel_id.to_string(),
            text: text.to_string(),
        };
        let resp = self
            .send_and_wait::<SendChatMessageSuccessEventPayload, SendChatMessageFailedEventPayload>(
                SEND_CHAT_MESSAGE_EVENT,
                &payload,
                SEND_CHAT_MESSAGE_FAILED_EVENT,
                SEND_CHAT_MESSAGE_SUCCESS_EVENT,
            )
            .await
            .map_err(|e| {
                tracing::error!(error = %e, "failed to send chat message");
                e
            })?;

        if !resp.is_success {
            let cause = resp
                .error
                .map(|e| e.cause)
                .unwrap_or_else(|| "unknown".to_string());
            tracing::error!(cause = %cause, "failed to send chat message");
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
