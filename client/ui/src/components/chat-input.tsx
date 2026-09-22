import React, { useRef, useState } from "react";
import type { ChatData } from "#/types/chat.ts";
import { Button } from "#/components/ui/button.tsx";
import { PlusIcon, SendIcon } from "lucide-react";
import { Textarea } from "#/components/ui/textarea.tsx";

type ChatInputProps = {
  disabled?: boolean;
  onSend(data: ChatData): void;
};

export default function ChatInput({ disabled, onSend }: ChatInputProps) {
  const [text, setText] = useState("");
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const handleInput = (_: React.InputEvent<HTMLTextAreaElement>) => {
    const textarea = textareaRef.current;
    if (textarea) {
      textarea.style.height = "auto";
      textarea.style.height = `${textarea.scrollHeight}px`;
    }
  };

  const handleSend = () => {
    onSend({
      text: text,
    });
  };

  return (
    <>
      <div className="flex justify-center drop-shadow-accent w-full px-3 py-2 backdrop-blur-3xl">
        <div className="flex items-end gap-3 border border-solid border-accent rounded-xl p-2 w-full max-w-3xl bg-background">
          <Button variant="outline" size="icon" className="mb-0.5 shrink-0">
            <PlusIcon />
          </Button>
          <Textarea
            ref={textareaRef}
            id="textarea-message"
            placeholder="Type your message here."
            rows={1}
            onInput={handleInput}
            className="rounded-lg p-3 border-none focus-visible:ring-0 shadow-none"
            defaultValue={text}
            onChange={(e) => setText(e.target.value)}
          />
          <Button
            variant="outline"
            size="icon"
            className="mb-0.5 shrink-0"
            disabled={disabled}
            onClick={handleSend}
          >
            <SendIcon />
          </Button>
        </div>
      </div>
    </>
  );
}
