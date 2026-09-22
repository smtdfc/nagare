import { createFileRoute } from "@tanstack/react-router";
import { Compass, Sparkles } from "lucide-react";
import { useChat } from "#/hooks/use-chat.ts";
import { useEffect } from "react";

export const Route = createFileRoute("/(dashboard)/chat/new")({
  component: RouteComponent,
});

function RouteComponent() {
  const reset = useChat((c) => c.reset);

  const suggestions = [
    {
      icon: <Sparkles className="w-4 h-4 text-amber-500" />,
      title: "Drafting creative ideas",
      desc: "Planning the new project",
    },
    {
      icon: <Compass className="w-4 h-4 text-emerald-500" />,
      title: "Discover knowledge",
      desc: "Answering any questions",
    },
  ];

  useEffect(() => {
    reset();
  }, []);

  return (
    <>
      <div className="space-y-2 text-center">
        <h1 className="text-3xl md:text-4xl font-bold tracking-tight">
          Nagare Here 👋
        </h1>
        <p className="text-muted-foreground text-sm md:text-base">
          How can I help you today? Choose a suggestion below or start a
          conversation right away.
        </p>
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 w-full pt-5 px-4">
        {suggestions.map((item, index) => (
          <div
            key={index}
            className="flex flex-col items-start p-4 rounded-xl border border-border bg-card hover:bg-accent/50 transition-all cursor-pointer text-left shadow-sm"
          >
            <div className="p-3 rounded-lg bg-background mb-4 border border-border">
              {item.icon}
            </div>
            <h3 className="font-medium text-sm">{item.title}</h3>
            <p className="text-xs text-muted-foreground mt-1">{item.desc}</p>
          </div>
        ))}
      </div>
    </>
  );
}
