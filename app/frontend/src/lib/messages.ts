import { MessageType } from "#/dto/messages.ts";
import type { AgentResponse, Text, ToolCall } from "#/dto/messages.ts";

export function isTextMessage(message: Message): message is Text {
    return message.type === MessageType.TEXT_MESSAGE;
}

export function isAgentResponseMessage(message: Message): message is AgentResponse {
    return message.type === MessageType.AGENT_RESPONSE;
}

export function isToolCallMessage(message: Message): message is ToolCall {
    return message.type === MessageType.TOOL_CALL_MESSAGE;
}