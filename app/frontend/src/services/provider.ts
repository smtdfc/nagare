import type { ApiResponse, GetListProviderResponse, GetProviderByIDResponse } from "#/dto/api.ts";
import { getAxiosInstance } from "#/lib/axios.ts";
import { handleError } from "#/lib/error.ts";
import { AxiosError } from "axios";



export class ProviderService {
    static async getListProvider() {
        try {
            const instance = await getAxiosInstance();
            const resp = (await instance.get<ApiResponse<GetListProviderResponse>>("/provider/list")).data;
            return resp.data!;
        } catch (e: unknown) { handleError(e); }
    }

    static async getProviderById(id: string) {
        try {
            const instance = await getAxiosInstance();
            const resp = (await instance.get<ApiResponse<GetProviderByIDResponse>>(`/provider/${id}/details`)).data;
            return resp.data!.provider;
        } catch (e: unknown) { handleError(e); }
    }
}