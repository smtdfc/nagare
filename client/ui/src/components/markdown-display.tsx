import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";

type MarkdownDisplayProps = {
  content: string;
};

export default function MarkdownDisplay({ content }: MarkdownDisplayProps) {
  return (
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
      {content}
    </Markdown>
  );
}
