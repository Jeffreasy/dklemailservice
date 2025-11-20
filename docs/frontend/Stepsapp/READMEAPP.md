# DKL Steps App Documentation

Complete documentation for the DKL Steps App - a React Native mobile application for real-time step tracking and gamification during DKL events.

## Table of Contents

- [Overview](#overview)
- [Target Users](#target-users)
- [Key Features](#key-features)
- [Technology Stack](#technology-stack)
- [Architecture](#architecture)
- [API Integration](#api-integration)
- [Security](#security)
- [User Interface](#user-interface)
- [Development Setup](#development-setup)
- [Deployment](#deployment)
- [Testing](#testing)

## Overview

The DKL Steps App is a cross-platform mobile application built with React Native that enables participants to track their steps in real-time during DKL walking events. The app provides gamification features, social interaction, and live event updates through WebSocket connections.

### Purpose

The Steps App transforms walking events into engaging experiences by:
- Real-time step tracking with GPS integration
- Gamification with achievements and leaderboards
- Social features for community interaction
- Live event updates and notifications
- Offline functionality for rural areas

## Target Users

### Primary Users
- **Event Participants**: Full account holders with app access
- **Casual Walkers**: Temporary account users (limited features)
- **Supporters**: Family and friends following participants

### Account Types
- **Full Account Users**: Complete app access with all features
- **Temporary Account Users**: Limited access for event-only participation
- **Guest Viewers**: Read-only access for following participants

## Key Features

### Step Tracking & GPS
- **Real-time GPS Tracking**: Continuous location monitoring during events
- **Geofence Detection**: Automatic start/finish detection at event boundaries
- **Step Counting**: Pedometer integration with health data
- **Route Validation**: GPS-based route verification and distance calculation

### Gamification System
- **Achievements**: Unlock badges for milestones and challenges
- **Leaderboards**: Real-time ranking system with social features
- **Challenges**: Daily and event-specific walking challenges
- **Points System**: Earn points for steps, achievements, and social engagement

### Social Features
- **Participant Profiles**: Public profiles with stats and achievements
- **Team Support**: Cheer for friends and family members
- **Live Updates**: Real-time progress sharing
- **Community Feed**: Event-wide activity and celebrations

### Event Integration
- **Live Event Data**: Real-time event status and participant counts
- **Checkpoint Tracking**: Progress through event checkpoints
- **Emergency Features**: SOS functionality and emergency contacts
- **Event Notifications**: Push notifications for important updates

### Offline Functionality
- **Offline Tracking**: Continue tracking without internet connection
- **Data Synchronization**: Automatic sync when connection restored
- **Cached Content**: Access to previously loaded data offline
- **Background Processing**: Continue tracking in background

## Technology Stack

### Mobile Framework
- **React Native 0.72+**: Cross-platform mobile development
- **Expo**: Managed workflow for easier development and deployment
- **TypeScript**: Type-safe development with comprehensive type definitions

### Navigation & State Management
- **React Navigation**: Native navigation with stack, tab, and drawer navigators
- **Zustand**: Lightweight state management for global app state
- **AsyncStorage**: Persistent local storage for user preferences
- **React Query**: Server state management and offline caching

### Device Integration
- **React Native Location**: GPS and geofencing capabilities
- **Expo Location**: High-accuracy location services
- **Expo Pedometer**: Step counting and health data integration
- **Expo Notifications**: Push notifications and local alerts

### Networking & Real-time
- **Axios**: HTTP client with interceptors and offline support
- **Socket.IO Client**: WebSocket connections for real-time updates
- **NetInfo**: Network connectivity monitoring
- **Background Tasks**: Background location and data sync

### UI Components & Styling
- **React Native Paper**: Material Design components for iOS/Android
- **Expo Vector Icons**: Comprehensive icon library
- **React Native Reanimated**: Smooth animations and transitions
- **React Native Gesture Handler**: Native touch gestures

### Development Tools
- **Expo CLI**: Development server and build tools
- **ESLint + Prettier**: Code quality and formatting
- **Jest + React Native Testing Library**: Unit and integration testing
- **Detox**: End-to-end testing for mobile apps

## Architecture

### App Structure

```
mobile/
├── src/
│   ├── components/          # Reusable UI components
│   │   ├── common/         # Generic components
│   │   ├── tracking/       # GPS and step tracking
│   │   ├── gamification/   # Achievements and leaderboards
│   │   └── social/         # Social features
│   ├── screens/            # Screen components
│   │   ├── auth/           # Authentication screens
│   │   ├── tracking/       # Step tracking screens
│   │   ├── social/         # Social interaction screens
│   │   └── settings/       # App settings
│   ├── navigation/         # Navigation configuration
│   ├── services/           # API and external services
│   ├── stores/             # Zustand state stores
│   ├── hooks/              # Custom React hooks
│   ├── utils/              # Utility functions
│   ├── types/              # TypeScript definitions
│   └── constants/          # App constants
├── assets/                 # Images, fonts, and media
├── e2e/                    # End-to-end tests
└── __tests__/              # Unit tests
```

### State Management Strategy

#### Global State (Zustand)
```typescript
// Auth store
interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => void;
}

// Tracking store
interface TrackingState {
  isTracking: boolean;
  currentSteps: number;
  currentLocation: Location | null;
  startTracking: () => void;
  stopTracking: () => void;
}
```

#### Server State (React Query)
```typescript
// API data with offline support
const { data: leaderboard, isLoading } = useQuery({
  queryKey: ['leaderboard'],
  queryFn: () => api.getLeaderboard(),
  staleTime: 5 * 60 * 1000, // 5 minutes
  cacheTime: 30 * 60 * 1000, // 30 minutes
  networkMode: 'offlineFirst', // Enable offline support
});
```

#### Local State (React Hooks)
```typescript
function StepTracker() {
  const [steps, setSteps] = useState(0);
  const [isActive, setIsActive] = useState(false);

  // Component-specific state
}
```

## API Integration

### Authentication Flow

```typescript
// Secure token storage for mobile
import * as SecureStore from 'expo-secure-store';

class TokenManager {
  static async setTokens(accessToken: string, refreshToken: string) {
    await SecureStore.setItemAsync('access_token', accessToken);
    await SecureStore.setItemAsync('refresh_token', refreshToken);
  }

  static async getAccessToken(): Promise<string | null> {
    return await SecureStore.getItemAsync('access_token');
  }

  static async refreshAccessToken(): Promise<string> {
    const refreshToken = await SecureStore.getItemAsync('refresh_token');
    if (!refreshToken) return null;

    const response = await api.post('/auth/refresh', {
      refresh_token: refreshToken
    });

    const { token, refresh_token: newRefreshToken } = response.data;
    await this.setTokens(token, newRefreshToken);
    return token;
  }
}
```

### Real-time WebSocket Integration

```typescript
// WebSocket client for real-time updates
import io from 'socket.io-client';

class WebSocketManager {
  private socket: Socket | null = null;

  connect(token: string) {
    this.socket = io('wss://api.dklemailservice.com', {
      auth: { token },
      transports: ['websocket'],
    });

    this.socket.on('connect', () => {
      console.log('Connected to WebSocket');
    });

    this.socket.on('leaderboard:update', (data) => {
      // Update leaderboard in real-time
      queryClient.invalidateQueries(['leaderboard']);
    });

    this.socket.on('achievement:unlocked', (achievement) => {
      // Show achievement notification
      showAchievementToast(achievement);
    });
  }

  disconnect() {
    this.socket?.disconnect();
  }
}
```

### GPS Tracking Service

```typescript
// Background location tracking
import * as Location from 'expo-location';
import * as TaskManager from 'expo-task-manager';

const LOCATION_TASK_NAME = 'background-location-task';

TaskManager.defineTask(LOCATION_TASK_NAME, ({ data, error }) => {
  if (error) {
    console.error('Location task error:', error);
    return;
  }

  const { locations } = data as { locations: Location.LocationObject[] };
  const location = locations[0];

  if (location) {
    // Send location update to API
    api.post('/registration/location', {
      lat: location.coords.latitude,
      long: location.coords.longitude,
      accuracy: location.coords.accuracy,
      timestamp: location.timestamp,
    }).catch(error => {
      // Store for later sync if offline
      storeOfflineLocation(location);
    });
  }
});

class LocationService {
  static async startTracking() {
    const { status } = await Location.requestBackgroundPermissionsAsync();
    if (status !== 'granted') {
      throw new Error('Location permission denied');
    }

    await Location.startLocationUpdatesAsync(LOCATION_TASK_NAME, {
      accuracy: Location.Accuracy.High,
      timeInterval: 10000, // 10 seconds
      distanceInterval: 10, // 10 meters
      showsBackgroundLocationIndicator: true,
    });
  }

  static async stopTracking() {
    await Location.stopLocationUpdatesAsync(LOCATION_TASK_NAME);
  }
}
```

### Offline Data Synchronization

```typescript
// Offline queue for API requests
class OfflineQueue {
  private queue: QueuedRequest[] = [];

  async addRequest(request: QueuedRequest) {
    this.queue.push(request);
    await AsyncStorage.setItem('offline_queue', JSON.stringify(this.queue));
  }

  async processQueue() {
    const netInfo = await NetInfo.fetch();
    if (netInfo.isConnected) {
      for (const request of this.queue) {
        try {
          await api.request(request.config);
          // Remove from queue on success
          this.queue = this.queue.filter(r => r.id !== request.id);
        } catch (error) {
          // Keep in queue for retry
          console.error('Failed to sync request:', error);
        }
      }
      await AsyncStorage.setItem('offline_queue', JSON.stringify(this.queue));
    }
  }
}
```

## Security

### Mobile-Specific Security

#### iOS Keychain Integration
```typescript
// iOS secure storage
import * as Keychain from 'react-native-keychain';

class SecureTokenManager {
  static async storeTokens(accessToken: string, refreshToken: string) {
    await Keychain.setGenericPassword(
      'accessToken',
      accessToken,
      {
        service: 'com.dkl.stepsapp',
        accessControl: Keychain.ACCESS_CONTROL.BIOMETRY_ANY,
        accessible: Keychain.ACCESSIBLE.WHEN_UNLOCKED
      }
    );

    await Keychain.setGenericPassword(
      'refreshToken',
      refreshToken,
      {
        service: 'com.dkl.stepsapp.refresh',
        accessible: Keychain.ACCESSIBLE.AFTER_FIRST_UNLOCK
      }
    );
  }
}
```

#### Android Keystore Integration
```typescript
// Android secure storage
import EncryptedStorage from 'react-native-encrypted-storage';

class SecureTokenManager {
  static async storeTokens(accessToken: string, refreshToken: string) {
    await EncryptedStorage.setItem(
      'access_token',
      JSON.stringify({
        token: accessToken,
        timestamp: Date.now()
      })
    );

    await EncryptedStorage.setItem(
      'refresh_token',
      JSON.stringify({
        token: refreshToken,
        timestamp: Date.now()
      })
    );
  }
}
```

### Biometric Authentication

```typescript
// Face ID / Touch ID integration
import * as LocalAuthentication from 'expo-local-authentication';

class BiometricAuth {
  static async authenticate(): Promise<boolean> {
    const hasHardware = await LocalAuthentication.hasHardwareAsync();
    const supportedTypes = await LocalAuthentication.supportedAuthenticationTypesAsync();

    if (!hasHardware || supportedTypes.length === 0) {
      return false;
    }

    const result = await LocalAuthentication.authenticateAsync({
      promptMessage: 'Authenticate to access DKL Steps App',
      fallbackLabel: 'Use passcode',
    });

    return result.success;
  }
}
```

### Certificate Pinning

```typescript
// SSL certificate pinning for API calls
import { Platform } from 'react-native';
import axios from 'axios';

if (Platform.OS === 'ios') {
  // iOS certificate pinning
  const config = {
    // Certificate pinning configuration
  };
} else {
  // Android certificate pinning
  const config = {
    // Certificate pinning configuration
  };
}

const api = axios.create({
  baseURL: 'https://api.dklemailservice.com',
  ...config,
});
```

## User Interface

### Design System

#### Color Palette
```typescript
const colors = {
  primary: '#1976d2',
  secondary: '#dc004e',
  success: '#4caf50',
  warning: '#ff9800',
  error: '#f44336',
  background: '#f5f5f5',
  surface: '#ffffff',
  text: '#212121',
  textSecondary: '#757575',
};
```

#### Typography Scale
- **H1**: 32px - Screen titles
- **H2**: 24px - Section headers
- **H3**: 20px - Card titles
- **Body**: 16px - Regular text
- **Caption**: 14px - Secondary text
- **Small**: 12px - Metadata

### Navigation Structure

#### Tab Navigation
- **Home**: Dashboard with current progress and quick actions
- **Track**: Step tracking interface with live GPS data
- **Leaderboard**: Real-time rankings and social features
- **Profile**: User profile, achievements, and settings

#### Stack Navigation
- **Authentication**: Login, registration, password reset
- **Event Details**: Event information and participant list
- **Achievement Details**: Achievement information and sharing
- **Settings**: App preferences and account management

### Key UI Components

#### Step Tracking Interface
```typescript
function StepTracker() {
  const { steps, distance, isTracking } = useTracking();

  return (
    <View style={styles.container}>
      <AnimatedCircularProgress
        size={200}
        width={10}
        fill={(steps / targetSteps) * 100}
        tintColor={colors.primary}
        backgroundColor={colors.background}
      >
        {() => (
          <View style={styles.stepDisplay}>
            <Text style={styles.stepCount}>{steps.toLocaleString()}</Text>
            <Text style={styles.stepLabel}>Steps</Text>
          </View>
        )}
      </AnimatedCircularProgress>

      <View style={styles.stats}>
        <StatCard title="Distance" value={`${distance.toFixed(1)} km`} />
        <StatCard title="Time" value={formatDuration(duration)} />
        <StatCard title="Pace" value={`${pace} min/km`} />
      </View>
    </View>
  );
}
```

#### Leaderboard Component
```typescript
function Leaderboard() {
  const { data: leaderboard } = useLeaderboard();

  return (
    <FlatList
      data={leaderboard}
      keyExtractor={(item) => item.id}
      renderItem={({ item, index }) => (
        <LeaderboardItem
          rank={index + 1}
          participant={item}
          isCurrentUser={item.id === currentUser.id}
        />
      )}
    />
  );
}
```

#### Achievement System
```typescript
function AchievementBadge({ achievement }: { achievement: Achievement }) {
  return (
    <TouchableOpacity style={styles.badge}>
      <Image source={{ uri: achievement.iconUrl }} style={styles.icon} />
      <Text style={styles.title}>{achievement.title}</Text>
      <Text style={styles.description}>{achievement.description}</Text>
      {achievement.unlockedAt && (
        <Text style={styles.unlockedText}>Unlocked!</Text>
      )}
    </TouchableOpacity>
  );
}
```

## Development Setup

### Prerequisites
- **Node.js**: 18+ LTS
- **Expo CLI**: `npm install -g @expo/cli`
- **iOS Simulator**: Xcode (macOS only)
- **Android Emulator**: Android Studio
- **Physical Device**: For testing push notifications

### Installation

```bash
# Clone the repository
git clone https://github.com/dkl/dkl-steps-app.git
cd dkl-steps-app

# Install dependencies
npm install

# Install iOS dependencies (macOS only)
cd ios && pod install && cd ..

# Start Expo development server
npx expo start
```

### Environment Configuration

```javascript
// app.config.js
export default {
  name: 'DKL Steps App',
  slug: 'dkl-steps-app',
  version: '1.0.0',
  orientation: 'portrait',
  icon: './assets/icon.png',
  userInterfaceStyle: 'light',
  splash: {
    image: './assets/splash.png',
    resizeMode: 'contain',
    backgroundColor: '#ffffff'
  },
  assetBundlePatterns: ['**/*'],
  ios: {
    supportsTablet: false,
    bundleIdentifier: 'com.dkl.stepsapp'
  },
  android: {
    adaptiveIcon: {
      foregroundImage: './assets/adaptive-icon.png',
      backgroundColor: '#ffffff'
    },
    package: 'com.dkl.stepsapp'
  },
  plugins: [
    'expo-location',
    'expo-notifications',
    [
      'expo-build-properties',
      {
        ios: {
          deploymentTarget: '13.0'
        },
        android: {
          compileSdkVersion: 33,
          targetSdkVersion: 33,
          buildToolsVersion: '33.0.0'
        }
      }
    ]
  ],
  extra: {
    apiUrl: process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080/api',
    wsUrl: process.env.EXPO_PUBLIC_WS_URL || 'ws://localhost:8080'
  }
};
```

### Development Scripts

```json
{
  "scripts": {
    "start": "expo start",
    "android": "expo start --android",
    "ios": "expo start --ios",
    "web": "expo start --web",
    "eject": "expo eject",
    "test": "jest",
    "lint": "eslint . --ext .js,.jsx,.ts,.tsx",
    "format": "prettier --write ."
  }
}
```

### Platform-Specific Setup

#### iOS Development
```bash
# Install CocoaPods
sudo gem install cocoapods

# Install iOS dependencies
cd ios && pod install && cd ..

# Run on iOS simulator
npm run ios
```

#### Android Development
```bash
# Start Android emulator
# Then run:
npm run android
```

## Deployment

### Build Process

#### Expo Application Services (EAS)
```bash
# Install EAS CLI
npm install -g @expo/eas-cli

# Login to Expo
eas login

# Configure build profiles
# eas.json
{
  "build": {
    "development": {
      "developmentClient": true,
      "distribution": "internal"
    },
    "preview": {
      "distribution": "internal"
    },
    "production": {
      "channel": "production"
    }
  },
  "submit": {
    "production": {}
  }
}

# Build for production
eas build --platform ios
eas build --platform android
```

### App Store Deployment

#### iOS App Store
```bash
# Build for App Store
eas build --platform ios --profile production

# Submit to App Store
eas submit --platform ios
```

#### Google Play Store
```bash
# Build for Play Store
eas build --platform android --profile production

# Submit to Play Store
eas submit --platform android
```

### Over-the-Air Updates

```typescript
// Configure OTA updates
import * as Updates from 'expo-updates';

async function checkForUpdates() {
  try {
    const update = await Updates.checkForUpdateAsync();
    if (update.isAvailable) {
      await Updates.fetchUpdateAsync();
      await Updates.reloadAsync();
    }
  } catch (error) {
    console.error('Error checking for updates:', error);
  }
}
```

## Testing

### Testing Strategy

#### Unit Tests
```typescript
// Component testing
import { render, fireEvent } from '@testing-library/react-native';
import { StepCounter } from '../StepCounter';

test('displays step count correctly', () => {
  const { getByText } = render(<StepCounter steps={1500} />);
  expect(getByText('1,500')).toBeTruthy();
});

test('calls onIncrement when plus button is pressed', () => {
  const mockOnIncrement = jest.fn();
  const { getByTestId } = render(
    <StepCounter steps={0} onIncrement={mockOnIncrement} />
  );

  fireEvent.press(getByTestId('increment-button'));
  expect(mockOnIncrement).toHaveBeenCalled();
});
```

#### Integration Tests
```typescript
// API integration testing
import { renderHook, waitFor } from '@testing-library/react-native';
import { useLeaderboard } from '../hooks/useLeaderboard';

test('fetches leaderboard data', async () => {
  const { result } = renderHook(() => useLeaderboard());

  await waitFor(() => {
    expect(result.current.isSuccess).toBe(true);
  });

  expect(result.current.data).toBeDefined();
  expect(Array.isArray(result.current.data)).toBe(true);
});
```

#### End-to-End Tests (Detox)
```javascript
// e2e/firstTest.e2e.js
describe('Login Flow', () => {
  beforeEach(async () => {
    await device.reloadReactNative();
  });

  it('should login successfully', async () => {
    await element(by.id('email-input')).typeText('user@example.com');
    await element(by.id('password-input')).typeText('password123');
    await element(by.id('login-button')).tap();

    await expect(element(by.id('dashboard-screen'))).toBeVisible();
  });
});
```

### Test Configuration

#### Jest Configuration
```javascript
// jest.config.js
module.exports = {
  preset: 'jest-expo',
  setupFilesAfterEnv: ['<rootDir>/jest.setup.js'],
  testPathIgnorePatterns: ['/node_modules/', '/e2e/'],
  collectCoverageFrom: [
    'src/**/*.{ts,tsx}',
    '!src/**/*.d.ts',
  ],
  coverageThreshold: {
    global: {
      branches: 70,
      functions: 80,
      lines: 80,
      statements: 80,
    },
  },
};
```

#### Detox Configuration
```javascript
// .detoxrc.js
module.exports = {
  testRunner: 'jest',
  runnerConfig: 'e2e/config.json',
  configurations: {
    ios: {
      type: 'ios.simulator',
      device: {
        type: 'iPhone 13',
      },
    },
    android: {
      type: 'android.emulator',
      device: {
        avdName: 'Pixel_3_API_29',
      },
    },
  },
};
```

### CI/CD Pipeline

```yaml
# GitHub Actions workflow
name: Mobile CI/CD
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - run: npm ci
      - run: npm run lint
      - run: npm run test -- --coverage
      - run: npm run build

  build-ios:
    needs: test
    runs-on: macos-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v3
      - uses: expo/expo-github-action@v8
        with:
          expo-version: latest
          token: ${{ secrets.EXPO_TOKEN }}
      - run: eas build --platform ios --profile preview --non-interactive

  build-android:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v3
      - uses: expo/expo-github-action@v8
        with:
          expo-version: latest
          token: ${{ secrets.EXPO_TOKEN }}
      - run: eas build --platform android --profile preview --non-interactive
```

---

This documentation provides a comprehensive overview of the DKL Steps App architecture and development practices. For specific implementation details, refer to the inline code comments and API documentation.