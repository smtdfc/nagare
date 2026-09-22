import type { ApiError as ApiResponseError } from "@nagare-app/dtos";
import { AxiosError } from "axios";

export class ApiError extends Error {
  constructor(public apiError: ApiResponseError) {
    super(`${apiError.code}: ${apiError.message}`);
  }
}

export class WebsocketApiError<D extends { cause: string }> extends Error {
  constructor(public apiError: D) {
    super(`${apiError.cause}`);
  }
}
export function catchError(e: unknown): never {
  if (e instanceof WebsocketApiError) {
    throw new ApiError({
      statusCode: 400,
      code: "WS_ERROR",
      message: e.message || "Cannot connect to the server.",
    } as ApiResponseError);
  }

  if (e instanceof AxiosError) {
    if (e.response && e.response.data && e.response.data.isSuccess === false) {
      throw new ApiError(e.response.data.error as ApiResponseError);
    }

    throw new ApiError({
      statusCode: e.response?.status || 500,
      code: e.code || "NETWORK_ERROR",
      message: e.message || "Cannot connect to the server.",
    } as ApiResponseError);
  }

  throw new ApiError({
    code: "UNKNOWN_ERROR",
    message: e instanceof Error ? e.message : "Unknown error",
  } as ApiResponseError);
}
