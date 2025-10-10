<template>
  <div class="image-upload-demo">
    <h1>DKL Image Upload Demo</h1>

    <!-- Single Image Upload -->
    <div class="upload-section">
      <h2>Single Image Upload</h2>
      <ImageUploader
        @upload-success="onSingleUploadSuccess"
        @upload-error="onUploadError"
        :max-files="1"
      />
    </div>

    <!-- Batch Image Upload -->
    <div class="upload-section">
      <h2>Batch Image Upload (Max 5)</h2>
      <ImageUploader
        @upload-success="onBatchUploadSuccess"
        @upload-error="onUploadError"
        :max-files="5"
      />
    </div>

    <!-- Chat Image Upload -->
    <div class="upload-section">
      <h2>Chat Image Upload</h2>
      <ChatImageUploader
        :channel-id="currentChannelId"
        @message-sent="onMessageSent"
      />
    </div>

    <!-- Uploaded Images Gallery -->
    <div v-if="uploadedImages.length > 0" class="gallery-section">
      <h2>Uploaded Images Gallery</h2>
      <div class="image-gallery">
        <div
          v-for="image in uploadedImages"
          :key="image.public_id"
          class="gallery-item"
        >
          <ResponsiveImage
            :public-id="image.public_id"
            :alt="image.filename"
            class="gallery-image"
            :cloud-name="cloudName"
          />
          <div class="image-meta">
            <h3>{{ image.filename }}</h3>
            <p>{{ formatFileSize(image.bytes) }}</p>
            <p>{{ image.width }} × {{ image.height }}</p>
            <button
              @click="deleteImage(image.public_id)"
              class="delete-btn"
            >
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Error Display -->
    <div v-if="error" class="error-message">
      {{ error }}
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ImageUploadClient, ResponsiveImageGenerator } from './image-upload-client.js'

// Components
import ImageUploader from './components/ImageUploader.vue'
import ChatImageUploader from './components/ChatImageUploader.vue'
import ResponsiveImage from './components/ResponsiveImage.vue'

// Reactive data
const uploadedImages = ref([])
const error = ref(null)
const currentChannelId = ref('demo-channel-123')
const cloudName = ref('your-cloud-name') // Replace with actual cloud name

// Computed
const client = computed(() => new ImageUploadClient({
  apiBaseUrl: '/api',
  authToken: localStorage.getItem('authToken') // Get from your auth system
}))

// Methods
const onSingleUploadSuccess = (result) => {
  uploadedImages.value.push(result.data)
  error.value = null
}

const onBatchUploadSuccess = (result) => {
  uploadedImages.value.push(...result.data)
  error.value = null
}

const onUploadError = (err) => {
  error.value = err.message
}

const onMessageSent = (message) => {
  console.log('Chat message sent:', message)
  // Handle chat message in your chat system
}

const deleteImage = async (publicId) => {
  try {
    await client.value.deleteImage(publicId)
    uploadedImages.value = uploadedImages.value.filter(
      img => img.public_id !== publicId
    )
  } catch (err) {
    error.value = `Failed to delete image: ${err.message}`
  }
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}
</script>

<style scoped>
.image-upload-demo {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.upload-section {
  margin-bottom: 40px;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
}

.upload-section h2 {
  margin-top: 0;
  color: #333;
}

.gallery-section {
  margin-top: 40px;
}

.image-gallery {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.gallery-item {
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
  background: white;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.gallery-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.image-meta {
  padding: 15px;
}

.image-meta h3 {
  margin: 0 0 10px 0;
  font-size: 16px;
  font-weight: 500;
}

.image-meta p {
  margin: 5px 0;
  color: #666;
  font-size: 14px;
}

.delete-btn {
  background: #dc3545;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  margin-top: 10px;
  transition: background-color 0.2s;
}

.delete-btn:hover {
  background: #c82333;
}

.error-message {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
  border-radius: 4px;
  padding: 12px;
  margin-top: 20px;
}
</style>