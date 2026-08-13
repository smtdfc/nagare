import type { ApiResponse } from "#/dto/api.ts";
import { AxiosError } from "axios";

export function generateErrorCode(err: unknown): string {
    const errorMetadata = {
        stack: "",
        name: "",
        message: "",
        timestamp: Date.now()
    }

    if (err instanceof Error) {
        errorMetadata.stack = err.stack || "No stack provided";
        errorMetadata.name = err.name;
        errorMetadata.message = err.message;
    } else {
        errorMetadata.message = String(err);
        errorMetadata.name = "UnknownError";
    }
    const jsonString = JSON.stringify(errorMetadata);
    if (typeof window === "undefined") {
        return Buffer.from(jsonString).toString("base64");
    } else {
        return btoa(unescape(encodeURIComponent(jsonString)));
    }
}

export function handleError(e: unknown): never {
    const errorCode = generateErrorCode(e);
    if (e instanceof AxiosError) {
        const resp = e.response?.data as ApiResponse<any>;
        if (!resp) throw new Error(`Network error! Code: ${errorCode}`);
        if (resp.error) {
            const apiError = resp.error;
            throw new Error(apiError.message);
        }
    }

    throw new Error(`Unknown error! Code: ${errorCode}`);
}

export function getErrorMessage(err: unknown): string {
    if (err instanceof Error) {
        return err.message;
    }

    if (typeof err === "string") {
        return err;
    }

    if (err && typeof err === "object" && "message" in err && typeof (err as Record<string, unknown>).message === "string") {
        return (err as { message: string }).message;
    }

    try {
        return JSON.stringify(err);
    } catch {
        return "An unknown error occurred";
    }
}