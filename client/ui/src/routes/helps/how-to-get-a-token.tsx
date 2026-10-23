import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Button } from "#/components/ui/button";
import { KeyIcon, ExternalLinkIcon } from "lucide-react";

export const Route = createFileRoute('/helps/how-to-get-a-token')({
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
            <h1 className="text-xl font-bold">How to get a token</h1>
          </div>
          <p className="text-sm text-muted-foreground leading-relaxed">
            A token is required to sign in to the Nagare web dashboard. Generate
            one with the Nagare CLI, then enter it together with your Gateway
            URL on the sign-in page.
          </p>

          <section className="flex flex-col gap-4">
            <h2 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground">
              Requirements
            </h2>
            <ul className="list-disc list-inside text-sm text-muted-foreground space-y-1">
              <li>
                Nagare is installed (<span className="font-mono">nagare --version</span>)
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

              <Step number={2} title="Generate a token">
                Run:
                <CodeBlock>nagare auth token</CodeBlock>
                <p className="pt-2 text-xs text-muted-foreground">
                  Copy the printed token (for example:{" "}
                  <span className="font-mono">Your token is: abc123...</span>).
                </p>
              </Step>

              <Step number={3} title="Sign in to the dashboard">
                Open <a href="http://localhost:3005" className="underline">
                  http://localhost:3005
                </a>{" "}
                and sign in with:
                <ul className="list-disc list-inside text-xs text-muted-foreground space-y-1 mt-4">
                  <li>
                    <b>Gateway URL:</b> <span className="font-mono">http://localhost:9832</span>
                  </li>
                  <li>
                    <b>Token:</b> the token copied in the previous step
                  </li>
                </ul>
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
        <div >{children}</div>
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
