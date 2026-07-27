import type { WsEvent, WsMessage } from "#/dto/api.ts";

type EventCallback<T> = (data: T) => void;

export class WebSocketHelper {
    private url: string;
    private ws: WebSocket | null = null;
    private listeners: Map<string, Set<EventCallback<any>>> = new Map();
    private isManualClose = false;
    private reconnectInterval = 3000;

    constructor(url: string) {
        this.url = url;
    }


    public connect(): void {
        this.isManualClose = false;
        this.ws = new WebSocket(this.url);

        this.ws.onopen = (event) => {
            console.log("WebSocket connected:", this.url);
            this.trigger("open", event);
        };

        this.ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data) as WsMessage<any>;
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

    public on<TE extends WsEvent, TP>(eventName: TE, callback: EventCallback<TP>): void {
        if (!this.listeners.has(eventName)) {
            this.listeners.set(eventName, new Set());
        }
        this.listeners.get(eventName)?.add(callback);
    }

    public off<TE extends WsEvent, TP>(eventName: TE, callback?: EventCallback<TP>): void {
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
            this.ws.send(JSON.stringify({ event: eventName, payload: data }));
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