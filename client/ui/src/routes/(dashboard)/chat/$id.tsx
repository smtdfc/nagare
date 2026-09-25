import { createFileRoute, useParams } from "@tanstack/react-router";
import { Bubble, BubbleContent } from "@/components/ui/bubble";
import {
  MessageScroller,
  MessageScrollerButton,
  MessageScrollerContent,
  MessageScrollerItem,
  MessageScrollerProvider,
  MessageScrollerViewport,
} from "@/components/ui/message-scroller";
import { useChat } from "#/hooks/use-chat.ts";
import { useEffect, useState, useRef, useCallback } from "react";
import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { ChatService } from "@nagare-app/services";
import { type Message as ChatMessage } from "@nagare-app/messages";
import {
  getMessageRole,
  isAgentCompletedMessage,
  isTextMessage,
} from "#/lib/message";
import { Spinner } from "#/components/ui/spinner";

export const Route = createFileRoute("/(dashboard)/chat/$id")({
  component: RouteComponent,
});

function renderMessage(message: ChatMessage) {
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
            <Markdown
              remarkPlugins={[remarkGfm]}
              components={{
                table: ({ node, ...props }) => (
                  <div className="my-4 w-full overflow-x-auto rounded-lg border border-border">
                    <table
                      className="w-full border-collapse text-left text-sm"
                      {...props}
                    />
                  </div>
                ),
                thead: ({ node, ...props }) => (
                  <thead
                    className="bg-muted/50 font-semibold border-b border-border"
                    {...props}
                  />
                ),
                th: ({ node, ...props }) => (
                  <th
                    className="px-4 py-2 border-r border-border last:border-r-0"
                    {...props}
                  />
                ),
                td: ({ node, ...props }) => (
                  <td
                    className="px-4 py-2 border-t border-r border-border last:border-r-0 align-top"
                    {...props}
                  />
                ),

                br: () => <br className="my-1 block" />,
              }}
            >
              {message?.content}
            </Markdown>
          </div>
        )}
      </div>
    );
  }
  return null;
}

function RouteComponent() {
  const { id } = useParams({ from: "/(dashboard)/chat/$id" });
  const [messages, setMessages] = useState<ChatMessage[]>([]);

  const pendingMessages = useChat((c) => c.pendingMessages);
  const setPendingMessages = useChat((c) => c.setPendingMessages);
  const setChatSession = useChat((c) => c.setChatSession);
  const setIsConnected = useChat((c) => c.setIsConnected);
  const isConnected = useChat((c) => c.isConnected);
  const isProcessing = useChat((c) => c.isProcessing);
  const setIsProcessing = useChat((c) => c.setIsProcessing);

  const messagesRef = useRef(messages);
  messagesRef.current = messages;

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

  useEffect(() => {
    let isMounted = true;
    let unListener: (() => void) | undefined;

    const fetchData = async () => {
      try {
        const session = await ChatService.getChatSession(id);
        if (!isMounted) return;
        setChatSession(session);

        const history = await ChatService.getChatHistory(id);
        if (!isMounted) return;
        if (history) {
          setMessages((prev) => [...history, ...prev]);
        }

        unListener = await ChatService.listenChatSession(id, handleMessage);
        if (!isMounted) {
          unListener?.();
          return;
        }
        setIsConnected(true);
      } catch (error) {
        console.error("Error:", error);
      }
    };

    fetchData();

    return () => {
      isMounted = false;
      unListener?.();
      setIsConnected(false);
    };
  }, [id, setChatSession, setIsConnected, handleMessage]);

  useEffect(() => {
    const processPending = async () => {
      if (pendingMessages && pendingMessages.length > 0 && isConnected) {
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
    <div className="flex flex-col h-[calc(100vh-4rem)] w-full overflow-hidden">
      <MessageScrollerProvider>
        <MessageScroller className="flex-1 overflow-hidden p-0">
          <MessageScrollerViewport className="h-full overflow-y-auto">
            <MessageScrollerContent className="flex flex-col max-w-3xl mx-auto w-full px-4 py-6">
              {messages.map((message) => (
                <MessageScrollerItem key={message.id} messageId={message.id}>
                  {renderMessage(message)}
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
