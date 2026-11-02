use thiserror::Error;

#[derive(Debug, Error)]
pub enum ModuleError {
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
    #[error("init logger error: {0}")]
    InitLogger(String),
    #[error("logger error: {0}")]
    Logger(String),
    #[error("ipc error: {0}")]
    IPC(String),
    #[error("channel closed")]
    ChannelClosed,
    #[error("{0}")]
    Other(String),

    #[error("timeout")]
    Timeout,
}

impl From<nagare_ipc::IpcError> for ModuleError {
    fn from(e: nagare_ipc::IpcError) -> Self {
        match e {
            nagare_ipc::IpcError::ConnectionFailed(msg) => Self::ConnectionFailed(msg),
            nagare_ipc::IpcError::ConnectionNotReady => Self::ConnectionNotReady,
            nagare_ipc::IpcError::Io(msg) => Self::Io(msg),
            nagare_ipc::IpcError::Json(msg) => Self::Json(msg),
            nagare_ipc::IpcError::MessageTooLarge => Self::IPC(e.to_string()),
        }
    }
}

impl From<std::io::Error> for ModuleError {
    fn from(e: std::io::Error) -> Self {
        Self::Io(e.to_string())
    }
}

impl From<serde_json::Error> for ModuleError {
    fn from(e: serde_json::Error) -> Self {
        Self::Json(e.to_string())
    }
}

impl From<fern::InitError> for ModuleError {
    fn from(e: fern::InitError) -> Self {
        Self::InitLogger(e.to_string())
    }
}

impl From<log::SetLoggerError> for ModuleError {
    fn from(e: log::SetLoggerError) -> Self {
        Self::Logger(e.to_string())
    }
}

impl From<tokio::sync::oneshot::error::RecvError> for ModuleError {
    fn from(_: tokio::sync::oneshot::error::RecvError) -> Self {
        Self::ChannelClosed
    }
}
