use crate::errors::PluginError;
use serde::de::DeserializeOwned;
use serde::Serialize;

use crate::dtos::Payload;

/// Mirror of Go `GetData[T]`: re-decode the generic `data` field into `T`.
pub fn get_data<T: DeserializeOwned>(payload: &Payload<serde_json::Value>) -> Result<T, PluginError> {
    if payload.data.is_null() {
        return Err(PluginError::NilPayload);
    }
    serde_json::from_value(payload.data.clone()).map_err(PluginError::from)
}

/// Mirror of Go `helpers.MarshalJson` / `UnmarshalJson`.
pub fn unmarshal_json<T: DeserializeOwned>(raw: &str) -> Result<T, PluginError> {
    serde_json::from_str(raw).map_err(PluginError::from)
}

pub fn marshal_json<T: Serialize>(data: &T) -> Result<String, PluginError> {
    serde_json::to_string(data).map_err(PluginError::from)
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
            PluginError::NilPayload
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
