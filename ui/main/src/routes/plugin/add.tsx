import { createFileRoute } from "@tanstack/react-router";
import { useState, useRef } from "react";
import {
  UploadCloud,
  FileBox,
  CheckCircle2,
  AlertCircle,
  Loader2,
  FolderOpen,
} from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Input } from "@/components/ui/input";
import { toast } from "#/components/ui/toast";
import { OpenPluginSelectDialog } from "@nagare-agent/service-bindings";
import { PluginService } from "#/services/plugin";
import { getErrorMessage } from "#/lib/error";

export const Route = createFileRoute("/plugin/add")({
  component: RouteComponent,
  staticData: {
    breadcrumb: "Install Plugin",
  },
});

function RouteComponent() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [filePath, setFilePath] = useState<string>("");
  const [isUploading, setIsUploading] = useState(false);
  const [uploadStatus, setUploadStatus] = useState<
    "idle" | "success" | "error"
  >("idle");
  const [error, setError] = useState<string | null>(null);

  const fileInputRef = useRef<HTMLInputElement>(null);

  const isDesktop =
    typeof window !== "undefined" && window.NagareUI === "desktop";

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      setSelectedFile(file);
      setUploadStatus("idle");
    }
  };

  const handleTriggerFileSelect = () => {
    fileInputRef.current?.click();
  };

  const handleOpenDesktopFileDialog = async () => {
    try {
      const path = await OpenPluginSelectDialog();
      if (path) {
        setFilePath(path);
        setUploadStatus("idle");
      }
    } catch (err) {
      toast.add({
        title: "Error",
        description: getErrorMessage(err),
        type: "error",
      });
    }
  };

  const handleInstallPlugin = async () => {
    if (isDesktop ? !filePath.trim() : !selectedFile) return;

    setIsUploading(true);
    setUploadStatus("idle");

    try {
      if (isDesktop) {
        await PluginService.installLocalPlugin({ path: filePath });
      }
      setUploadStatus("success");
    } catch (err) {
      toast.add({
        title: "Error",
        description: getErrorMessage(err),
        type: "error",
      });
      setUploadStatus("error");
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <div className="flex justify-center items-center min-h-[80vh] p-4">
      <Card className="w-full max-w-xl">
        <CardHeader>
          <CardTitle className="text-2xl font-bold flex items-center gap-2">
            <FileBox className="w-6 h-6 text-primary" />
            Install Plugin
          </CardTitle>
        </CardHeader>

        <CardContent className="space-y-6">
          {isDesktop ? (
            <div className="space-y-2">
              <label className="text-sm font-medium text-foreground flex items-center gap-2">
                <FolderOpen className="w-4 h-4 text-primary" />
                Plugin path (.nagare_plugin)
              </label>
              <div className="flex gap-2">
                <Input
                  type="text"
                  placeholder="..."
                  value={filePath}
                  onChange={(e) => {
                    setFilePath(e.target.value);
                    setUploadStatus("idle");
                  }}
                  disabled={isUploading}
                />
                <Button
                  type="button"
                  variant="secondary"
                  onClick={handleOpenDesktopFileDialog}
                  disabled={isUploading}
                  className="shrink-0 gap-2"
                >
                  <FolderOpen className="w-4 h-4" />
                  Browse
                </Button>
              </div>
            </div>
          ) : (
            <>
              <input
                type="file"
                ref={fileInputRef}
                onChange={handleFileChange}
                accept=".nagare_plugin"
                className="hidden"
              />

              <div
                onClick={handleTriggerFileSelect}
                className="border-2 border-dashed border-muted-foreground/25 hover:border-primary/50 rounded-xl p-8 text-center cursor-pointer transition-all bg-muted/30 hover:bg-muted/50 flex flex-col items-center justify-center gap-3 group"
              >
                <div className="p-4 rounded-full bg-primary/10 text-primary group-hover:scale-110 transition-transform">
                  <UploadCloud className="w-8 h-8" />
                </div>
                {selectedFile ? (
                  <div className="space-y-1">
                    <p className="font-medium text-foreground">
                      {selectedFile.name}
                    </p>
                    <p className="text-xs text-muted-foreground">
                      {(selectedFile.size / (1024 * 1024)).toFixed(2)} MB
                    </p>
                  </div>
                ) : (
                  <div className="space-y-1">
                    <p className="text-sm font-medium text-foreground">
                      Click to select a plugin package
                    </p>
                  </div>
                )}
              </div>
            </>
          )}

          {uploadStatus === "success" && (
            <Alert className="bg-emerald-500/10 border-emerald-500/20 text-emerald-600 dark:text-emerald-400">
              <CheckCircle2 className="w-4 h-4" />
              <AlertTitle>Success!</AlertTitle>
              <AlertDescription>
                Plugin installed successfully.
              </AlertDescription>
            </Alert>
          )}

          {uploadStatus === "error" && (
            <Alert variant="destructive">
              <AlertCircle className="w-4 h-4" />
              <AlertTitle>Failed</AlertTitle>
              <AlertDescription>
                Cannot install plugin. Check the server logs for more
                information.
              </AlertDescription>
            </Alert>
          )}
        </CardContent>

        <CardFooter className="flex justify-end gap-3">
          <Button
            variant="outline"
            onClick={() => {
              setSelectedFile(null);
              setFilePath("");
              setUploadStatus("idle");
            }}
            disabled={isUploading || (isDesktop ? !filePath : !selectedFile)}
          >
            Unselect
          </Button>

          <Button
            onClick={handleInstallPlugin}
            disabled={
              isDesktop
                ? !filePath.trim() || isUploading
                : !selectedFile || isUploading
            }
          >
            {isUploading && <Loader2 className="w-4 h-4 mr-2 animate-spin" />}
            {isUploading ? "Installing..." : "Install Plugin"}
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}
