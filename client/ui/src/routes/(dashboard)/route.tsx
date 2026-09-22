import { createFileRoute, Outlet, useNavigate } from "@tanstack/react-router";
import { AppSidebar } from "#/components/app-sidebar.tsx";
import { SidebarInset, SidebarProvider } from "#/components/ui/sidebar.tsx";
import { useEffect, useState } from "react";
import { AuthService, websocket } from "@nagare-app/services";
import { showToastError } from "#/lib/toast";
import LoadingLayout from "#/components/loading-layout";
import { useAuth } from "#/hooks/use-auth.ts";

export const Route = createFileRoute("/(dashboard)")({
  component: RouteComponent,
});

function RouteComponent() {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(true);
  const isAuthenticated = useAuth((a) => a.isAuthenticated);
  const setAuth = useAuth((a) => a.setIsAuthenticated);

  const redirectToAuth = async () => {
    if (window.location.pathname.startsWith("/auth")) return;
    await navigate({
      to: "/auth",
      search: { next: window.location.pathname },
    });
  };

  useEffect(() => {
    let isMounted = true;

    const initDashboard = async () => {
      try {
        const isAuth = isAuthenticated || (await AuthService.isAuthenticated());

        if (!isAuth) {
          setAuth(false);
          await redirectToAuth();
          return;
        }

        setAuth(true);

        if (!websocket.isConnected()) {
          await websocket.connect();
        }
        await AuthService.websocketAuth();
      } catch (e) {
        setAuth(false);
        showToastError(e);
        await redirectToAuth();
      } finally {
        if (isMounted) setIsLoading(false);
      }
    };

    initDashboard();

    return () => {
      isMounted = false;
    };
  }, [isAuthenticated]);

  if (isLoading) {
    return <LoadingLayout text="Loading..." />;
  }

  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <Outlet />
      </SidebarInset>
    </SidebarProvider>
  );
}
