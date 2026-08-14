import type { ApiResponse, GetAllPluginsResponse } from "@nagare-agent/dto";
import { getAxiosInstance } from "#/lib/axios";
import { handleError } from "#/lib/error";

export class PluginService {
  static async getAllPlugin() {
    try {
      const instance = await getAxiosInstance();
      const resp = (
        await instance.get<ApiResponse<GetAllPluginsResponse>>("/plugin/list")
      ).data;
      return resp.data!.plugins;
    } catch (e: unknown) {
      handleError(e);
    }
  }
}
