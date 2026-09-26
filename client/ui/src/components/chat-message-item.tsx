import { getMessageRole, isTextMessage } from "#/lib/message.ts";
import { Bubble, BubbleContent } from "#/components/ui/bubble.tsx";
import MarkdownDisplay from "#/components/markdown-display.tsx";
import { type Message as ChatMessage } from "@nagare-app/messages";

type ChatMessageItemProps = {
  message: ChatMessage;
};

export default function ChatMessageItem({ message }: ChatMessageItemProps) {
  const role = getMessageRole(message);
  const isUser = role === "user" || role === "client";

  if (isTextMessage(message)) {
    return (
      <div
        className={`flex w-full  ${isUser ? "justify-end" : "justify-start"}`}
      >
        {isUser ? (
          <Bubble
            variant="default"
            className="max-w-[75%] p-3 text-sm rounded-2xl bg-primary text-primary-foreground rounded-br-xs shadow-md whitespace-pre-wrap"
          >
            <BubbleContent>{message?.content}</BubbleContent>
          </Bubble>
        ) : (
          <div className="max-w-[85%] px-1 py-2.5 text-sm text-foreground leading-relaxed prose dark:prose-invert">
            <MarkdownDisplay content={message.content} />
          </div>
        )}
      </div>
    );
  }
  return null;
}
