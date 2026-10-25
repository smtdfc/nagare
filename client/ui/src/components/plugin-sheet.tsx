import { FolderOpen, Loader2, Upload, X } from "lucide-react";
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
  mode?: "install" | "update";
  pluginName?: string;
  localPath: string;
  setLocalPath: (path: string) => void;
  onInstall: () => void;
  onBrowse?: () => void;
  onFileSelect?: (file: File) => void;
  isDesktop?: boolean;
  isInstalling: boolean;
  signatureFileName?: string;
  onSignatureFileSelect?: (file: File | null) => void;
}

export function PluginSheet({
  open,
  onOpenChange,
  mode = "install",
  pluginName,
  localPath,
  setLocalPath,
  onInstall,
  onBrowse,
  onFileSelect,
  isDesktop = true,
  isInstalling,
  signatureFileName,
  onSignatureFileSelect,
}: PluginSheetProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const signatureInputRef = useRef<HTMLInputElement>(null);
  const isUpdate = mode === "update";
  const showSignatureUpload = isUpdate || !isDesktop;

  const handleButtonClick = () => {
    if (isDesktop && !isUpdate && onBrowse) {
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
            {isUpdate
              ? "Update Plugin"
              : isDesktop
                ? "Install Local Plugin"
                : "Upload Plugin Package"}
          </SheetTitle>
          <SheetDescription className="text-xs text-muted-foreground">
            {isUpdate
              ? `Upload a new package${pluginName ? ` for "${pluginName}"` : ""}. Your sessions are preserved.`
              : isDesktop
                ? "Provide the absolute path to a local plugin package directory or archive."
                : "Upload your .nagare_plugin archive. You can also attach its signature file."}
          </SheetDescription>
        </SheetHeader>

        <div className="flex flex-col gap-4 px-4 py-2">
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium">
              {isUpdate
                ? "Plugin File *"
                : isDesktop
                  ? "Plugin Path *"
                  : "Plugin File *"}
            </label>
            <p className="text-xs text-muted-foreground">
              {isDesktop && !isUpdate
                ? "Absolute path to the plugin directory on this machine."
                : "Select the .nagare_plugin file from your computer."}
            </p>

            <input
              type="file"
              ref={fileInputRef}
              className="hidden"
              accept=".nagare_plugin"
              onChange={handleFileChange}
            />

            <div className="flex gap-2 mt-1">
              <Input
                placeholder={
                  isDesktop && !isUpdate
                    ? "/home/user/my-plugin"
                    : "Select file..."
                }
                value={localPath}
                readOnly={!isDesktop || isUpdate}
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
                {isDesktop && !isUpdate ? (
                  <FolderOpen className="w-4 h-4" />
                ) : (
                  <Upload className="w-4 h-4" />
                )}
              </Button>
            </div>
          </div>

          {showSignatureUpload ? (
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium">
                Signature File (optional)
              </label>
              <p className="text-xs text-muted-foreground">
                Upload the detached .sig file to verify the publisher. A signed
                archive without a matching trusted key installs with a warning.
              </p>

              <input
                type="file"
                ref={signatureInputRef}
                className="hidden"
                accept=".sig"
                onChange={(e) =>
                  onSignatureFileSelect?.(e.target.files?.[0] ?? null)
                }
              />

              <div className="flex gap-2 mt-1">
                <Input
                  placeholder="Select signature file (optional)..."
                  value={signatureFileName ?? ""}
                  readOnly
                  onClick={() => signatureInputRef.current?.click()}
                />
                <Button
                  type="button"
                  variant="outline"
                  size="icon"
                  className="shrink-0"
                  aria-label="Select signature file"
                  onClick={() => signatureInputRef.current?.click()}
                >
                  <Upload className="w-4 h-4" />
                </Button>
                {signatureFileName && (
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    className="shrink-0"
                    aria-label="Clear signature file"
                    onClick={() => {
                      onSignatureFileSelect?.(null);
                      if (signatureInputRef.current) {
                        signatureInputRef.current.value = "";
                      }
                    }}
                  >
                    <X className="w-4 h-4" />
                  </Button>
                )}
              </div>
            </div>
          ) : (
            <p className="text-xs text-muted-foreground">
              If the archive is signed, keep the .sig file next to it. The
              gateway picks it up automatically.
            </p>
          )}
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
            {isUpdate ? "Update" : "Install"}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
