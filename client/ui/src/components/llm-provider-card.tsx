import type { LLMProvider } from "@nagare-app/dtos";
import { LLMProviderService } from "@nagare-app/services";
import { useState } from "react";
import { toast } from "./ui/toast";
import { showToastError } from "#/lib/toast";
import { Brain, ChevronDown, ChevronUp, Loader2, Trash2 } from "lucide-react";
import { Button } from "./ui/button";
import { Separator } from "./ui/separator";

export function LLMProviderCard({
  provider,
  onDelete,
}: {
  provider: LLMProvider;
  onDelete: (id: string) => void;
}) {
  const [expanded, setExpanded] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  const handleDelete = async () => {
    setIsDeleting(true);
    try {
      await LLMProviderService.delete(provider.id);
      onDelete(provider.id);
      toast.add({
        title: "Deleted",
        type: "success",
        priority: "high",
        description: `Provider "${provider.name}" has been removed.`,
      });
    } catch (e) {
      showToastError(e);
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <div className="rounded-xl border bg-card shadow-sm overflow-hidden">
      <div className="flex items-center justify-between p-4 gap-3">
        <div className="flex items-center gap-3 min-w-0">
          <div className="p-2 rounded-lg bg-primary/10 text-primary shrink-0">
            <Brain className="w-4 h-4" />
          </div>
          <div className="min-w-0">
            <p className="text-sm font-medium truncate">{provider.name}</p>
            <p className="text-xs text-muted-foreground capitalize">
              {provider.compatible}
            </p>
          </div>
        </div>
        <div className="flex items-center gap-1 shrink-0">
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={() => setExpanded((v) => !v)}
            aria-label="Toggle details"
          >
            {expanded ? (
              <ChevronUp className="w-4 h-4" />
            ) : (
              <ChevronDown className="w-4 h-4" />
            )}
          </Button>
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={handleDelete}
            disabled={isDeleting}
            aria-label="Delete provider"
          >
            {isDeleting ? (
              <Loader2 className="w-4 h-4 animate-spin" />
            ) : (
              <Trash2 className="w-4 h-4 text-destructive" />
            )}
          </Button>
        </div>
      </div>

      {expanded && (
        <>
          <Separator />
          <div className="p-4 flex flex-col gap-3 text-sm">
            {provider.baseUrl && (
              <div>
                <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                  Base URL
                </span>
                <p className="mt-0.5 font-mono text-xs break-all">
                  {provider.baseUrl}
                </p>
              </div>
            )}
            {provider.models.length > 0 && (
              <div>
                <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                  Models ({provider.models.length})
                </span>
                <div className="mt-1.5 flex flex-wrap gap-1.5">
                  {provider.models.map((m) => (
                    <span
                      key={m}
                      className="px-2 py-0.5 rounded-md bg-muted text-xs font-mono"
                    >
                      {m}
                    </span>
                  ))}
                </div>
              </div>
            )}
          </div>
        </>
      )}
    </div>
  );
}
