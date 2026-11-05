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
import {
  MessageType,
  Role,
  type Message as ChatMessage,
} from "@nagare-app/messages";
import {
  getMessageRole,
  isAgentCompletedMessage,
  isAgentErrorMessage,
  isApprovalRequestMessage,
  isConversationSummaryMessage,
  isTextMessage,
  isToolCallMessage,
} from "#/lib/message";
import { Spinner } from "#/components/ui/spinner";
import ChatMessageItem from "#/components/chat-message-item.tsx";
import { Button } from "@/components/ui/button";
import { HistoryIcon, MessageSquareIcon } from "lucide-react";
import { showToastError } from "#/lib/toast";
import { toast } from "#/components/ui/toast.tsx";
import {
  clearLiveChatSnapshot,
  loadLiveChatSnapshot,
  saveLiveChatSnapshot,
  type LiveChatSnapshot,
} from "#/lib/live-chat";

const HISTORY_PAGE_SIZE = 20;
// Polling is the safety net for a reload that happens while a turn finishes:
// the completion event was missed, so the DB is the only source of truth.
const LIVE_CHAT_POLL_INTERVAL_MS = 4_000;
const LIVE_CHAT_POLL_TIMEOUT_MS = 3 * 60 * 1000;
// Writing sessionStorage on every streamed chunk is O(n^2) for long replies.
const LIVE_CHAT_SAVE_THROTTLE_MS = 250;

export const Route = createFileRoute("/(dashboard)/chat/$id")({
  component: RouteComponent,
});

const allowShowMessageTypes = [
  MessageType.TextMessageType,
  MessageType.ApprovalRequestMessageType,
  MessageType.ConversationSummaryMessageType,
  MessageType.AgentErrorMessageType,
];
const allowHandleMessageTypes = [
  MessageType.TextMessageType,
  MessageType.ToolCallMessageType,
  MessageType.AgentCompletedMessageType,
  MessageType.AgentErrorMessageType,
  MessageType.ApprovalRequestMessageType,
  MessageType.ConversationSummaryMessageType,
  MessageType.CompactSuccessMessageType,
];

function renderMessageContent(message: ChatMessage) {
  if (!allowShowMessageTypes.includes(message.type)) {
    return null;
  }

  return (
    <MessageScrollerItem key={message.id} messageId={message.id}>
      <ChatMessageItem message={message} />
    </MessageScrollerItem>
  );
}

function lastIndexWhere<T>(
  items: T[],
  predicate: (item: T) => boolean,
): number {
  for (let index = items.length - 1; index >= 0; index -= 1) {
    const item = items[index];
    if (item !== undefined && predicate(item)) {
      return index;
    }
  }

  return -1;
}

// The snapshot is only useful while the turn is still missing from the DB. If
// the turn was persisted before the reload, restoring it would duplicate it.
function isSnapshotPersisted(
  snapshot: LiveChatSnapshot,
  history: ChatMessage[],
): boolean {
  const userText = snapshot.userText.trim();
  if (!userText) return false;

  const lastUserIndex = lastIndexWhere(
    history,
    (message) => isTextMessage(message) && getMessageRole(message) === "user",
  );
  const userIndex = lastIndexWhere(
    history,
    (message) =>
      isTextMessage(message) &&
      getMessageRole(message) === "user" &&
      message.content.trim() === userText,
  );

  if (userIndex === -1 || userIndex !== lastUserIndex) {
    return false;
  }

  const assistantText = snapshot.assistantText.trim();
  if (!assistantText) {
    return true;
  }

  const persistedText = history
    .slice(userIndex + 1)
    .filter(isTextMessage)
    .filter((message) => getMessageRole(message) !== "user")
    .map((message) => message.content)
    .join("\n");

  // Whitespace can differ between streamed chunks and persisted messages
  // (tool calls split a reply into several rows), so compare ignoring it.
  const normalize = (value: string) => value.replace(/\s+/g, "");
  return normalize(persistedText).includes(normalize(assistantText));
}

function ChatEmptyState() {
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

function RouteComponent() {
  const { id } = useParams({ from: "/(dashboard)/chat/$id" });
  const [statusText, setStatusText] = useState("Processing...");
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [isLoadingHistory, setIsLoadingHistory] = useState(false);
  const [isInitialLoading, setIsInitialLoading] = useState(true);
  const [canLoadOlderHistory, setCanLoadOlderHistory] = useState(false);
  const [restoredSnapshot, setRestoredSnapshot] =
    useState<LiveChatSnapshot | null>(null);

  const pendingMessages = useChat((c) => c.pendingMessages);
  const setPendingMessages = useChat((c) => c.setPendingMessages);
  const setChatSession = useChat((c) => c.setChatSession);
  const setIsConnected = useChat((c) => c.setIsConnected);
  const setInvokeID = useChat((c) => c.setInvokeID);
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
  const liveSaveTimerRef = useRef<number | undefined>(undefined);
  const liveSnapshotRef = useRef<LiveChatSnapshot | null>(null);

  const handleMessage = useCallback(async (chunk: ChatMessage) => {
    if (!allowHandleMessageTypes.includes(chunk.type)) {
      return;
    }

    // Compare the type directly: negating a bare {id,type,invokeID}
    // guard would exclude every other message from the union.
    if (chunk.type === MessageType.CompactSuccessMessageType) {
      toast.add({
        title: "Conversation compacted",
        type: "success",
        description: "Your conversation history has been summarized.",
      });
      return;
    }

    setMessages((prev) => {
      const lastMsg = prev[prev.length - 1];
      const isChunkText = isTextMessage(chunk);
      const isLastText = lastMsg && isTextMessage(lastMsg);

      if (isChunkText && isLastText) {
        setStatusText(`Responding...`);
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

      // A tool/approval event may be the latest visible item. The next text
      // chunk must start a new message instead of being dropped.
      if (isChunkText) {
        setStatusText("Responding...");
        return [...prev, chunk];
      }

      if (isAgentCompletedMessage(chunk)) {
        setStatusText(`Finishing ...`);
        setIsProcessing(false);
        // Keep the snapshot until the poll sees the saved turn. This covers a
        // reload in the short window between completion and SaveHistory.
        return prev;
      }

      if (isToolCallMessage(chunk)) {
        setStatusText(`Using ${chunk.name}`);
        return prev;
      }

      if (isApprovalRequestMessage(chunk)) {
        setStatusText("Waiting for approval...");
        return [...prev, chunk];
      }

      if (isAgentErrorMessage(chunk)) {
        setIsProcessing(false);
        setStatusText("");
        if (liveSaveTimerRef.current !== undefined) {
          window.clearTimeout(liveSaveTimerRef.current);
          liveSaveTimerRef.current = undefined;
        }
        liveSnapshotRef.current = null;
        clearLiveChatSnapshot(historySessionRef.current);
        setRestoredSnapshot(null);
        return [...prev, chunk];
      }

      if (isConversationSummaryMessage(chunk)) {
        return [...prev, chunk];
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
      if (!history || history.messages.length === 0) {
        historyCursorRef.current = null;
        setCanLoadOlderHistory(false);
        return;
      }

      historyCursorRef.current = history.nextCursor;
      setCanLoadOlderHistory(Boolean(history.nextCursor));

      const existingMessageIds = new Set(
        messagesRef.current.map((message) => message.id),
      );
      const olderMessages = history.messages.filter(
        (message) => !existingMessageIds.has(message.id),
      );

      setMessages((prev) => [...olderMessages, ...prev]);
      requestAnimationFrame(() => {
        const currentViewport = viewportRef.current;
        if (!currentViewport) return;
        currentViewport.scrollTop =
          previousScrollTop +
          (currentViewport.scrollHeight - previousScrollHeight);
      });
    } catch (error) {
      showToastError(error);
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
    setIsInitialLoading(true);
    setRestoredSnapshot(null);
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

        const savedMessages = history?.messages ?? [];
        setMessages(savedMessages);

        const snapshot = loadLiveChatSnapshot(id);
        if (snapshot && !isSnapshotPersisted(snapshot, savedMessages)) {
          const restored: ChatMessage[] = [
            {
              id: `live-user-${id}`,
              type: MessageType.TextMessageType,
              role: Role.USER,
              invokeID: "",
              content: snapshot.userText,
            },
          ];

          if (snapshot.assistantText.trim()) {
            restored.push({
              id: `live-assistant-${id}`,
              type: MessageType.TextMessageType,
              role: Role.AGENT,
              invokeID: "",
              content: snapshot.assistantText,
            });
          }

          setMessages([...savedMessages, ...restored]);
          setStatusText(snapshot.status || "Responding...");
          setIsProcessing(true);
          setRestoredSnapshot(snapshot);
        } else if (snapshot) {
          clearLiveChatSnapshot(id);
        }

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
        showToastError(error);
      } finally {
        if (isMounted) {
          setIsInitialLoading(false);
        }
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

  // Keep a tab-local copy of the in-progress turn so a reload can restore it.
  useEffect(() => {
    if (!isProcessing) return;

    const lastUserIndex = lastIndexWhere(
      messages,
      (message) => isTextMessage(message) && getMessageRole(message) === "user",
    );
    if (lastUserIndex === -1) return;

    const userMessage = messages[lastUserIndex];
    if (!userMessage || !isTextMessage(userMessage)) return;

    const assistantText = messages
      .slice(lastUserIndex + 1)
      .filter(isTextMessage)
      .filter((message) => getMessageRole(message) !== "user")
      .map((message) => message.content)
      .join("");

    liveSnapshotRef.current = {
      sessionID: id,
      userText: userMessage.content,
      assistantText,
      status: statusText,
      updatedAt: Date.now(),
    };

    if (liveSaveTimerRef.current !== undefined) return;

    liveSaveTimerRef.current = window.setTimeout(() => {
      liveSaveTimerRef.current = undefined;
      const pending = liveSnapshotRef.current;
      if (pending) {
        saveLiveChatSnapshot(pending);
      }
    }, LIVE_CHAT_SAVE_THROTTLE_MS);
  }, [id, messages, isProcessing, statusText]);

  // Flush the last snapshot when leaving the chat or switching sessions so the
  // tail of an in-progress reply is not lost.
  useEffect(() => {
    return () => {
      if (liveSaveTimerRef.current !== undefined) {
        window.clearTimeout(liveSaveTimerRef.current);
        liveSaveTimerRef.current = undefined;
      }

      if (liveSnapshotRef.current) {
        saveLiveChatSnapshot(liveSnapshotRef.current);
        liveSnapshotRef.current = null;
      }
    };
  }, [id]);

  // If the turn finished while the tab was reloading, the completion event is
  // gone. Poll until the saved turn appears, then replace the optimistic
  // messages with the persisted ones.
  useEffect(() => {
    if (!restoredSnapshot) return;

    let cancelled = false;
    let timer: number | undefined;
    const startedAt = Date.now();

    const reconcile = async () => {
      if (cancelled) return;

      try {
        const history = await ChatService.getChatHistory(
          id,
          HISTORY_PAGE_SIZE,
          "",
        );
        if (cancelled || historySessionRef.current !== id) return;

        const latest = history?.messages ?? [];
        if (isSnapshotPersisted(restoredSnapshot, latest)) {
          if (liveSaveTimerRef.current !== undefined) {
            window.clearTimeout(liveSaveTimerRef.current);
            liveSaveTimerRef.current = undefined;
          }
          liveSnapshotRef.current = null;
          clearLiveChatSnapshot(id);
          setRestoredSnapshot(null);
          setMessages(latest);
          historyCursorRef.current = latest.length ? history.nextCursor : null;
          setCanLoadOlderHistory(Boolean(latest.length && history.nextCursor));
          setIsProcessing(false);
          requestAnimationFrame(() => {
            const viewport = viewportRef.current;
            if (viewport) viewport.scrollTop = viewport.scrollHeight;
          });
          return;
        }
      } catch {
        // A failed poll is harmless: keep the restored bubble and retry.
      }

      if (!cancelled && Date.now() - startedAt < LIVE_CHAT_POLL_TIMEOUT_MS) {
        timer = window.setTimeout(
          () => void reconcile(),
          LIVE_CHAT_POLL_INTERVAL_MS,
        );
      }
    };

    timer = window.setTimeout(
      () => void reconcile(),
      LIVE_CHAT_POLL_INTERVAL_MS,
    );

    return () => {
      cancelled = true;
      if (timer !== undefined) {
        window.clearTimeout(timer);
      }
    };
  }, [id, restoredSnapshot]);

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
          setStatusText("Processing...");
          const response = await ChatService.sendMessage(
            id,
            firstMessage.content,
          );
          if (response?.invokeID) {
            setMessages((prev) =>
              prev.map((message) =>
                message.id === firstMessage.id
                  ? { ...message, invokeID: response.invokeID }
                  : message,
              ),
            );
            setInvokeID(response.invokeID);
          }
          setIsProcessing(true);
          setPendingMessages(pendingMessages.slice(1));
        }
      }
    };

    processPending();
  }, [pendingMessages, setPendingMessages, isConnected, id]);

  const visibleMessageCount = messages.filter((message) =>
    allowShowMessageTypes.includes(message.type),
  ).length;

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
