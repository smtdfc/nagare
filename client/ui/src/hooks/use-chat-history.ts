import {
  type Dispatch,
  type RefObject,
  type SetStateAction,
  type UIEvent,
  useCallback,
  useEffect,
  useRef,
  useState,
} from "react";
import { ChatService } from "@nagare-app/services";
import type { Message as ChatMessage } from "@nagare-app/messages";
import { isVisibleChatMessage } from "#/lib/message";
import { showToastError } from "#/lib/toast";

const HISTORY_PAGE_SIZE = 20;

export type ChatHistoryController = {
  messages: ChatMessage[];
  setMessages: Dispatch<SetStateAction<ChatMessage[]>>;
  isLoadingHistory: boolean;
  canLoadOlderHistory: boolean;
  visibleMessageCount: number;
  viewportRef: RefObject<HTMLDivElement | null>;
  handleViewportScroll: (event: UIEvent<HTMLDivElement>) => void;
  loadOlderHistory: () => Promise<void>;
  loadInitialHistory: () => Promise<void>;
  resetHistory: () => void;
};

export function useChatHistory(id: string): ChatHistoryController {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [isLoadingHistory, setIsLoadingHistory] = useState(false);
  const [canLoadOlderHistory, setCanLoadOlderHistory] = useState(false);

  const messagesRef = useRef(messages);
  messagesRef.current = messages;
  const historyCursorRef = useRef<string | null>("");
  const historySessionRef = useRef(id);
  const isLoadingHistoryRef = useRef(false);
  const viewportRef = useRef<HTMLDivElement>(null);

  const resetHistory = useCallback(() => {
    setMessages([]);
    setCanLoadOlderHistory(false);
    isLoadingHistoryRef.current = false;
    setIsLoadingHistory(false);
    historyCursorRef.current = "";
    historySessionRef.current = id;
    viewportRef.current?.scrollTo({ top: 0 });
  }, [id]);

  const loadInitialHistory = useCallback(async () => {
    const history = await ChatService.getChatHistory(id, HISTORY_PAGE_SIZE, "");
    if (historySessionRef.current !== id) return;

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
  }, [id]);

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

  const visibleMessageCount = messages.filter(isVisibleChatMessage).length;

  return {
    messages,
    setMessages,
    isLoadingHistory,
    canLoadOlderHistory,
    visibleMessageCount,
    viewportRef,
    handleViewportScroll,
    loadOlderHistory,
    loadInitialHistory,
    resetHistory,
  };
}
