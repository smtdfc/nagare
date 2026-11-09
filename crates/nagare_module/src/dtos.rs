use serde::{Deserialize, Serialize, de::DeserializeOwned};

pub const INVOKE_PREFIX: &str = "module:invoke";
pub const REPLY_PREFIX: &str = "module:reply";

pub fn make_invoke_event(action: &str) -> String {
    if action.is_empty() {
        INVOKE_PREFIX.to_string()
    } else {
        format!("{INVOKE_PREFIX}:{action}")
    }
}

pub fn make_reply_event(action: &str) -> String {
    if action.is_empty() {
        REPLY_PREFIX.to_string()
    } else {
        format!("{REPLY_PREFIX}:{action}")
    }
}

pub fn parse_action_event(event: &str) -> Option<(&str, &str)> {
    if let Some(rest) = event.strip_prefix(INVOKE_PREFIX) {
        let action = rest.strip_prefix(':').unwrap_or("");
        Some((INVOKE_PREFIX, action))
    } else if let Some(rest) = event.strip_prefix(REPLY_PREFIX) {
        let action = rest.strip_prefix(':').unwrap_or("");
        Some((REPLY_PREFIX, action))
    } else {
        None
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Payload<T> {
    pub event: String,
    #[serde(rename = "requestID")]
    pub request_id: String,
    pub data: T,
}

impl<T> Payload<T> {
    pub fn action(&self) -> Option<&str> {
        parse_action_event(&self.event).map(|(_, action)| action)
    }

    pub fn is_invoke(&self) -> bool {
        self.event.starts_with(INVOKE_PREFIX)
    }

    pub fn is_reply(&self) -> bool {
        self.event.starts_with(REPLY_PREFIX)
    }
}

impl Payload<serde_json::Value> {
    pub fn bind_data<U: DeserializeOwned>(&self) -> Result<U, serde_json::Error> {
        serde_json::from_value(self.data.clone())
    }
}
