import { envConfig } from "../env";
import type { Event as AppEvent, Payload } from "@nagare-app/dtos";
import { WebsocketApiError } from "./errors";

type WebSocketEventListener = (data: any) => void;
type WebSocketReconnectHandler = () => void | Promise<void>;

export type WebSocketEventName =
  | AppEvent
  | "open"
  | "close"
  | "error"
  | "message";

interface WebSocketHelperOptions {
  reconnectInterval?: number;
  requestTimeout?: number;
}

export class WebSocketHelper {
  private url: string;
  private ws: WebSocket | null = null;
  private listeners: Map<string, WebSocketEventListener[]> = new Map();
  private reconnectHandlers: WebSocketReconnectHandler[] = [];
  private reconnectInterval: number;
  private requestTimeout: number;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private connectPromise: Promise<void> | null = null;
  private hasConnectedOnce = false;
  private shouldReconnect: boolean = true;

  constructor(url: string, options: WebSocketHelperOptions = {}) {
    this.url = url;
    this.reconnectInterval = options.reconnectInterval || 3000;
    this.requestTimeout = options.requestTimeout || 15000;
  }

  public isConnected() {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
  }

  public connect(): Promise<void> {
    if (this.ws?.readyState === WebSocket.OPEN) {
      return Promise.resolve();
    }
    if (this.connectPromise) {
      return this.connectPromise;
    }

    this.shouldReconnect = true;
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }

    const promise = new Promise<void>((resolve, reject) => {
      const socket = new WebSocket(this.url);
      this.ws = socket;

      socket.onopen = (event) => {
        if (this.ws !== socket) return;

        const isReconnect = this.hasConnectedOnce;
        this.hasConnectedOnce = true;

        console.log("WebSocket connected");
        this._trigger("open", event);
        resolve();

        if (isReconnect) {
          void this._runReconnectHandlers();
        }
      };

      socket.onmessage = (event: MessageEvent) => {
        if (this.ws !== socket) return;

        try {
          const data = JSON.parse(event.data) as Payload<any>;
          if (data.event) {
            this._trigger(data.event, data);
          }
          this._trigger("message", data);
        } catch {
          this._trigger("message", event.data);
        }
      };

      socket.onerror = (error) => {
        if (this.ws !== socket) return;

        console.error("WebSocket error:", error);
        this._trigger("error", error);
        reject(error);
      };

      socket.onclose = (event) => {
        if (this.ws !== socket) return;

        this.ws = null;
        console.log("WebSocket closed");
        this._trigger("close", event);

        if (this.shouldReconnect) {
          this.reconnectTimer = setTimeout(() => {
            this.reconnectTimer = null;
            void this.connect().catch(() => {});
          }, this.reconnectInterval);
        }
      };
    }).finally(() => {
      this.connectPromise = null;
    });

    this.connectPromise = promise;
    return promise;
  }

  public on(event: WebSocketEventName, callback: WebSocketEventListener): void {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, []);
    }
    this.listeners.get(event)!.push(callback);
  }

  public off(event: WebSocketEventName, callback?: WebSocketEventListener): void {
    if (!this.listeners.has(event)) return;

    if (!callback) {
      this.listeners.delete(event);
      return;
    }

    const filtered = this.listeners.get(event)!.filter((cb) => cb !== callback);
    this.listeners.set(event, filtered);
  }

  /**
   * Registers a handler that runs after an automatic reconnect. It is not
   * called for the first connection, because callers already authenticate
   * explicitly at startup. Handlers run sequentially, so an authentication
   * handler can finish before room subscriptions are restored.
   */
  public onReconnect(handler: WebSocketReconnectHandler): () => void {
    this.reconnectHandlers.push(handler);

    return () => {
      this.reconnectHandlers = this.reconnectHandlers.filter(
        (item) => item !== handler,
      );
    };
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

    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }

    if (this.ws) {
      const socket = this.ws;
      this.ws = null;
      socket.close();
    }
  }

  private _trigger(event: WebSocketEventName, data: any): void {
    if (this.listeners.has(event)) {
      this.listeners.get(event)!.forEach((callback) => callback(data));
    }
  }

  private async _runReconnectHandlers(): Promise<void> {
    for (const handler of [...this.reconnectHandlers]) {
      try {
        await handler();
      } catch (error) {
        console.error("WebSocket reconnect handler failed:", error);
      }
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
      let timeout: ReturnType<typeof setTimeout> | undefined;

      const cleanup = () => {
        if (timeout) {
          clearTimeout(timeout);
        }
        this.off(successEvent, successHandler);
        this.off(failEvent, failHandler);
      };

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

      timeout = setTimeout(() => {
        cleanup();
        reject(new Error("WebSocket request timed out"));
      }, this.requestTimeout);

      this.on(successEvent, successHandler);
      this.on(failEvent, failHandler);

      if (this.send(sendEvent, data, requestID) === null) {
        cleanup();
        reject(new Error("WebSocket is not connected"));
      }
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
