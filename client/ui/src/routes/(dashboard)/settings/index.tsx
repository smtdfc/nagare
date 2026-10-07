import { SettingItem, type SettingSection } from "#/components/setting-item";
import { createFileRoute } from "@tanstack/react-router";
import {
  Settings,
  Palette,
  Puzzle,
  Brain,
  Database,
  FileText,
  KeyRound,
} from "lucide-react";

export const Route = createFileRoute("/(dashboard)/settings/")({
  component: RouteComponent,
  staticData: {
    breadcrumb: "All Settings",
  },
});

const settingSections: SettingSection[] = [
  {
    category: "System & Personalization",
    items: [
      {
        title: "General Settings",
        description: "Manage current model and provider.",
        icon: Settings,
        href: "/settings/general",
      },
    ],
  },
  {
    category: "Providers & Credentials",
    items: [
      {
        title: "Credentials",
        description: "Store and manage reusable provider API keys.",
        icon: KeyRound,
        href: "/settings/credentials",
      },
      {
        title: "LLM Providers",
        description:
          "Manage and configure your LLM providers (OpenAI, Anthropic, etc.).",
        icon: Brain,
        href: "/settings/llm-providers",
      },
    ],
  },
  {
    category: "Plugins & Integrations",
    items: [
      {
        title: "Plugins",
        description:
          "Enable, disable, or configure integrated plugins (such as Google Calendar, Vector DB...).",
        icon: Puzzle,
        href: "/settings/plugins",
      },
    ],
  },
  {
    category: "Data",
    items: [
      {
        title: "Database",
        description: "Manage your database and its configuration.",
        icon: Database,
        href: "/settings/database",
      },
      {
        title: "Logs",
        description: "View and manage system logs.",
        icon: FileText,
        href: "/settings/logs",
      },
    ],
  },
];

function RouteComponent() {
  return (
    <div className="flex flex-col gap-8 px-4  py-6  max-w-5xl mx-auto w-full">
      <div>
        <h3 className="text-lg font-bold tracking-tight">Settings</h3>
        <p className="text-xs text-muted-foreground mt-1">
          Manage system options, plugins, and your Nagare configuration.
        </p>
      </div>

      <div className="flex flex-col gap-6">
        {settingSections.map((section, index) => (
          <div key={index} className="flex flex-col gap-3">
            <h2 className="text-xs font-semibold uppercase tracking-wider text-muted-foreground px-1">
              {section.category}
            </h2>
            <div className="grid grid-cols-1 gap-3">
              {section.items.map((item, itemIndex) => (
                <SettingItem key={itemIndex} item={item} />
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
