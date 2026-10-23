use std::io;
use std::sync::{
    Arc, RwLock,
    atomic::{AtomicBool, Ordering},
};

use serde::Serialize;
use tokio::sync::{Mutex, Notify};

use crate::dtos::Payload;
use crate::errors::PluginError;
use crate::ipc::{MAX_MESSAGE_SIZE, encode_message};

/// Callback invoked for every incoming payload.
pub type OnEvent = Arc<dyn Fn(Payload<serde_json::Value>) + Send + Sync>;

/// Length-prefixed framed transport + listen loop.
#[derive(Clone)]
pub struct Connector {
    inner: Arc<ConnectorInner>,
}

struct ConnectorInner {
    #[cfg(not(windows))]
    write_half: Mutex<Option<tokio::net::unix::OwnedWriteHalf>>,

    #[cfg(windows)]
    write_half: RwLock<Option<Arc<tokio::net::windows::named_pipe::NamedPipeClient>>>,

    #[cfg(windows)]
    write_lock: Mutex<()>,

    on_event: RwLock<Option<OnEvent>>,
    closed: AtomicBool,
    notify_closed: Notify,

    #[cfg(windows)]
    shutdown: Notify,
}

impl Connector {
    pub fn new(on_event: Option<OnEvent>) -> Self {
        Self {
            inner: Arc::new(ConnectorInner {
                #[cfg(not(windows))]
                write_half: Mutex::new(None),

                #[cfg(windows)]
                write_half: RwLock::new(None),

                #[cfg(windows)]
                write_lock: Mutex::new(()),

                on_event: RwLock::new(on_event),
                closed: AtomicBool::new(false),
                notify_closed: Notify::new(),

                #[cfg(windows)]
                shutdown: Notify::new(),
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

    /// Connect to the plugin host.
    pub async fn connect(&self, path: &str) -> Result<(), PluginError> {
        #[cfg(not(windows))]
        {
            let stream = tokio::net::UnixStream::connect(path)
                .await
                .map_err(|e| PluginError::ConnectionFailed(e.to_string()))?;

            let (read_half, write_half) = stream.into_split();

            *self.inner.write_half.lock().await = Some(write_half);
            self.inner.closed.store(false, Ordering::SeqCst);

            let connector = self.clone();

            tokio::spawn(async move {
                connector.read_loop_unix(read_half).await;
            });

            Ok(())
        }

        #[cfg(windows)]
        {
            let client = Arc::new(
                tokio::net::windows::named_pipe::ClientOptions::new()
                    .open(path)
                    .map_err(|e| PluginError::ConnectionFailed(e.to_string()))?,
            );

            *self.inner.write_half.write().unwrap() = Some(Arc::clone(&client));
            self.inner.closed.store(false, Ordering::SeqCst);

            let connector = self.clone();

            tokio::spawn(async move {
                connector.read_loop_windows(client).await;
            });

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

            if msg_len == 0 || msg_len > MAX_MESSAGE_SIZE {
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

    #[cfg(windows)]
    async fn read_pipe_exact(
        client: &tokio::net::windows::named_pipe::NamedPipeClient,
        mut buf: &mut [u8],
    ) -> io::Result<()> {
        while !buf.is_empty() {
            client.readable().await?;

            match client.try_read(buf) {
                Ok(0) => {
                    return Err(io::Error::new(
                        io::ErrorKind::UnexpectedEof,
                        "named pipe closed",
                    ));
                }
                Ok(n) => {
                    let (_, remaining) = buf.split_at_mut(n);
                    buf = remaining;
                }
                Err(e) if e.kind() == io::ErrorKind::WouldBlock => continue,
                Err(e) => return Err(e),
            }
        }

        Ok(())
    }

    #[cfg(windows)]
    async fn read_loop_windows(
        &self,
        client: Arc<tokio::net::windows::named_pipe::NamedPipeClient>,
    ) {
        loop {
            if self.is_closed() {
                break;
            }

            let mut len_buf = [0u8; 4];

            let read_result = tokio::select! {
                _ = self.inner.shutdown.notified() => break,
                result = Self::read_pipe_exact(&client, &mut len_buf) => result,
            };

            if read_result.is_err() {
                break;
            }

            let msg_len = u32::from_be_bytes(len_buf) as usize;

            if msg_len == 0 || msg_len > MAX_MESSAGE_SIZE {
                break;
            }

            let mut buf = vec![0u8; msg_len];

            let read_result = tokio::select! {
                _ = self.inner.shutdown.notified() => break,
                result = Self::read_pipe_exact(&client, &mut buf) => result,
            };

            if read_result.is_err() {
                break;
            }

            let Ok(payload) = serde_json::from_slice::<Payload<serde_json::Value>>(&buf) else {
                continue;
            };

            self.dispatch(payload);
        }

        let _ = self.close().await;
    }

    /// Wait until the connection closes.
    pub async fn wait(&self) {
        let notified = self.inner.notify_closed.notified();
        tokio::pin!(notified);
        notified.as_mut().enable();

        if self.is_closed() {
            return;
        }

        notified.await;
    }

    /// Send a framed event.
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

            write_half
                .write_all(&frame)
                .await
                .map_err(PluginError::from)?;

            Ok(())
        }

        #[cfg(windows)]
        {
            let client = self
                .inner
                .write_half
                .read()
                .unwrap()
                .as_ref()
                .cloned()
                .ok_or(PluginError::ConnectionNotReady)?;

            let _write_guard = self.inner.write_lock.lock().await;

            if self.is_closed() {
                return Err(PluginError::ConnectionNotReady);
            }

            let mut remaining = frame.as_slice();

            while !remaining.is_empty() {
                client.writable().await.map_err(PluginError::from)?;

                match client.try_write(remaining) {
                    Ok(0) => {
                        return Err(PluginError::from(io::Error::new(
                            io::ErrorKind::WriteZero,
                            "failed to write to named pipe",
                        )));
                    }
                    Ok(n) => remaining = &remaining[n..],
                    Err(e) if e.kind() == io::ErrorKind::WouldBlock => continue,
                    Err(e) => return Err(PluginError::from(e)),
                }
            }

            Ok(())
        }
    }

    /// Close the connection.
    pub async fn close(&self) -> Result<(), PluginError> {
        #[cfg(not(windows))]
        {
            use tokio::io::AsyncWriteExt;

            let mut guard = self.inner.write_half.lock().await;

            if let Some(mut write_half) = guard.take() {
                let _ = write_half.shutdown().await;
            }
        }

        #[cfg(windows)]
        {
            self.inner.shutdown.notify_one();
            self.inner.write_half.write().unwrap().take();
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
        let connector = Connector::new(None);

        let err = connector
            .send("test.event", &"data", "req-1")
            .await
            .unwrap_err();

        assert!(matches!(err, PluginError::ConnectionNotReady));

        connector.close().await.unwrap();
    }

    #[tokio::test]
    async fn wait_returns_after_close() {
        let connector = Connector::new(None);
        let other = connector.clone();

        tokio::spawn(async move {
            tokio::time::sleep(std::time::Duration::from_millis(20)).await;
            other.close().await.unwrap();
        });

        tokio::time::timeout(std::time::Duration::from_secs(2), connector.wait())
            .await
            .expect("wait should resolve after close");
    }
}
