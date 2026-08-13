import ChatInput from '#/components/chat-input.tsx';
import ChatMessage from '#/components/chat-message.tsx';
import WelcomeHeader from '#/components/welcome-header.tsx';
import { AgentResponseStatus } from '#/dto/messages.ts';
import { isAgentResponseMessage, isTextMessage, isToolCallMessage } from '#/lib/messages.ts';
import { initChatWebsocketConnection } from '#/lib/websocket.ts';
import { ChatService } from '#/services/chat.ts';
import { createFileRoute } from '@tanstack/react-router';
import { useState, useRef, useEffect, useCallback } from 'react';

export const Route = createFileRoute('/')({
  component: Home,
  staticData: {
    breadcrumb: 'Chat',
  },
})

interface MessageItem {
  id: string;
  role: 'user' | 'agent' | 'system';
  content: string;
}

enum Status {
  CONNECTING,
  IDLE,
  PENDING
}

function Home() {
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [status, setStatus] = useState<Status>(Status.CONNECTING);
  const [messages, setMessages] = useState<MessageItem[]>([]);
  const [prompt, setPrompt] = useState<string>("");
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const sessionIdRef = useRef(sessionId);
  sessionIdRef.current = sessionId;

  const addMessage = useCallback((message: MessageItem, forceNew: boolean = false) => {
    setMessages((prevMessages) => {
      if (prevMessages.length === 0) {
        return [message];
      }
      const lastMessage = prevMessages[prevMessages.length - 1];
      if (lastMessage.role !== message.role || forceNew) {
        return [...prevMessages, message];
      }

      return [
        ...prevMessages.slice(0, -1),
        {
          ...lastMessage,
          content: lastMessage.content + message.content,
        }
      ];
    });
  }, []);

  const handleSubmit = async () => {
    if (!prompt.trim() || status === Status.PENDING) return;

    const userMsgContent = prompt;
    setPrompt("");
    setStatus(Status.PENDING);

    const userMessageId = `user_${Date.now().toString()}`;
    addMessage({
      id: userMessageId,
      role: "user",
      content: userMsgContent
    });

    setTimeout(() => {
      messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, 50);

    try {
      let id = sessionIdRef.current;
      if (!id) {
        id = await ChatService.createChatSession();
        setSessionId(id);
      }

      await ChatService.sendMessage(id, userMsgContent);
    } catch (e) {
      console.log(e);
      setStatus(Status.IDLE);
    }
  };

  useEffect(() => {
    let isMounted = true;
    let unsubscribe: (() => void) | undefined;

    const setupChat = async () => {
      try {
        await initChatWebsocketConnection();

        if (!isMounted) return;
        setStatus(Status.IDLE);
        let mustNewMessage = false;

        unsubscribe = ChatService.listenMessage((_isSuccess: boolean, messageChunk?: Message, err?: string) => {
          if (!messageChunk) {
            addMessage({
              id: `agent_${Date.now().toString()}`,
              role: "agent",
              content: err || "Unknown error"
            }, mustNewMessage);
            return;
          }

          if (isTextMessage(messageChunk)) {
            addMessage({
              id: `agent_${Date.now().toString()}`,
              role: "agent",
              content: messageChunk.content
            }, mustNewMessage);
            mustNewMessage = false;
          }

          if (isToolCallMessage(messageChunk)) {
            addMessage({
              id: `agent_${Date.now().toString()}`,
              role: "agent",
              content: `Call: ${messageChunk.name}`
            }, true);
            mustNewMessage = true;
          }

          if (isAgentResponseMessage(messageChunk) && messageChunk.status === AgentResponseStatus.AGENT_RESPONSE_COMPLETED) {
            mustNewMessage = true;
            setTimeout(() => {
              messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
            }, 50);

            setStatus(Status.IDLE);
          }
        });

      } catch (error) {
        console.error("Connection Error:", error);
      }
    };

    setupChat();

    return () => {
      isMounted = false;
      if (unsubscribe) {
        unsubscribe();
      }
    };
  }, [addMessage]);

  return (
    <div className="flex flex-col h-full w-full overflow-hidden relative">
      <div className="flex-1 overflow-y-auto px-4 py-6 pb-36">
        <div className="max-w-3xl mx-auto w-full flex flex-col gap-4" style={{ willChange: 'scroll-position', contain: 'content' }}>
          {messages.length === 0 ? (
            <div className="my-auto py-20 flex flex-col items-center justify-center">
              <WelcomeHeader />
            </div>
          ) : (
            messages.map((m) => (
              <ChatMessage key={m.id} content={m.content} index={m.id} role={m.role} />
            ))
          )}
          <div ref={messagesEndRef} className="h-1" />
        </div>
      </div>

      <ChatInput
        value={prompt}
        disable={status === Status.PENDING}
        onChange={(e) => setPrompt(e.target.value)}
        onSubmit={handleSubmit}
      />
    </div>
  );
}