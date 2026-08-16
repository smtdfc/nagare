import { GetRestApiConnect, GetToken } from "@nagare-agent/service-bindings";
import type { Axios } from "axios";
import axios from "axios";

let token: string | null = null;
let instance: Axios;

export async function getAxiosInstance(
  forceGetLatestToken = false,
): Promise<Axios> {
  if (instance && !forceGetLatestToken) return instance;

  if (token === null || forceGetLatestToken) {
    token = await GetToken();
    console.log("token", token);
  }

  instance = axios.create({
    baseURL: await GetRestApiConnect(),
    timeout: 5000,
    headers: { "X-Nagare-Secure": token },
  });

  return instance;
}
