use nagare_module::Payload;
use serde_json::Value;
use xa11y::{Key, input_sim};

use crate::module::ComputerUseModule;
use crate::types::*;

fn parse_key(k: &str) -> Option<Key> {
    match k.to_lowercase().as_str() {
        "enter" | "return" => Some(Key::Enter),
        "escape" | "esc" => Some(Key::Escape),
        "backspace" => Some(Key::Backspace),
        "tab" => Some(Key::Tab),
        "space" => Some(Key::Space),
        "delete" => Some(Key::Delete),
        "insert" => Some(Key::Insert),
        "up" | "arrowup" => Some(Key::ArrowUp),
        "down" | "arrowdown" => Some(Key::ArrowDown),
        "left" | "arrowleft" => Some(Key::ArrowLeft),
        "right" | "arrowright" => Some(Key::ArrowRight),
        "home" => Some(Key::Home),
        "end" => Some(Key::End),
        "pageup" => Some(Key::PageUp),
        "pagedown" => Some(Key::PageDown),
        "ctrl" | "control" => Some(Key::Ctrl),
        "alt" => Some(Key::Alt),
        "shift" => Some(Key::Shift),
        "meta" | "win" | "super" | "cmd" => Some(Key::Meta),
        f if f.starts_with('f') && f.len() > 1 => f[1..].parse::<u8>().ok().map(Key::F),
        c if c.chars().count() == 1 => c.chars().next().map(Key::Char),
        _ => None,
    }
}

impl ComputerUseModule {
    pub async fn on_keyboard_shortcut(&self, payload: &Payload<Value>) {
        let req: KeyboardShortcutRequest = match payload.bind_data() {
            Ok(req) => req,
            Err(e) => {
                let _ = self
                    .client
                    .reply(
                        "keyboard_shortcut",
                        &payload.request_id,
                        &KeyboardShortcutResponse {
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
                        "keyboard_shortcut",
                        &payload.request_id,
                        &KeyboardShortcutResponse {
                            is_error: true,
                            error: Some(e.to_string()),
                        },
                    )
                    .await;
                return;
            }
        };

        let main_key = match parse_key(&req.key) {
            Some(k) => k,
            None => {
                let _ = self
                    .client
                    .reply(
                        "keyboard_shortcut",
                        &payload.request_id,
                        &KeyboardShortcutResponse {
                            is_error: true,
                            error: Some(format!("Unknown key: {}", req.key)),
                        },
                    )
                    .await;
                return;
            }
        };

        let modifiers: Vec<Key> = req
            .modifiers
            .unwrap_or_default()
            .iter()
            .filter_map(|m| parse_key(m))
            .collect();

        let kb = sim.keyboard();
        let res = if modifiers.is_empty() {
            kb.press(main_key)
        } else {
            kb.chord(main_key, &modifiers)
        };

        let resp = match res {
            Ok(_) => KeyboardShortcutResponse {
                is_error: false,
                error: None,
            },
            Err(e) => KeyboardShortcutResponse {
                is_error: true,
                error: Some(e.to_string()),
            },
        };

        let _ = self
            .client
            .reply("keyboard_shortcut", &payload.request_id, &resp)
            .await;
    }
}
