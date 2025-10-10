/**
 * React Image Upload Components
 *
 * Complete React implementation with hooks, progress tracking,
 * drag-and-drop, and responsive image display.
 */

import React, { useState, useCallback, useRef } from 'react';
import { ImageUploadClient, ResponsiveImageGenerator } from './image-upload-client.js';

// Main Image Upload Component
export const ImageUploader = ({
  onUploadSuccess,
  onUploadError,
  maxFiles = 10,
  acceptedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'],
  apiBaseUrl = '/api',
  authToken
}) => {
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [error, setError] = useState(null);
  const [uploadedImages, setUploadedImages] = useState([]);
  const [dragActive, setDragActive] = useState(false);

  const fileInputRef = useRef(null);
  const client = React.useMemo(() =>
    new ImageUploadClient({ apiBaseUrl, authToken }), [apiBaseUrl, authToken]
  );

  const handleFiles = useCallback(async (files) => {
    const validFiles = Array.from(files).filter(file => {
      if (!acceptedTypes.includes(file.type)) {
        setError(`Invalid file type: ${file.name}`);
        return false;
      }
      if (file.size > 10 * 1024 * 1024) { // 10MB
        setError(`File too large: ${file.name}`);
        return false;
      }
      return true;
    });

    if (validFiles.length === 0) return;
    if (validFiles.length > maxFiles) {
      setError(`Too many files. Maximum ${maxFiles} files allowed.`);
      return;
    }

    setUploading(true);
    setProgress(0);
    setError(null);

    try {
      let result;
      if (validFiles.length === 1) {
        result = await client.uploadImage(validFiles[0], {
          onProgress: setProgress
        });
        setUploadedImages(prev => [...prev, result.data]);
      } else {
        result = await client.uploadBatchImages(validFiles, {
          onProgress: setProgress
        });
        setUploadedImages(prev => [...prev, ...result.data]);
      }

      onUploadSuccess?.(result);
    } catch (err) {
      setError(err.message);
      onUploadError?.(err);
    } finally {
      setUploading(false);
      setProgress(0);
    }
  }, [client, acceptedTypes, maxFiles, onUploadSuccess, onUploadError]);

  const handleDrag = useCallback((e) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  }, []);

  const handleDrop = useCallback((e) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);

    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      handleFiles(e.dataTransfer.files);
    }
  }, [handleFiles]);

  const handleFileSelect = useCallback((e) => {
    if (e.target.files && e.target.files[0]) {
      handleFiles(e.target.files);
    }
  }, [handleFiles]);

  const openFileDialog = () => {
    fileInputRef.current?.click();
  };

  return (
    <div className="image-uploader">
      <div
        className={`upload-zone ${dragActive ? 'drag-active' : ''} ${uploading ? 'uploading' : ''}`}
        onDragEnter={handleDrag}
        onDragLeave={handleDrag}
        onDragOver={handleDrag}
        onDrop={handleDrop}
        onClick={openFileDialog}
      >
        <input
          ref={fileInputRef}
          type="file"
          multiple
          accept={acceptedTypes.join(',')}
          onChange={handleFileSelect}
          style={{ display: 'none' }}
        />

        <div className="upload-content">
          {uploading ? (
            <div className="upload-progress">
              <div className="progress-bar">
                <div
                  className="progress-fill"
                  style={{ width: `${progress}%` }}
                />
              </div>
              <div className="progress-text">
                Uploading... {progress}%
              </div>
            </div>
          ) : (
            <>
              <div className="upload-icon">📁</div>
              <div className="upload-text">
                Drag & drop images here or click to select
              </div>
              <div className="upload-hint">
                Supports JPEG, PNG, GIF, WebP (max 10MB each)
              </div>
            </>
          )}
        </div>
      </div>

      {error && (
        <div className="upload-error">
          {error}
        </div>
      )}

      {uploadedImages.length > 0 && (
        <div className="uploaded-images">
          <h3>Uploaded Images</h3>
          <div className="image-grid">
            {uploadedImages.map((image, index) => (
              <UploadedImage
                key={image.public_id}
                image={image}
                onDelete={() => {
                  setUploadedImages(prev => prev.filter((_, i) => i !== index));
                }}
              />
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

// Responsive Image Component
export const ResponsiveImage = ({
  publicId,
  alt = '',
  className = '',
  sizes = '(max-width: 768px) 100vw, 50vw',
  cloudName
}) => {
  const generator = React.useMemo(() =>
    new ResponsiveImageGenerator(cloudName), [cloudName]
  );

  const responsiveData = React.useMemo(() =>
    generator.generateResponsiveUrls(publicId), [generator, publicId]
  );

  return (
    <img
      src={responsiveData.src}
      srcSet={responsiveData.srcSet}
      sizes={sizes}
      alt={alt}
      className={className}
      loading="lazy"
    />
  );
};

// Thumbnail Component
export const ImageThumbnail = ({
  publicId,
  alt = '',
  className = '',
  size = 200,
  cloudName
}) => {
  const generator = React.useMemo(() =>
    new ResponsiveImageGenerator(cloudName), [cloudName]
  );

  const thumbnailUrl = React.useMemo(() =>
    generator.generateThumbnailUrl(publicId, { width: size, height: size }),
    [generator, publicId, size]
  );

  return (
    <img
      src={thumbnailUrl}
      alt={alt}
      className={className}
      loading="lazy"
    />
  );
};

// Individual Uploaded Image Component
const UploadedImage = ({ image, onDelete }) => {
  const [showDelete, setShowDelete] = useState(false);

  return (
    <div
      className="uploaded-image"
      onMouseEnter={() => setShowDelete(true)}
      onMouseLeave={() => setShowDelete(false)}
    >
      <ImageThumbnail
        publicId={image.public_id}
        alt={image.filename}
        className="image-preview"
        cloudName="your-cloud-name" // Replace with actual cloud name
      />

      <div className="image-info">
        <div className="image-name">{image.filename}</div>
        <div className="image-size">
          {(image.bytes / 1024 / 1024).toFixed(2)} MB
        </div>
        <div className="image-dimensions">
          {image.width} × {image.height}
        </div>
      </div>

      {showDelete && (
        <button
          className="delete-button"
          onClick={onDelete}
          title="Delete image"
        >
          ×
        </button>
      )}
    </div>
  );
};

// Chat Image Upload Component
export const ChatImageUploader = ({
  channelId,
  onMessageSent,
  apiBaseUrl = '/api',
  authToken
}) => {
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [error, setError] = useState(null);
  const [caption, setCaption] = useState('');

  const client = React.useMemo(() =>
    new ImageUploadClient({ apiBaseUrl, authToken }), [apiBaseUrl, authToken]
  );

  const handleImageSelect = async (file) => {
    if (!file) return;

    setUploading(true);
    setProgress(0);
    setError(null);

    try {
      const result = await client.sendChatImage(channelId, file, caption, {
        onProgress: setProgress
      });

      onMessageSent?.(result.message);
      setCaption('');
    } catch (err) {
      setError(err.message);
    } finally {
      setUploading(false);
      setProgress(0);
    }
  };

  return (
    <div className="chat-image-uploader">
      <div className="caption-input">
        <input
          type="text"
          placeholder="Add a caption (optional)"
          value={caption}
          onChange={(e) => setCaption(e.target.value)}
          maxLength={500}
        />
      </div>

      <div className="image-input">
        <input
          type="file"
          accept="image/*"
          onChange={(e) => handleImageSelect(e.target.files[0])}
          disabled={uploading}
        />

        {uploading && (
          <div className="upload-progress">
            <div className="progress-bar">
              <div
                className="progress-fill"
                style={{ width: `${progress}%` }}
              />
            </div>
            <div className="progress-text">
              Sending image... {progress}%
            </div>
          </div>
        )}
      </div>

      {error && (
        <div className="upload-error">
          {error}
        </div>
      )}
    </div>
  );
};

// CSS Styles (can be moved to separate CSS file)
const styles = `
.image-uploader {
  max-width: 600px;
  margin: 0 auto;
}

.upload-zone {
  border: 2px dashed #ccc;
  border-radius: 8px;
  padding: 40px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #fafafa;
}

.upload-zone:hover,
.upload-zone.drag-active {
  border-color: #007bff;
  background: #f0f8ff;
}

.upload-zone.uploading {
  pointer-events: none;
  opacity: 0.7;
}

.upload-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.upload-icon {
  font-size: 48px;
}

.upload-text {
  font-size: 18px;
  font-weight: 500;
  color: #333;
}

.upload-hint {
  font-size: 14px;
  color: #666;
}

.upload-progress {
  width: 100%;
  max-width: 300px;
}

.progress-bar {
  width: 100%;
  height: 8px;
  background: #e9ecef;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 8px;
}

.progress-fill {
  height: 100%;
  background: #007bff;
  transition: width 0.3s ease;
}

.progress-text {
  text-align: center;
  font-size: 14px;
  color: #666;
}

.upload-error {
  color: #dc3545;
  background: #f8d7da;
  border: 1px solid #f5c6cb;
  border-radius: 4px;
  padding: 12px;
  margin-top: 16px;
}

.uploaded-images {
  margin-top: 24px;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
  margin-top: 16px;
}

.uploaded-image {
  position: relative;
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
  background: white;
}

.image-preview {
  width: 100%;
  height: 150px;
  object-fit: cover;
  display: block;
}

.image-info {
  padding: 12px;
}

.image-name {
  font-weight: 500;
  font-size: 14px;
  margin-bottom: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.image-size,
.image-dimensions {
  font-size: 12px;
  color: #666;
}

.delete-button {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: rgba(220, 53, 69, 0.9);
  color: white;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: bold;
  transition: background-color 0.2s;
}

.delete-button:hover {
  background: rgba(220, 53, 69, 1);
}

.chat-image-uploader {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  border-top: 1px solid #eee;
}

.caption-input input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.caption-input input:focus {
  outline: none;
  border-color: #007bff;
}

.image-input {
  display: flex;
  align-items: center;
  gap: 12px;
}

.image-input input[type="file"] {
  flex: 1;
}
`;

// Inject styles (in a real app, you'd use CSS modules or styled-components)
if (typeof document !== 'undefined') {
  const styleSheet = document.createElement('style');
  styleSheet.textContent = styles;
  document.head.appendChild(styleSheet);
}

export default {
  ImageUploader,
  ResponsiveImage,
  ImageThumbnail,
  ChatImageUploader
};