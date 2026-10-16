use thiserror::Error;

#[derive(Debug, Error)]
pub enum PluginError {
    #[error("connection failed: {0}")]
    ConnectionFailed(String),
    #[error("connection not found")]
    ConnectionNotFound,
    #[error("connection not ready")]
    ConnectionNotReady,
    #[error("create payload failed: {0}")]
    CreatePayloadFailed(String),
    #[error("handshake failed: {0}")]
    HandshakeFailed(String),
    #[error("handshake timeout")]
    HandshakeTimeout,
    #[error("handshake cancelled")]
    HandshakeCancelled,
    #[error("prepare chat session failed: {0}")]
    PrepareChatSessionFailed(String),
    #[error("incorrect tool args")]
    IncorrectToolArgs,
    #[error("marshal results failed: {0}")]
    MarshalToolResultFailed(String),
    #[error("timeout")]
    Timeout,
    #[error("action cancelled")]
    ActionCancelled,
    #[error("unexpected event received: {0}")]
    UnexpectedEvent(String),
    #[error("payload or data is nil")]
    NilPayload,
    #[error("io error: {0}")]
    Io(String),
    #[error("json error: {0}")]
    Json(String),
    #[error("channel closed")]
    ChannelClosed,
    #[error("{0}")]
    Other(String),
}

impl From<std::io::Error> for PluginError {
    fn from(e: std::io::Error) -> Self {
        Self::Io(e.to_string())
    }
}

impl From<serde_json::Error> for PluginError {
    fn from(e: serde_json::Error) -> Self {
        Self::Json(e.to_string())
    }
}

impl From<tokio::sync::oneshot::error::RecvError> for PluginError {
    fn from(_: tokio::sync::oneshot::error::RecvError) -> Self {
        Self::ChannelClosed
    }
}
