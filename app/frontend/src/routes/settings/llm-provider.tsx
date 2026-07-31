import { DataTable } from '#/components/data-table.tsx';
import { Button } from '#/components/ui/button.tsx';
import { Spinner } from '#/components/ui/spinner.tsx';
import type { Provider } from '#/dto/api.ts';
import { ProviderService } from '#/services/provider.ts';
import { createFileRoute } from '@tanstack/react-router'
import { useEffect, useMemo, useState } from 'react';
import { Plus, } from 'lucide-react';
import { createColumns } from '#/components/columns/provider-columns.tsx';

export const Route = createFileRoute('/settings/llm-provider')({
  component: RouteComponent,
  staticData: {
    breadcrumb: 'LLM Provider settings',
  },
})


function RouteComponent() {
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [providers, setProviders] = useState<Provider[]>([]);

  useEffect(() => {
    (async () => {
      try {
        const data = await ProviderService.getListProvider();
        setProviders(data.providers || []);
      } catch (err) {
        console.error("Failed to load providers:", err);
      } finally {
        setIsLoading(false);
      }
    })();
  }, []);

  const columns = useMemo(() => {
    return createColumns(
      (provider) => {
        console.log("Edit provider:", provider);


      },
      (provider) => {
        console.log("Delete provider:", provider);

      }
    );
  }, []);

  return (
    <div className="container mx-auto py-8 px-6 max-w-6xl">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">LLM Providers</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Manage your AI provider endpoints and configurations.
          </p>
        </div>
        <Button className="flex items-center gap-2 ">
          <Plus className="size-4" />
          Add Provider
        </Button>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center h-[50vh] w-full rounded-lg border border-dashed">
          <Spinner className="size-8 text-primary animate-spin" />
        </div>
      ) : (
        <div className=" overflow-hidden">
          <DataTable columns={columns} data={providers} />
        </div>
      )}
    </div>
  )
}