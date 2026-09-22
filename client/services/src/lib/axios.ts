import axios from "axios";
import { envConfig } from "#/env.ts";

export function createAxiosWithCurrentEnv() {
  const instance = axios.create({
    baseURL: envConfig.baseRestApiUrl,
    timeout: 1000 * 30,
  });

  instance.interceptors.request.use(async (config) => {
    const token = await window.bindings.getToken();
    if (token) {
      config.headers.set("Authorization", `Bearer ${token}`);
    }
    return config;
  });

  return instance;
}

let instance = createAxiosWithCurrentEnv();

export function refreshInstance() {
  instance = createAxiosWithCurrentEnv();
}

export { instance };
