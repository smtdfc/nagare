use std::sync::{
    Arc, RwLock,
    atomic::{AtomicBool, Ordering},
};

#[cfg(windows)]
use std::io;

use serde::{Serialize, de::DeserializeOwned};
use tokio::sync::{Mutex, Notify};

use crate::{
    errors::IpcError,
    ipc::{MAX_MESSAGE_SIZE, encode_message},
};

/// Callback invoked for every incoming payload.
pub type OnEvent<T> = Arc<dyn Fn(T) + Send + Sync>;

/// Length-prefixed framed transport + listen loop.
///
/// The connector only knows about bytes and frames. The payload type is owned
/// by the SDK that uses it (`nagare_plugin` and `nagare_module` intentionally
/// keep separate payload shapes).
pub struct Connector<T: DeserializeOwned + 'static> {
    inner: Arc<ConnectorInner<T>>,
}

impl<T: DeserializeOwned + 'static> Clone for Connector<T> {
    fn clone(&self) -> Self {
        Self {
            inner: Arc::clone(&self.inner),
        }
    }
}

struct ConnectorInner<T> {
    #[cfg(not(windows))]
    write_half: Mutex<Option<tokio::net::unix::OwnedWriteHalf>>,

    #[cfg(windows)]
    write_half: RwLock<Option<Arc<tokio::net::windows::named_pipe::NamedPipeClient>>>,

    #[cfg(windows)]
    write_lock: Mutex<()>,

    on_event: RwLock<Option<OnEvent<T>>>,
    closed: AtomicBool,
    notify_closed: Notify,
    shutdown: Notify,
}

impl<T: DeserializeOwned + 'static> Connector<T> {
    pub fn new(on_event: Option<OnEvent<T>>) -> Self {
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
                shutdown: Notify::new(),
            }),
        }
    }

    pub fn set_on_event(&self, on_event: OnEvent<T>) {
        *self.inner.on_event.write().unwrap() = Some(on_event);
    }

    fn mark_closed(&self) {
        self.inner.closed.store(true, Ordering::SeqCst);
        self.inner.notify_closed.notify_waiters();
    }

    pub fn is_closed(&self) -> bool {
        self.inner.closed.load(Ordering::SeqCst)
    }

    fn dispatch(&self, payload: T) {
        let handler = self.inner.on_event.read().unwrap().clone();

        if let Some(handler) = handler {
            handler(payload);
        }
    }

    /// Completes once the connection is closed. Registering before checking
    /// `is_closed` keeps `close()` from racing past this future.
    async fn shutdown_signal(&self) {
        let notified = self.inner.shutdown.notified();
        tokio::pin!(notified);
        notified.as_mut().enable();

        if self.is_closed() {
            return;
        }

        notified.await;
    }

    /// Connect to the host.
    pub async fn connect(&self, path: &str) -> Result<(), IpcError> {
        #[cfg(not(windows))]
        {
            let stream = tokio::net::UnixStream::connect(path)
                .await
                .map_err(|e| IpcError::ConnectionFailed(e.to_string()))?;

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
                    .map_err(|e| IpcError::ConnectionFailed(e.to_string()))?,
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
            if self.is_closed() {
                break;
            }

            let mut len_buf = [0u8; 4];

            let read_result = tokio::select! {
                _ = self.shutdown_signal() => break,
                result = read_half.read_exact(&mut len_buf) => result,
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
                _ = self.shutdown_signal() => break,
                result = read_half.read_exact(&mut buf) => result,
            };

            if read_result.is_err() {
                break;
            }

            match serde_json::from_slice::<T>(&buf) {
                Ok(payload) => self.dispatch(payload),
                Err(e) => log::warn!("dropping malformed ipc payload: {e}"),
            }
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
                _ = self.shutdown_signal() => break,
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
                _ = self.shutdown_signal() => break,
                result = Self::read_pipe_exact(&client, &mut buf) => result,
            };

            if read_result.is_err() {
                break;
            }

            match serde_json::from_slice::<T>(&buf) {
                Ok(payload) => self.dispatch(payload),
                Err(e) => log::warn!("dropping malformed ipc payload: {e}"),
            }
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

    /// Send a framed JSON payload. The caller owns the payload shape.
    pub async fn send<P: Serialize + ?Sized>(&self, payload: &P) -> Result<(), IpcError> {
        if self.is_closed() {
            return Err(IpcError::ConnectionNotReady);
        }

        let bytes = serde_json::to_vec(payload).map_err(IpcError::from)?;
        let frame = encode_message(&bytes)?;

        #[cfg(not(windows))]
        {
            use tokio::io::AsyncWriteExt;

            let mut guard = self.inner.write_half.lock().await;

            let Some(write_half) = guard.as_mut() else {
                return Err(IpcError::ConnectionNotReady);
            };

            write_half.write_all(&frame).await.map_err(IpcError::from)?;

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
                .ok_or(IpcError::ConnectionNotReady)?;

            let _write_guard = self.inner.write_lock.lock().await;

            if self.is_closed() {
                return Err(IpcError::ConnectionNotReady);
            }

            let mut remaining = frame.as_slice();

            while !remaining.is_empty() {
                client.writable().await.map_err(IpcError::from)?;

                match client.try_write(remaining) {
                    Ok(0) => {
                        return Err(IpcError::from(io::Error::new(
                            io::ErrorKind::WriteZero,
                            "failed to write to named pipe",
                        )));
                    }
                    Ok(n) => remaining = &remaining[n..],
                    Err(e) if e.kind() == io::ErrorKind::WouldBlock => continue,
                    Err(e) => return Err(IpcError::from(e)),
                }
            }

            Ok(())
        }
    }

    /// Close the connection and stop the read loop.
    pub async fn close(&self) -> Result<(), IpcError> {
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
            self.inner.write_half.write().unwrap().take();
        }

        self.inner.shutdown.notify_waiters();
        self.mark_closed();

        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde_json::{Value, json};

    #[tokio::test]
    async fn send_without_connection_returns_not_ready() {
        let connector: Connector<Value> = Connector::new(None);

        let err = connector.send(&json!({"event": "test"})).await.unwrap_err();

        assert!(matches!(err, IpcError::ConnectionNotReady));

        connector.close().await.unwrap();
    }

    #[tokio::test]
    async fn wait_returns_after_close() {
        let connector: Connector<Value> = Connector::new(None);
        let other = connector.clone();

        tokio::spawn(async move {
            tokio::time::sleep(std::time::Duration::from_millis(20)).await;
            other.close().await.unwrap();
        });

        tokio::time::timeout(std::time::Duration::from_secs(2), connector.wait())
            .await
            .expect("wait should resolve after close");
    }

    #[cfg(unix)]
    #[tokio::test]
    async fn roundtrip_over_unix_socket() {
        use tokio::io::{AsyncReadExt, AsyncWriteExt};

        let unique = std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)
            .unwrap()
            .as_nanos();
        let path =
            std::env::temp_dir().join(format!("nagare-ipc-{}-{unique}.sock", std::process::id()));
        let listener = tokio::net::UnixListener::bind(&path).unwrap();

        let server = tokio::spawn(async move {
            let (mut stream, _) = listener.accept().await.unwrap();

            let mut len_buf = [0u8; 4];
            stream.read_exact(&mut len_buf).await.unwrap();
            let mut buf = vec![0u8; u32::from_be_bytes(len_buf) as usize];
            stream.read_exact(&mut buf).await.unwrap();

            let received: Value = serde_json::from_slice(&buf).unwrap();
            assert_eq!(received["event"], "test:ping");

            let reply = encode_message(
                &serde_json::to_vec(&json!({"event": "test:pong", "data": {"ok": true}})).unwrap(),
            )
            .unwrap();
            stream.write_all(&reply).await.unwrap();
            stream.flush().await.unwrap();

            // Keep the connection open until the client closes it.
            let mut eof = [0u8; 1];
            let _ = stream.read(&mut eof).await;
        });

        let (tx, mut rx) = tokio::sync::mpsc::unbounded_channel();
        let connector: Connector<Value> = Connector::new(Some(Arc::new(move |payload| {
            let _ = tx.send(payload);
        })));

        connector.connect(path.to_str().unwrap()).await.unwrap();
        connector
            .send(&json!({"event": "test:ping", "data": {"n": 1}}))
            .await
            .unwrap();

        let received = tokio::time::timeout(std::time::Duration::from_secs(2), rx.recv())
            .await
            .expect("timed out waiting for the reply")
            .expect("connection closed before the reply");

        assert_eq!(received["event"], "test:pong");
        assert_eq!(received["data"]["ok"], true);

        connector.close().await.unwrap();
        let _ = server.await;
        let _ = std::fs::remove_file(&path);
    }
}
