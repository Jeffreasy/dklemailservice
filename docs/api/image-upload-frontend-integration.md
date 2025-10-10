# Frontend Image Upload Integration Guide

Complete guide for integrating Cloudinary image upload functionality into your frontend application.

## Table of Contents

- [Quick Start](#quick-start)
- [API Overview](#api-overview)
- [Authentication](#authentication)
- [JavaScript Client Setup](#javascript-client-setup)
- [React Integration](#react-integration)
- [Vue.js Integration](#vuejs-integration)
- [Vanilla JavaScript Examples](#vanilla-javascript-examples)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)
- [Complete Examples](#complete-examples)

## Quick Start

### 1. Include the JavaScript Client

```html
<script src="/path/to/image-upload-client.js"></script>
```

### 2. Initialize and Use

```javascript
// Initialize client
const client = new ImageUploadClient({
  apiBaseUrl: '/api',
  authToken: 'your-jwt-token'
});

// Upload single image
const result = await client.uploadImage(file, {
  onProgress: (percent) => console.log(`${percent}%`)
});

// Upload multiple images
const results = await client.uploadBatchImages(files, {
  onProgress: (percent) => console.log(`${percent}%`)
});
```

## API Overview

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/images/upload` | Upload single image |
| `POST` | `/api/images/batch-upload` | Upload multiple images |
| `POST` | `/api/images/batch-upload?mode=sequential` | Sequential batch upload |
| `GET` | `/api/images/{public_id}` | Get image metadata |
| `DELETE` | `/api/images/{public_id}` | Delete image |
| `POST` | `/api/chat/channels/{id}/messages` | Send chat image message |

### Request/Response Formats

#### Single Image Upload

**Request:**
```javascript
const formData = new FormData();
formData.append('image', file);

fetch('/api/images/upload', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`
  },
  body: formData
});
```

**Response:**
```json
{
  "success": true,
  "data": {
    "public_id": "dkl_images/user123_abc123def.jpg",
    "url": "https://res.cloudinary.com/.../image.jpg",
    "secure_url": "https://res.cloudinary.com/.../image.jpg",
    "width": 1920,
    "height": 1080,
    "format": "jpg",
    "bytes": 245760
  }
}
```

#### Batch Image Upload

**Request:**
```javascript
const formData = new FormData();
files.forEach(file => formData.append('images', file));

fetch('/api/images/batch-upload?mode=parallel', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`
  },
  body: formData
});
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "public_id": "dkl_images/user123_file1.jpg",
      "url": "https://res.cloudinary.com/.../file1.jpg",
      "secure_url": "https://res.cloudinary.com/.../file1.jpg",
      "width": 800,
      "height": 600,
      "format": "jpg",
      "bytes": 102400
    }
  ],
  "uploaded_count": 3,
  "total_files": 3,
  "mode": "parallel"
}
```

## Authentication

All image upload endpoints require JWT authentication:

```javascript
const headers = {
  'Authorization': `Bearer ${jwtToken}`
};
```

### Token Management

```javascript
// Store token after login
localStorage.setItem('authToken', token);

// Retrieve token for API calls
const token = localStorage.getItem('authToken');

// Include in requests
const headers = {
  'Authorization': `Bearer ${token}`
};
```

## JavaScript Client Setup

### Installation

1. **Download the client:**
   ```bash
   curl -o image-upload-client.js https://your-domain.com/docs/api/image-upload-client.js
   ```

2. **Include in your HTML:**
   ```html
   <script src="image-upload-client.js"></script>
   ```

3. **Or use as ES6 module:**
   ```javascript
   import { ImageUploadClient, ResponsiveImageGenerator } from './image-upload-client.js';
   ```

### Client Configuration

```javascript
const client = new ImageUploadClient({
  // Required
  authToken: 'your-jwt-token',

  // Optional
  apiBaseUrl: '/api',                    // Default: '/api'
  maxFileSize: 10 * 1024 * 1024,        // Default: 10MB
  maxBatchSize: 10,                     // Default: 10 files

  // Allowed file types
  allowedTypes: [
    'image/jpeg',
    'image/png',
    'image/gif',
    'image/webp'
  ]
});
```

### Client Methods

#### Single Image Upload

```javascript
async uploadImage(file, options = {})
```

**Parameters:**
- `file`: File object from input or drag-drop
- `options`: Configuration object

**Options:**
- `onProgress(percent)`: Progress callback (0-100)
- `onSuccess(result)`: Success callback
- `onError(error)`: Error callback

**Example:**
```javascript
try {
  const result = await client.uploadImage(file, {
    onProgress: (percent) => {
      progressBar.style.width = `${percent}%`;
      progressText.textContent = `${percent}%`;
    }
  });

  console.log('Upload successful:', result.data.secure_url);
} catch (error) {
  console.error('Upload failed:', error.message);
}
```

#### Batch Image Upload

```javascript
async uploadBatchImages(files, options = {})
async uploadBatchImagesSequential(files, options = {})
```

**Parameters:**
- `files`: Array of File objects
- `options`: Configuration object

**Options:**
- `onProgress(percent)`: Overall progress (0-100)
- `onFileProgress({fileIndex, fileName, progress})`: Per-file progress (sequential mode)
- `onBatchProgress({completed, total, currentFile})`: Batch completion progress
- `onSuccess(result)`: Success callback
- `onError(error)`: Error callback

**Examples:**

```javascript
// Parallel upload (default - faster)
const result = await client.uploadBatchImages(files, {
  onProgress: (percent) => console.log(`Overall: ${percent}%`),
  onBatchProgress: ({completed, total}) =>
    console.log(`Completed ${completed}/${total} files`)
});

// Sequential upload (more reliable for large files)
const result = await client.uploadBatchImagesSequential(files, {
  onFileProgress: ({fileIndex, fileName, progress}) =>
    console.log(`File ${fileName}: ${progress}%`),
  onBatchProgress: ({completed, total, currentFile}) =>
    console.log(`Completed ${completed}/${total} files`)
});
```

#### Chat Image Messages

```javascript
async sendChatImage(channelId, file, caption = '', options = {})
```

**Example:**
```javascript
const result = await client.sendChatImage('channel-123', file, 'Check this out!', {
  onProgress: (percent) => console.log(`${percent}%`)
});
```

#### Image Management

```javascript
// Get image metadata
const metadata = await client.getImageMetadata('public_id_123');

// Delete image
await client.deleteImage('public_id_123');
```

## React Integration

### Hook Usage

```javascript
import { useImageUpload } from './image-upload-client.js';

function ImageUploader() {
  const { uploadImage, uploadBatch, uploading, progress, error } = useImageUpload('/api', token);

  const handleSingleUpload = async (file) => {
    try {
      const result = await uploadImage(file);
      console.log('Uploaded:', result.data.secure_url);
    } catch (err) {
      console.error('Upload failed:', err);
    }
  };

  const handleBatchUpload = async (files, mode = 'parallel') => {
    try {
      const result = await uploadBatch(files, mode);
      console.log(`Uploaded ${result.uploaded_count}/${result.total_files} files`);
    } catch (err) {
      console.error('Batch upload failed:', err);
    }
  };

  return (
    <div>
      <input
        type="file"
        onChange={(e) => handleSingleUpload(e.target.files[0])}
        disabled={uploading}
      />

      <input
        type="file"
        multiple
        onChange={(e) => handleBatchUpload(Array.from(e.target.files))}
        disabled={uploading}
      />

      {uploading && (
        <div>
          <progress value={progress} max="100" />
          <span>{progress}%</span>
        </div>
      )}

      {error && <div className="error">{error}</div>}
    </div>
  );
}
```

### Component Example

```javascript
import React, { useState, useCallback } from 'react';
import { ImageUploadClient } from './image-upload-client.js';

function ImageGallery() {
  const [images, setImages] = useState([]);
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);

  const client = useMemo(() => new ImageUploadClient({
    apiBaseUrl: '/api',
    authToken: localStorage.getItem('authToken')
  }), []);

  const uploadImages = useCallback(async (files) => {
    setUploading(true);
    setProgress(0);

    try {
      // Auto-detect mode based on file sizes
      const totalSize = files.reduce((sum, file) => sum + file.size, 0);
      const avgSize = totalSize / files.length;
      const mode = avgSize > 2 * 1024 * 1024 ? 'sequential' : 'parallel'; // 2MB threshold

      const result = await client.uploadBatchImages(files, {
        onProgress: setProgress,
        onBatchProgress: ({completed, total}) => {
          console.log(`Completed ${completed}/${total} files`);
        }
      });

      setImages(prev => [...prev, ...result.data]);
    } catch (error) {
      console.error('Upload failed:', error);
    } finally {
      setUploading(false);
      setProgress(0);
    }
  }, [client]);

  return (
    <div className="image-gallery">
      <div className="upload-zone">
        <input
          type="file"
          multiple
          accept="image/*"
          onChange={(e) => uploadImages(Array.from(e.target.files))}
          disabled={uploading}
        />
        {uploading && (
          <div className="progress">
            <div className="progress-bar" style={{width: `${progress}%`}} />
            <span>{progress}%</span>
          </div>
        )}
      </div>

      <div className="image-grid">
        {images.map((image, index) => (
          <img
            key={index}
            src={image.secure_url}
            alt={image.filename}
            loading="lazy"
          />
        ))}
      </div>
    </div>
  );
}
```

## Vue.js Integration

### Composable Usage

```javascript
import { useImageUpload } from './image-upload-client.js';

export default {
  setup() {
    const {
      uploadImage,
      uploadBatch,
      uploading,
      progress,
      error
    } = useImageUpload('/api', authToken);

    const handleUpload = async (files) => {
      try {
        const result = await uploadBatch(files, 'parallel');
        console.log('Upload successful:', result);
      } catch (err) {
        console.error('Upload failed:', err);
      }
    };

    return {
      uploadImage,
      uploadBatch,
      uploading,
      progress,
      error,
      handleUpload
    };
  }
};
```

### Component Example

```vue
<template>
  <div class="image-uploader">
    <div class="upload-controls">
      <input
        ref="fileInput"
        type="file"
        multiple
        accept="image/*"
        @change="handleFileSelect"
        :disabled="uploading"
      />

      <button @click="triggerFileSelect" :disabled="uploading">
        {{ uploading ? 'Uploading...' : 'Select Images' }}
      </button>
    </div>

    <div v-if="uploading" class="progress-container">
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: progress + '%' }" />
      </div>
      <span class="progress-text">{{ progress }}%</span>
    </div>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <div class="image-preview">
      <div
        v-for="(image, index) in uploadedImages"
        :key="index"
        class="image-item"
      >
        <img :src="image.secure_url" :alt="image.filename" />
        <button @click="deleteImage(image.public_id)">×</button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue';
import { ImageUploadClient } from './image-upload-client.js';

export default {
  name: 'ImageUploader',
  setup() {
    const fileInput = ref(null);
    const uploadedImages = ref([]);

    const client = new ImageUploadClient({
      apiBaseUrl: '/api',
      authToken: localStorage.getItem('authToken')
    });

    const uploadState = reactive({
      uploading: false,
      progress: 0,
      error: null
    });

    const triggerFileSelect = () => {
      fileInput.value.click();
    };

    const handleFileSelect = async (event) => {
      const files = Array.from(event.target.files);
      if (files.length === 0) return;

      uploadState.uploading = true;
      uploadState.progress = 0;
      uploadState.error = null;

      try {
        // Choose mode based on file characteristics
        const totalSize = files.reduce((sum, file) => sum + file.size, 0);
        const mode = totalSize > 10 * 1024 * 1024 ? 'sequential' : 'parallel'; // 10MB threshold

        const result = await client.uploadBatchImages(files, {
          onProgress: (percent) => uploadState.progress = percent,
          onBatchProgress: ({completed, total}) => {
            console.log(`Completed ${completed}/${total} files`);
          }
        });

        uploadedImages.value.push(...result.data);
      } catch (error) {
        uploadState.error = error.message;
      } finally {
        uploadState.uploading = false;
        uploadState.progress = 0;
        event.target.value = ''; // Reset input
      }
    };

    const deleteImage = async (publicId) => {
      try {
        await client.deleteImage(publicId);
        uploadedImages.value = uploadedImages.value.filter(
          img => img.public_id !== publicId
        );
      } catch (error) {
        console.error('Delete failed:', error);
      }
    };

    return {
      fileInput,
      uploadedImages,
      uploadState,
      triggerFileSelect,
      handleFileSelect,
      deleteImage
    };
  }
};
</script>

<style scoped>
.upload-controls {
  margin-bottom: 20px;
}

.progress-container {
  margin: 10px 0;
}

.progress-bar {
  width: 100%;
  height: 20px;
  background-color: #f0f0f0;
  border-radius: 10px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background-color: #4CAF50;
  transition: width 0.3s ease;
}

.error-message {
  color: #f44336;
  margin: 10px 0;
}

.image-preview {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 10px;
  margin-top: 20px;
}

.image-item {
  position: relative;
}

.image-item img {
  width: 100%;
  height: 150px;
  object-fit: cover;
  border-radius: 4px;
}

.image-item button {
  position: absolute;
  top: 5px;
  right: 5px;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  border: none;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  cursor: pointer;
}
</style>
```

## Vanilla JavaScript Examples

### Basic Single File Upload

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Image Upload</title>
    <style>
        .upload-container { max-width: 500px; margin: 50px auto; }
        .progress-bar { width: 100%; height: 20px; background: #f0f0f0; border-radius: 10px; overflow: hidden; margin: 10px 0; }
        .progress-fill { height: 100%; background: #4CAF50; transition: width 0.3s; }
        .image-preview { margin-top: 20px; }
        .image-preview img { max-width: 100%; border-radius: 4px; }
    </style>
</head>
<body>
    <div class="upload-container">
        <h2>Upload Image</h2>

        <input type="file" id="fileInput" accept="image/*" />
        <button id="uploadBtn">Upload</button>

        <div id="progressContainer" style="display: none;">
            <div class="progress-bar">
                <div class="progress-fill" id="progressFill" style="width: 0%"></div>
            </div>
            <span id="progressText">0%</span>
        </div>

        <div id="result"></div>
        <div id="imagePreview" class="image-preview"></div>
    </div>

    <script src="image-upload-client.js"></script>
    <script>
        // Initialize client
        const client = new ImageUploadClient({
            apiBaseUrl: '/api',
            authToken: localStorage.getItem('authToken')
        });

        const fileInput = document.getElementById('fileInput');
        const uploadBtn = document.getElementById('uploadBtn');
        const progressContainer = document.getElementById('progressContainer');
        const progressFill = document.getElementById('progressFill');
        const progressText = document.getElementById('progressText');
        const result = document.getElementById('result');
        const imagePreview = document.getElementById('imagePreview');

        uploadBtn.addEventListener('click', async () => {
            const file = fileInput.files[0];
            if (!file) {
                result.textContent = 'Please select a file';
                return;
            }

            try {
                result.textContent = '';
                progressContainer.style.display = 'block';
                progressFill.style.width = '0%';
                progressText.textContent = '0%';
                uploadBtn.disabled = true;

                const uploadResult = await client.uploadImage(file, {
                    onProgress: (percent) => {
                        progressFill.style.width = `${percent}%`;
                        progressText.textContent = `${percent}%`;
                    }
                });

                result.innerHTML = `<div style="color: green;">Upload successful!</div>`;
                imagePreview.innerHTML = `
                    <img src="${uploadResult.data.secure_url}" alt="Uploaded image" />
                    <p>URL: ${uploadResult.data.secure_url}</p>
                    <p>Size: ${uploadResult.data.width}x${uploadResult.data.height}</p>
                `;

            } catch (error) {
                result.innerHTML = `<div style="color: red;">Upload failed: ${error.message}</div>`;
            } finally {
                progressContainer.style.display = 'none';
                uploadBtn.disabled = false;
            }
        });
    </script>
</body>
</html>
```

### Advanced Batch Upload with Drag & Drop

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Batch Image Upload</title>
    <style>
        .upload-zone {
            border: 2px dashed #ccc;
            border-radius: 8px;
            padding: 40px;
            text-align: center;
            margin: 20px 0;
            transition: border-color 0.3s;
        }
        .upload-zone.dragover { border-color: #4CAF50; background: #f9fff9; }
        .upload-zone.uploading { border-color: #2196F3; background: #f3f9ff; }
        .file-list { margin: 20px 0; }
        .file-item { display: flex; justify-content: space-between; padding: 8px; border: 1px solid #ddd; margin: 4px 0; border-radius: 4px; }
        .file-progress { flex: 1; margin: 0 10px; }
        .progress-bar { height: 8px; background: #f0f0f0; border-radius: 4px; overflow: hidden; }
        .progress-fill { height: 100%; background: #4CAF50; transition: width 0.3s; }
        .image-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(150px, 1fr)); gap: 10px; margin-top: 20px; }
        .image-item { position: relative; }
        .image-item img { width: 100%; height: 150px; object-fit: cover; border-radius: 4px; }
        .delete-btn { position: absolute; top: 5px; right: 5px; background: rgba(255,0,0,0.8); color: white; border: none; border-radius: 50%; width: 24px; height: 24px; cursor: pointer; }
    </style>
</head>
<body>
    <div class="upload-container" style="max-width: 800px; margin: 50px auto;">
        <h2>Batch Image Upload</h2>

        <div class="upload-zone" id="uploadZone">
            <p>Drag & drop images here or click to select</p>
            <input type="file" id="fileInput" multiple accept="image/*" style="display: none;" />
            <button id="selectBtn">Select Files</button>
        </div>

        <div class="controls">
            <select id="uploadMode">
                <option value="parallel">Parallel Upload (Fast)</option>
                <option value="sequential">Sequential Upload (Reliable)</option>
            </select>
            <button id="uploadBtn" disabled>Upload Images</button>
            <button id="clearBtn">Clear</button>
        </div>

        <div id="fileList" class="file-list"></div>

        <div id="overallProgress" style="display: none;">
            <div class="progress-bar">
                <div class="progress-fill" id="overallProgressFill" style="width: 0%"></div>
            </div>
            <span id="overallProgressText">0%</span>
        </div>

        <div id="result"></div>

        <div id="imageGrid" class="image-grid"></div>
    </div>

    <script src="image-upload-client.js"></script>
    <script>
        // Initialize client
        const client = new ImageUploadClient({
            apiBaseUrl: '/api',
            authToken: localStorage.getItem('authToken')
        });

        // DOM elements
        const uploadZone = document.getElementById('uploadZone');
        const fileInput = document.getElementById('fileInput');
        const selectBtn = document.getElementById('selectBtn');
        const uploadBtn = document.getElementById('uploadBtn');
        const clearBtn = document.getElementById('clearBtn');
        const uploadMode = document.getElementById('uploadMode');
        const fileList = document.getElementById('fileList');
        const overallProgress = document.getElementById('overallProgress');
        const overallProgressFill = document.getElementById('overallProgressFill');
        const overallProgressText = document.getElementById('overallProgressText');
        const result = document.getElementById('result');
        const imageGrid = document.getElementById('imageGrid');

        let selectedFiles = [];

        // Drag & drop functionality
        uploadZone.addEventListener('dragover', (e) => {
            e.preventDefault();
            uploadZone.classList.add('dragover');
        });

        uploadZone.addEventListener('dragleave', () => {
            uploadZone.classList.remove('dragover');
        });

        uploadZone.addEventListener('drop', (e) => {
            e.preventDefault();
            uploadZone.classList.remove('dragover');

            const files = Array.from(e.dataTransfer.files).filter(file =>
                file.type.startsWith('image/')
            );
            handleFiles(files);
        });

        // File selection
        selectBtn.addEventListener('click', () => fileInput.click());
        fileInput.addEventListener('change', (e) => handleFiles(Array.from(e.target.files)));

        function handleFiles(files) {
            selectedFiles = files;
            updateFileList();
            uploadBtn.disabled = files.length === 0;
        }

        function updateFileList() {
            fileList.innerHTML = '';
            selectedFiles.forEach((file, index) => {
                const item = document.createElement('div');
                item.className = 'file-item';
                item.innerHTML = `
                    <span>${file.name}</span>
                    <span>${(file.size / 1024 / 1024).toFixed(2)} MB</span>
                    <div class="file-progress">
                        <div class="progress-bar">
                            <div class="progress-fill" style="width: 0%"></div>
                        </div>
                    </div>
                    <button onclick="removeFile(${index})">×</button>
                `;
                fileList.appendChild(item);
            });
        }

        function removeFile(index) {
            selectedFiles.splice(index, 1);
            updateFileList();
            uploadBtn.disabled = selectedFiles.length === 0;
        }

        // Upload functionality
        uploadBtn.addEventListener('click', async () => {
            if (selectedFiles.length === 0) return;

            const mode = uploadMode.value;
            const uploadMethod = mode === 'sequential' ? 'uploadBatchImagesSequential' : 'uploadBatchImages';

            try {
                result.textContent = '';
                overallProgress.style.display = 'block';
                overallProgressFill.style.width = '0%';
                overallProgressText.textContent = '0%';
                uploadZone.classList.add('uploading');
                uploadBtn.disabled = true;

                const uploadResult = await client[uploadMethod](selectedFiles, {
                    onProgress: (percent) => {
                        overallProgressFill.style.width = `${percent}%`;
                        overallProgressText.textContent = `${percent}%`;
                    },
                    onFileProgress: ({fileIndex, fileName, progress}) => {
                        updateFileProgress(fileIndex, progress);
                    },
                    onBatchProgress: ({completed, total, currentFile}) => {
                        console.log(`Completed ${completed}/${total} files`);
                        if (currentFile) {
                            console.log(`Currently uploading: ${currentFile.name}`);
                        }
                    }
                });

                result.innerHTML = `
                    <div style="color: green; margin: 10px 0;">
                        Upload successful! ${uploadResult.uploaded_count}/${uploadResult.total_files} files uploaded.
                    </div>
                `;

                // Display uploaded images
                uploadResult.data.forEach(image => {
                    const imgDiv = document.createElement('div');
                    imgDiv.className = 'image-item';
                    imgDiv.innerHTML = `
                        <img src="${image.secure_url}" alt="${image.filename}" loading="lazy" />
                        <button class="delete-btn" onclick="deleteImage('${image.public_id}')">×</button>
                    `;
                    imageGrid.appendChild(imgDiv);
                });

                // Clear file list
                selectedFiles = [];
                updateFileList();

            } catch (error) {
                result.innerHTML = `<div style="color: red; margin: 10px 0;">Upload failed: ${error.message}</div>`;
            } finally {
                overallProgress.style.display = 'none';
                uploadZone.classList.remove('uploading');
                uploadBtn.disabled = false;
            }
        });

        function updateFileProgress(fileIndex, progress) {
            const progressBars = fileList.querySelectorAll('.progress-fill');
            if (progressBars[fileIndex]) {
                progressBars[fileIndex].style.width = `${progress}%`;
            }
        }

        // Delete image
        async function deleteImage(publicId) {
            try {
                await client.deleteImage(publicId);
                // Remove from grid
                const imageItems = imageGrid.querySelectorAll('.image-item');
                imageItems.forEach(item => {
                    const deleteBtn = item.querySelector('.delete-btn');
                    if (deleteBtn && deleteBtn.onclick.toString().includes(publicId)) {
                        item.remove();
                    }
                });
            } catch (error) {
                console.error('Delete failed:', error);
            }
        }

        // Clear functionality
        clearBtn.addEventListener('click', () => {
            selectedFiles = [];
            updateFileList();
            imageGrid.innerHTML = '';
            result.textContent = '';
            uploadBtn.disabled = true;
        });
    </script>
</body>
</html>
```

## Error Handling

### Client-Side Validation

```javascript
function validateFiles(files) {
    const errors = [];

    files.forEach((file, index) => {
        // Check file type
        if (!client.allowedTypes.includes(file.type)) {
            errors.push(`File ${index + 1}: Invalid file type`);
        }

        // Check file size
        if (file.size > client.maxFileSize) {
            errors.push(`File ${index + 1}: File too large (max ${client.maxFileSize / 1024 / 1024}MB)`);
        }

        // Check file name
        if (file.name.length > 255) {
            errors.push(`File ${index + 1}: Filename too long`);
        }
    });

    return errors;
}
```

### Network Error Handling

```javascript
async function uploadWithRetry(file, maxRetries = 3) {
    let lastError;

    for (let attempt = 1; attempt <= maxRetries; attempt++) {
        try {
            return await client.uploadImage(file, {
                onProgress: (percent) => console.log(`Attempt ${attempt}: ${percent}%`)
            });
        } catch (error) {
            lastError = error;

            if (attempt < maxRetries) {
                // Exponential backoff
                const delay = Math.pow(2, attempt) * 1000;
                console.log(`Attempt ${attempt} failed, retrying in ${delay}ms...`);
                await new Promise(resolve => setTimeout(resolve, delay));
            }
        }
    }

    throw lastError;
}
```

### Batch Upload Error Handling

```javascript
async function uploadBatchWithFallback(files) {
    try {
        // Try parallel upload first
        return await client.uploadBatchImages(files, {
            onProgress: (percent) => console.log(`Parallel: ${percent}%`)
        });
    } catch (error) {
        console.warn('Parallel upload failed, trying sequential mode:', error.message);

        try {
            // Fallback to sequential upload
            return await client.uploadBatchImagesSequential(files, {
                onProgress: (percent) => console.log(`Sequential: ${percent}%`),
                onFileProgress: ({fileIndex, fileName, progress}) =>
                    console.log(`File ${fileName}: ${progress}%`)
            });
        } catch (sequentialError) {
            console.error('Both upload modes failed:', sequentialError.message);
            throw sequentialError;
        }
    }
}
```

## Best Practices

### 1. File Validation

```javascript
function validateAndPrepareFiles(files) {
    return files.filter(file => {
        // Client-side validation
        if (!file.type.startsWith('image/')) {
            console.warn(`Skipping non-image file: ${file.name}`);
            return false;
        }

        if (file.size > 10 * 1024 * 1024) { // 10MB
            console.warn(`Skipping large file: ${file.name}`);
            return false;
        }

        return true;
    });
}
```

### 2. Progress Feedback

```javascript
function createProgressUI() {
    const container = document.createElement('div');
    container.innerHTML = `
        <div class="upload-progress">
            <div class="progress-bar">
                <div class="progress-fill" style="width: 0%"></div>
            </div>
            <div class="progress-text">0%</div>
            <div class="file-status"></div>
        </div>
    `;
    return container;
}

function updateProgress(ui, percent, status = '') {
    const fill = ui.querySelector('.progress-fill');
    const text = ui.querySelector('.progress-text');
    const statusDiv = ui.querySelector('.file-status');

    fill.style.width = `${percent}%`;
    text.textContent = `${percent}%`;
    if (status) statusDiv.textContent = status;
}
```

### 3. Memory Management

```javascript
function cleanupUploads() {
    // Clear file inputs
    const inputs = document.querySelectorAll('input[type="file"]');
    inputs.forEach(input => input.value = '');

    // Clear preview images
    const previews = document.querySelectorAll('.image-preview img');
    previews.forEach(img => {
        URL.revokeObjectURL(img.src); // Free memory
    });
}
```

### 4. Responsive Images

```javascript
// Generate responsive image URLs
function createResponsiveImage(src, publicId) {
    const generator = new ResponsiveImageGenerator('your-cloud-name');

    const { src: responsiveSrc, srcSet, sizes } = generator.generateResponsiveUrls(publicId);

    const img = document.createElement('img');
    img.src = responsiveSrc;
    img.srcset = srcSet;
    img.sizes = sizes;
    img.loading = 'lazy';

    return img;
}
```

### 5. Upload Queue Management

```javascript
class UploadQueue {
    constructor(client, maxConcurrent = 3) {
        this.client = client;
        this.maxConcurrent = maxConcurrent;
        this.queue = [];
        this.active = 0;
        this.results = [];
    }

    add(files) {
        this.queue.push(...files);
        this.process();
    }

    async process() {
        while (this.active < this.maxConcurrent && this.queue.length > 0) {
            const file = this.queue.shift();
            this.active++;

            try {
                const result = await this.client.uploadImage(file);
                this.results.push(result);
            } catch (error) {
                console.error(`Upload failed for ${file.name}:`, error);
            } finally {
                this.active--;
                this.process(); // Process next item
            }
        }
    }

    getResults() {
        return this.results;
    }
}
```

## Complete Examples

### React Image Gallery Component

```javascript
import React, { useState, useCallback, useRef } from 'react';
import { ImageUploadClient, ResponsiveImageGenerator } from './image-upload-client.js';

const ImageGallery = () => {
  const [images, setImages] = useState([]);
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [error, setError] = useState(null);
  const fileInputRef = useRef();

  const client = new ImageUploadClient({
    apiBaseUrl: '/api',
    authToken: localStorage.getItem('authToken')
  });

  const generator = new ResponsiveImageGenerator('your-cloud-name');

  const handleFileSelect = useCallback(async (event) => {
    const files = Array.from(event.target.files);
    if (files.length === 0) return;

    setUploading(true);
    setProgress(0);
    setError(null);

    try {
      // Auto-detect upload mode based on file characteristics
      const totalSize = files.reduce((sum, file) => sum + file.size, 0);
      const avgSize = totalSize / files.length;
      const mode = avgSize > 2 * 1024 * 1024 ? 'sequential' : 'parallel'; // 2MB threshold

      const result = await client.uploadBatchImages(files, {
        onProgress: setProgress,
        onBatchProgress: ({completed, total}) => {
          console.log(`Completed ${completed}/${total} files`);
        }
      });

      // Add responsive image data
      const imagesWithResponsive = result.data.map(img => ({
        ...img,
        responsive: generator.generateResponsiveUrls(img.public_id)
      }));

      setImages(prev => [...prev, ...imagesWithResponsive]);
    } catch (err) {
      setError(err.message);
    } finally {
      setUploading(false);
      setProgress(0);
      event.target.value = ''; // Reset input
    }
  }, [client, generator]);

  const deleteImage = useCallback(async (publicId) => {
    try {
      await client.deleteImage(publicId);
      setImages(prev => prev.filter(img => img.public_id !== publicId));
    } catch (err) {
      console.error('Delete failed:', err);
    }
  }, [client]);

  return (
    <div className="image-gallery">
      <div className="upload-controls">
        <input
          ref={fileInputRef}
          type="file"
          multiple
          accept="image/*"
          onChange={handleFileSelect}
          style={{ display: 'none' }}
        />
        <button
          onClick={() => fileInputRef.current.click()}
          disabled={uploading}
        >
          {uploading ? 'Uploading...' : 'Select Images'}
        </button>
      </div>

      {uploading && (
        <div className="progress-container">
          <div className="progress-bar">
            <div className="progress-fill" style={{ width: `${progress}%` }} />
          </div>
          <span>{progress}%</span>
        </div>
      )}

      {error && <div className="error">{error}</div>}

      <div className="image-grid">
        {images.map((image) => (
          <div key={image.public_id} className="image-item">
            <img
              src={image.responsive.src}
              srcSet={image.responsive.srcSet}
              sizes={image.responsive.sizes}
              alt={image.filename}
              loading="lazy"
            />
            <button onClick={() => deleteImage(image.public_id)}>×</button>
          </div>
        ))}
      </div>
    </div>
  );
};

export default ImageGallery;
```

### Vue.js Chat Image Component

```vue
<template>
  <div class="chat-image-uploader">
    <div class="image-input">
      <input
        ref="fileInput"
        type="file"
        accept="image/*"
        @change="handleFileSelect"
        style="display: none"
      />
      <button @click="$refs.fileInput.click()" :disabled="uploading">
        📎 Attach Image
      </button>
    </div>

    <div v-if="selectedFile" class="image-preview">
      <img :src="previewUrl" alt="Preview" />
      <div class="caption-input">
        <input
          v-model="caption"
          placeholder="Add a caption..."
          maxlength="500"
        />
      </div>
      <div class="actions">
        <button @click="sendImage" :disabled="uploading">
          {{ uploading ? 'Sending...' : 'Send' }}
        </button>
        <button @click="cancel" :disabled="uploading">Cancel</button>
      </div>
    </div>

    <div v-if="uploading" class="progress">
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: progress + '%' }" />
      </div>
      <span>{{ progress }}%</span>
    </div>

    <div v-if="error" class="error">
      {{ error }}
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue';
import { ImageUploadClient } from './image-upload-client.js';

export default {
  name: 'ChatImageUploader',
  props: {
    channelId: {
      type: String,
      required: true
    }
  },
  emits: ['image-sent'],
  setup(props, { emit }) {
    const fileInput = ref(null);
    const selectedFile = ref(null);
    const previewUrl = ref('');
    const caption = ref('');
    const uploading = ref(false);
    const progress = ref(0);
    const error = ref('');

    const client = computed(() => new ImageUploadClient({
      apiBaseUrl: '/api',
      authToken: localStorage.getItem('authToken')
    }));

    const handleFileSelect = (event) => {
      const file = event.target.files[0];
      if (!file) return;

      // Validate file
      if (!file.type.startsWith('image/')) {
        error.value = 'Please select an image file';
        return;
      }

      if (file.size > 10 * 1024 * 1024) { // 10MB
        error.value = 'File too large (max 10MB)';
        return;
      }

      selectedFile.value = file;
      previewUrl.value = URL.createObjectURL(file);
      error.value = '';
    };

    const sendImage = async () => {
      if (!selectedFile.value) return;

      uploading.value = true;
      progress.value = 0;
      error.value = '';

      try {
        const result = await client.value.sendChatImage(
          props.channelId,
          selectedFile.value,
          caption.value,
          {
            onProgress: (percent) => progress.value = percent
          }
        );

        emit('image-sent', result.message);
        cancel(); // Reset form
      } catch (err) {
        error.value = err.message;
      } finally {
        uploading.value = false;
        progress.value = 0;
      }
    };

    const cancel = () => {
      selectedFile.value = null;
      if (previewUrl.value) {
        URL.revokeObjectURL(previewUrl.value);
        previewUrl.value = '';
      }
      caption.value = '';
      error.value = '';
      if (fileInput.value) {
        fileInput.value.value = '';
      }
    };

    // Cleanup on unmount
    const cleanup = () => {
      if (previewUrl.value) {
        URL.revokeObjectURL(previewUrl.value);
      }
    };

    return {
      fileInput,
      selectedFile,
      previewUrl,
      caption,
      uploading,
      progress,
      error,
      handleFileSelect,
      sendImage,
      cancel,
      cleanup
    };
  },
  unmounted() {
    this.cleanup();
  }
};
</script>

<style scoped>
.chat-image-uploader {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 16px;
  margin: 8px 0;
}

.image-input {
  margin-bottom: 12px;
}

.image-preview {
  margin: 12px 0;
  padding: 12px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  background: #f9f9f9;
}

.image-preview img {
  max-width: 200px;
  max-height: 200px;
  border-radius: 4px;
  margin-bottom: 8px;
}

.caption-input {
  margin: 8px 0;
}

.caption-input input {
  width: 100%;
  padding: 6px 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 14px;
}

.actions {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

.actions button {
  padding: 6px 12px;
  border: 1px solid #ccc;
  border-radius: 4px;
  background: white;
  cursor: pointer;
}

.actions button:first-child {
  background: #007bff;
  color: white;
  border-color: #007bff;
}

.actions button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.progress {
  margin: 8px 0;
}

.progress-bar {
  width: 100%;
  height: 6px;
  background: #f0f0f0;
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 4px;
}

.progress-fill {
  height: 100%;
  background: #28a745;
  transition: width 0.3s ease;
}

.error {
  color: #dc3545;
  font-size: 14px;
  margin: 8px 0;
}
</style>
```

---

## 🚀 **Ready to Start Implementing!**

This comprehensive frontend integration guide provides everything you need to:

- ✅ **Set up the JavaScript client** with proper configuration
- ✅ **Implement single and batch image uploads** with progress tracking
- ✅ **Handle authentication and error scenarios**
- ✅ **Create responsive React and Vue.js components**
- ✅ **Build complete working examples** with drag & drop
- ✅ **Follow best practices** for performance and reliability
- ✅ **Integrate with chat systems** for image messaging

**Start with the basic examples and gradually add the advanced features as needed!** 🎯</result>
</attempt_completion>