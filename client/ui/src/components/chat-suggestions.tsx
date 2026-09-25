type ChatSuggestionsProps = {
  onSelect: (prompt: string) => void;
  suggestions?: {
    icon: React.ReactNode;
    title: string;
    desc: string;
    prompt: string;
  }[];
};

export function ChatSuggestions({
  onSelect,
  suggestions,
}: ChatSuggestionsProps) {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 w-full max-w-2xl mx-auto my-6 px-4">
      {suggestions?.map((item, index) => (
        <button
          key={index}
          onClick={() => onSelect(item.prompt)}
          className="flex flex-col text-left p-4 rounded-xl border border-border/60 bg-card/50 hover:bg-accent/40 hover:border-primary/40 transition-all duration-200 group shadow-2xs hover:shadow-md cursor-pointer"
        >
          <div className="p-2 w-fit rounded-lg bg-background border border-border/40 mb-3 group-hover:scale-110 transition-transform">
            {item.icon}
          </div>
          <span className="text-sm font-medium text-foreground mb-1 group-hover:text-primary transition-colors">
            {item.title}
          </span>
          <span className="text-xs text-muted-foreground line-clamp-1">
            {item.desc}
          </span>
        </button>
      ))}
    </div>
  );
}
