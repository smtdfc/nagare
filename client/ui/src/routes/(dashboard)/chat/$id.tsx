import { createFileRoute, useParams } from "@tanstack/react-router";
import { ChatMessageList } from "#/components/chat-message-list";
import { useChatConversation } from "#/hooks/use-chat-conversation.ts";

export const Route = createFileRoute("/(dashboard)/chat/$id")({
  component: RouteComponent,
});

function RouteComponent() {
  const { id } = useParams({ from: "/(dashboard)/chat/$id" });
  const conversation = useChatConversation(id);

  return <ChatMessageList sessionID={id} conversation={conversation} />;
}
