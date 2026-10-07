import { instance } from "./lib/axios";
import {
  AddCredentialEndpoint,
  type AddCredentialRequest,
  type AddCredentialResponse,
  type ApiResponse,
  type Credential,
  DeleteCredentialEndpoint,
  type DeleteCredentialRequest,
  type GetListCredentialResponse,
  ListCredentialsEndpoint,
  UpdateCredentialEndpoint,
  type UpdateCredentialRequest,
  type UpdateCredentialResponse,
} from "@nagare-app/dtos";
import { catchError } from "./lib/errors";

export class CredentialService {
  static async list() {
    try {
      const response = await instance.get<
        ApiResponse<GetListCredentialResponse>
      >(ListCredentialsEndpoint);
      return (response.data.data.credentials ?? []).filter(
        (credential): credential is Credential => credential !== null,
      );
    } catch (e) {
      catchError(e);
    }
  }

  static async add(payload: AddCredentialRequest) {
    try {
      const response = await instance.post<ApiResponse<AddCredentialResponse>>(
        AddCredentialEndpoint,
        payload,
      );
      return response.data.data.credential!;
    } catch (e) {
      catchError(e);
    }
  }

  static async update(payload: UpdateCredentialRequest) {
    try {
      const response = await instance.post<
        ApiResponse<UpdateCredentialResponse>
      >(UpdateCredentialEndpoint, payload);
      return response.data.data.credential!;
    } catch (e) {
      catchError(e);
    }
  }

  static async delete(id: string) {
    try {
      const payload: DeleteCredentialRequest = { id };
      await instance.post<ApiResponse<void>>(DeleteCredentialEndpoint, payload);
    } catch (e) {
      catchError(e);
    }
  }
}
