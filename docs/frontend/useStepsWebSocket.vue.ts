/**
 * useStepsWebSocket - Vue 3 Composable voor WebSocket stappen tracking
 * 
 * Features:
 * - Automatic connection management
 * - Reactive state with Vue refs
 * - Cleanup on unmount
 * - Type-safe
 * 
 * Usage:
 * ```vue
 * <script setup lang="ts">
 * import { useStepsWebSocket } from './composables/useStepsWebSocket';
 * 
 * const props = defineProps<{
 *   userId: string;
 *   participantId: string;
 * }>();
 * 
 * const {
 *   connected,
 *   latestUpdate,
 *   totalSteps,
 *   leaderboard
 * } = useStepsWebSocket(props.userId, props.participantId);
 * </script>
 * 
 * <template>
 *   <div>
 *     <p>Connected: {{ connected }}</p>
 *     <p>Steps: {{ latestUpdate?.steps || 0 }}</p>
 *     <p>Total: {{ totalSteps }}</p>
 *   </div>
 * </template>
 * ```
 */

import { ref, onMounted, onUnmounted, Ref } from 'vue';
import {
  StepsWebSocketClient,
  StepUpdateMessage,
  TotalUpdateMessage,
  LeaderboardUpdateMessage,
  BadgeEarnedMessage,
  ConnectionState,
  StepsWebSocketConfig,
} from './steps-websocket-client';

export interface UseStepsWebSocketReturn {
  connected: Ref<boolean>;
  connectionState: Ref<ConnectionState>;
  latestUpdate: Ref<StepUpdateMessage | null>;
  totalSteps: Ref<number>;
  leaderboard: Ref<LeaderboardUpdateMessage | null>;
  latestBadge: Ref<BadgeEarnedMessage | null>;
  subscribe: (channels: string[]) => void;
  unsubscribe: (channels: string[]) => void;
  reconnect: () => void;
  disconnect: () => void;
}

/**
 * Composable voor WebSocket connectie naar steps tracking
 */
export function useStepsWebSocket(
  userId: string | Ref<string>,
  participantId?: string | Ref<string>,
  config?: StepsWebSocketConfig
): UseStepsWebSocketReturn {
  // Reactive state
  const connected = ref<boolean>(false);
  const connectionState = ref<ConnectionState>(ConnectionState.DISCONNECTED);
  const latestUpdate = ref<StepUpdateMessage | null>(null);
  const totalSteps = ref<number>(0);
  const leaderboard = ref<LeaderboardUpdateMessage | null>(null);
  const latestBadge = ref<BadgeEarnedMessage | null>(null);

  let client: StepsWebSocketClient | null = null;

  // Get WebSocket URL
  const getWebSocketUrl = (): string => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = window.location.host;
    return `${protocol}//${host}/ws/steps`;
  };

  // Get token from localStorage
  const getToken = (): string => {
    return localStorage.getItem('token') || '';
  };

  // Initialize client
  const initializeClient = () => {
    const url = getWebSocketUrl();
    const token = getToken();
    const uid = typeof userId === 'string' ? userId : userId.value;
    const pid = participantId 
      ? (typeof participantId === 'string' ? participantId : participantId.value)
      : undefined;

    if (!uid) {
      console.warn('[useStepsWebSocket] No userId provided');
      return;
    }

    client = new StepsWebSocketClient(
      url,
      token,
      uid,
      pid,
      {
        debug: import.meta.env.DEV,
        ...config,
      }
    );

    // Setup event listeners
    client.onStateChange((state) => {
      connectionState.value = state;
      connected.value = state === ConnectionState.CONNECTED;
    });

    client.on<StepUpdateMessage>('step_update', (data) => {
      latestUpdate.value = data;
    });

    client.on<TotalUpdateMessage>('total_update', (data) => {
      totalSteps.value = data.total_steps;
    });

    client.on<LeaderboardUpdateMessage>('leaderboard_update', (data) => {
      leaderboard.value = data;
    });

    client.on<BadgeEarnedMessage>('badge_earned', (data) => {
      latestBadge.value = data;
    });

    // Connect
    client.connect();
  };

  // Subscribe to channels
  const subscribe = (channels: string[]) => {
    client?.subscribe(channels);
  };

  // Unsubscribe from channels
  const unsubscribe = (channels: string[]) => {
    client?.unsubscribe(channels);
  };

  // Manual reconnect
  const reconnect = () => {
    client?.disconnect();
    client?.connect();
  };

  // Manual disconnect
  const disconnect = () => {
    client?.disconnect();
  };

  // Lifecycle hooks
  onMounted(() => {
    initializeClient();
  });

  onUnmounted(() => {
    client?.disconnect();
    client = null;
  });

  return {
    connected,
    connectionState,
    latestUpdate,
    totalSteps,
    leaderboard,
    latestBadge,
    subscribe,
    unsubscribe,
    reconnect,
    disconnect,
  };
}

/**
 * Composable voor participant dashboard met auto-subscription
 */
export function useParticipantDashboard(
  userId: string | Ref<string>,
  participantId: string | Ref<string>,
  config?: StepsWebSocketConfig
): UseStepsWebSocketReturn {
  const composable = useStepsWebSocket(userId, participantId, config);

  // Auto-subscribe on connection
  onMounted(() => {
    if (composable.connected.value) {
      composable.subscribe(['step_updates', 'badge_earned']);
    }
  });

  return composable;
}

/**
 * Composable voor public leaderboard
 */
export function useLeaderboard(
  config?: StepsWebSocketConfig
): Pick<UseStepsWebSocketReturn, 'connected' | 'connectionState' | 'totalSteps' | 'leaderboard' | 'reconnect'> {
  const composable = useStepsWebSocket('public', undefined, config);

  // Auto-subscribe to public channels
  onMounted(() => {
    if (composable.connected.value) {
      composable.subscribe(['total_updates', 'leaderboard_updates']);
    }
  });

  return {
    connected: composable.connected,
    connectionState: composable.connectionState,
    totalSteps: composable.totalSteps,
    leaderboard: composable.leaderboard,
    reconnect: composable.reconnect,
  };
}

/**
 * Composable voor admin monitoring
 */
export function useStepsMonitoring(
  userId: string | Ref<string>,
  config?: StepsWebSocketConfig
): UseStepsWebSocketReturn {
  const composable = useStepsWebSocket(userId, undefined, config);

  // Auto-subscribe to all channels
  onMounted(() => {
    if (composable.connected.value) {
      composable.subscribe([
        'step_updates',
        'total_updates',
        'leaderboard_updates',
        'badge_earned',
      ]);
    }
  });

  return composable;
}

export default useStepsWebSocket;