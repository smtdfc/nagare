import {
  MessageType,
  type AgentCompletedMessage,
  type AgentErrorMessage,
  type TextMessage,
  type ToolCallMessage,
} from "@nagare-app/messages";
import { type Message as ChatMessage } from "@nagare-app/messages";

export function isTextMessage(message: ChatMessage): message is TextMessage {
  return message.type === MessageType.TextMessageType;
}

export function isAgentCompletedMessage(
  message: ChatMessage,
): message is AgentCompletedMessage {
  return message.type === MessageType.AgentCompletedMessageType;
}

export function isToolCallMessage(
  message: ChatMessage,
): message is ToolCallMessage {
  return message.type === MessageType.ToolCallMessageType;
}

export function isAgentErrorMessage(
  message: ChatMessage,
): message is AgentErrorMessage {
  return message.type === MessageType.AgentErrorMessageType;
}

export function getMessageRole(message: ChatMessage) {
  if (isTextMessage(message)) {
    return (message.role as string).toLowerCase();
  }

  return "agent";
}
