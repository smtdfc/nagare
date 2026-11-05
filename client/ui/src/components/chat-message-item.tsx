import {
  getMessageRole,
  isAgentErrorMessage,
  isApprovalRequestMessage,
  isConversationSummaryMessage,
  isTextMessage,
} from "#/lib/message.ts";
import { describeAgentError } from "#/lib/agent-error.ts";
import { Bubble, BubbleContent } from "#/components/ui/bubble.tsx";
import MarkdownDisplay from "#/components/markdown-display.tsx";
import {
  type AgentErrorMessage,
  type ApprovalRequestMessage,
  type Message as ChatMessage,
} from "@nagare-app/messages";
import { ChatService } from "@nagare-app/services";
import { useState } from "react";
import { Button } from "./ui/button";
import { Layers, TriangleAlertIcon } from "lucide-react";

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

  if (isConversationSummaryMessage(message)) {
    return <ConversationSummary message={message} />;
  }

  if (isAgentErrorMessage(message)) {
    return <AgentError message={message} />;
  }

  return null;
}

function AgentError({ message }: { message: AgentErrorMessage }) {
  const { title, hint, detail } = describeAgentError(
    message.code,
    message.error,
  );

  return (
    <div className="my-1 w-full rounded-lg border border-destructive/40 bg-destructive/10 p-3 text-sm">
      <div className="flex items-start gap-2">
        <TriangleAlertIcon className="mt-0.5 size-4 shrink-0 text-destructive" />
        <div className="min-w-0 flex-1">
          <p className="font-medium text-destructive">{title}</p>
          {detail && (
            <p className="mt-1 break-words text-foreground/80">{detail}</p>
          )}
          {hint && <p className="mt-1 text-xs text-muted-foreground">{hint}</p>}
        </div>
      </div>
    </div>
  );
}

function ToolApproval({ message }: { message: ApprovalRequestMessage }) {
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

function ConversationSummary() {
  return (
    <div className="my-6 w-full">
      <div className="relative flex items-center gap-3">
        <div className="h-px flex-1 bg-border/60" />
        <div className="flex shrink-0 items-center gap-1.5 rounded-full border border-border/60 bg-muted/50 px-3 py-1 text-xs font-medium text-muted-foreground">
          <Layers className="size-3" />
          <span>Context compacted</span>
        </div>
        <div className="h-px flex-1 bg-border/60" />
      </div>
    </div>
  );
}
