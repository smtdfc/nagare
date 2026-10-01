import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/(dashboard)/settings/llm-providers")({
  component: RouteComponent,
});

function RouteComponent() {
  return <Outlet />;
}
