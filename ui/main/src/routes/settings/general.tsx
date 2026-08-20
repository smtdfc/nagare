import { Button } from "#/components/ui/button";
import { Field, FieldGroup, FieldLabel, FieldSet } from "#/components/ui/field";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "#/components/ui/select";
import { Spinner } from "#/components/ui/spinner";
import { toast } from "#/components/ui/toast";
import type { GeneralSettings, ProviderInfo } from "@nagare-agent/dto";
import { ProviderService } from "#/services/provider";
import { SettingsService } from "#/services/settings";
import { createFileRoute, useRouter } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { getErrorMessage } from "#/lib/error";

export const Route = createFileRoute("/settings/general")({
  component: RouteComponent,
  staticData: {
    breadcrumb: "General settings",
  },
});

function RouteComponent() {
  const router = useRouter();
  const [currentConf, setCurrentConf] = useState<Partial<GeneralSettings>>();
  const [isLoading, setIsLoading] = useState(true);
  const [providers, setProviders] = useState<ProviderInfo[]>([]);
  const [models, setModels] = useState<string[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const settings = await SettingsService.getGeneralSettings();
        setCurrentConf(settings!);
        const data = await ProviderService.getListProvider();
        setProviders(data.providers || []);
        const currentProvider = data.providers?.find(
          (p) => p.id === settings!.current_provider,
        );
        setModels(currentProvider?.available_models || []);
      } catch (err) {
        toast.add({
          title: "Error",
          description: getErrorMessage(err),
          type: "error",
        });
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, []);

  useEffect(() => {
    const currentProvider = providers?.find(
      (p) => p.id === currentConf!.current_provider,
    );
    setModels(currentProvider?.available_models || []);
  }, [currentConf]);

  const handleSubmit = () => {
    SettingsService.saveGeneralSettings({
      settings: currentConf! as any,
    })
      .then(() => {
        toast.add({
          title: "Success",
          description: "Settings saved successfully",
          type: "success",
        });
      })
      .catch((err) => {
        toast.add({
          title: "Error",
          description: getErrorMessage(err),
          type: "error",
        });
      });
  };
  const handleCancel = () => {
    router.history.back();
  };

  return (
    <div className="container mx-auto py-8 px-6 max-w-6xl">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">
            General settings
          </h1>
          <p className="text-sm text-muted-foreground mt-1">General settings</p>
        </div>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center h-[50vh] w-full rounded-lg">
          <Spinner className="size-8 text-primary animate-spin" />
        </div>
      ) : (
        <div className=" overflow-hidden">
          <form onSubmit={(e) => e.preventDefault()}>
            <FieldGroup>
              <FieldSet>
                <FieldGroup>
                  <Field>
                    <FieldLabel htmlFor="comp-base-url">
                      Current Provider:
                    </FieldLabel>
                    <Select
                      value={currentConf?.current_provider}
                      onValueChange={(value) =>
                        setCurrentConf({
                          ...currentConf,
                          current_provider: value!,
                        })
                      }
                    >
                      <SelectTrigger className="w-45">
                        <SelectValue placeholder="Select a provider">
                          {providers.find(
                            (p) => p.id === currentConf?.current_provider,
                          )?.name || "Select a provider"}
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
                  </Field>

                  <Field>
                    <FieldLabel htmlFor="comp-base-url">
                      Current Model:
                    </FieldLabel>
                    <Select
                      defaultValue={models.find(
                        (m) => m === currentConf?.current_model,
                      )}
                      onValueChange={(value) =>
                        setCurrentConf({
                          ...currentConf,
                          current_model: value!,
                        })
                      }
                    >
                      <SelectTrigger className="w-45">
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
                  </Field>
                  <div className="flex items-center justify-end gap-3 pt-4">
                    <Button
                      type="button"
                      variant="outline"
                      onClick={handleCancel}
                    >
                      Cancel
                    </Button>
                    <Button onClick={handleSubmit}>Save</Button>
                  </div>
                </FieldGroup>
              </FieldSet>
            </FieldGroup>
          </form>
        </div>
      )}
    </div>
  );
}
