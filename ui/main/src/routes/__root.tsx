import { Outlet, createRootRoute } from "@tanstack/react-router";
import { AppSidebar } from "@/components/app-sidebar";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { DynamicBreadcrumbs } from "@/components/dynamic-breadcrumbs.tsx";
import { useEffect, useState } from "react";
import { isAuth } from "#/lib/auth";
import { AuthForm } from "#/components/auth-form";
import { AuthService } from "#/services/auth";
import { toast } from "#/components/ui/toast";
import { SaveToken } from "@nagare-agent/service-bindings";
import { useAuth } from "#/hooks/use-auth";
import { Spinner } from "#/components/ui/spinner";
import { Separator } from "@base-ui/react";

export const Route = createRootRoute({
  component: RootComponent,
});

function RootLayout() {
  return (
    <>
      <SidebarProvider className="h-screen w-screen overflow-hidden">
        <AppSidebar />
        <SidebarInset className="relative h-screen w-full overflow-hidden flex flex-col">
          <header className="flex h-16 shrink-0 items-center gap-2">
            <div className="flex items-center gap-2 px-4">
              <SidebarTrigger className="-ml-1" />
              <Separator
                orientation="vertical"
                className="mr-2 data-[orientation=vertical]:h-4"
              />
              <DynamicBreadcrumbs />
            </div>
          </header>
          <div className="flex-1 relative flex flex-col overflow-y-auto">
            <Outlet />
          </div>
        </SidebarInset>
      </SidebarProvider>
    </>
  );
}

function AuthLayout() {
  const setAuthState = useAuth((s) => s.setAuthState);
  const handleSubmit = async (e: { host: string; token: string }) => {
    try {
      await SaveToken(e.token);
      const profile = await AuthService.profile();
      if (!profile) throw new Error("Unauthorized");
      setAuthState(profile);
    } catch (e) {
      toast.add({
        title: "Error",
        description: e instanceof Error ? e.message : String(e),
        type: "error",
      });
    }
  };
  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-background p-6 md:p-10">
      <div className="w-full max-w-sm">
        <AuthForm onSubmit={handleSubmit} />
      </div>
    </div>
  );
}

export function RootComponent() {
  const auth = useAuth((s) => s.auth);
  const setAuthState = useAuth((s) => s.setAuthState);
  const [isLoading, setIsLoading] = useState(true); // Thêm trạng thái kiểm tra

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const profile = await AuthService.profile();
        if (!profile) throw new Error("Unauthorized");
        setAuthState(profile);
      } catch (e) {
        toast.add({
          title: "Error",
          description: e instanceof Error ? e.message : String(e),
          type: "error",
        });
      } finally {
        setIsLoading(false);
      }
    };
    checkAuth();
  }, [setAuthState]);

  if (isLoading) {
    return (
      <div className="flex h-screen w-full items-center justify-center">
        <Spinner className="size-10" />
      </div>
    );
  }

  return auth ? <RootLayout /> : <AuthLayout />;
}
