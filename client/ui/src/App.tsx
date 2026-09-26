import { routeTree } from "./routeTree.gen";
import { createRouter, RouterProvider } from "@tanstack/react-router";
import { Toaster } from "#/components/ui/toast.tsx";
import { TooltipProvider } from "#/components/ui/tooltip.tsx";
import type { ChatData } from "#/types/chat.ts";

const router = createRouter({
  routeTree,
  defaultPreload: "intent",
  scrollRestoration: true,
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }

  interface HistoryState {
    initialChat?: ChatData;
  }
}

export function App() {
  return (
    <TooltipProvider>
      <Toaster />
      <RouterProvider router={router} />
    </TooltipProvider>
  );
}
