use crate::errors::PluginError;

pub const MAX_MESSAGE_SIZE: usize = 64 * 1024 * 1024;

/// Mirror of Go `ipc.GetDefaultSocketPath`.
pub fn default_socket_path() -> String {
    if let Ok(p) = std::env::var("NAGARE_PLUGIN_SOCKET_PATH") {
        if !p.is_empty() {
            return p;
        }
    }
    #[cfg(windows)]
    {
        return r"\\.\pipe\nagare-plugin".to_string();
    }
    #[cfg(not(windows))]
    {
        if let Some(home) = dirs_home() {
            return format!("{home}/.nagare/nagare.sock");
        }
        "/tmp/nagare.sock".to_string()
    }
}

#[cfg(not(windows))]
fn dirs_home() -> Option<String> {
    std::env::var("HOME").ok()
}

/// Encode a frame: 4-byte big-endian length prefix + body.
pub fn encode_message(data: &[u8]) -> Result<Vec<u8>, PluginError> {
    if data.len() > MAX_MESSAGE_SIZE {
        return Err(PluginError::Other("ipc message exceeds maximum size".into()));
    }
    let mut buf = Vec::with_capacity(4 + data.len());
    buf.extend_from_slice(&(data.len() as u32).to_be_bytes());
    buf.extend_from_slice(data);
    Ok(buf)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn encode_roundtrip_len_prefix() {
        let buf = encode_message(b"hello").unwrap();
        assert_eq!(buf.len(), 9);
        let len = u32::from_be_bytes(buf[0..4].try_into().unwrap()) as usize;
        assert_eq!(len, 5);
        assert_eq!(&buf[4..], b"hello");
    }
}
