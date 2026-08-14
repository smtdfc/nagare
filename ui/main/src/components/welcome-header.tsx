export default function WelcomeHeader() {
  return (
    <div className="flex-1 flex flex-col items-center justify-center text-center max-w-2xl mx-auto space-y-6 my-auto">
      <div className="space-y-2">
        <h1 className="text-3xl md:text-4xl font-bold tracking-tight">
          Nagare here 👋
        </h1>
        <p className="text-muted-foreground text-sm md:text-base">
          How can I help you today? Choose a suggestion below or start a
          conversation.
        </p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 w-full text-left">
        <button className="p-3 rounded-xl border bg-card hover:bg-accent/50 transition-colors text-sm flex flex-col gap-1 shadow-sm">
          <span className="font-medium">Open an app</span>
          <span className="text-xs text-muted-foreground">
            Run an application at your request.
          </span>
        </button>

        <button className="p-3 rounded-xl border bg-card hover:bg-accent/50 transition-colors text-sm flex flex-col gap-1 shadow-sm">
          <span className="font-medium">System error analysis</span>
          <span className="text-xs text-muted-foreground">
            Check the logs and find a way to fix the issue.
          </span>
        </button>

        <button className="p-3 rounded-xl border bg-card hover:bg-accent/50 transition-colors text-sm flex flex-col gap-1 shadow-sm">
          <span className="font-medium">Find a document</span>
          <span className="text-xs text-muted-foreground">
            Read and extract key points from the file.
          </span>
        </button>

        <button className="p-3 rounded-xl border bg-card hover:bg-accent/50 transition-colors text-sm flex flex-col gap-1 shadow-sm">
          <span className="font-medium">Create a workflow</span>
          <span className="text-xs text-muted-foreground">
            Create and edit workflows
          </span>
        </button>
      </div>
    </div>
  );
}
