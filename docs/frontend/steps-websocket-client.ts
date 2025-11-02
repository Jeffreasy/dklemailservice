/**
 * StepsWebSocketClient - TypeScript WebSocket client voor stappen tracking
 * 
 * Features:
 * - Auto-reconnect met exponential backoff
 * - Subscription-based event handling
 * - Type-safe message types
 * - Connection state management
 * - Ping/pong keep-alive
 * 
 * Usage:
 * ```typescript
 * const client = new StepsWebSocketClient(
 *   'ws://localhost:8080/ws/steps',
 *   'jwt-token',
 *   'user-id',
 *   'participant-id'
 * );
 * 
 * client.on('step_update', (data) => {
 *   console.log('Steps updated:', data.steps);
 * });
 * 
 * client.connect();
 * ```
 */

// Message Types
export type MessageType =
  | 'step_update'
  | 'total_update'
  | 'leaderboard_update'
  | 'badge_earned'
  | 'subscribe'
  | 'unsubscribe'
  | 'ping'
  | 'pong';

// Message Interfaces
export interface StepUpdateMessage {
  type: 'step_update';
  participant_id: string;
  naam: string;
  steps: number;
  delta: number;
  route: string;
  allocated_funds: number;
  timestamp: number;
}

export interface TotalUpdateMessage {
  type: 'total_update';
  total_steps: number;
  year: number;
  timestamp: number;
}

export interface LeaderboardEntry {
  rank: number;
  participant_id: string;
  naam: string;
  steps: number;
  achievement_points: number;
  total_score: number;
  route: string;
  badge_count: number;
}

export interface LeaderboardUpdateMessage {
  type: 'leaderboard_update';
  top_n: number;
  entries: LeaderboardEntry[];
  timestamp: number;
}

export interface BadgeEarnedMessage {
  type: 'badge_earned';
  participant_id: string;
  badge_name: string;
  badge_icon: string;
  points: number;
  timestamp: number;
}

export interface PingMessage {
  type: 'ping';
  timestamp: number;
}

export interface PongMessage {
  type: 'pong';
  timestamp: number;
}

export type WebSocketMessage =
  | StepUpdateMessage
  | TotalUpdateMessage
  | LeaderboardUpdateMessage
  | BadgeEarnedMessage
  | PingMessage
  | PongMessage;

// Connection States
export enum ConnectionState {
  DISCONNECTED = 'DISCONNECTED',
  CONNECTING = 'CONNECTING',
  CONNECTED = 'CONNECTED',
  RECONNECTING = 'RECONNECTING',
  FAILED = 'FAILED',
}

// Event Callback Type
type EventCallback<T = any> = (data: T) => void;

// Client Configuration
export interface StepsWebSocketConfig {
  reconnectInterval?: number;
  maxReconnectInterval?: number;
  reconnectDecay?: number;
  maxReconnectAttempts?: number;
  pingInterval?: number;
  debug?: boolean;
}

const DEFAULT_CONFIG: Required<StepsWebSocketConfig> = {
  reconnectInterval: 1000,
  maxReconnectInterval: 30000,
  reconnectDecay: 1.5,
  maxReconnectAttempts: Infinity,
  pingInterval: 30000,
  debug: false,
};

/**
 * StepsWebSocketClient - Main WebSocket client class
 */
export class StepsWebSocketClient {
  private ws: WebSocket | null = null;
  private url: string;
  private token: string;
  private userId: string;
  private participantId?: string;
  
  private config: Required<StepsWebSocketConfig>;
  private callbacks: Map<MessageType, EventCallback[]> = new Map();
  private stateCallbacks: Map<ConnectionState, EventCallback[]> = new Map();
  
  private reconnectAttempts: number = 0;
  private reconnectTimeout: number;
  private pingInterval: number | null = null;
  
  private _state: ConnectionState = ConnectionState.DISCONNECTED;
  private subscriptions: Set<string> = new Set();

  constructor(
    url: string,
    token: string,
    userId: string,
    participantId?: string,
    config?: StepsWebSocketConfig
  ) {
    this.url = url;
    this.token = token;
    this.userId = userId;
    this.participantId = participantId;
    this.config = { ...DEFAULT_CONFIG, ...config };
    this.reconnectTimeout = this.config.reconnectInterval;
  }

  /**
   * Current connection state
   */
  get state(): ConnectionState {
    return this._state;
  }

  private setState(newState: ConnectionState): void {
    if (this._state !== newState) {
      this._state = newState;
      this.log(`State changed to: ${newState}`);
      this.emitStateChange(newState);
    }
  }

  /**
   * Check if client is connected
   */
  get isConnected(): boolean {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
  }

  /**
   * Connect to WebSocket server
   */
  connect(): void {
    if (this.ws !== null && this.ws.readyState === WebSocket.OPEN) {
      this.log('Already connected');
      return;
    }

    this.setState(
      this.reconnectAttempts > 0
        ? ConnectionState.RECONNECTING
        : ConnectionState.CONNECTING
    );

    const queryParams = new URLSearchParams({
      user_id: this.userId,
      ...(this.token && { token: this.token }),
      ...(this.participantId && { participant_id: this.participantId }),
    });

    const wsUrl = `${this.url}?${queryParams}`;
    this.log(`Connecting to: ${wsUrl}`);

    try {
      this.ws = new WebSocket(wsUrl);
      this.setupEventHandlers();
    } catch (error) {
      this.log(`Connection error: ${error}`, true);
      this.handleConnectionError();
    }
  }

  /**
   * Disconnect from WebSocket server
   */
  disconnect(): void {
    this.log('Disconnecting...');
    this.stopPingInterval();
    
    if (this.ws) {
      this.ws.close(1000, 'Client disconnect');
      this.ws = null;
    }
    
    this.setState(ConnectionState.DISCONNECTED);
  }

  /**
   * Subscribe to message types/channels
   */
  subscribe(channels: string[]): void {
    if (!this.isConnected) {
      this.log('Cannot subscribe, not connected', true);
      return;
    }

    channels.forEach(channel => this.subscriptions.add(channel));

    this.send({
      type: 'subscribe',
      channels: channels,
    });

    this.log(`Subscribed to: ${channels.join(', ')}`);
  }

  /**
   * Unsubscribe from message types/channels
   */
  unsubscribe(channels: string[]): void {
    if (!this.isConnected) {
      return;
    }

    channels.forEach(channel => this.subscriptions.delete(channel));

    this.send({
      type: 'unsubscribe',
      channels: channels,
    });

    this.log(`Unsubscribed from: ${channels.join(', ')}`);
  }

  /**
   * Register event listener for specific message type
   */
  on<T extends WebSocketMessage>(
    messageType: T['type'],
    callback: EventCallback<T>
  ): void {
    if (!this.callbacks.has(messageType)) {
      this.callbacks.set(messageType, []);
    }
    this.callbacks.get(messageType)!.push(callback as EventCallback);
  }

  /**
   * Register event listener for connection state changes
   */
  onStateChange(callback: EventCallback<ConnectionState>): void {
    // Store in a special 'state' key
    const stateKey = 'state' as any;
    if (!this.callbacks.has(stateKey)) {
      this.callbacks.set(stateKey, []);
    }
    this.callbacks.get(stateKey)!.push(callback);
  }

  /**
   * Remove event listener
   */
  off(messageType: MessageType, callback: EventCallback): void {
    const callbacks = this.callbacks.get(messageType);
    if (callbacks) {
      const index = callbacks.indexOf(callback);
      if (index > -1) {
        callbacks.splice(index, 1);
      }
    }
  }

  /**
   * Remove all event listeners for a message type
   */
  removeAllListeners(messageType?: MessageType): void {
    if (messageType) {
      this.callbacks.delete(messageType);
    } else {
      this.callbacks.clear();
    }
  }

  /**
   * Send raw message to server
   */
  private send(message: any): void {
    if (!this.isConnected) {
      this.log('Cannot send, not connected', true);
      return;
    }

    try {
      this.ws!.send(JSON.stringify(message));
    } catch (error) {
      this.log(`Send error: ${error}`, true);
    }
  }

  /**
   * Setup WebSocket event handlers
   */
  private setupEventHandlers(): void {
    if (!this.ws) return;

    this.ws.onopen = () => {
      this.log('Connected');
      this.setState(ConnectionState.CONNECTED);
      this.reconnectAttempts = 0;
      this.reconnectTimeout = this.config.reconnectInterval;
      
      // Resubscribe to previous subscriptions
      if (this.subscriptions.size > 0) {
        this.subscribe(Array.from(this.subscriptions));
      }
      
      // Start ping interval
      this.startPingInterval();
    };

    this.ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data) as WebSocketMessage;
        this.handleMessage(message);
      } catch (error) {
        this.log(`Parse error: ${error}`, true);
      }
    };

    this.ws.onerror = (error) => {
      this.log(`WebSocket error: ${error}`, true);
    };

    this.ws.onclose = (event) => {
      this.log(`Connection closed: ${event.code} - ${event.reason}`);
      this.stopPingInterval();
      this.ws = null;
      
      if (event.code !== 1000) {
        // Abnormal closure, attempt reconnect
        this.handleConnectionError();
      } else {
        this.setState(ConnectionState.DISCONNECTED);
      }
    };
  }

  /**
   * Handle incoming WebSocket message
   */
  private handleMessage(message: WebSocketMessage): void {
    this.log(`Received: ${message.type}`);

    // Handle pong
    if (message.type === 'pong') {
      this.log('Pong received');
      return;
    }

    // Emit to registered callbacks
    const callbacks = this.callbacks.get(message.type);
    if (callbacks) {
      callbacks.forEach(callback => {
        try {
          callback(message);
        } catch (error) {
          this.log(`Callback error: ${error}`, true);
        }
      });
    }
  }

  /**
   * Handle connection error and reconnect logic
   */
  private handleConnectionError(): void {
    if (this.reconnectAttempts >= this.config.maxReconnectAttempts) {
      this.log('Max reconnect attempts reached', true);
      this.setState(ConnectionState.FAILED);
      return;
    }

    this.reconnectAttempts++;
    this.setState(ConnectionState.RECONNECTING);

    const delay = Math.min(
      this.reconnectTimeout,
      this.config.maxReconnectInterval
    );

    this.log(`Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`);

    setTimeout(() => {
      this.reconnectTimeout = Math.min(
        this.reconnectTimeout * this.config.reconnectDecay,
        this.config.maxReconnectInterval
      );
      this.connect();
    }, delay);
  }

  /**
   * Start ping interval to keep connection alive
   */
  private startPingInterval(): void {
    this.stopPingInterval();
    
    this.pingInterval = setInterval(() => {
      if (this.isConnected) {
        this.send({
          type: 'ping',
          timestamp: Date.now(),
        });
      }
    }, this.config.pingInterval);
  }

  /**
   * Stop ping interval
   */
  private stopPingInterval(): void {
    if (this.pingInterval) {
      clearInterval(this.pingInterval);
      this.pingInterval = null;
    }
  }

  /**
   * Emit state change to listeners
   */
  private emitStateChange(state: ConnectionState): void {
    const stateKey = 'state' as any;
    const callbacks = this.callbacks.get(stateKey);
    if (callbacks) {
      callbacks.forEach(callback => callback(state));
    }
  }

  /**
   * Log helper
   */
  private log(message: string, isError: boolean = false): void {
    if (this.config.debug) {
      const prefix = '[StepsWebSocket]';
      if (isError) {
        console.error(prefix, message);
      } else {
        console.log(prefix, message);
      }
    }
  }
}

export default StepsWebSocketClient;