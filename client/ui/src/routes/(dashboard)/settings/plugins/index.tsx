import { createFileRoute } from "@tanstack/react-router";
import { Plus, Plug, Package } from "lucide-react";
import { useEffect, useState, useCallback, useRef } from "react";
import { Button } from "#/components/ui/button.tsx";
import { Skeleton } from "#/components/ui/skeleton.tsx";
import { envConfig, Environment, PluginService } from "@nagare-app/services";
import { showToastError } from "#/lib/toast.ts";
import { toast } from "#/components/ui/toast.tsx";
import type {
  InstallLocalPluginResponse,
  InstallPluginFromAttachmentResponse,
  Plugin,
} from "@nagare-app/dtos";
import { PluginCard } from "#/components/plugin-card";
import { PluginSheet } from "#/components/plugin-sheet";
import { EmptyState } from "#/components/empty-state";
import { SettingsHeader } from "#/components/settings-header";

export const Route = createFileRoute("/(dashboard)/settings/plugins/")({
  component: RouteComponent,
  staticData: {
    breadcrumb: "Plugins settings",
  },
});

function RouteComponent() {
  const [plugins, setPlugins] = useState<Plugin[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [sheetOpen, setSheetOpen] = useState(false);
  const [localPath, setLocalPath] = useState("");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isInstalling, setIsInstalling] = useState(false);
  const [updatingPluginId, setUpdatingPluginId] = useState<string | null>(null);
  const [isSubmittingUpdate, setIsSubmittingUpdate] = useState(false);
  const updateFileInputRef = useRef<HTMLInputElement>(null);

  const isDesktop = envConfig.current === Environment.Desktop;

  const loadPlugins = useCallback(async () => {
    setIsLoading(true);
    try {
      const list = await PluginService.list();
      setPlugins(list ?? []);
    } catch (e) {
      showToastError(e);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadPlugins();
  }, [loadPlugins]);

  const handleToggle = (id: string, isActive: boolean) => {
    setPlugins((prev) =>
      prev.map((p) => (p.id === id ? { ...p, isActive } : p)),
    );
  };

  const handleUninstall = (id: string) => {
    setPlugins((prev) => prev.filter((p) => p.id !== id));
  };

  const handleUpdatePlugin = (plugin: Plugin) => {
    // Reset input value trước khi mở dialog để người dùng có thể chọn lại
    // cùng một file (tránh trình duyệt không trigger onChange khi tên file trùng).
    if (updateFileInputRef.current) {
      updateFileInputRef.current.value = "";
    }
    setUpdatingPluginId(plugin.id);
    updateFileInputRef.current?.click();
  };

  const handleUpdateFileChange = async (
    e: React.ChangeEvent<HTMLInputElement>,
  ) => {
    const file = e.target.files?.[0];
    if (!file || !updatingPluginId) return;

    setIsSubmittingUpdate(true);
    try {
      const { attachmentId } = await PluginService.upload(file);
      const result = await PluginService.update(attachmentId);
      const updated = result?.plugin;
      if (!updated) throw new Error("Failed to update plugin");

      setPlugins((prev) =>
        prev.map((p) =>
          p.id === updatingPluginId ? { ...p, ...updated } : p,
        ),
      );

      toast.add({
        title: "Updated",
        type: "success",
        priority: "high",
        description: `Plugin "${updated.name}" has been updated to v${updated.version}. Your sessions are preserved.`,
      });

      if (result?.warning) {
        toast.add({
          title: "Plugin warning",
          type: "warning",
          priority: "high",
          description: result.warning,
        });
      }
    } catch (e) {
      showToastError(e);
    } finally {
      setIsSubmittingUpdate(false);
      setUpdatingPluginId(null);
      if (updateFileInputRef.current) {
        updateFileInputRef.current.value = "";
      }
    }
  };

  const handleInstall = async () => {
    if (isDesktop && !localPath.trim()) return;
    if (!isDesktop && !selectedFile) {
      toast.add({
        title: "No file selected",
        type: "error",
        priority: "high",
        description: "Please select a plugin package to upload.",
      });
      return;
    }

    setIsInstalling(true);
    try {
      let result:
        | InstallLocalPluginResponse
        | InstallPluginFromAttachmentResponse
        | null = null;

      if (isDesktop) {
        result = await PluginService.installLocal(localPath.trim());
      } else if (selectedFile) {
        const { attachmentId } = await PluginService.upload(selectedFile);
        result = await PluginService.installFromAttachment(attachmentId);
      }

      const plugin = result?.plugin ?? null;
      if (!plugin) throw new Error("Failed to install plugin");

      setPlugins((prev) => [...prev, plugin]);
      setSheetOpen(false);
      setLocalPath("");
      setSelectedFile(null);

      toast.add({
        title: "Success",
        type: "success",
        priority: "high",
        description: "Plugin installed successfully.",
      });

      if (result?.warning) {
        toast.add({
          title: "Plugin warning",
          type: "warning",
          priority: "high",
          description: result.warning,
        });
      }
    } catch (e) {
      showToastError(e);
    } finally {
      setIsInstalling(false);
    }
  };

  return (
    <div className="flex flex-col items-center w-full p-6">
      <div className="flex flex-col w-full max-w-2xl mx-auto gap-6">
        <SettingsHeader
          title="Plugins"
          description="Manage your installed plugins. You can activate, deactivate, or uninstall plugins as needed."
          action={
            <Button
              variant="secondary"
              size="sm"
              onClick={() => setSheetOpen(true)}
            >
              <Plus className="w-4 h-4 mr-2" />
              Install Plugin
            </Button>
          }
        />

        {isLoading ? (
          <div className="flex flex-col gap-3">
            {[...Array(4)].map((_, i) => (
              <Skeleton key={i} className="h-20 rounded-xl" />
            ))}
          </div>
        ) : plugins.length === 0 ? (
          <EmptyState
            icon={Plug}
            title="No plugins installed"
            description="Install a local plugin to extend Nagare's capabilities."
            actionLabel="Install Plugin"
            actionIcon={Package}
            onAction={() => setSheetOpen(true)}
          />
        ) : (
          <div className="flex flex-col gap-3">
            {plugins.map((p) => (
              <PluginCard
                key={p.id}
                plugin={p}
                onToggle={handleToggle}
                onUninstall={handleUninstall}
                onUpdate={handleUpdatePlugin}
              />
            ))}
          </div>
        )}

        <PluginSheet
          open={sheetOpen}
          onOpenChange={setSheetOpen}
          localPath={localPath}
          setLocalPath={setLocalPath}
          onInstall={handleInstall}
          isInstalling={isInstalling}
          onBrowse={() => {}}
          isDesktop={isDesktop}
          onFileSelect={setSelectedFile}
        />

        {/* Hidden file input used to pick a new package when updating a plugin. */}
        <input
          ref={updateFileInputRef}
          type="file"
          className="hidden"
          accept=".nagare_plugin,.zip"
          onChange={handleUpdateFileChange}
        />
      </div>
    </div>
  );
}
