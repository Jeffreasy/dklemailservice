# Frontend Videos API Documentatie

Deze documentatie beschrijft hoe je videos kunt ophalen voor gebruik in de frontend.

## Base URLs

- **Development (Docker):** `http://localhost:8082`
- **Production:** `https://dklemailservice.onrender.com`

## Beschikbare Endpoints

### 1. Alle Zichtbare Videos Ophalen

**Endpoint:** `GET /api/videos`

Haalt alle zichtbare videos op, gesorteerd op `order_number`.

**Request:**
```typescript
const response = await fetch('https://dklemailservice.onrender.com/api/videos');
const videos = await response.json();
```

**Response:**
```typescript
interface Video {
  id: string;              // UUID van de video
  video_id: string;        // Streamable video ID (bijv. "q9ngqu")
  url: string;             // Embed URL (bijv. "https://streamable.com/e/q9ngqu")
  title: string;           // Titel
  description: string | null; // Beschrijving
  thumbnail_url: string | null; // Thumbnail URL (kan null zijn)
  visible: boolean;        // Of video zichtbaar is
  order_number: number;    // Volgorde (1, 2, 3...)
  created_at: string;      // ISO timestamp
  updated_at: string;      // ISO timestamp
}

// Response is een array van Video objecten
Video[]
```

**Voorbeeld Response:**
```json
[
  {
    "id": "14ee164b-50e3-4f59-a3e9-3a54312af9cd",
    "video_id": "q9ngqu",
    "url": "https://streamable.com/e/q9ngqu",
    "title": "De koninklijkeloop!",
    "description": "Preview!!!",
    "thumbnail_url": null,
    "visible": true,
    "order_number": 1,
    "created_at": "2025-03-28T21:45:53Z",
    "updated_at": "2025-04-21T10:49:03.744Z"
  },
  {
    "id": "87502f84-91db-419f-9766-4071ade3e94f",
    "video_id": "tt6k80",
    "url": "https://streamable.com/e/0o2qf9",
    "title": "Highlights Koninklijke Loop 2024 (Wandelevenement Apeldoorn)",
    "description": "Herbeleef de mooiste momenten en de sfeer van De Koninklijke Loop 2024 in deze highlight video.",
    "thumbnail_url": null,
    "visible": true,
    "order_number": 2,
    "created_at": "2024-12-22T20:23:06.654419Z",
    "updated_at": "2024-12-26T23:25:42.699Z"
  }
]
```

---

## Streamable Video Embedding

### Basis Embed Code

De `url` field bevat de embed URL. Deze kan direct gebruikt worden in een iframe:

```html
<iframe 
  src="https://streamable.com/e/q9ngqu" 
  frameborder="0" 
  width="100%" 
  height="360" 
  allowfullscreen
  allow="autoplay"
></iframe>
```

### Video ID Extractie

Als je alleen de video ID hebt (`video_id` field):

```typescript
function getStreamableEmbedUrl(videoId: string): string {
  return `https://streamable.com/e/${videoId}`;
}

// Gebruik:
const embedUrl = getStreamableEmbedUrl('q9ngqu');
// Returns: "https://streamable.com/e/q9ngqu"
```

### Thumbnail Generatie

Streamable thumbnails kun je genereren via hun URL pattern:

```typescript
function getStreamableThumbnail(videoId: string): string {
  return `https://cdn-cf-east.streamable.com/image/${videoId}.jpg`;
}

// Gebruik:
const thumbnail = getStreamableThumbnail('q9ngqu');
// Returns: "https://cdn-cf-east.streamable.com/image/q9ngqu.jpg"
```

---

## React Voorbeelden

### Video Gallery Component

```typescript
import React, { useEffect, useState } from 'react';

interface Video {
  id: string;
  video_id: string;
  url: string;
  title: string;
  description: string | null;
  thumbnail_url: string | null;
  order_number: number;
}

const VideoGallery: React.FC = () => {
  const [videos, setVideos] = useState<Video[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchVideos = async () => {
      try {
        const response = await fetch(
          'https://dklemailservice.onrender.com/api/videos'
        );
        
        if (!response.ok) {
          throw new Error('Failed to fetch videos');
        }
        
        const data = await response.json();
        setVideos(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error');
      } finally {
        setLoading(false);
      }
    };

    fetchVideos();
  }, []);

  const getThumbnail = (videoId: string): string => {
    return `https://cdn-cf-east.streamable.com/image/${videoId}.jpg`;
  };

  if (loading) return <div>Loading videos...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div className="video-gallery">
      <h2>Video Gallery</h2>
      <div className="video-grid">
        {videos.map((video) => (
          <div key={video.id} className="video-card">
            <div className="video-thumbnail">
              <img
                src={video.thumbnail_url || getThumbnail(video.video_id)}
                alt={video.title}
                loading="lazy"
              />
            </div>
            <h3>{video.title}</h3>
            {video.description && <p>{video.description}</p>}
          </div>
        ))}
      </div>
    </div>
  );
};

export default VideoGallery;
```

### Video Player Component

```typescript
import React, { useState } from 'react';

interface VideoPlayerProps {
  video: {
    id: string;
    video_id: string;
    url: string;
    title: string;
    description: string | null;
  };
}

const VideoPlayer: React.FC<VideoPlayerProps> = ({ video }) => {
  const [isPlaying, setIsPlaying] = useState(false);

  const getThumbnail = (videoId: string): string => {
    return `https://cdn-cf-east.streamable.com/image/${videoId}.jpg`;
  };

  return (
    <div className="video-player">
      {!isPlaying ? (
        <div 
          className="video-thumbnail-container"
          onClick={() => setIsPlaying(true)}
        >
          <img
            src={getThumbnail(video.video_id)}
            alt={video.title}
            className="thumbnail"
          />
          <div className="play-button">▶</div>
        </div>
      ) : (
        <iframe
          src={video.url}
          frameBorder="0"
          width="100%"
          height="360"
          allowFullScreen
          allow="autoplay"
          title={video.title}
        />
      )}
      <div className="video-info">
        <h3>{video.title}</h3>
        {video.description && <p>{video.description}</p>}
      </div>
    </div>
  );
};

export default VideoPlayer;
```

### Video Grid With Modal

```typescript
import React, { useEffect, useState } from 'react';

interface Video {
  id: string;
  video_id: string;
  url: string;
  title: string;
  description: string | null;
  order_number: number;
}

const VideoGridWithModal: React.FC = () => {
  const [videos, setVideos] = useState<Video[]>([]);
  const [selectedVideo, setSelectedVideo] = useState<Video | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchVideos = async () => {
      try {
        const response = await fetch(
          'https://dklemailservice.onrender.com/api/videos'
        );
        const data = await response.json();
        setVideos(data);
      } catch (err) {
        console.error('Failed to fetch videos:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchVideos();
  }, []);

  const getThumbnail = (videoId: string): string => {
    return `https://cdn-cf-east.streamable.com/image/${videoId}.jpg`;
  };

  if (loading) return <div>Loading...</div>;

  return (
    <>
      <div className="video-grid">
        {videos.map((video) => (
          <div
            key={video.id}
            className="video-card"
            onClick={() => setSelectedVideo(video)}
          >
            <img
              src={getThumbnail(video.video_id)}
              alt={video.title}
              loading="lazy"
            />
            <h4>{video.title}</h4>
          </div>
        ))}
      </div>

      {selectedVideo && (
        <div className="modal" onClick={() => setSelectedVideo(null)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <button 
              className="close-button"
              onClick={() => setSelectedVideo(null)}
            >
              ×
            </button>
            <iframe
              src={selectedVideo.url}
              frameBorder="0"
              width="100%"
              height="450"
              allowFullScreen
              allow="autoplay"
              title={selectedVideo.title}
            />
            <h2>{selectedVideo.title}</h2>
            {selectedVideo.description && <p>{selectedVideo.description}</p>}
          </div>
        </div>
      )}
    </>
  );
};

export default VideoGridWithModal;
```

---

## Next.js App Router Voorbeeld

```typescript
// app/videos/page.tsx
import { Video } from './types';

async function getVideos(): Promise<Video[]> {
  const response = await fetch(
    'https://dklemailservice.onrender.com/api/videos',
    { 
      cache: 'no-store' // Of use revalidate: 3600 for 1 hour cache
    }
  );
  
  if (!response.ok) {
    throw new Error('Failed to fetch videos');
  }
  
  return response.json();
}

function getThumbnail(videoId: string): string {
  return `https://cdn-cf-east.streamable.com/image/${videoId}.jpg`;
}

export default async function VideosPage() {
  const videos = await getVideos();

  return (
    <div className="container">
      <h1>Video Gallery</h1>
      <div className="grid">
        {videos.map((video) => (
          <a 
            key={video.id} 
            href={`/videos/${video.id}`}
            className="video-card"
          >
            <img 
              src={getThumbnail(video.video_id)}
              alt={video.title}
            />
            <h2>{video.title}</h2>
            <p>{video.description}</p>
          </a>
        ))}
      </div>
    </div>
  );
}
```

```typescript
// app/videos/[id]/page.tsx
import { Video } from '../types';

async function getVideo(id: string): Promise<Video> {
  const response = await fetch(
    `https://dklemailservice.onrender.com/api/videos`,
    { cache: 'no-store' }
  );
  
  const videos: Video[] = await response.json();
  const video = videos.find(v => v.id === id);
  
  if (!video) {
    throw new Error('Video not found');
  }
  
  return video;
}

export default async function VideoPage({ 
  params 
}: { 
  params: { id: string } 
}) {
  const video = await getVideo(params.id);

  return (
    <div className="container">
      <div className="video-player">
        <iframe
          src={video.url}
          frameBorder="0"
          width="100%"
          height="500"
          allowFullScreen
          allow="autoplay"
          title={video.title}
        />
      </div>
      <h1>{video.title}</h1>
      {video.description && <p>{video.description}</p>}
    </div>
  );
}
```

---

## Belangrijke Opmerkingen

### 🎥 Video Platform
- **Alle videos worden gehost op Streamable**
- URLs zijn direct te embedden in iframes
- Geen API key nodig voor publieke videos
- Autoplay wordt ondersteund

### 🖼️ Thumbnails
- `thumbnail_url` kan `null` zijn in de database
- Fallback: gebruik Streamable's thumbnail URL pattern:
  ```
  https://cdn-cf-east.streamable.com/image/{video_id}.jpg
  ```
- Thumbnails zijn 1280x720 pixels (16:9 aspect ratio)

### ⚡ Performance Tips
1. **Lazy load thumbnails** in grid views
2. **Click-to-play** pattern (thumbnail eerst, dan iframe)
3. **Modal of lightbox** voor betere UX
4. **Responsive iframes** met aspect ratio padding trick:
   ```css
   .video-container {
     position: relative;
     padding-bottom: 56.25%; /* 16:9 aspect ratio */
     height: 0;
     overflow: hidden;
   }
   
   .video-container iframe {
     position: absolute;
     top: 0;
     left: 0;
     width: 100%;
     height: 100%;
   }
   ```

### 🔐 Authenticatie
- **Geen authenticatie nodig** voor publieke videos
- Dit is een public endpoint

### 📊 Huidige Data (Status November 2024)
- **Totaal videos:** 5
- Alle videos zijn **Streamable embeds**
- Videos zijn gesorteerd op `order_number`

### 🎬 Video Lijst (Production)
1. **De koninklijkeloop!** - Preview
2. **Highlights Koninklijke Loop 2024** - Wandelevenement
3. **De spannende start** - Start van deelnemers
4. **Promotie: Flyers verspreiden** - Vrijwilligers in actie
5. **Hoofdevenement** - Sfeerimpressie met muziek

---

## CSS Voorbeelden

### Responsive Video Container

```css
.video-container {
  position: relative;
  padding-bottom: 56.25%; /* 16:9 Aspect Ratio */
  height: 0;
  overflow: hidden;
  max-width: 100%;
  background: #000;
}

.video-container iframe {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border: 0;
}
```

### Video Grid Layout

```css
.video-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
  padding: 2rem;
}

.video-card {
  cursor: pointer;
  transition: transform 0.2s;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.video-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.2);
}

.video-card img {
  width: 100%;
  height: auto;
  display: block;
}

.video-card h4 {
  padding: 1rem;
  margin: 0;
}
```

### Play Button Overlay

```css
.video-thumbnail-container {
  position: relative;
  cursor: pointer;
}

.video-thumbnail-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  transition: background 0.3s;
}

.video-thumbnail-container:hover::before {
  background: rgba(0, 0, 0, 0.5);
}

.play-button {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 4rem;
  color: white;
  text-shadow: 0 2px 4px rgba(0,0,0,0.5);
  transition: transform 0.2s;
}

.video-thumbnail-container:hover .play-button {
  transform: translate(-50%, -50%) scale(1.1);
}
```

---

## TypeScript Types

```typescript
// types/videos.ts
export interface Video {
  id: string;
  video_id: string;
  url: string;
  title: string;
  description: string | null;
  thumbnail_url: string | null;
  visible: boolean;
  order_number: number;
  created_at: string;
  updated_at: string;
}

export interface VideoPlayerProps {
  video: Video;
  autoplay?: boolean;
  controls?: boolean;
}

export interface VideoGridProps {
  videos: Video[];
  onVideoClick?: (video: Video) => void;
}
```

---

## Error Handling

```typescript
async function fetchVideos() {
  try {
    const response = await fetch(
      'https://dklemailservice.onrender.com/api/videos'
    );
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const videos = await response.json();
    return videos;
  } catch (error) {
    console.error('Error fetching videos:', error);
    // Show fallback UI or error message
    return [];
  }
}
```

---

## Streamable API Reference

Voor meer informatie over Streamable embeds, zie:
- [Streamable Embed Documentation](https://support.streamable.com/hc/en-us/articles/360027932033-Embedding-videos)

---

## Vragen?

Voor vragen of problemen met de API, neem contact op met het backend team of check de API documentatie in [`handlers/video_handler.go`](../handlers/video_handler.go).