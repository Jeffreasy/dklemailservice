# Image Upload API Documentation

This document describes the image upload API endpoints for the DKL Email Service, providing frontend developers with complete integration guidelines.

## Overview

The image upload system supports:
- Single and batch image uploads
- Cloudinary CDN optimization
- Automatic thumbnail generation
- Chat image integration
- User-specific image management

## Authentication

All image upload endpoints require JWT authentication. Include the JWT token in the `Authorization` header:

```
Authorization: Bearer <jwt_token>
```

## API Endpoints

### 1. Single Image Upload

Upload a single image file.

**Endpoint:** `POST /api/images/upload`

**Content-Type:** `multipart/form-data`

**Request Body:**
```javascript
const formData = new FormData();
formData.append('image', file); // File object
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "public_id": "dkl_images/user123_abc123def456",
    "url": "https://res.cloudinary.com/your_cloud/image/upload/v1234567890/dkl_images/user123_abc123def456.jpg",
    "secure_url": "https://res.cloudinary.com/your_cloud/image/upload/v1234567890/dkl_images/user123_abc123def456.jpg",
    "width": 1920,
    "height": 1080,
    "format": "jpg",
    "bytes": 245760,
    "thumbnail_url": "https://res.cloudinary.com/your_cloud/image/upload/w_200,h_200,c_thumb,g_face/v1234567890/dkl_images/user123_abc123def456.jpg"
  }
}
```

**Error Responses:**

**400 Bad Request - Invalid file type:**
```json
{
  "error": "Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed."
}
```

**413 Payload Too Large - File too large:**
```json
{
  "error": "File too large. Maximum size is 10MB."
}
```

**401 Unauthorized - Missing authentication:**
```json
{
  "error": "Authentication required"
}
```

**500 Internal Server Error - Upload failed:**
```json
{
  "error": "Failed to upload image"
}
```

### 2. Batch Image Upload

Upload multiple images in a single request (max 10 images).

**Endpoint:** `POST /api/images/batch-upload`

**Content-Type:** `multipart/form-data`

**Request Body:**
```javascript
const formData = new FormData();
formData.append('images', file1);
formData.append('images', file2);
// ... up to 10 files
```

**Success Response (200):**
```json
{
  "success": true,
  "data": [
    {
      "public_id": "dkl_images/user123_file1_abc123",
      "url": "https://res.cloudinary.com/your_cloud/image/upload/v1234567890/dkl_images/user123_file1_abc123.jpg",
      "secure_url": "https://res.cloudinary.com/your_cloud/image/upload/v1234567890/dkl_images/user123_file1_abc123.jpg",
      "width": 800,
      "height": 600,
      "format": "jpg",
      "bytes": 102400,
      "thumbnail_url": "https://res.cloudinary.com/your_cloud/image/upload/w_200,h_200,c_thumb,g_face/v1234567890/dkl_images/user123_file1_abc123.jpg"
    },
    {
      "public_id": "dkl_images/user123_file2_def456",
      "url": "https://res.cloudinary.com/your_cloud/image/upload/v1234567890/dkl_images/user123_file2_def456.png",
      "secure_url": "https://res.cloudinary.com/your_cloud/image/upload/v1234567890/dkl_images/user123_file2_def456.png",
      "width": 1024,
      "height": 768,
      "format": "png",
      "bytes": 204800,
      "thumbnail_url": "https://res.cloudinary.com/your_cloud/image/upload/w_200,h_200,c_thumb,g_face/v1234567890/dkl_images/user123_file2_def456.png"
    }
  ]
}
```

**Error Response (400) - Invalid number of files:**
```json
{
  "error": "Provide 1-10 image files"
}
```

### 3. Get Image Metadata

Retrieve metadata for a specific image.

**Endpoint:** `GET /api/images/:public_id`

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "public_id": "dkl_images/user123_abc123def456",
    "url": "https://res.cloudinary.com/your_cloud/image/upload/v1234567890/dkl_images/user123_abc123def456.jpg"
  }
}
```

**Error Response (404) - Image not found:**
```json
{
  "error": "Image not found"
}
```

### 4. Delete Image

Delete an image from both Cloudinary and the database.

**Endpoint:** `DELETE /api/images/:public_id`

**Success Response (200):**
```json
{
  "success": true,
  "message": "Image deleted successfully"
}
```

**Error Response (500) - Delete failed:**
```json
{
  "error": "Failed to delete image"
}
```

## Chat Image Integration

### Send Image Message

Send an image message in a chat channel.

**Endpoint:** `POST /api/chat/channels/:channel_id/messages`

**Content-Type:** `multipart/form-data`

**Request Body:**
```javascript
const formData = new FormData();
formData.append('image', imageFile);
formData.append('content', 'Optional caption text'); // Optional
```

**Success Response (200):**
```json
{
  "success": true,
  "message": {
    "id": "uuid",
    "channel_id": "uuid",
    "user_id": "uuid",
    "content": "Optional caption",
    "message_type": "image",
    "file_url": "https://res.cloudinary.com/your_cloud/image/upload/v1234567890/chat_images/user123_abc123.jpg",
    "file_name": "photo.jpg",
    "file_size": 245760,
    "thumbnail_url": "https://res.cloudinary.com/your_cloud/image/upload/w_200,h_200,c_thumb,g_face/v1234567890/chat_images/user123_abc123.jpg",
    "created_at": "2025-01-01T12:00:00Z"
  }
}
```

## File Requirements

- **Supported formats:** JPEG, PNG, GIF, WebP
- **Maximum file size:** 10MB per file
- **Maximum batch size:** 10 files per request
- **Automatic optimization:** Images are optimized for web delivery

## Cloudinary URL Structure

### Base URL Format
```
https://res.cloudinary.com/{cloud_name}/image/upload/{transformations}/{public_id}.{format}
```

### Available Transformations

#### Thumbnail Generation
- `w_200,h_200,c_thumb,g_face`: 200x200 thumbnail with face detection
- `w_400,h_400,c_fill`: 400x400 cropped fill
- `w_800,c_limit`: Max width 800px, maintain aspect ratio

#### Quality Optimization
- `q_auto`: Automatic quality optimization
- `f_auto`: Automatic format conversion (WebP/AVIF when supported)

#### Responsive Images
- `w_480,dpr_auto`: 480px width, automatic device pixel ratio
- `w_768,dpr_auto`: 768px width, automatic device pixel ratio
- `w_1024,dpr_auto`: 1024px width, automatic device pixel ratio

## Error Handling

### Client-Side Validation

```javascript
function validateImageFile(file) {
  const maxSize = 10 * 1024 * 1024; // 10MB
  const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];

  if (!allowedTypes.includes(file.type)) {
    throw new Error('Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed.');
  }

  if (file.size > maxSize) {
    throw new Error('File too large. Maximum size is 10MB.');
  }

  return true;
}
```

### Network Error Handling

```javascript
async function uploadImage(file) {
  try {
    const response = await fetch('/api/images/upload', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authToken}`
      },
      body: formData
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || `HTTP ${response.status}`);
    }

    const result = await response.json();
    return result.data;

  } catch (error) {
    if (error.name === 'TypeError' && error.message.includes('fetch')) {
      throw new Error('Network error. Please check your connection.');
    }
    throw error;
  }
}
```

## Rate Limiting

- **Image uploads:** Limited per user (configurable)
- **Global limits:** Protection against abuse
- **Rate limit headers:** Included in responses when limits are approached

## Security Considerations

- All uploads require valid JWT authentication
- Files are scanned for malware by Cloudinary
- User isolation prevents access to other users' images
- HTTPS-only URLs for secure delivery
- CORS configured for approved frontend domains

## CDN Optimization

Images are automatically optimized by Cloudinary:
- Format conversion (WebP/AVIF for modern browsers)
- Quality optimization based on content
- Global CDN delivery (25,000+ edge locations)
- Automatic compression and caching

## Usage Examples

See the JavaScript client implementation in `docs/api/image-upload-client.js` for complete working examples including progress tracking, error handling, and responsive image generation.