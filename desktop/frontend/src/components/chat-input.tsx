import { PlusIcon, SendIcon } from "lucide-react";
import { Button } from "./ui/button";
import { Textarea } from "./ui/textarea";

interface ChatInputProps {
    disable?: boolean,
    value: string;
    onChange: (e: React.ChangeEvent<HTMLTextAreaElement>) => void;
    onSubmit: () => void;
}

export default function ChatInput({ disable, value, onChange, onSubmit }: ChatInputProps) {
    const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            onSubmit();
        }
    };

    return (
        <div className="absolute bottom-0 left-0 right-0 p-4 bg-linear-to-t from-background via-background/80 to-transparent z-20">
            <div className="max-w-3xl mx-auto flex items-center gap-3 bg-background/90 backdrop-blur-md p-3 rounded-2xl border shadow-xl">
                <Button disabled={disable} className="rounded-full w-10 h-10 p-0 shrink-0 flex items-center justify-center" variant="outline" type="button">
                    <PlusIcon />
                </Button>

                <div className="flex-1 overflow-hidden flex items-center">
                    <Textarea
                        id="chat"
                        value={value}
                        onChange={onChange}
                        onKeyDown={handleKeyDown}
                        placeholder="Ask Nagare something ..."
                        rows={1}
                        className="w-full resize-none border-0 focus-visible:ring-0 shadow-none bg-transparent whitespace-pre-wrap wrap-break-word min-h-6 max-h-32 m-0 py-1"
                    />
                </div>

                <Button
                    disabled={disable}
                    className="rounded-full w-10 h-10 p-0 shrink-0 flex items-center justify-center"
                    type="button"
                    onClick={onSubmit}
                >
                    <SendIcon />
                </Button>
            </div>
        </div>
    );
}