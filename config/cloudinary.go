package config

import (
	"os"

	dkllogger "dklautomationgo/logger"
)

// CloudinaryConfig bevat alle Cloudinary configuratie
type CloudinaryConfig struct {
	CloudName    string
	APIKey       string
	APISecret    string
	UploadFolder string
	UploadPreset string
	Secure       bool
	IsTest       bool
}

// LoadCloudinaryConfig laadt Cloudinary configuratie uit environment variables
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
		dkllogger.Info("Cloudinary configuratie niet gevonden, Cloudinary features zijn uitgeschakeld")
		return nil // Cloudinary not configured
	}

	uploadFolder := getEnvOrDefault("CLOUDINARY_UPLOAD_FOLDER", "dkl_images")
	uploadPreset := os.Getenv("CLOUDINARY_UPLOAD_PRESET")
	secure := getEnvOrDefault("CLOUDINARY_SECURE", "true") == "true"

	dkllogger.Info("Cloudinary configuratie geladen",
		"cloud_name", cloudName,
		"upload_folder", uploadFolder,
		"upload_preset", uploadPreset != "",
		"secure", secure,
		"is_test", isTest)

	return &CloudinaryConfig{
		CloudName:    cloudName,
		APIKey:       apiKey,
		APISecret:    apiSecret,
		UploadFolder: uploadFolder,
		UploadPreset: uploadPreset,
		Secure:       secure,
		IsTest:       isTest,
	}
}

// getEnvOrDefault haalt een omgevingsvariabele op met een standaardwaarde
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
