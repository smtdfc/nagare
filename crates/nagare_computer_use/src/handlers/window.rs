use nagare_module::Payload;
use serde_json::Value;


use crate::module::ComputerUseModule;
use crate::types::*;

impl ComputerUseModule {
    pub async fn on_window_control(&self, payload: &Payload<Value>) {
        let req: WindowControlRequest = match payload.bind_data() {
            Ok(req) => req,
            Err(e) => {
                let _ = self
                    .client
                    .reply(
                        "window_control",
                        &payload.request_id,
                        &WindowControlResponse {
                            is_error: true,
                            error: Some(format!("Invalid request: {e}")),
                        },
                    )
                    .await;
                return;
            }
        };

        let app_opt = {
            self.app_instances
                .read()
                .unwrap()
                .get(&req.instance_id)
                .cloned()
        };

        let app = match app_opt {
            Some(a) => a,
            None => {
                let _ = self
                    .client
                    .reply(
                        "window_control",
                        &payload.request_id,
                        &WindowControlResponse {
                            is_error: true,
                            error: Some(format!("App instance {} not found", req.instance_id)),
                        },
                    )
                    .await;
                return;
            }
        };

        let windows_res = app.windows();
        let res: Result<(), String> = match windows_res {
            Ok(windows) => {
                if let Some(win) = windows.first() {
                    match req.action.as_str() {
                        "activate" => win.activate().map_err(|e| e.to_string()),
                        "minimize" => win.minimize().map_err(|e| e.to_string()),
                        "maximize" => win.maximize().map_err(|e| e.to_string()),
                        "restore" => win.restore().map_err(|e| e.to_string()),
                        "close" => win.close().map_err(|e| e.to_string()),
                        other => Err(format!("Unknown window action: {other}")),
                    }
                } else {
                    match req.action.as_str() {
                        "activate" => app.as_element().activate().map_err(|e| e.to_string()),
                        other => Err(format!(
                            "No top-level window found to perform action: {other}"
                        )),
                    }
                }
            }
            Err(e) => Err(e.to_string()),
        };

        let resp = match res {
            Ok(_) => WindowControlResponse {
                is_error: false,
                error: None,
            },
            Err(e) => WindowControlResponse {
                is_error: true,
                error: Some(e),
            },
        };

        let _ = self
            .client
            .reply("window_control", &payload.request_id, &resp)
            .await;
    }
}
