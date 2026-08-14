import { createFileRoute } from '@tanstack/react-router'
import { useState, useRef } from 'react'
import {
  UploadCloud,
  FileBox,
  CheckCircle2,
  AlertCircle,
  Loader2,
} from 'lucide-react'

import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'

export const Route = createFileRoute('/plugin/add')({
  component: RouteComponent,
})

function RouteComponent() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [isUploading, setIsUploading] = useState(false)
  const [uploadStatus, setUploadStatus] = useState<
    'idle' | 'success' | 'error'
  >('idle')
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0]
      setSelectedFile(file)
      setUploadStatus('idle')
    }
  }

  const handleTriggerFileSelect = () => {
    fileInputRef.current?.click()
  }

  const handleInstallPlugin = async () => {
    if (!selectedFile) return

    setIsUploading(true)
    setUploadStatus('idle')

    try {
      await new Promise((resolve) => setTimeout(resolve, 2000)) // Fake delay
      setUploadStatus('success')
    } catch (err) {
      setUploadStatus('error')
    } finally {
      setIsUploading(false)
    }
  }

  return (
    <div className="flex justify-center items-center min-h-[80vh] p-4">
      <Card className="w-full max-w-xl  ">
        <CardHeader>
          <CardTitle className="text-2xl font-bold flex items-center gap-2">
            <FileBox className="w-6 h-6 text-primary" />
            Install Plugin
          </CardTitle>
          <CardDescription>
            Select a plugin package (.nagare_plugin file) to integrate into your
            Nagare system.
          </CardDescription>
        </CardHeader>

        <CardContent className="space-y-6">
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

          {uploadStatus === 'success' && (
            <Alert className="bg-emerald-500/10 border-emerald-500/20 text-emerald-600 dark:text-emerald-400">
              <CheckCircle2 className="w-4 h-4" />
              <AlertTitle>Success!</AlertTitle>
              <AlertDescription>
                Plugin installed successfully.
              </AlertDescription>
            </Alert>
          )}

          {uploadStatus === 'error' && (
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
              setSelectedFile(null)
              setUploadStatus('idle')
            }}
            disabled={isUploading || !selectedFile}
          >
            Hủy chọn
          </Button>

          <Button
            onClick={handleInstallPlugin}
            disabled={!selectedFile || isUploading}
          >
            {isUploading && <Loader2 className="w-4 h-4 mr-2 animate-spin" />}
            {isUploading ? 'Installing...' : 'Install Plugin'}
          </Button>
        </CardFooter>
      </Card>
    </div>
  )
}
