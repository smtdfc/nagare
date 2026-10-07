import { Loader2 } from "lucide-react";
import { useState } from "react";
import { Button } from "#/components/ui/button.tsx";
import { Input } from "#/components/ui/input.tsx";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from "#/components/ui/sheet.tsx";

export interface CredentialForm {
  name: string;
  apiKey: string;
}

interface CredentialSheetProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  form: CredentialForm;
  setForm: React.Dispatch<React.SetStateAction<CredentialForm>>;
  onSave: () => void;
  isSaving: boolean;
  mode: "add" | "edit";
}

export function CredentialSheet({
  open,
  onOpenChange,
  form,
  setForm,
  onSave,
  isSaving,
  mode,
}: CredentialSheetProps) {
  const [showApiKey, setShowApiKey] = useState(false);
  const isEdit = mode === "edit";

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-full overflow-y-auto sm:max-w-md">
        <SheetHeader>
          <SheetTitle className="text-lg font-semibold">
            {isEdit ? "Edit Credential" : "Add Credential"}
          </SheetTitle>
          <SheetDescription className="text-xs text-muted-foreground">
            {isEdit
              ? "Update the credential name or replace its API key."
              : "Store an API key for reuse across LLM providers."}
          </SheetDescription>
        </SheetHeader>

        <div className="flex flex-col gap-4 px-4 py-2">
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium">Name *</label>
            <Input
              placeholder="e.g. Personal OpenAI"
              value={form.name}
              onChange={(event) =>
                setForm((previous) => ({
                  ...previous,
                  name: event.target.value,
                }))
              }
            />
          </div>

          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium">
              API Key{isEdit ? " (optional)" : " *"}
            </label>
            <div className="relative">
              <Input
                type={showApiKey ? "text" : "password"}
                placeholder={
                  isEdit ? "Leave blank to keep existing key" : "sk-..."
                }
                value={form.apiKey}
                onChange={(event) =>
                  setForm((previous) => ({
                    ...previous,
                    apiKey: event.target.value,
                  }))
                }
                className="pr-20"
              />
              <button
                type="button"
                className="absolute inset-y-0 right-2 text-xs text-muted-foreground hover:text-foreground"
                onClick={() => setShowApiKey((visible) => !visible)}
              >
                {showApiKey ? "Hide" : "Show"}
              </button>
            </div>
          </div>
        </div>

        <SheetFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            onClick={onSave}
            disabled={
              isSaving || !form.name.trim() || (!isEdit && !form.apiKey)
            }
          >
            {isSaving && <Loader2 className="mr-1 h-4 w-4 animate-spin" />}
            {isEdit ? "Save Changes" : "Save Credential"}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
