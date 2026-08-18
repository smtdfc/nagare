import type { WsEvent, WsMessage } from "@nagare-agent/dto";
import { GetWebsocketConnect } from "@nagare-agent/service-bindings";
import { waitForWails } from "./wails";

type EventCallback<T> = (data: T) => void;

type BuiltInWebsocket = {
  chat: WebSocketHelper | null;
};

const builtInWs: BuiltInWebsocket = {
  chat: null,
};

export async function initChatWebsocketConnection() {
  await waitForWails();
  const wsUrl = await GetWebsocketConnect();
  builtInWs.chat = new WebSocketHelper(wsUrl);
  builtInWs.chat.connect();
}

export { builtInWs };

export class WebSocketHelper {
  private url: string;
  private ws: WebSocket | null = null;
  private listeners: Map<string, Set<EventCallback<any>>> = new Map();
  private isManualClose = false;
  private reconnectInterval = 3000;

  constructor(url: string) {
    this.url = url;
  }

  public connect() {
    this.isManualClose = false;
    this.ws = new WebSocket(this.url);

    this.ws.onopen = (event) => {
      console.log("WebSocket connected:", this.url);
      this.trigger("open", event);
    };

    this.ws.onmessage = async (event) => {
      try {
        let rawData = event.data;
        if (rawData instanceof Blob) {
          rawData = await rawData.text();
        }

        const message = JSON.parse(rawData) as WsMessage<any>;
        const { event: eventName, payload } = message;
        if (eventName && this.listeners.has(eventName)) {
          this.listeners.get(eventName)?.forEach((callback) => {
            callback(payload);
          });
        }
      } catch (error) {
        console.error("Failed to parse WebSocket message:", event.data);
      }
    };

    this.ws.onerror = (error) => {
      console.error("WebSocket error:", error);
      this.trigger("error", error);
    };

    this.ws.onclose = (event) => {
      console.log("WebSocket closed:", event);
      this.trigger("close", event);

      if (!this.isManualClose) {
        setTimeout(() => {
          this.connect();
        }, this.reconnectInterval);
      }
    };
  }

  public on<TE extends WsEvent, TP>(
    eventName: TE,
    callback: EventCallback<TP>,
  ): void {
    if (!this.listeners.has(eventName)) {
      this.listeners.set(eventName, new Set());
    }
    this.listeners.get(eventName)?.add(callback);
  }

  public off<TE extends WsEvent, TP>(
    eventName: TE,
    callback?: EventCallback<TP>,
  ): void {
    if (!this.listeners.has(eventName)) return;

    if (callback) {
      this.listeners.get(eventName)?.delete(callback);
      if (this.listeners.get(eventName)?.size === 0) {
        this.listeners.delete(eventName);
      }
    } else {
      this.listeners.delete(eventName);
    }
  }

  public send<T>(eventName: string, data: T): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      const message = JSON.stringify({ event: eventName, payload: data });
      this.ws.send(message);
    } else {
      console.warn("WebSocket is not open. Cannot send message.");
    }
  }

  public disconnect(): void {
    this.isManualClose = true;
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  private trigger<T>(eventName: string, data: T): void {
    if (this.listeners.has(eventName)) {
      this.listeners.get(eventName)?.forEach((callback) => {
        callback(data);
      });
    }
  }
}

export function wsRequest<
  TResponse,
  TError = any,
  TSendEvent extends WsEvent = WsEvent,
  TSuccessEvent extends WsEvent = WsEvent,
  TFailedEvent extends WsEvent = WsEvent,
>(
  ws: WebSocketHelper,
  sendEvent: TSendEvent,
  successEvent: TSuccessEvent,
  failureEvent: TFailedEvent,
  payload: any = {},
  timeoutMs: number = 10000,
): Promise<TResponse> {
  return new Promise((resolve, reject) => {
    let timeoutId: ReturnType<typeof setTimeout> | null = null;

    const cleanup = () => {
      ws.off(successEvent, handleSuccess);
      ws.off(failureEvent, handleFailure);
      if (timeoutId) clearTimeout(timeoutId);
    };

    const handleSuccess = (data: TResponse) => {
      cleanup();
      resolve(data);
    };

    const handleFailure = (data: TError) => {
      cleanup();
      reject(data);
    };

    ws.on(successEvent, handleSuccess);
    ws.on(failureEvent, handleFailure);

    timeoutId = setTimeout(() => {
      cleanup();
      reject(new Error(`WebSocket request timeout for event: ${sendEvent}`));
    }, timeoutMs);

    ws.send(sendEvent, payload);
  });
}
