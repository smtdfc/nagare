import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/(dashboard)/settings/model-providers")({
  component: RouteComponent,
});

function RouteComponent() {
  return <Outlet />;
}
