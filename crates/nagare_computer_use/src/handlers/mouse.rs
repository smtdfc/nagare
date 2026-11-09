use nagare_module::Payload;
use serde_json::Value;
use xa11y::{MouseButton, ScrollDelta, input_sim};

use crate::module::ComputerUseModule;
use crate::types::*;

impl ComputerUseModule {
    pub async fn on_mouse_action(&self, payload: &Payload<Value>) {
        let req: MouseActionRequest = match payload.bind_data() {
            Ok(req) => req,
            Err(e) => {
                let _ = self
                    .client
                    .reply(
                        "mouse_action",
                        &payload.request_id,
                        &MouseActionResponse {
                            is_error: true,
                            error: Some(format!("Invalid request: {e}")),
                        },
                    )
                    .await;
                return;
            }
        };

        let sim = match input_sim() {
            Ok(s) => s,
            Err(e) => {
                let _ = self
                    .client
                    .reply(
                        "mouse_action",
                        &payload.request_id,
                        &MouseActionResponse {
                            is_error: true,
                            error: Some(e.to_string()),
                        },
                    )
                    .await;
                return;
            }
        };

        let mouse = sim.mouse();
        let point = match (req.x, req.y) {
            (Some(x), Some(y)) => (x as i32, y as i32),
            _ => (0, 0),
        };

        let res = match req.action.as_str() {
            "move" => mouse.move_to(point),
            "click" => {
                if req.x.is_some() && req.y.is_some() {
                    mouse.click(point)
                } else {
                    mouse
                        .down(MouseButton::Left)
                        .and_then(|_| mouse.up(MouseButton::Left))
                }
            }
            "right_click" => {
                if req.x.is_some() && req.y.is_some() {
                    mouse.right_click(point)
                } else {
                    mouse
                        .down(MouseButton::Right)
                        .and_then(|_| mouse.up(MouseButton::Right))
                }
            }
            "double_click" => {
                if req.x.is_some() && req.y.is_some() {
                    mouse.double_click(point)
                } else {
                    mouse
                        .down(MouseButton::Left)
                        .and_then(|_| mouse.up(MouseButton::Left))
                        .and_then(|_| mouse.down(MouseButton::Left))
                        .and_then(|_| mouse.up(MouseButton::Left))
                }
            }
            "scroll" => {
                let dx = req.scroll_delta_x.unwrap_or(0.0) as i32;
                let dy = req.scroll_delta_y.unwrap_or(0.0) as i32;
                mouse.scroll(point, ScrollDelta::new(dx, dy))
            }
            other => {
                let _ = self
                    .client
                    .reply(
                        "mouse_action",
                        &payload.request_id,
                        &MouseActionResponse {
                            is_error: true,
                            error: Some(format!("Unknown mouse action: {other}")),
                        },
                    )
                    .await;
                return;
            }
        };

        let resp = match res {
            Ok(_) => MouseActionResponse {
                is_error: false,
                error: None,
            },
            Err(e) => MouseActionResponse {
                is_error: true,
                error: Some(e.to_string()),
            },
        };

        let _ = self
            .client
            .reply("mouse_action", &payload.request_id, &resp)
            .await;
    }
}
