package handlers

import (
	"dklautomationgo/models"
	"dklautomationgo/services"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

type EmailHandler struct {
	emailService *services.EmailService
}

func NewEmailHandler(emailService *services.EmailService) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

func (h *EmailHandler) HandleContactEmail(c *fiber.Ctx) error {
	var contact models.ContactFormulier
	if err := c.BodyParser(&contact); err != nil {
		log.Printf("Error parsing contact form: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	log.Printf("Sending admin email to: %s", os.Getenv("ADMIN_EMAIL"))

	// Stuur email naar admin
	adminEmailData := &models.ContactEmailData{
		ToAdmin:    true,
		Contact:    &contact,
		AdminEmail: os.Getenv("ADMIN_EMAIL"),
	}
	if err := h.emailService.SendContactEmail(adminEmailData); err != nil {
		log.Printf("Error sending admin email: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send admin notification",
		})
	}

	log.Printf("Successfully sent admin email to: %s", os.Getenv("ADMIN_EMAIL"))

	// Stuur bevestigingsemail naar gebruiker
	userEmailData := &models.ContactEmailData{
		ToAdmin: false,
		Contact: &contact,
	}
	if err := h.emailService.SendContactEmail(userEmailData); err != nil {
		log.Printf("Error sending user email: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send confirmation email",
		})
	}

	log.Printf("Successfully sent confirmation email to: %s", contact.Email)

	return c.JSON(fiber.Map{
		"message": "Emails sent successfully",
	})
}

// Temporarily disabled until templates are ready
/*
func (h *EmailHandler) HandleAanmeldingEmail(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error": "Registration email service temporarily disabled",
	})
}
*/
