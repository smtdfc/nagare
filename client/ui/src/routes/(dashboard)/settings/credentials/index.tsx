import { createFileRoute } from "@tanstack/react-router";
import { KeyRound, Plus } from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import type { Credential } from "@nagare-app/dtos";
import { CredentialService } from "@nagare-app/services";
import { Button } from "#/components/ui/button.tsx";
import { CredentialCard } from "#/components/credential-card.tsx";
import {
  CredentialSheet,
  type CredentialForm,
} from "#/components/credential-sheet.tsx";
import { EmptyState } from "#/components/empty-state.tsx";
import { SettingsHeader } from "#/components/settings-header.tsx";
import { Skeleton } from "#/components/ui/skeleton.tsx";
import { showToastError } from "#/lib/toast.ts";
import { toast } from "#/components/ui/toast.tsx";

export const Route = createFileRoute("/(dashboard)/settings/credentials/")({
  component: RouteComponent,
  staticData: {
    breadcrumb: "Credentials settings",
  },
});

const emptyForm = (): CredentialForm => ({
  name: "",
  apiKey: "",
});

function RouteComponent() {
  const [credentials, setCredentials] = useState<Credential[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [deletingID, setDeletingID] = useState<string | null>(null);
  const [sheetOpen, setSheetOpen] = useState(false);
  const [editingCredential, setEditingCredential] = useState<Credential | null>(
    null,
  );
  const [form, setForm] = useState<CredentialForm>(emptyForm());

  const loadCredentials = useCallback(async () => {
    setIsLoading(true);
    try {
      setCredentials((await CredentialService.list()) ?? []);
    } catch (error) {
      showToastError(error);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadCredentials();
  }, [loadCredentials]);

  const openCreateSheet = () => {
    setEditingCredential(null);
    setForm(emptyForm());
    setSheetOpen(true);
  };

  const openEditSheet = (credential: Credential) => {
    setEditingCredential(credential);
    setForm({ name: credential.name, apiKey: "" });
    setSheetOpen(true);
  };

  const handleSave = async () => {
    if (!form.name.trim() || (!editingCredential && !form.apiKey)) return;

    setIsSaving(true);
    try {
      if (editingCredential) {
        const updated = await CredentialService.update({
          id: editingCredential.id,
          name: form.name,
          apiKey: form.apiKey,
        });
        if (!updated) return;
        setCredentials((previous) =>
          previous.map((credential) =>
            credential.id === updated.id ? updated : credential,
          ),
        );
        toast.add({
          title: "Updated",
          type: "success",
          priority: "high",
          description: `Credential "${updated.name}" has been updated.`,
        });
      } else {
        const created = await CredentialService.add(form);
        if (!created) return;
        setCredentials((previous) => [...previous, created]);
        toast.add({
          title: "Added",
          type: "success",
          priority: "high",
          description: `Credential "${created.name}" has been added.`,
        });
      }
      setSheetOpen(false);
      setForm(emptyForm());
    } catch (error) {
      showToastError(error);
    } finally {
      setIsSaving(false);
    }
  };

  const handleDelete = async (id: string) => {
    setDeletingID(id);
    try {
      await CredentialService.delete(id);
      setCredentials((previous) =>
        previous.filter((credential) => credential.id !== id),
      );
      toast.add({
        title: "Deleted",
        type: "success",
        priority: "high",
        description: "Credential has been removed.",
      });
    } catch (error) {
      showToastError(error);
    } finally {
      setDeletingID(null);
    }
  };

  return (
    <div className="flex w-full flex-col items-center p-6">
      <div className="mx-auto flex w-full max-w-2xl flex-col gap-6">
        <SettingsHeader
          title="Credentials"
          description="Store reusable provider credentials without exposing their API keys."
          action={
            <Button size="sm" onClick={openCreateSheet}>
              <Plus className="mr-1 h-4 w-4" />
              Add
            </Button>
          }
        />

        {isLoading ? (
          <div className="flex flex-col gap-3">
            {[...Array(3)].map((_, index) => (
              <Skeleton key={index} className="h-16 rounded-xl" />
            ))}
          </div>
        ) : credentials.length === 0 ? (
          <EmptyState
            icon={KeyRound}
            title="No credentials yet"
            description="Add a credential to reuse it across your LLM providers."
            actionLabel="Add Credential"
            actionIcon={Plus}
            onAction={openCreateSheet}
          />
        ) : (
          <div className="flex flex-col gap-3">
            {credentials.map((credential) => (
              <CredentialCard
                key={credential.id}
                credential={credential}
                onEdit={openEditSheet}
                onDelete={handleDelete}
                isDeleting={deletingID === credential.id}
              />
            ))}
          </div>
        )}

        <CredentialSheet
          open={sheetOpen}
          onOpenChange={setSheetOpen}
          form={form}
          setForm={setForm}
          onSave={handleSave}
          isSaving={isSaving}
          mode={editingCredential ? "edit" : "add"}
        />
      </div>
    </div>
  );
}
