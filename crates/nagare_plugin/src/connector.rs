use std::sync::{
    atomic::{AtomicBool, Ordering},
    Arc, RwLock,
};

use serde::Serialize;
use tokio::sync::{Mutex, Notify};

use crate::dtos::Payload;
use crate::errors::PluginError;
use crate::ipc::{encode_message, MAX_MESSAGE_SIZE};

/// Callback invoked for every incoming payload.
pub type OnEvent = Arc<dyn Fn(Payload<serde_json::Value>) + Send + Sync>;

/// Mirror of Go `Connector`: length-prefixed framed transport + listen loop.
#[derive(Clone)]
pub struct Connector {
    inner: Arc<ConnectorInner>,
}

struct ConnectorInner {
    #[cfg(not(windows))]
    write_half: Mutex<Option<tokio::net::unix::OwnedWriteHalf>>,
    #[cfg(windows)]
    write_half: Mutex<Option<tokio::net::windows::named_pipe::NamedPipeClient>>,
    on_event: RwLock<Option<OnEvent>>,
    closed: AtomicBool,
    notify_closed: Notify,
}

impl Connector {
    pub fn new(on_event: Option<OnEvent>) -> Self {
        Self {
            inner: Arc::new(ConnectorInner {
                #[cfg(not(windows))]
                write_half: Mutex::new(None),
                #[cfg(windows)]
                write_half: Mutex::new(None),
                on_event: RwLock::new(on_event),
                closed: AtomicBool::new(false),
                notify_closed: Notify::new(),
            }),
        }
    }

    pub fn set_on_event(&self, on_event: OnEvent) {
        *self.inner.on_event.write().unwrap() = Some(on_event);
    }

    fn mark_closed(&self) {
        self.inner.closed.store(true, Ordering::SeqCst);
        self.inner.notify_closed.notify_waiters();
    }

    pub fn is_closed(&self) -> bool {
        self.inner.closed.load(Ordering::SeqCst)
    }

    fn dispatch(&self, payload: Payload<serde_json::Value>) {
        let handler = self.inner.on_event.read().unwrap().clone();
        if let Some(handler) = handler {
            handler(payload);
        }
    }

    /// Mirror of Go `Connect(ctx, path)`.
    pub async fn connect(&self, path: &str) -> Result<(), PluginError> {
        #[cfg(not(windows))]
        {
            let stream = tokio::net::UnixStream::connect(path)
                .await
                .map_err(|e| PluginError::ConnectionFailed(e.to_string()))?;
            let (read_half, write_half) = stream.into_split();
            *self.inner.write_half.lock().await = Some(write_half);
            self.inner.closed.store(false, Ordering::SeqCst);
            let zelf = self.clone();
            tokio::spawn(async move { zelf.read_loop_unix(read_half).await });
            Ok(())
        }
        #[cfg(windows)]
        {
            let client = tokio::net::windows::named_pipe::ClientOptions::new()
                .open(path)
                .map_err(|e| PluginError::ConnectionFailed(e.to_string()))?;
            *self.inner.write_half.lock().await = Some(client);
            self.inner.closed.store(false, Ordering::SeqCst);
            Ok(())
        }
    }

    #[cfg(not(windows))]
    async fn read_loop_unix(&self, mut read_half: tokio::net::unix::OwnedReadHalf) {
        use tokio::io::AsyncReadExt;
        loop {
            let mut len_buf = [0u8; 4];
            if read_half.read_exact(&mut len_buf).await.is_err() {
                break;
            }
            let msg_len = u32::from_be_bytes(len_buf) as usize;
            if msg_len > MAX_MESSAGE_SIZE {
                break;
            }
            let mut buf = vec![0u8; msg_len];
            if read_half.read_exact(&mut buf).await.is_err() {
                break;
            }
            let Ok(payload) = serde_json::from_slice::<Payload<serde_json::Value>>(&buf) else {
                continue;
            };
            self.dispatch(payload);
        }
        let _ = self.close().await;
    }

    /// Mirror of Go `Wait(ctx)`: resolves when the connection closes.
    pub async fn wait(&self) {
        if self.is_closed() {
            return;
        }
        self.inner.notify_closed.notified().await;
    }

    /// Mirror of Go `Send(eventType, data, requestID)`.
    pub async fn send<T: Serialize>(
        &self,
        event: &str,
        data: &T,
        request_id: &str,
    ) -> Result<(), PluginError> {
        let payload = Payload {
            event: event.to_string(),
            request_id: request_id.to_string(),
            data,
        };
        let bytes = serde_json::to_vec(&payload)
            .map_err(|e| PluginError::CreatePayloadFailed(e.to_string()))?;
        let frame = encode_message(&bytes)?;

        #[cfg(not(windows))]
        {
            use tokio::io::AsyncWriteExt;
            let mut guard = self.inner.write_half.lock().await;
            let Some(write_half) = guard.as_mut() else {
                return Err(PluginError::ConnectionNotReady);
            };
            write_half.write_all(&frame).await.map_err(PluginError::from)?;
            Ok(())
        }
        #[cfg(windows)]
        {
            use tokio::io::AsyncWriteExt;
            let mut guard = self.inner.write_half.lock().await;
            let Some(client) = guard.as_mut() else {
                return Err(PluginError::ConnectionNotReady);
            };
            client.write_all(&frame).await.map_err(PluginError::from)?;
            Ok(())
        }
    }

    /// Mirror of Go `Close()`.
    pub async fn close(&self) -> Result<(), PluginError> {
        #[cfg(not(windows))]
        {
            use tokio::io::AsyncWriteExt;
            let mut guard = self.inner.write_half.lock().await;
            if let Some(mut w) = guard.take() {
                let _ = w.shutdown().await;
            }
        }
        #[cfg(windows)]
        {
            let mut guard = self.inner.write_half.lock().await;
            guard.take();
        }
        self.mark_closed();
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn send_without_connection_returns_not_ready() {
        let c = Connector::new(None);
        let err = c.send("test.event", &"data", "req-1").await.unwrap_err();
        assert!(matches!(err, PluginError::ConnectionNotReady));
        // Close when already closed should succeed.
        c.close().await.unwrap();
    }

    #[tokio::test]
    async fn wait_returns_after_close() {
        let c = Connector::new(None);
        let c2 = c.clone();
        tokio::spawn(async move {
            tokio::time::sleep(std::time::Duration::from_millis(20)).await;
            c2.close().await.unwrap();
        });
        tokio::time::timeout(std::time::Duration::from_secs(2), c.wait())
            .await
            .expect("wait should resolve after close");
    }
}
