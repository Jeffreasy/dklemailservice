/**
 * useStepsWebSocket - React Hook voor WebSocket stappen tracking
 * 
 * Features:
 * - Automatic connection management
 * - Real-time step updates
 * - Type-safe state management
 * - Cleanup on unmount
 * 
 * Usage:
 * ```tsx
 * function Dashboard() {
 *   const {
 *     connected,
 *     latestUpdate,
 *     totalSteps,
 *     leaderboard,
 *     subscribe,
 *     unsubscribe
 *   } = useStepsWebSocket('user-123', 'participant-456');
 *   
 *   return <div>Steps: {latestUpdate?.steps || 0}</div>;
 * }
 * ```
 */

import { useEffect, useState, useRef, useCallback } from 'react';
import {
  StepsWebSocketClient,
  StepUpdateMessage,
  TotalUpdateMessage,
  LeaderboardUpdateMessage,
  BadgeEarnedMessage,
  ConnectionState,
  StepsWebSocketConfig,
} from './steps-websocket-client';

// Hook State Interface
export interface StepsWebSocketState {
  connected: boolean;
  connectionState: ConnectionState;
  latestUpdate: StepUpdateMessage | null;
  totalSteps: number;
  leaderboard: LeaderboardUpdateMessage | null;
  latestBadge: BadgeEarnedMessage | null;
}

// Hook Return Type
export interface UseStepsWebSocketReturn extends StepsWebSocketState {
  subscribe: (channels: string[]) => void;
  unsubscribe: (channels: string[]) => void;
  reconnect: () => void;
  disconnect: () => void;
}

/**
 * Custom hook for WebSocket connection to steps tracking
 */
export function useStepsWebSocket(
  userId: string,
  participantId?: string,
  config?: StepsWebSocketConfig
): UseStepsWebSocketReturn {
  const [state, setState] = useState<StepsWebSocketState>({
    connected: false,
    connectionState: ConnectionState.DISCONNECTED,
    latestUpdate: null,
    totalSteps: 0,
    leaderboard: null,
    latestBadge: null,
  });

  const clientRef = useRef<StepsWebSocketClient | null>(null);

  // Get WebSocket URL and token
  const getWebSocketUrl = useCallback(() => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    return `${protocol}//${host}/ws/steps`;
  }, []);

  const getToken = useCallback(() => {
    return localStorage.getItem('token') || '';
  }, []);

  // Initialize WebSocket client
  useEffect(() => {
    if (!userId) {
      console.warn('[useStepsWebSocket] No userId provided');
      return;
    }

    const url = getWebSocketUrl();
    const token = getToken();

    const client = new StepsWebSocketClient(
      url,
      token,
      userId,
      participantId,
      {
        debug: process.env.NODE_ENV === 'development',
        ...config,
      }
    );

    // Setup event listeners
    client.onStateChange((connectionState) => {
      setState((prev) => ({
        ...prev,
        connected: connectionState === ConnectionState.CONNECTED,
        connectionState,
      }));
    });

    client.on<StepUpdateMessage>('step_update', (data) => {
      setState((prev) => ({
        ...prev,
        latestUpdate: data,
      }));
    });

    client.on<TotalUpdateMessage>('total_update', (data) => {
      setState((prev) => ({
        ...prev,
        totalSteps: data.total_steps,
      }));
    });

    client.on<LeaderboardUpdateMessage>('leaderboard_update', (data) => {
      setState((prev) => ({
        ...prev,
        leaderboard: data,
      }));
    });

    client.on<BadgeEarnedMessage>('badge_earned', (data) => {
      setState((prev) => ({
        ...prev,
        latestBadge: data,
      }));
    });

    // Connect
    client.connect();

    // Store reference
    clientRef.current = client;

    // Cleanup on unmount
    return () => {
      client.disconnect();
      clientRef.current = null;
    };
  }, [userId, participantId, getWebSocketUrl, getToken, config]);

  // Subscribe to channels
  const subscribe = useCallback((channels: string[]) => {
    clientRef.current?.subscribe(channels);
  }, []);

  // Unsubscribe from channels
  const unsubscribe = useCallback((channels: string[]) => {
    clientRef.current?.unsubscribe(channels);
  }, []);

  // Manual reconnect
  const reconnect = useCallback(() => {
    clientRef.current?.disconnect();
    clientRef.current?.connect();
  }, []);

  // Manual disconnect
  const disconnect = useCallback(() => {
    clientRef.current?.disconnect();
  }, []);

  return {
    ...state,
    subscribe,
    unsubscribe,
    reconnect,
    disconnect,
  };
}

/**
 * Hook for participant dashboard with auto-subscription
 */
export function useParticipantDashboard(
  userId: string,
  participantId: string,
  config?: StepsWebSocketConfig
): UseStepsWebSocketReturn {
  const hook = useStepsWebSocket(userId, participantId, config);

  // Auto-subscribe to participant-relevant channels
  useEffect(() => {
    if (hook.connected) {
      hook.subscribe(['step_updates', 'badge_earned']);
    }
  }, [hook.connected, hook.subscribe]);

  return hook;
}

/**
 * Hook for public leaderboard view
 */
export function useLeaderboard(
  config?: StepsWebSocketConfig
): Pick<UseStepsWebSocketReturn, 'connected' | 'connectionState' | 'totalSteps' | 'leaderboard' | 'reconnect'> {
  const hook = useStepsWebSocket('public', undefined, config);

  // Auto-subscribe to public channels
  useEffect(() => {
    if (hook.connected) {
      hook.subscribe(['total_updates', 'leaderboard_updates']);
    }
  }, [hook.connected, hook.subscribe]);

  return {
    connected: hook.connected,
    connectionState: hook.connectionState,
    totalSteps: hook.totalSteps,
    leaderboard: hook.leaderboard,
    reconnect: hook.reconnect,
  };
}

/**
 * Hook for admin monitoring
 */
export function useStepsMonitoring(
  userId: string,
  config?: StepsWebSocketConfig
): UseStepsWebSocketReturn {
  const hook = useStepsWebSocket(userId, undefined, config);

  // Auto-subscribe to all channels for monitoring
  useEffect(() => {
    if (hook.connected) {
      hook.subscribe([
        'step_updates',
        'total_updates',
        'leaderboard_updates',
        'badge_earned',
      ]);
    }
  }, [hook.connected, hook.subscribe]);

  return hook;
}

export default useStepsWebSocket;