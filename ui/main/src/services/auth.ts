import { getAxiosInstance } from "#/lib/axios";
import { handleError } from "#/lib/error";
import { builtInWs, wsRequest } from "#/lib/websocket";
import {
  WsEvent,
  type ApiResponse,
  type GetProfileResponse,
  type WsAuthSuccess,
} from "@nagare-agent/dto";
import { ClearToken, GetToken } from "@nagare-agent/service-bindings";

export class AuthService {
  static async profile() {
    try {
      const instance = await getAxiosInstance(true);
      const resp = (
        await instance.get<ApiResponse<GetProfileResponse>>("/auth/me")
      ).data;
      return resp.data!.profile;
    } catch (e: unknown) {
      handleError(e);
    }
  }

  static async logout() {
    try {
      // const instance = await getAxiosInstance(true);
      // await instance.post("/auth/logout");
      await ClearToken();
    } catch (e: unknown) {
      handleError(e);
    }
  }

  static async websocketAuth() {
    if (!builtInWs.chat) throw new Error("Websocket is not ready");
    return await wsRequest<WsAuthSuccess>(
      builtInWs.chat,
      WsEvent.WS_AUTH_REQUEST,
      WsEvent.WS_AUTH_SUCCESS,
      WsEvent.WS_AUTH_FAILED,
      {
        token: await GetToken(),
      },
    );
  }
}
