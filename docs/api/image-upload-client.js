/**
 * DKL Email Service - Image Upload Client
 *
 * Complete JavaScript client for image upload functionality
 * with progress tracking, error handling, and responsive image support.
 *
 * Usage Examples:
 *
 * // Parallel batch upload (default - faster for small files)
 * const result = await client.uploadBatchImages(files, {
 *   onProgress: (percent) => console.log(`${percent}%`),
 *   onBatchProgress: (progress) => console.log(progress)
 * });
 *
 * // Sequential batch upload (better for large files/slow connections)
 * const result = await client.uploadBatchImagesSequential(files, {
 *   onProgress: (percent) => console.log(`${percent}%`),
 *   onFileProgress: (fileProgress) => console.log(fileProgress)
 * });
 */

class ImageUploadClient {
  constructor(options = {}) {
    this.apiBaseUrl = options.apiBaseUrl || '/api';
    this.authToken = options.authToken;
    this.maxFileSize = options.maxFileSize || 10 * 1024 * 1024; // 10MB
    this.allowedTypes = options.allowedTypes || [
      'image/jpeg',
      'image/png',
      'image/gif',
      'image/webp'
    ];
    this.maxBatchSize = options.maxBatchSize || 10;
  }

  /**
   * Set authentication token
   * @param {string} token - JWT token
   */
  setAuthToken(token) {
    this.authToken = token;
  }

  /**
   * Validate image file
   * @param {File} file - File to validate
   * @throws {Error} Validation error
   */
  validateImageFile(file) {
    if (!this.allowedTypes.includes(file.type)) {
      throw new Error('Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed.');
    }

    if (file.size > this.maxFileSize) {
      throw new Error('File too large. Maximum size is 10MB.');
    }

    return true;
  }

  /**
   * Upload single image with progress tracking
   * @param {File} file - Image file to upload
   * @param {Object} options - Upload options
   * @param {Function} options.onProgress - Progress callback (percent)
   * @param {Function} options.onSuccess - Success callback (result)
   * @param {Function} options.onError - Error callback (error)
   * @returns {Promise<Object>} Upload result
   */
  async uploadImage(file, options = {}) {
    this.validateImageFile(file);

    const formData = new FormData();
    formData.append('image', file);

    return this._uploadWithProgress('/images/upload', formData, options);
  }

  /**
   * Upload multiple images with progress tracking (parallel mode - default)
   * @param {File[]} files - Array of image files
   * @param {Object} options - Upload options
   * @param {Function} options.onProgress - Progress callback (percent)
   * @param {Function} options.onBatchProgress - Batch progress callback ({completed, total, currentFile})
   * @param {Function} options.onSuccess - Success callback (results)
   * @param {Function} options.onError - Error callback (error)
   * @returns {Promise<Object[]>} Array of upload results
   */
  async uploadBatchImages(files, options = {}) {
    return this._uploadBatchImages(files, { ...options, mode: 'parallel' });
  }

  /**
   * Upload multiple images sequentially (better for large files/slow connections)
   * @param {File[]} files - Array of image files
   * @param {Object} options - Upload options
   * @param {Function} options.onProgress - Progress callback (percent)
   * @param {Function} options.onFileProgress - File-specific progress callback ({fileIndex, fileName, progress})
   * @param {Function} options.onBatchProgress - Batch progress callback ({completed, total, currentFile})
   * @param {Function} options.onSuccess - Success callback (results)
   * @param {Function} options.onError - Error callback (error)
   * @returns {Promise<Object[]>} Array of upload results
   */
  async uploadBatchImagesSequential(files, options = {}) {
    return this._uploadBatchImages(files, { ...options, mode: 'sequential' });
  }

  /**
   * Internal batch upload implementation
   * @private
   */
  async _uploadBatchImages(files, options = {}) {
    if (files.length === 0 || files.length > this.maxBatchSize) {
      throw new Error(`Provide 1-${this.maxBatchSize} image files`);
    }

    // Validate all files first
    files.forEach(file => this.validateImageFile(file));

    const { mode = 'parallel' } = options;
    const results = [];
    const errors = [];

    if (mode === 'sequential') {
      // Sequential upload: process files one by one
      for (let i = 0; i < files.length; i++) {
        const file = files[i];

        try {
          if (options.onBatchProgress) {
            options.onBatchProgress({
              completed: results.length,
              total: files.length,
              currentFile: { index: i, name: file.name }
            });
          }

          const result = await this.uploadImage(file, {
            onProgress: (percent) => {
              if (options.onFileProgress) {
                options.onFileProgress({
                  fileIndex: i,
                  fileName: file.name,
                  progress: percent
                });
              }
              if (options.onProgress) {
                // Calculate overall progress
                const baseProgress = (i / files.length) * 100;
                const currentFileProgress = percent / files.length;
                options.onProgress(Math.round(baseProgress + currentFileProgress));
              }
            }
          });

          results.push(result.data);

          if (options.onBatchProgress) {
            options.onBatchProgress({
              completed: results.length,
              total: files.length,
              currentFile: null
            });
          }
        } catch (error) {
          errors.push({
            fileIndex: i,
            fileName: file.name,
            error: error.message
          });

          if (options.onError) {
            options.onError(error);
          }

          // Continue with next file in sequential mode
        }
      }
    } else {
      // Parallel upload: process all files simultaneously
      const formData = new FormData();
      files.forEach(file => formData.append('images', file));

      const batchOptions = {
        ...options,
        onProgress: (percent) => {
          if (options.onProgress) options.onProgress(percent);
        },
        onSuccess: (result) => {
          results.push(...result.data);
          if (options.onBatchProgress) {
            options.onBatchProgress({
              completed: results.length,
              total: files.length,
              currentFile: null
            });
          }
        }
      };

      await this._uploadWithProgress(`/images/batch-upload?mode=${mode}`, formData, batchOptions);
    }

    // Return comprehensive result
    const finalResult = {
      success: results.length > 0,
      data: results,
      uploaded_count: results.length,
      total_files: files.length,
      mode: mode
    };

    if (errors.length > 0) {
      finalResult.errors = errors;
      finalResult.errors_count = errors.length;
    }

    if (options.onSuccess) {
      options.onSuccess(finalResult);
    }

    return finalResult;
  }

  /**
   * Send image message in chat
   * @param {string} channelId - Chat channel ID
   * @param {File} file - Image file
   * @param {string} caption - Optional caption text
   * @param {Object} options - Upload options
   * @returns {Promise<Object>} Chat message result
   */
  async sendChatImage(channelId, file, caption = '', options = {}) {
    this.validateImageFile(file);

    const formData = new FormData();
    formData.append('image', file);
    if (caption.trim()) {
      formData.append('content', caption.trim());
    }

    return this._uploadWithProgress(`/chat/channels/${channelId}/messages`, formData, options);
  }

  /**
   * Delete image
   * @param {string} publicId - Cloudinary public ID
   * @returns {Promise<Object>} Delete result
   */
  async deleteImage(publicId) {
    const response = await this._makeRequest(`/images/${publicId}`, {
      method: 'DELETE'
    });

    return response;
  }

  /**
   * Get image metadata
   * @param {string} publicId - Cloudinary public ID
   * @returns {Promise<Object>} Image metadata
   */
  async getImageMetadata(publicId) {
    const response = await this._makeRequest(`/images/${publicId}`, {
      method: 'GET'
    });

    return response.data;
  }

  /**
   * Internal upload method with progress tracking
   * @private
   */
  async _uploadWithProgress(endpoint, formData, options = {}) {
    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest();

      // Progress tracking
      xhr.upload.addEventListener('progress', (event) => {
        if (event.lengthComputable && options.onProgress) {
          const percentComplete = (event.loaded / event.total) * 100;
          options.onProgress(Math.round(percentComplete));
        }
      });

      // Load completion
      xhr.addEventListener('load', () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          try {
            const response = JSON.parse(xhr.responseText);
            if (options.onSuccess) options.onSuccess(response);
            resolve(response);
          } catch (e) {
            const error = new Error('Invalid response format');
            if (options.onError) options.onError(error);
            reject(error);
          }
        } else {
          let errorMessage = `Upload failed with status ${xhr.status}`;
          try {
            const errorData = JSON.parse(xhr.responseText);
            errorMessage = errorData.error || errorMessage;
          } catch (e) {
            // Use default error message
          }
          const error = new Error(errorMessage);
          if (options.onError) options.onError(error);
          reject(error);
        }
      });

      // Error handling
      xhr.addEventListener('error', () => {
        const error = new Error('Network error during upload');
        if (options.onError) options.onError(error);
        reject(error);
      });

      // Abort handling
      xhr.addEventListener('abort', () => {
        const error = new Error('Upload was cancelled');
        if (options.onError) options.onError(error);
        reject(error);
      });

      // Configure request
      xhr.open('POST', `${this.apiBaseUrl}${endpoint}`);
      xhr.setRequestHeader('Authorization', `Bearer ${this.authToken}`);
      xhr.send(formData);
    });
  }

  /**
   * Internal request method for non-upload operations
   * @private
   */
  async _makeRequest(endpoint, options = {}) {
    const url = `${this.apiBaseUrl}${endpoint}`;
    const config = {
      method: options.method || 'GET',
      headers: {
        'Authorization': `Bearer ${this.authToken}`,
        'Content-Type': 'application/json',
        ...options.headers
      },
      ...options
    };

    const response = await fetch(url, config);

    if (!response.ok) {
      let errorMessage = `Request failed with status ${response.status}`;
      try {
        const errorData = await response.json();
        errorMessage = errorData.error || errorMessage;
      } catch (e) {
        // Use default error message
      }
      throw new Error(errorMessage);
    }

    return response.json();
  }
}

/**
 * Responsive Image URL Generator
 * Generates optimized Cloudinary URLs for different screen sizes
 */
class ResponsiveImageGenerator {
  constructor(cloudName) {
    this.cloudName = cloudName;
    this.baseUrl = `https://res.cloudinary.com/${cloudName}/image/upload/`;
  }

  /**
   * Generate responsive image URLs
   * @param {string} publicId - Cloudinary public ID
   * @param {Object} options - Generation options
   * @returns {Object} Responsive image data
   */
  generateResponsiveUrls(publicId, options = {}) {
    const {
      sizes = [
        { width: 480, dpr: 1 },
        { width: 768, dpr: 1 },
        { width: 1024, dpr: 1 },
        { width: 480, dpr: 2 },
        { width: 768, dpr: 2 }
      ],
      quality = 'auto',
      format = 'auto'
    } = options;

    const src = this._generateUrl(publicId, { quality, format });
    const srcSet = sizes.map(size =>
      `${this._generateUrl(publicId, { ...size, quality, format })} ${size.width}w`
    ).join(', ');

    return {
      src,
      srcSet,
      sizes: '(max-width: 768px) 100vw, 50vw'
    };
  }

  /**
   * Generate thumbnail URL
   * @param {string} publicId - Cloudinary public ID
   * @param {Object} options - Thumbnail options
   * @returns {string} Thumbnail URL
   */
  generateThumbnailUrl(publicId, options = {}) {
    const {
      width = 200,
      height = 200,
      crop = 'thumb',
      gravity = 'face'
    } = options;

    return this._generateUrl(publicId, {
      width,
      height,
      crop,
      gravity,
      quality: 'auto',
      format: 'auto'
    });
  }

  /**
   * Generate custom transformation URL
   * @param {string} publicId - Cloudinary public ID
   * @param {Object} transformations - Transformation parameters
   * @returns {string} Transformed URL
   */
  generateCustomUrl(publicId, transformations = {}) {
    return this._generateUrl(publicId, transformations);
  }

  /**
   * Internal URL generation
   * @private
   */
  _generateUrl(publicId, transformations) {
    const params = [];

    // Add transformations
    if (transformations.width) params.push(`w_${transformations.width}`);
    if (transformations.height) params.push(`h_${transformations.height}`);
    if (transformations.crop) params.push(`c_${transformations.crop}`);
    if (transformations.gravity) params.push(`g_${transformations.gravity}`);
    if (transformations.quality) params.push(`q_${transformations.quality}`);
    if (transformations.format) params.push(`f_${transformations.format}`);
    if (transformations.dpr) params.push(`dpr_${transformations.dpr}`);

    const transformationString = params.length > 0 ? params.join(',') + '/' : '';
    return `${this.baseUrl}${transformationString}${publicId}`;
  }
}

/**
 * React Hook for Image Upload (React example)
 */
function useImageUpload(apiBaseUrl, authToken) {
  const [uploading, setUploading] = React.useState(false);
  const [progress, setProgress] = React.useState(0);
  const [error, setError] = React.useState(null);

  const client = React.useMemo(() =>
    new ImageUploadClient({ apiBaseUrl, authToken }), [apiBaseUrl, authToken]
  );

  const uploadImage = React.useCallback(async (file) => {
    setUploading(true);
    setProgress(0);
    setError(null);

    try {
      const result = await client.uploadImage(file, {
        onProgress: setProgress
      });
      return result;
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setUploading(false);
      setProgress(0);
    }
  }, [client]);

  const uploadBatch = React.useCallback(async (files, mode = 'parallel') => {
    setUploading(true);
    setProgress(0);
    setError(null);

    try {
      const uploadMethod = mode === 'sequential' ? 'uploadBatchImagesSequential' : 'uploadBatchImages';
      const result = await client[uploadMethod](files, {
        onProgress: setProgress,
        onBatchProgress: (progress) => {
          // Optional: handle individual file progress
          console.log('Batch progress:', progress);
        }
      });
      return result;
    } catch (err) {
      setError(err.message);
      throw err;
    } finally {
      setUploading(false);
      setProgress(0);
    }
  }, [client]);

  return {
    uploadImage,
    uploadBatch,
    uploading,
    progress,
    error
  };
}

/**
 * Vue Composable for Image Upload (Vue example)
 */
function useImageUpload(apiBaseUrl, authToken) {
  const uploading = Vue.ref(false);
  const progress = Vue.ref(0);
  const error = Vue.ref(null);

  const client = Vue.computed(() =>
    new ImageUploadClient({ apiBaseUrl: Vue.unref(apiBaseUrl), authToken: Vue.unref(authToken) })
  );

  const uploadImage = async (file) => {
    uploading.value = true;
    progress.value = 0;
    error.value = null;

    try {
      const result = await client.value.uploadImage(file, {
        onProgress: (percent) => progress.value = percent
      });
      return result;
    } catch (err) {
      error.value = err.message;
      throw err;
    } finally {
      uploading.value = false;
      progress.value = 0;
    }
  };

  const uploadBatch = async (files, mode = 'parallel') => {
    uploading.value = true;
    progress.value = 0;
    error.value = null;

    try {
      const uploadMethod = mode === 'sequential' ? 'uploadBatchImagesSequential' : 'uploadBatchImages';
      const result = await client.value[uploadMethod](files, {
        onProgress: (percent) => progress.value = percent,
        onBatchProgress: (progress) => {
          // Optional: handle individual file progress
          console.log('Batch progress:', progress);
        }
      });
      return result;
    } catch (err) {
      error.value = err.message;
      throw err;
    } finally {
      uploading.value = false;
      progress.value = 0;
    }
  };

  return {
    uploadImage: Vue.markRaw(uploadImage),
    uploadBatch: Vue.markRaw(uploadBatch),
    uploading: Vue.readonly(uploading),
    progress: Vue.readonly(progress),
    error: Vue.readonly(error)
  };
}

// Export for different module systems
if (typeof module !== 'undefined' && module.exports) {
  // CommonJS
  module.exports = { ImageUploadClient, ResponsiveImageGenerator, useImageUpload };
} else if (typeof define === 'function' && define.amd) {
  // AMD
  define([], function() {
    return { ImageUploadClient, ResponsiveImageGenerator, useImageUpload };
  });
} else if (typeof window !== 'undefined') {
  // Browser global
  window.DKLImageUpload = { ImageUploadClient, ResponsiveImageGenerator, useImageUpload };
}