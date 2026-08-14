import type {
  AgentResponse,
  ResponseCompleted,
  ResponseFailed,
  ResponseStarted,
  Text,
  ToolCall,
  ToolCallResult,
} from "./messages";

export * from "./api";
export * from "./messages";

declare global {
  type Message =
    | Text
    | AgentResponse
    | ResponseStarted
    | ResponseCompleted
    | ResponseFailed
    | ToolCall
    | ToolCallResult;
}
