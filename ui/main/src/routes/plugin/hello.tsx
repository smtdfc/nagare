import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/plugin/hello")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/plugin/hello"!</div>;
}
