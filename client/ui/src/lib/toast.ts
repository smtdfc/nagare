import { ApiError } from "@nagare-app/services";
import { toast } from "#/components/ui/toast.tsx";

export function showToastError(err: unknown) {
  if (err instanceof ApiError) {
    toast.add({
      title: "Error",
      type: "error",
      priority: "high",
      description: err.apiError.message,
    });
  } else if (err instanceof Error) {
    toast.add({
      title: "Error",
      type: "error",
      priority: "high",
      description: err.message,
    });
  } else {
    toast.add({
      title: "Error",
      type: "error",
      priority: "high",
      description: "An unknown error occurred.",
    });
  }
}
