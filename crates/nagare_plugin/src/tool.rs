use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use std::time::Duration;

use schemars::JsonSchema;
use serde::{Serialize, de::DeserializeOwned};

use crate::dtos::{
    PluginToolCallEventPayload, PluginToolCallResultEventPayload, PluginToolInfo,
    REGISTER_PLUGIN_TOOL_EVENT, REGISTER_PLUGIN_TOOL_FAILED_EVENT,
    REGISTER_PLUGIN_TOOL_SUCCESS_EVENT, REGISTER_TOOL_CATEGORIES_EVENT,
    REGISTER_TOOL_CATEGORIES_FAILED_EVENT, REGISTER_TOOL_CATEGORIES_SUCCESS_EVENT,
    RegisterPluginToolEventPayload, RegisterPluginToolFailedEventPayload,
    RegisterPluginToolSuccessEventPayload, RegisterToolCategoriesEventPayload,
    RegisterToolCategoriesFailedEventPayload, RegisterToolCategoriesSuccessEventPayload,
};
use crate::errors::PluginError;
use crate::helpers::{marshal_json, unmarshal_json};

/// Mirror of Go `PluginTool` interface.
pub trait PluginTool: Send + Sync {
    fn name(&self) -> &str;
    fn description(&self) -> &str;
    fn args_schema(&self) -> String;
    fn categories(&self) -> Vec<String>;
    fn execute(&self, args_raw: &str) -> Result<String, PluginError>;
}

/// Mirror of Go `BaseToolCallback[I, O]`: `Fn(I) -> Result<O, PluginError>`.
pub type ToolCallback<I, O> = Arc<dyn Fn(I) -> Result<O, PluginError> + Send + Sync>;

/// Mirror of Go `BaseTool[I, O]`.
pub struct BaseTool<I, O> {
    pub tool_name: String,
    pub tool_description: String,
    pub callback: ToolCallback<I, O>,
    pub tool_categories: Vec<String>,
    schema_cache: Mutex<Option<String>>,
    _phantom: std::marker::PhantomData<(I, O)>,
}

impl<I, O> BaseTool<I, O>
where
    I: DeserializeOwned + JsonSchema,
    O: Serialize,
{
    pub fn new(
        name: impl Into<String>,
        description: impl Into<String>,
        callback: impl Fn(I) -> Result<O, PluginError> + Send + Sync + 'static,
        categories: Vec<String>,
    ) -> Self {
        Self {
            tool_name: name.into(),
            tool_description: description.into(),
            callback: Arc::new(callback),
            tool_categories: categories,
            schema_cache: Mutex::new(None),
            _phantom: std::marker::PhantomData,
        }
    }
}

impl<I, O> PluginTool for BaseTool<I, O>
where
    I: DeserializeOwned + JsonSchema + Send + Sync,
    O: Serialize + Send + Sync,
{
    fn name(&self) -> &str {
        &self.tool_name
    }

    fn description(&self) -> &str {
        &self.tool_description
    }

    fn args_schema(&self) -> String {
        {
            let guard = self.schema_cache.lock().unwrap();
            if let Some(cached) = guard.as_ref() {
                return cached.clone();
            }
        }
        let schema = schemars::schema_for!(I);
        let schema_str = serde_json::to_string(&schema).unwrap_or_else(|_| "{}".to_string());
        *self.schema_cache.lock().unwrap() = Some(schema_str.clone());
        schema_str
    }

    fn categories(&self) -> Vec<String> {
        self.tool_categories.clone()
    }

    fn execute(&self, args_raw: &str) -> Result<String, PluginError> {
        let args: I = unmarshal_json(args_raw).map_err(|_| PluginError::IncorrectToolArgs)?;
        let result = (self.callback)(args)?;
        marshal_json(&result).map_err(|e| PluginError::MarshalToolResultFailed(e.to_string()))
    }
}

/// Mirror of Go `DefineTool[I, O]`.
pub fn define_tool<I, O>(
    name: impl Into<String>,
    description: impl Into<String>,
    callback: impl Fn(I) -> Result<O, PluginError> + Send + Sync + 'static,
    categories: Vec<String>,
) -> Arc<BaseTool<I, O>>
where
    I: DeserializeOwned + JsonSchema + Send + Sync + 'static,
    O: Serialize + Send + Sync + 'static,
{
    Arc::new(BaseTool::new(name, description, callback, categories))
}

// ---- PluginClient tool methods (mirror of Go tool.go) ----

use crate::client::PluginClient;

impl PluginClient {
    fn tool_call_timeout() -> Duration {
        Duration::from_secs(60)
    }

    pub async fn register_plugin_tool(&self, tool: Arc<dyn PluginTool>) -> Result<(), PluginError> {
        let payload = RegisterPluginToolEventPayload {
            tool: PluginToolInfo {
                name: tool.name().to_string(),
                description: tool.description().to_string(),
                args: tool.args_schema(),
                categories: tool.categories(),
            },
        };
        let tool_name = tool.name().to_string();
        let resp = self
            .send_and_wait::<
                RegisterPluginToolSuccessEventPayload,
                RegisterPluginToolFailedEventPayload,
            >(
                REGISTER_PLUGIN_TOOL_EVENT,
                &payload,
                REGISTER_PLUGIN_TOOL_FAILED_EVENT,
                REGISTER_PLUGIN_TOOL_SUCCESS_EVENT,
            )
            .await
            .map_err(|e| {
                tracing::error!(error = %e, tool = %tool_name, "failed to register plugin tool");
                e
            })?;

        if !resp.is_success {
            let cause = resp
                .error
                .map(|e| e.cause)
                .unwrap_or_else(|| "unknown".to_string());
            tracing::error!(tool = %tool_name, cause = %cause, "failed to register plugin tool");
            return Err(PluginError::Other(cause));
        }

        self.inner.tools.lock().unwrap().push(tool);
        Ok(())
    }

    pub async fn register_tool_categories(
        &self,
        categories: HashMap<String, String>,
    ) -> Result<(), PluginError> {
        let payload = RegisterToolCategoriesEventPayload { categories };
        let resp = self
            .send_and_wait::<
                RegisterToolCategoriesSuccessEventPayload,
                RegisterToolCategoriesFailedEventPayload,
            >(
                REGISTER_TOOL_CATEGORIES_EVENT,
                &payload,
                REGISTER_TOOL_CATEGORIES_FAILED_EVENT,
                REGISTER_TOOL_CATEGORIES_SUCCESS_EVENT,
            )
            .await
            .map_err(|e| {
                tracing::error!(error = %e, "failed to register tool categories");
                e
            })?;

        if !resp.is_success {
            let cause = resp
                .error
                .map(|e| e.cause)
                .unwrap_or_else(|| "unknown".to_string());
            tracing::error!(cause = %cause, "failed to register tool categories");
            return Err(PluginError::Other(cause));
        }
        Ok(())
    }

    pub(crate) async fn handle_plugin_tool_call(
        &self,
        request_id: String,
        call: PluginToolCallEventPayload,
    ) {
        let tool: Option<Arc<dyn PluginTool>> = {
            let guard = self.inner.tools.lock().unwrap();
            guard.iter().find(|t| t.name() == call.name).cloned()
        };

        let Some(tool) = tool else {
            let msg = format!("plugin tool {:?} is not registered", call.name);
            self.send_plugin_tool_result(&request_id, "", &msg).await;
            return;
        };

        // Execute on a blocking thread with a 60s timeout (mirrors Go context.WithTimeout).
        let args = call.args.clone();
        let exec: tokio::task::JoinHandle<Result<String, PluginError>> =
            tokio::task::spawn_blocking(move || tool.execute(args.as_str()));
        let result: Result<
            Result<Result<String, PluginError>, tokio::task::JoinError>,
            tokio::time::error::Elapsed,
        > = tokio::time::timeout(Self::tool_call_timeout(), exec).await;

        match result {
            Err(_) => {
                self.send_plugin_tool_result(&request_id, "", "timeout")
                    .await;
            }
            Ok(Err(join_err)) => {
                self.send_plugin_tool_result(&request_id, "", &join_err.to_string())
                    .await;
            }
            Ok(Ok(Err(e))) => {
                self.send_plugin_tool_result(&request_id, "", &e.to_string())
                    .await;
            }
            Ok(Ok(Ok(output))) => {
                self.send_plugin_tool_result(&request_id, output.as_str(), "").await;
            }
        }
    }

    pub(crate) async fn send_plugin_tool_result(
        &self,
        request_id: &str,
        result: &str,
        call_error: &str,
    ) {
        use crate::dtos::PLUGIN_TOOL_CALL_RESULT_EVENT;
        let payload = PluginToolCallResultEventPayload {
            result: result.to_string(),
            error: call_error.to_string(),
        };
        if let Err(e) = self
            .inner
            .connector
            .send(PLUGIN_TOOL_CALL_RESULT_EVENT, &payload, request_id)
            .await
        {
            tracing::error!(error = %e, request_id = %request_id, "failed to send plugin tool result");
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use schemars::JsonSchema;
    use serde::{Deserialize, Serialize};

    #[derive(Debug, Deserialize, JsonSchema)]
    struct SampleArgs {
        query: String,
    }

    #[derive(Debug, Serialize)]
    struct SampleResult {
        answer: String,
    }

    #[test]
    fn define_tool_and_execution() {
        let tool = define_tool(
            "search_tool",
            "searches things",
            |args: SampleArgs| -> Result<SampleResult, PluginError> {
                Ok(SampleResult {
                    answer: format!("result for {}", args.query),
                })
            },
            vec!["search".to_string()],
        );

        assert_eq!(tool.name(), "search_tool");
        assert_eq!(tool.description(), "searches things");
        assert_eq!(tool.categories(), vec!["search".to_string()]);

        let schema1 = tool.args_schema();
        assert!(!schema1.is_empty() && schema1 != "{}");
        let schema2 = tool.args_schema();
        assert_eq!(schema1, schema2, "schema cache mismatch");

        let res = tool.execute(r#"{"query":"golang"}"#).unwrap();
        assert_eq!(res, r#"{"answer":"result for golang"}"#);

        let err = tool.execute("{invalid_json").unwrap_err();
        assert!(matches!(err, PluginError::IncorrectToolArgs));
    }
}
