/**
 * WebSocket Service
 * Handles real-time bidirectional communication with the server
 * Used for: Bidding updates, Notifications, Live market prices
 */

type MessageHandler = (data: any) => void;

interface WebSocketMessage {
  type: string;
  data: any;
  timestamp?: string;
}

class WebSocketService {
  private ws: WebSocket | null = null;
  private url: string = 'ws://localhost:8080/ws';
  private reconnectAttempts: number = 0;
  private maxReconnectAttempts: number = 5;
  private reconnectDelay: number = 3000;
  private messageHandlers: Map<string, MessageHandler[]> = new Map();
  private isConnecting: boolean = false;

  /**
   * Connect to WebSocket server with JWT token
   */
  connect(token: string): Promise<void> {
    return new Promise((resolve, reject) => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        console.log('✅ WebSocket already connected');
        resolve();
        return;
      }

      if (this.isConnecting) {
        console.log('⏳ WebSocket connection already in progress');
        return;
      }

      this.isConnecting = true;
      const wsUrl = `${this.url}?token=${token}`;

      console.log('🔌 Connecting to WebSocket:', wsUrl);

      try {
        this.ws = new WebSocket(wsUrl);

        this.ws.onopen = () => {
          console.log('✅ WebSocket connected');
          this.reconnectAttempts = 0;
          this.isConnecting = false;
          resolve();
        };

        this.ws.onmessage = (event) => {
          try {
            const message: WebSocketMessage = JSON.parse(event.data);
            console.log('📨 WebSocket message:', message);
            this.handleMessage(message);
          } catch (error) {
            console.error('❌ Failed to parse WebSocket message:', error);
          }
        };

        this.ws.onerror = (error) => {
          console.error('❌ WebSocket error:', error);
          this.isConnecting = false;
          reject(error);
        };

        this.ws.onclose = () => {
          console.log('🔌 WebSocket disconnected');
          this.isConnecting = false;
          this.ws = null;
          this.attemptReconnect(token);
        };
      } catch (error) {
        console.error('❌ Failed to create WebSocket:', error);
        this.isConnecting = false;
        reject(error);
      }
    });
  }

  /**
   * Disconnect from WebSocket server
   */
  disconnect(): void {
    if (this.ws) {
      console.log('🔌 Disconnecting WebSocket...');
      this.ws.close();
      this.ws = null;
    }
    this.messageHandlers.clear();
  }

  /**
   * Send message to server
   */
  send(type: string, data: any): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      const message: WebSocketMessage = {
        type,
        data,
        timestamp: new Date().toISOString(),
      };
      this.ws.send(JSON.stringify(message));
      console.log('📤 Sent WebSocket message:', message);
    } else {
      console.error('❌ WebSocket is not connected');
    }
  }

  /**
   * Subscribe to specific message type
   */
  on(type: string, handler: MessageHandler): void {
    if (!this.messageHandlers.has(type)) {
      this.messageHandlers.set(type, []);
    }
    this.messageHandlers.get(type)?.push(handler);
    console.log(`📡 Subscribed to "${type}" messages`);
  }

  /**
   * Unsubscribe from specific message type
   */
  off(type: string, handler?: MessageHandler): void {
    if (!handler) {
      this.messageHandlers.delete(type);
      console.log(`📡 Unsubscribed from all "${type}" messages`);
    } else {
      const handlers = this.messageHandlers.get(type);
      if (handlers) {
        const index = handlers.indexOf(handler);
        if (index > -1) {
          handlers.splice(index, 1);
        }
      }
      console.log(`📡 Unsubscribed from "${type}" message`);
    }
  }

  /**
   * Handle incoming message
   */
  private handleMessage(message: WebSocketMessage): void {
    const handlers = this.messageHandlers.get(message.type);
    if (handlers) {
      handlers.forEach((handler) => handler(message.data));
    }
  }

  /**
   * Attempt to reconnect
   */
  private attemptReconnect(token: string): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('❌ Max reconnect attempts reached');
      return;
    }

    this.reconnectAttempts++;
    console.log(
      `🔄 Reconnecting... Attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts}`
    );

    setTimeout(() => {
      this.connect(token).catch((error) => {
        console.error('❌ Reconnect failed:', error);
      });
    }, this.reconnectDelay);
  }

  /**
   * Check if WebSocket is connected
   */
  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }

  /**
   * Get connection state
   */
  getState(): string {
    if (!this.ws) return 'CLOSED';
    switch (this.ws.readyState) {
      case WebSocket.CONNECTING:
        return 'CONNECTING';
      case WebSocket.OPEN:
        return 'OPEN';
      case WebSocket.CLOSING:
        return 'CLOSING';
      case WebSocket.CLOSED:
        return 'CLOSED';
      default:
        return 'UNKNOWN';
    }
  }
}

// Singleton instance
export const websocketService = new WebSocketService();

// Message types (matches backend)
export const WS_MESSAGE_TYPES = {
  // Bidding
  BID_PLACED: 'BID_PLACED',
  ASK_PLACED: 'ASK_PLACED',
  BID_UPDATED: 'BID_UPDATED',
  ASK_UPDATED: 'ASK_UPDATED',
  MATCH_CREATED: 'MATCH_CREATED',
  MARKET_PRICE_UPDATED: 'MARKET_PRICE_UPDATED',
  
  // Notifications
  NOTIFICATION: 'NOTIFICATION',
  
  // Orders
  ORDER_STATUS_UPDATED: 'ORDER_STATUS_UPDATED',
  
  // Products
  PRODUCT_UPDATED: 'PRODUCT_UPDATED',
  NEW_PRODUCT: 'NEW_PRODUCT',
} as const;

export type WSMessageType = typeof WS_MESSAGE_TYPES[keyof typeof WS_MESSAGE_TYPES];
