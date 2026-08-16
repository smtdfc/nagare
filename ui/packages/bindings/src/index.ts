type AppService = {
  GenerateToken(): Promise<string>;
  GetRestApiConnect(): Promise<string>;
  GetWebsocketConnect(): Promise<string>;
  IsServerRunning(): Promise<boolean>;
  ShowErrorDialog(title: string, message: string): Promise<void>;
  OpenPluginSelectDialog(): Promise<string>;
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

export async function IsServerRunning(): Promise<boolean> {
  if (window.NagareUI === "web") {
    return true;
  }
  return window.AppService.IsServerRunning();
}

export async function ShowErrorDialog(
  title: string,
  message: string,
): Promise<void> {
  if (window.NagareUI === "web") {
    alert(`${title}: ${message}`);
    return;
  }
  return window.AppService.ShowErrorDialog(title, message);
}

export async function OpenPluginSelectDialog(): Promise<string> {
  if (window.NagareUI === "web") {
    return "";
  }
  return window.AppService.OpenPluginSelectDialog();
}
