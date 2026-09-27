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
        h1: ({ node, ...props }) => (
          <h1 className="text-2xl font-bold my-4" {...props} />
        ),

        h2: ({ node, ...props }) => (
          <h2 className="text-xl font-semibold my-3" {...props} />
        ),

        h3: ({ node, ...props }) => (
          <h3 className="text-lg font-semibold my-2" {...props} />
        ),
        h4: ({ node, ...props }) => (
          <h4 className="text-md font-semibold my-2" {...props} />
        ),
        h5: ({ node, ...props }) => (
          <h5 className="text-sm font-semibold my-1" {...props} />
        ),
        h6: ({ node, ...props }) => (
          <h6 className="text-xs font-semibold my-1" {...props} />
        ),
        p: ({ node, ...props }) => (
          <p className="my-2 leading-relaxed" {...props} />
        ),
        a: ({ node, ...props }) => (
          <a
            className="text-primary underline hover:text-primary/80"
            target="_blank"
            rel="noopener noreferrer"
            {...props}
          />
        ),
        ul: ({ node, ...props }) => (
          <ul className="list-disc list-inside my-2" {...props} />
        ),
        ol: ({ node, ...props }) => (
          <ol className="list-decimal list-inside my-2" {...props} />
        ),
        li: ({ node, ...props }) => <li className="my-1" {...props} />,
        img: ({ node, ...props }) => (
          <img
            className="my-2 max-w-full rounded-lg border border-border"
            {...props}
          />
        ),
        code: ({ node, className, children, ...props }) => {
          const match = /language-(\w+)/.exec(className || "");
          return match ? (
            <pre
              className="my-2 overflow-x-auto rounded-lg border border-border bg-muted p-4 text-sm"
              {...(props as React.HTMLAttributes<HTMLPreElement>)}
            >
              <code className={className}>{children}</code>
            </pre>
          ) : (
            <code className="rounded bg-muted px-1 py-0.5 text-sm" {...props}>
              {children}
            </code>
          );
        },
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
