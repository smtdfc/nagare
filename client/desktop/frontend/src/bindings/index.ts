import "@nagare-app/services";
import { envConfig } from "@nagare-app/services";

async function getToken() {
  return localStorage.getItem("token");
}

async function setToken(token: string | null): Promise<void> {
  localStorage.setItem("token", token!);
}

async function loadConfig(): Promise<void> {
  const configJson = localStorage.getItem("config");
  if (configJson) {
    const config = JSON.parse(configJson);
    if (config.baseRestApiUrl) {
      envConfig.baseRestApiUrl = config.baseRestApiUrl;
    }
    if (config.baseWebsocketUrl) {
      envConfig.baseWebsocketUrl = config.baseWebsocketUrl;
    }
  }
}

async function saveConfig(): Promise<void> {
  const config = {
    baseRestApiUrl: envConfig.baseRestApiUrl,
    baseWebsocketUrl: envConfig.baseWebsocketUrl,
  };
  localStorage.setItem("config", JSON.stringify(config));
}

window.bindings = {
  getToken,
  setToken,
  loadConfig,
  saveConfig,
};
