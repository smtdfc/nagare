export type AgentErrorDescription = {
  title: string;
  hint?: string;
  detail?: string;
};

const ERROR_META: Record<string, { title: string; hint?: string }> = {
  MODEL_QUOTA_EXCEED: {
    title: "Quota exceeded",
    hint: "The provider is rate limiting this account. Wait a moment and try again.",
  },
  MODEL_PAYMENT_REQUIRED: {
    title: "Payment required",
    hint: "The provider account needs a valid billing setup.",
  },
  MODEL_PROVIDER_ADAPTER_ERROR: {
    title: "Model provider error",
    hint: "The provider rejected the request or is temporarily unavailable. Try again, or switch model.",
  },
  MODEL_PROVIDER_NOT_SUPPORTED: {
    title: "Provider not supported",
    hint: "Choose a supported provider in Settings.",
  },
  MODEL_PROVIDER_NOT_FOUND: {
    title: "Provider not found",
    hint: "The configured provider no longer exists. Check Settings → Model Providers.",
  },
  MODEL_RESPONSE_EMPTY: {
    title: "Empty response",
    hint: "The model returned no content. Try sending the message again.",
  },
  MESSAGE_STREAM: {
    title: "Response interrupted",
    hint: "The response stream ended unexpectedly. Try again.",
  },
  MESSAGE_NOT_SUPPORT: {
    title: "Unsupported conversation",
    hint: "This conversation contains a message type the selected model cannot process. Start a new chat or switch model.",
  },
  CURRENT_PROVIDER_NOT_SETUP: {
    title: "No provider configured",
    hint: "Open Settings → Model Providers and choose a default provider.",
  },
  CURRENT_MODEL_NOT_SETUP: {
    title: "No model selected",
    hint: "Pick a model for this chat before sending a message.",
  },
  AGENT_EXCEED_MAX_ITERATIONS: {
    title: "Agent stopped",
    hint: "The agent looped too many times. Rephrase or split the request.",
  },
  AGENT_LOOP_ERROR: {
    title: "Agent error",
    hint: "The agent stopped before finishing the response. Try again.",
  },
  AGENT_CHAT_ERROR: {
    title: "Chat failed",
    hint: "The chat pipeline failed before the response could start.",
  },
  INVALID_TOOL_SCHEMA: {
    title: "Invalid tool configuration",
    hint: "A tool could not be sent to the model because its schema is invalid.",
  },
  INCORRECT_TOOL_ARGS: {
    title: "Invalid tool arguments",
    hint: "The model produced arguments a tool could not accept.",
  },
  GET_LIST_MODEL_FAILED: {
    title: "Could not load models",
    hint: "The provider model list could not be fetched. Check the provider configuration.",
  },
  UNKNOWN_ERROR: {
    title: "Something went wrong",
    hint: "The request failed unexpectedly. Try again.",
  },
};

export function describeAgentError(
  code: string,
  error: string,
): AgentErrorDescription {
  const meta = ERROR_META[code];
  const detail = error.trim();

  return {
    title: meta?.title ?? "Something went wrong",
    hint: meta?.hint,
    detail: detail && detail !== meta?.title ? detail : undefined,
  };
}
