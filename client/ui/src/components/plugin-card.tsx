import type { Plugin } from "@nagare-app/dtos";
import { PluginService } from "@nagare-app/services";
import { useState } from "react";
import { toast } from "./ui/toast";
import { showToastError } from "#/lib/toast";
import { Loader2, Plug, PowerOff, Trash2, Zap } from "lucide-react";
import { Button } from "./ui/button";

export function PluginCard({
  plugin,
  onToggle,
  onUninstall,
}: {
  plugin: Plugin;
  onToggle: (id: string, active: boolean) => void;
  onUninstall: (id: string) => void;
}) {
  const [isTogglingActive, setIsTogglingActive] = useState(false);
  const [isUninstalling, setIsUninstalling] = useState(false);

  const handleToggle = async () => {
    setIsTogglingActive(true);
    try {
      if (plugin.isActive) {
        await PluginService.deactivate(plugin.id);
        onToggle(plugin.id, false);
        toast.add({
          title: "Deactivated",
          type: "info",
          description: `Plugin "${plugin.name}" has been deactivated.`,
        });
      } else {
        await PluginService.activate(plugin.id);
        onToggle(plugin.id, true);
        toast.add({
          title: "Activated",
          type: "success",
          description: `Plugin "${plugin.name}" is now active.`,
        });
      }
    } catch (e) {
      showToastError(e);
    } finally {
      setIsTogglingActive(false);
    }
  };

  const handleUninstall = async () => {
    setIsUninstalling(true);
    try {
      await PluginService.uninstall(plugin.id);
      onUninstall(plugin.id);
      toast.add({
        title: "Uninstalled",
        type: "success",
        priority: "high",
        description: `Plugin "${plugin.name}" has been uninstalled.`,
      });
    } catch (e) {
      showToastError(e);
    } finally {
      setIsUninstalling(false);
    }
  };

  return (
    <div className="rounded-xl border bg-card shadow-sm p-4 flex items-start gap-4">
      <div
        className={`p-2.5 rounded-lg shrink-0 ${plugin.isActive ? "bg-primary/10 text-primary" : "bg-muted text-muted-foreground"}`}
      >
        <Plug className="w-4 h-4" />
      </div>

      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 flex-wrap">
          <p className="text-sm font-medium truncate">{plugin.name}</p>
          <span
            className={`px-1.5 py-0.5 rounded-full text-[10px] font-medium ${
              plugin.isActive
                ? "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400"
                : "bg-muted text-muted-foreground"
            }`}
          >
            {plugin.isActive ? "Active" : "Inactive"}
          </span>
        </div>
        <p className="text-xs text-muted-foreground mt-0.5">
          by <span className="font-medium">{plugin.author}</span> &middot; v
          {plugin.version}
        </p>
      </div>

      <div className="flex items-center gap-1 shrink-0">
        <Button
          variant={plugin.isActive ? "outline" : "default"}
          size="sm"
          onClick={handleToggle}
          disabled={isTogglingActive}
          aria-label={plugin.isActive ? "Deactivate" : "Activate"}
        >
          {isTogglingActive ? (
            <Loader2 className="w-3.5 h-3.5 animate-spin" />
          ) : plugin.isActive ? (
            <PowerOff className="w-3.5 h-3.5" />
          ) : (
            <Zap className="w-3.5 h-3.5" />
          )}
          {plugin.isActive ? "Deactivate" : "Activate"}
        </Button>
        <Button
          variant="ghost"
          size="icon-sm"
          onClick={handleUninstall}
          disabled={isUninstalling}
          aria-label="Uninstall plugin"
        >
          {isUninstalling ? (
            <Loader2 className="w-4 h-4 animate-spin" />
          ) : (
            <Trash2 className="w-4 h-4 text-destructive" />
          )}
        </Button>
      </div>
    </div>
  );
}
