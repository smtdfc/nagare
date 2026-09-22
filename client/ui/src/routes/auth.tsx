import { LoginForm, type LoginInput } from "#/components/login-form";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import {
  AuthService,
  envConfig,
  refresh,
  websocket,
} from "@nagare-app/services";
import { showToastError } from "#/lib/toast.ts";
import { useAuth } from "#/hooks/use-auth.ts";

export const Route = createFileRoute("/auth")({
  component: RouteComponent,
});

function RouteComponent() {
  const navigate = useNavigate();
  const setAuth = useAuth((a) => a.setIsAuthenticated);

  const redirect = async () => {
    const params = new URLSearchParams(window.location.search);
    const nextUrl = params.get("next");

    await navigate({
      to: nextUrl || "/",
    });
  };

  const handleSubmit = async (input: LoginInput) => {
    envConfig.baseRestApiUrl = input.gateway!;
    envConfig.baseWebsocketUrl = `${input.gateway}/ws`;
    window.bindings.setToken(input.token!);
    refresh();
    try {
      const isAuth = await AuthService.isAuthenticated();
      await websocket.connect();
      await AuthService.websocketAuth();
      if (isAuth) {
        setAuth(true);
        await redirect();
      } else {
        setAuth(false);
        showToastError(
          "Authentication failed; please check the Gateway or Token!",
        );
      }
    } catch (e) {
      setAuth(false);
      showToastError(e);
    }
  };

  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-background p-6 md:p-10">
      <div className="w-full max-w-sm">
        <LoginForm onSubmit={handleSubmit} />
      </div>
    </div>
  );
}
