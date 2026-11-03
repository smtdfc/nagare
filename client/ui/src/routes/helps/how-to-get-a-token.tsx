import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Button } from "#/components/ui/button";
import { KeyIcon, ExternalLinkIcon } from "lucide-react";

export const Route = createFileRoute("/helps/how-to-get-a-token")({
  component: RouteComponent,
});

function RouteComponent() {
  const navigate = useNavigate();

  return (
    <div className="flex min-h-svh flex-col items-center justify-center bg-background p-6 md:p-10">
      <div className="w-full max-w-md">
        <div className="flex flex-col gap-6">
          <div className="flex items-center gap-3">
            <KeyIcon className="size-5 text-primary" />
            <h1 className="text-xl font-bold">How to sign in</h1>
          </div>
          <p className="text-sm text-muted-foreground leading-relaxed">
            The easiest way to sign in is to let the CLI open the dashboard for
            you. It generates a token, attaches the Gateway URL, and the
            dashboard signs in automatically.
          </p>

          <section className="flex flex-col gap-4">
            <h2 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground">
              Requirements
            </h2>
            <ul className="list-disc list-inside text-sm text-muted-foreground space-y-1">
              <li>
                Nagare is installed (
                <span className="font-mono">nagare --version</span>)
              </li>
            </ul>
          </section>

          <section className="flex flex-col gap-5">
            <h2 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground">
              Steps
            </h2>

            <div className="flex flex-col gap-3">
              <Step number={1} title="Start the Gateway">
                Open a terminal and run:
                <CodeBlock>nagare gateway start</CodeBlock>
                <p className="pt-2 text-xs text-muted-foreground">
                  The Gateway runs by default at{" "}
                  <span className="font-mono">http://localhost:9832</span>.
                </p>
              </Step>

              <Step number={2} title="Start the Web UI">
                Run:
                <CodeBlock>nagare web</CodeBlock>
                <p className="pt-2 text-xs text-muted-foreground">
                  The dashboard opens in your browser already signed in. Use{" "}
                  <span className="font-mono">--no-open</span> to print the
                  sign-in URL instead, or{" "}
                  <span className="font-mono">--gateway</span> when the Gateway
                  runs on another host.
                </p>
              </Step>

              <Step number={3} title="Sign in manually (optional)">
                If you open the dashboard yourself, generate a token and paste
                it on the sign-in page:
                <CodeBlock>nagare auth token</CodeBlock>
                <p className="pt-2 text-xs text-muted-foreground">
                  The gateway defaults to{" "}
                  <span className="font-mono">http://localhost:9832</span>. Use
                  the <b>Advanced</b> link on the sign-in page to point at
                  another gateway.
                </p>
              </Step>
            </div>
          </section>

          <div className="flex items-center justify-between gap-2 pt-4 border-t border-border">
            <Button
              variant="outline"
              size="sm"
              onClick={() => navigate({ to: "/auth" })}
            >
              <ExternalLinkIcon className="size-3.5 mr-1.5" />
              Go to sign in
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}

function Step({
  number,
  title,
  children,
}: {
  number: number;
  title: string;
  children: React.ReactNode;
}) {
  return (
    <div className="flex gap-3">
      <div className="flex size-6 shrink-0 items-center justify-center rounded-full bg-primary text-primary-foreground text-xs font-semibold">
        {number}
      </div>
      <div className="flex flex-col gap-3">
        <h3 className="text-sm font-semibold">{title}</h3>
        <div>{children}</div>
      </div>
    </div>
  );
}

function CodeBlock({ children }: { children: string }) {
  return (
    <div className="mt-3">
      <pre className="my-0 overflow-x-auto rounded-lg border border-border bg-muted p-3 text-xs font-mono text-foreground">
        <code>{children}</code>
      </pre>
    </div>
  );
}
