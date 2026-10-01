import { instance } from "./lib/axios";
import {
  type ApiResponse,
  type GeneralSettings,
  type GetGeneralSettingsResponse,
  GetGeneralSettings,
  type SetGeneralSettingsRequest,
  SetGeneralSettings,
} from "@nagare-app/dtos";
import { catchError } from "./lib/errors";

export class SettingsService {
  static async getGeneralSettings() {
    try {
      const response =
        await instance.get<ApiResponse<GetGeneralSettingsResponse>>(
          GetGeneralSettings,
        );
      return response.data.data.generalSettings!;
    } catch (e) {
      catchError(e);
    }
  }

  static async setGeneralSettings(settings: GeneralSettings) {
    try {
      const payload: SetGeneralSettingsRequest = { generalSettings: settings };
      await instance.post<ApiResponse<void>>(SetGeneralSettings, payload);
    } catch (e) {
      catchError(e);
    }
  }
}
