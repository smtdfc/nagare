use thiserror::Error;

/// Transport-level errors only. Handshake, logger and protocol errors belong to
/// the SDK that sits on top of this crate (`nagare_plugin` / `nagare_module`).
#[derive(Debug, Error)]
pub enum IpcError {
    #[error("connection failed: {0}")]
    ConnectionFailed(String),
    #[error("connection not ready")]
    ConnectionNotReady,
    #[error("ipc message exceeds maximum size")]
    MessageTooLarge,
    #[error("io error: {0}")]
    Io(String),
    #[error("json error: {0}")]
    Json(String),
}

impl From<std::io::Error> for IpcError {
    fn from(e: std::io::Error) -> Self {
        Self::Io(e.to_string())
    }
}

impl From<serde_json::Error> for IpcError {
    fn from(e: serde_json::Error) -> Self {
        Self::Json(e.to_string())
    }
}
