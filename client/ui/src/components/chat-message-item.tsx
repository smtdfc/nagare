import {
  getMessageRole,
  isApprovalRequestMessage,
  isTextMessage,
} from "#/lib/message.ts";
import { Bubble, BubbleContent } from "#/components/ui/bubble.tsx";
import MarkdownDisplay from "#/components/markdown-display.tsx";
import {
  type ApprovalRequestMessage,
  type Message as ChatMessage,
} from "@nagare-app/messages";
import { ChatService } from "@nagare-app/services";
import { useState } from "react";
import { Button } from "./ui/button";

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

  if (isApprovalRequestMessage(message)) {
    return <ToolApproval message={message} />;
  }
  return null;
}

function ToolApproval({
  message,
}: {
  message: ApprovalRequestMessage;
}) {
  const [decision, setDecision] = useState<boolean | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const submit = async (accept: boolean) => {
    setIsSubmitting(true);
    try {
      await ChatService.submitToolApproval({
        requestID: message.requestID,
        sessionID: message.sessionID,
        accept,
      });
      setDecision(accept);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="my-3 rounded-lg border border-amber-500/40 bg-amber-500/10 p-4 text-sm">
      <p className="font-medium">{message.title}</p>
      {decision === null ? (
        <div className="mt-3 flex gap-2">
          <Button
            disabled={isSubmitting}
            size="sm"
            onClick={() => void submit(true)}
          >
            Approve
          </Button>
          <Button
            disabled={isSubmitting}
            size="sm"
            variant="destructive"
            onClick={() => void submit(false)}
          >
            Reject
          </Button>
        </div>
      ) : (
        <p className="mt-3 text-muted-foreground">
          {decision ? "Approved" : "Rejected"}
        </p>
      )}
    </div>
  );
}
