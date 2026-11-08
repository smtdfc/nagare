import { Loader2 } from "lucide-react";
import { Button } from "#/components/ui/button.tsx";
import { Input } from "#/components/ui/input.tsx";
import { ModelListEditor } from "#/components/model-list-editor";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
  SheetFooter,
} from "#/components/ui/sheet.tsx";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "#/components/ui/select.tsx";
import type { AddModelProviderRequest, Credential } from "@nagare-app/dtos";
import { useEffect, useState } from "react";

const COMPATIBLE_OPTIONS = [
  { value: "OpenAI", label: "OpenAI" },
  { value: "Google", label: "Google (Gemini)" },
];

interface ModelProviderSheetProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  form: AddModelProviderRequest;
  setForm: React.Dispatch<React.SetStateAction<AddModelProviderRequest>>;
  onSave: () => void;
  isSaving: boolean;
  credentials: Credential[];
  isCredentialsLoading: boolean;
  mode?: "add" | "edit";
}

export function ModelProviderSheet({
  open,
  onOpenChange,
  form,
  setForm,
  onSave,
  isSaving,
  credentials,
  isCredentialsLoading,
  mode = "add",
}: ModelProviderSheetProps) {
  const [authMode, setAuthMode] = useState<"apiKey" | "credential">(
    "credential",
  );

  const handleCompatibleChange = (value: string | null) => {
    if (!value) return;
    setForm((prev) => ({ ...prev, compatible: value }));
  };

  const isEdit = mode === "edit";
  const selectedCredential = credentials.find(
    (credential) => credential.id === form.credentialId,
  );

  useEffect(() => {
    if (open) {
      setAuthMode(form.credentialId ? "credential" : "apiKey");
    }
  }, [open]);

  const handleAuthModeChange = (value: string | null) => {
    if (value !== "apiKey" && value !== "credential") return;
    setAuthMode(value);
    setForm((prev) => ({
      ...prev,
      apiKey: value === "apiKey" ? prev.apiKey : "",
      credentialId: value === "credential" ? prev.credentialId : "",
    }));
  };

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-full sm:max-w-md overflow-y-auto">
        <SheetHeader>
          <SheetTitle className="text-lg font-semibold">
            {isEdit ? "Edit Model Provider" : "Add Model Provider"}
          </SheetTitle>
          <SheetDescription className="text-xs text-muted-foreground">
            {isEdit
              ? "Update the configuration for this model provider."
              : "Connect a new model provider to Nagare."}
          </SheetDescription>
        </SheetHeader>

        <div className="flex flex-col gap-4 px-4 py-2">
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium">Name *</label>
            <Input
              placeholder="e.g. My OpenAI"
              value={form.name}
              onChange={(e) =>
                setForm((prev) => ({ ...prev, name: e.target.value }))
              }
            />
          </div>

          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium">Compatible with</label>
            <Select
              value={form.compatible}
              onValueChange={handleCompatibleChange}
            >
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Select provider type" />
              </SelectTrigger>
              <SelectContent>
                {COMPATIBLE_OPTIONS.map((opt) => (
                  <SelectItem key={opt.value} value={opt.value}>
                    {opt.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium">Authentication</label>
            <Select value={authMode} onValueChange={handleAuthModeChange}>
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Select authentication method">
                  {authMode === "credential"
                    ? "Use saved credential"
                    : "Enter API key directly"}
                </SelectValue>
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="credential">Use saved credential</SelectItem>
                <SelectItem value="apiKey">Enter API key directly</SelectItem>
              </SelectContent>
            </Select>
          </div>

          {authMode === "credential" ? (
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium">Credential *</label>
              <Select
                value={form.credentialId}
                onValueChange={(value) => {
                  if (value) {
                    setForm((prev) => ({ ...prev, credentialId: value }));
                  }
                }}
                disabled={isCredentialsLoading || credentials.length === 0}
              >
                <SelectTrigger className="w-full">
                  <SelectValue
                    placeholder={
                      isCredentialsLoading
                        ? "Loading credentials..."
                        : "Select a credential"
                    }
                  >
                    {selectedCredential?.name}
                  </SelectValue>
                </SelectTrigger>
                <SelectContent>
                  {credentials.map((credential) => (
                    <SelectItem key={credential.id} value={credential.id}>
                      {credential.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {!isCredentialsLoading && credentials.length === 0 && (
                <p className="text-xs text-muted-foreground">
                  Add a credential before creating a provider.
                </p>
              )}
            </div>
          ) : (
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium">API Key *</label>
              <Input
                type="password"
                placeholder="sk-..."
                value={form.apiKey}
                onChange={(event) =>
                  setForm((prev) => ({ ...prev, apiKey: event.target.value }))
                }
              />
            </div>
          )}

          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium">Base URL</label>
            <p className="text-xs text-muted-foreground">
              Leave empty to use the default provider endpoint.
            </p>
            <Input
              placeholder="https://api.openai.com/v1"
              value={form.baseUrl}
              onChange={(e) =>
                setForm((prev) => ({ ...prev, baseUrl: e.target.value }))
              }
            />
          </div>

          <ModelListEditor
            models={form.models}
            onChange={(models) => setForm((prev) => ({ ...prev, models }))}
          />
        </div>

        <SheetFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            onClick={onSave}
            disabled={
              isSaving ||
              !form.name.trim() ||
              (authMode === "credential" && !form.credentialId) ||
              (authMode === "apiKey" && !form.apiKey)
            }
          >
            {isSaving && <Loader2 className="w-4 h-4 animate-spin" />}
            {isEdit ? "Save Changes" : "Save Provider"}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
