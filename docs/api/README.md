# DKL Email Service - Frontend Integration Guide

This directory contains comprehensive documentation and examples for integrating the DKL Email Service image upload functionality into frontend applications.

## 📋 Table of Contents

- [API Documentation](#api-documentation)
- [JavaScript Client Library](#javascript-client-library)
- [Framework Examples](#framework-examples)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)

## 📖 API Documentation

### [Complete API Reference](./image-upload-api.md)
Detailed documentation of all image upload endpoints including:
- Request/response formats
- Authentication requirements
- Error handling
- Rate limiting
- File restrictions

### Key Endpoints
- `POST /api/images/upload` - Single image upload
- `POST /api/images/batch-upload` - Multiple image upload (max 10)
- `GET /api/images/:public_id` - Get image metadata
- `DELETE /api/images/:public_id` - Delete image
- `POST /api/chat/channels/:id/messages` - Send chat image

## 🚀 JavaScript Client Library

### [ImageUploadClient](./image-upload-client.js)
Complete JavaScript client with:
- Progress tracking
- Error handling
- Batch uploads
- Chat integration
- Responsive image generation

### Features
- **Progress Callbacks**: Real-time upload progress
- **Error Recovery**: Automatic retry logic
- **File Validation**: Client-side validation
- **Authentication**: JWT token management
- **Cross-browser**: Works in all modern browsers

### Basic Usage
```javascript
import { ImageUploadClient } from './image-upload-client.js';

const client = new ImageUploadClient({
  apiBaseUrl: '/api',
  authToken: 'your-jwt-token'
});

// Upload single image
const result = await client.uploadImage(file, {
  onProgress: (percent) => console.log(`${percent}% uploaded`)
});

// Upload multiple images
const results = await client.uploadBatchImages([file1, file2], {
  onProgress: (percent) => updateProgressBar(percent)
});
```

## 🎨 Framework Examples

### React Integration
- **[React Components](./react-image-upload-example.jsx)**: Complete React implementation
- Drag & drop interface
- Progress indicators
- Responsive image display
- Error handling

### Vue.js Integration
- **[Vue Components](./vue-image-upload-example.vue)**: Vue 3 Composition API
- Reactive progress tracking
- Gallery management
- Chat integration

### Vanilla JavaScript
- **[HTML Demo](./vanilla-js-example.html)**: Pure JavaScript implementation
- No framework dependencies
- Complete working example
- Copy-paste ready code

## 🏃 Quick Start

### 1. Include the Client Library
```html
<script src="./image-upload-client.js"></script>
```

### 2. Initialize Client
```javascript
const client = new ImageUploadClient({
  apiBaseUrl: '/api',
  authToken: localStorage.getItem('authToken')
});
```

### 3. Upload an Image
```javascript
const fileInput = document.getElementById('file-input');

fileInput.addEventListener('change', async (e) => {
  const file = e.target.files[0];
  if (!file) return;

  try {
    const result = await client.uploadImage(file, {
      onProgress: (percent) => {
        console.log(`Upload progress: ${percent}%`);
      }
    });

    console.log('Upload successful:', result.data);
    // Display the image using result.data.secure_url
  } catch (error) {
    console.error('Upload failed:', error.message);
  }
});
```

## ⚙️ Configuration

### Environment Variables
```env
# Required for image uploads
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret

# Optional configuration
CLOUDINARY_UPLOAD_FOLDER=dkl_images
CLOUDINARY_UPLOAD_PRESET=your_preset
CLOUDINARY_SECURE=true
```

### Client Configuration
```javascript
const client = new ImageUploadClient({
  apiBaseUrl: '/api',           // API base URL
  authToken: 'jwt-token',       // Authentication token
  maxFileSize: 10485760,        // 10MB max file size
  allowedTypes: ['image/jpeg', 'image/png', 'image/gif', 'image/webp'],
  maxBatchSize: 10              // Max files in batch upload
});
```

## 🚨 Error Handling

### Network Errors
```javascript
try {
  const result = await client.uploadImage(file);
} catch (error) {
  if (error.message.includes('Network error')) {
    // Handle network issues
    showRetryButton();
  } else if (error.message.includes('Invalid file type')) {
    // Handle validation errors
    showFileTypeError();
  } else {
    // Handle other errors
    showGenericError(error.message);
  }
}
```

### Authentication Errors
```javascript
// Check for 401 responses
if (error.message.includes('Authentication required')) {
  // Redirect to login or refresh token
  redirectToLogin();
}
```

### File Validation
```javascript
function validateFile(file) {
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

## 📱 Responsive Images

### Automatic Responsive URLs
```javascript
import { ResponsiveImageGenerator } from './image-upload-client.js';

const generator = new ResponsiveImageGenerator('your-cloud-name');

const responsiveData = generator.generateResponsiveUrls('public_id', {
  sizes: [
    { width: 480, dpr: 1 },
    { width: 768, dpr: 1 },
    { width: 1024, dpr: 1 }
  ]
});

// Use in HTML
<img
  src={responsiveData.src}
  srcSet={responsiveData.srcSet}
  sizes="(max-width: 768px) 100vw, 50vw"
  alt="Responsive image"
/>
```

### Custom Transformations
```javascript
// Generate thumbnail
const thumbnailUrl = generator.generateThumbnailUrl('public_id', {
  width: 200,
  height: 200
});

// Custom transformation
const customUrl = generator.generateCustomUrl('public_id', {
  width: 800,
  crop: 'fill',
  gravity: 'face',
  quality: 'auto',
  format: 'auto'
});
```

## 💡 Best Practices

### Performance
- **Lazy Loading**: Use `loading="lazy"` for images below the fold
- **Progressive JPEG**: Enable progressive loading where possible
- **CDN Optimization**: Leverage Cloudinary's global CDN
- **Caching**: Implement proper cache headers

### User Experience
- **Progress Indicators**: Always show upload progress
- **Drag & Drop**: Support drag-and-drop for better UX
- **Preview**: Show image previews before upload
- **Error Recovery**: Provide retry options for failed uploads

### Security
- **File Validation**: Validate files on both client and server
- **Authentication**: Require valid JWT tokens for all uploads
- **CORS**: Configure proper CORS policies
- **Rate Limiting**: Respect API rate limits

### Accessibility
- **Alt Text**: Provide meaningful alt text for images
- **Keyboard Navigation**: Ensure keyboard accessibility
- **Screen Readers**: Test with screen reading software
- **Color Contrast**: Maintain proper contrast ratios

## 🔧 Advanced Usage

### Custom Upload Presets
```javascript
// Use Cloudinary upload presets for advanced processing
const result = await client.uploadImage(file, {
  // Additional parameters can be passed to Cloudinary
  transformation: 'w_800,h_600,c_fill,g_auto'
});
```

### Integration with Form Libraries
```javascript
// React Hook Form integration
import { useForm } from 'react-hook-form';

function ImageUploadForm() {
  const { register, handleSubmit, setValue } = useForm();

  const onSubmit = async (data) => {
    const result = await client.uploadImage(data.image[0]);
    setValue('imageUrl', result.data.secure_url);
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <input {...register('image')} type="file" accept="image/*" />
      <button type="submit">Upload</button>
    </form>
  );
}
```

### Real-time Chat Integration
```javascript
// WebSocket integration for real-time chat
const socket = new WebSocket('ws://localhost:8080/chat');

socket.onmessage = (event) => {
  const message = JSON.parse(event.data);
  if (message.type === 'image') {
    displayChatImage(message);
  }
};

async function sendChatImage(channelId, file, caption) {
  const result = await client.sendChatImage(channelId, file, caption);
  socket.send(JSON.stringify({
    type: 'image',
    message: result.message
  }));
}
```

## 🐛 Troubleshooting

### Common Issues

**401 Unauthorized**
- Check JWT token validity
- Ensure token is properly formatted
- Verify token expiration

**413 Payload Too Large**
- Check file size limits (10MB default)
- Compress images before upload
- Split large batches into smaller chunks

**415 Unsupported Media Type**
- Verify file type is supported
- Check MIME type detection
- Ensure proper file extensions

**429 Too Many Requests**
- Implement exponential backoff
- Respect rate limiting headers
- Add upload queuing for bulk operations

### Debug Mode
```javascript
// Enable debug logging
const client = new ImageUploadClient({
  apiBaseUrl: '/api',
  authToken: token,
  debug: true  // Enable detailed logging
});
```

## 📚 Additional Resources

- [Cloudinary Documentation](https://cloudinary.com/documentation)
- [MDN File API](https://developer.mozilla.org/en-US/docs/Web/API/File)
- [Web APIs - XMLHttpRequest](https://developer.mozilla.org/en-US/docs/Web/API/XMLHttpRequest)
- [Fetch API](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API)

## 🤝 Contributing

When contributing to the frontend integration:
1. Test across multiple browsers
2. Include error handling for all async operations
3. Provide TypeScript definitions if applicable
4. Update documentation for new features
5. Follow the existing code style and patterns

---

**Need Help?** Check the [main API documentation](../api/rest-api.md) or create an issue in the project repository.