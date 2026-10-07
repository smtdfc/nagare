import { Loader2, X } from "lucide-react";
import { Button } from "#/components/ui/button.tsx";
import { Input } from "#/components/ui/input.tsx";
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
import type { AddLLMProviderRequest, Credential } from "@nagare-app/dtos";
import { useEffect, useState } from "react";

const COMPATIBLE_OPTIONS = [{ value: "OpenAI", label: "OpenAI" }];

interface LLMProviderSheetProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  form: AddLLMProviderRequest;
  setForm: React.Dispatch<React.SetStateAction<AddLLMProviderRequest>>;
  onSave: () => void;
  isSaving: boolean;
  credentials: Credential[];
  isCredentialsLoading: boolean;
  mode?: "add" | "edit";
}

export function LLMProviderSheet({
  open,
  onOpenChange,
  form,
  setForm,
  onSave,
  isSaving,
  credentials,
  isCredentialsLoading,
  mode = "add",
}: LLMProviderSheetProps) {
  const [modelInput, setModelInput] = useState("");
  const [authMode, setAuthMode] = useState<"apiKey" | "credential">(
    "credential",
  );

  const handleAddModel = () => {
    const trimmed = modelInput.trim();
    if (!trimmed || form.models.includes(trimmed)) return;
    setForm((prev) => ({ ...prev, models: [...prev.models, trimmed] }));
    setModelInput("");
  };

  const handleRemoveModel = (model: string) => {
    setForm((prev) => ({
      ...prev,
      models: prev.models.filter((m) => m !== model),
    }));
  };

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
            {isEdit ? "Edit LLM Provider" : "Add LLM Provider"}
          </SheetTitle>
          <SheetDescription className="text-xs text-muted-foreground">
            {isEdit
              ? "Update the configuration for this LLM provider."
              : "Connect a new language model provider to Nagare."}
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

          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium">Models</label>
            <p className="text-xs text-muted-foreground">
              Add the model IDs this provider supports.
            </p>
            <div className="flex gap-2">
              <Input
                placeholder="gpt-4o"
                value={modelInput}
                onChange={(e) => setModelInput(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    e.preventDefault();
                    handleAddModel();
                  }
                }}
              />
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={handleAddModel}
              >
                Add
              </Button>
            </div>
            {form.models.length > 0 && (
              <div className="mt-1.5 flex flex-wrap gap-1.5">
                {form.models.map((m) => (
                  <span
                    key={m}
                    className="flex items-center gap-1 px-2 py-0.5 rounded-md bg-muted text-xs font-mono"
                  >
                    {m}
                    <button
                      type="button"
                      onClick={() => handleRemoveModel(m)}
                      className="text-muted-foreground hover:text-foreground"
                      aria-label={`Remove model ${m}`}
                    >
                      <X className="w-3 h-3" />
                    </button>
                  </span>
                ))}
              </div>
            )}
          </div>
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
