# DKL Admin Panel Documentation

Complete documentation for the DKL Admin Panel - a comprehensive web-based administration interface for managing the DKL Email Service system.

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

The DKL Admin Panel is a modern web-based administration interface built with React that provides comprehensive management capabilities for the DKL Email Service system. It serves as the central control center for administrators and staff to manage all aspects of the event management platform.

### Purpose

The admin panel enables efficient management of:
- User accounts and permissions
- Event configuration and monitoring
- Content management (photos, videos, partners, sponsors)
- Participant registrations and tracking
- Email communications and templates
- System monitoring and analytics

## Target Users

### Primary Users
- **System Administrators**: Full access to all system features
- **Event Staff**: Event management, participant tracking, content management
- **Content Managers**: CMS functionality, media management, newsletter management

### Access Levels
- **Admin**: Complete system access with all permissions
- **Staff**: Limited administrative access based on assigned roles
- **Moderator**: Content moderation and basic management functions

## Key Features

### User Management
- **Account Management**: Create, edit, and deactivate user accounts
- **Role Assignment**: Assign roles and permissions via RBAC system
- **Password Management**: Reset passwords and manage account security
- **Session Monitoring**: View active sessions and device management

### Event Management
- **Event Creation**: Configure events with geofencing and routing
- **Registration Management**: Monitor and manage participant registrations
- **Real-time Tracking**: Live participant tracking during events
- **Statistics & Reporting**: Event analytics and performance metrics

### Content Management System (CMS)
- **Media Library**: Upload and manage photos and videos
- **Album Management**: Organize media into themed collections
- **Partner & Sponsor Management**: Manage business relationships
- **Newsletter System**: Create and send email newsletters

### Communication Tools
- **Email Templates**: Manage and customize email templates
- **Contact Management**: Handle user inquiries and support tickets
- **Bulk Communications**: Send targeted emails to user groups
- **Auto-responses**: Configure automated email responses

### Analytics & Monitoring
- **System Health**: Monitor API performance and database status
- **User Activity**: Track user engagement and system usage
- **Event Analytics**: Real-time event participation statistics
- **Security Logs**: Audit trail of administrative actions

## Technology Stack

### Frontend Framework
- **React 18+**: Modern React with hooks and functional components
- **TypeScript**: Type-safe development with comprehensive type definitions
- **Vite**: Fast build tool and development server with SWC

### UI Components & Styling
- **Headless UI**: Unstyled, accessible UI components (@headlessui/react)
- **Heroicons**: Beautiful hand-crafted SVG icons (@heroicons/react)
- **Tailwind CSS**: Utility-first CSS framework with forms plugin
- **Tailwind Variants**: Component variant utilities
- **Tabler Icons**: Additional icon library (@tabler/icons-react)

### Data Fetching & State Management
- **TanStack Query**: Powerful data fetching and caching (@tanstack/react-query)
- **React Hook Form**: Performant forms with validation (@hookform/resolvers)
- **Zod**: TypeScript-first schema validation
- **Axios**: HTTP client with interceptors for authentication

### Rich Text Editing
- **Tiptap**: Extensible rich text editor framework
- **Tiptap Extensions**: Blockquote, color, heading, highlight, image, link, text alignment
- **Mantine Tiptap**: Pre-built Tiptap components

### Drag & Drop
- **@dnd-kit**: Modern drag and drop toolkit for React
- **@hello-pangea/dnd**: Alternative drag and drop library

### Routing & Navigation
- **React Router DOM**: Client-side routing and navigation
- **React Window**: Virtualized lists and tables for performance

### Utilities & Libraries
- **Date-fns**: Modern JavaScript date utility library
- **Lodash**: Utility functions for common programming tasks
- **JWT Decode**: Decode JWT tokens
- **DOMPurify**: XSS sanitization
- **EmailJS**: Client-side email sending

### Development Tools
- **ESLint**: Code linting with TypeScript and React rules
- **Vitest**: Fast unit testing framework with UI
- **Playwright**: End-to-end testing framework
- **Testing Library**: React component testing utilities
- **MSW**: API mocking for testing

## Architecture

### Component Architecture

```
src/
├── components/          # Reusable UI components
│   ├── ui/             # Headless UI components (buttons, inputs, modals)
│   ├── forms/          # Form components with react-hook-form
│   ├── layout/         # Layout components (header, sidebar, navigation)
│   ├── tables/         # Data table components with virtualization
│   ├── editors/        # Rich text editors (Tiptap-based)
│   └── drag-drop/      # Drag and drop components (@dnd-kit)
├── pages/              # Page components (dashboard, users, events, etc.)
├── hooks/              # Custom React hooks
│   ├── useAuth.ts      # Authentication hooks
│   ├── useForm.ts      # Form validation hooks
│   └── useQuery.ts     # Data fetching hooks
├── services/           # API service layer
│   ├── api.ts          # Axios client configuration
│   ├── auth.ts         # Authentication services
│   └── endpoints/      # API endpoint definitions
├── lib/                # Utility libraries
│   ├── validations/    # Zod schemas
│   ├── utils/          # Helper functions
│   └── constants/      # Application constants
├── types/              # TypeScript type definitions
└── test/               # Test utilities and mocks
```

### State Management Strategy

#### Server State (TanStack Query)
- **Data fetching**: Centralized API calls with caching
- **Background refetching**: Automatic data updates
- **Optimistic updates**: Immediate UI feedback
- **Error handling**: Retry logic and error boundaries
- **Cache management**: Intelligent cache invalidation

```typescript
// Query configuration
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 minutes
      gcTime: 10 * 60 * 1000, // 10 minutes
      retry: (failureCount, error) => {
        // Custom retry logic
        if (error?.status === 401) return false;
        return failureCount < 3;
      },
    },
  },
});
```

#### Form State (React Hook Form + Zod)
- **Validation**: Client-side validation with Zod schemas
- **Performance**: Minimal re-renders with controlled components
- **Integration**: Seamless API integration with TanStack Query

```typescript
const schema = z.object({
  name: z.string().min(2, 'Name is required'),
  email: z.string().email('Invalid email'),
});

function UserForm() {
  const { register, handleSubmit, formState: { errors } } = useForm({
    resolver: zodResolver(schema),
  });

  const mutation = useMutation({
    mutationFn: createUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
    },
  });
}
```

#### Component State (React Hooks)
- **Local state**: useState for component-specific data
- **Effects**: useEffect for side effects and lifecycle
- **Memoization**: useMemo and useCallback for performance

## API Integration

### Authentication Flow

```typescript
// Login process
const login = async (credentials: LoginCredentials) => {
  try {
    const response = await authService.login(credentials);
    // Store tokens securely
    tokenManager.setTokens(response.token, response.refresh_token);
    // Update global auth state
    authStore.setUser(response.user);
    // Redirect to dashboard
    navigate('/dashboard');
  } catch (error) {
    // Handle authentication errors
    showErrorToast('Invalid credentials');
  }
};
```

### API Service Layer

```typescript
// Service class pattern
class EventService {
  private api = axios.create({
    baseURL: '/api',
    timeout: 10000,
  });

  constructor() {
    // Add JWT interceptor
    this.api.interceptors.request.use((config) => {
      const token = tokenManager.getAccessToken();
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    // Add response interceptor for token refresh
    this.api.interceptors.response.use(
      (response) => response,
      async (error) => {
        if (error.response?.status === 401) {
          // Attempt token refresh
          const newToken = await tokenManager.refreshToken();
          if (newToken) {
            // Retry original request
            return this.api.request(error.config);
          }
        }
        return Promise.reject(error);
      }
    );
  }

  async getEvents(params?: EventQueryParams): Promise<Event[]> {
    const response = await this.api.get('/events', { params });
    return response.data;
  }

  async createEvent(eventData: EventCreateData): Promise<Event> {
    const response = await this.api.post('/events', eventData);
    return response.data;
  }
}
```

### Real-time Updates

```typescript
// WebSocket integration for live updates
import { useWebSocket } from '@/hooks/useWebSocket';
import { useQueryClient } from '@tanstack/react-query';

function EventDashboard() {
  const queryClient = useQueryClient();
  const { isConnected, subscribe, unsubscribe } = useWebSocket();

  useEffect(() => {
    // Subscribe to event updates
    subscribe('event:updates', (data) => {
      // Update local state with real-time data
      queryClient.invalidateQueries({ queryKey: ['events'] });
    });

    return () => {
      unsubscribe('event:updates');
    };
  }, [subscribe, unsubscribe, queryClient]);

  // Component renders with real-time data
}
```

### Form Handling with Validation

```typescript
// Form handling with react-hook-form and zod
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useMutation, useQueryClient } from '@tanstack/react-query';

const userSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.string().email('Invalid email address'),
  role: z.enum(['admin', 'staff', 'moderator']),
});

type UserFormData = z.infer<typeof userSchema>;

function CreateUserForm() {
  const queryClient = useQueryClient();
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<UserFormData>({
    resolver: zodResolver(userSchema),
  });

  const createUserMutation = useMutation({
    mutationFn: (data: UserFormData) => api.post('/users', data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      reset();
      toast.success('User created successfully');
    },
    onError: (error) => {
      toast.error('Failed to create user');
    },
  });

  const onSubmit = (data: UserFormData) => {
    createUserMutation.mutate(data);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div>
        <label htmlFor="name">Name</label>
        <input
          {...register('name')}
          id="name"
          className="w-full px-3 py-2 border rounded-md"
        />
        {errors.name && (
          <p className="text-red-500 text-sm">{errors.name.message}</p>
        )}
      </div>

      <div>
        <label htmlFor="email">Email</label>
        <input
          {...register('email')}
          id="email"
          type="email"
          className="w-full px-3 py-2 border rounded-md"
        />
        {errors.email && (
          <p className="text-red-500 text-sm">{errors.email.message}</p>
        )}
      </div>

      <button
        type="submit"
        disabled={isSubmitting}
        className="px-4 py-2 bg-blue-500 text-white rounded-md disabled:opacity-50"
      >
        {isSubmitting ? 'Creating...' : 'Create User'}
      </button>
    </form>
  );
}
```

## Security

### Token Storage Security

#### Web Browser Environment
- **Access Token**: Stored in memory only (never in localStorage)
- **Refresh Token**: HttpOnly, secure, SameSite cookies (preferred)
- **Fallback**: Encrypted localStorage with automatic cleanup

#### Implementation
```typescript
class TokenManager {
  // Access token - memory only
  private accessToken: string | null = null;

  // Refresh token - httpOnly cookie
  async refreshAccessToken(): Promise<string> {
    const response = await fetch('/api/auth/refresh', {
      method: 'POST',
      credentials: 'include' // Include httpOnly cookies
    });

    const data = await response.json();
    this.accessToken = data.token;
    return data.token;
  }

  getAccessToken(): string | null {
    return this.accessToken;
  }
}
```

### Permission-Based UI Rendering

```typescript
// Permission-based component rendering
function AdminPanel({ children, requiredPermission }: AdminPanelProps) {
  const { user } = useAuth();
  const hasPermission = usePermission(requiredPermission);

  if (!hasPermission) {
    return <AccessDenied />;
  }

  return <div className="admin-panel">{children}</div>;
}

// Usage
<AdminPanel requiredPermission="admin:access">
  <UserManagement />
</AdminPanel>
```

### Security Best Practices

- **Input Validation**: All user inputs validated on client and server
- **XSS Protection**: Sanitized HTML rendering, CSP headers
- **CSRF Protection**: Token-based CSRF protection
- **Rate Limiting**: API rate limiting with progressive delays
- **Audit Logging**: All administrative actions logged
- **Session Management**: Automatic logout on inactivity

## User Interface

### Design System

#### Color Palette
```css
:root {
  --primary-color: #1976d2;
  --secondary-color: #dc004e;
  --success-color: #4caf50;
  --warning-color: #ff9800;
  --error-color: #f44336;
  --background-color: #f5f5f5;
  --surface-color: #ffffff;
}
```

#### Typography Scale
- **H1**: 2.5rem (40px) - Page titles
- **H2**: 2rem (32px) - Section headers
- **H3**: 1.5rem (24px) - Subsection headers
- **Body**: 1rem (16px) - Regular text
- **Caption**: 0.875rem (14px) - Secondary text

### Responsive Design

#### Breakpoints
- **Mobile**: < 768px
- **Tablet**: 768px - 1024px
- **Desktop**: > 1024px

#### Layout Components
- **Sidebar Navigation**: Collapsible navigation menu
- **Top Bar**: User menu, notifications, search
- **Main Content**: Responsive grid layout
- **Data Tables**: Sortable, filterable, paginated tables

### Key UI Components

#### Dashboard
- **Metrics Cards**: Key performance indicators with Headless UI
- **Charts**: Event participation, user activity graphs
- **Recent Activity**: Live feed of system events
- **Quick Actions**: Shortcuts to common tasks

#### Data Management
- **Virtualized Tables**: Performance-optimized tables with react-window
- **Advanced Filters**: Multi-field filtering and search with Zod validation
- **Bulk Actions**: Select multiple items for batch operations
- **Drag & Drop**: Reorder items with @dnd-kit
- **Export Functions**: CSV/PDF export capabilities

#### Forms
- **React Hook Form**: High-performance form handling
- **Zod Validation**: Type-safe form validation
- **File Upload**: React-dropzone for drag-and-drop uploads
- **Rich Text Editor**: Tiptap-based WYSIWYG editor with extensions

#### Rich Text Editor Implementation
```typescript
// Tiptap rich text editor with extensions
import { useEditor, EditorContent } from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import Image from '@tiptap/extension-image';
import Link from '@tiptap/extension-link';
import TextAlign from '@tiptap/extension-text-align';
import Highlight from '@tiptap/extension-highlight';

function RichTextEditor({ content, onChange }: RichTextEditorProps) {
  const editor = useEditor({
    extensions: [
      StarterKit,
      Image,
      Link.configure({
        openOnClick: false,
      }),
      TextAlign.configure({
        types: ['heading', 'paragraph'],
      }),
      Highlight,
    ],
    content,
    onUpdate: ({ editor }) => {
      onChange(editor.getHTML());
    },
  });

  return (
    <div className="border rounded-md">
      <EditorToolbar editor={editor} />
      <EditorContent editor={editor} className="p-4 min-h-[200px]" />
    </div>
  );
}
```

#### Drag & Drop Components
```typescript
// Drag and drop with @dnd-kit
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core';
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';

function SortableList({ items, onReorder }: SortableListProps) {
  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  function handleDragEnd(event) {
    const { active, over } = event;

    if (active.id !== over.id) {
      const oldIndex = items.findIndex((item) => item.id === active.id);
      const newIndex = items.findIndex((item) => item.id === over.id);

      onReorder(arrayMove(items, oldIndex, newIndex));
    }
  }

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragEnd={handleDragEnd}
    >
      <SortableContext items={items} strategy={verticalListSortingStrategy}>
        {items.map((item) => (
          <SortableItem key={item.id} item={item} />
        ))}
      </SortableContext>
    </DndContext>
  );
}
```

## Development Setup

### Prerequisites
- **Node.js**: 18+ LTS
- **npm** or **yarn**: Package manager
- **Git**: Version control

### Installation

```bash
# Clone the repository
git clone https://github.com/dkl/dkl25adminpanel.git
cd dkl25adminpanel

# Install dependencies
npm install

# Copy environment variables
cp .env.example .env.local

# Start development server
npm run dev
```

### Environment Configuration

```env
# API Configuration
VITE_API_BASE_URL=http://localhost:8080/api
VITE_WS_BASE_URL=ws://localhost:8080

# Authentication
VITE_JWT_REFRESH_BUFFER=300000

# Feature Flags
VITE_ENABLE_ANALYTICS=true
VITE_ENABLE_WEBSOCKET=true
```

### Development Scripts

```json
{
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "build:analyze": "tsc && vite build --mode analyze",
    "lint": "eslint .",
    "preview": "vite preview",
    "test": "vitest",
    "test:ui": "vitest --ui",
    "test:coverage": "vitest --coverage",
    "test:e2e": "playwright test",
    "test:e2e:ui": "playwright test --ui",
    "test:e2e:headed": "playwright test --headed",
    "test:e2e:debug": "playwright test --debug",
    "test:all": "npm test -- --run && npm run test:e2e"
  }
}
```

### Code Quality

#### ESLint Configuration
```javascript
module.exports = {
  extends: [
    'react-app',
    'react-app/jest',
    '@typescript-eslint/recommended',
    'prettier'
  ],
  rules: {
    // Custom rules for admin panel
    '@typescript-eslint/no-unused-vars': 'error',
    'react-hooks/exhaustive-deps': 'warn',
    'no-console': 'warn'
  }
};
```

## Deployment

### Build Process

```bash
# Production build
npm run build

# Build artifacts in dist/ directory
# - index.html
# - assets/ (JS, CSS, images)
# - Static files
```

### Deployment Options

#### Docker Deployment
```dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

#### Static Hosting
- **Netlify**: Automatic deployments from Git
- **Vercel**: Serverless deployment with edge functions
- **AWS S3 + CloudFront**: Scalable static hosting
- **Render**: Full-stack deployment with API

### Environment Variables for Production

```env
# Production API endpoints
VITE_API_BASE_URL=https://api.dklemailservice.com/api
VITE_WS_BASE_URL=wss://api.dklemailservice.com

# Security
VITE_ENABLE_HTTPS=true
VITE_CSP_ENABLED=true

# Analytics
VITE_GA_TRACKING_ID=GA_MEASUREMENT_ID
```

## Testing

### Testing Strategy

#### Unit Tests (Vitest)
```typescript
// Component testing with React Testing Library
import { render, screen, fireEvent } from '@testing-library/react';
import { UserForm } from './UserForm';

test('renders user form correctly', () => {
  render(<UserForm />);
  expect(screen.getByLabelText('Name')).toBeInTheDocument();
  expect(screen.getByLabelText('Email')).toBeInTheDocument();
});

test('validates email input', async () => {
  render(<UserForm />);
  const emailInput = screen.getByLabelText('Email');

  fireEvent.change(emailInput, { target: { value: 'invalid-email' } });
  fireEvent.blur(emailInput);

  expect(await screen.findByText('Invalid email format')).toBeInTheDocument();
});
```

#### Integration Tests (Vitest + MSW)
```typescript
// API integration testing with MSW
import { rest } from 'msw';
import { setupServer } from 'msw/node';
import { render, screen, waitFor } from '@testing-library/react';
import { UserList } from './UserList';

const server = setupServer(
  rest.get('/api/users', (req, res, ctx) => {
    return res(ctx.json([
      { id: 1, name: 'John Doe', email: 'john@example.com' }
    ]));
  })
);

test('loads and displays users', async () => {
  server.listen();
  render(<UserList />);

  await waitFor(() => {
    expect(screen.getByText('John Doe')).toBeInTheDocument();
  });

  server.close();
});
```

#### E2E Tests (Playwright)
```typescript
// Playwright E2E testing
import { test, expect } from '@playwright/test';

test('admin can create new user', async ({ page }) => {
  // Login as admin
  await page.goto('/login');
  await page.fill('[data-testid="email"]', 'admin@example.com');
  await page.fill('[data-testid="password"]', 'password');
  await page.click('[data-testid="login-button"]');

  // Navigate to user management
  await page.click('[data-testid="users-menu"]');

  // Create new user
  await page.click('[data-testid="create-user-button"]');
  await page.fill('[data-testid="user-name"]', 'Jane Smith');
  await page.fill('[data-testid="user-email"]', 'jane@example.com');
  await page.click('[data-testid="save-user-button"]');

  // Verify user was created
  await expect(page.locator('[data-testid="user-list"]')).toContainText('Jane Smith');
});
```

### Test Coverage Goals
- **Unit Tests**: 80%+ coverage for components and utilities
- **Integration Tests**: All critical user journeys
- **E2E Tests**: Core admin workflows

### Vitest Configuration
```typescript
/// <reference types="vitest" />
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react-swc';

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    coverage: {
      reporter: ['text', 'json', 'html'],
      exclude: ['node_modules/', 'src/test/'],
    },
  },
});
```

### Playwright Configuration
```typescript
// playwright.config.ts
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
  ],
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
  },
});
```

### CI/CD Pipeline
```yaml
# GitHub Actions workflow
name: CI/CD
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
      - run: npm run test:coverage
      - run: npm run test:e2e
      - run: npm run build
```

---

This documentation provides a comprehensive overview of the DKL Admin Panel. For specific implementation details, refer to the inline code comments and API documentation.