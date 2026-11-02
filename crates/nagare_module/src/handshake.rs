use crate::{ModuleClient, ModuleError};
use serde::{Deserialize, Serialize};

const HANDSHAKE_EVENT: &str = "external_module:handshake";
const HANDSHAKE_FAILED_EVENT: &str = "external_module:handshake:failed";
const HANDSHAKE_SUCCESS_EVENT: &str = "external_module:handshake:success";

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct HandshakeEventPayload {
    #[serde(rename = "moduleID")]
    pub module_id: String,
    #[serde(rename = "connectCode")]
    pub connect_code: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct HandshakeSuccessEventPayload {}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct HandshakeFailedEventPayload {
    pub cause: String,
}

impl ModuleClient {
    pub async fn handshake(&self) -> Result<(), ModuleError> {
        let connect_config = self.connect_config();
        let payload = HandshakeEventPayload {
            connect_code: connect_config.connect_code,
            module_id: connect_config.module_id,
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
            return Err(ModuleError::HandshakeFailed(cause));
        }
        tracing::info!("handshake successful");
        Ok(())
    }
}
