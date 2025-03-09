package services

import (
	"bytes"
	"dklautomationgo/models"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

const emailStyle = `
<style>
	body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
	.container { max-width: 600px; margin: 0 auto; padding: 20px; }
	h2 { color: #2c5282; margin-bottom: 20px; }
	h3 { color: #2d3748; margin-top: 20px; }
	.info-item { margin: 10px 0; }
	.info-label { font-weight: bold; color: #4a5568; }
	ul { list-style-type: none; padding-left: 0; }
	li { margin: 10px 0; }
	.footer { margin-top: 30px; padding-top: 20px; border-top: 1px solid #e2e8f0; }
	.empty-field { color: #718096; font-style: italic; }
</style>`

type EmailService struct {
	templates map[string]*template.Template
}

func NewEmailService() (*EmailService, error) {
	templates := make(map[string]*template.Template)

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %v", err)
	}

	// Load contact email templates
	contactAdminTemplate, err := template.ParseFiles(fmt.Sprintf("%s/templates/contact_admin_email.html", cwd))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contact admin template: %v", err)
	}
	templates["contact_admin"] = contactAdminTemplate

	contactUserTemplate, err := template.ParseFiles(fmt.Sprintf("%s/templates/contact_email.html", cwd))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contact user template: %v", err)
	}
	templates["contact_user"] = contactUserTemplate

	// Load aanmelding email templates
	aanmeldingAdminTemplate, err := template.ParseFiles(fmt.Sprintf("%s/templates/aanmelding_admin_email.html", cwd))
	if err != nil {
		return nil, fmt.Errorf("failed to parse aanmelding admin template: %v", err)
	}
	templates["aanmelding_admin"] = aanmeldingAdminTemplate

	aanmeldingUserTemplate, err := template.ParseFiles(fmt.Sprintf("%s/templates/aanmelding_email.html", cwd))
	if err != nil {
		return nil, fmt.Errorf("failed to parse aanmelding user template: %v", err)
	}
	templates["aanmelding_user"] = aanmeldingUserTemplate

	return &EmailService{
		templates: templates,
	}, nil
}

func (s *EmailService) SendContactEmail(data *models.ContactEmailData) error {
	var templateName string
	var subject string
	var recipient string

	if data.ToAdmin {
		templateName = "contact_admin"
		subject = "Nieuw contactformulier ontvangen"
		recipient = data.AdminEmail
		log.Printf("Sending admin email to: %s using template: %s", recipient, templateName)
	} else {
		templateName = "contact_user"
		subject = "Bedankt voor je bericht"
		recipient = data.Contact.Email
		log.Printf("Sending user email to: %s using template: %s", recipient, templateName)
	}

	template := s.templates[templateName]
	if template == nil {
		log.Printf("Template not found: %s", templateName)
		return fmt.Errorf("template not found: %s", templateName)
	}

	var body bytes.Buffer
	if err := template.Execute(&body, data); err != nil {
		log.Printf("Failed to execute template: %v", err)
		return fmt.Errorf("failed to execute template: %v", err)
	}

	log.Printf("Successfully generated email body for template: %s", templateName)
	return s.sendEmail(recipient, subject, body.String())
}

func (s *EmailService) SendAanmeldingEmail(data *models.AanmeldingEmailData) error {
	var templateName string
	var subject string
	var recipient string

	if data.ToAdmin {
		templateName = "aanmelding_admin"
		subject = "Nieuwe aanmelding ontvangen"
		recipient = data.AdminEmail
	} else {
		templateName = "aanmelding_user"
		subject = "Bedankt voor je aanmelding"
		recipient = data.Aanmelding.Email
	}

	template := s.templates[templateName]
	if template == nil {
		return fmt.Errorf("template not found: %s", templateName)
	}

	var body bytes.Buffer
	if err := template.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return s.sendEmail(recipient, subject, body.String())
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	log.Printf("Attempting to send email to: %s with subject: %s", to, subject)

	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@dekoninklijkeloop.nl")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	log.Printf("SMTP Configuration - Host: %s, Port: %s, Username: %s", smtpHost, smtpPortStr, smtpUsername)

	if smtpHost == "" || smtpPortStr == "" || smtpUsername == "" || smtpPassword == "" {
		return fmt.Errorf("missing SMTP configuration - Host: %s, Port: %s, Username: %s", smtpHost, smtpPortStr, smtpUsername)
	}

	// Parse SMTP port
	smtpPort := 587 // Default to 587 if parsing fails
	if port, err := strconv.Atoi(smtpPortStr); err == nil {
		smtpPort = port
	} else {
		log.Printf("Failed to parse SMTP port, using default 587: %v", err)
	}

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	// Enable SSL/TLS
	d.SSL = true

	log.Printf("Attempting to connect to SMTP server: %s:%d", smtpHost, smtpPort)
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("Successfully sent email to: %s", to)
	return nil
}
