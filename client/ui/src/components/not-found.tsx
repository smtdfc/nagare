import { Link } from "@tanstack/react-router";
import { ArrowLeft } from "lucide-react";
import { Button } from "#/components/ui/button.tsx";

export function NotFoundComponent() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[80vh] w-full p-6 text-center">
      <div className="flex flex-col items-center max-w-md w-full gap-6">
        <div className="flex flex-col gap-2">
          <span className="text-7xl  uppercase tracking-widest text-muted-foreground">
            404
          </span>
          <h1 className="text-2xl font-semibold tracking-tight">
            Page not found
          </h1>
          <p className="text-sm text-muted-foreground">
            Sorry, the page you are looking for doesn't exist or has been moved.
          </p>
        </div>

        <div className="flex items-center gap-3 w-full justify-center pt-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => window.history.back()}
            className="gap-2"
          >
            <ArrowLeft className="w-4 h-4" />
            Go back
          </Button>

          <Button size="sm" className="gap-2">
            <Link to="/">Back to Dashboard</Link>
          </Button>
        </div>
      </div>
    </div>
  );
}
