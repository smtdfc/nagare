import type { ApiResponse, GetGeneralSettingsResponse } from '#/dto/api'
import { getAxiosInstance } from '#/lib/axios'
import { handleError } from '#/lib/error'

export class SettingsService {
  static async getGeneralSettings() {
    try {
      const instance = await getAxiosInstance()
      const resp = (
        await instance.get<ApiResponse<GetGeneralSettingsResponse>>(
          `/settings/general`,
        )
      ).data
      return resp.data!.settings
    } catch (e: unknown) {
      handleError(e)
    }
  }

  static async saveGeneralSettings(settings: GetGeneralSettingsResponse) {
    try {
      const instance = await getAxiosInstance()
      await instance.post(`/settings/general/save`, settings)
    } catch (e: unknown) {
      handleError(e)
    }
  }
}
