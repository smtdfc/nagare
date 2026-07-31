import type { ApiResponse, GetListProviderResponse } from "#/dto/api.ts";
import { getAxiosInstance } from "#/lib/axios.ts";

export class ProviderService {
    static async getListProvider() {
        const instance = await getAxiosInstance();
        const resp = (await instance.get<ApiResponse<GetListProviderResponse>>("/provider/list")).data;
        return resp.data!;
    }
}