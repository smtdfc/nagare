import { createFileRoute } from "@tanstack/react-router";
import { Sparkles, Compass, Lightbulb, Laptop } from "lucide-react";
import { useChat } from "#/hooks/use-chat.ts";
import { useEffect } from "react";
import { ChatSuggestions } from "#/components/chat-suggestions";

export const Route = createFileRoute("/(dashboard)/chat/new")({
  component: RouteComponent,
});

function RouteComponent() {
  const reset = useChat((c) => c.reset);

  const suggestions = [
    {
      icon: <Sparkles className="w-4 h-4 text-amber-500" />,
      title: "Draft creative ideas",
      desc: "Plan and outline a new project",
      prompt:
        "Help me brainstorm ideas and create a project plan for a task management web app.",
    },
    {
      icon: <Compass className="w-4 h-4 text-emerald-500" />,
      title: "Explore knowledge",
      desc: "Get answers to science and life questions",
      prompt:
        "Explain Einstein's theory of relativity in a simple, easy-to-understand way.",
    },
    {
      icon: <Lightbulb className="w-4 h-4 text-purple-500" />,
      title: "Get expert advice",
      desc: "Receive professional insights and solutions",
      prompt:
        "Suggest some effective time management methods for software developers.",
    },
    {
      icon: <Laptop className="w-4 h-4 text-blue-500" />,
      title: "PC & Workflow Assistant",
      desc: "Manage tasks, organize files, and optimize workflow",
      prompt:
        "Help me organize my workspace, write a script to clean up files, or optimize my daily computer workflow.",
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
      <ChatSuggestions
        onSelect={(prompt) => {
          reset();
        }}
        suggestions={suggestions}
      />
    </>
  );
}
