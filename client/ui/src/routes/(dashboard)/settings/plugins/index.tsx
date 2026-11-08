import { createFileRoute } from "@tanstack/react-router";
import { Plus, Plug, Package } from "lucide-react";
import { Button } from "#/components/ui/button.tsx";
import { Skeleton } from "#/components/ui/skeleton.tsx";
import { PluginCard } from "#/components/plugin-card";
import { PluginSheet } from "#/components/plugin-sheet";
import { EmptyState } from "#/components/empty-state";
import { SettingsHeader } from "#/components/settings-header";
import { usePluginManagement } from "#/hooks/use-plugin-management.ts";

export const Route = createFileRoute("/(dashboard)/settings/plugins/")({
  component: RouteComponent,
  staticData: {
    breadcrumb: "Plugins settings",
  },
});

function RouteComponent() {
  const {
    plugins,
    isLoading,
    sheetOpen,
    localPath,
    setLocalPath,
    selectedSignatureFile,
    setSelectedSignatureFile,
    isInstalling,
    isDesktop,
    isUpdateMode,
    updatingPlugin,
    handleToggle,
    handleUninstall,
    handleSheetOpenChange,
    openInstallSheet,
    handleUpdatePlugin,
    handleInstall,
    setSelectedFile,
  } = usePluginManagement();

  return (
    <div className="flex flex-col items-center w-full p-6">
      <div className="flex flex-col w-full max-w-2xl mx-auto gap-6">
        <SettingsHeader
          title="Plugins"
          description="Manage your installed plugins. You can activate, deactivate, or uninstall plugins as needed."
          action={
            <Button variant="secondary" size="sm" onClick={openInstallSheet}>
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
