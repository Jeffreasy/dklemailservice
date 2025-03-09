package services

import (
	"bytes"
	"dklautomationgo/models"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

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
	dialer    *gomail.Dialer
	templates map[string]*template.Template
}

func NewEmailService() (*EmailService, error) {
	// Initialize email dialer
	dialer := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		587,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
	)

	// Load email templates
	templates := make(map[string]*template.Template)
	templateDir := "templates"

	// Only load contact form templates for now
	templateFiles := []string{
		"contact_email.html",
		"contact_admin_email.html",
	}

	for _, file := range templateFiles {
		tmpl, err := template.ParseFiles(filepath.Join(templateDir, file))
		if err != nil {
			return nil, fmt.Errorf("error loading template %s: %v", file, err)
		}
		templates[file] = tmpl
	}

	return &EmailService{
		dialer:    dialer,
		templates: templates,
	}, nil
}

type templateData struct {
	ToAdmin     bool
	Contact     *models.ContactFormulier
	Aanmelding  *models.Aanmelding
	CurrentYear int
}

func formatEmptyField(value string) string {
	if value == "" {
		return `<span class="empty-field">Niet opgegeven</span>`
	}
	return value
}

func (s *EmailService) SendContactEmail(data *models.ContactEmailData) error {
	var templateName string
	var subject string
	var toEmail string

	if data.ToAdmin {
		templateName = "contact_admin_email.html"
		subject = fmt.Sprintf("Nieuw contactformulier van %s", data.Contact.Naam)
		toEmail = data.AdminEmail
	} else {
		templateName = "contact_email.html"
		subject = "Bedankt voor je bericht - De Koninklijke Loop"
		toEmail = data.Contact.Email
	}

	// Execute template
	var body bytes.Buffer
	if err := s.templates[templateName].Execute(&body, data); err != nil {
		log.Printf("Template execution error: %v", err)
		return fmt.Errorf("error executing template: %v", err)
	}

	// Create email message
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("De Koninklijke Loop <%s>", os.Getenv("SMTP_FROM")))
	m.SetHeader("To", toEmail)
	if !data.ToAdmin {
		m.SetHeader("Reply-To", os.Getenv("ADMIN_EMAIL"))
	}
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	// Send email
	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("Error sending email: %v", err)
		return fmt.Errorf("error sending email: %v", err)
	}

	return nil
}

// Tijdelijk uitgeschakeld tot we de aanmelding templates hebben
/*
func (s *EmailService) SendAanmeldingEmail(data *models.AanmeldingEmailData) error {
	return fmt.Errorf("aanmelding email service temporarily disabled")
}
*/
