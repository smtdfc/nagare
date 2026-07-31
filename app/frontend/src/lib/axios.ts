import { GenerateToken, GetRestApiConnect } from "@wails/go/main/App";
import type { Axios } from "axios";
import axios from "axios";
import { waitForWails } from "./wails";

let token: string | null = null;
let instance: Axios;

export async function getAxiosInstance(): Promise<Axios> {
    if (instance) return instance;

    if (token === null) {
        await waitForWails(10, 2000);
        token = await GenerateToken();
    }

    instance = axios.create({
        baseURL: await GetRestApiConnect(),
        timeout: 5000,
        headers: { "X-Nagare-Secure": token },
    });

    return instance;
}
