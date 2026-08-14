import type {
  Text,
  AgentResponse,
  ResponseStarted,
  ResponseCompleted,
  ResponseFailed,
  ToolCall,
  ToolCallResult,
} from '@/dto/messages'

export {}

declare global {
  type Message =
    | Text
    | AgentResponse
    | ResponseStarted
    | ResponseCompleted
    | ResponseFailed
    | ToolCall
    | ToolCallResult
}
