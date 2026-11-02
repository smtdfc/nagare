use serde::{Deserialize, Serialize};

pub const ADD_VECTOR_EVENT: &str = "internal:vector:add";

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AddVectorEventPayload {
    #[serde(rename = "requestID")]
    pub request_id: String,
    #[serde(rename = "vector")]
    pub vector: Vec<f32>,
    #[serde(rename = "ids")]
    pub ids: Vec<u32>,
}
