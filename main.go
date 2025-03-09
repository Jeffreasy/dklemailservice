package main

import (
	"dklautomationgo/handlers"
	"dklautomationgo/services"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables")
	}

	// Initialize services
	emailService, err := services.NewEmailService()
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}

	// Initialize handlers
	emailHandler := handlers.NewEmailHandler(emailService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Printf("Error handling request: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Er is een fout opgetreden bij het verwerken van je verzoek",
			})
		},
	})

	// Configure CORS
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	if len(allowedOrigins) == 0 || (len(allowedOrigins) == 1 && allowedOrigins[0] == "") {
		allowedOrigins = []string{"https://www.dekoninklijkeloop.nl", "https://dekoninklijkeloop.nl"}
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(allowedOrigins, ","),
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,OPTIONS",
	}))

	// Root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "DKL Email Service API",
			"version": "1.0.0",
		})
	})

	// API routes group
	api := app.Group("/api")

	// Health check endpoint
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	// Email routes
	api.Post("/contact-email", emailHandler.HandleContactEmail)
	api.Post("/aanmelding-email", emailHandler.HandleAanmeldingEmail)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 for web traffic
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
