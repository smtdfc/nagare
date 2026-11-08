import { MessageSquareIcon } from "lucide-react";

export function ChatEmptyState() {
  return (
    <div className="flex min-h-[50vh] flex-col items-center justify-center gap-2 text-center">
      <MessageSquareIcon className="size-8 text-muted-foreground/60" />
      <p className="text-sm font-medium">No messages yet</p>
      <p className="max-w-sm text-sm text-muted-foreground">
        Send a message below to start the conversation.
      </p>
    </div>
  );
}
