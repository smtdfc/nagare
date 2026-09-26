import { instance } from "#/lib/axios.ts";
import {
  type ApiResponse,
  CreateChatSessionEndpoint,
  type CreateChatSessionResponse,
  GetChatHistoryEndpoint,
  type GetChatHistoryResponse,
  GetChatSessionEndpoint,
  type GetChatSessionResponse,
  ListChatSessionsEndpoint,
  type ListChatSessionsResponse,
  type Payload,
  type ReceivedChatMessageEventPayload,
  type RegisterChatListenerFailEventPayload,
  type RegisterChatListenerSuccessEventEventPayload,
  type RegisterChatMessageListenerEventPayload,
  SendChatMessageEndpoint,
  type Session,
} from "@nagare-app/dtos";
import { catchError } from "#/lib/errors.ts";
import type { Message } from "@nagare-app/messages";
import { websocket } from "./lib/websocket";
import { Event as AppEvent } from "@nagare-app/dtos";

type ChatHandler = (message: Message) => Promise<void>;

export class ChatService {
  static async createChatSession(title: string) {
    try {
      const response = await instance.post<
        ApiResponse<CreateChatSessionResponse>
      >(CreateChatSessionEndpoint, { title });
      const body = response.data;
      return body.data.session!;
    } catch (e) {
      catchError(e);
    }
  }

  static async getChatSession(id: string) {
    try {
      const response = await instance.get<ApiResponse<GetChatSessionResponse>>(
        GetChatSessionEndpoint.replace(":id", id),
      );
      const body = response.data;
      return body.data.session!;
    } catch (e) {
      catchError(e);
    }
  }

  static async listSessions(limit = 50, offset = 0) {
    try {
      const response = await instance.get<
        ApiResponse<ListChatSessionsResponse>
      >(ListChatSessionsEndpoint, { params: { limit, offset } });
      const body = response.data;
      return body.data.sessions! as Session[];
    } catch (e) {
      catchError(e);
    }
  }

  static async getChatHistory(id: string, limit = 50, beforeID = "") {
    try {
      const response = await instance.get<ApiResponse<GetChatHistoryResponse>>(
        GetChatHistoryEndpoint.replace(":id", id),
        { params: { limit, beforeID } },
      );
      const body = response.data;
      return {
        messages: (body.data.messages as Message[])!,
        nextCursor: body.data.nextCursor,
      };
    } catch (error) {
      catchError(error);
    }
  }

  static async listenChatSession(id: string, handler: ChatHandler) {
    const handlePayload = (
      payload: Payload<ReceivedChatMessageEventPayload>,
    ) => {
      const message = JSON.parse(payload.data.message) as Message;
      handler(message);
    };
    try {
      await websocket.request<
        RegisterChatMessageListenerEventPayload,
        RegisterChatListenerSuccessEventEventPayload,
        RegisterChatListenerFailEventPayload
      >(
        AppEvent.RegisterChatListenerEvent,
        {
          id: "",
          sessionID: id,
        },
        AppEvent.RegisterChatListenerSuccessEvent,
        AppEvent.RegisterChatListenerFailEvent,
      );

      websocket.on(AppEvent.ReceivedChatMessageEvent, handlePayload);

      return () => {
        websocket.off(AppEvent.ReceivedChatMessageEvent, handlePayload);
        websocket.send(AppEvent.UnregisterChatListenerEvent, {
          id: "",
          sessionID: id,
        });
      };
    } catch (error) {
      catchError(error);
    }
  }

  static async sendMessage(id: string, text: string) {
    try {
      await instance.post<ApiResponse<any>>(SendChatMessageEndpoint, {
        sessionID: id,
        text,
      });
    } catch (error) {
      catchError(error);
    }
  }
}
