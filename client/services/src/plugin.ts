import { instance } from "./lib/axios";
import {
  type ApiResponse,
  type GetListPluginResponse,
  GetListPluginEndpoint,
  type InstallLocalPluginRequest,
  InstallLocalPluginEndpoint,
  type UninstallPluginRequest,
  UninstallPluginEndpoint,
  type ActivatePluginRequest,
  ActivatePluginEndpoint,
  type DeactivatePluginRequest,
  DeactivatePluginEndpoint,
  type GetPluginStatusRequest,
  type GetPluginStatusResponse,
  GetPluginStatusEndpoint,
  type UploadPluginResponse,
  UploadPluginEndpoint,
  type Plugin,
  type InstallPluginFromAttachmentResponse,
  InstallPluginFromAttachmentEndpoint,
  type InstallLocalPluginResponse,
} from "@nagare-app/dtos";
import { catchError } from "./lib/errors";

export class PluginService {
  static async list() {
    try {
      const response = await instance.get<ApiResponse<GetListPluginResponse>>(
        GetListPluginEndpoint,
      );
      return (response.data.data.plugins ?? []) as Plugin[];
    } catch (e) {
      catchError(e);
    }
  }

  static async installLocal(path: string): Promise<Plugin | null> {
    try {
      const payload: InstallLocalPluginRequest = { path };
      const data = await instance.post<ApiResponse<InstallLocalPluginResponse>>(
        InstallLocalPluginEndpoint,
        payload,
      );
      return data.data.data.plugin;
    } catch (e) {
      catchError(e);
    }
  }

  static async uninstall(id: string) {
    try {
      const payload: UninstallPluginRequest = { id };
      await instance.post<ApiResponse<void>>(UninstallPluginEndpoint, payload);
    } catch (e) {
      catchError(e);
    }
  }

  static async activate(id: string) {
    try {
      const payload: ActivatePluginRequest = { id };
      await instance.post<ApiResponse<void>>(ActivatePluginEndpoint, payload);
    } catch (e) {
      catchError(e);
    }
  }

  static async deactivate(id: string) {
    try {
      const payload: DeactivatePluginRequest = { id };
      await instance.post<ApiResponse<void>>(DeactivatePluginEndpoint, payload);
    } catch (e) {
      catchError(e);
    }
  }

  static async getStatus(id: string) {
    try {
      const payload: GetPluginStatusRequest = { id };
      const response = await instance.get<ApiResponse<GetPluginStatusResponse>>(
        GetPluginStatusEndpoint,
        { params: payload },
      );
      return response.data.data.status ?? null;
    } catch (e) {
      catchError(e);
    }
  }

  static async upload(file: File): Promise<UploadPluginResponse> {
    try {
      const formData = new FormData();
      formData.append("file", file);

      const response = await instance.post<ApiResponse<UploadPluginResponse>>(
        UploadPluginEndpoint,
        formData,
        {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        },
      );
      return response.data.data;
    } catch (e) {
      catchError(e);
    }
  }

  static async installFromAttachment(
    attachmentId: string,
  ): Promise<Plugin | null> {
    try {
      const payload = { attachmentId };
      const response = await instance.post<
        ApiResponse<InstallPluginFromAttachmentResponse>
      >(InstallPluginFromAttachmentEndpoint, payload);
      return response.data.data.plugin;
    } catch (e) {
      catchError(e);
    }
  }
}
