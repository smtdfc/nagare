import { FolderOpen, Loader2, Upload } from "lucide-react";
import { useRef } from "react";
import { Button } from "#/components/ui/button.tsx";
import { Input } from "#/components/ui/input.tsx";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
  SheetFooter,
} from "#/components/ui/sheet.tsx";

interface PluginSheetProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  localPath: string;
  setLocalPath: (path: string) => void;
  onInstall: () => void;
  onBrowse?: () => void;
  onFileSelect?: (file: File) => void;
  isDesktop?: boolean;
  isInstalling: boolean;
}

export function PluginSheet({
  open,
  onOpenChange,
  localPath,
  setLocalPath,
  onInstall,
  onBrowse,
  onFileSelect,
  isDesktop = true,
  isInstalling,
}: PluginSheetProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleButtonClick = () => {
    if (isDesktop && onBrowse) {
      onBrowse();
    } else {
      fileInputRef.current?.click();
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setLocalPath(file.name);
      if (onFileSelect) {
        onFileSelect(file);
      }
    }
  };

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-full sm:max-w-md overflow-y-auto">
        <SheetHeader>
          <SheetTitle className="text-lg font-semibold">
            {isDesktop ? "Install Local Plugin" : "Upload Plugin Package"}
          </SheetTitle>
          <SheetDescription className="text-xs text-muted-foreground">
            {isDesktop
              ? "Provide the absolute path to a local plugin package directory or archive."
              : "Upload your plugin archive file (.zip or package format)."}
          </SheetDescription>
        </SheetHeader>

        <div className="flex flex-col gap-4 px-4 py-2">
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium">
              {isDesktop ? "Plugin Path *" : "Plugin File *"}
            </label>
            <p className="text-xs text-muted-foreground">
              {isDesktop
                ? "Absolute path to the plugin directory on this machine."
                : "Select the plugin file from your computer."}
            </p>

            <input
              type="file"
              ref={fileInputRef}
              className="hidden"
              accept=".zip,.tar,.gz"
              onChange={handleFileChange}
            />

            <div className="flex gap-2 mt-1">
              <Input
                placeholder={
                  isDesktop ? "/home/user/my-plugin" : "Select file..."
                }
                value={localPath}
                readOnly={!isDesktop}
                onChange={(e) => setLocalPath(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && onInstall()}
              />
              <Button
                type="button"
                variant="outline"
                size="icon"
                className="shrink-0"
                aria-label="Browse"
                onClick={handleButtonClick}
              >
                {isDesktop ? (
                  <FolderOpen className="w-4 h-4" />
                ) : (
                  <Upload className="w-4 h-4" />
                )}
              </Button>
            </div>
          </div>
        </div>

        <SheetFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            onClick={onInstall}
            disabled={isInstalling || !localPath.trim()}
          >
            {isInstalling && <Loader2 className="w-4 h-4 animate-spin mr-1" />}
            Install
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
