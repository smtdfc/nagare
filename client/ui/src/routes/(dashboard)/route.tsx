import { createFileRoute, Outlet, useNavigate } from "@tanstack/react-router";
import { AppSidebar } from "#/components/app-sidebar.tsx";
import { SidebarInset, SidebarProvider } from "#/components/ui/sidebar.tsx";
import { useEffect, useState } from "react";
import { ApiError, AuthService, websocket } from "@nagare-app/services";
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

    const connectWebsocket = async () => {
      try {
        if (!websocket.isConnected()) {
          await websocket.connect();
        }
        await AuthService.websocketAuth();
      } catch (e) {
        // The REST session is already valid, so a socket problem must not log
        // the user out. Reconnect handlers re-authenticate once it is back.
        showToastError(e);
      }
    };

    const initDashboard = async () => {
      try {
        let isAuth = isAuthenticated;

        if (!isAuth) {
          try {
            isAuth = await AuthService.isAuthenticated();
          } catch (e) {
            // A missing or expired token is the normal first-load path, so
            // redirect silently instead of showing an unauthorized toast.
            if (!(e instanceof ApiError && e.apiError.statusCode === 401)) {
              showToastError(e);
            }
            isAuth = false;
          }
        }

        if (!isAuth) {
          setAuth(false);
          await redirectToAuth();
          return;
        }

        setAuth(true);
      } finally {
        if (isMounted) setIsLoading(false);
      }

      await connectWebsocket();
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
