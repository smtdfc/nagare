import { createFileRoute } from "@tanstack/react-router";
import { Save, Loader2 } from "lucide-react";
import { useEffect, useState, useCallback } from "react";
import { Button } from "#/components/ui/button.tsx";
import { LLMProviderService, SettingsService } from "@nagare-app/services";
import { showToastError } from "#/lib/toast.ts";
import { toast } from "#/components/ui/toast.tsx";
import type { GeneralSettings, LLMProvider } from "@nagare-app/dtos";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "#/components/ui/select.tsx";
import { Skeleton } from "#/components/ui/skeleton";
import { Separator } from "#/components/ui/separator";
import { SettingsHeader } from "#/components/settings-header";

export const Route = createFileRoute("/(dashboard)/settings/general/")({
  component: RouteComponent,
  staticData: {
    breadcrumb: "General settings",
  },
});

function RouteComponent() {
  const [settings, setSettings] = useState<GeneralSettings | null>(null);
  const [providers, setProviders] = useState<LLMProvider[]>([]);
  const [models, setModels] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [isLoadingModels, setIsLoadingModels] = useState(false);

  const loadData = useCallback(async () => {
    setIsLoading(true);
    try {
      const [generalSettings, providerList] = await Promise.all([
        SettingsService.getGeneralSettings(),
        LLMProviderService.list(),
      ]);
      setSettings(generalSettings);
      setProviders(providerList);

      if (generalSettings?.currentProvider) {
        const providerModels = await LLMProviderService.getModels({
          id: generalSettings.currentProvider,
        });
        setModels(providerModels);
      }
    } catch (e) {
      showToastError(e);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadData();
  }, [loadData]);

  const handleProviderChange = async (providerId: string | null) => {
    if (!providerId) return;
    setSettings((prev) =>
      prev ? { ...prev, currentProvider: providerId, currentModel: "" } : prev,
    );
    setModels([]);

    setIsLoadingModels(true);
    try {
      const providerModels = await LLMProviderService.getModels({
        id: providerId,
      });
      setModels(providerModels);
    } catch (e) {
      showToastError(e);
    } finally {
      setIsLoadingModels(false);
    }
  };

  const handleModelChange = (model: string | null) => {
    setSettings((prev) =>
      prev ? { ...prev, currentModel: model ?? "" } : prev,
    );
  };

  const handleSave = async () => {
    if (!settings) return;
    setIsSaving(true);
    try {
      await SettingsService.setGeneralSettings(settings);
      toast.add({
        title: "Saved",
        type: "success",
        priority: "high",
        description: "General settings updated successfully.",
      });
    } catch (e) {
      showToastError(e);
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <div className="flex flex-col items-center w-full p-6">
      <div className="flex flex-col w-full max-w-2xl mx-auto gap-6">
        <SettingsHeader
          title="General Settings"
          description="Manage your default conversation settings and active AI models."
        />

        {isLoading ? (
          <div className="flex flex-col gap-5">
            <Skeleton className="h-16 rounded-xl" />
            <Skeleton className="h-16 rounded-xl" />
          </div>
        ) : (
          <div className="flex flex-col gap-5">
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium">
                Current LLM provider
              </label>
              <p className="text-xs text-muted-foreground">
                Select the active LLM provider to use for conversations.
              </p>
              <Select
                value={settings?.currentProvider ?? ""}
                onValueChange={handleProviderChange}
              >
                <SelectTrigger className="w-full mt-1.5">
                  <SelectValue placeholder="Select a provider">
                    {providers.find((p) => p.id === settings?.currentProvider)
                      ?.name ?? "Select a provider"}
                  </SelectValue>
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    {providers.map((item) => (
                      <SelectItem key={item.id} value={item.id}>
                        {item.name}
                      </SelectItem>
                    ))}
                  </SelectGroup>
                </SelectContent>
              </Select>
            </div>

            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium">Current model</label>
              <p className="text-xs text-muted-foreground">
                Select the default model for the chosen provider.
              </p>
              <div className="relative mt-1.5">
                {isLoadingModels && (
                  <Loader2 className="absolute right-8 top-2.5 w-4 h-4 animate-spin text-muted-foreground z-10" />
                )}
                <Select
                  value={settings?.currentModel ?? ""}
                  disabled={isLoadingModels || models.length === 0}
                  onValueChange={handleModelChange}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="Select a model" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      {models.map((m) => (
                        <SelectItem key={m} value={m}>
                          {m}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="flex justify-end pt-5">
              <Button onClick={handleSave} disabled={isSaving}>
                {isSaving ? (
                  <Loader2 className="w-4 h-4 animate-spin" />
                ) : (
                  <Save className="w-4 h-4" />
                )}
                Save Changes
              </Button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
