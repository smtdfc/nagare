#[derive(Debug, Clone, Default)]
pub struct ConnectConfig {
    pub module_id: String,
    pub socket_path: String,
    pub connect_code: String,
}
