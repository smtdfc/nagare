use std::{
    collections::HashMap,
    path::PathBuf,
    sync::{Arc, RwLock},
};

use log::{error, info};
use nagare_module::ModuleClient;
use xa11y::App;

pub struct ComputerUseModule {
    pub client: Arc<ModuleClient>,
    pub app_instances: RwLock<HashMap<String, Arc<App>>>,
    pub element_cache: RwLock<HashMap<String, HashMap<String, xa11y::Element>>>,
}

impl ComputerUseModule {
    pub fn new() -> Result<Arc<Self>, String> {
        let data_path = PathBuf::from(
            std::env::var("NAGARE_MODULE_DATA_DIR").unwrap_or_else(|_| "./data".to_string()),
        );

        std::fs::create_dir_all(&data_path)
            .map_err(|e| format!("Failed to create data directory: {e}"))?;

        Ok(Arc::new(Self {
            client: ModuleClient::new_with_logger(),
            app_instances: RwLock::new(HashMap::new()),
            element_cache: RwLock::new(HashMap::new()),
        }))
    }

    pub async fn start(self: Arc<Self>) {
        let client = self.client.clone();
        let zelf_on_start = self.clone();

        let result = client
            .start(|| async move {
                info!("starting module ....");
                match zelf_on_start.client.handshake().await {
                    Ok(_) => {
                        info!("handshake success")
                    }
                    Err(e) => {
                        error!("Failed to handshake: {e}");
                        zelf_on_start.shutdown(false);
                    }
                }

                let zelf = zelf_on_start.clone();
                zelf_on_start.client.on_invoke(move |action, payload| {
                    let zself = zelf.clone();
                    tokio::spawn(async move {
                        match action.as_str() {
                            "get_app_instance" => zself.on_get_app_instance(&payload).await,
                            "list_running_apps" => zself.on_list_running_apps(&payload).await,
                            "open_app" => zself.on_open_app(&payload).await,
                            "get_ui_tree" => zself.on_get_ui_tree(&payload).await,
                            "click_element" => zself.on_click_element(&payload).await,
                            "input_text" => zself.on_input_text(&payload).await,
                            "mouse_action" => zself.on_mouse_action(&payload).await,
                            "keyboard_shortcut" => zself.on_keyboard_shortcut(&payload).await,
                            "window_control" => zself.on_window_control(&payload).await,
                            "query_element_by_text" => {
                                zself.on_query_element_by_text(&payload).await
                            }
                            _ => {}
                        }
                    });
                });
            })
            .await;

        match result {
            Ok(_) => {}
            Err(err) => {
                eprintln!("Error: {}", err.to_string());
                self.shutdown(false);
            }
        }
    }

    pub fn shutdown(&self, is_success: bool) {
        std::process::exit(if is_success { 0 } else { 1 });
    }
}
