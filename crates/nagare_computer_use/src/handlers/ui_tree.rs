use std::collections::HashMap;

use nagare_module::Payload;
use serde_json::Value;


use crate::module::ComputerUseModule;
use crate::types::*;

impl ComputerUseModule {
    pub async fn on_get_ui_tree(&self, payload: &Payload<Value>) {
        let req: GetUITreeRequest = match payload.bind_data() {
            Ok(req) => req,
            Err(e) => {
                let _ = self
                    .client
                    .reply(
                        "get_ui_tree",
                        &payload.request_id,
                        &GetUITreeResponse {
                            is_error: true,
                            tree: None,
                            summary: String::new(),
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
                        "get_ui_tree",
                        &payload.request_id,
                        &GetUITreeResponse {
                            is_error: true,
                            tree: None,
                            summary: String::new(),
                            error: Some(format!("App instance {} not found", req.instance_id)),
                        },
                    )
                    .await;
                return;
            }
        };

        let max_depth = req.max_depth.unwrap_or(8);
        let mut element_map: HashMap<String, xa11y::Element> = HashMap::new();
        let mut counter: usize = 1;

        fn build_simplified_tree(
            elem: &xa11y::Element,
            current_depth: usize,
            max_depth: usize,
            counter: &mut usize,
            map: &mut HashMap<String, xa11y::Element>,
        ) -> Option<UINode> {
            if current_depth > max_depth {
                return None;
            }

            let data = elem.data();
            let role_str = data.role.to_snake_case().to_string();
            let name_val = data.name.clone().unwrap_or_default();
            let has_value = data.value.is_some();

            let is_interactive = matches!(
                role_str.as_str(),
                "button"
                    | "check_box"
                    | "combo_box"
                    | "menu_item"
                    | "radio_button"
                    | "text_field"
                    | "password_text"
                    | "slider"
                    | "spin_button"
                    | "tab"
                    | "link"
                    | "window"
                    | "dialog"
            ) || !name_val.is_empty()
                || has_value;

            let mut children_nodes = Vec::new();
            if let Ok(children) = elem.children() {
                for child in children {
                    if let Some(child_node) =
                        build_simplified_tree(&child, current_depth + 1, max_depth, counter, map)
                    {
                        children_nodes.push(child_node);
                    }
                }
            }

            if !is_interactive && children_nodes.is_empty() {
                return None;
            }

            let ref_id = format!("@e{}", *counter);
            *counter += 1;
            map.insert(ref_id.clone(), elem.clone());

            let bounds = data.bounds.unwrap_or(xa11y::Rect {
                x: 0,
                y: 0,
                width: 0,
                height: 0,
            });

            Some(UINode {
                ref_id,
                role: role_str,
                name: name_val,
                value: data.value.clone(),
                x: bounds.x,
                y: bounds.y,
                width: bounds.width,
                height: bounds.height,
                children: children_nodes,
            })
        }

        let root = app.as_element();
        let tree_node = build_simplified_tree(&root, 0, max_depth, &mut counter, &mut element_map);

        // Store updated element cache
        {
            let mut cache = self.element_cache.write().unwrap();
            cache.insert(req.instance_id.clone(), element_map);
        }

        fn format_summary(node: &UINode, indent: usize, out: &mut String) {
            let padding = "  ".repeat(indent);
            let val_str = node
                .value
                .as_ref()
                .map(|v| format!(" value=\"{v}\""))
                .unwrap_or_default();
            out.push_str(&format!(
                "{}[{}] {} \"{}\"{} [bounds: {},{} {}x{}]\n",
                padding,
                node.ref_id,
                node.role,
                node.name,
                val_str,
                node.x,
                node.y,
                node.width,
                node.height
            ));
            for child in &node.children {
                format_summary(child, indent + 1, out);
            }
        }

        let mut summary = String::new();
        if let Some(ref root_node) = tree_node {
            format_summary(root_node, 0, &mut summary);
        }

        let resp = GetUITreeResponse {
            is_error: false,
            tree: tree_node,
            summary,
            error: None,
        };

        let _ = self
            .client
            .reply("get_ui_tree", &payload.request_id, &resp)
            .await;
    }

    pub async fn on_query_element_by_text(&self, payload: &Payload<Value>) {
        let req: QueryElementByTextRequest = match payload.bind_data() {
            Ok(req) => req,
            Err(e) => {
                let _ = self
                    .client
                    .reply(
                        "query_element_by_text",
                        &payload.request_id,
                        &QueryElementByTextResponse {
                            is_error: true,
                            matches: Vec::new(),
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
                        "query_element_by_text",
                        &payload.request_id,
                        &QueryElementByTextResponse {
                            is_error: true,
                            matches: Vec::new(),
                            error: Some(format!("App instance {} not found", req.instance_id)),
                        },
                    )
                    .await;
                return;
            }
        };

        let is_exact = req.exact.unwrap_or(false);
        let query_text = req.text.to_lowercase();

        let mut matches = Vec::new();
        let mut matched_elements: Vec<xa11y::Element> = Vec::new();

        fn traverse(
            elem: &xa11y::Element,
            query: &str,
            exact: bool,
            matches: &mut Vec<ElementMatch>,
            matched_elements: &mut Vec<xa11y::Element>,
        ) {
            let data = elem.data();
            let mut is_match = false;

            if let Some(name) = &data.name {
                let n = name.to_lowercase();
                if exact && n == query {
                    is_match = true;
                } else if !exact && n.contains(query) {
                    is_match = true;
                }
            }

            if !is_match {
                if let Some(val) = &data.value {
                    let v = val.to_lowercase();
                    if exact && v == query {
                        is_match = true;
                    } else if !exact && v.contains(query) {
                        is_match = true;
                    }
                }
            }

            if is_match {
                let role_str = data.role.to_snake_case().to_string();
                let name_val = data.name.clone().unwrap_or_default();
                let bounds = data.bounds.unwrap_or(xa11y::Rect {
                    x: 0,
                    y: 0,
                    width: 0,
                    height: 0,
                });

                matches.push(ElementMatch {
                    ref_id: String::new(), // populated below under lock
                    role: role_str,
                    name: name_val,
                    value: data.value.clone(),
                    x: bounds.x,
                    y: bounds.y,
                    width: bounds.width,
                    height: bounds.height,
                });
                matched_elements.push(elem.clone());
            }

            if let Ok(children) = elem.children() {
                for child in children {
                    traverse(&child, query, exact, matches, matched_elements);
                }
            }
        }

        let root = app.as_element();
        traverse(
            &root,
            &query_text,
            is_exact,
            &mut matches,
            &mut matched_elements,
        );

        // Assign ref_ids and cache elements (drop the guard before await)
        {
            let mut cache_guard = self.element_cache.write().unwrap();
            let inst_cache = cache_guard
                .entry(req.instance_id.clone())
                .or_insert_with(HashMap::new);

            let mut next_ref_num = inst_cache.len() + 1;

            for (m, elem) in matches.iter_mut().zip(matched_elements.into_iter()) {
                let elem_handle = elem.handle;
                let existing_ref = inst_cache.iter().find_map(|(k, v)| {
                    if v.handle == elem_handle {
                        Some(k.clone())
                    } else {
                        None
                    }
                });

                let ref_id = match existing_ref {
                    Some(r) => r,
                    None => {
                        let new_ref = format!("@e{}", next_ref_num);
                        next_ref_num += 1;
                        inst_cache.insert(new_ref.clone(), elem);
                        new_ref
                    }
                };

                m.ref_id = ref_id;
            }
        }

        let resp = QueryElementByTextResponse {
            is_error: false,
            matches,
            error: None,
        };

        let _ = self
            .client
            .reply("query_element_by_text", &payload.request_id, &resp)
            .await;
    }
}
