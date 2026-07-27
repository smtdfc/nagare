import { WebSocketHelper } from "#/lib/websocket.js";
import { GetWebsocketConnect } from "../../wailsjs/go/main/App.js"

let ws: WebSocketHelper;

export async function initWebsocketConnection() {
    const wsUrl = await GetWebsocketConnect()
    ws = new WebSocketHelper(wsUrl);
    ws.connect();
}

export class ChatService {
    static async createChatSession() {
    }
}