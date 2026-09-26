import { createFileRoute, useParams } from "@tanstack/react-router";
import {
  MessageScroller,
  MessageScrollerButton,
  MessageScrollerContent,
  MessageScrollerItem,
  MessageScrollerProvider,
  MessageScrollerViewport,
} from "@/components/ui/message-scroller";
import { useChat } from "#/hooks/use-chat.ts";
import { type UIEvent, useEffect, useState, useRef, useCallback } from "react";

import { ChatService } from "@nagare-app/services";
import { type Message as ChatMessage } from "@nagare-app/messages";
import {
  getMessageRole,
  isAgentCompletedMessage,
  isTextMessage,
} from "#/lib/message";
import { Spinner } from "#/components/ui/spinner";
import ChatMessageItem from "#/components/chat-message-item.tsx";
import { Button } from "@/components/ui/button";
import { HistoryIcon } from "lucide-react";

const HISTORY_PAGE_SIZE = 20;

export const Route = createFileRoute("/(dashboard)/chat/$id")({
  component: RouteComponent,
});

function RouteComponent() {
  const { id } = useParams({ from: "/(dashboard)/chat/$id" });
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [isLoadingHistory, setIsLoadingHistory] = useState(false);
  const [canLoadOlderHistory, setCanLoadOlderHistory] = useState(false);

  const pendingMessages = useChat((c) => c.pendingMessages);
  const setPendingMessages = useChat((c) => c.setPendingMessages);
  const setChatSession = useChat((c) => c.setChatSession);
  const setIsConnected = useChat((c) => c.setIsConnected);
  const isConnected = useChat((c) => c.isConnected);
  const isProcessing = useChat((c) => c.isProcessing);
  const setIsProcessing = useChat((c) => c.setIsProcessing);

  const messagesRef = useRef(messages);
  messagesRef.current = messages;
  const historyCursorRef = useRef<string | null>("");
  const historySessionRef = useRef(id);
  const isLoadingHistoryRef = useRef(false);
  const isSessionReadyRef = useRef(false);
  const viewportRef = useRef<HTMLDivElement>(null);

  const handleMessage = useCallback(async (chunk: ChatMessage) => {
    setMessages((prev) => {
      const lastMsg = prev[prev.length - 1];
      const isChunkText = isTextMessage(chunk);
      const isLastText = lastMsg && isTextMessage(lastMsg);

      if (isChunkText && isLastText) {
        const lastRole = getMessageRole(lastMsg);
        const chunkRole = getMessageRole(chunk);
        if (lastRole !== chunkRole) {
          return [...prev, chunk];
        }

        if (lastMsg.type === chunk.type && lastRole === chunkRole) {
          const updated = [...prev];
          updated[updated.length - 1] = {
            ...lastMsg,
            content: lastMsg.content + chunk.content,
          };
          return updated;
        }
      }

      if (isAgentCompletedMessage(chunk)) {
        setIsProcessing(false);
      }

      return prev;
    });
  }, []);

  const loadOlderHistory = useCallback(async () => {
    if (
      isLoadingHistoryRef.current ||
      historyCursorRef.current === null ||
      historySessionRef.current !== id
    ) {
      return;
    }

    const viewport = viewportRef.current;
    if (!viewport) return;

    isLoadingHistoryRef.current = true;
    setIsLoadingHistory(true);
    const previousScrollHeight = viewport.scrollHeight;
    const previousScrollTop = viewport.scrollTop;

    try {
      const history = await ChatService.getChatHistory(
        id,
        HISTORY_PAGE_SIZE,
        historyCursorRef.current,
      );
      if (historySessionRef.current !== id) return;
      if (!history) return;

      if (history.messages.length === 0) {
        historyCursorRef.current = null;
        setCanLoadOlderHistory(false);
        return;
      }

      historyCursorRef.current = history.nextCursor;
      setCanLoadOlderHistory(Boolean(history.nextCursor));
      setMessages((prev) => {
        const existingMessageIds = new Set(prev.map((message) => message.id));
        const olderMessages = history.messages.filter(
          (message) => !existingMessageIds.has(message.id),
        );

        return [...olderMessages, ...prev];
      });
      requestAnimationFrame(() => {
        const currentViewport = viewportRef.current;
        if (!currentViewport) return;
        currentViewport.scrollTop =
          previousScrollTop +
          (currentViewport.scrollHeight - previousScrollHeight);
      });
    } catch (error) {
      console.error("Error loading older chat history:", error);
    } finally {
      if (historySessionRef.current === id) {
        isLoadingHistoryRef.current = false;
        setIsLoadingHistory(false);
      }
    }
  }, [id]);

  const handleViewportScroll = useCallback(
    (event: UIEvent<HTMLDivElement>) => {
      if (event.currentTarget.scrollTop <= 80) {
        void loadOlderHistory();
      }
    },
    [loadOlderHistory],
  );

  useEffect(() => {
    const viewport = viewportRef.current;
    if (
      !viewport ||
      messages.length === 0 ||
      historyCursorRef.current === null ||
      isLoadingHistoryRef.current
    ) {
      return;
    }

    if (viewport.scrollHeight <= viewport.clientHeight) {
      void loadOlderHistory();
    }
  }, [id, loadOlderHistory, messages.length]);

  useEffect(() => {
    let isMounted = true;
    let unListener: (() => void) | undefined;

    setMessages([]);
    setCanLoadOlderHistory(false);
    setChatSession(null);
    setIsConnected(false);
    setIsProcessing(false);
    isLoadingHistoryRef.current = false;
    isSessionReadyRef.current = false;
    setIsLoadingHistory(false);
    historyCursorRef.current = "";
    historySessionRef.current = id;
    viewportRef.current?.scrollTo({ top: 0 });

    const fetchData = async () => {
      try {
        const session = await ChatService.getChatSession(id);
        if (!isMounted) return;
        setChatSession(session);

        const history = await ChatService.getChatHistory(
          id,
          HISTORY_PAGE_SIZE,
          "",
        );
        if (!isMounted) return;
        historyCursorRef.current = history?.messages.length
          ? history.nextCursor
          : null;
        setCanLoadOlderHistory(
          Boolean(history?.messages.length && history.nextCursor),
        );
        setMessages(history?.messages ?? []);
        requestAnimationFrame(() => {
          const viewport = viewportRef.current;
          if (!viewport || historySessionRef.current !== id) return;
          viewport.scrollTop = viewport.scrollHeight;
        });

        unListener = await ChatService.listenChatSession(id, handleMessage);
        if (!isMounted) {
          unListener?.();
          return;
        }
        isSessionReadyRef.current = true;
        setIsConnected(true);
      } catch (error) {
        console.error("Error:", error);
      }
    };

    fetchData();

    return () => {
      isMounted = false;
      isSessionReadyRef.current = false;
      unListener?.();
      setIsConnected(false);
    };
  }, [id, setChatSession, setIsConnected, handleMessage]);

  useEffect(() => {
    const processPending = async () => {
      if (
        pendingMessages &&
        pendingMessages.length > 0 &&
        isConnected &&
        isSessionReadyRef.current
      ) {
        const firstMessage = pendingMessages[0];

        if (firstMessage && isTextMessage(firstMessage)) {
          setMessages((prev) => [...prev, firstMessage]);
          await ChatService.sendMessage(id, firstMessage.content);
          setIsProcessing(true);
          setPendingMessages(pendingMessages.slice(1));
        }
      }
    };

    processPending();
  }, [pendingMessages, setPendingMessages, isConnected, id]);

  return (
    <div className="flex flex-col h-full min-h-0 w-full overflow-hidden">
      <MessageScrollerProvider key={id}>
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

              {messages.map((message) => (
                <MessageScrollerItem key={message.id} messageId={message.id}>
                  <ChatMessageItem message={message} />
                </MessageScrollerItem>
              ))}

              {isProcessing && (
                <div className="flex w-full justify-start">
                  <Spinner />
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
