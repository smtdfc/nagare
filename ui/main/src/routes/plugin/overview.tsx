import { createColumns } from "#/components/columns/plugin-columns";
import { DataTable } from "#/components/data-table";
import { Button } from "#/components/ui/button";
import { Spinner } from "#/components/ui/spinner";
import type { PluginInfo } from "@nagare-agent/dto";
import { PluginService } from "#/services/plugin";
import { createFileRoute, useRouter } from "@tanstack/react-router";
import { Plus } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

export const Route = createFileRoute("/plugin/overview")({
  component: RouteComponent,
  staticData: {
    breadcrumb: "Plugin Overview",
  },
});

function RouteComponent() {
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [plugins, setPlugins] = useState<PluginInfo[]>([]);
  const router = useRouter();
  useEffect(() => {
    (async () => {
      try {
        const data = await PluginService.getAllPlugin();
        setPlugins(data || []);
      } catch (err) {
        console.error("Failed to load providers:", err);
      } finally {
        setIsLoading(false);
      }
    })();
  }, []);

  const columns = useMemo(() => {
    return createColumns(
      (plugin) => {},
      async (plugin) => {},
    );
  }, []);

  return (
    <div className="container mx-auto py-8 px-6 max-w-6xl">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Manage Plugins</h1>
          <p className="text-sm text-muted-foreground mt-1">
            Manage your plugins.
          </p>
        </div>
        <Button
          className="flex items-center gap-2 "
          onClick={() =>
            router.navigate({
              to: "/plugin/add",
            })
          }
        >
          <Plus className="size-4" />
          Add Plugin
        </Button>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center h-[50vh] w-full rounded-lg border border-dashed">
          <Spinner className="size-8 text-primary animate-spin" />
        </div>
      ) : (
        <div className=" overflow-hidden">
          <DataTable columns={columns} data={plugins} />
        </div>
      )}
    </div>
  );
}
