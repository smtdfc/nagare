use nagare_plugin::PluginClient;
use tokio::io::{AsyncReadExt, AsyncWriteExt};

async fn read_frame(stream: &mut tokio::net::UnixStream) -> serde_json::Value {
    let mut len_buf = [0u8; 4];
    stream.read_exact(&mut len_buf).await.unwrap();
    let len = u32::from_be_bytes(len_buf) as usize;
    let mut buf = vec![0u8; len];
    stream.read_exact(&mut buf).await.unwrap();
    serde_json::from_slice(&buf).unwrap()
}

async fn write_frame(stream: &mut tokio::net::UnixStream, v: &serde_json::Value) {
    let bytes = serde_json::to_vec(v).unwrap();
    stream
        .write_all(&(bytes.len() as u32).to_be_bytes())
        .await
        .unwrap();
    stream.write_all(&bytes).await.unwrap();
}

/// Full handshake round trip against a mock gateway speaking the same
/// length-prefixed JSON framing as Go `pkgs/ipc`.
#[tokio::test]
async fn handshake_round_trip_over_unix_socket() {
    let sock = format!(
        "{}/nagare-test-{}.sock",
        std::env::temp_dir().display(),
        std::process::id()
    );
    let _ = std::fs::remove_file(&sock);
    let listener = tokio::net::UnixListener::bind(&sock).unwrap();

    let server = tokio::spawn(async move {
        let (mut stream, _) = listener.accept().await.unwrap();
        let req: serde_json::Value = read_frame(&mut stream).await;
        assert_eq!(req["event"], "plugin:handshake");
        assert_eq!(req["data"]["packageName"], "nagare.plugin.test");
        let request_id = req["requestID"].as_str().unwrap().to_string();
        write_frame(
            &mut stream,
            &serde_json::json!({
                "event": "plugin:handshake:success",
                "requestID": request_id,
                "data": {},
            }),
        )
        .await;
        // Keep the connection open briefly so the client can finish.
        tokio::time::sleep(std::time::Duration::from_millis(200)).await;
    });

    let client = PluginClient::new();
    client.set_socket_path(&sock);
    client.set_connect_code("code-123");
    client
        .load_metadata(
            r#"{"package_name":"nagare.plugin.test","name":"t","version":"1","api_version":"1","author":"a"}"#,
        )
        .unwrap();
    client.connect().await.expect("connect failed");
    client.handshake().await.expect("handshake failed");

    server.await.unwrap();
    let _ = std::fs::remove_file(&sock);
}
