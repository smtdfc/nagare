import { GetToken } from "@nagare-agent/service-bindings";

export async function isAuth() {
  let token = await GetToken();
  if (token) {
    return true;
  }
  return false;
}
