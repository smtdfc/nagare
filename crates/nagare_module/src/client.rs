use std::collections::HashMap;
use std::sync::{Arc, Mutex, RwLock};
use std::time::Duration;

use serde::{Serialize, de::DeserializeOwned};
use tokio::sync::oneshot;
use uuid::Uuid;

use nagare_ipc::Connector;

use crate::config::ConnectConfig;
use crate::dtos::Payload;

use crate::errors::ModuleError;
use crate::helpers::{default_socket_path, get_data};
use crate::logger::init_logging_once;

/// Result of a request/response round trip (mirror of Go `contractsActionResult`).
#[derive(Debug)]
pub struct ActionResult<S, F> {
    pub is_success: bool,
    pub data: Option<S>,
    pub error: Option<F>,
}

pub type OnOtherMessage = Arc<dyn Fn(String, Payload<serde_json::Value>) + Send + Sync>;

pub(crate) struct ModuleInner {
    pub(crate) connect_config: RwLock<ConnectConfig>,
    pub(crate) config_dir: RwLock<String>,
    pub(crate) data_dir: RwLock<String>,
    pub(crate) pending_requests:
        Mutex<HashMap<String, oneshot::Sender<Payload<serde_json::Value>>>>,
    pub(crate) on_message_callback: RwLock<OnOtherMessage>,
    pub(crate) connector: Connector<Payload<serde_json::Value>>,
}

/// Client for the external-module protocol.
#[derive(Clone)]
pub struct ModuleClient {
    pub(crate) inner: Arc<ModuleInner>,
}

impl ModuleClient {
    /// Mirror of Go `NewModule()`.
    pub fn new() -> Arc<Self> {
        Self::new_with_logger()
    }

    pub fn new_with_logger() -> Arc<Self> {
        let log_file = std::env::var("NAGARE_MODULE_LOG_FILE")
            .expect("NAGARE_MODULE_LOG_FILE environment variable is not set");

        init_logging_once(&log_file);

        let connector = Connector::new(None);
        let this = Arc::new(Self {
            inner: Arc::new(ModuleInner {
                connect_config: RwLock::new(ConnectConfig {
                    module_id: std::env::var("NAGARE_MODULE_ID").unwrap_or_default(),
                    socket_path: std::env::var("NAGARE_MODULE_SOCKET_PATH")
                        .expect("NAGARE_MODULE_SOCKET_PATH environment variable is not set"),
                    connect_code: std::env::var("NAGARE_MODULE_CONNECT_CODE")
                        .expect("NAGARE_MODULE_CONNECT_CODE environment variable is not set"),
                }),
                config_dir: RwLock::new(
                    std::env::var("NAGARE_MODULE_CONFIG_DIR")
                        .expect("NAGARE_MODULE_CONFIG_DIR environment variable is not set"),
                ),
                data_dir: RwLock::new(
                    std::env::var("NAGARE_MODULE_DATA_DIR")
                        .expect("NAGARE_MODULE_DATA_DIR environment variable is not set"),
                ),
                pending_requests: Mutex::new(HashMap::new()),
                on_message_callback: RwLock::new(Arc::new(
                    |_: String, _: Payload<serde_json::Value>| {},
                )),
                connector,
            }),
        });

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

    pub fn connect_config(&self) -> ConnectConfig {
        self.inner.connect_config.read().unwrap().clone()
    }

    pub fn set_connect_config(&self, cfg: ConnectConfig) {
        *self.inner.connect_config.write().unwrap() = cfg;
    }

    pub fn set_socket_path(&self, path: impl Into<String>) {
        self.inner.connect_config.write().unwrap().socket_path = path.into();
    }

    pub fn config_dir(&self) -> String {
        self.inner.config_dir.read().unwrap().clone()
    }

    pub fn data_dir(&self) -> String {
        self.inner.data_dir.read().unwrap().clone()
    }

    pub fn on_message(
        &self,
        cb: impl Fn(String, Payload<serde_json::Value>) + Send + Sync + 'static,
    ) {
        *self.inner.on_message_callback.write().unwrap() = Arc::new(cb);
    }

    fn resolve_socket_path(&self) -> String {
        let mut cfg = self.inner.connect_config.write().unwrap();
        if cfg.socket_path.is_empty() {
            cfg.socket_path = std::env::var("NAGARE_MODULE_SOCKET_PATH")
                .expect("NAGARE_MODULE_SOCKET_PATH environment variable is not set");
        }
        if cfg.socket_path.is_empty() {
            cfg.socket_path = default_socket_path();
        }
        cfg.socket_path.clone()
    }

    /// Connect without blocking (split of Go `Start`).
    pub async fn connect(&self) -> Result<(), ModuleError> {
        let path = self.resolve_socket_path();
        self.inner.connector.connect(&path).await.map_err(|e| {
            tracing::error!(error = %e, "Start module error");
            ModuleError::from(e)
        })
    }

    /// Wait until the connection closes (mirror of Go `Connector.Wait`).
    pub async fn wait(&self) {
        self.inner.connector.wait().await;
    }

    /// Mirror of Go `Start(ctx, onStart)`: connect, run `on_start`, then wait.
    ///
    /// `on_start` is an async closure that typically performs the handshake.
    pub async fn start<F, Fut>(&self, on_start: F) -> Result<(), ModuleError>
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
        data: &impl Serialize,
        failed_event: &str,
        success_event: &str,
    ) -> Result<ActionResult<S, F>, ModuleError>
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

        let send_result = self
            .inner
            .connector
            .send(&Payload {
                event: event.to_string(),
                request_id: request_id.clone(),
                data: serde_json::to_value(data)?,
            })
            .await;
        if let Err(e) = send_result {
            self.inner
                .pending_requests
                .lock()
                .unwrap()
                .remove(&request_id);
            tracing::error!(event = %event, error = %e, "failed to send event");
            return Err(e.into());
        }

        // Cleanup on every exit path, mirrors Go
        // `defer delete(pendingRequests, requestID)`.
        let resp = match rx.await {
            Ok(resp) => resp,
            Err(e) => {
                self.inner
                    .pending_requests
                    .lock()
                    .unwrap()
                    .remove(&request_id);
                return Err(e.into());
            }
        };
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
        Err(ModuleError::UnexpectedEvent(resp.event))
    }

    /// Same as [`Self::send_and_wait`] with a timeout (Go `ctx` equivalent).
    pub async fn send_and_wait_timeout<S, F>(
        &self,
        event: &str,
        payload: &impl Serialize,
        failed_event: &str,
        success_event: &str,
        timeout: Duration,
    ) -> Result<ActionResult<S, F>, ModuleError>
    where
        S: DeserializeOwned,
        F: DeserializeOwned,
    {
        tokio::time::timeout(
            timeout,
            self.send_and_wait(event, payload, failed_event, success_event),
        )
        .await
        .map_err(|_| ModuleError::Timeout)?
    }

    /// Route responses back to their `send_and_wait` caller. Events that do
    /// not match a pending request are handed to the `on_other_message`
    /// callback, which is how a module receives host-initiated events.
    fn handle_event(self: Arc<Self>, payload: Payload<serde_json::Value>) {
        let pending = self
            .inner
            .pending_requests
            .lock()
            .unwrap()
            .remove(&payload.request_id);

        if let Some(tx) = pending {
            let _ = tx.send(payload);
            return;
        }

        let cb = self.inner.on_message_callback.read().unwrap().clone();
        cb(payload.event.clone(), payload);
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn new_module_dirs() {
        let p = ModuleClient::new();
        assert!(!p.config_dir().is_empty());
        assert!(!p.data_dir().is_empty());
    }

    #[test]
    fn connect_config_defaults() {
        let p = ModuleClient::new();
        let cfg = p.connect_config();
        assert_eq!(
            cfg.socket_path,
            std::env::var("NAGARE_MODULE_SOCKET_PATH").unwrap()
        );
        p.set_socket_path("/tmp/nagare-external-module.sock");
        assert_eq!(
            p.connect_config().socket_path,
            "/tmp/nagare-external-module.sock"
        );
    }

    #[tokio::test]
    async fn send_without_connection_fails() {
        let p = ModuleClient::new();
        p.set_socket_path("/non/existent/path/socket.sock");

        let res = p
            .send_and_wait::<serde_json::Value, serde_json::Value>(
                "module:test",
                &serde_json::json!({}),
                "module:test:failed",
                "module:test:success",
            )
            .await;

        assert!(res.is_err());
        assert!(p.inner.pending_requests.lock().unwrap().is_empty());
    }

    #[tokio::test]
    async fn handle_event_resolves_pending_request() {
        let p = ModuleClient::new();
        let (tx, rx) = oneshot::channel();
        p.inner
            .pending_requests
            .lock()
            .unwrap()
            .insert("req-1".to_string(), tx);

        Arc::clone(&p).handle_event(Payload {
            event: "module:test:success".to_string(),
            request_id: "req-1".to_string(),
            data: serde_json::json!({"ok": true}),
        });

        let resp = rx.await.expect("pending request should be resolved");
        assert_eq!(resp.event, "module:test:success");
        assert!(p.inner.pending_requests.lock().unwrap().is_empty());
    }

    #[tokio::test]
    async fn handle_event_dispatches_other_message() {
        let p = ModuleClient::new();
        let received = Arc::new(Mutex::new(None));
        let sink = Arc::clone(&received);

        p.on_message(move |event, payload| {
            *sink.lock().unwrap() = Some((event, payload));
        });

        Arc::clone(&p).handle_event(Payload {
            event: "module:vector:add".to_string(),
            request_id: "host-1".to_string(),
            data: serde_json::json!({"ids": [1, 2]}),
        });

        let captured = received
            .lock()
            .unwrap()
            .take()
            .expect("callback not called");
        assert_eq!(captured.0, "module:vector:add");
        assert_eq!(captured.1.request_id, "host-1");
    }
}
