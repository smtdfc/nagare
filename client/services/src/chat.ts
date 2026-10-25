import { instance } from "./lib/axios.js";
import {
  type ApiResponse,
  CreateChatSessionEndpoint,
  type CreateChatSessionResponse,
  ArchiveChatSessionEndpoint,
  type ArchiveChatSessionRequest,
  type ArchiveChatSessionResponse,
  DeleteChatSessionEndpoint,
  DuplicateChatSessionEndpoint,
  type DuplicateChatSessionResponse,
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
  type SubmitToolApprovalEventPayload,
  type SubmitToolApprovalSuccessEventPayload,
  type SubmitToolApprovalFailEventPayload,
  SendChatMessageEndpoint,
  type SendChatMessageResponse,
  type Session,
  UpdateChatSessionModelSettingsEndpoint,
  type UpdateChatSessionModelSettingsRequest,
  type UpdateChatSessionModelSettingsResponse,
} from "@nagare-app/dtos";
import { catchError } from "./lib/errors.js";
import type { Message } from "@nagare-app/messages";
import { websocket } from "./lib/websocket";
import { Event as AppEvent } from "@nagare-app/dtos";
import { AuthService } from "./auth";

type ChatHandler = (message: Message) => Promise<void>;

const sessionHandlers = new Map<string, ChatHandler>();
let reconnectHandlingInstalledOn: typeof websocket | null = null;

function handleReceivedChatMessage(
  payload: Payload<ReceivedChatMessageEventPayload>,
) {
  const handler = sessionHandlers.get(payload.data.sessionID);
  if (!handler) return;

  const message = JSON.parse(payload.data.message) as Message;
  void handler(message).catch((error) => {
    console.error("Failed to handle chat message:", error);
  });
}

async function registerChatListener(sessionID: string) {
  await websocket.request<
    RegisterChatMessageListenerEventPayload,
    RegisterChatListenerSuccessEventEventPayload,
    RegisterChatListenerFailEventPayload
  >(
    AppEvent.RegisterChatListenerEvent,
    {
      id: "",
      sessionID,
    },
    AppEvent.RegisterChatListenerSuccessEvent,
    AppEvent.RegisterChatListenerFailEvent,
  );
}

function installReconnectHandling() {
  if (reconnectHandlingInstalledOn === websocket) {
    return;
  }
  reconnectHandlingInstalledOn = websocket;

  websocket.on(AppEvent.ReceivedChatMessageEvent, handleReceivedChatMessage);

  websocket.onReconnect(async () => {
    try {
      await AuthService.websocketAuth();
    } catch (error) {
      console.error("Failed to re-authenticate WebSocket:", error);
      return;
    }

    for (const sessionID of [...sessionHandlers.keys()]) {
      try {
        await registerChatListener(sessionID);
      } catch (error) {
        console.error(
          `Failed to restore chat listener for session ${sessionID}:`,
          error,
        );
      }
    }
  });
}

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

  static async updateChatSessionModelSettings(
    id: string,
    payload: UpdateChatSessionModelSettingsRequest,
  ) {
    try {
      const response = await instance.post<
        ApiResponse<UpdateChatSessionModelSettingsResponse>
      >(UpdateChatSessionModelSettingsEndpoint.replace(":id", id), payload);
      return response.data.data.session!;
    } catch (error) {
      catchError(error);
    }
  }

  /** @deprecated use updateChatSessionModelSettings instead */
  static async updateChatSessionLLMSettings(
    id: string,
    payload: UpdateChatSessionModelSettingsRequest,
  ) {
    return ChatService.updateChatSessionModelSettings(id, payload);
  }

  static async deleteChatSession(id: string) {
    try {
      await instance.delete(DeleteChatSessionEndpoint.replace(":id", id));
    } catch (error) {
      catchError(error);
    }
  }

  static async archiveChatSession(
    id: string,
    payload: ArchiveChatSessionRequest,
  ) {
    try {
      const response = await instance.post<
        ApiResponse<ArchiveChatSessionResponse>
      >(ArchiveChatSessionEndpoint.replace(":id", id), payload);
      return response.data.data.session!;
    } catch (error) {
      catchError(error);
    }
  }

  static async duplicateChatSession(id: string) {
    try {
      const response = await instance.post<
        ApiResponse<DuplicateChatSessionResponse>
      >(DuplicateChatSessionEndpoint.replace(":id", id));
      return response.data.data.session!;
    } catch (error) {
      catchError(error);
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
    installReconnectHandling();

    try {
      await registerChatListener(id);
      sessionHandlers.set(id, handler);

      return () => {
        sessionHandlers.delete(id);
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
      const response = await instance.post<
        ApiResponse<SendChatMessageResponse>
      >(SendChatMessageEndpoint, {
        sessionID: id,
        text,
      });
      return response.data.data;
    } catch (error) {
      catchError(error);
    }
  }

  static async submitToolApproval(
    payload: SubmitToolApprovalEventPayload,
  ) {
    return websocket.request<
      SubmitToolApprovalEventPayload,
      SubmitToolApprovalSuccessEventPayload,
      SubmitToolApprovalFailEventPayload
    >(
      AppEvent.SubmitToolApprovalEvent,
      payload,
      AppEvent.SubmitToolApprovalSuccessEvent,
      AppEvent.SubmitToolApprovalFailEvent,
    );
  }
}
