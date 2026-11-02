use serde::{Serialize, de::DeserializeOwned};

use crate::{dtos::Payload, errors::ModuleError};

/// Mirror of Go `ipc.GetDefaultSocketPath` for external modules.
pub fn default_socket_path() -> String {
    if let Ok(p) = std::env::var("NAGARE_MODULE_SOCKET_PATH") {
        if !p.is_empty() {
            return p;
        }
    }
    #[cfg(windows)]
    {
        return r"\\.\pipe\nagare-external-module".to_string();
    }
    #[cfg(not(windows))]
    {
        if let Some(home) = dirs_home() {
            return format!("{home}/.nagare/nagare-external-module.sock");
        }
        "/tmp/nagare-external-module.sock".to_string()
    }
}

#[cfg(not(windows))]
fn dirs_home() -> Option<String> {
    std::env::var("HOME").ok()
}

/// Mirror of Go `GetData[T]`: re-decode the generic `data` field into `T`.
pub fn get_data<T: DeserializeOwned>(
    payload: &Payload<serde_json::Value>,
) -> Result<T, ModuleError> {
    if payload.data.is_null() {
        return Err(ModuleError::NilPayload);
    }
    serde_json::from_value(payload.data.clone()).map_err(ModuleError::from)
}

/// Mirror of Go `helpers.MarshalJson` / `UnmarshalJson`.
pub fn unmarshal_json<T: DeserializeOwned>(raw: &str) -> Result<T, ModuleError> {
    serde_json::from_str(raw).map_err(ModuleError::from)
}

pub fn marshal_json<T: Serialize>(data: &T) -> Result<String, ModuleError> {
    serde_json::to_string(data).map_err(ModuleError::from)
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::dtos::Payload;
    use serde::Deserialize;

    #[derive(Debug, Deserialize, PartialEq)]
    struct SampleData {
        message: String,
    }

    #[test]
    fn get_data_null_payload() {
        let payload = Payload {
            event: "test.event".to_string(),
            request_id: "r1".to_string(),
            data: serde_json::Value::Null,
        };
        assert!(matches!(
            get_data::<SampleData>(&payload).unwrap_err(),
            ModuleError::NilPayload
        ));
    }

    #[test]
    fn get_data_valid_payload() {
        let payload = Payload {
            event: "test.event".to_string(),
            request_id: "r1".to_string(),
            data: serde_json::json!({"message": "hello"}),
        };
        let res = get_data::<SampleData>(&payload).unwrap();
        assert_eq!(
            res,
            SampleData {
                message: "hello".to_string()
            }
        );
    }
}
