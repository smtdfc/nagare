use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Payload<T> {
    pub event: String,
    #[serde(rename = "requestID")]
    pub request_id: String,
    pub data: T,
}
