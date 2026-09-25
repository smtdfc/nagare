import { instance } from "#/lib/axios.ts";
import {
  type ApiResponse,
  CreateChatSessionEndpoint,
  type CreateChatSessionResponse,
  GetChatHistoryEndpoint,
  type GetChatHistoryResponse,
  GetChatSessionEndpoint,
  type GetChatSessionResponse,
  type Payload,
  type ReceivedChatMessageEventPayload,
  type RegisterChatListenerFailEventPayload,
  type RegisterChatListenerSuccessEventEventPayload,
  type RegisterChatMessageListenerEventPayload,
  SendChatMessageEndpoint,
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

  static async getChatHistory(id: string) {
    try {
      const response = await instance.get<ApiResponse<GetChatHistoryResponse>>(
        GetChatHistoryEndpoint.replace(":id", id),
        {},
      );
      const body = response.data;
      return (body.data.messages as Message[])!;
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
