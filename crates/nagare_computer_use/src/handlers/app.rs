use std::sync::Arc;
use std::time::Duration;

use log::{error, info};
use nagare_module::Payload;
use serde_json::Value;
use uuid::Uuid;
use xa11y::{App, AppExt};

use crate::module::ComputerUseModule;
use crate::types::*;

impl ComputerUseModule {
    pub async fn on_get_app_instance(&self, payload: &Payload<Value>) {
        let req: GetAppInstanceRequest = match payload.bind_data() {
            Ok(req) => req,
            Err(e) => {
                error!("Failed to parse GetAppInstanceRequest: {e}");
                return;
            }
        };

        info!("Getting app instance for: {}", req.name);

        let app = match App::by_name(&req.name, Duration::from_secs(5)) {
            Ok(a) => a,
            Err(e) => {
                if let Err(e) = self
                    .client
                    .reply(
                        "get_app_instance",
                        &payload.request_id,
                        &GetAppInstanceResponse {
                            is_error: true,
                            instance_id: "".to_owned(),
                            error: Some(e.to_string()),
                        },
                    )
                    .await
                {
                    error!("Failed to reply get_app_instance: {e}");
                }
                return;
            }
        };

        let id = Uuid::new_v4().to_string();
        self.app_instances
            .write()
            .unwrap()
            .insert(id.clone(), Arc::new(app));

        if let Err(e) = self
            .client
            .reply(
                "get_app_instance",
                &payload.request_id,
                &GetAppInstanceResponse {
                    is_error: false,
                    instance_id: id,
                    error: None,
                },
            )
            .await
        {
            error!("Failed to reply get_app_instance: {e}");
        }
    }

    pub async fn on_list_running_apps(&self, payload: &Payload<Value>) {
        info!("Listing running accessibility applications");
        let apps_result = App::list();
        let resp = match apps_result {
            Ok(apps) => {
                let list = apps
                    .into_iter()
                    .map(|a| AppInfo {
                        name: a.name,
                        pid: a.pid,
                    })
                    .collect();
                ListRunningAppsResponse {
                    is_error: false,
                    apps: list,
                    error: None,
                }
            }
            Err(e) => {
                error!("Failed to list running apps: {e}");
                ListRunningAppsResponse {
                    is_error: true,
                    apps: Vec::new(),
                    error: Some(e.to_string()),
                }
            }
        };

        if let Err(e) = self
            .client
            .reply("list_running_apps", &payload.request_id, &resp)
            .await
        {
            error!("Failed to reply list_running_apps: {e}");
        }
    }

    pub async fn on_open_app(&self, payload: &Payload<Value>) {
        let req: OpenAppRequest = match payload.bind_data() {
            Ok(req) => req,
            Err(e) => {
                error!("Failed to parse OpenAppRequest: {e}");
                let _ = self
                    .client
                    .reply(
                        "open_app",
                        &payload.request_id,
                        &OpenAppResponse {
                            is_error: true,
                            pid: None,
                            error: Some(format!("Invalid request: {e}")),
                        },
                    )
                    .await;
                return;
            }
        };

        info!("Opening app: {}", req.name);
        let mut cmd = std::process::Command::new(&req.name);
        if let Some(args) = req.args {
            cmd.args(args);
        }

        let resp = match cmd.spawn() {
            Ok(child) => OpenAppResponse {
                is_error: false,
                pid: Some(child.id()),
                error: None,
            },
            Err(e) => {
                error!("Failed to open app {}: {e}", req.name);
                OpenAppResponse {
                    is_error: true,
                    pid: None,
                    error: Some(e.to_string()),
                }
            }
        };

        if let Err(e) = self
            .client
            .reply("open_app", &payload.request_id, &resp)
            .await
        {
            error!("Failed to reply open_app: {e}");
        }
    }
}
