import type { LucideIcon } from "lucide-react";
import { Button } from "#/components/ui/button.tsx";

interface EmptyStateProps {
  icon: LucideIcon;
  title: string;
  description: string;
  actionLabel?: string;
  actionIcon?: LucideIcon;
  onAction?: () => void;
}

export function EmptyState({
  icon: Icon,
  title,
  description,
  actionLabel,
  actionIcon: ActionIcon,
  onAction,
}: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-16 gap-3 text-center border rounded-xl bg-card/50">
      <div className="p-4 rounded-full bg-muted">
        <Icon className="w-6 h-6 text-muted-foreground" />
      </div>
      <p className="text-sm font-medium">{title}</p>
      <p className="text-xs text-muted-foreground max-w-xs">{description}</p>
      {actionLabel && onAction && (
        <Button size="sm" variant="outline" className="mt-2" onClick={onAction}>
          {ActionIcon && <ActionIcon className="w-4 h-4 mr-1" />}
          {actionLabel}
        </Button>
      )}
    </div>
  );
}
