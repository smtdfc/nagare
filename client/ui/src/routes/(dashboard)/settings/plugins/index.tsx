import { createFileRoute } from "@tanstack/react-router";
import { Plus, Plug, Package } from "lucide-react";
import { useEffect, useState, useCallback } from "react";
import { Button } from "#/components/ui/button.tsx";
import { Skeleton } from "#/components/ui/skeleton.tsx";
import { envConfig, Environment, PluginService } from "@nagare-app/services";
import { showToastError } from "#/lib/toast.ts";
import { toast } from "#/components/ui/toast.tsx";
import type {
  InstallLocalPluginResponse,
  InstallPluginFromAttachmentResponse,
  Plugin,
  UpdatePluginResponse,
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
  const [selectedSignatureFile, setSelectedSignatureFile] =
    useState<File | null>(null);
  const [isInstalling, setIsInstalling] = useState(false);
  const [updatingPluginId, setUpdatingPluginId] = useState<string | null>(null);

  const isDesktop = envConfig.current === Environment.Desktop;
  const isUpdateMode = updatingPluginId !== null;
  const updatingPlugin = plugins.find((p) => p.id === updatingPluginId);

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

  const resetSheetState = useCallback(() => {
    setLocalPath("");
    setSelectedFile(null);
    setSelectedSignatureFile(null);
    setUpdatingPluginId(null);
  }, []);

  const handleSheetOpenChange = (open: boolean) => {
    setSheetOpen(open);
    if (!open) resetSheetState();
  };

  const openInstallSheet = () => {
    resetSheetState();
    setSheetOpen(true);
  };

  const handleUpdatePlugin = (plugin: Plugin) => {
    resetSheetState();
    setUpdatingPluginId(plugin.id);
    setSheetOpen(true);
  };

  const showMutationResult = (
    result:
      | InstallLocalPluginResponse
      | InstallPluginFromAttachmentResponse
      | UpdatePluginResponse
      | null,
    isUpdate: boolean,
  ) => {
    const plugin = result?.plugin ?? null;
    if (!plugin) {
      throw new Error(
        isUpdate ? "Failed to update plugin" : "Failed to install plugin",
      );
    }

    if (isUpdate && updatingPluginId) {
      setPlugins((prev) =>
        prev.map((p) => (p.id === updatingPluginId ? { ...p, ...plugin } : p)),
      );

      toast.add({
        title: "Updated",
        type: "success",
        priority: "high",
        description: `Plugin "${plugin.name}" has been updated to v${plugin.version}. Your sessions are preserved.`,
      });
    } else {
      setPlugins((prev) => [...prev, plugin]);

      toast.add({
        title: "Success",
        type: "success",
        priority: "high",
        description: "Plugin installed successfully.",
      });
    }

    if (result?.warning) {
      toast.add({
        title: "Plugin warning",
        type: "warning",
        priority: "high",
        description: result.warning,
      });
    }

    setSheetOpen(false);
    resetSheetState();
  };

  const uploadSignatureIfSelected = async () => {
    if (!selectedSignatureFile) {
      return undefined;
    }

    const { attachmentId } = await PluginService.upload(selectedSignatureFile);
    return attachmentId;
  };

  const handleInstall = async () => {
    const uploadFile = isUpdateMode || !isDesktop;

    if (uploadFile && !selectedFile) {
      toast.add({
        title: "No file selected",
        type: "error",
        priority: "high",
        description: "Please select a plugin package to upload.",
      });
      return;
    }

    if (!uploadFile && !localPath.trim()) {
      return;
    }

    setIsInstalling(true);
    try {
      if (uploadFile && selectedFile) {
        const { attachmentId } = await PluginService.upload(selectedFile);
        const signatureAttachmentId = await uploadSignatureIfSelected();

        const result = isUpdateMode
          ? await PluginService.update(attachmentId, signatureAttachmentId)
          : await PluginService.installFromAttachment(
              attachmentId,
              signatureAttachmentId,
            );

        showMutationResult(result, isUpdateMode);
        return;
      }

      const result = await PluginService.installLocal(localPath.trim());
      showMutationResult(result, false);
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
              onClick={openInstallSheet}
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
            onAction={openInstallSheet}
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
          onOpenChange={handleSheetOpenChange}
          mode={isUpdateMode ? "update" : "install"}
          pluginName={updatingPlugin?.name}
          localPath={localPath}
          setLocalPath={setLocalPath}
          onInstall={handleInstall}
          isInstalling={isInstalling}
          onBrowse={() => {}}
          isDesktop={isDesktop}
          onFileSelect={setSelectedFile}
          signatureFileName={selectedSignatureFile?.name}
          onSignatureFileSelect={setSelectedSignatureFile}
        />
      </div>
    </div>
  );
}
