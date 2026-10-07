import type { Credential } from "@nagare-app/dtos";
import { Edit2, KeyRound, Loader2, Trash2 } from "lucide-react";
import { Button } from "#/components/ui/button.tsx";

interface CredentialCardProps {
  credential: Credential;
  onEdit: (credential: Credential) => void;
  onDelete: (id: string) => void;
  isDeleting: boolean;
}

export function CredentialCard({
  credential,
  onEdit,
  onDelete,
  isDeleting,
}: CredentialCardProps) {
  return (
    <div className="flex items-center justify-between gap-3 overflow-hidden rounded-xl border bg-card p-4 shadow-sm">
      <div className="flex min-w-0 items-center gap-3">
        <div className="shrink-0 rounded-lg bg-primary/10 p-2 text-primary">
          <KeyRound className="h-4 w-4" />
        </div>
        <div className="min-w-0">
          <p className="truncate text-sm font-medium">{credential.name}</p>
          <p className="text-xs text-muted-foreground">
            API key stored securely
          </p>
        </div>
      </div>
      <div className="flex shrink-0 items-center gap-1">
        <Button
          variant="ghost"
          size="icon-sm"
          onClick={() => onEdit(credential)}
          aria-label={`Edit ${credential.name}`}
        >
          <Edit2 className="h-4 w-4" />
        </Button>
        <Button
          variant="ghost"
          size="icon-sm"
          onClick={() => onDelete(credential.id)}
          disabled={isDeleting}
          aria-label={`Delete ${credential.name}`}
        >
          {isDeleting ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <Trash2 className="h-4 w-4 text-destructive" />
          )}
        </Button>
      </div>
    </div>
  );
}
