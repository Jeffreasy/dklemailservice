package services

import (
	"bytes"
	"dklautomationgo/models"
	"fmt"
	"html/template"
	"os"

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

	// Load contact email templates
	contactAdminTemplate, err := template.ParseFiles("templates/contact_admin_email.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse contact admin template: %v", err)
	}
	templates["contact_admin"] = contactAdminTemplate

	contactUserTemplate, err := template.ParseFiles("templates/contact_email.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse contact user template: %v", err)
	}
	templates["contact_user"] = contactUserTemplate

	// Load aanmelding email templates
	aanmeldingAdminTemplate, err := template.ParseFiles("templates/aanmelding_admin_email.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse aanmelding admin template: %v", err)
	}
	templates["aanmelding_admin"] = aanmeldingAdminTemplate

	aanmeldingUserTemplate, err := template.ParseFiles("templates/aanmelding_email.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse aanmelding user template: %v", err)
	}
	templates["aanmelding_user"] = aanmeldingUserTemplate

	return &EmailService{
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
	var recipient string

	if data.ToAdmin {
		templateName = "contact_admin"
		subject = "Nieuw contactformulier ontvangen"
		recipient = data.AdminEmail
	} else {
		templateName = "contact_user"
		subject = "Bedankt voor je bericht"
		recipient = data.Contact.Email
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
	m := gomail.NewMessage()
	m.SetHeader("From", "noreply@dekoninklijkeloop.nl")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" || smtpPortStr == "" || smtpUsername == "" || smtpPassword == "" {
		return fmt.Errorf("missing SMTP configuration")
	}

	d := gomail.NewDialer(smtpHost, 587, smtpUsername, smtpPassword)
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
