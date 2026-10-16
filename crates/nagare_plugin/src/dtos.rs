use serde::{Deserialize, Serialize};

// ---- Event names (mirror Go plugin/dtos) ----

pub const RECEIVED_CHAT_MESSAGE_EVENT: &str = "ReceivedChatMessageEvent";

pub const HANDSHAKE_EVENT: &str = "plugin:handshake";
pub const HANDSHAKE_SUCCESS_EVENT: &str = "plugin:handshake:success";
pub const HANDSHAKE_FAILED_EVENT: &str = "plugin:handshake:failed";

pub const PREPARE_CHAT_SESSION_EVENT: &str = "plugin:prepare_chat_session";
pub const PREPARE_CHAT_SESSION_SUCCESS_EVENT: &str = "plugin:prepare_chat_session:success";
pub const PREPARE_CHAT_SESSION_FAILED_EVENT: &str = "plugin:prepare_chat_session:failed";

pub const SEND_CHAT_MESSAGE_EVENT: &str = "plugin:send_chat_message";
pub const SEND_CHAT_MESSAGE_SUCCESS_EVENT: &str = "plugin:send_chat_message:success";
pub const SEND_CHAT_MESSAGE_FAILED_EVENT: &str = "plugin:send_chat_message:failed";

pub const RESET_CHAT_CHANNEL_EVENT: &str = "plugin:reset_chat_channel";
pub const RESET_CHAT_CHANNEL_SUCCESS_EVENT: &str = "plugin:reset_chat_channel:success";
pub const RESET_CHAT_CHANNEL_FAILED_EVENT: &str = "plugin:reset_chat_channel:failed";

pub const REGISTER_PLUGIN_TOOL_EVENT: &str = "plugin:register_plugin_tool";
pub const REGISTER_PLUGIN_TOOL_SUCCESS_EVENT: &str = "plugin:register_plugin_tool:success";
pub const REGISTER_PLUGIN_TOOL_FAILED_EVENT: &str = "plugin:register_plugin_tool:failed";

pub const REGISTER_TOOL_CATEGORIES_EVENT: &str = "plugin:register_tool_categories";
pub const REGISTER_TOOL_CATEGORIES_SUCCESS_EVENT: &str = "plugin:register_tool_categories:success";
pub const REGISTER_TOOL_CATEGORIES_FAILED_EVENT: &str = "plugin:register_tool_categories:failed";

pub const PLUGIN_TOOL_CALL_EVENT: &str = "plugin:tool_call";
pub const PLUGIN_TOOL_CALL_RESULT_EVENT: &str = "plugin:tool_call_result";

// ---- Base payload ----

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Payload<T> {
    pub event: String,
    #[serde(rename = "requestID")]
    pub request_id: String,
    pub data: T,
}

// ---- Chat DTOs ----

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ReceivedChatMessageEventPayload {
    #[serde(rename = "sessionID")]
    pub session_id: String,
    #[serde(rename = "channelID")]
    pub channel_id: String,
    pub message: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct HandshakeEventPayload {
    #[serde(default)]
    pub id: String,
    #[serde(rename = "packageName")]
    pub package_name: String,
    #[serde(rename = "connectCode")]
    pub connect_code: String,
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct HandshakeSuccessEventPayload {
    #[serde(default)]
    pub id: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct HandshakeFailedEventPayload {
    #[serde(default)]
    pub id: String,
    pub cause: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PrepareChatSessionEventPayload {
    #[serde(rename = "channelID")]
    pub channel_id: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PrepareChatSessionSuccessEventPayload {
    #[serde(rename = "channelID")]
    pub channel_id: String,
    #[serde(rename = "sessionID")]
    pub session_id: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PrepareChatSessionFailedEventPayload {
    #[serde(rename = "channelID")]
    pub channel_id: String,
    pub cause: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SendChatMessageEventPayload {
    #[serde(rename = "sessionID")]
    pub session_id: String,
    pub text: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SendChatMessageSuccessEventPayload {
    #[serde(rename = "sessionID")]
    pub session_id: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SendChatMessageFailedEventPayload {
    pub cause: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ResetChatChannelEventPayload {
    #[serde(rename = "channelID")]
    pub channel_id: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ResetChatChannelSuccessEventPayload {
    #[serde(rename = "channelID")]
    pub channel_id: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ResetChatChannelFailedEventPayload {
    pub cause: String,
}

// ---- Tool DTOs ----

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PluginToolInfo {
    pub name: String,
    pub args: String,
    pub description: String,
    pub categories: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RegisterPluginToolEventPayload {
    pub tool: PluginToolInfo,
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct RegisterPluginToolSuccessEventPayload {}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RegisterPluginToolFailedEventPayload {
    pub cause: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RegisterToolCategoriesEventPayload {
    pub categories: std::collections::HashMap<String, String>,
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct RegisterToolCategoriesSuccessEventPayload {}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RegisterToolCategoriesFailedEventPayload {
    pub cause: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PluginToolCallEventPayload {
    pub name: String,
    pub args: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct PluginToolCallResultEventPayload {
    pub result: String,
    #[serde(default, skip_serializing_if = "String::is_empty")]
    pub error: String,
}
