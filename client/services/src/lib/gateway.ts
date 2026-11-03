export type GatewayURLs = {
  rest: string;
  ws: string;
};

/**
 * Normalizes the gateway URL typed by the user before it is stored in the
 * environment config. A trailing slash used to produce "//ws", which the
 * gateway does not route, so the REST check passed while the websocket
 * handshake failed.
 */
export function normalizeGateway(input: string): GatewayURLs {
  let rest = input.trim().replace(/\/+$/, "");
  rest = rest.replace(/^ws:/i, "http:").replace(/^wss:/i, "https:");

  if (!rest) {
    return { rest: "", ws: "" };
  }

  if (!/^https?:\/\//i.test(rest)) {
    rest = `http://${rest}`;
  }

  const ws = rest.replace(/^http:/i, "ws:").replace(/^https:/i, "wss:");

  return { rest, ws: `${ws}/ws` };
}
