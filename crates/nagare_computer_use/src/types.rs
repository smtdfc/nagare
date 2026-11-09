#[derive(serde::Deserialize)]
pub struct GetAppInstanceRequest {
    pub name: String,
}

#[derive(serde::Serialize)]
pub struct GetAppInstanceResponse {
    pub is_error: bool,
    pub instance_id: String,
    pub error: Option<String>,
}

#[allow(dead_code)]
#[derive(serde::Deserialize)]
pub struct ListRunningAppsRequest {}

#[derive(serde::Serialize, Clone)]
pub struct AppInfo {
    pub name: String,
    pub pid: Option<u32>,
}

#[derive(serde::Serialize)]
pub struct ListRunningAppsResponse {
    pub is_error: bool,
    pub apps: Vec<AppInfo>,
    pub error: Option<String>,
}

#[derive(serde::Deserialize)]
pub struct OpenAppRequest {
    pub name: String,
    pub args: Option<Vec<String>>,
}

#[derive(serde::Serialize)]
pub struct OpenAppResponse {
    pub is_error: bool,
    pub pid: Option<u32>,
    pub error: Option<String>,
}

#[derive(serde::Deserialize)]
pub struct GetUITreeRequest {
    pub instance_id: String,
    pub max_depth: Option<usize>,
}

#[derive(serde::Serialize, Clone)]
pub struct UINode {
    pub ref_id: String,
    pub role: String,
    pub name: String,
    pub value: Option<String>,
    pub x: i32,
    pub y: i32,
    pub width: u32,
    pub height: u32,
    pub children: Vec<UINode>,
}

#[derive(serde::Serialize)]
pub struct GetUITreeResponse {
    pub is_error: bool,
    pub tree: Option<UINode>,
    pub summary: String,
    pub error: Option<String>,
}

#[derive(serde::Deserialize)]
pub struct ClickElementRequest {
    pub instance_id: String,
    pub ref_id: String,
}

#[derive(serde::Serialize)]
pub struct ClickElementResponse {
    pub is_error: bool,
    pub error: Option<String>,
}

#[derive(serde::Deserialize)]
pub struct InputTextRequest {
    pub instance_id: Option<String>,
    pub ref_id: Option<String>,
    pub text: String,
}

#[derive(serde::Serialize)]
pub struct InputTextResponse {
    pub is_error: bool,
    pub error: Option<String>,
}

#[derive(serde::Deserialize)]
pub struct MouseActionRequest {
    pub action: String, // "move", "click", "right_click", "double_click", "scroll"
    pub x: Option<f64>,
    pub y: Option<f64>,
    pub scroll_delta_x: Option<f64>,
    pub scroll_delta_y: Option<f64>,
}

#[derive(serde::Serialize)]
pub struct MouseActionResponse {
    pub is_error: bool,
    pub error: Option<String>,
}

#[derive(serde::Deserialize)]
pub struct KeyboardShortcutRequest {
    pub key: String,
    pub modifiers: Option<Vec<String>>, // "ctrl", "alt", "shift", "meta"
}

#[derive(serde::Serialize)]
pub struct KeyboardShortcutResponse {
    pub is_error: bool,
    pub error: Option<String>,
}

#[derive(serde::Deserialize)]
pub struct WindowControlRequest {
    pub instance_id: String,
    pub action: String, // "activate", "minimize", "maximize", "restore", "close"
}

#[derive(serde::Serialize)]
pub struct WindowControlResponse {
    pub is_error: bool,
    pub error: Option<String>,
}

#[derive(serde::Deserialize)]
pub struct QueryElementByTextRequest {
    pub instance_id: String,
    pub text: String,
    pub exact: Option<bool>,
}

#[derive(serde::Serialize)]
pub struct ElementMatch {
    pub ref_id: String,
    pub role: String,
    pub name: String,
    pub value: Option<String>,
    pub x: i32,
    pub y: i32,
    pub width: u32,
    pub height: u32,
}

#[derive(serde::Serialize)]
pub struct QueryElementByTextResponse {
    pub is_error: bool,
    pub matches: Vec<ElementMatch>,
    pub error: Option<String>,
}
