import { Separator } from "#/components/ui/separator.tsx";
import type { ReactNode } from "react";

interface SettingsHeaderProps {
  title: string;
  description: string;
  action?: ReactNode;
}

export function SettingsHeader({
  title,
  description,
  action,
}: SettingsHeaderProps) {
  return (
    <div className="flex flex-col gap-6 w-full">
      <div className="flex items-center justify-between">
        <div className="flex flex-col gap-1">
          <h3 className="text-lg font-semibold tracking-tight">{title}</h3>
          <p className="text-xs text-muted-foreground">{description}</p>
        </div>
        {action && <div>{action}</div>}
      </div>
      <Separator />
    </div>
  );
}
