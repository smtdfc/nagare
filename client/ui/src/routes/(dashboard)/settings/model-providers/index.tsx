import { createFileRoute } from "@tanstack/react-router";
import { Plus, Brain } from "lucide-react";
import { useEffect, useState, useCallback } from "react";
import { Button } from "#/components/ui/button.tsx";
import { Skeleton } from "#/components/ui/skeleton.tsx";
import { CredentialService } from "@nagare-app/services";
import { showToastError } from "#/lib/toast.ts";
import { toast } from "#/components/ui/toast.tsx";
import type {
  Credential,
  AddModelProviderRequest,
} from "@nagare-app/dtos";
import { useModelProviders } from "#/hooks/use-model-providers.ts";
import { ModelProviderSheet } from "#/components/model-provider-sheet.tsx";
import { ModelProviderCard } from "#/components/model-provider-card.tsx";
import { EmptyState } from "#/components/empty-state";
import { SettingsHeader } from "#/components/settings-header";

export const Route = createFileRoute("/(dashboard)/settings/model-providers/")({
  component: RouteComponent,
  staticData: {
    breadcrumb: "Model Providers settings",
  },
});

const emptyForm = (): AddModelProviderRequest => ({
  name: "",
  compatible: "OpenAI",
  apiKey: "",
  credentialId: "",
  models: [],
  baseUrl: "",
});

function RouteComponent() {
  const providers = useModelProviders((s) => s.providers);
  const isLoading = useModelProviders((s) => s.isLoading);
  const isSaving = useModelProviders((s) => s.isSaving);
  const fetchProviders = useModelProviders((s) => s.fetchProviders);
  const addProvider = useModelProviders((s) => s.addProvider);
  const removeProvider = useModelProviders((s) => s.removeProvider);

  const [credentials, setCredentials] = useState<Credential[]>([]);
  const [isCredentialsLoading, setIsCredentialsLoading] = useState(true);
  const [sheetOpen, setSheetOpen] = useState(false);
  const [form, setForm] = useState<AddModelProviderRequest>(emptyForm());

  const loadCredentials = useCallback(async () => {
    setIsCredentialsLoading(true);
    try {
      setCredentials((await CredentialService.list()) ?? []);
    } catch (e) {
      showToastError(e);
    } finally {
      setIsCredentialsLoading(false);
    }
  }, []);

  useEffect(() => {
    void fetchProviders().catch(showToastError);
    void loadCredentials();
  }, [fetchProviders, loadCredentials]);

  const handleSave = async () => {
    if (!form.name.trim() || (!form.credentialId && !form.apiKey)) return;
    try {
      const created = await addProvider(form);
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
    }
  };

  const handleDelete = (id: string) => {
    void removeProvider(id).catch(showToastError);
  };

  return (
    <div className="flex flex-col items-center w-full p-6">
      <div className="flex flex-col w-full max-w-2xl mx-auto gap-6">
        <SettingsHeader
          title="Model Providers"
          description="Configure and manage your connected AI model providers."
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
            description="Add a model provider to start configuring your AI models."
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
              <ModelProviderCard
                key={p.id}
                provider={p}
                onDelete={handleDelete}
              />
            ))}
          </div>
        )}

        <ModelProviderSheet
          open={sheetOpen}
          onOpenChange={setSheetOpen}
          form={form}
          setForm={setForm}
          onSave={handleSave}
          isSaving={isSaving}
          credentials={credentials}
          isCredentialsLoading={isCredentialsLoading}
          mode="add"
        />
      </div>
    </div>
  );
}
