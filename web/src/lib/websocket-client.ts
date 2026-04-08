type MessageHandler = (data: any) => void;
export type ConnectionStatus = 'connected' | 'reconnecting' | 'disconnected';
type StatusHandler = (status: ConnectionStatus) => void;

const WS_URL = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws';

export class WebSocketClient {
  private ws: WebSocket | null = null;
  private handlers: Map<string, MessageHandler[]> = new Map();
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 10;
  private token: string | null = null;
  private pingInterval: ReturnType<typeof setInterval> | null = null;
  private _status: ConnectionStatus = 'disconnected';
  private statusHandlers: StatusHandler[] = [];

  get status() {
    return this._status;
  }

  private setStatus(status: ConnectionStatus) {
    this._status = status;
    this.statusHandlers.forEach((h) => h(status));
  }

  onStatusChange(handler: StatusHandler) {
    this.statusHandlers.push(handler);
    return () => {
      this.statusHandlers = this.statusHandlers.filter((h) => h !== handler);
    };
  }

  connect(token: string) {
    this.token = token;

    // Skip WebSocket if URL is not configured or points to a non-WS endpoint
    if (!WS_URL || WS_URL === 'ws://localhost:8080/ws') {
      this.setStatus('disconnected');
      return;
    }

    this.ws = new WebSocket(`${WS_URL}?user_id=${token}`);

    this.ws.onopen = () => {
      this.reconnectAttempts = 0;
      this.setStatus('connected');
      this.subscribe('market');
      this.subscribe('portfolio');
      this.startPing();
    };

    this.ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        const handlers = this.handlers.get(msg.type) || [];
        handlers.forEach((h) => h(msg.data));
      } catch {
        // Ignore parse errors
      }
    };

    this.ws.onclose = () => {
      this.stopPing();
      this.attemptReconnect();
    };

    this.ws.onerror = () => {
      this.ws?.close();
    };
  }

  disconnect() {
    this.stopPing();
    this.ws?.close();
    this.ws = null;
    this.setStatus('disconnected');
  }

  subscribe(channel: string) {
    this.send({ type: 'subscribe', channel });
  }

  unsubscribe(channel: string) {
    this.send({ type: 'unsubscribe', channel });
  }

  on(type: string, handler: MessageHandler) {
    if (!this.handlers.has(type)) {
      this.handlers.set(type, []);
    }
    this.handlers.get(type)!.push(handler);
  }

  off(type: string, handler: MessageHandler) {
    const handlers = this.handlers.get(type);
    if (handlers) {
      this.handlers.set(type, handlers.filter((h) => h !== handler));
    }
  }

  private send(data: object) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data));
    }
  }

  private attemptReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts || !this.token) {
      this.setStatus('disconnected');
      return;
    }
    this.setStatus('reconnecting');
    this.reconnectAttempts++;
    const delay = Math.min(Math.pow(2, this.reconnectAttempts) * 1000, 30000);
    setTimeout(() => this.connect(this.token!), delay);
  }

  private startPing() {
    this.stopPing();
    this.pingInterval = setInterval(() => {
      this.send({ type: 'ping' });
    }, 30000);
  }

  private stopPing() {
    if (this.pingInterval) {
      clearInterval(this.pingInterval);
      this.pingInterval = null;
    }
  }
}

export const wsClient = new WebSocketClient();
