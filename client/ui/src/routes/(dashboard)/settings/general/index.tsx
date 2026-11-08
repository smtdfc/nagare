import { createFileRoute } from "@tanstack/react-router";
import { Save, Loader2 } from "lucide-react";
import { Button } from "#/components/ui/button.tsx";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "#/components/ui/select.tsx";
import { Skeleton } from "#/components/ui/skeleton";
import { SettingsHeader } from "#/components/settings-header";
import {
  USE_CONVERSATION_MODEL,
  useGeneralSettings,
} from "#/hooks/use-general-settings.ts";

export const Route = createFileRoute("/(dashboard)/settings/general/")({
  component: RouteComponent,
  staticData: {
    breadcrumb: "General settings",
  },
});

function RouteComponent() {
  const {
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
  } = useGeneralSettings();

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
                Default model provider
              </label>
              <p className="text-xs text-muted-foreground">
                Select the active model provider to use for conversations.
              </p>
              <Select
                value={settings?.defaultModelProvider ?? ""}
                onValueChange={handleProviderChange}
              >
                <SelectTrigger className="w-full mt-1.5">
                  <SelectValue placeholder="Select a provider">
                    {providers.find(
                      (p) => p.id === settings?.defaultModelProvider,
                    )?.name ?? "Select a provider"}
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
              <label className="text-sm font-medium">Default model</label>
              <p className="text-xs text-muted-foreground">
                Select the default model for the chosen provider.
              </p>
              <div className="relative mt-1.5">
                {isLoadingModels && (
                  <Loader2 className="absolute right-8 top-2.5 w-4 h-4 animate-spin text-muted-foreground z-10" />
                )}
                <Select
                  value={settings?.defaultModel ?? ""}
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

            <div className="flex flex-col gap-5 border-t pt-5">
              <div className="flex flex-col gap-1.5">
                <label className="text-sm font-medium">
                  Conversation summarization
                </label>
                <p className="text-xs text-muted-foreground">
                  Select the model used to summarize conversations when
                  compacting. Leave it on the conversation model to reuse the
                  model of the current session.
                </p>
              </div>

              <div className="flex flex-col gap-1.5">
                <label className="text-sm font-medium">
                  Summary model provider
                </label>
                <p className="text-xs text-muted-foreground">
                  Select the provider for the summarization model.
                </p>
                <Select
                  value={
                    settings?.summaryModelProvider || USE_CONVERSATION_MODEL
                  }
                  onValueChange={handleSummaryProviderChange}
                >
                  <SelectTrigger className="w-full mt-1.5">
                    <SelectValue placeholder="Use conversation model">
                      {settings?.summaryModelProvider
                        ? (providers.find(
                            (p) => p.id === settings.summaryModelProvider,
                          )?.name ?? "Select a provider")
                        : "Use conversation model"}
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    <SelectGroup>
                      <SelectItem value={USE_CONVERSATION_MODEL}>
                        Use conversation model
                      </SelectItem>
                      {providers.map((item) => (
                        <SelectItem key={item.id} value={item.id}>
                          {item.name}
                        </SelectItem>
                      ))}
                    </SelectGroup>
                  </SelectContent>
                </Select>
              </div>

              {settings?.summaryModelProvider && (
                <div className="flex flex-col gap-1.5">
                  <label className="text-sm font-medium">Summary model</label>
                  <p className="text-xs text-muted-foreground">
                    Select the model used to summarize the conversation.
                  </p>
                  <div className="relative mt-1.5">
                    {isLoadingSummaryModels && (
                      <Loader2 className="absolute right-8 top-2.5 w-4 h-4 animate-spin text-muted-foreground z-10" />
                    )}
                    <Select
                      value={settings?.summaryModel ?? ""}
                      disabled={
                        isLoadingSummaryModels || summaryModels.length === 0
                      }
                      onValueChange={handleSummaryModelChange}
                    >
                      <SelectTrigger className="w-full">
                        <SelectValue placeholder="Select a model" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectGroup>
                          {summaryModels.map((m) => (
                            <SelectItem key={m} value={m}>
                              {m}
                            </SelectItem>
                          ))}
                        </SelectGroup>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
              )}
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
