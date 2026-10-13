import { instance } from "./lib/axios";
import {
  type ApiResponse,
  type AddModelProviderRequest,
  type AddModelProviderResponse,
  AddModelProviderEndpoint,
  type DeleteModelProviderRequest,
  DeleteModelProviderEndpoint,
  type GetListModelProviderResponse,
  ListModelProvidersEndpoint,
  type GetModelProviderDetailsResponse,
  GetModelProviderDetailsEndpoint,
  type GetModelProviderModelsRequest,
  type GetModelProviderModelsResponse,
  GetModelProviderModelsEndpoint,
  type ModelProvider,
} from "@nagare-app/dtos";
import { catchError } from "./lib/errors";

export type AddModelProviderPayload = AddModelProviderRequest;

export class ModelProviderService {
  static async list() {
    try {
      const response = await instance.get<
        ApiResponse<GetListModelProviderResponse>
      >(ListModelProvidersEndpoint);
      return (response.data.data.providers ?? []) as ModelProvider[];
    } catch (e) {
      catchError(e);
    }
  }

  static async getDetails(id: string) {
    try {
      const response = await instance.get<
        ApiResponse<GetModelProviderDetailsResponse>
      >(GetModelProviderDetailsEndpoint, { params: { id } });
      return response.data.data.provider!;
    } catch (e) {
      catchError(e);
    }
  }

  static async add(payload: AddModelProviderPayload) {
    try {
      const response = await instance.post<ApiResponse<AddModelProviderResponse>>(
        AddModelProviderEndpoint,
        payload,
      );
      return response.data.data.provider!;
    } catch (e) {
      catchError(e);
    }
  }

  static async delete(id: string) {
    try {
      const payload: DeleteModelProviderRequest = { id };
      await instance.post<ApiResponse<void>>(
        DeleteModelProviderEndpoint,
        payload,
      );
    } catch (e) {
      catchError(e);
    }
  }

  static async getModels(payload: GetModelProviderModelsRequest) {
    try {
      const response = await instance.post<
        ApiResponse<GetModelProviderModelsResponse>
      >(GetModelProviderModelsEndpoint, payload);
      return response.data.data.models ?? [];
    } catch (e) {
      catchError(e);
    }
  }
}
