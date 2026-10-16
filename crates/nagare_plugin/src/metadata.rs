use serde::{Deserialize, Serialize};
use std::collections::HashMap;

use crate::errors::PluginError;

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct PluginMetadata {
    #[serde(default)]
    pub package_name: String,
    #[serde(default)]
    pub name: String,
    #[serde(default)]
    pub version: String,
    #[serde(default)]
    pub api_version: String,
    #[serde(default)]
    pub author: String,
    #[serde(default)]
    pub bin: HashMap<String, String>,
    #[serde(default)]
    pub features: Vec<String>,
}

impl PluginMetadata {
    pub fn validate(&self) -> Result<(), PluginError> {
        if self.package_name.trim().is_empty() {
            return Err(PluginError::Other(
                "plugin metadata 'package_name' is required".into(),
            ));
        }
        if self.name.trim().is_empty() {
            return Err(PluginError::Other(
                "plugin metadata 'name' is required".into(),
            ));
        }
        if self.version.trim().is_empty() {
            return Err(PluginError::Other(
                "plugin metadata 'version' is required".into(),
            ));
        }
        if self.api_version.trim().is_empty() {
            return Err(PluginError::Other(
                "plugin metadata 'api_version' is required".into(),
            ));
        }
        if self.author.trim().is_empty() {
            return Err(PluginError::Other(
                "plugin metadata 'author' is required".into(),
            ));
        }
        Ok(())
    }

    pub fn get_features(&self) -> String {
        self.features.join(",")
    }
}
