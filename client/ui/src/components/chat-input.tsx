import { useRef, useState, useEffect, type SyntheticEvent } from "react";
import type { ChatData } from "#/types/chat.ts";
import { Button } from "#/components/ui/button.tsx";
import { PlusIcon, SendIcon, SquareIcon } from "lucide-react";
import { Textarea } from "#/components/ui/textarea.tsx";

type ChatInputProps = {
  isProcess?: boolean;
  currentChatData?: ChatData;
  onSend(data: ChatData): void;
  onCancel(): void;
};

export default function ChatInput({
  currentChatData,
  onSend,
  onCancel,
  isProcess,
}: ChatInputProps) {
  const [text, setText] = useState(currentChatData?.text || "");
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  useEffect(() => {
    if (currentChatData?.text !== undefined) {
      setText(currentChatData.text);
    }
  }, [currentChatData]);

  const handleInput = (_: SyntheticEvent<HTMLTextAreaElement>) => {
    const textarea = textareaRef.current;
    if (textarea) {
      textarea.style.height = "auto";
      textarea.style.height = `${textarea.scrollHeight}px`;
    }
  };

  const handleSend = () => {
    if (isProcess) {
      onCancel();
      return;
    }

    if (!text.trim()) return;
    onSend({
      text: text,
    });
    setText("");
    if (textareaRef.current) {
      textareaRef.current.style.height = "auto";
    }
  };

  return (
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
          className="rounded-lg p-3 border-none focus-visible:ring-0 shadow-none resize-none"
          value={text}
          onChange={(e) => setText(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              handleSend();
            }
          }}
        />
        <Button
          variant={isProcess ? "destructive" : "outline"}
          size="icon"
          className="mb-0.5 shrink-0"
          onClick={handleSend}
        >
          {isProcess ? <SquareIcon className="h-4 w-4 fill-current" /> : <SendIcon />}
        </Button>
      </div>
    </div>
  );
}
