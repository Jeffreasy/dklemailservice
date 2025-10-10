# Cloudinary CDN Integration Plan for DKL Email Service

## Application Analysis Summary

### Current Architecture
The DKL Email Service is a comprehensive Go backend application built with:
- **Framework**: Fiber (Go web framework)
- **Database**: PostgreSQL with GORM ORM
- **Caching**: Redis (optional, for RBAC and chat features)
- **Architecture**: Clean Architecture with handlers, services, repositories, and models
- **Features**: Email handling, authentication, contact forms, registrations, chat system, RBAC, newsletters

### Current Image Handling
- **Static Files**: Serves favicon.ico from `./public` directory
- **Chat System**: Supports image message types with `FileURL` field in `ChatMessage` model
- **No Upload Capability**: Currently no image upload endpoints or storage mechanisms
- **Frontend Serving**: Static files served via Fiber's static middleware

### Data Flow to Frontend
- **API-First**: Frontend consumes JSON APIs
- **Static Assets**: Served from `./public` directory
- **CORS**: Configured for specific frontend domains
- **Authentication**: JWT-based with protected endpoints

## Cloudinary Integration Requirements

### Cloudinary Service Overview
Cloudinary provides cloud-based image management including:
- Image upload and storage
- CDN delivery with optimization
- Image transformations (resize, crop, format conversion)
- Secure URLs with access controls
- Built-in malware scanning and moderation
- Webhook notifications for upload events

### Go SDK Requirements
- **Package**: `github.com/cloudinary/cloudinary-go/v2` (pin to v2.3.0 for stability)
- **Configuration**: Cloud name, API key, API secret, upload presets
- **Features Needed**: Upload, URL generation, transformations, async operations, webhooks

### Dependencies and Alternatives
**Why Cloudinary vs Alternatives:**
- Superior image optimization and CDN performance
- Built-in transformations reduce backend processing
- Malware scanning and moderation features
- Comprehensive Go SDK support
- Global CDN with edge locations

**Alternatives Considered:**
- AWS S3 + CloudFront: More infrastructure management, less optimization features
- Cloudflare Images: Similar but less mature Go SDK
- Self-hosted: Higher maintenance, scaling challenges

**Migration Path:** ImageService designed with interfaces for easy provider switching if needed.

## Integration Plan

### Phase 1: Infrastructure Setup (Estimated: 2-4 hours)

#### 1.1 Add Dependencies
```bash
go get github.com/cloudinary/cloudinary-go/v2@v2.3.0
```

#### 1.2 Environment Configuration
Add to `.env` and `.env.example`:
```env
# Cloudinary Configuration
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret
CLOUDINARY_UPLOAD_FOLDER=dkl_images
CLOUDINARY_UPLOAD_PRESET=dkl_preset
CLOUDINARY_SECURE=true

# Test/Development Environment (separate credentials)
CLOUDINARY_TEST_CLOUD_NAME=test_cloud_name
CLOUDINARY_TEST_API_KEY=test_api_key
CLOUDINARY_TEST_API_SECRET=test_api_secret
```

#### 1.3 Configuration Loading
Extend `config/config.go` to include Cloudinary configuration:
```go
type CloudinaryConfig struct {
    CloudName    string
    APIKey       string
    APISecret    string
    UploadFolder string
    UploadPreset string
    Secure       bool
    IsTest       bool
}

func LoadCloudinaryConfig() *CloudinaryConfig {
    isTest := os.Getenv("APP_ENV") == "test" || os.Getenv("APP_ENV") == "development"

    var cloudName, apiKey, apiSecret string
    if isTest {
        cloudName = os.Getenv("CLOUDINARY_TEST_CLOUD_NAME")
        apiKey = os.Getenv("CLOUDINARY_TEST_API_KEY")
        apiSecret = os.Getenv("CLOUDINARY_TEST_API_SECRET")
    } else {
        cloudName = os.Getenv("CLOUDINARY_CLOUD_NAME")
        apiKey = os.Getenv("CLOUDINARY_API_KEY")
        apiSecret = os.Getenv("CLOUDINARY_API_SECRET")
    }

    if cloudName == "" || apiKey == "" || apiSecret == "" {
        return nil // Cloudinary not configured
    }

    return &CloudinaryConfig{
        CloudName:    cloudName,
        APIKey:       apiKey,
        APISecret:    apiSecret,
        UploadFolder: getEnvOrDefault("CLOUDINARY_UPLOAD_FOLDER", "dkl_images"),
        UploadPreset: os.Getenv("CLOUDINARY_UPLOAD_PRESET"),
        Secure:       getEnvOrDefault("CLOUDINARY_SECURE", "true") == "true",
        IsTest:       isTest,
    }
}
```

### Phase 2: Core Service Implementation (Estimated: 4-6 hours)

#### 2.1 Image Service Creation
Create `services/image_service.go`:
```go
import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "fmt"
    "mime/multipart"
    "strings"
    "time"

    "github.com/cloudinary/cloudinary-go/v2/api/uploader"
    "github.com/google/uuid"
)

// Custom errors
var (
    ErrUploadFailed     = errors.New("image upload failed")
    ErrInvalidFileType  = errors.New("invalid file type")
    ErrFileTooLarge     = errors.New("file too large")
    ErrInvalidFile      = errors.New("invalid file")
)

type UploadResult struct {
    PublicID   string            `json:"public_id"`
    URL        string            `json:"url"`
    SecureURL  string            `json:"secure_url"`
    Width      int               `json:"width"`
    Height     int               `json:"height"`
    Format     string            `json:"format"`
    Bytes      int               `json:"bytes"`
    ThumbnailURL string          `json:"thumbnail_url,omitempty"`
}

type ImageService struct {
    cld    *cloudinary.Cloudinary
    config *config.CloudinaryConfig
}

func NewImageService(config *config.CloudinaryConfig) (*ImageService, error) {
    cld, err := cloudinary.NewFromParams(config.CloudName, config.APIKey, config.APISecret)
    if err != nil {
        return nil, err
    }

    return &ImageService{
        cld:    cld,
        config: config,
    }, nil
}

func (s *ImageService) UploadImage(ctx context.Context, file multipart.File, filename string, folder string, userID string) (*UploadResult, error) {
    // Generate unique public ID
    publicID := s.generatePublicID(userID, filename)

    // Upload parameters
    uploadParams := uploader.UploadParams{
        PublicID:     publicID,
        Folder:       folder,
        ResourceType: "image",
        Format:       "auto", // Auto-detect format
    }

    // Add upload preset if configured
    if s.config.UploadPreset != "" {
        uploadParams.UploadPreset = s.config.UploadPreset
    }

    // Add eager transformations for thumbnails
    uploadParams.Eager = []uploader.Eager{
        {Width: 200, Height: 200, Crop: "thumb", Gravity: "face"},
    }

    // Perform upload
    resp, err := s.cld.Upload.Upload(ctx, file, uploadParams)
    if err != nil {
        return nil, fmt.Errorf("%w: %v", ErrUploadFailed, err)
    }

    result := &UploadResult{
        PublicID:  resp.PublicID,
        URL:       resp.URL,
        SecureURL: resp.SecureURL,
        Width:     resp.Width,
        Height:    resp.Height,
        Format:    resp.Format,
        Bytes:     resp.Bytes,
    }

    // Generate thumbnail URL if eager transformation exists
    if len(resp.Eager) > 0 {
        result.ThumbnailURL = resp.Eager[0].SecureURL
    }

    return result, nil
}

func (s *ImageService) UploadImageAsync(ctx context.Context, file multipart.File, filename string, folder string, userID string, callbackURL string) (string, error) {
    publicID := s.generatePublicID(userID, filename)

    uploadParams := uploader.UploadParams{
        PublicID:     publicID,
        Folder:       folder,
        ResourceType: "image",
        Async:        true,
        CallbackURL:  callbackURL,
    }

    resp, err := s.cld.Upload.Upload(ctx, file, uploadParams)
    if err != nil {
        return "", fmt.Errorf("%w: %v", ErrUploadFailed, err)
    }

    return resp.PublicID, nil
}

func (s *ImageService) DeleteImage(ctx context.Context, publicID string) error {
    _, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
        PublicID:     publicID,
        ResourceType: "image",
        Invalidate:   true, // Clear CDN cache
    })
    return err
}

func (s *ImageService) GetImageURL(publicID string, transformations map[string]string) string {
    img, err := s.cld.Image(publicID)
    if err != nil {
        return ""
    }

    // Apply transformations
    for key, value := range transformations {
        switch strings.ToLower(key) {
        case "width":
            img.Width(value)
        case "height":
            img.Height(value)
        case "crop":
            img.Crop(value)
        case "quality":
            img.Quality(value)
        case "format":
            img.Format(value)
        }
    }

    url, _ := img.Signed(s.config.Secure).ToURL()
    return url
}

func (s *ImageService) generatePublicID(userID, filename string) string {
    // Generate unique ID combining user ID and random bytes
    randomBytes := make([]byte, 8)
    rand.Read(randomBytes)
    randomStr := hex.EncodeToString(randomBytes)

    // Extract file extension
    parts := strings.Split(filename, ".")
    ext := ""
    if len(parts) > 1 {
        ext = parts[len(parts)-1]
    }

    return fmt.Sprintf("%s/%s_%s", s.config.UploadFolder, userID, randomStr)
}
```

#### 2.2 Service Factory Integration
Extend `services/factory.go` to include ImageService:
```go
type ServiceFactory struct {
    // ... existing services
    ImageService *ImageService
}

// NewServiceFactory updates to initialize ImageService
func NewServiceFactory(repoFactory *repository.Repository) *ServiceFactory {
    // ... existing initialization
    
    if config := config.LoadCloudinaryConfig(); config != nil {
        imageService, err := NewImageService(config)
        if err != nil {
            logger.Warn("Failed to initialize Cloudinary service", "error", err)
        } else {
            factory.ImageService = imageService
        }
    }
    
    return factory
}
```

### Phase 3: API Endpoints (Estimated: 3-5 hours)

#### 3.1 Image Handler Creation
Create `handlers/image_handler.go`:
```go
import (
    "context"
    "strings"
    "time"

    "dklautomationgo/services"
    "github.com/gofiber/fiber/v2"
)

type ImageHandler struct {
    imageService *services.ImageService
    authService  *services.AuthService
}

func NewImageHandler(imageService *services.ImageService, authService *services.AuthService) *ImageHandler {
    return &ImageHandler{
        imageService: imageService,
        authService:  authService,
    }
}

// Validation middleware for file uploads
func (h *ImageHandler) ValidateImageUpload(c *fiber.Ctx) error {
    // Check file size (max 10MB)
    if c.Get("Content-Length") != "" {
        if contentLength := c.Get("Content-Length"); contentLength > "10485760" { // 10MB
            return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
                "error": "File too large. Maximum size is 10MB.",
            })
        }
    }

    // Check content type
    contentType := c.Get("Content-Type")
    if !strings.Contains(contentType, "multipart/form-data") {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid content type. Expected multipart/form-data.",
        })
    }

    return c.Next()
}

func (h *ImageHandler) UploadImage(c *fiber.Ctx) error {
    // Get user from context (set by auth middleware)
    userID := c.Locals("user_id").(string)

    // Parse multipart form
    form, err := c.MultipartForm()
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Failed to parse form data",
        })
    }

    files := form.File["image"]
    if len(files) == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "No image file provided",
        })
    }

    file := files[0]

    // Validate file type
    if !h.isValidImageType(file.Filename) {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed.",
        })
    }

    // Open file
    src, err := file.Open()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to open uploaded file",
        })
    }
    defer src.Close()

    // Upload to Cloudinary
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    result, err := h.imageService.UploadImage(ctx, src, file.Filename, h.imageService.config.UploadFolder, userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to upload image",
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data":    result,
    })
}

func (h *ImageHandler) UploadBatchImages(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(string)

    form, err := c.MultipartForm()
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Failed to parse form data",
        })
    }

    files := form.File["images"]
    if len(files) == 0 || len(files) > 10 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Provide 1-10 image files",
        })
    }

    results := make([]*services.UploadResult, 0, len(files))

    for _, file := range files {
        if !h.isValidImageType(file.Filename) {
            continue // Skip invalid files
        }

        src, err := file.Open()
        if err != nil {
            continue
        }

        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

        result, err := h.imageService.UploadImage(ctx, src, file.Filename, h.imageService.config.UploadFolder, userID)
        cancel()
        src.Close()

        if err == nil {
            results = append(results, result)
        }
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data":    results,
    })
}

func (h *ImageHandler) GetImageMetadata(c *fiber.Ctx) error {
    publicID := c.Params("public_id")
    if publicID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Public ID required",
        })
    }

    // Get image URL without transformations for metadata
    url := h.imageService.GetImageURL(publicID, map[string]string{})
    if url == "" {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Image not found",
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data": fiber.Map{
            "public_id": publicID,
            "url":       url,
        },
    })
}

func (h *ImageHandler) DeleteImage(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(string)
    publicID := c.Params("public_id")

    // TODO: Check ownership via database if tracking images
    // For now, allow authenticated users to delete their own images

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err := h.imageService.DeleteImage(ctx, publicID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to delete image",
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Image deleted successfully",
    })
}

func (h *ImageHandler) isValidImageType(filename string) bool {
    validTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
    filename = strings.ToLower(filename)

    for _, ext := range validTypes {
        if strings.HasSuffix(filename, ext) {
            return true
        }
    }
    return false
}

func (h *ImageHandler) RegisterRoutes(app *fiber.App) {
    api := app.Group("/api")

    // Protected routes
    protected := api.Group("/images", handlers.AuthMiddleware(h.authService))
    protected.Post("/upload", h.ValidateImageUpload, h.UploadImage)
    protected.Post("/batch-upload", h.ValidateImageUpload, h.UploadBatchImages)
    protected.Get("/:public_id", h.GetImageMetadata)
    protected.Delete("/:public_id", h.DeleteImage)
}
```

### Phase 4: Chat System Enhancement (Estimated: 3-4 hours)

#### 4.1 Extend Chat Message Upload
Update `handlers/chat_handler.go` to support image uploads:
```go
func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(string)
    channelID := c.Params("channel_id")

    // Check if this is a multipart form (image upload)
    contentType := c.Get("Content-Type")
    if strings.Contains(contentType, "multipart/form-data") {
        return h.handleImageMessage(c, userID, channelID)
    }

    // Handle text message as before
    return h.handleTextMessage(c, userID, channelID)
}

func (h *ChatHandler) handleImageMessage(c *fiber.Ctx, userID, channelID string) error {
    // Parse form
    form, err := c.MultipartForm()
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid form data",
        })
    }

    // Get image file
    files := form.File["image"]
    if len(files) == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "No image provided",
        })
    }

    file := files[0]
    content := form.Value["content"]
    messageContent := ""
    if len(content) > 0 {
        messageContent = content[0]
    }

    // Validate file type
    if !h.isValidImageType(file.Filename) {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid image type",
        })
    }

    // Open and upload image
    src, err := file.Open()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to process image",
        })
    }
    defer src.Close()

    // Upload to Cloudinary
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    result, err := h.imageService.UploadImage(ctx, src, file.Filename, "chat_images", userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to upload image",
        })
    }

    // Create chat message
    message := &models.ChatMessage{
        ChannelID:   channelID,
        UserID:      userID,
        Content:     messageContent,
        MessageType: "image",
        FileURL:     result.SecureURL,
        FileName:    file.Filename,
        FileSize:    result.Bytes,
    }

    // Add thumbnail URL if available
    if result.ThumbnailURL != "" {
        // Store thumbnail URL in a custom field or extend model
        message.Content += fmt.Sprintf("\n[thumbnail:%s]", result.ThumbnailURL)
    }

    // Save message
    err = h.chatService.CreateMessage(message)
    if err != nil {
        // Cleanup uploaded image on failure
        h.imageService.DeleteImage(context.Background(), result.PublicID)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to save message",
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": message,
    })
}

func (h *ChatHandler) EditMessage(c *fiber.Ctx) error {
    // ... existing edit logic ...

    // If message type is image and being deleted/edited, cleanup Cloudinary
    messageID := c.Params("id")
    existingMessage, err := h.chatService.GetMessage(messageID)
    if err == nil && existingMessage.MessageType == "image" && existingMessage.FileURL != "" {
        // Extract public ID from URL and delete from Cloudinary
        publicID := h.extractPublicIDFromURL(existingMessage.FileURL)
        if publicID != "" {
            h.imageService.DeleteImage(context.Background(), publicID)
        }
    }
}

func (h *ChatHandler) DeleteMessage(c *fiber.Ctx) error {
    // ... existing delete logic ...

    // Cleanup associated images
    messageID := c.Params("id")
    existingMessage, err := h.chatService.GetMessage(messageID)
    if err == nil && existingMessage.MessageType == "image" && existingMessage.FileURL != "" {
        publicID := h.extractPublicIDFromURL(existingMessage.FileURL)
        if publicID != "" {
            h.imageService.DeleteImage(context.Background(), publicID)
        }
    }
}

func (h *ChatHandler) extractPublicIDFromURL(url string) string {
    // Extract public ID from Cloudinary URL
    // URL format: https://res.cloudinary.com/{cloud_name}/image/upload/v{version}/{public_id}.{format}
    parts := strings.Split(url, "/")
    if len(parts) >= 8 {
        filename := parts[len(parts)-1]
        return strings.Split(filename, ".")[0]
    }
    return ""
}

func (h *ChatHandler) isValidImageType(filename string) bool {
    validTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
    filename = strings.ToLower(filename)

    for _, ext := range validTypes {
        if strings.HasSuffix(filename, ext) {
            return true
        }
    }
    return false
}
```

#### 4.2 Chat Image Permissions
Ensure RBAC permissions for image uploads in chat:
- `chat:upload_image` permission for image uploads
- `chat:delete_image` permission for deleting image messages
- Validate against channel membership and ownership

#### 4.3 Thumbnail Support
Extend `ChatMessage` model to include thumbnail URL:
```go
type ChatMessage struct {
    // ... existing fields ...
    ThumbnailURL *string `gorm:"type:text" json:"thumbnail_url,omitempty"`
}
```

### Phase 5: Database Extensions (Optional) (Estimated: 2-3 hours)

#### 5.1 Image Metadata Model
If needed for tracking uploaded images:
```go
type UploadedImage struct {
    ID          string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    UserID      string     `gorm:"type:uuid;index;not null;foreignKey:UserID;references:Users(ID)" json:"user_id"`
    PublicID    string     `gorm:"type:text;unique;not null" json:"public_id"`
    URL         string     `gorm:"type:text;not null" json:"url"`
    SecureURL   string     `gorm:"type:text;not null" json:"secure_url"`
    Filename    string     `gorm:"type:text;not null" json:"filename"`
    Size        int64      `gorm:"type:bigint;not null" json:"size"`
    MimeType    string     `gorm:"type:text;not null" json:"mime_type"`
    Width       int        `gorm:"type:integer" json:"width"`
    Height      int        `gorm:"type:integer" json:"height"`
    Folder      string     `gorm:"type:text;index;not null" json:"folder"`
    ThumbnailURL *string   `gorm:"type:text" json:"thumbnail_url,omitempty"`
    DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
    CreatedAt   time.Time  `gorm:"default:now()" json:"created_at"`
    UpdatedAt   time.Time  `gorm:"default:now()" json:"updated_at"`
}

func (UploadedImage) TableName() string {
    return "uploaded_images"
}
```

#### 5.2 Repository for Image Tracking
```go
type UploadedImageRepository interface {
    Create(ctx context.Context, image *models.UploadedImage) error
    GetByID(ctx context.Context, id string) (*models.UploadedImage, error)
    GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*models.UploadedImage, error)
    GetByPublicID(ctx context.Context, publicID string) (*models.UploadedImage, error)
    Delete(ctx context.Context, id string) error
    SoftDelete(ctx context.Context, id string) error // For GDPR compliance
}
```

#### 5.3 Migration
Add migration for image tracking if implemented:
```sql
-- Migration: V1_29__create_uploaded_images_table.sql
CREATE TABLE uploaded_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES gebruikers(id) ON DELETE CASCADE,
    public_id TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL,
    secure_url TEXT NOT NULL,
    filename TEXT NOT NULL,
    size BIGINT NOT NULL,
    mime_type TEXT NOT NULL,
    width INTEGER,
    height INTEGER,
    folder TEXT NOT NULL,
    thumbnail_url TEXT,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_uploaded_images_user_id ON uploaded_images(user_id);
CREATE INDEX idx_uploaded_images_folder ON uploaded_images(folder);
CREATE INDEX idx_uploaded_images_deleted_at ON uploaded_images(deleted_at);
CREATE INDEX idx_uploaded_images_created_at ON uploaded_images(created_at DESC);
```

### Phase 6: Frontend Integration (Estimated: 4-6 hours)

#### 6.1 API Endpoints for Frontend
```http
POST /api/images/upload
Content-Type: multipart/form-data
Authorization: Bearer <jwt_token>

Response:
{
    "success": true,
    "data": {
        "public_id": "dkl_images/user123_abc123def",
        "url": "https://res.cloudinary.com/.../image.jpg",
        "secure_url": "https://res.cloudinary.com/.../image.jpg",
        "width": 800,
        "height": 600,
        "format": "jpg",
        "bytes": 245760,
        "thumbnail_url": "https://res.cloudinary.com/.../thumbnail.jpg"
    }
}

Error Responses:
413 Payload Too Large: {"error": "File too large. Maximum size is 10MB."}
400 Bad Request: {"error": "Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed."}
401 Unauthorized: {"error": "Authentication required"}
500 Internal Server Error: {"error": "Failed to upload image"}
```

#### 6.2 Image Upload with Progress Feedback
```javascript
class ImageUploader {
    constructor(apiBaseUrl, authToken) {
        this.apiBaseUrl = apiBaseUrl;
        this.authToken = authToken;
    }

    async uploadImage(file, onProgress) {
        const formData = new FormData();
        formData.append('image', file);

        return new Promise((resolve, reject) => {
            const xhr = new XMLHttpRequest();

            xhr.upload.addEventListener('progress', (event) => {
                if (event.lengthComputable && onProgress) {
                    const percentComplete = (event.loaded / event.total) * 100;
                    onProgress(percentComplete);
                }
            });

            xhr.addEventListener('load', () => {
                if (xhr.status >= 200 && xhr.status < 300) {
                    try {
                        const response = JSON.parse(xhr.responseText);
                        resolve(response);
                    } catch (e) {
                        reject(new Error('Invalid response format'));
                    }
                } else {
                    try {
                        const error = JSON.parse(xhr.responseText);
                        reject(new Error(error.error || 'Upload failed'));
                    } catch (e) {
                        reject(new Error(`Upload failed with status ${xhr.status}`));
                    }
                }
            });

            xhr.addEventListener('error', () => {
                reject(new Error('Network error during upload'));
            });

            xhr.open('POST', `${this.apiBaseUrl}/images/upload`);
            xhr.setRequestHeader('Authorization', `Bearer ${this.authToken}`);
            xhr.send(formData);
        });
    }

    async uploadBatchImages(files, onProgress) {
        const formData = new FormData();
        files.forEach((file, index) => {
            formData.append('images', file);
        });

        // Similar XMLHttpRequest implementation for batch upload
        // with progress tracking
    }
}

// Usage example
const uploader = new ImageUploader('/api', authToken);

try {
    const result = await uploader.uploadImage(file, (progress) => {
        console.log(`Upload progress: ${progress}%`);
        // Update UI progress bar
    });

    // Use result.data.secure_url for displaying the image
    console.log('Upload successful:', result.data);
} catch (error) {
    console.error('Upload failed:', error.message);
    // Show error to user
}
```

#### 6.3 Chat Image Upload
Extend chat message sending to support file uploads:
```javascript
async function sendChatImage(channelId, file, message = '') {
    const formData = new FormData();
    formData.append('image', file);
    if (message) {
        formData.append('content', message);
    }

    const response = await fetch(`/api/chat/channels/${channelId}/messages`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${authToken}`
        },
        body: formData
    });

    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Failed to send image');
    }

    return await response.json();
}

// Usage in chat component
const handleImageSend = async (file) => {
    try {
        setUploading(true);
        const result = await sendChatImage(channelId, file, 'Check this out!');

        // Add message to chat UI
        addMessageToChat(result.message);

        // Image URLs are available in result.message.file_url and thumbnail_url
    } catch (error) {
        showError('Failed to send image: ' + error.message);
    } finally {
        setUploading(false);
    }
};
```

#### 6.4 Responsive Images with Cloudinary
```javascript
// Generate responsive image URLs
function getResponsiveImageUrl(publicId, options = {}) {
    const baseUrl = `https://res.cloudinary.com/${cloudName}/image/upload/`;

    const transformations = [];

    // Auto format and quality
    transformations.push('f_auto', 'q_auto');

    // Responsive width with device pixel ratio
    if (options.width) {
        transformations.push(`w_${options.width}`);
    }
    if (options.dpr) {
        transformations.push(`dpr_${options.dpr}`);
    } else {
        transformations.push('dpr_auto'); // Auto DPR
    }

    // Crop if specified
    if (options.crop) {
        transformations.push(`c_${options.crop}`);
    }

    return `${baseUrl}${transformations.join(',')}/${publicId}`;
}

// React component example
const ResponsiveImage = ({ publicId, alt, className, sizes = '(max-width: 768px) 100vw, 50vw' }) => {
    const src = getResponsiveImageUrl(publicId);
    const srcSet = `
        ${getResponsiveImageUrl(publicId, { width: 480, dpr: 1 })} 480w,
        ${getResponsiveImageUrl(publicId, { width: 768, dpr: 1 })} 768w,
        ${getResponsiveImageUrl(publicId, { width: 1024, dpr: 1 })} 1024w,
        ${getResponsiveImageUrl(publicId, { width: 480, dpr: 2 })} 960w,
        ${getResponsiveImageUrl(publicId, { width: 768, dpr: 2 })} 1536w
    `;

    return (
        <img
            src={src}
            srcSet={srcSet}
            sizes={sizes}
            alt={alt}
            className={className}
            loading="lazy"
        />
    );
};
```

### Phase 7: Security & Performance (Estimated: 4-6 hours)

#### 7.1 Security Measures
- **File Validation**: Check file types, sizes (10MB limit), dimensions
- **Malware Scanning**: Enable Cloudinary's built-in virus scanning via upload presets
- **Access Control**: JWT authentication required for all upload operations
- **Rate Limiting**: Upload rate limits (10 per hour per user, 100 per hour global)
- **CORS**: Restrict upload origins to configured frontend domains
- **Secure URLs**: Use HTTPS URLs with signed requests for private images
- **Input Sanitization**: Validate filenames and metadata

#### 7.2 Signed URLs for Private Images
```go
// For private/sensitive images
func (s *ImageService) GetSignedImageURL(publicID string, expiry time.Duration) string {
    img, err := s.cld.Image(publicID)
    if err != nil {
        return ""
    }

    // Create signed URL with expiration
    signedURL, err := img.Signed(true).ToURL()
    if err != nil {
        return ""
    }

    return signedURL
}
```

#### 7.3 Webhook Integration
```go
// Handle Cloudinary webhooks for upload notifications
func (h *ImageHandler) HandleCloudinaryWebhook(c *fiber.Ctx) error {
    // Verify webhook signature
    signature := c.Get("X-Cld-Signature")
    timestamp := c.Get("X-Cld-Timestamp")
    body := c.Body()

    if !h.verifyWebhookSignature(signature, timestamp, body) {
        return c.Status(fiber.StatusUnauthorized).SendString("Invalid signature")
    }

    var webhookData struct {
        NotificationType string `json:"notification_type"`
        PublicID         string `json:"public_id"`
        UploadResult     struct {
            PublicID string `json:"public_id"`
            URL      string `json:"url"`
            ModerationStatus string `json:"moderation_status,omitempty"`
        } `json:"upload_result,omitempty"`
    }

    if err := c.BodyParser(&webhookData); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid webhook data")
    }

    // Handle different webhook types
    switch webhookData.NotificationType {
    case "upload":
        // Update database with upload completion
        h.handleUploadCompletion(webhookData.UploadResult)
    case "moderation":
        // Handle content moderation results
        h.handleModerationResult(webhookData.UploadResult)
    }

    return c.SendString("OK")
}

func (h *ImageHandler) verifyWebhookSignature(signature, timestamp string, body []byte) bool {
    // Implement signature verification using Cloudinary API secret
    // Follow Cloudinary webhook security documentation
    return true // Placeholder
}

func (h *ImageHandler) RegisterWebhookRoutes(app *fiber.App) {
    app.Post("/webhooks/cloudinary", h.HandleCloudinaryWebhook)
}
```

#### 7.4 CDN Purging and Cache Management
```go
func (s *ImageService) PurgeImageCache(ctx context.Context, publicID string) error {
    // Purge specific image from CDN cache
    _, err := s.cld.Admin.PurgeAsset(ctx, admin.PurgeAssetParams{
        PublicID: publicID,
        ResourceType: "image",
    })
    return err
}

func (s *ImageService) InvalidateTransformations(ctx context.Context, publicID string) error {
    // Invalidate derived images (thumbnails, transformations)
    _, err := s.cld.Admin.PurgeAsset(ctx, admin.PurgeAssetParams{
        PublicID: publicID,
        ResourceType: "image",
        Invalidate: true,
    })
    return err
}
```

#### 7.5 Performance Optimizations
- **Image Optimization**: Auto-format (f_auto), quality (q_auto), and responsive images
- **CDN Caching**: Leverage Cloudinary's global CDN with 25,000+ edge locations
- **Async Uploads**: Background processing for large files with webhook notifications
- **Progress Tracking**: Real-time upload progress via XMLHttpRequest
- **Lazy Loading**: Frontend lazy loading for image galleries
- **Preloading**: Critical images preloaded for better UX
- **Compression**: Automatic WebP/AVIF conversion for modern browsers

### Phase 8: Monitoring & Logging (Estimated: 3-4 hours)

#### 8.1 Prometheus Metrics
Add to `services/prometheus_metrics.go`:
```go
// Image upload metrics
var (
	imageUploadsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "image_uploads_total",
			Help: "Total number of image uploads",
		},
		[]string{"status", "user_type"}, // success/failure, authenticated/anonymous
	)

	imageUploadDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "image_upload_duration_seconds",
			Help:    "Time taken for image uploads",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"size_range"}, // small/medium/large
	)

	imageStorageBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "image_storage_bytes_total",
			Help: "Total image storage used in bytes",
		},
	)

	imageBandwidthBytes = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "image_bandwidth_bytes_total",
			Help: "Total bandwidth used for image delivery",
		},
	)
)

// Cloudinary-specific metrics
var (
	cloudinaryAPICallsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudinary_api_calls_total",
			Help: "Total Cloudinary API calls",
		},
		[]string{"endpoint", "status"}, // upload/delete/transform
	)

	cloudinaryAPILatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cloudinary_api_latency_seconds",
			Help:    "Cloudinary API call latency",
			Buckets: []float64{.1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"endpoint"},
	)

	cloudinaryQuotaUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cloudinary_quota_usage_percent",
			Help: "Cloudinary quota usage percentage",
		},
		[]string{"resource"}, // storage, bandwidth, transformations
	)
)

func init() {
	prometheus.MustRegister(
		imageUploadsTotal,
		imageUploadDuration,
		imageStorageBytes,
		imageBandwidthBytes,
		cloudinaryAPICallsTotal,
		cloudinaryAPILatency,
		cloudinaryQuotaUsage,
	)
}
```

#### 8.2 Cloudinary Usage Monitoring
```go
// Periodic quota checking
func (s *ImageService) monitorCloudinaryUsage(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.updateCloudinaryMetrics(ctx)
		}
	}
}

func (s *ImageService) updateCloudinaryMetrics(ctx context.Context) {
	// Get usage statistics from Cloudinary Admin API
	usage, err := s.cld.Admin.Usage(ctx, admin.UsageParams{})
	if err != nil {
		logger.Warn("Failed to get Cloudinary usage", "error", err)
		return
	}

	// Update Prometheus metrics
	if usage.Storage != nil {
		imageStorageBytes.Set(float64(usage.Storage.Used))
		if usage.Storage.Limit > 0 {
			usagePercent := (float64(usage.Storage.Used) / float64(usage.Storage.Limit)) * 100
			cloudinaryQuotaUsage.WithLabelValues("storage").Set(usagePercent)
		}
	}

	if usage.Bandwidth != nil {
		imageBandwidthBytes.Reset() // Reset counter for new period
		// Note: Admin API may not provide cumulative bandwidth
	}
}
```

#### 8.3 Structured Logging
Add to `services/image_service.go`:
```go
func (s *ImageService) UploadImage(ctx context.Context, file multipart.File, filename string, folder string, userID string) (*UploadResult, error) {
	start := time.Now()

	result, err := s.uploadImageInternal(ctx, file, filename, folder, userID)

	duration := time.Since(start)
	sizeRange := s.getSizeRange(result.Bytes)

	// Update metrics
	if err != nil {
		imageUploadsTotal.WithLabelValues("failure", "authenticated").Inc()
		logger.Error("Image upload failed",
			"error", err,
			"user_id", userID,
			"filename", filename,
			"duration_ms", duration.Milliseconds(),
		)
	} else {
		imageUploadsTotal.WithLabelValues("success", "authenticated").Inc()
		imageUploadDuration.WithLabelValues(sizeRange).Observe(duration.Seconds())

		logger.Info("Image uploaded successfully",
			"public_id", result.PublicID,
			"user_id", userID,
			"filename", filename,
			"size_bytes", result.Bytes,
			"duration_ms", duration.Milliseconds(),
		)
	}

	return result, err
}

func (s *ImageService) getSizeRange(bytes int) string {
	switch {
	case bytes < 100000: // < 100KB
		return "small"
	case bytes < 1000000: // < 1MB
		return "medium"
	default:
		return "large"
	}
}
```

#### 8.4 Alerting Configuration
Add alerting rules for Grafana/Prometheus Alertmanager:
```yaml
# Alert on high upload failure rate
- alert: HighImageUploadFailureRate
  expr: rate(image_uploads_total{status="failure"}[5m]) / rate(image_uploads_total[5m]) > 0.1
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "High image upload failure rate"
    description: "Image upload failure rate is {{ $value | printf "%.2f" }}%"

# Alert on Cloudinary quota usage
- alert: CloudinaryQuotaNearLimit
  expr: cloudinary_quota_usage_percent > 90
  for: 15m
  labels:
    severity: warning
  annotations:
    summary: "Cloudinary quota near limit"
    description: "{{ $labels.resource }} usage is {{ $value | printf "%.1f" }}%"

# Alert on slow uploads
- alert: SlowImageUploads
  expr: histogram_quantile(0.95, rate(image_upload_duration_seconds_bucket[5m])) > 30
  for: 10m
  labels:
    severity: warning
  annotations:
    summary: "Slow image uploads detected"
    description: "95th percentile upload time is {{ $value | printf "%.1f" }}s"
```

#### 8.5 Audit Logging
```go
type AuditEvent struct {
	Timestamp   time.Time `json:"timestamp"`
	UserID      string    `json:"user_id"`
	Action      string    `json:"action"` // upload, delete, view
	ResourceID  string    `json:"resource_id"`
	ResourceType string   `json:"resource_type"` // image, chat_message
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

func (s *ImageService) logAuditEvent(event AuditEvent) {
	logger.Info("Audit event",
		"user_id", event.UserID,
		"action", event.Action,
		"resource_id", event.ResourceID,
		"resource_type", event.ResourceType,
		"ip_address", event.IPAddress,
	)
	// Could also write to separate audit log file or database
}
```

### Phase 9: Testing (Estimated: 6-8 hours)

#### 9.1 Unit Tests
```go
// services/image_service_test.go
func TestImageService_UploadImage(t *testing.T) {
    // Mock Cloudinary client
    mockCld := &mocks.Cloudinary{}
    mockUpload := &mocks.Uploader{}

    service := &ImageService{
        cld: mockCld,
        config: &config.CloudinaryConfig{
            CloudName: "test",
            UploadFolder: "test_folder",
        },
    }

    // Mock successful upload
    mockUpload.On("Upload", mock.Anything, mock.AnythingOfType("uploader.UploadParams")).
        Return(&uploader.UploadResult{
            PublicID: "test_public_id",
            URL: "https://test.cloudinary.com/test.jpg",
            SecureURL: "https://test.cloudinary.com/test.jpg",
            Width: 800,
            Height: 600,
            Format: "jpg",
            Bytes: 102400,
        }, nil)

    mockCld.On("Upload").Return(mockUpload)

    result, err := service.UploadImage(context.Background(), file, "test.jpg", "test_folder", "user123")

    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "test_public_id", result.PublicID)
}

func TestImageService_UploadImage_InvalidFile(t *testing.T) {
    service := &ImageService{
        config: &config.CloudinaryConfig{},
    }

    _, err := service.UploadImage(context.Background(), nil, "", "test", "user123")

    assert.Error(t, err)
    assert.Equal(t, ErrInvalidFile, err)
}
```

#### 9.2 Mocking Cloudinary for Tests
```go
// Use testify/mock or create interfaces for testability
type CloudinaryClient interface {
    Upload() Uploader
    Image(publicID string) (Image, error)
    Admin() Admin
}

type Uploader interface {
    Upload(ctx context.Context, file io.Reader, params uploader.UploadParams) (*uploader.UploadResult, error)
    Destroy(ctx context.Context, params uploader.DestroyParams) (*uploader.DestroyResult, error)
}

// Test with real Cloudinary but test credentials
func createTestImageService() *ImageService {
    config := &config.CloudinaryConfig{
        CloudName:   os.Getenv("CLOUDINARY_TEST_CLOUD_NAME"),
        APIKey:      os.Getenv("CLOUDINARY_TEST_API_KEY"),
        APISecret:   os.Getenv("CLOUDINARY_TEST_API_SECRET"),
        UploadFolder: "test_uploads",
        IsTest:      true,
    }

    service, _ := NewImageService(config)
    return service
}
```

#### 9.3 Handler Tests
```go
// handlers/image_handler_test.go
func TestImageHandler_UploadImage(t *testing.T) {
    app := fiber.New()

    mockService := &mocks.ImageService{}
    mockAuth := &mocks.AuthService{}

    handler := NewImageHandler(mockService, mockAuth)

    // Mock successful upload
    mockService.On("UploadImage", mock.Anything, mock.Anything, mock.Anything, mock.Anything, "user123").
        Return(&UploadResult{
            PublicID: "test_id",
            URL: "https://test.com/image.jpg",
        }, nil)

    app.Post("/upload", handler.ValidateImageUpload, handler.UploadImage)

    // Create test request
    req := httptest.NewRequest("POST", "/upload", createMultipartBody())
    req.Header.Set("Content-Type", "multipart/form-data")
    req.Header.Set("Authorization", "Bearer test_token")

    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
}

func TestImageHandler_UploadImage_FileTooLarge(t *testing.T) {
    app := fiber.New()
    handler := NewImageHandler(nil, nil)

    app.Post("/upload", handler.ValidateImageUpload, func(c *fiber.Ctx) error {
        return c.SendString("Should not reach here")
    })

    // Create request with large content-length
    req := httptest.NewRequest("POST", "/upload", nil)
    req.Header.Set("Content-Length", "11534336") // 11MB > 10MB limit
    req.Header.Set("Content-Type", "multipart/form-data")

    resp, err := app.Test(req)
    assert.NoError(t, err)
    assert.Equal(t, 413, resp.StatusCode) // Payload Too Large
}
```

#### 9.4 Integration Tests
```go
// tests/image_integration_test.go
func TestImageUploadWorkflow(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    // Setup test database and services
    db := setupTestDB(t)
    defer db.Close()

    service := createTestImageService()
    handler := NewImageHandler(service, createTestAuthService())

    // Create test file
    file, err := os.Open("testdata/test_image.jpg")
    require.NoError(t, err)
    defer file.Close()

    // Test upload
    result, err := service.UploadImage(context.Background(), file, "test.jpg", "integration_test", "test_user")
    require.NoError(t, err)
    assert.NotEmpty(t, result.PublicID)
    assert.Contains(t, result.URL, "cloudinary")

    // Test retrieval
    url := service.GetImageURL(result.PublicID, map[string]string{"width": "300"})
    assert.Contains(t, url, "w_300")

    // Test deletion
    err = service.DeleteImage(context.Background(), result.PublicID)
    assert.NoError(t, err)
}
```

#### 9.5 End-to-End Tests
```javascript
// tests/e2e/image_upload.spec.js (using Playwright)
const { test, expect } = require('@playwright/test');

test('should upload image successfully', async ({ page }) => {
  // Login
  await page.goto('/login');
  await page.fill('[data-testid="email"]', 'test@example.com');
  await page.fill('[data-testid="password"]', 'password');
  await page.click('[data-testid="login-button"]');

  // Navigate to upload page
  await page.goto('/upload');

  // Upload file
  const fileInput = page.locator('input[type="file"]');
  await fileInput.setInputFiles('./testdata/test_image.jpg');

  // Wait for upload completion
  await page.waitForSelector('[data-testid="upload-success"]');

  // Verify image is displayed
  const image = page.locator('[data-testid="uploaded-image"]');
  await expect(image).toBeVisible();

  // Verify image URL is from Cloudinary
  const src = await image.getAttribute('src');
  expect(src).toContain('res.cloudinary.com');
});

test('should show progress during upload', async ({ page }) => {
  // Similar setup, then check progress bar
  const progressBar = page.locator('[data-testid="upload-progress"]');
  await expect(progressBar).toBeVisible();

  // Wait for completion
  await page.waitForSelector('[data-testid="progress-complete"]');
});
```

#### 9.6 Load Testing
```javascript
// tests/load/image_load_test.js (using k6)
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '2m', target: 10 },  // Ramp up to 10 users
    { duration: '5m', target: 50 },  // Ramp up to 50 users
    { duration: '2m', target: 100 }, // Ramp up to 100 users
    { duration: '5m', target: 100 }, // Stay at 100 users
    { duration: '2m', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
    http_req_failed: ['rate<0.1'],    // Error rate should be below 10%
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TOKEN = __ENV.AUTH_TOKEN;

export default function () {
  const imageData = open('./testdata/large_image.jpg', 'b');

  const formData = {
    image: http.file(imageData, 'test_image.jpg'),
  };

  const response = http.post(`${BASE_URL}/api/images/upload`, formData, {
    headers: {
      'Authorization': `Bearer ${TOKEN}`,
    },
  });

  check(response, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
    'contains image URL': (r) => r.json().data.url.includes('cloudinary'),
  });

  sleep(1);
}
```

#### 9.7 Test Coverage Goals
- **Unit Tests**: >90% coverage for ImageService and ImageHandler
- **Integration Tests**: Full upload-to-delete workflow
- **E2E Tests**: Frontend upload experience
- **Load Tests**: 100 concurrent uploads, <500ms p95 response time

### Phase 10: Deployment & Configuration (Estimated: 4-6 hours)

#### 10.1 Environment Setup
- Configure Cloudinary credentials in production environment
- Set up dedicated upload folders per environment (prod/dev/staging)
- Configure upload presets for optimization and security
- Set up webhook endpoints for async notifications
- Configure monitoring dashboards

#### 10.2 CI/CD Integration
```yaml
# .github/workflows/deploy.yml
name: Deploy with Cloudinary
on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Test with Cloudinary mocks
        run: go test ./... -v
        env:
          CLOUDINARY_TEST_CLOUD_NAME: ${{ secrets.CLOUDINARY_TEST_CLOUD_NAME }}
          CLOUDINARY_TEST_API_KEY: ${{ secrets.CLOUDINARY_TEST_API_KEY }}
          CLOUDINARY_TEST_API_SECRET: ${{ secrets.CLOUDINARY_TEST_API_SECRET }}

      - name: Build
        run: go build -ldflags="-s -w" -o app

  deploy-staging:
    needs: test
    runs-on: ubuntu-latest
    environment: staging
    steps:
      - name: Deploy to staging
        run: |
          echo "Deploying to staging with Cloudinary integration"
          # Deployment commands here

  deploy-production:
    needs: deploy-staging
    runs-on: ubuntu-latest
    environment: production
    steps:
      - name: Deploy to production
        run: |
          echo "Deploying to production"
          # Set production Cloudinary credentials
        env:
          CLOUDINARY_CLOUD_NAME: ${{ secrets.CLOUDINARY_CLOUD_NAME }}
          CLOUDINARY_API_KEY: ${{ secrets.CLOUDINARY_API_KEY }}
          CLOUDINARY_API_SECRET: ${{ secrets.CLOUDINARY_API_SECRET }}
```

#### 10.3 Feature Flags for Gradual Rollout
```go
// Add to config
type FeatureFlags struct {
    EnableCloudinary    bool `json:"enable_cloudinary"`
    EnableImageUploads  bool `json:"enable_image_uploads"`
    EnableChatImages    bool `json:"enable_chat_images"`
    EnableImageTracking bool `json:"enable_image_tracking"`
}

// Usage in handlers
func (h *ImageHandler) UploadImage(c *fiber.Ctx) error {
    if !config.FeatureFlags.EnableImageUploads {
        return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
            "error": "Image uploads are temporarily disabled",
        })
    }
    // ... rest of implementation
}
```

#### 10.4 Rollback Plan
**Immediate Rollback (Feature Flag):**
```bash
# Disable Cloudinary features instantly
export ENABLE_CLOUDINARY=false
export ENABLE_IMAGE_UPLOADS=false
export ENABLE_CHAT_IMAGES=false
# Restart application
```

**Full Rollback (Code Revert):**
```bash
# Git revert to pre-Cloudinary commit
git revert <cloudinary-integration-commit>
git push origin main

# Clean up Cloudinary resources
# - Delete uploaded images via Cloudinary dashboard
# - Remove database records if image tracking was enabled
# - Update frontend to remove image upload UI
```

**Partial Rollback Scenarios:**
- **High Error Rates**: Disable uploads, keep existing images accessible
- **Cost Issues**: Reduce upload limits, implement stricter validation
- **Performance Issues**: Switch to local storage temporarily

#### 10.5 Migration Strategy
- **Blue-Green Deployment**: Deploy new version alongside old, test thoroughly
- **Canary Release**: Roll out to 10% of users first, monitor metrics
- **Database Migration**: Run image tracking migration only after successful testing
- **Data Migration**: For existing chat messages with image URLs, update to Cloudinary format

#### 10.6 Cost Monitoring and Optimization
```go
// services/cost_monitor.go
type CostMonitor struct {
    cloudinaryClient *cloudinary.Cloudinary
    alertThreshold   float64 // e.g., 80% of monthly quota
}

func (cm *CostMonitor) CheckUsageAndAlert(ctx context.Context) error {
    usage, err := cm.cloudinaryClient.Admin.Usage(ctx, admin.UsageParams{})
    if err != nil {
        return err
    }

    // Check storage usage
    if usage.Storage != nil && usage.Storage.Limit > 0 {
        usagePercent := (float64(usage.Storage.Used) / float64(usage.Storage.Limit)) * 100
        if usagePercent > cm.alertThreshold {
            cm.sendAlert("Storage usage", usagePercent)
        }
    }

    // Check bandwidth usage
    if usage.Bandwidth != nil && usage.Bandwidth.Limit > 0 {
        usagePercent := (float64(usage.Bandwidth.Used) / float64(usage.Bandwidth.Limit)) * 100
        if usagePercent > cm.alertThreshold {
            cm.sendAlert("Bandwidth usage", usagePercent)
        }
    }

    return nil
}

func (cm *CostMonitor) sendAlert(resourceType string, usagePercent float64) {
    logger.Warn("Cloudinary quota alert",
        "resource", resourceType,
        "usage_percent", usagePercent,
        "threshold", cm.alertThreshold,
    )

    // Send email alert, Slack notification, etc.
    // Could integrate with services/notification_service.go
}
```

**Cost Estimation (Monthly):**
- **Free Tier**: 25GB storage, 25GB monthly bandwidth
- **Expected Usage**: 1000 images/month = ~50GB storage, ~200GB bandwidth
- **Cost**: ~$50-100/month depending on transformations
- **Optimization**: Use auto-format, appropriate quality settings, cache headers

#### 10.7 Production Checklist
- [ ] Cloudinary account created and configured
- [ ] Upload presets created for security/optimization
- [ ] Webhook endpoints configured and tested
- [ ] Environment variables set in production
- [ ] Feature flags configured for gradual rollout
- [ ] Monitoring dashboards set up
- [ ] Cost alerts configured
- [ ] Rollback procedures documented
- [ ] Team trained on new functionality
- [ ] Frontend updated and deployed

## Implementation Order

1. **Phase 1**: Infrastructure setup (2-4 hours) - Dependencies, config, environment setup
2. **Phase 2**: Core image service implementation (4-6 hours) - ImageService with upload/delete/URL generation
3. **Phase 3**: API endpoints (3-5 hours) - Handlers, validation, batch uploads
4. **Phase 4**: Chat system enhancement (3-4 hours) - Image message support, thumbnails
5. **Phase 5**: Database extensions (2-3 hours) - Optional image tracking model
6. **Phase 6**: Frontend integration (4-6 hours) - Upload components, progress feedback, responsive images
7. **Phase 7**: Security and performance hardening (4-6 hours) - Malware scanning, signed URLs, webhooks
8. **Phase 8**: Monitoring and logging (3-4 hours) - Prometheus metrics, Cloudinary monitoring
9. **Phase 9**: Comprehensive testing (6-8 hours) - Unit, integration, E2E, load tests
10. **Phase 10**: Production deployment (4-6 hours) - CI/CD, feature flags, cost monitoring

**Total Estimated Time**: 34-52 hours

**Milestones**:
- **Week 1**: Phases 1-3 (Infrastructure and basic upload functionality)
- **Week 2**: Phases 4-6 (Chat integration and frontend)
- **Week 3**: Phases 7-8 (Security, monitoring, testing)
- **Week 4**: Phase 9-10 (Production deployment and optimization)

**Code Review Points**: Review after each phase completion before proceeding to next phase.

## Risk Assessment

### Technical Risks
- **Cloudinary API Limits**: Monitor usage and costs with automated alerts
- **Network Dependencies**: Handle Cloudinary outages gracefully with retry logic and fallbacks
- **File Size Limits**: Implement client and server-side validation (10MB limit)
- **SDK Breaking Changes**: Pin to specific version (v2.3.0) and monitor for updates
- **Concurrent Upload Limits**: Implement queuing for high-traffic scenarios

### Security Risks
- **Unauthorized Uploads**: JWT authentication required for all operations
- **Malicious Files**: File type validation, size limits, and Cloudinary's built-in malware scanning
- **Data Exposure**: Secure credential management with separate test/prod environments
- **Credential Leaks**: Environment variables with proper secret management
- **Content Moderation**: Implement checks for inappropriate content

### Operational Risks
- **Cost Management**: Monitor Cloudinary usage costs with automated alerts and quotas
- **Performance Impact**: CDN optimization for global users with proper caching strategies
- **Downtime Handling**: Fallback mechanisms for upload failures and graceful degradation
- **Vendor Lock-in**: Design with interfaces for potential provider migration
- **Legal Risks**: Copyright infringement monitoring and GDPR compliance for user-uploaded content
- **Internationalization**: Multi-language error messages and regional CDN optimization

### Business Risks
- **Budget Overruns**: Cost estimation and monitoring to stay within budget
- **Service Availability**: Dependency on Cloudinary uptime (99.9% SLA)
- **Data Privacy**: Compliance with GDPR for user-uploaded images and metadata
- **Scalability**: Performance testing for expected user growth

## Success Criteria

- ✅ Images can be uploaded via API endpoints with progress feedback
- ✅ Chat system supports image messages with thumbnail generation
- ✅ Images are delivered via Cloudinary CDN with global optimization
- ✅ Proper authentication and authorization with RBAC integration
- ✅ Image optimization and responsive delivery (auto-format, quality, DPR)
- ✅ Comprehensive error handling and logging with structured events
- ✅ Performance meets application requirements (<500ms p95 upload time)
- ✅ Security standards maintained with malware scanning and signed URLs
- ✅ 90%+ test coverage for ImageService and ImageHandler
- ✅ Cost monitoring and alerts for quota management
- ✅ Zero unauthorized uploads in security testing
- ✅ Successful load testing (100 concurrent uploads)
- ✅ GDPR compliance for user-uploaded content
- ✅ Smooth user experience with progress indicators and error feedback

## Next Steps

1. Review and approve this integration plan
2. Set up Cloudinary account and obtain credentials
3. Begin implementation with Phase 1 (Infrastructure Setup)
4. Implement and test incrementally
5. Deploy to staging environment for frontend integration testing
6. Full production deployment with monitoring