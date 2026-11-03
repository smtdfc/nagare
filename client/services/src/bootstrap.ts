import { envConfig } from "./env";
import { refreshInstance } from "./lib/axios";
import { normalizeGateway } from "./lib/gateway";
import { refreshWebsocket } from "./lib/websocket";

export const BOOTSTRAP_TOKEN_PARAM = "token";
export const BOOTSTRAP_GATEWAY_PARAM = "gateway";

function readParam(params: URLSearchParams, name: string): string | undefined {
  const value = params.get(name)?.trim();
  return value ? value : undefined;
}

/**
 * Applies connection settings handed over by the launcher (`nagare web` or the
 * desktop shell) through the page URL. The token and gateway are persisted to
 * the local bindings and removed from the address bar so they do not end up in
 * browser history, server logs, or referrers.
 *
 * Returns true when the page was opened with bootstrap parameters.
 */
export async function applyBootstrapFromURL(): Promise<boolean> {
  const params = new URLSearchParams(window.location.search);
  const token = readParam(params, BOOTSTRAP_TOKEN_PARAM);
  const gatewayInput = readParam(params, BOOTSTRAP_GATEWAY_PARAM);
  if (!token && !gatewayInput) {
    return false;
  }

  if (gatewayInput) {
    const gateway = normalizeGateway(gatewayInput);
    if (gateway.rest) {
      envConfig.baseRestApiUrl = gateway.rest;
      envConfig.baseWebsocketUrl = gateway.ws;
      await window.bindings.saveConfig();
    }
  }

  if (token) {
    await window.bindings.setToken(token);
  }

  params.delete(BOOTSTRAP_TOKEN_PARAM);
  params.delete(BOOTSTRAP_GATEWAY_PARAM);
  const query = params.toString();
  const hash = window.location.hash;
  const next = `${window.location.pathname}${query ? `?${query}` : ""}${hash}`;
  window.history.replaceState(window.history.state, "", next);

  refreshInstance();
  refreshWebsocket();

  return true;
}
