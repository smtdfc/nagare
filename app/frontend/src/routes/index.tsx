import ChatInput from '#/components/chat-input.tsx';
import ChatMessage from '#/components/chat-message.tsx';
import WelcomeHeader from '#/components/welcome-header.tsx';
import { createFileRoute } from '@tanstack/react-router'
import { useState, useRef, useEffect } from 'react';

export const Route = createFileRoute('/')({ component: Home })

interface MessageItem {
  id: string;
  role: 'user' | 'agent' | 'system';
  content: string;
}

enum Status {
  IDLE,
  PENDING
}

function Home() {
  const [status, setStatus] = useState<Status>(Status.IDLE);
  const [messages, setMessages] = useState<MessageItem[]>([]);
  const [prompt, setPrompt] = useState<string>("");
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleSubmit = () => {
    setStatus(Status.PENDING);
    setMessages([{
      id: "1",
      role: "user",
      content: prompt
    }]);
    setPrompt("");
  }

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  return (
    <div className="flex flex-col h-full w-full overflow-hidden relative">
      <div className="flex-1 overflow-y-auto px-4 py-6 pb-36">
        <div className="max-w-3xl mx-auto w-full flex flex-col gap-4">
          {messages.length === 0 ? (
            <div className="my-auto py-20 flex flex-col items-center justify-center">
              <WelcomeHeader />
            </div>
          ) : (
            messages.map((m, _) => (
              <ChatMessage content={m.content} index={m.id} role={m.role} />
            ))
          )}
          <div ref={messagesEndRef} className="h-1" />
        </div>
      </div>

      <ChatInput
        value={prompt}
        disable={status == Status.PENDING}
        onChange={(e) => setPrompt(e.target.value)}
        onSubmit={() => handleSubmit()}
      />
    </div>
  );
}