# DKL Website Frontend Documentation

Documentation for the planned DKL Website frontend - a modern, responsive web application that will serve as the public face of DKL events and provide participant registration functionality through integration with the DKL Email Service backend API.

## Table of Contents

- [Overview](#overview)
- [Target Users](#target-users)
- [Key Features](#key-features)
- [Technology Stack](#technology-stack)
- [Architecture](#architecture)
- [API Integration](#api-integration)
- [Security](#security)
- [User Interface](#user-interface)
- [SEO & Performance](#seo--performance)
- [Development Setup](#development-setup)
- [Deployment](#deployment)
- [Testing](#testing)
- [Backend API Reference](#backend-api-reference)

## Overview

The DKL Website is a planned Vite + React-based web application that will serve as the primary public interface for the DKL Email Service backend API. It will provide event information, participant registration, news updates, and serve as a gateway to the admin panel and mobile app.

### Current Status

**Note**: This frontend application has not been implemented yet. This document outlines the planned architecture and requirements for the frontend that will consume the existing DKL Email Service backend API.

### Purpose

The website will fulfill multiple roles:
- **Information Hub**: Comprehensive event information and updates via backend API
- **Registration Portal**: Public registration for events with optional account creation
- **Content Management**: Display of photos, videos, and event content (when CMS features are implemented)
- **Community Engagement**: News, newsletters, and social features
- **Administrative Access**: Gateway to admin panel for staff (redirect to separate admin interface)

## Target Users

### Primary User Groups
- **Event Participants**: People interested in registering for DKL events
- **General Public**: Visitors seeking information about DKL and events
- **Media & Partners**: Journalists and business partners
- **Staff & Administrators**: Event staff accessing admin panel

### User Journeys
- **First-time Visitors**: Discover DKL, learn about events, register interest
- **Returning Participants**: Check event status, update information, access app
- **Staff Members**: Access admin panel, manage content, monitor events
- **Partners**: View partnership information, contact details

## Key Features

### Public Information Pages
- **Event Listings**: Upcoming and past events with detailed information
- **Event Details**: Comprehensive event pages with schedules and routes
- **About DKL**: Organization history, mission, and team information
- **News & Updates**: Blog posts, newsletters, and announcements
- **Photo Galleries**: Event photo albums and media content
- **Video Content**: Event videos and promotional content

### Registration System
- **Public Registration**: Event registration with optional account creation
- **Account Management**: Profile updates and password management
- **Email Verification**: Secure account verification process
- **QR Code Generation**: Registration confirmation with QR codes

### Content Management Integration
- **Dynamic Content**: CMS-driven pages and content blocks
- **Media Galleries**: Photo and video galleries from events
- **Partner Showcases**: Business partner information and logos
- **Sponsor Recognition**: Event sponsor displays and information

### Administrative Access
- **Staff Login**: Secure access to admin panel
- **Role-based Redirects**: Automatic redirection based on user permissions
- **Session Management**: Secure session handling for admin access

### Social & Community Features
- **Newsletter Signup**: Email subscription for updates
- **Social Media Integration**: Links to social media platforms
- **Contact Forms**: General inquiries and support requests
- **Live Updates**: Real-time event status and announcements

## Technology Stack

### Frontend Framework
- **Vite 7+**: Fast build tool and development server with esbuild
- **React 18+**: Modern React with concurrent features and hooks
- **TypeScript**: Type-safe development with comprehensive type definitions
- **React Router DOM v7**: Client-side routing with data loading

### UI Framework & Styling
- **Material-UI (MUI)**: Comprehensive component library with theming
- **Emotion**: CSS-in-JS styling with styled components
- **Framer Motion**: Smooth animations and page transitions
- **Tailwind CSS**: Utility-first CSS with typography plugin
- **Heroicons**: Beautiful hand-crafted SVG icons

### Forms & Validation
- **React Hook Form**: Performant forms with validation
- **Zod**: TypeScript-first schema validation
- **@hookform/resolvers**: Integration between React Hook Form and Zod

### HTTP & Data Fetching
- **Axios**: HTTP client with interceptors for authentication
- **React Hot Toast**: Beautiful toast notifications
- **Lodash.debounce**: Debounced functions for performance

### Analytics & Monitoring
- **Vercel Analytics**: Real user monitoring and performance metrics
- **Vercel Speed Insights**: Core Web Vitals tracking
- **React GA4**: Google Analytics 4 integration

### Utilities
- **React Helmet Async**: Document head management for SEO
- **UUID**: Unique identifier generation
- **QRCode**: QR code generation for registrations
- **Focus Trap React**: Accessibility focus management

### Development Tools
- **ESLint**: Code linting with TypeScript and React rules
- **Prettier**: Code formatting with strict configuration
- **TypeScript**: Strict type checking
- **Vite Plugin PWA**: Progressive Web App support

## Architecture

### Project Structure

```
dkl25main/
├── src/
│   ├── components/       # Reusable React components
│   │   ├── ui/          # Material-UI based components
│   │   ├── forms/       # Form components with react-hook-form
│   │   ├── layout/      # Layout components (Header, Footer, Navigation)
│   │   ├── sections/    # Page sections (Hero, EventCard, etc.)
│   │   └── common/      # Shared components
│   ├── pages/           # Page components
│   │   ├── Home.tsx     # Home page
│   │   ├── Events.tsx   # Event listings
│   │   ├── EventDetail.tsx # Event details
│   │   ├── About.tsx    # About page
│   │   ├── Contact.tsx  # Contact page
│   │   ├── Register.tsx # Registration page
│   │   └── Admin.tsx    # Admin panel redirect
│   ├── hooks/           # Custom React hooks
│   ├── lib/             # Utility libraries
│   │   ├── api.ts       # Axios client configuration
│   │   ├── validations/ # Zod validation schemas
│   │   ├── utils/       # Helper functions
│   │   └── theme.ts     # Material-UI theme configuration
│   ├── types/           # TypeScript definitions
│   ├── App.tsx          # Main App component
│   ├── main.tsx         # Application entry point
│   └── index.css        # Global styles
├── public/              # Static assets
│   ├── favicon.ico
│   ├── manifest.json    # PWA manifest
│   └── icons/           # PWA icons
├── server.ts            # Development API server
├── index.html           # HTML template
├── vite.config.ts       # Vite configuration
├── tailwind.config.js   # Tailwind CSS configuration
├── tsconfig.json        # TypeScript configuration
├── package.json         # Dependencies and scripts
└── docker/              # Docker configuration
```

### React Router DOM v7 Structure

#### Public Routes
```typescript
// src/App.tsx
import { RouterProvider, createBrowserRouter } from 'react-router-dom';

const router = createBrowserRouter([
  {
    path: '/',
    element: <HomePage />,
  },
  {
    path: '/events',
    element: <EventsPage />,
  },
  {
    path: '/events/:id',
    element: <EventDetailPage />,
  },
  {
    path: '/about',
    element: <AboutPage />,
  },
  {
    path: '/contact',
    element: <ContactPage />,
  },
  {
    path: '/register',
    element: <RegisterPage />,
  },
  {
    path: '/admin',
    element: <AdminRedirect />,
  },
]);
```

#### Route-based Code Splitting
```typescript
// Lazy loading for performance
const EventsPage = lazy(() => import('./pages/Events'));
const EventDetailPage = lazy(() => import('./pages/EventDetail'));

// With React Router v7 data loading
const router = createBrowserRouter([
  {
    path: '/events',
    element: <EventsPage />,
    loader: eventsLoader, // Data loading function
  },
  {
    path: '/events/:id',
    element: <EventDetailPage />,
    loader: eventDetailLoader, // Data loading with params
  },
]);
```

### Axios Client Configuration

```typescript
// src/lib/api.ts
import axios, { AxiosInstance, AxiosResponse } from 'axios';
import { toast } from 'react-hot-toast';

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api',
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Request interceptor for authentication
    this.client.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('auth_token');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor for error handling
    this.client.interceptors.response.use(
      (response: AxiosResponse) => response,
      (error) => {
        if (error.response?.status === 401) {
          // Handle unauthorized access
          localStorage.removeItem('auth_token');
          window.location.href = '/login';
        } else if (error.response?.status >= 500) {
          toast.error('Server error. Please try again later.');
        }
        return Promise.reject(error);
      }
    );
  }

  // Public API methods
  async getEvents() {
    const response = await this.client.get('/events');
    return response.data;
  }

  async getEvent(id: string) {
    const response = await this.client.get(`/events/${id}`);
    return response.data;
  }

  async registerParticipant(data: RegistrationData) {
    const response = await this.client.post('/public/aanmelden', data);
    return response.data;
  }

  async login(credentials: LoginCredentials) {
    const response = await this.client.post('/auth/login', credentials);
    return response.data;
  }
}

export const api = new ApiClient();
```

### Form Handling with React Hook Form + Zod

```typescript
// src/lib/validations/schemas.ts
import { z } from 'zod';

export const registrationSchema = z.object({
  naam: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.string().email('Invalid email address'),
  telefoon: z.string().regex(/^\+?[0-9\s\-\(\)]+$/, 'Invalid phone number'),
  rol: z.enum(['Deelnemer', 'Begeleider', 'Vrijwilliger'], {
    errorMap: () => ({ message: 'Please select a valid role' }),
  }),
  afstand: z.enum(['2.5 KM', '6 KM', '10 KM', '15 KM']),
  ondersteuning: z.enum(['Ja', 'Nee', 'Anders']),
  heeft_vervoer: z.boolean(),
  terms: z.boolean().refine(val => val, 'You must accept the terms'),
  event_id: z.string().uuid('Invalid event ID'),
  want_account: z.boolean(),
  wachtwoord: z.string().optional(),
}).refine(
  (data) => !data.want_account || (data.wachtwoord && data.wachtwoord.length >= 8),
  {
    message: 'Password must be at least 8 characters when creating an account',
    path: ['wachtwoord'],
  }
);

export type RegistrationFormData = z.infer<typeof registrationSchema>;
```

```typescript
// src/components/forms/RegistrationForm.tsx
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { TextField, Button, Checkbox, FormControlLabel } from '@mui/material';
import { toast } from 'react-hot-toast';
import { useNavigate } from 'react-router-dom';
import { api } from '@/lib/api';
import { registrationSchema, RegistrationFormData } from '@/lib/validations/schemas';

interface RegistrationFormProps {
  eventId: string;
  onSuccess: (data: any) => void;
}

export function RegistrationForm({ eventId, onSuccess }: RegistrationFormProps) {
  const navigate = useNavigate();
  const {
    register,
    handleSubmit,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<RegistrationFormData>({
    resolver: zodResolver(registrationSchema),
    defaultValues: {
      event_id: eventId,
      heeft_vervoer: false,
      terms: false,
      want_account: false,
    },
  });

  const wantAccount = watch('want_account');

  const onSubmit = async (data: RegistrationFormData) => {
    try {
      const response = await api.registerParticipant(data);
      toast.success(response.message);
      onSuccess(response);
    } catch (error) {
      toast.error('Registration failed. Please try again.');
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <TextField
        {...register('naam')}
        label="Full Name"
        fullWidth
        error={!!errors.naam}
        helperText={errors.naam?.message}
      />

      <TextField
        {...register('email')}
        label="Email"
        type="email"
        fullWidth
        error={!!errors.email}
        helperText={errors.email?.message}
      />

      <TextField
        {...register('telefoon')}
        label="Phone"
        fullWidth
        error={!!errors.telefoon}
        helperText={errors.telefoon?.message}
      />

      {/* Other form fields */}

      <FormControlLabel
        control={<Checkbox {...register('want_account')} />}
        label="Create a full account for the DKL Steps App"
      />

      {wantAccount && (
        <TextField
          {...register('wachtwoord')}
          label="Password"
          type="password"
          fullWidth
          error={!!errors.wachtwoord}
          helperText={errors.wachtwoord?.message}
        />
      )}

      <FormControlLabel
        control={<Checkbox {...register('terms')} />}
        label="I accept the terms and conditions"
      />
      {errors.terms && (
        <p className="text-red-500 text-sm">{errors.terms.message}</p>
      )}

      <Button
        type="submit"
        variant="contained"
        fullWidth
        disabled={isSubmitting}
      >
        {isSubmitting ? 'Registering...' : 'Register'}
      </Button>
    </form>
  );
}
```

## API Integration

### Backend API Overview

The frontend will integrate with the DKL Email Service backend API, which provides the following key endpoints:

#### Public Endpoints (No Authentication Required)
- `GET /api/health` - Service health check
- `POST /api/public/aanmelden` - Public event registration
- `POST /api/contact-email` - Contact form submission
- `POST /api/auth/login` - User authentication
- `POST /api/auth/forgot-password` - Password reset request
- `POST /api/auth/reset-password-with-token` - Password reset with token

#### Protected Endpoints (JWT Authentication Required)
- `GET /api/auth/profile` - User profile information
- `POST /api/auth/logout` - User logout
- Various admin endpoints for contact/registration management

### API Client Architecture

```typescript
// src/lib/api.ts
class ApiClient {
  private baseURL: string;
  private token: string | null = null;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  setToken(token: string) {
    this.token = token;
    localStorage.setItem('auth_token', token);
  }

  getToken(): string | null {
    return this.token || localStorage.getItem('auth_token');
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    const headers = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    const token = this.getToken();
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }

    const response = await fetch(url, {
      ...options,
      headers,
    });

    if (!response.ok) {
      // Handle authentication errors
      if (response.status === 401) {
        this.setToken('');
        window.location.href = '/login';
      }
      throw new Error(`API Error: ${response.status}`);
    }

    return response.json();
  }

  // Public endpoints
  async getHealth() {
    return this.request<{status: string, version: string}>('/health');
  }

  // Registration endpoint
  async registerParticipant(data: RegistrationData) {
    return this.request('/public/aanmelden', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  // Contact form endpoint
  async submitContactForm(data: ContactData) {
    return this.request('/contact-email', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  // Authentication
  async login(credentials: LoginCredentials) {
    const response = await this.request<{token: string, user: any}>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });
    this.setToken(response.token);
    return response;
  }
}

export const api = new ApiClient(
  import.meta.env.VITE_API_URL || 'http://localhost:8080/api'
);
```

### Data Fetching Patterns

#### Server-Side Data Fetching
```typescript
// Server Component with SWR for client updates
import useSWR from 'swr';

export default function EventsPage() {
  const { data: events, error } = useSWR('/events', api.getEvents, {
    revalidateOnFocus: false,
    dedupingInterval: 300000, // 5 minutes
  });

  if (error) return <div>Failed to load events</div>;
  if (!events) return <div>Loading...</div>;

  return (
    <div>
      {events.map(event => (
        <EventCard key={event.id} event={event} />
      ))}
    </div>
  );
}
```

#### Form Handling with Validation
```typescript
// Client Component with React Hook Form
'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

const registrationSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.string().email('Invalid email address'),
  eventId: z.string().uuid('Invalid event ID'),
  acceptTerms: z.boolean().refine(val => val, 'You must accept the terms'),
});

type RegistrationForm = z.infer<typeof registrationSchema>;

export function RegistrationForm({ eventId }: { eventId: string }) {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<RegistrationForm>({
    resolver: zodResolver(registrationSchema),
    defaultValues: {
      eventId,
      acceptTerms: false,
    },
  });

  const onSubmit = async (data: RegistrationForm) => {
    try {
      await api.register(data);
      // Handle success
    } catch (error) {
      // Handle error
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      {/* Form fields with validation */}
    </form>
  );
}
```

### Authentication Integration

```typescript
// lib/auth/client.ts
export class AuthClient {
  async login(credentials: LoginCredentials) {
    const response = await api.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify(credentials),
    });

    const { token, refresh_token, user } = response.data;

    // Store tokens securely
    this.setTokens(token, refresh_token);

    return { token, user };
  }

  async refreshToken() {
    const refreshToken = this.getRefreshToken();
    if (!refreshToken) throw new Error('No refresh token');

    const response = await api.request('/auth/refresh', {
      method: 'POST',
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    const { token, refresh_token: newRefreshToken } = response.data;
    this.setTokens(token, newRefreshToken);

    return token;
  }

  private setTokens(accessToken: string, refreshToken: string) {
    // Store in httpOnly cookies for security
    document.cookie = `access_token=${accessToken}; path=/; secure; samesite=strict`;
    document.cookie = `refresh_token=${refreshToken}; path=/; secure; samesite=strict; httponly`;
  }

  private getRefreshToken(): string | null {
    // Read from httpOnly cookie (server-side)
    return null; // Cookies are handled server-side
  }
}
```

## Backend API Reference

This section provides key API endpoints that the frontend will need to integrate with. For complete API documentation, see the Swagger documentation at `/swagger/*` when the backend is running.

### Authentication Endpoints

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "wachtwoord": "password"
}
```

**Response:**
```json
{
  "token": "jwt_token_here",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "naam": "User Name"
  }
}
```

#### Password Reset
```http
POST /api/auth/forgot-password
Content-Type: application/json

{
  "email": "user@example.com"
}
```

### Public Registration

#### Event Registration
```http
POST /api/public/aanmelden
Content-Type: application/json

{
  "naam": "John Doe",
  "email": "john@example.com",
  "telefoon": "+31612345678",
  "rol": "Deelnemer",
  "afstand": "10 KM",
  "ondersteuning": "Nee",
  "heeft_vervoer": false,
  "terms": true,
  "event_id": "event-uuid",
  "want_account": false,
  "wachtwoord": "optional_password_if_want_account_true"
}
```

### Contact Form

#### Submit Contact Form
```http
POST /api/contact-email
Content-Type: application/json

{
  "naam": "John Doe",
  "email": "john@example.com",
  "bericht": "Contact message",
  "privacy_akkoord": true
}
```

### Health Check

#### Service Health
```http
GET /api/health
```

**Response:**
```json
{
  "service": "DKL Email Service API",
  "version": "1.0",
  "status": "running",
  "environment": "development",
  "timestamp": "2024-03-20T10:30:00Z",
  "checks": {
    "database": true,
    "redis": true,
    "smtp": true
  }
}
```

### Error Handling

The API returns standard HTTP status codes:
- `200` - Success
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid/missing token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `429` - Too Many Requests (rate limited)
- `500` - Internal Server Error

Error responses include a JSON body:
```json
{
  "error": "Error message description"
}
```

### Rate Limiting

The API implements rate limiting:
- Contact forms: 100/hour globally, 5/hour per IP
- Registration: 200/hour globally, 10/hour per IP
- Authentication: 5 attempts per 5 minutes per IP

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1640995200
```

## Security

### Authentication & Authorization

#### JWT Token Management
- **HttpOnly Cookies**: Secure token storage preventing XSS attacks
- **CSRF Protection**: SameSite cookie attributes and CSRF tokens
- **Token Refresh**: Automatic token renewal before expiration
- **Secure Headers**: Content Security Policy and security headers

#### Server-Side Authentication
```typescript
// middleware.ts
import { NextRequest, NextResponse } from 'next/server';
import { jwtVerify } from 'jose';

export async function middleware(request: NextRequest) {
  const token = request.cookies.get('access_token')?.value;

  if (!token) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  try {
    const { payload } = await jwtVerify(
      token,
      new TextEncoder().encode(process.env.JWT_SECRET)
    );

    // Add user info to request headers
    const requestHeaders = new Headers(request.headers);
    requestHeaders.set('x-user-id', payload.sub as string);
    requestHeaders.set('x-user-roles', JSON.stringify(payload.roles));

    return NextResponse.next({
      request: {
        headers: requestHeaders,
      },
    });
  } catch (error) {
    // Token invalid, redirect to login
    return NextResponse.redirect(new URL('/login', request.url));
  }
}

export const config = {
  matcher: ['/admin/:path*', '/profile/:path*'],
};
```

### Input Validation & Sanitization

#### Client-Side Validation
```typescript
// lib/validations/schemas.ts
import { z } from 'zod';

export const contactFormSchema = z.object({
  name: z.string()
    .min(2, 'Name must be at least 2 characters')
    .max(100, 'Name must be less than 100 characters')
    .regex(/^[a-zA-Z\s]+$/, 'Name can only contain letters and spaces'),

  email: z.string()
    .email('Please enter a valid email address')
    .max(254, 'Email address is too long'),

  message: z.string()
    .min(10, 'Message must be at least 10 characters')
    .max(1000, 'Message must be less than 1000 characters'),

  acceptPrivacy: z.boolean()
    .refine(val => val, 'You must accept the privacy policy'),
});

export type ContactFormData = z.infer<typeof contactFormSchema>;
```

#### Server-Side Validation
```typescript
// app/api/contact/route.ts
import { NextRequest, NextResponse } from 'next/server';
import { contactFormSchema } from '@/lib/validations/schemas';
import { sanitizeHtml } from '@/lib/utils/sanitize';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();

    // Validate input
    const validatedData = contactFormSchema.parse(body);

    // Sanitize HTML content
    const sanitizedMessage = sanitizeHtml(validatedData.message);

    // Process the contact form
    await processContactForm({
      ...validatedData,
      message: sanitizedMessage,
    });

    return NextResponse.json({ success: true });
  } catch (error) {
    if (error instanceof z.ZodError) {
      return NextResponse.json(
        { error: 'Validation failed', details: error.errors },
        { status: 400 }
      );
    }

    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    );
  }
}
```

### Security Headers & CSP

```typescript
// next.config.js
module.exports = {
  async headers() {
    return [
      {
        source: '/(.*)',
        headers: [
          {
            key: 'X-Frame-Options',
            value: 'DENY',
          },
          {
            key: 'X-Content-Type-Options',
            value: 'nosniff',
          },
          {
            key: 'Referrer-Policy',
            value: 'strict-origin-when-cross-origin',
          },
          {
            key: 'Content-Security-Policy',
            value: [
              "default-src 'self'",
              "script-src 'self' 'unsafe-inline' 'unsafe-eval'",
              "style-src 'self' 'unsafe-inline'",
              "img-src 'self' data: https:",
              "font-src 'self'",
              "connect-src 'self' https://api.dklemailservice.com",
              "frame-ancestors 'none'",
            ].join('; '),
          },
        ],
      },
    ];
  },
};
```

## User Interface

### Design System

#### Color Palette
```css
:root {
  --color-primary: #1976d2;
  --color-secondary: #dc004e;
  --color-success: #4caf50;
  --color-warning: #ff9800;
  --color-error: #f44336;
  --color-background: #ffffff;
  --color-surface: #f8fafc;
  --color-text-primary: #1e293b;
  --color-text-secondary: #64748b;
  --color-border: #e2e8f0;
}
```

#### Typography Scale
```css
:root {
  --font-size-xs: 0.75rem;    /* 12px */
  --font-size-sm: 0.875rem;   /* 14px */
  --font-size-base: 1rem;     /* 16px */
  --font-size-lg: 1.125rem;   /* 18px */
  --font-size-xl: 1.25rem;    /* 20px */
  --font-size-2xl: 1.5rem;    /* 24px */
  --font-size-3xl: 1.875rem;  /* 30px */
  --font-size-4xl: 2.25rem;   /* 36px */
  --font-size-5xl: 3rem;      /* 48px */
}
```

### Responsive Design

#### Breakpoint System
```css
/* Mobile First Approach */
.container {
  width: 100%;
  padding-left: 1rem;
  padding-right: 1rem;
  margin-left: auto;
  margin-right: auto;
}

/* Small devices (tablets, 640px and up) */
@media (min-width: 640px) {
  .container {
    max-width: 640px;
  }
}

/* Medium devices (desktops, 768px and up) */
@media (min-width: 768px) {
  .container {
    max-width: 768px;
  }
}

/* Large devices (large desktops, 1024px and up) */
@media (min-width: 1024px) {
  .container {
    max-width: 1024px;
  }
}

/* Extra large devices (extra large desktops, 1280px and up) */
@media (min-width: 1280px) {
  .container {
    max-width: 1280px;
  }
}
```

### Component Architecture

#### Atomic Design Pattern
```
components/
├── atoms/          # Basic elements (Button, Input, Icon)
├── molecules/      # Simple combinations (FormField, Card, Badge)
├── organisms/      # Complex components (Header, Footer, EventCard)
└── templates/      # Page layouts (HomePage, EventPage)
```

#### Example Components

##### Button Component
```typescript
// components/ui/Button.tsx
interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  loading?: boolean;
  fullWidth?: boolean;
}

export function Button({
  variant = 'primary',
  size = 'md',
  loading = false,
  fullWidth = false,
  className,
  children,
  ...props
}: ButtonProps) {
  return (
    <button
      className={cn(
        'inline-flex items-center justify-center rounded-md font-medium transition-colors',
        'focus:outline-none focus:ring-2 focus:ring-offset-2',
        'disabled:opacity-50 disabled:pointer-events-none',
        {
          'bg-primary text-primary-foreground hover:bg-primary/90': variant === 'primary',
          'bg-secondary text-secondary-foreground hover:bg-secondary/80': variant === 'secondary',
          'border border-input bg-background hover:bg-accent': variant === 'outline',
          'hover:bg-accent hover:text-accent-foreground': variant === 'ghost',
          'h-9 px-3 text-sm': size === 'sm',
          'h-10 px-4 py-2': size === 'md',
          'h-11 px-8': size === 'lg',
          'w-full': fullWidth,
        },
        className
      )}
      disabled={loading}
      {...props}
    >
      {loading && <Spinner className="mr-2 h-4 w-4" />}
      {children}
    </button>
  );
}
```

##### Event Card Component
```typescript
// components/organisms/EventCard.tsx
interface EventCardProps {
  event: Event;
  showRegisterButton?: boolean;
}

export function EventCard({ event, showRegisterButton = true }: EventCardProps) {
  const router = useRouter();

  return (
    <Card className="overflow-hidden transition-shadow hover:shadow-lg">
      <div className="aspect-video relative">
        <Image
          src={event.imageUrl}
          alt={event.name}
          fill
          className="object-cover"
        />
        <Badge className="absolute top-2 right-2">
          {event.status}
        </Badge>
      </div>

      <CardContent className="p-6">
        <div className="space-y-2">
          <h3 className="text-xl font-semibold">{event.name}</h3>
          <p className="text-muted-foreground">{event.description}</p>

          <div className="flex items-center space-x-4 text-sm text-muted-foreground">
            <div className="flex items-center">
              <CalendarIcon className="mr-1 h-4 w-4" />
              {formatDate(event.startTime)}
            </div>
            <div className="flex items-center">
              <MapPinIcon className="mr-1 h-4 w-4" />
              {event.location}
            </div>
          </div>
        </div>

        {showRegisterButton && (
          <div className="mt-4 flex space-x-2">
            <Button
              onClick={() => router.push(`/events/${event.id}`)}
              className="flex-1"
            >
              View Details
            </Button>
            <Button
              variant="outline"
              onClick={() => router.push(`/register?event=${event.id}`)}
            >
              Register
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
```

## SEO & Performance

### Search Engine Optimization

#### Meta Tags & Structured Data
```typescript
// app/events/[id]/page.tsx
import { Metadata } from 'next';

interface EventPageProps {
  params: { id: string };
}

export async function generateMetadata({ params }: EventPageProps): Promise<Metadata> {
  const event = await api.getEvent(params.id);

  return {
    title: `${event.name} | DKL Events`,
    description: event.description,
    openGraph: {
      title: event.name,
      description: event.description,
      images: [event.imageUrl],
      type: 'website',
    },
    jsonLd: {
      '@context': 'https://schema.org',
      '@type': 'Event',
      name: event.name,
      description: event.description,
      startDate: event.startTime,
      endDate: event.endTime,
      location: {
        '@type': 'Place',
        name: event.location,
        address: event.address,
      },
      organizer: {
        '@type': 'Organization',
        name: 'De Koninklijke Loop',
        url: 'https://dekoninklijkeloop.nl',
      },
    },
  };
}
```

#### Dynamic Sitemap Generation
```typescript
// app/sitemap.ts
import { MetadataRoute } from 'next';

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const events = await api.getEvents();
  const news = await api.getNews();

  const eventUrls = events.map(event => ({
    url: `https://dekoninklijkeloop.nl/events/${event.id}`,
    lastModified: new Date(event.updatedAt),
    changeFrequency: 'weekly' as const,
    priority: 0.8,
  }));

  const newsUrls = news.map(item => ({
    url: `https://dekoninklijkeloop.nl/news/${item.id}`,
    lastModified: new Date(item.updatedAt),
    changeFrequency: 'daily' as const,
    priority: 0.6,
  }));

  return [
    {
      url: 'https://dekoninklijkeloop.nl',
      lastModified: new Date(),
      changeFrequency: 'daily',
      priority: 1,
    },
    ...eventUrls,
    ...newsUrls,
  ];
}
```

### Performance Optimization

#### Image Optimization
```typescript
// components/ui/OptimizedImage.tsx
import Image from 'next/image';
import { useState } from 'react';

interface OptimizedImageProps {
  src: string;
  alt: string;
  width: number;
  height: number;
  priority?: boolean;
}

export function OptimizedImage({ src, alt, width, height, priority }: OptimizedImageProps) {
  const [isLoading, setIsLoading] = useState(true);

  return (
    <div className="relative overflow-hidden rounded-lg">
      {isLoading && (
        <div className="absolute inset-0 animate-pulse bg-gray-200" />
      )}
      <Image
        src={src}
        alt={alt}
        width={width}
        height={height}
        priority={priority}
        onLoad={() => setIsLoading(false)}
        className="transition-opacity duration-300"
        style={{ opacity: isLoading ? 0 : 1 }}
      />
    </div>
  );
}
```

#### Code Splitting & Lazy Loading
```typescript
// app/events/page.tsx
import dynamic from 'next/dynamic';
import { Suspense } from 'react';

// Lazy load heavy components
const EventFilters = dynamic(() => import('@/components/EventFilters'), {
  loading: () => <div>Loading filters...</div>,
});

const EventMap = dynamic(() => import('@/components/EventMap'), {
  loading: () => <div>Loading map...</div>,
  ssr: false, // Disable SSR for map component
});

export default function EventsPage() {
  return (
    <div>
      <Suspense fallback={<div>Loading...</div>}>
        <EventFilters />
      </Suspense>

      <Suspense fallback={<div>Loading map...</div>}>
        <EventMap />
      </Suspense>

      {/* Other content */}
    </div>
  );
}
```

#### Caching Strategy
```typescript
// lib/cache.ts
import { unstable_cache } from 'next/cache';

export const getCachedEvents = unstable_cache(
  async () => {
    return api.getEvents();
  },
  ['events'],
  {
    revalidate: 300, // 5 minutes
    tags: ['events'],
  }
);

export const getCachedEvent = unstable_cache(
  async (id: string) => {
    return api.getEvent(id);
  },
  ['event'],
  {
    revalidate: 60, // 1 minute
    tags: [`event-${id}`],
  }
);
```

## Development Setup

### Prerequisites
- **Node.js**: 18.17+ LTS
- **npm** or **yarn**: Package manager
- **Git**: Version control

### Installation

```bash
# Clone the repository
git clone https://github.com/dkl/dkl25main.git
cd dkl25main

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
VITE_API_URL=http://localhost:8080/api

# Analytics
VITE_GA_TRACKING_ID=GA_MEASUREMENT_ID

# Development
VITE_NODE_ENV=development
```

### Development Scripts

```json
{
  "scripts": {
    "dev": "vite",
    "build": "tsc --noEmit && vite build",
    "build:client": "vite build",
    "type-check": "tsc --noEmit",
    "preview": "vite preview",
    "lint": "eslint . --ext ts,tsx --report-unused-disable-directives --max-warnings 0",
    "format": "prettier --write \"src/**/*.{ts,tsx,js,jsx,json,css,md}\"",
    "format:check": "prettier --check \"src/**/*.{ts,tsx,js,jsx,json,css,md}\"",
    "clean": "rm -rf dist node_modules/.vite",
    "api": "ts-node server.ts",
    "docker:dev": "docker-compose -f docker-compose.dev.yml up",
    "docker:dev:build": "docker-compose -f docker-compose.dev.yml up --build",
    "docker:dev:down": "docker-compose -f docker-compose.dev.yml down",
    "docker:prod": "docker-compose up",
    "docker:prod:build": "docker-compose up --build",
    "docker:prod:down": "docker-compose down",
    "docker:build": "docker build -t dkl-frontend .",
    "docker:clean": "docker-compose down -v && docker system prune -f"
  }
}
```

### Code Quality Tools

#### ESLint Configuration
```javascript
// .eslintrc.js
module.exports = {
  extends: [
    'next/core-web-vitals',
    '@typescript-eslint/recommended',
    'prettier',
  ],
  rules: {
    '@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
    '@typescript-eslint/no-explicit-any': 'warn',
    'react-hooks/exhaustive-deps': 'error',
    'no-console': 'warn',
  },
};
```

#### Prettier Configuration
```javascript
// .prettierrc.js
module.exports = {
  semi: true,
  trailingComma: 'es5',
  singleQuote: true,
  printWidth: 80,
  tabWidth: 2,
  useTabs: false,
  plugins: ['prettier-plugin-tailwindcss'],
};
```

## Deployment

### Vercel Deployment (Recommended)

#### Automatic Deployment
```javascript
// vercel.json
{
  "buildCommand": "npm run build",
  "outputDirectory": "dist",
  "framework": null,
  "regions": ["fra1"],
  "headers": [
    {
      "source": "/api/(.*)",
      "headers": [
        {
          "key": "Cache-Control",
          "value": "no-cache"
        }
      ]
    }
  ],
  "rewrites": [
    {
      "source": "/api/(.*)",
      "destination": "https://api.dklemailservice.com/api/$1"
    }
  ]
}
```

#### Environment Variables
```bash
# Set environment variables in Vercel dashboard
VITE_API_URL=https://api.dklemailservice.com/api
VITE_GA_TRACKING_ID=GA_MEASUREMENT_ID
```

### Docker Deployment

#### Development Docker Setup
```yaml
# docker-compose.dev.yml
version: '3.8'
services:
  frontend:
    build:
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "5173:5173"
    volumes:
      - .:/app
      - /app/node_modules
    environment:
      - VITE_API_URL=http://host.docker.internal:8080/api
    command: npm run dev
```

#### Production Docker Setup
```yaml
# docker-compose.yml
version: '3.8'
services:
  frontend:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "80:80"
    environment:
      - VITE_API_URL=https://api.dklemailservice.com/api
      - VITE_GA_TRACKING_ID=GA_MEASUREMENT_ID
```

#### Dockerfile
```dockerfile
# Dockerfile
FROM node:22-alpine AS base

# Install dependencies
FROM base AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci --only=production

# Build the application
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

RUN npm run build

# Production image
FROM nginx:alpine AS runner
WORKDIR /app

# Copy built application
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

#### Nginx Configuration
```nginx
# nginx.conf
events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    server {
        listen 80;
        server_name localhost;
        root /usr/share/nginx/html;
        index index.html;

        # Enable gzip compression
        gzip on;
        gzip_vary on;
        gzip_min_length 1024;
        gzip_types text/plain text/css text/xml text/javascript application/javascript application/xml+rss application/json;

        # Security headers
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header X-XSS-Protection "1; mode=block" always;
        add_header Referrer-Policy "strict-origin-when-cross-origin" always;

        # Handle client-side routing
        location / {
            try_files $uri $uri/ /index.html;
        }

        # Cache static assets
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }
}
```

### CI/CD Pipeline

```yaml
# .github/workflows/deploy.yml
name: Deploy to Production
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'npm'
      - run: npm ci
      - run: npm run lint
      - run: npm run type-check
      - run: npm run test -- --coverage
      - run: npm run build

  deploy:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v3
      - uses: Vercel-Deploy-Action@v1
        with:
          vercel-token: ${{ secrets.VERCEL_TOKEN }}
          vercel-org-id: ${{ secrets.VERCEL_ORG_ID }}
          vercel-project-id: ${{ secrets.VERCEL_PROJECT_ID }}
```

## Testing

### Current Testing Setup

The project currently does not have a testing framework configured. To add comprehensive testing, consider the following setup:

### Recommended Testing Stack

#### Unit & Integration Testing (Vitest)
```bash
npm install -D vitest @testing-library/react @testing-library/jest-dom @testing-library/user-event jsdom
```

#### E2E Testing (Playwright)
```bash
npm install -D @playwright/test
npx playwright install
```

### Testing Configuration

#### Vitest Configuration
```typescript
/// <reference types="vitest" />
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

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

#### Sample Test Structure
```typescript
// src/components/__tests__/RegistrationForm.test.tsx
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { RegistrationForm } from '../RegistrationForm';

// Mock the API
const mockApi = vi.fn();
vi.mock('@/lib/api', () => ({
  api: {
    registerParticipant: mockApi,
  },
}));

describe('RegistrationForm', () => {
  it('submits form data correctly', async () => {
    const user = userEvent.setup();
    mockApi.mockResolvedValue({ success: true });

    render(<RegistrationForm eventId="event-123" onSuccess={() => {}} />);

    await user.type(screen.getByLabelText('Full Name'), 'John Doe');
    await user.type(screen.getByLabelText('Email'), 'john@example.com');
    await user.click(screen.getByLabelText('I accept the terms and conditions'));
    await user.click(screen.getByRole('button', { name: 'Register' }));

    await waitFor(() => {
      expect(mockApi).toHaveBeenCalledWith({
        naam: 'John Doe',
        email: 'john@example.com',
        event_id: 'event-123',
        terms: true,
        heeft_vervoer: false,
        want_account: false,
      });
    });
  });
});
```

#### Playwright E2E Tests
```typescript
// e2e/registration.spec.ts
import { test, expect } from '@playwright/test';

test('complete registration flow', async ({ page }) => {
  // Navigate to registration page
  await page.goto('/register');

  // Fill out registration form
  await page.fill('[data-testid="name-input"]', 'Jane Smith');
  await page.fill('[data-testid="email-input"]', 'jane@example.com');
  await page.check('[data-testid="terms-checkbox"]');
  await page.click('[data-testid="register-button"]');

  // Verify success message
  await expect(page.locator('[data-testid="success-message"]')).toBeVisible();
  await expect(page.locator('[data-testid="success-message"]')).toContainText(
    'Registration successful'
  );
});
```

### Performance Testing

#### Lighthouse CI
```yaml
# .github/workflows/lighthouse.yml
name: Lighthouse CI
on: [pull_request]

jobs:
  lighthouse:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '22'
      - run: npm ci
      - run: npm run build
      - run: npm run preview &
      - run: npx lhci autorun
```

---

This documentation provides a comprehensive overview of the DKL Website architecture and development practices. For specific implementation details, refer to the inline code comments and component documentation.