import {
  MessageScroller,
  MessageScrollerButton,
  MessageScrollerContent,
  MessageScrollerItem,
  MessageScrollerProvider,
  MessageScrollerViewport,
} from "@/components/ui/message-scroller";
import { Button } from "@/components/ui/button";
import { Spinner } from "#/components/ui/spinner";
import { HistoryIcon } from "lucide-react";
import ChatMessageItem from "#/components/chat-message-item.tsx";
import { ChatEmptyState } from "#/components/chat-empty-stare";
import type { ChatConversationController } from "#/hooks/use-chat-conversation.ts";
import { isVisibleChatMessage } from "#/lib/message";
import type { Message as ChatMessage } from "@nagare-app/messages";

type ChatMessageListProps = {
  sessionID: string;
  conversation: ChatConversationController;
};

function renderMessageContent(message: ChatMessage) {
  if (!isVisibleChatMessage(message)) {
    return null;
  }

  return (
    <MessageScrollerItem key={message.id} messageId={message.id}>
      <ChatMessageItem message={message} />
    </MessageScrollerItem>
  );
}

export function ChatMessageList({
  sessionID,
  conversation,
}: ChatMessageListProps) {
  const {
    messages,
    statusText,
    isProcessing,
    isInitialLoading,
    isLoadingHistory,
    canLoadOlderHistory,
    visibleMessageCount,
    viewportRef,
    handleViewportScroll,
    loadOlderHistory,
  } = conversation;

  return (
    <div className="flex flex-col h-full min-h-0 w-full overflow-hidden">
      <MessageScrollerProvider key={sessionID}>
        <MessageScroller className="flex-1 overflow-hidden p-0">
          {isLoadingHistory && (
            <div className="pointer-events-none absolute inset-x-0 top-2 z-10 flex justify-center">
              <div className="rounded-full border bg-background/95 p-2 shadow-sm">
                <Spinner />
              </div>
            </div>
          )}
          <MessageScrollerViewport
            ref={viewportRef}
            onScroll={handleViewportScroll}
            className="h-full overflow-y-auto"
          >
            <MessageScrollerContent className="flex flex-col max-w-3xl mx-auto w-full px-4 py-6">
              {canLoadOlderHistory && (
                <div className="flex justify-center pb-2">
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    disabled={isLoadingHistory}
                    onClick={() => void loadOlderHistory()}
                  >
                    <HistoryIcon />
                    Load older messages
                  </Button>
                </div>
              )}

              {!isInitialLoading &&
                !isProcessing &&
                visibleMessageCount === 0 && <ChatEmptyState />}

              {messages.map((message) => renderMessageContent(message))}
              {isProcessing && (
                <div className="flex w-full justify-start items-center gap-1.5">
                  <Spinner />
                  <span className="text-sm">{statusText}</span>
                </div>
              )}
            </MessageScrollerContent>
          </MessageScrollerViewport>
          <MessageScrollerButton />
        </MessageScroller>
      </MessageScrollerProvider>
    </div>
  );
}
