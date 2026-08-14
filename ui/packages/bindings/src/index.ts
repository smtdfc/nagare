type AppService = {
  GenerateToken(): Promise<string>;
  GetRestApiConnect(): Promise<string>;
  GetWebsocketConnect(): Promise<string>;
};

declare global {
  interface Window {
    AppService: AppService;
    NagareUI: "web" | "desktop";
  }
}

export function GenerateToken(): Promise<string> {
  if (window.NagareUI === "desktop") {
    return window.AppService.GenerateToken();
  }

  throw new Error("GenerateToken is not supported on web");
}

export async function GetRestApiConnect(): Promise<string> {
  if (window.NagareUI === "web") {
    return `${window.location.origin}/api/v1`;
  }
  return window.AppService.GetRestApiConnect();
}

export async function GetWebsocketConnect(): Promise<string> {
  if (window.NagareUI === "web") {
    return `${window.location.origin}/ws`;
  }
  return window.AppService.GetWebsocketConnect();
}

export {};
