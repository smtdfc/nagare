import * as React from "react";

import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ChatService, ModelProviderService } from "@nagare-app/services";
import type { ModelProvider, Session } from "@nagare-app/dtos";
import { showToastError } from "#/lib/toast";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Check, Loader2, SlidersHorizontal } from "lucide-react";

type ChatModelSelectorProps = {
  session: Session;
  onSessionUpdated: (session: Session) => void;
};

export function ChatModelSelector({
  session,
  onSessionUpdated,
}: ChatModelSelectorProps) {
  const [isOpen, setIsOpen] = React.useState(false);
  const [providers, setProviders] = React.useState<ModelProvider[]>([]);
  const [models, setModels] = React.useState<string[]>([]);
  const [providerID, setProviderID] = React.useState(
    session.currentModelProvider?.id ?? "",
  );
  const [model, setModel] = React.useState(session.currentModel ?? "");
  const [isLoading, setIsLoading] = React.useState(false);
  const [isSaving, setIsSaving] = React.useState(false);

  React.useEffect(() => {
    setProviderID(session.currentModelProvider?.id ?? "");
    setModel(session.currentModel ?? "");
  }, [session]);

  React.useEffect(() => {
    let isMounted = true;
    void ModelProviderService.list()
      .then((items) => {
        if (isMounted) setProviders(items ?? []);
      })
      .catch(showToastError);
    return () => {
      isMounted = false;
    };
  }, []);

  const loadModels = async (id: string) => {
    setIsLoading(true);
    try {
      setModels((await ModelProviderService.getModels({ id })) ?? []);
    } catch (error) {
      showToastError(error);
    } finally {
      setIsLoading(false);
    }
  };

  React.useEffect(() => {
    if (providerID) {
      void loadModels(providerID);
    }
  }, [providerID]);

  const handleProviderChange = (value: string | null) => {
    if (!value) return;
    setProviderID(value);
    setModel("");
  };

  const handleSave = async () => {
    if (!providerID || !model) return;
    setIsSaving(true);
    try {
      const updatedSession = await ChatService.updateChatSessionModelSettings(
        session.id,
        {
          currentModelProvider: providerID,
          currentModel: model,
        },
      );
      onSessionUpdated(updatedSession);
      setIsOpen(false);
    } catch (error) {
      showToastError(error);
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <Popover open={isOpen} onOpenChange={setIsOpen}>
      <PopoverTrigger
        render={
          <Button
            variant="ghost"
            size="sm"
            className="h-8 gap-1.5 px-2 data-open:bg-accent"
            aria-label="Change provider and model"
          />
        }
      >
        <SlidersHorizontal className="h-4 w-4" />
        <span className="hidden max-w-28 truncate sm:inline">
          {session.currentModelProvider?.name ?? "AI"}
        </span>
      </PopoverTrigger>
      <PopoverContent className="w-80" align="end">
        <div className="flex flex-col gap-4">
          <div>
            <h3 className="font-medium">Chat model</h3>
            <p className="text-xs text-muted-foreground">
              Choose the provider and model for this session.
            </p>
          </div>
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium">Provider</label>
            <Select value={providerID} onValueChange={handleProviderChange}>
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Select a provider">
                  {providers.find((provider) => provider.id === providerID)
                    ?.name ?? session.currentModelProvider?.name}
                </SelectValue>
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  {providers.map((provider) => (
                    <SelectItem key={provider.id} value={provider.id}>
                      {provider.name}
                    </SelectItem>
                  ))}
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium">Model</label>
            <div className="relative">
              {isLoading && (
                <Loader2 className="absolute right-8 top-2.5 z-10 h-4 w-4 animate-spin text-muted-foreground" />
              )}
              <Select
                value={model}
                disabled={isLoading || !providerID || models.length === 0}
                onValueChange={(value) => setModel(value ?? "")}
              >
                <SelectTrigger className="w-full">
                  <SelectValue placeholder="Select a model" />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    {models.map((item) => (
                      <SelectItem key={item} value={item}>
                        {item}
                      </SelectItem>
                    ))}
                  </SelectGroup>
                </SelectContent>
              </Select>
            </div>
          </div>
          <Button
            onClick={handleSave}
            disabled={isSaving || !providerID || !model}
          >
            {isSaving ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Check className="h-4 w-4" />
            )}
            Save
          </Button>
        </div>
      </PopoverContent>
    </Popover>
  );
}
