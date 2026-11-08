import {
  type RefObject,
  type UIEvent,
  useEffect,
  useRef,
  useState,
} from "react";
import { ChatService } from "@nagare-app/services";
import type { Message as ChatMessage } from "@nagare-app/messages";
import { useChat } from "#/hooks/use-chat.ts";
import { useChatHistory } from "#/hooks/use-chat-history.ts";
import { useChatMessageStream } from "#/hooks/use-chat-message-stream.ts";
import { isTextMessage } from "#/lib/message";
import { showToastError } from "#/lib/toast";

export type ChatConversationController = {
  messages: ChatMessage[];
  statusText: string;
  isProcessing: boolean;
  isInitialLoading: boolean;
  isLoadingHistory: boolean;
  canLoadOlderHistory: boolean;
  visibleMessageCount: number;
  viewportRef: RefObject<HTMLDivElement | null>;
  handleViewportScroll: (event: UIEvent<HTMLDivElement>) => void;
  loadOlderHistory: () => Promise<void>;
};

export function useChatConversation(id: string): ChatConversationController {
  const [isInitialLoading, setIsInitialLoading] = useState(true);

  const pendingMessages = useChat((c) => c.pendingMessages);
  const setPendingMessages = useChat((c) => c.setPendingMessages);
  const setChatSession = useChat((c) => c.setChatSession);
  const setIsConnected = useChat((c) => c.setIsConnected);
  const setInvokeID = useChat((c) => c.setInvokeID);
  const isConnected = useChat((c) => c.isConnected);
  const isProcessing = useChat((c) => c.isProcessing);
  const setIsProcessing = useChat((c) => c.setIsProcessing);

  const history = useChatHistory(id);
  const { statusText, setStatusText, currentInvokeIDRef, handleMessage } =
    useChatMessageStream(history.setMessages);

  const isSessionReadyRef = useRef(false);

  useEffect(() => {
    let isMounted = true;
    let unListener: (() => void) | undefined;

    history.resetHistory();
    setChatSession(null);
    setIsConnected(false);
    setIsProcessing(false);
    setIsInitialLoading(true);
    currentInvokeIDRef.current = null;
    setInvokeID(null);
    isSessionReadyRef.current = false;

    const fetchData = async () => {
      try {
        const session = await ChatService.getChatSession(id);
        if (!isMounted) return;
        setChatSession(session);

        await history.loadInitialHistory();
        if (!isMounted) return;

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
      currentInvokeIDRef.current = null;
      setInvokeID(null);
    };
  }, [
    id,
    setChatSession,
    setIsConnected,
    history.resetHistory,
    history.loadInitialHistory,
    handleMessage,
    currentInvokeIDRef,
    setInvokeID,
    setIsProcessing,
  ]);

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
          history.setMessages((prev) => [...prev, firstMessage]);
          setStatusText("Processing...");
          const response = await ChatService.sendMessage(
            id,
            firstMessage.content,
          );
          if (response?.invokeID) {
            currentInvokeIDRef.current = response.invokeID;
            history.setMessages((prev) =>
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
  }, [
    pendingMessages,
    setPendingMessages,
    isConnected,
    id,
    history.setMessages,
    setStatusText,
    currentInvokeIDRef,
    setInvokeID,
    setIsProcessing,
  ]);

  return {
    messages: history.messages,
    statusText,
    isProcessing,
    isInitialLoading,
    isLoadingHistory: history.isLoadingHistory,
    canLoadOlderHistory: history.canLoadOlderHistory,
    visibleMessageCount: history.visibleMessageCount,
    viewportRef: history.viewportRef,
    handleViewportScroll: history.handleViewportScroll,
    loadOlderHistory: history.loadOlderHistory,
  };
}
