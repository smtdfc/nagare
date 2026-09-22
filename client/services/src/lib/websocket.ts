import { envConfig } from "#/env";
import type { Event as AppEvent, Payload } from "@nagare-app/dtos";
import { WebsocketApiError } from "./errors";

type WebSocketEventListener = (data: any) => void;

interface WebSocketHelperOptions {
  reconnectInterval?: number;
}

export class WebSocketHelper {
  private url: string;
  private ws: WebSocket | null = null;
  private listeners: Map<string, WebSocketEventListener[]> = new Map();
  private reconnectInterval: number;
  private shouldReconnect: boolean = true;

  constructor(url: string, options: WebSocketHelperOptions = {}) {
    this.url = url;
    this.reconnectInterval = options.reconnectInterval || 3000;
  }

  public isConnected() {
    return this.ws && this.ws.readyState === WebSocket.OPEN;
  }

  public connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      this.shouldReconnect = true;
      this.ws = new WebSocket(this.url);

      this.ws.onopen = (event) => {
        console.log("WebSocket connected");
        this._trigger("open", event);
        resolve();
      };

      this.ws.onmessage = (event: MessageEvent) => {
        try {
          const data = JSON.parse(event.data) as Payload<any>;
          if (data.event) {
            this._trigger(data.event, data);
          }
          this._trigger("message", data);
        } catch (e) {
          this._trigger("message", event.data);
        }
      };

      this.ws.onerror = (error) => {
        console.error("WebSocket error:", error);
        this._trigger("error", error);
        reject(error);
      };

      this.ws.onclose = (event) => {
        console.log("WebSocket closed");
        this._trigger("close", event);

        if (this.shouldReconnect) {
          setTimeout(() => this.connect(), this.reconnectInterval);
        }
      };
    });
  }

  public on(event: AppEvent, callback: WebSocketEventListener): void {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, []);
    }
    this.listeners.get(event)!.push(callback);
  }

  public off(event: AppEvent, callback?: WebSocketEventListener): void {
    if (!this.listeners.has(event)) return;

    if (!callback) {
      this.listeners.delete(event);
      return;
    }

    const filtered = this.listeners.get(event)!.filter((cb) => cb !== callback);
    this.listeners.set(event, filtered);
  }

  public send<D = any>(
    event: AppEvent,
    data: D,
    requestID?: string,
  ): string | null {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      const id = requestID
        ? requestID
        : Math.random().toString(36).substring(2);
      const payload: Payload<D> = {
        requestID: id,
        event,
        data,
      };
      this.ws.send(JSON.stringify(payload));
      return id;
    } else {
      console.warn("WebSocket is not connected. Cannot send data.");
      return null;
    }
  }

  public disconnect(): void {
    this.shouldReconnect = false;
    if (this.ws) {
      this.ws.close();
    }
  }

  private _trigger(event: string, data: any): void {
    if (this.listeners.has(event)) {
      this.listeners.get(event)!.forEach((callback) => callback(data));
    }
  }

  public request<D = any, T = any, F = any>(
    sendEvent: AppEvent,
    data: D,
    successEvent: AppEvent,
    failEvent: AppEvent,
  ): Promise<T> {
    return new Promise((resolve, reject) => {
      const requestID = Math.random().toString(36).substring(2);
      const successHandler = (resData: Payload<T> | any) => {
        if (resData && resData.requestID === requestID) {
          cleanup();
          resolve(resData.data !== undefined ? resData.data : resData);
        }
      };

      const failHandler = (errData: Payload<F>) => {
        if (errData && errData.requestID === requestID) {
          cleanup();
          reject(new WebsocketApiError(errData.data as any));
        }
      };

      const cleanup = () => {
        this.off(successEvent, successHandler);
        this.off(failEvent, failHandler);
      };

      this.on(successEvent, successHandler);
      this.on(failEvent, failHandler);

      this.send(sendEvent, data, requestID);
    });
  }
}

export function createWebsocketWithCurrentEnv() {
  return new WebSocketHelper(envConfig.baseWebsocketUrl);
}

export function refreshWebsocket() {
  if (websocket) {
    websocket.disconnect();
  }
  websocket = createWebsocketWithCurrentEnv();
}

let websocket = createWebsocketWithCurrentEnv();

export { websocket };
