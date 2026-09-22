export enum Environment {
  Desktop,
  Web,
}

export interface EnvBindings {
  getToken(): Promise<string | null>;
  setToken: (token: string | null) => void;
}

export const envConfig = {
  current: Environment.Web,
  baseRestApiUrl: "http://localhost:9832",
  baseWebsocketUrl: "http://localhost:9832/ws",
};

declare global {
  interface Window {
    bindings: EnvBindings;
  }
}
