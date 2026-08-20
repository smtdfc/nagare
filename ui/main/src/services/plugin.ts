import type {
  ApiResponse,
  GetAllPluginsResponse,
  InstallLocalPluginRequest,
  RemovePluginRequest,
} from "@nagare-agent/dto";
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

  static async installLocalPlugin(data: InstallLocalPluginRequest) {
    try {
      const instance = await getAxiosInstance();
      await instance.post("/plugin/install-local", data);
    } catch (e: unknown) {
      handleError(e);
    }
  }

  static async removePlugin(data: RemovePluginRequest) {
    try {
      const instance = await getAxiosInstance();
      await instance.post("/plugin/remove", data);
    } catch (e: unknown) {
      handleError(e);
    }
  }
}
