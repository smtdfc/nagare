import { WsEvent } from "#/dto/api.js";
import type { AgentResponse, CreateSessionSuccess, InvokeAgent, InvokeAgentFailed } from "#/dto/api.js";
import { WebSocketHelper, wsRequest } from "#/lib/websocket.js";
import { GetWebsocketConnect } from "@wails/go/main/App.js"


let ws: WebSocketHelper;

export async function initWebsocketConnection() {
    const wsUrl = await GetWebsocketConnect()
    ws = new WebSocketHelper(wsUrl);
    ws.connect();
}

export class ChatService {

    static listenMessage(cb: (isSuccess: boolean, data?: Message, err?: string) => void) {

        const onSuccess = (d: AgentResponse) => {
            cb(false, d.message as Message, undefined);
        };

        const onFailed = (d: InvokeAgentFailed) => {
            cb(true, undefined, d.cause);
        };

        ws.on(WsEvent.WS_INVOKE_AGENT_FAILED, onFailed);
        ws.on(WsEvent.WS_AGENT_RESPONSE, onSuccess);

        return () => {
            ws.off(WsEvent.WS_AGENT_RESPONSE, onSuccess);
            ws.off(WsEvent.WS_INVOKE_AGENT_FAILED, onFailed);
        }
    }

    static async sendMessage(id: string, message: string) {
        ws.send<InvokeAgent>(WsEvent.WS_INVOKE_AGENT, {
            id: Date.now().toString(32),
            session_id: id,
            text: message,
        })
    }

    static async createChatSession() {
        return (await wsRequest<CreateSessionSuccess>(
            ws,
            WsEvent.WS_CREATE_SESSION,
            WsEvent.WS_CREATE_SESSION_SUCCESS,
            WsEvent.WS_CREATE_SESSION_FAILED,
            {}
        )).id;
    }
}