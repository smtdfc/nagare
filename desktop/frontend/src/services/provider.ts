import type {
  ApiResponse,
  CreateProviderRequest,
  CreateProviderResponse,
  DeleteProviderResponse,
  FetchModelRequest,
  FetchModelResponse,
  GetListProviderResponse,
  GetProviderByIDResponse,
  Provider,
  UpdateProviderRequest,
  UpdateProviderResponse,
} from '#/dto/api.ts'
import { getAxiosInstance } from '#/lib/axios.ts'
import { handleError } from '#/lib/error.ts'

export class ProviderService {
  static async getListProvider() {
    try {
      const instance = await getAxiosInstance()
      const resp = (
        await instance.get<ApiResponse<GetListProviderResponse>>(
          '/provider/list',
        )
      ).data
      return resp.data!
    } catch (e: unknown) {
      handleError(e)
    }
  }

  static async getProviderById(id: string) {
    try {
      const instance = await getAxiosInstance()
      const resp = (
        await instance.get<ApiResponse<GetProviderByIDResponse>>(
          `/provider/${id}/details`,
        )
      ).data
      return resp.data!.provider
    } catch (e: unknown) {
      handleError(e)
    }
  }

  static async updateProvider(data: UpdateProviderRequest) {
    try {
      const instance = await getAxiosInstance()
      const resp = (
        await instance.post<ApiResponse<UpdateProviderResponse>>(
          `/provider/update`,
          data,
        )
      ).data
    } catch (e: unknown) {
      handleError(e)
    }
  }

  static async deleteProvider(id: string) {
    try {
      const instance = await getAxiosInstance()
      const resp = (
        await instance.post<ApiResponse<DeleteProviderResponse>>(
          `/provider/delete`,
          { id },
        )
      ).data
      return resp.data
    } catch (e: unknown) {
      handleError(e)
    }
  }

  static async createProvider(data: CreateProviderRequest) {
    try {
      const instance = await getAxiosInstance()
      const resp = (
        await instance.post<ApiResponse<CreateProviderResponse>>(
          `/provider/create`,
          data,
        )
      ).data
      return resp.data
    } catch (e: unknown) {
      handleError(e)
    }
  }

  static async fetchModel(data: FetchModelRequest) {
    try {
      const instance = await getAxiosInstance()
      const resp = (
        await instance.post<ApiResponse<FetchModelResponse>>(
          `/provider/fetch-model`,
          data,
        )
      ).data
      return resp.data?.models ?? []
    } catch (e: unknown) {
      handleError(e)
    }
  }
}
