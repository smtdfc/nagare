import { useCallback, useEffect, useState } from "react";
import { ModelProviderService, SettingsService } from "@nagare-app/services";
import type { GeneralSettings, ModelProvider } from "@nagare-app/dtos";
import { showToastError } from "#/lib/toast.ts";
import { toast } from "#/components/ui/toast.tsx";

export const USE_CONVERSATION_MODEL = "__conversation_model__";

export function useGeneralSettings() {
  const [settings, setSettings] = useState<GeneralSettings | null>(null);
  const [providers, setProviders] = useState<ModelProvider[]>([]);
  const [models, setModels] = useState<string[]>([]);
  const [summaryModels, setSummaryModels] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [isLoadingModels, setIsLoadingModels] = useState(false);
  const [isLoadingSummaryModels, setIsLoadingSummaryModels] = useState(false);

  const loadData = useCallback(async () => {
    setIsLoading(true);
    try {
      const [generalSettings, providerList] = await Promise.all([
        SettingsService.getGeneralSettings(),
        ModelProviderService.list(),
      ]);
      setSettings(generalSettings);
      setProviders(providerList);

      const [defaultModels, summaryModels] = await Promise.all([
        generalSettings?.defaultModelProvider
          ? ModelProviderService.getModels({
              id: generalSettings.defaultModelProvider,
            })
          : Promise.resolve([]),
        generalSettings?.summaryModelProvider
          ? ModelProviderService.getModels({
              id: generalSettings.summaryModelProvider,
            })
          : Promise.resolve([]),
      ]);
      setModels(defaultModels);
      setSummaryModels(summaryModels);
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
      prev
        ? { ...prev, defaultModelProvider: providerId, defaultModel: "" }
        : prev,
    );
    setModels([]);

    setIsLoadingModels(true);
    try {
      const providerModels = await ModelProviderService.getModels({
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
      prev ? { ...prev, defaultModel: model ?? "" } : prev,
    );
  };

  const handleSummaryProviderChange = async (providerId: string | null) => {
    if (!providerId) return;

    if (providerId === USE_CONVERSATION_MODEL) {
      setSettings((prev) =>
        prev ? { ...prev, summaryModelProvider: "", summaryModel: "" } : prev,
      );
      setSummaryModels([]);
      return;
    }

    setSettings((prev) =>
      prev
        ? { ...prev, summaryModelProvider: providerId, summaryModel: "" }
        : prev,
    );
    setSummaryModels([]);

    setIsLoadingSummaryModels(true);
    try {
      const providerModels = await ModelProviderService.getModels({
        id: providerId,
      });
      setSummaryModels(providerModels);
    } catch (e) {
      showToastError(e);
    } finally {
      setIsLoadingSummaryModels(false);
    }
  };

  const handleSummaryModelChange = (model: string | null) => {
    setSettings((prev) =>
      prev ? { ...prev, summaryModel: model ?? "" } : prev,
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

  return {
    settings,
    providers,
    models,
    summaryModels,
    isLoading,
    isSaving,
    isLoadingModels,
    isLoadingSummaryModels,
    handleProviderChange,
    handleModelChange,
    handleSummaryProviderChange,
    handleSummaryModelChange,
    handleSave,
  };
}
