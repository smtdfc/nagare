import {
  type Dispatch,
  type RefObject,
  type SetStateAction,
  useCallback,
  useRef,
  useState,
} from "react";
import { MessageType, type Message as ChatMessage } from "@nagare-app/messages";
import { useChat } from "#/hooks/use-chat.ts";
import {
  getMessageRole,
  isAgentCompletedMessage,
  isAgentErrorMessage,
  isApprovalRequestMessage,
  isConversationSummaryMessage,
  isCronJobTriggerMessage,
  isHandledChatMessage,
  isTextMessage,
  isToolCallMessage,
} from "#/lib/message";
import { toast } from "#/components/ui/toast.tsx";

export type ChatMessageStreamController = {
  statusText: string;
  setStatusText: Dispatch<SetStateAction<string>>;
  currentInvokeIDRef: RefObject<string | null>;
  handleMessage: (chunk: ChatMessage) => Promise<void>;
};

export function useChatMessageStream(
  setMessages: Dispatch<SetStateAction<ChatMessage[]>>,
): ChatMessageStreamController {
  const [statusText, setStatusText] = useState("Processing...");
  const setInvokeID = useChat((c) => c.setInvokeID);
  const invokeID = useChat((c) => c.invokeID);
  const setIsProcessing = useChat((c) => c.setIsProcessing);

  const currentInvokeIDRef = useRef<string | null>(null);
  currentInvokeIDRef.current = invokeID;

  const handleMessage = useCallback(async (chunk: ChatMessage) => {
    if (!isHandledChatMessage(chunk)) {
      return;
    }

    if (
      !currentInvokeIDRef.current ||
      chunk.invokeID !== currentInvokeIDRef.current
    ) {
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
      setIsProcessing(false);
      currentInvokeIDRef.current = null;
      setInvokeID(null);
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
        currentInvokeIDRef.current = null;
        setInvokeID(null);
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
        currentInvokeIDRef.current = null;
        setInvokeID(null);
        return [...prev, chunk];
      }

      if (
        isConversationSummaryMessage(chunk) ||
        isCronJobTriggerMessage(chunk)
      ) {
        return [...prev, chunk];
      }

      return prev;
    });
  }, []);

  return {
    statusText,
    setStatusText,
    currentInvokeIDRef,
    handleMessage,
  };
}
