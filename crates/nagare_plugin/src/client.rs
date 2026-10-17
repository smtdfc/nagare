use std::collections::HashMap;
use std::sync::{Arc, Mutex, RwLock};
use std::time::Duration;

use serde::{Serialize, de::DeserializeOwned};
use tokio::sync::oneshot;
use uuid::Uuid;

use crate::config::ConnectConfig;
use crate::connector::Connector;
use crate::dtos::{
    HANDSHAKE_FAILED_EVENT, HANDSHAKE_SUCCESS_EVENT, PLUGIN_TOOL_CALL_EVENT,
    PREPARE_CHAT_SESSION_FAILED_EVENT, PREPARE_CHAT_SESSION_SUCCESS_EVENT, Payload,
    PluginToolCallEventPayload, RECEIVED_CHAT_MESSAGE_EVENT, REGISTER_PLUGIN_TOOL_FAILED_EVENT,
    REGISTER_PLUGIN_TOOL_SUCCESS_EVENT, REGISTER_TOOL_CATEGORIES_FAILED_EVENT,
    REGISTER_TOOL_CATEGORIES_SUCCESS_EVENT, RESET_CHAT_CHANNEL_FAILED_EVENT,
    RESET_CHAT_CHANNEL_SUCCESS_EVENT, ReceivedChatMessageEventPayload,
    SEND_CHAT_MESSAGE_FAILED_EVENT, SEND_CHAT_MESSAGE_SUCCESS_EVENT,
};
use crate::errors::PluginError;
use crate::helpers::get_data;
use crate::ipc::default_socket_path;
use crate::logger::init_logging_once;
use crate::metadata::PluginMetadata;
use crate::tool::PluginTool;

/// Result of a request/response round trip (mirror of Go `plugin_dtosActionResult`).
#[derive(Debug)]
pub struct ActionResult<S, F> {
    pub is_success: bool,
    pub data: Option<S>,
    pub error: Option<F>,
}

/// Callback for incoming chat chunks.
pub type OnReceivedChatMessage = Arc<dyn Fn(String, String, String) + Send + Sync>;

pub type OnOtherMessage = Arc<dyn Fn(String, Payload<serde_json::Value>) + Send + Sync>;

pub(crate) struct PluginInner {
    pub(crate) metadata: RwLock<Option<PluginMetadata>>,
    pub(crate) connect_config: RwLock<ConnectConfig>,
    pub(crate) config_dir: RwLock<String>,
    pub(crate) pending_requests:
        Mutex<HashMap<String, oneshot::Sender<Payload<serde_json::Value>>>>,
    pub(crate) tools: Mutex<Vec<Arc<dyn PluginTool>>>,
    pub(crate) on_received_chat_message: RwLock<OnReceivedChatMessage>,
    pub(crate) on_other_message: RwLock<OnOtherMessage>,
    pub(crate) connector: Connector,
}

/// Mirror of Go `PluginClient`.
#[derive(Clone)]
pub struct PluginClient {
    pub(crate) inner: Arc<PluginInner>,
}

impl PluginClient {
    /// Mirror of Go `NewPlugin()`.
    pub fn new() -> Arc<Self> {
        Self::new_with_logger()
    }

    pub fn new_with_logger() -> Arc<Self> {
        let log_file = std::env::var("NAGARE_PLUGIN_LOG_FILE")
            .expect("NAGARE_PLUGIN_LOG_FILE environment variable is not set");

        init_logging_once(&log_file);

        let connector = Connector::new(None);
        let this = Arc::new(Self {
            inner: Arc::new(PluginInner {
                metadata: RwLock::new(None),
                connect_config: RwLock::new(ConnectConfig::default()),
                config_dir: RwLock::new(
                    std::env::var("NAGARE_PLUGIN_CONFIG_DIR")
                        .expect("NAGARE_PLUGIN_CONFIG_DIR environment variable is not set"),
                ),
                pending_requests: Mutex::new(HashMap::new()),
                tools: Mutex::new(Vec::new()),
                on_received_chat_message: RwLock::new(Arc::new(
                    |_: String, _: String, _: String| {},
                )),
                on_other_message: RwLock::new(Arc::new(
                    |_: String, _: Payload<serde_json::Value>| {},
                )),
                connector,
            }),
        });

        // Route connector events back into `handle_event`, mirroring
        // `NewConnector(p.handleEvent)` in Go.
        let weak = Arc::downgrade(&this);
        let handler: Arc<dyn Fn(Payload<serde_json::Value>) + Send + Sync> =
            Arc::new(move |payload| {
                if let Some(strong) = weak.upgrade() {
                    strong.handle_event(payload);
                }
            });
        this.inner.connector.set_on_event(handler);
        this
    }

    pub fn metadata(&self) -> Option<PluginMetadata> {
        self.inner.metadata.read().unwrap().clone()
    }

    pub fn connect_config(&self) -> ConnectConfig {
        self.inner.connect_config.read().unwrap().clone()
    }

    pub fn set_connect_config(&self, cfg: ConnectConfig) {
        *self.inner.connect_config.write().unwrap() = cfg;
    }

    pub fn set_socket_path(&self, path: impl Into<String>) {
        self.inner.connect_config.write().unwrap().socket_path = path.into();
    }

    pub fn set_connect_code(&self, code: impl Into<String>) {
        self.inner.connect_config.write().unwrap().connect_code = code.into();
    }

    pub fn config_dir(&self) -> String {
        self.inner.config_dir.read().unwrap().clone()
    }

    pub fn set_on_received_chat_message(
        &self,
        cb: impl Fn(String, String, String) + Send + Sync + 'static,
    ) {
        *self.inner.on_received_chat_message.write().unwrap() = Arc::new(cb);
    }

    pub fn set_on_other_message(
        &self,
        cb: impl Fn(String, Payload<serde_json::Value>) + Send + Sync + 'static,
    ) {
        *self.inner.on_other_message.write().unwrap() = Arc::new(cb);
    }

    /// Mirror of Go `LoadMetadata(raw)`.
    pub fn load_metadata(&self, raw: &str) -> Result<PluginMetadata, PluginError> {
        let metadata: PluginMetadata = serde_json::from_str(raw)?;
        *self.inner.metadata.write().unwrap() = Some(metadata.clone());
        Ok(metadata)
    }

    fn resolve_socket_path(&self) -> String {
        let mut cfg = self.inner.connect_config.write().unwrap();
        if cfg.socket_path.is_empty() {
            cfg.socket_path = std::env::var("NAGARE_PLUGIN_SOCKET_PATH")
                .expect("NAGARE_PLUGIN_SOCKET_PATH environment variable is not set");
        }
        if cfg.socket_path.is_empty() {
            cfg.socket_path = default_socket_path();
        }
        if cfg.connect_code.is_empty() {
            cfg.connect_code = std::env::var("NAGARE_PLUGIN_CONNECT_CODE")
                .expect("NAGARE_PLUGIN_CONNECT_CODE environment variable is not set");
        }
        cfg.socket_path.clone()
    }

    /// Connect without blocking (split of Go `Start`).
    pub async fn connect(&self) -> Result<(), PluginError> {
        let path = self.resolve_socket_path();
        self.inner.connector.connect(&path).await.map_err(|e| {
            tracing::error!(error = %e, "Start plugin error");
            e
        })
    }

    /// Wait until the connection closes (mirror of Go `Connector.Wait`).
    pub async fn wait(&self) {
        self.inner.connector.wait().await;
    }

    /// Mirror of Go `Start(ctx, onStart)`: connect, run `on_start`, then wait.
    ///
    /// `on_start` is an async closure that typically performs the handshake;
    /// it captures a clone of the client, e.g.:
    /// ```ignore
    /// let c2 = client.clone();
    /// client.start(move || async move { let _ = c2.handshake().await; }).await?;
    /// ```
    pub async fn start<F, Fut>(&self, on_start: F) -> Result<(), PluginError>
    where
        F: FnOnce() -> Fut,
        Fut: std::future::Future<Output = ()>,
    {
        self.connect().await?;
        on_start().await;
        self.wait().await;
        Ok(())
    }

    /// Generic request/response helper (mirror of Go `sendAndWait`).
    pub async fn send_and_wait<S, F>(
        &self,
        event: &str,
        payload: &impl Serialize,
        failed_event: &str,
        success_event: &str,
    ) -> Result<ActionResult<S, F>, PluginError>
    where
        S: DeserializeOwned,
        F: DeserializeOwned,
    {
        let request_id = Uuid::new_v4().to_string();
        let (tx, rx) = oneshot::channel();
        {
            self.inner
                .pending_requests
                .lock()
                .unwrap()
                .insert(request_id.clone(), tx);
        }

        let send_result = self.inner.connector.send(event, payload, &request_id).await;
        if let Err(e) = send_result {
            self.inner
                .pending_requests
                .lock()
                .unwrap()
                .remove(&request_id);
            tracing::error!(event = %event, error = %e, "failed to send event");
            return Err(e);
        }

        let resp = rx.await?;
        // Cleanup, mirrors Go `defer delete(pendingRequests, requestID)`.
        self.inner
            .pending_requests
            .lock()
            .unwrap()
            .remove(&request_id);

        let mut result = ActionResult {
            is_success: false,
            data: None,
            error: None,
        };
        if resp.event == failed_event {
            result.is_success = false;
            result.error = get_data::<F>(&resp).ok();
            return Ok(result);
        }
        if !success_event.is_empty() && resp.event == success_event {
            result.is_success = true;
            result.data = get_data::<S>(&resp).ok();
            return Ok(result);
        }
        tracing::error!(event = %resp.event, "unexpected event received");
        Err(PluginError::UnexpectedEvent(resp.event))
    }

    /// Same as [`Self::send_and_wait`] with a timeout (Go `ctx` equivalent).
    pub async fn send_and_wait_timeout<S, F>(
        &self,
        event: &str,
        payload: &impl Serialize,
        failed_event: &str,
        success_event: &str,
        timeout: Duration,
    ) -> Result<ActionResult<S, F>, PluginError>
    where
        S: DeserializeOwned,
        F: DeserializeOwned,
    {
        tokio::time::timeout(
            timeout,
            self.send_and_wait(event, payload, failed_event, success_event),
        )
        .await
        .map_err(|_| PluginError::Timeout)?
    }

    fn handle_event(self: Arc<Self>, payload: Payload<serde_json::Value>) {
        match payload.event.as_str() {
            HANDSHAKE_SUCCESS_EVENT
            | HANDSHAKE_FAILED_EVENT
            | PREPARE_CHAT_SESSION_FAILED_EVENT
            | PREPARE_CHAT_SESSION_SUCCESS_EVENT
            | REGISTER_PLUGIN_TOOL_SUCCESS_EVENT
            | REGISTER_PLUGIN_TOOL_FAILED_EVENT
            | SEND_CHAT_MESSAGE_SUCCESS_EVENT
            | SEND_CHAT_MESSAGE_FAILED_EVENT
            | RESET_CHAT_CHANNEL_SUCCESS_EVENT
            | RESET_CHAT_CHANNEL_FAILED_EVENT
            | REGISTER_TOOL_CATEGORIES_SUCCESS_EVENT
            | REGISTER_TOOL_CATEGORIES_FAILED_EVENT => {
                if let Some(tx) = self
                    .inner
                    .pending_requests
                    .lock()
                    .unwrap()
                    .remove(&payload.request_id)
                {
                    let _ = tx.send(payload);
                }
            }
            PLUGIN_TOOL_CALL_EVENT => {
                let zelf = Arc::clone(&self);
                let request_id = payload.request_id.clone();
                match get_data::<PluginToolCallEventPayload>(&payload) {
                    Err(e) => {
                        let zelf2 = zelf.clone();
                        tokio::spawn(async move {
                            zelf2
                                .send_plugin_tool_result(&request_id, "", &e.to_string())
                                .await;
                        });
                    }
                    Ok(call) => {
                        tokio::spawn(async move {
                            zelf.handle_plugin_tool_call(request_id, call).await;
                        });
                    }
                }
            }
            RECEIVED_CHAT_MESSAGE_EVENT => {
                match get_data::<ReceivedChatMessageEventPayload>(&payload) {
                    Err(e) => {
                        tracing::warn!(error = %e, "failed to parse ReceivedChatMessageEvent");
                    }
                    Ok(d) => {
                        let cb = self.inner.on_received_chat_message.read().unwrap().clone();
                        cb(d.session_id, d.channel_id, d.message);
                    }
                }
            }
            other => {
                let cb = self.inner.on_other_message.read().unwrap().clone();
                cb(other.to_owned(), payload);
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::dtos::RECEIVED_CHAT_MESSAGE_EVENT;
    use std::sync::{Arc, Mutex};

    #[test]
    fn new_plugin_load_metadata() {
        let p = PluginClient::new();
        assert!(p.metadata().is_none());
        let raw = r#"{"package_name":"nagare.plugin.test","name":"t","version":"1","api_version":"1","author":"a"}"#;
        let m = p.load_metadata(raw).unwrap();
        assert_eq!(m.package_name, "nagare.plugin.test");
        assert!(p.metadata().is_some());
        assert!(p.load_metadata("{invalid").is_err());
    }

    #[test]
    fn connect_config_defaults() {
        let p = PluginClient::new();
        let cfg = p.connect_config();
        assert!(cfg.socket_path.is_empty());
        p.set_socket_path("/tmp/nagare.sock");
        p.set_connect_code("secret-code-123");
        let cfg = p.connect_config();
        assert_eq!(cfg.socket_path, "/tmp/nagare.sock");
        assert_eq!(cfg.connect_code, "secret-code-123");
    }

    #[tokio::test]
    async fn operations_without_connection_fail() {
        let p = PluginClient::new();
        assert!(p.handshake().await.is_err());
        assert!(p.prepare_chat_session("chat-123").await.is_err());
        assert!(p.send_chat_message("session-1", "hello").await.is_err());
        assert!(p.reset_chat_channel("chat-123").await.is_err());
    }

    #[tokio::test]
    async fn start_with_invalid_socket_returns_error() {
        let p = PluginClient::new();
        p.set_socket_path("/non/existent/path/socket.sock");
        let res = tokio::time::timeout(std::time::Duration::from_secs(5), p.start(|| async {}))
            .await
            .expect("start should resolve");
        assert!(res.is_err());
    }

    #[test]
    fn handle_event_dispatches_received_chat_message() {
        let p = PluginClient::new();
        let received = Arc::new(Mutex::new(false));
        let r2 = received.clone();
        p.set_on_received_chat_message(move |s, c, m| {
            if s == "sess-1" && c == "chan-1" && m == "test chunk" {
                *r2.lock().unwrap() = true;
            }
        });
        p.clone().handle_event(Payload {
            event: RECEIVED_CHAT_MESSAGE_EVENT.to_string(),
            request_id: "r1".to_string(),
            data: serde_json::json!({
                "sessionID": "sess-1",
                "channelID": "chan-1",
                "message": "test chunk",
            }),
        });
        assert!(*received.lock().unwrap());
    }

    #[test]
    fn handle_event_dispatches_other_message() {
        let p = PluginClient::new();
        let received = Arc::new(Mutex::new(false));
        let r2 = received.clone();
        p.set_on_other_message(move |s, c| {
            if s == "other-message" {
                *r2.lock().unwrap() = true;
            }
        });

        p.clone().handle_event(Payload {
            event: "other-message".to_owned(),
            request_id: "r1".to_string(),
            data: serde_json::json!({
                "sessionID": "sess-1",
                "channelID": "chan-1",
                "message": "test chunk",
            }),
        });
        assert!(*received.lock().unwrap());
    }
}
