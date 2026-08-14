import { WsEvent } from "@nagare-agent/dto";
import type {
  AgentOuput,
  CreateSessionSuccess,
  InvokeAgent,
  InvokeAgentFailed,
} from "@nagare-agent/dto";
import { wsRequest, builtInWs } from "#/lib/websocket.ts";

export class ChatService {
  static listenMessage(
    cb: (isSuccess: boolean, data?: Message, err?: string) => void,
  ) {
    const onSuccess = (d: AgentOuput) => {
      cb(false, d.message as Message, undefined);
    };

    const onFailed = (d: InvokeAgentFailed) => {
      cb(true, undefined, d.cause);
    };

    builtInWs.chat!.on(WsEvent.WS_INVOKE_AGENT_FAILED, onFailed);
    builtInWs.chat!.on(WsEvent.WS_AGENT_RESPONSE, onSuccess);

    return () => {
      builtInWs.chat!.off(WsEvent.WS_AGENT_RESPONSE, onSuccess);
      builtInWs.chat!.off(WsEvent.WS_INVOKE_AGENT_FAILED, onFailed);
    };
  }

  static async sendMessage(id: string, message: string) {
    builtInWs.chat!.send<InvokeAgent>(WsEvent.WS_INVOKE_AGENT, {
      id: Date.now().toString(32),
      session_id: id,
      text: message,
    });
  }

  static async createChatSession() {
    if (!builtInWs.chat) throw new Error("Websocket is not ready");
    return (
      await wsRequest<CreateSessionSuccess>(
        builtInWs.chat,
        WsEvent.WS_CREATE_SESSION,
        WsEvent.WS_CREATE_SESSION_SUCCESS,
        WsEvent.WS_CREATE_SESSION_FAILED,
        {},
      )
    ).id;
  }
}
