import {
  MessageType,
  type AgentCompletedMessage,
  type AgentErrorMessage,
  type ApprovalRequestMessage,
  type CompactSuccessMessage,
  type ConversationSummaryMessage,
  type CronJobTriggerMessage,
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

export function isApprovalRequestMessage(
  message: ChatMessage,
): message is ApprovalRequestMessage {
  return message.type === MessageType.ApprovalRequestMessageType;
}

export function isConversationSummaryMessage(
  message: ChatMessage,
): message is ConversationSummaryMessage {
  return message.type === MessageType.ConversationSummaryMessageType;
}

export function isCronJobTriggerMessage(
  message: ChatMessage,
): message is CronJobTriggerMessage {
  return message.type === MessageType.CronJobTriggerMessageType;
}

export function isCompactSuccessMessage(
  message: ChatMessage,
): message is CompactSuccessMessage {
  return message.type === MessageType.CompactSuccessMessageType;
}

export function getMessageRole(message: ChatMessage) {
  if (isTextMessage(message)) {
    return (message.role as string).toLowerCase();
  }

  if (isAgentCompletedMessage(message) || isAgentErrorMessage(message)) {
    return "system";
  }

  if (isToolCallMessage(message)) {
    return "tool";
  }

  if (isApprovalRequestMessage(message)) {
    return "system";
  }

  return "agent";
}
