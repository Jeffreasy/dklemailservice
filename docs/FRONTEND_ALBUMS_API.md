# Frontend Albums API Documentatie

Deze documentatie beschrijft hoe je albums en photos kunt ophalen voor gebruik in de frontend.

## Base URLs

- **Development (Docker):** `http://localhost:8082`
- **Production:** `https://dklemailservice.onrender.com`

## Beschikbare Endpoints

### 1. Alle Zichtbare Albums Ophalen

**Endpoint:** `GET /api/albums`

Haalt alle zichtbare albums op, gesorteerd op `order_number`.

**Request:**
```typescript
const response = await fetch('https://dklemailservice.onrender.com/api/albums');
const albums = await response.json();
```

**Response:**
```typescript
interface Album {
  id: string;                    // UUID van het album
  title: string;                 // Titel (bijv. "DKL 2025")
  description: string | null;    // Beschrijving
  cover_photo_id: string | null; // UUID van cover photo
  visible: boolean;              // Of album zichtbaar is
  order_number: number;          // Volgorde (1, 2, 3...)
  created_at: string;           // ISO timestamp
  updated_at: string;           // ISO timestamp
}

// Response is een array van Album objecten
Album[]
```

**Voorbeeld Response:**
```json
[
  {
    "id": "d51cff45-b958-4370-a983-51e650ffa43e",
    "title": "DKL 2025",
    "description": "De koninklijke Loop 2025!",
    "cover_photo_id": "08362e92-340a-432a-b306-153ad27ee686",
    "visible": true,
    "order_number": 1,
    "created_at": "2025-05-17T20:10:00Z",
    "updated_at": "2025-05-17T20:10:06.643082Z"
  },
  {
    "id": "ce8df963-f118-4296-9f5c-e33308dc7bfa",
    "title": "DKL-2024",
    "description": "DKL 2024",
    "cover_photo_id": "8a4d5c20-ea73-4336-8c6b-af7a197ef7c2",
    "visible": true,
    "order_number": 2,
    "created_at": "2024-12-26T15:17:54.268Z",
    "updated_at": "2025-02-17T13:32:50.793Z"
  }
]
```

---

### 2. Albums Met Cover Photos Ophalen

**Endpoint:** `GET /api/albums?include_covers=true`

Haalt albums op inclusief de complete cover photo informatie.

**Request:**
```typescript
const response = await fetch('https://dklemailservice.onrender.com/api/albums?include_covers=true');
const albums = await response.json();
```

**Response:**
```typescript
interface AlbumWithCover {
  id: string;
  title: string;
  description: string | null;
  cover_photo_id: string | null;
  visible: boolean;
  order_number: number;
  created_at: string;
  updated_at: string;
  cover_photo: Photo | null;     // ⭐ Extra: volledige photo info
}

interface Photo {
  id: string;
  title: string;
  description: string | null;
  url: string;                   // Cloudinary image URL
  thumbnail_url: string | null;  // Thumbnail URL
  cloudinary_id: string;
  visible: boolean;
  width: number | null;
  height: number | null;
  format: string | null;
  created_at: string;
  updated_at: string;
}
```

**Voorbeeld Response:**
```json
[
  {
    "id": "d51cff45-b958-4370-a983-51e650ffa43e",
    "title": "DKL 2025",
    "description": "De koninklijke Loop 2025!",
    "cover_photo_id": "08362e92-340a-432a-b306-153ad27ee686",
    "visible": true,
    "order_number": 1,
    "cover_photo": {
      "id": "08362e92-340a-432a-b306-153ad27ee686",
      "title": "WhatsApp Image 2025-05-17 at 19",
      "url": "https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513127/thhu8mxkqhfhxjydi5ai.jpg",
      "thumbnail_url": "https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513127/thhu8mxkqhfhxjydi5ai.jpg",
      "cloudinary_id": "thhu8mxkqhfhxjydi5ai",
      "visible": true,
      "width": 1600,
      "height": 1200,
      "format": "jpg"
    }
  }
]
```

---

### 3. Photos Van Specifiek Album Ophalen

**Endpoint:** `GET /api/albums/:albumId/photos`

Haalt alle photos van een specifiek album op.

**Request:**
```typescript
const albumId = 'd51cff45-b958-4370-a983-51e650ffa43e';
const response = await fetch(`https://dklemailservice.onrender.com/api/albums/${albumId}/photos`);
const photos = await response.json();
```

**Response:**
```typescript
interface PhotoWithAlbumInfo {
  id: string;
  title: string;
  description: string | null;
  url: string;                   // Volledige Cloudinary URL
  thumbnail_url: string | null;  // Thumbnail (200x200)
  cloudinary_id: string;
  visible: boolean;
  width: number | null;
  height: number | null;
  format: string | null;         // "jpg", "png", "webp", etc.
  created_at: string;
  updated_at: string;
  order_number: number | null;   // Volgorde binnen album
}
```

**Voorbeeld Response:**
```json
[
  {
    "id": "08362e92-340a-432a-b306-153ad27ee686",
    "title": "WhatsApp Image 2025-05-17 at 19",
    "description": null,
    "url": "https://res.cloudinary.com/dgfuv7wif/image/upload/v1747513127/thhu8mxkqhfhxjydi5ai.jpg",
    "thumbnail_url": "https://res.cloudinary.com/dgfuv7wif/image/upload/c_thumb,w_200,g_face/v1747513127/thhu8mxkqhfhxjydi5ai.jpg",
    "cloudinary_id": "thhu8mxkqhfhxjydi5ai",
    "visible": true,
    "width": 1600,
    "height": 1200,
    "format": "jpg",
    "created_at": "2025-05-17T20:12:07Z",
    "updated_at": "2025-05-17T20:12:07Z",
    "order_number": 1
  }
]
```

---

## React Voorbeeld

### Albums Grid Met Cover Photos

```typescript
import React, { useEffect, useState } from 'react';

interface AlbumWithCover {
  id: string;
  title: string;
  description: string | null;
  order_number: number;
  cover_photo: {
    url: string;
    thumbnail_url: string | null;
  } | null;
}

const AlbumsGrid: React.FC = () => {
  const [albums, setAlbums] = useState<AlbumWithCover[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchAlbums = async () => {
      try {
        const response = await fetch(
          'https://dklemailservice.onrender.com/api/albums?include_covers=true'
        );
        
        if (!response.ok) {
          throw new Error('Failed to fetch albums');
        }
        
        const data = await response.json();
        setAlbums(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error');
      } finally {
        setLoading(false);
      }
    };

    fetchAlbums();
  }, []);

  if (loading) return <div>Loading albums...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div className="albums-grid">
      {albums.map((album) => (
        <div key={album.id} className="album-card">
          {album.cover_photo && (
            <img
              src={album.cover_photo.thumbnail_url || album.cover_photo.url}
              alt={album.title}
              loading="lazy"
            />
          )}
          <h3>{album.title}</h3>
          <p>{album.description}</p>
        </div>
      ))}
    </div>
  );
};

export default AlbumsGrid;
```

### Album Photo Gallery

```typescript
import React, { useEffect, useState } from 'react';

interface Photo {
  id: string;
  title: string;
  url: string;
  thumbnail_url: string | null;
  order_number: number | null;
}

interface PhotoGalleryProps {
  albumId: string;
}

const PhotoGallery: React.FC<PhotoGalleryProps> = ({ albumId }) => {
  const [photos, setPhotos] = useState<Photo[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchPhotos = async () => {
      try {
        const response = await fetch(
          `https://dklemailservice.onrender.com/api/albums/${albumId}/photos`
        );
        const data = await response.json();
        setPhotos(data);
      } catch (err) {
        console.error('Failed to fetch photos:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchPhotos();
  }, [albumId]);

  if (loading) return <div>Loading photos...</div>;

  return (
    <div className="photo-gallery">
      {photos.map((photo) => (
        <div key={photo.id} className="photo-item">
          <img
            src={photo.thumbnail_url || photo.url}
            alt={photo.title}
            loading="lazy"
            onClick={() => window.open(photo.url, '_blank')}
          />
        </div>
      ))}
    </div>
  );
};

export default PhotoGallery;
```

---

## Next.js App Router Voorbeeld

```typescript
// app/albums/page.tsx
import { Album } from './types';

async function getAlbums(): Promise<Album[]> {
  const response = await fetch(
    'https://dklemailservice.onrender.com/api/albums?include_covers=true',
    { 
      cache: 'no-store' // of 'force-cache' met revalidate
    }
  );
  
  if (!response.ok) {
    throw new Error('Failed to fetch albums');
  }
  
  return response.json();
}

export default async function AlbumsPage() {
  const albums = await getAlbums();

  return (
    <div className="container">
      <h1>Photo Albums</h1>
      <div className="grid">
        {albums.map((album) => (
          <a 
            key={album.id} 
            href={`/albums/${album.id}`}
            className="album-card"
          >
            {album.cover_photo && (
              <img 
                src={album.cover_photo.thumbnail_url || album.cover_photo.url}
                alt={album.title}
              />
            )}
            <h2>{album.title}</h2>
            <p>{album.description}</p>
          </a>
        ))}
      </div>
    </div>
  );
}
```

---

## Belangrijke Opmerkingen

### 🔐 Authenticatie
- **Geen authenticatie nodig** voor het ophalen van zichtbare albums en photos
- Deze zijn publieke endpoints

### 🖼️ Cloudinary Images
- Alle images worden gehost op Cloudinary
- URLs zijn direct bruikbaar in `<img>` tags
- Thumbnail URLs zijn geoptimaliseerd (200x200px)
- Voor custom sizes, gebruik Cloudinary transformaties:
  ```
  https://res.cloudinary.com/dgfuv7wif/image/upload/w_400,h_300,c_fill/[cloudinary_id].jpg
  ```

### ⚡ Performance Tips
1. **Gebruik thumbnails** voor grid views (sneller laden)
2. **Lazy loading** implementeren voor images
3. **Caching** toepassen (albums veranderen niet vaak)
4. **Skeleton screens** tonen tijdens laden

### 🔄 Data Freshness
- Albums worden gesorteerd op `order_number` (1, 2, 3...)
- Photos binnen een album ook gesorteerd op `order_number`
- Production data: **3 albums, 45 totale photos**

### 📊 Huidige Data (Status November 2024)
- **DKL 2025:** 30 photos
- **DKL-2024:** 10 photos
- **Voorbereidingen:** 5 photos

---

## Error Handling

```typescript
async function fetchAlbums() {
  try {
    const response = await fetch('https://dklemailservice.onrender.com/api/albums');
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const albums = await response.json();
    return albums;
  } catch (error) {
    console.error('Error fetching albums:', error);
    // Toon fallback UI of error message
    return [];
  }
}
```

---

## TypeScript Types

```typescript
// types/albums.ts
export interface Album {
  id: string;
  title: string;
  description: string | null;
  cover_photo_id: string | null;
  visible: boolean;
  order_number: number;
  created_at: string;
  updated_at: string;
}

export interface Photo {
  id: string;
  title: string;
  description: string | null;
  url: string;
  thumbnail_url: string | null;
  cloudinary_id: string;
  visible: boolean;
  width: number | null;
  height: number | null;
  format: string | null;
  created_at: string;
  updated_at: string;
}

export interface AlbumWithCover extends Album {
  cover_photo: Photo | null;
}

export interface PhotoWithAlbumInfo extends Photo {
  order_number: number | null;
}
```

---

## Vragen?

Voor vragen of problemen met de API, neem contact op met het backend team of check de API documentatie in [`handlers/album_handler.go`](../handlers/album_handler.go).