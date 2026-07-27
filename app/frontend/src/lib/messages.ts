import { MessageType } from "#/dto/messages.ts";
import type { Text } from "#/dto/messages.ts";

export function isTextMessage(message: Message): message is Text {
    return message.id === MessageType.TEXT_MESSAGE;
}