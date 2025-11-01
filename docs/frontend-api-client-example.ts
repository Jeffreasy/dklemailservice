/**
 * DKL API Client - Ready to use
 * 
 * Kopieer dit bestand naar je frontend project als:
 * - src/services/api.ts (Vite/React)
 * - lib/api.ts (Next.js)
 * - utils/api.ts (Create React App)
 */

import axios, { AxiosInstance, AxiosError } from 'axios';

// =============================================================================
// CONFIGURATIE
// =============================================================================

interface APIConfig {
  baseURL: string;
  timeout: number;
}

// Haal API URL uit environment variabelen (past zich aan aan je framework)
const getAPIBaseURL = (): string => {
  // Vite
  if (typeof import.meta !== 'undefined' && import.meta.env?.VITE_API_BASE_URL) {
    return import.meta.env.VITE_API_BASE_URL;
  }
  
  // Next.js
  if (process.env.NEXT_PUBLIC_API_BASE_URL) {
    return process.env.NEXT_PUBLIC_API_BASE_URL;
  }
  
  // Create React App
  if (process.env.REACT_APP_API_BASE_URL) {
    return process.env.REACT_APP_API_BASE_URL;
  }
  
  // Fallback naar development
  return 'http://localhost:8082/api';
};

const API_CONFIG: APIConfig = {
  baseURL: getAPIBaseURL(),
  timeout: 30000, // 30 seconden
};

// =============================================================================
// API CLIENT SETUP
// =============================================================================

export const apiClient: AxiosInstance = axios.create(API_CONFIG);

// Request Interceptor - Voeg JWT token toe
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response Interceptor - Handle errors
apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    // Token expired of unauthorized
    if (error.response?.status === 401) {
      localStorage.removeItem('auth_token');
      localStorage.removeItem('refresh_token');
      
      // Redirect to login (pas aan naar jouw routing)
      if (typeof window !== 'undefined') {
        window.location.href = '/login';
      }
    }
    
    return Promise.reject(error);
  }
);

// =============================================================================
// API METHODS
// =============================================================================

export const api = {
  // -------------------------
  // AUTHENTICATION
  // -------------------------
  auth: {
    login: async (email: string, password: string) => {
      const response = await apiClient.post('/auth/login', {
        email,
        wachtwoord: password  // API verwacht 'wachtwoord' (Nederlands)
      });
      
      // Sla tokens op
      if (response.data.token) {
        localStorage.setItem('auth_token', response.data.token);
      }
      if (response.data.refresh_token) {
        localStorage.setItem('refresh_token', response.data.refresh_token);
      }
      
      return response.data;
    },
    
    logout: async () => {
      await apiClient.post('/auth/logout');
      localStorage.removeItem('auth_token');
      localStorage.removeItem('refresh_token');
    },
    
    getProfile: async () => {
      const response = await apiClient.get('/auth/profile');
      return response.data;
    },
    
    resetPassword: async (currentPassword: string, newPassword: string) => {
      const response = await apiClient.post('/auth/reset-password', {
        huidig_wachtwoord: currentPassword,
        nieuw_wachtwoord: newPassword,
      });
      return response.data;
    },
    
    refreshToken: async () => {
      const refreshToken = localStorage.getItem('refresh_token');
      const response = await apiClient.post('/auth/refresh', {
        refresh_token: refreshToken,
      });
      
      if (response.data.token) {
        localStorage.setItem('auth_token', response.data.token);
      }
      
      return response.data;
    },
  },

  // -------------------------
  // CONTACT FORMULIEREN
  // -------------------------
  contacts: {
    list: async (limit = 50, offset = 0) => {
      const response = await apiClient.get(`/contact?limit=${limit}&offset=${offset}`);
      return response.data;
    },
    
    getById: async (id: string) => {
      const response = await apiClient.get(`/contact/${id}`);
      return response.data;
    },
    
    filterByStatus: async (status: string) => {
      const response = await apiClient.get(`/contact/status/${status}`);
      return response.data;
    },
    
    update: async (id: string, data: any) => {
      const response = await apiClient.put(`/contact/${id}`, data);
      return response.data;
    },
    
    delete: async (id: string) => {
      const response = await apiClient.delete(`/contact/${id}`);
      return response.data;
    },
    
    addReply: async (id: string, message: string) => {
      const response = await apiClient.post(`/contact/${id}/antwoord`, {
        bericht: message,
      });
      return response.data;
    },
  },

  // -------------------------
  // AANMELDINGEN
  // -------------------------
  registrations: {
    list: async (limit = 50, offset = 0) => {
      const response = await apiClient.get(`/aanmelding?limit=${limit}&offset=${offset}`);
      return response.data;
    },
    
    getById: async (id: string) => {
      const response = await apiClient.get(`/aanmelding/${id}`);
      return response.data;
    },
    
    filterByRole: async (role: string) => {
      const response = await apiClient.get(`/aanmelding/rol/${role}`);
      return response.data;
    },
    
    update: async (id: string, data: any) => {
      const response = await apiClient.put(`/aanmelding/${id}`, data);
      return response.data;
    },
    
    delete: async (id: string) => {
      const response = await apiClient.delete(`/aanmelding/${id}`);
      return response.data;
    },
    
    addReply: async (id: string, message: string) => {
      const response = await apiClient.post(`/aanmelding/${id}/antwoord`, {
        bericht: message,
      });
      return response.data;
    },
  },

  // -------------------------
  // ALBUMS
  // -------------------------
  albums: {
    // Public endpoint
    listPublic: async () => {
      const response = await apiClient.get('/albums');
      return response.data;
    },
    
    // Admin endpoints
    listAll: async (limit = 50, offset = 0) => {
      const response = await apiClient.get(`/albums/admin?limit=${limit}&offset=${offset}`);
      return response.data;
    },
    
    getById: async (id: string) => {
      const response = await apiClient.get(`/albums/${id}`);
      return response.data;
    },
    
    getPhotos: async (id: string) => {
      const response = await apiClient.get(`/albums/${id}/photos`);
      return response.data;
    },
    
    create: async (data: any) => {
      const response = await apiClient.post('/albums', data);
      return response.data;
    },
    
    update: async (id: string, data: any) => {
      const response = await apiClient.put(`/albums/${id}`, data);
      return response.data;
    },
    
    delete: async (id: string) => {
      const response = await apiClient.delete(`/albums/${id}`);
      return response.data;
    },
  },

  // -------------------------
  // PHOTOS
  // -------------------------
  photos: {
    listPublic: async (filters?: { year?: number; title?: string }) => {
      let url = '/photos';
      if (filters) {
        const params = new URLSearchParams();
        if (filters.year) params.append('year', filters.year.toString());
        if (filters.title) params.append('title', filters.title);
        url += '?' + params.toString();
      }
      const response = await apiClient.get(url);
      return response.data;
    },
    
    listAll: async (limit = 50, offset = 0) => {
      const response = await apiClient.get(`/photos/admin?limit=${limit}&offset=${offset}`);
      return response.data;
    },
    
    create: async (data: any) => {
      const response = await apiClient.post('/photos', data);
      return response.data;
    },
    
    update: async (id: string, data: any) => {
      const response = await apiClient.put(`/photos/${id}`, data);
      return response.data;
    },
    
    delete: async (id: string) => {
      const response = await apiClient.delete(`/photos/${id}`);
      return response.data;
    },
  },

  // -------------------------
  // VIDEOS
  // -------------------------
  videos: {
    listPublic: async () => {
      const response = await apiClient.get('/videos');
      return response.data;
    },
    
    listAll: async (limit = 50, offset = 0) => {
      const response = await apiClient.get(`/videos/admin?limit=${limit}&offset=${offset}`);
      return response.data;
    },
    
    create: async (data: any) => {
      const response = await apiClient.post('/videos', data);
      return response.data;
    },
    
    update: async (id: string, data: any) => {
      const response = await apiClient.put(`/videos/${id}`, data);
      return response.data;
    },
    
    delete: async (id: string) => {
      const response = await apiClient.delete(`/videos/${id}`);
      return response.data;
    },
  },

  // -------------------------
  // SPONSORS
  // -------------------------
  sponsors: {
    listPublic: async () => {
      const response = await apiClient.get('/sponsors');
      return response.data;
    },
    
    listAll: async (limit = 50, offset = 0) => {
      const response = await apiClient.get(`/sponsors/admin?limit=${limit}&offset=${offset}`);
      return response.data;
    },
    
    create: async (data: any) => {
      const response = await apiClient.post('/sponsors', data);
      return response.data;
    },
    
    update: async (id: string, data: any) => {
      const response = await apiClient.put(`/sponsors/${id}`, data);
      return response.data;
    },
    
    delete: async (id: string) => {
      const response = await apiClient.delete(`/sponsors/${id}`);
      return response.data;
    },
  },

  // -------------------------
  // STEPS TRACKING (NIEUW!)
  // -------------------------
  steps: {
    updateSteps: async (participantId: string, steps: number) => {
      const response = await apiClient.post(`/steps/${participantId}`, { steps });
      return response.data;
    },
    
    getParticipantDashboard: async (participantId: string) => {
      const response = await apiClient.get(`/participant/${participantId}/dashboard`);
      return response.data;
    },
    
    getTotalSteps: async () => {
      const response = await apiClient.get('/total-steps');
      return response.data;
    },
    
    getFundsDistribution: async () => {
      const response = await apiClient.get('/funds-distribution');
      return response.data;
    },
  },

  // -------------------------
  // USERS (Gebruikersbeheer)
  // -------------------------
  users: {
    list: async (limit = 50, offset = 0) => {
      const response = await apiClient.get(`/users?limit=${limit}&offset=${offset}`);
      return response.data;
    },
    
    getById: async (id: string) => {
      const response = await apiClient.get(`/users/${id}`);
      return response.data;
    },
    
    update: async (id: string, data: any) => {
      const response = await apiClient.put(`/users/${id}`, data);
      return response.data;
    },
    
    delete: async (id: string) => {
      const response = await apiClient.delete(`/users/${id}`);
      return response.data;
    },
    
    assignRole: async (userId: string, roleId: string) => {
      const response = await apiClient.post(`/users/${userId}/roles`, { role_id: roleId });
      return response.data;
    },
    
    removeRole: async (userId: string, roleId: string) => {
      const response = await apiClient.delete(`/users/${userId}/roles/${roleId}`);
      return response.data;
    },
  },

  // -------------------------
  // NEWSLETTER
  // -------------------------
  newsletter: {
    list: async (limit = 50, offset = 0) => {
      const response = await apiClient.get(`/newsletter?limit=${limit}&offset=${offset}`);
      return response.data;
    },
    
    getById: async (id: string) => {
      const response = await apiClient.get(`/newsletter/${id}`);
      return response.data;
    },
    
    create: async (data: any) => {
      const response = await apiClient.post('/newsletter', data);
      return response.data;
    },
    
    update: async (id: string, data: any) => {
      const response = await apiClient.put(`/newsletter/${id}`, data);
      return response.data;
    },
    
    delete: async (id: string) => {
      const response = await apiClient.delete(`/newsletter/${id}`);
      return response.data;
    },
    
    send: async (id: string) => {
      const response = await apiClient.post(`/newsletter/${id}/send`);
      return response.data;
    },
  },

  // -------------------------
  // HEALTH & MONITORING
  // -------------------------
  health: async () => {
    const response = await apiClient.get('/health');
    return response.data;
  },
};

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

/**
 * Check of user is ingelogd
 */
export const isAuthenticated = (): boolean => {
  return !!localStorage.getItem('auth_token');
};

/**
 * Get current auth token
 */
export const getAuthToken = (): string | null => {
  return localStorage.getItem('auth_token');
};

/**
 * Clear authentication
 */
export const clearAuth = (): void => {
  localStorage.removeItem('auth_token');
  localStorage.removeItem('refresh_token');
};

/**
 * Handle API errors met vriendelijke berichten
 */
export const handleAPIError = (error: any): string => {
  if (axios.isAxiosError(error)) {
    if (error.response) {
      // Server responded with error status
      switch (error.response.status) {
        case 400:
          return error.response.data?.error || 'Ongeldige invoer';
        case 401:
          return 'Niet geautoriseerd - log opnieuw in';
        case 403:
          return 'Geen toegang tot deze resource';
        case 404:
          return 'Niet gevonden';
        case 429:
          return 'Te veel verzoeken - probeer later opnieuw';
        case 500:
          return 'Server error - probeer later opnieuw';
        default:
          return error.response.data?.error || 'Er is iets misgegaan';
      }
    } else if (error.request) {
      // Request made but no response received
      return 'Geen verbinding met server - check of backend draait';
    }
  }
  
  return error.message || 'Onbekende fout';
};

/**
 * Check of API beschikbaar is
 */
export const checkAPIAvailability = async (): Promise<boolean> => {
  try {
    await api.health();
    return true;
  } catch {
    return false;
  }
};

// =============================================================================
// REACT HOOKS (Optioneel - gebruik indien nodig)
// =============================================================================

/**
 * useAPI Hook - Voor data fetching met loading/error states
 * 
 * Gebruik:
 * const { data, loading, error, refetch } = useAPI(() => api.contacts.list());
 */
export function useAPI<T>(
  apiCall: () => Promise<T>,
  deps: any[] = []
) {
  const [data, setData] = React.useState<T | null>(null);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  const fetchData = React.useCallback(async () => {
    setLoading(true);
    setError(null);
    
    try {
      const result = await apiCall();
      setData(result);
    } catch (err) {
      setError(handleAPIError(err));
    } finally {
      setLoading(false);
    }
  }, deps);

  React.useEffect(() => {
    fetchData();
  }, [fetchData]);

  return { data, loading, error, refetch: fetchData };
}

// =============================================================================
// TYPESCRIPT TYPES (Optioneel - voeg toe naar behoefte)
// =============================================================================

export interface LoginResponse {
  token: string;
  refresh_token: string;
  user: {
    id: string;
    email: string;
    naam: string;
    rol: string;
  };
}

export interface Contact {
  id: string;
  naam: string;
  email: string;
  bericht: string;
  status: 'nieuw' | 'in_behandeling' | 'beantwoord' | 'gesloten';
  beantwoord: boolean;
  created_at: string;
  updated_at: string;
}

export interface Registration {
  id: string;
  naam: string;
  email: string;
  telefoon: string;
  rol: string;
  afstand: string;
  status: 'nieuw' | 'bevestigd' | 'geannuleerd' | 'voltooid';
  steps?: number;
  created_at: string;
  updated_at: string;
}

export interface Album {
  id: string;
  title: string;
  description: string;
  cover_photo_id?: string;
  visible: boolean;
  order_number: number;
  created_at: string;
  updated_at: string;
}

export interface Photo {
  id: string;
  url: string;
  title: string;
  description?: string;
  year?: number;
  visible: boolean;
  cloudinary_public_id: string;
  created_at: string;
}

// =============================================================================
// USAGE EXAMPLES
// =============================================================================

/*
 * VOORBEELD 1: Login Component
 * 
 * import { api, handleAPIError } from './services/api';
 * 
 * const LoginForm = () => {
 *   const [email, setEmail] = useState('');
 *   const [password, setPassword] = useState('');
 *   const [error, setError] = useState('');
 * 
 *   const handleLogin = async (e) => {
 *     e.preventDefault();
 *     setError('');
 * 
 *     try {
 *       const response = await api.auth.login(email, password);
 *       console.log('Logged in:', response.user);
 *       // Redirect naar dashboard
 *       window.location.href = '/dashboard';
 *     } catch (err) {
 *       setError(handleAPIError(err));
 *     }
 *   };
 * 
 *   return <form onSubmit={handleLogin}>...</form>;
 * };
 */

/*
 * VOORBEELD 2: Contacts List Component
 * 
 * import { api, useAPI } from './services/api';
 * 
 * const ContactsList = () => {
 *   const { data, loading, error, refetch } = useAPI(() => api.contacts.list());
 * 
 *   if (loading) return <div>Laden...</div>;
 *   if (error) return <div>Error: {error}</div>;
 * 
 *   return (
 *     <div>
 *       {data?.map(contact => (
 *         <div key={contact.id}>{contact.naam}</div>
 *       ))}
 *       <button onClick={refetch}>Refresh</button>
 *     </div>
 *   );
 * };
 */

/*
 * VOORBEELD 3: Album Management
 * 
 * const AlbumManager = () => {
 *   const createAlbum = async () => {
 *     try {
 *       await api.albums.create({
 *         title: 'Nieuw Album 2025',
 *         description: 'Foto\'s van het evenement',
 *         visible: true,
 *         order_number: 1
 *       });
 *       alert('Album aangemaakt!');
 *     } catch (err) {
 *       alert(handleAPIError(err));
 *     }
 *   };
 * 
 *   return <button onClick={createAlbum}>Album Aanmaken</button>;
 * };
 */

/*
 * VOORBEELD 4: Steps Update
 * 
 * const StepsTracker = ({ participantId }) => {
 *   const updateSteps = async (steps: number) => {
 *     try {
 *       await api.steps.updateSteps(participantId, steps);
 *       alert('Steps bijgewerkt!');
 *     } catch (err) {
 *       alert(handleAPIError(err));
 *     }
 *   };
 * 
 *   return <button onClick={() => updateSteps(10000)}>Log 10k Steps</button>;
 * };
 */

export default api;