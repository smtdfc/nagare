import {
  MessageType,
  type AgentCompletedMessage,
  type TextMessage,
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

export function getMessageRole(message: ChatMessage) {
  if (isTextMessage(message)) {
    return (message.role as string).toLowerCase();
  }

  return "agent";
}
