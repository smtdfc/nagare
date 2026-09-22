import { refreshInstance } from "#/lib/axios.ts";
import { refreshWebsocket } from "./lib/websocket";

export function refresh() {
  refreshInstance();
  refreshWebsocket();
}

export * from "./auth";
export * from "./env";
export * from "./lib/errors";
export * from "./chat";
export * from "./lib/websocket";
