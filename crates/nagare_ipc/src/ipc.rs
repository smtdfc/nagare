use crate::errors::IpcError;

/// Mirrors `pkgs/ipc.MaxMessageSize` on the Go side.
pub const MAX_MESSAGE_SIZE: usize = 64 * 1024 * 1024;

/// Encode a frame: 4-byte big-endian length prefix + body.
pub fn encode_message(data: &[u8]) -> Result<Vec<u8>, IpcError> {
    if data.len() > MAX_MESSAGE_SIZE {
        return Err(IpcError::MessageTooLarge);
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

    #[test]
    fn encode_rejects_oversized_payload() {
        let data = vec![0u8; MAX_MESSAGE_SIZE + 1];
        assert!(matches!(
            encode_message(&data).unwrap_err(),
            IpcError::MessageTooLarge
        ));
    }
}
