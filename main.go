package main

import (
	"dklautomationgo/handlers"
	"dklautomationgo/services"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
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
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("ALLOWED_ORIGINS"),
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,OPTIONS",
	}))

	// Health check endpoint
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	// Routes
	app.Post("/contact-email", emailHandler.HandleContactEmail)
	// Temporarily disabled until templates are ready
	// app.Post("/aanmelding-email", emailHandler.HandleAanmeldingEmail)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
