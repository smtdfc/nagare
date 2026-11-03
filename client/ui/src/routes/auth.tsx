import { LoginForm, type LoginInput } from "#/components/login-form";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { AuthService, envConfig, refresh } from "@nagare-app/services";
import { showToastError } from "#/lib/toast.ts";
import { useAuth } from "#/hooks/use-auth.ts";
import { normalizeGateway } from "#/lib/gateway.ts";

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
    const gateway = normalizeGateway(input.gateway ?? "");
    if (!gateway.rest || !input.token) {
      setAuth(false);
      showToastError(new Error("Gateway URL and token are required."));
      return;
    }

    envConfig.baseRestApiUrl = gateway.rest;
    envConfig.baseWebsocketUrl = gateway.ws;
    await window.bindings.setToken(input.token);
    refresh();

    let isAuth = false;
    try {
      isAuth = await AuthService.isAuthenticated();
    } catch (e) {
      setAuth(false);
      showToastError(e);
      return;
    }

    if (!isAuth) {
      setAuth(false);
      showToastError(
        new Error("Authentication failed; please check the Gateway or Token!"),
      );
      return;
    }

    // The REST session is valid. The dashboard owns the websocket connection,
    // so a socket failure cannot turn a successful login into a failed one.
    setAuth(true);
    await window.bindings.saveConfig();
    await redirect();
  };

  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-background p-6 md:p-10">
      <div className="w-full max-w-sm">
        <LoginForm onSubmit={handleSubmit} />
      </div>
    </div>
  );
}
