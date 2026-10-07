import { instance } from "./lib/axios";
import {
  type ApiResponse,
  type AddLLMProviderRequest,
  type AddLLMProviderResponse,
  AddLLMProviderEndpoint,
  type DeleteLLMProviderRequest,
  DeleteLLMProviderEndpoint,
  type GetListLLMProviderResponse,
  ListLLMProvidersEndpoint,
  type GetLLMProviderDetailsResponse,
  GetLLMProviderDetailsEndpoint,
  type GetLLMProviderModelsRequest,
  type GetLLMProviderModelsResponse,
  GetLLMProviderModelsEndpoint,
  type LLMProvider,
} from "@nagare-app/dtos";
import { catchError } from "./lib/errors";

export type AddLLMProviderPayload = AddLLMProviderRequest;

export class LLMProviderService {
  static async list() {
    try {
      const response = await instance.get<
        ApiResponse<GetListLLMProviderResponse>
      >(ListLLMProvidersEndpoint);
      return (response.data.data.providers ?? []) as LLMProvider[];
    } catch (e) {
      catchError(e);
    }
  }

  static async getDetails(id: string) {
    try {
      const response = await instance.get<
        ApiResponse<GetLLMProviderDetailsResponse>
      >(GetLLMProviderDetailsEndpoint, { params: { id } });
      return response.data.data.provider!;
    } catch (e) {
      catchError(e);
    }
  }

  static async add(payload: AddLLMProviderPayload) {
    try {
      const response = await instance.post<ApiResponse<AddLLMProviderResponse>>(
        AddLLMProviderEndpoint,
        payload,
      );
      return response.data.data.provider!;
    } catch (e) {
      catchError(e);
    }
  }

  static async delete(id: string) {
    try {
      const payload: DeleteLLMProviderRequest = { id };
      await instance.post<ApiResponse<void>>(
        DeleteLLMProviderEndpoint,
        payload,
      );
    } catch (e) {
      catchError(e);
    }
  }

  static async getModels(payload: GetLLMProviderModelsRequest) {
    try {
      const response = await instance.post<
        ApiResponse<GetLLMProviderModelsResponse>
      >(GetLLMProviderModelsEndpoint, payload);
      return response.data.data.models ?? [];
    } catch (e) {
      catchError(e);
    }
  }
}
