import { createFileRoute } from "@tanstack/react-router";
import { Plus, Brain } from "lucide-react";
import { useEffect, useState, useCallback } from "react";
import { Button } from "#/components/ui/button.tsx";
import { Skeleton } from "#/components/ui/skeleton.tsx";
import { LLMProviderService } from "@nagare-app/services";
import { showToastError } from "#/lib/toast.ts";
import { toast } from "#/components/ui/toast.tsx";
import type { LLMProvider, AddLLMProviderRequest } from "@nagare-app/dtos";
import { LLMProviderSheet } from "#/components/llm-provider-sheet.tsx";
import { LLMProviderCard } from "#/components/llm-provider-card";
import { EmptyState } from "#/components/empty-state";
import { SettingsHeader } from "#/components/settings-header";

export const Route = createFileRoute("/(dashboard)/settings/llm-providers/")({
  component: RouteComponent,
  staticData: {
    breadcrumb: "LLM Providers settings",
  },
});

const emptyForm = (): AddLLMProviderRequest => ({
  name: "",
  compatible: "OpenAI",
  apiKey: "",
  models: [],
  baseUrl: "",
});

function RouteComponent() {
  const [providers, setProviders] = useState<LLMProvider[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [sheetOpen, setSheetOpen] = useState(false);
  const [form, setForm] = useState<AddLLMProviderRequest>(emptyForm());
  const [isSaving, setIsSaving] = useState(false);

  const loadProviders = useCallback(async () => {
    setIsLoading(true);
    try {
      const list = await LLMProviderService.list();
      setProviders(list ?? []);
    } catch (e) {
      showToastError(e);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadProviders();
  }, [loadProviders]);

  const handleSave = async () => {
    if (!form.name.trim()) return;
    setIsSaving(true);
    try {
      const created = await LLMProviderService.add(form);
      setProviders((prev) => [...prev, created]);
      setSheetOpen(false);
      setForm(emptyForm());
      toast.add({
        title: "Added",
        type: "success",
        priority: "high",
        description: `Provider "${created.name}" has been added.`,
      });
    } catch (e) {
      showToastError(e);
    } finally {
      setIsSaving(false);
    }
  };

  const handleDelete = (id: string) => {
    setProviders((prev) => prev.filter((p) => p.id !== id));
  };

  return (
    <div className="flex flex-col items-center w-full p-6">
      <div className="flex flex-col w-full max-w-2xl mx-auto gap-6">
        <SettingsHeader
          title="LLM Providers"
          description="Configure and manage your connected AI language model providers."
          action={
            <Button
              size="sm"
              onClick={() => {
                setForm(emptyForm());
                setSheetOpen(true);
              }}
            >
              <Plus className="w-4 h-4 mr-1" />
              Add
            </Button>
          }
        />

        {isLoading ? (
          <div className="flex flex-col gap-3">
            {[...Array(3)].map((_, i) => (
              <Skeleton key={i} className="h-16 rounded-xl" />
            ))}
          </div>
        ) : providers.length === 0 ? (
          <EmptyState
            icon={Brain}
            title="No providers yet"
            description="Add an LLM provider to start configuring your AI models."
            actionLabel="Add Provider"
            actionIcon={Plus}
            onAction={() => {
              setForm(emptyForm());
              setSheetOpen(true);
            }}
          />
        ) : (
          <div className="flex flex-col gap-3">
            {providers.map((p) => (
              <LLMProviderCard
                key={p.id}
                provider={p}
                onDelete={handleDelete}
              />
            ))}
          </div>
        )}

        <LLMProviderSheet
          open={sheetOpen}
          onOpenChange={setSheetOpen}
          form={form}
          setForm={setForm}
          onSave={handleSave}
          isSaving={isSaving}
          mode="add"
        />
      </div>
    </div>
  );
}
