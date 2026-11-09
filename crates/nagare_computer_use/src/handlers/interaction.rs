use std::time::Duration;

use nagare_module::Payload;
use serde_json::Value;
use xa11y::input_sim;

use crate::module::ComputerUseModule;
use crate::types::*;

impl ComputerUseModule {
    pub async fn on_click_element(&self, payload: &Payload<Value>) {
        let req: ClickElementRequest = match payload.bind_data() {
            Ok(req) => req,
            Err(e) => {
                let _ = self
                    .client
                    .reply(
                        "click_element",
                        &payload.request_id,
                        &ClickElementResponse {
                            is_error: true,
                            error: Some(format!("Invalid request: {e}")),
                        },
                    )
                    .await;
                return;
            }
        };

        let elem_opt = {
            let cache = self.element_cache.read().unwrap();
            cache
                .get(&req.instance_id)
                .and_then(|map| map.get(&req.ref_id).cloned())
        };

        let elem = match elem_opt {
            Some(e) => e,
            None => {
                let _ = self
                    .client
                    .reply(
                        "click_element",
                        &payload.request_id,
                        &ClickElementResponse {
                            is_error: true,
                            error: Some(format!(
                                "Element ref '{}' not found in cache. Call get_ui_tree_tool first to refresh UI state.",
                                req.ref_id
                            )),
                        },
                    )
                    .await;
                return;
            }
        };

        let resp = match elem.perform_action("press") {
            Ok(_) => ClickElementResponse {
                is_error: false,
                error: None,
            },
            Err(_) => {
                // Fallback: click via bounding box coordinates
                if let Some(bounds) = elem.data().bounds {
                    let center_x = bounds.x + (bounds.width as i32 / 2);
                    let center_y = bounds.y + (bounds.height as i32 / 2);
                    match input_sim() {
                        Ok(sim) => match sim.mouse().click((center_x, center_y)) {
                            Ok(_) => ClickElementResponse {
                                is_error: false,
                                error: None,
                            },
                            Err(e) => ClickElementResponse {
                                is_error: true,
                                error: Some(format!("Failed to click element: {e}")),
                            },
                        },
                        Err(e) => ClickElementResponse {
                            is_error: true,
                            error: Some(format!("Failed to init input sim: {e}")),
                        },
                    }
                } else {
                    ClickElementResponse {
                        is_error: true,
                        error: Some("Element has no bounds and action failed".to_string()),
                    }
                }
            }
        };

        let _ = self
            .client
            .reply("click_element", &payload.request_id, &resp)
            .await;
    }

    pub async fn on_input_text(&self, payload: &Payload<Value>) {
        let req: InputTextRequest = match payload.bind_data() {
            Ok(req) => req,
            Err(e) => {
                let _ = self
                    .client
                    .reply(
                        "input_text",
                        &payload.request_id,
                        &InputTextResponse {
                            is_error: true,
                            error: Some(format!("Invalid request: {e}")),
                        },
                    )
                    .await;
                return;
            }
        };

        // If a specific element ref is provided, target it directly
        if let (Some(inst_id), Some(ref_id)) = (&req.instance_id, &req.ref_id) {
            let elem_opt = {
                let cache = self.element_cache.read().unwrap();
                cache.get(inst_id).and_then(|map| map.get(ref_id).cloned())
            };

            if let Some(elem) = elem_opt {
                let resp = match elem.set_value(&req.text) {
                    Ok(_) => InputTextResponse {
                        is_error: false,
                        error: None,
                    },
                    Err(_) => {
                        // Fallback: click to focus then type
                        if let Some(bounds) = elem.data().bounds {
                            let center_x = bounds.x + (bounds.width as i32 / 2);
                            let center_y = bounds.y + (bounds.height as i32 / 2);
                            if let Ok(sim) = input_sim() {
                                let _ = sim.mouse().click((center_x, center_y));
                                tokio::time::sleep(Duration::from_millis(100)).await;
                                let _ = sim.keyboard().type_text(&req.text);
                            }
                        }
                        InputTextResponse {
                            is_error: false,
                            error: None,
                        }
                    }
                };
                let _ = self
                    .client
                    .reply("input_text", &payload.request_id, &resp)
                    .await;
                return;
            }
        }

        // Fallback: type into currently focused control
        let resp = match input_sim() {
            Ok(sim) => match sim.keyboard().type_text(&req.text) {
                Ok(_) => InputTextResponse {
                    is_error: false,
                    error: None,
                },
                Err(e) => InputTextResponse {
                    is_error: true,
                    error: Some(e.to_string()),
                },
            },
            Err(e) => InputTextResponse {
                is_error: true,
                error: Some(e.to_string()),
            },
        };

        let _ = self
            .client
            .reply("input_text", &payload.request_id, &resp)
            .await;
    }
}
