use crate::client::PluginClient;
use crate::dtos::{
    HandshakeEventPayload, HandshakeFailedEventPayload, HandshakeSuccessEventPayload,
    HANDSHAKE_EVENT, HANDSHAKE_FAILED_EVENT, HANDSHAKE_SUCCESS_EVENT,
};
use crate::errors::PluginError;

impl PluginClient {
    /// Mirror of Go `Handshake(ctx)`.
    pub async fn handshake(&self) -> Result<(), PluginError> {
        let (package_name, connect_code) = {
            let metadata = self.metadata().unwrap_or_default();
            let cfg = self.connect_config();
            (metadata.package_name, cfg.connect_code)
        };
        let payload = HandshakeEventPayload {
            id: String::new(),
            package_name,
            connect_code,
        };
        let resp = self
            .send_and_wait::<HandshakeSuccessEventPayload, HandshakeFailedEventPayload>(
                HANDSHAKE_EVENT,
                &payload,
                HANDSHAKE_FAILED_EVENT,
                HANDSHAKE_SUCCESS_EVENT,
            )
            .await
            .map_err(|e| {
                tracing::error!(error = %e, "failed to perform handshake");
                e
            })?;

        if !resp.is_success {
            let cause = resp
                .error
                .map(|e| e.cause)
                .unwrap_or_else(|| "unknown".to_string());
            tracing::error!(cause = %cause, "handshake failed");
            return Err(PluginError::HandshakeFailed(cause));
        }
        tracing::info!("handshake successful");
        Ok(())
    }
}
