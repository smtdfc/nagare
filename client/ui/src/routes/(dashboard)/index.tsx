import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect } from "react";
import LoadingLayout from "#/components/loading-layout.tsx";

export const Route = createFileRoute("/(dashboard)/")({
  component: RouteComponent,
});

function RouteComponent() {
  const navigate = useNavigate();

  useEffect(() => {
    navigate({
      to: "/chat/new",
    });
  }, []);
  return (
    <>
      <LoadingLayout />
    </>
  );
}
