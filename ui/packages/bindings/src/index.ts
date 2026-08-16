type AppService = {
  GetRestApiConnect(): Promise<string>;
  GetWebsocketConnect(): Promise<string>;
  IsServerRunning(): Promise<boolean>;
  ShowErrorDialog(title: string, message: string): Promise<void>;
  OpenPluginSelectDialog(): Promise<string>;
  GetToken(): Promise<string>;
  SaveToken(token: string): Promise<void>;
  GetHost(): Promise<string>;
  ClearToken(): Promise<void>;
};

declare global {
  interface Window {
    AppService: AppService;
    NagareUI: "web" | "desktop";
  }
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

export async function GetToken(): Promise<string | null> {
  if (window.NagareUI === "web") {
    return window.localStorage.getItem("token") ?? null;
  }
  let token = await window.AppService.GetToken();
  return token == "" ? null : token;
}

export async function SaveToken(token: string): Promise<void> {
  if (window.NagareUI === "web") {
    window.localStorage.setItem("token", token);
    return;
  }
  return window.AppService.SaveToken(token);
}

export async function GetHost(): Promise<string> {
  if (window.NagareUI === "web") {
    return window.location.origin;
  }
  return window.AppService.GetHost();
}

export async function ClearToken(): Promise<void> {
  if (window.NagareUI === "web") {
    window.localStorage.removeItem("token");
    return;
  }
  return window.AppService.ClearToken();
}
