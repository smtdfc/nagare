import LLMProviderEditForm from "#/components/llm-provider-edit-form.tsx";
import { toast } from "#/components/ui/toast.tsx";
import type { Provider } from "@nagare-agent/dto";
import { ProviderService } from "#/services/provider.ts";
import { createFileRoute, useRouter } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { getErrorMessage } from "#/lib/error";

export const Route = createFileRoute("/settings/llm-provider/edit/$id")({
  component: RouteComponent,
  staticData: {
    breadcrumb: "Edit LLM Provider",
  },
});

function RouteComponent() {
  const router = useRouter();
  const { id } = Route.useParams();
  const [currentProvider, setCurrentProvider] = useState<Provider>();

  useEffect(() => {
    const fetchProvider = async () => {
      try {
        const details = await ProviderService.getProviderById(id);
        setCurrentProvider(details);
      } catch (e) {
        toast.add({
          type: "error",
          description: getErrorMessage(e),
          priority: "high",
        });
      }
    };

    fetchProvider();
  }, [id]);

  const handleSubmit = (provider: Provider) => {
    ProviderService.updateProvider({
      id: provider.id,
      compatible: provider.compatible,
      name: provider.name,
      base_url: provider.base_url,
      is_enable: provider.is_enable,
      available_models: provider.available_models,
      api_key: provider.api_key,
      model_name: "",
    })
      .then(() => {
        toast.add({
          type: "success",
          description: "Provider updated successfully",
          priority: "high",
        });
        router.history.back();
      })
      .catch((e) => {
        toast.add({
          type: "error",
          description: getErrorMessage(e),
          priority: "high",
        });
      });
  };

  return (
    <div className="container mx-auto py-8 px-6 max-w-6xl">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Edit Provider</h1>
        </div>
      </div>
      {currentProvider ? (
        <LLMProviderEditForm
          currentProvider={currentProvider}
          onSubmit={handleSubmit}
          onCancel={() => router.history.back()}
        />
      ) : (
        "Error"
      )}
    </div>
  );
}
