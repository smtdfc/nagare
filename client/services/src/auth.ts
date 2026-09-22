import { instance } from "#/lib/axios.ts";
import {
  type ApiResponse,
  type AuthEventPayload,
  type AuthFailedEventPayload,
  type AuthSuccessEventPayload,
  CheckAuthStateEndpoint,
  type CheckAuthStateResponse,
} from "@nagare-app/dtos";
import { Event as AppEvent } from "@nagare-app/dtos";
import { catchError } from "#/lib/errors.ts";
import { websocket } from "./lib/websocket";

export class AuthService {
  static async isAuthenticated() {
    try {
      const response = await instance.get<ApiResponse<CheckAuthStateResponse>>(
        CheckAuthStateEndpoint,
      );
      const body = response.data;
      return body.data.isAuth;
    } catch (e) {
      catchError(e);
    }
  }

  static async websocketAuth() {
    try {
      const token = await window.bindings.getToken();
      if (!token) {
        throw new Error("Unauthorized");
      }

      await websocket.request<
        AuthEventPayload,
        AuthSuccessEventPayload,
        AuthFailedEventPayload
      >(
        AppEvent.AuthEvent,
        {
          token,
        },
        AppEvent.AuthSuccessEvent,
        AppEvent.AuthFailedEvent,
      );
    } catch (e) {
      catchError(e);
    }
  }
}
