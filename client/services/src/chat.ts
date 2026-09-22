import { instance } from "#/lib/axios.ts";
import {
  type ApiResponse,
  CreateChatSessionEndpoint,
  type CreateChatSessionResponse,
  GetChatSessionEndpoint,
  type GetChatSessionResponse,
} from "@nagare-app/dtos";
import { catchError } from "#/lib/errors.ts";

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
}
