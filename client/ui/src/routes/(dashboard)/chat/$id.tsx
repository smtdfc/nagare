import { createFileRoute, useParams } from "@tanstack/react-router";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Bubble, BubbleContent } from "@/components/ui/bubble";
import {
  Message,
  MessageAvatar,
  MessageContent,
} from "@/components/ui/message";
import { useChat } from "#/hooks/use-chat.ts";
import { useEffect, useState } from "react";
import { ChatService } from "@nagare-app/services";
import {
  type Message as ChatMessage,
  MessageType,
  type TextMessage,
} from "@nagare-app/messages";

export const Route = createFileRoute("/(dashboard)/chat/$id")({
  component: RouteComponent,
});

function isTextMessage(message: ChatMessage): message is TextMessage {
  return message.type === MessageType.TextMessageType;
}
function renderMessage(message: ChatMessage) {
  if (isTextMessage(message)) {
    return (
      <Message>
        <MessageAvatar>
          <Avatar>
            <AvatarImage src="https://github.com/shadcn.png" alt="@shadcn" />
            <AvatarFallback>CN</AvatarFallback>
          </Avatar>
        </MessageAvatar>
        <MessageContent>
          <Bubble>
            <BubbleContent>{message?.content}</BubbleContent>
          </Bubble>
        </MessageContent>
      </Message>
    );
  }

  return "";
}

function RouteComponent() {
  const { id } = useParams({ from: "/(dashboard)/chat/$id" });
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const pendingMessages = useChat((c) => c.pendingMessages);
  const setPendingMessages = useChat((c) => c.setPendingMessages);
  const setChatSession = useChat((c) => c.setChatSession);

  useEffect(() => {
    const fetchData = async () => {
      const session = await ChatService.getChatSession(id);
      setChatSession(session);
    };

    fetchData();
  }, [id, setChatSession]);

  useEffect(() => {
    if (pendingMessages && pendingMessages.length > 0) {
      setMessages((prev) => [...prev, ...pendingMessages]);
    }
  }, [pendingMessages, setPendingMessages]);

  return <>{messages.map((message) => renderMessage(message))}</>;
}
