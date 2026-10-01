import { createFileRoute, Outlet } from "@tanstack/react-router";

export const Route = createFileRoute("/(dashboard)/settings/general")({
  component: RouteComponent,
});

function RouteComponent() {
  return <Outlet />;
}
