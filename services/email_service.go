package services

import (
	"bytes"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"fmt"
	"html/template"
	"path/filepath"
	"time"
)

var RetryDelayFactor = 100 // milliseconden

// SMTPClient interface definieert de methode voor het verzenden van emails
type SMTPClient interface {
	Send(msg *EmailMessage) error
	SendRegistration(msg *EmailMessage) error
	SendEmail(to, subject, body string) error
}

// EmailService beheert email templates en afhandeling van email verzending
type EmailService struct {
	templates         map[string]*template.Template
	rateLimiter       RateLimiterInterface
	smtpClient        SMTPClient
	metrics           *EmailMetrics
	prometheusMetrics PrometheusMetricsInterface
}

// EmailMessage representeert een te verzenden email
type EmailMessage struct {
	To      string
	Subject string
	Body    string
}

// NewEmailService maakt een nieuwe EmailService met de opgegeven SMTP client
// Laadt templates uit de configureerbare template directory
// en configureert rate limiting op basis van omgevingsvariabelen
func NewEmailService(smtpClient SMTPClient, metrics *EmailMetrics, rateLimiter RateLimiterInterface, prometheusMetrics PrometheusMetricsInterface) *EmailService {
	// Laad alle templates bij initialisatie
	templates := make(map[string]*template.Template)
	templateFiles := []string{
		"contact_admin_email",
		"contact_email",
		"aanmelding_admin_email",
		"aanmelding_email",
	}

	for _, name := range templateFiles {
		templatePath := filepath.Join("templates", name+".html")
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			logger.Error("Failed to load template", "template", name, "error", err)
			continue
		}
		templates[name] = tmpl
		logger.Info("Template loaded successfully", "template", name)
	}

	return &EmailService{
		templates:         templates,
		smtpClient:        smtpClient,
		metrics:           metrics,
		rateLimiter:       rateLimiter,
		prometheusMetrics: prometheusMetrics,
	}
}

func (s *EmailService) SendContactEmail(data *models.ContactEmailData) error {
	// Log appropriately
	if data.ToAdmin {
		logger.Debug("Contact admin email wordt voorbereid",
			"naam", data.Contact.Naam,
			"email", data.Contact.Email)
		return s.sendEmailWithTemplate("contact_admin", data.AdminEmail, "Nieuw contactformulier", data)
	}

	logger.Debug("Contact bevestigingsemail wordt voorbereid",
		"naam", data.Contact.Naam,
		"email", data.Contact.Email)
	return s.sendEmailWithTemplate("contact_user", data.Contact.Email, "Bedankt voor je bericht", data)
}

func (s *EmailService) SendAanmeldingEmail(data *models.AanmeldingEmailData) error {
	var templateName string
	var subject string
	var recipient string

	if data.ToAdmin {
		templateName = "aanmelding_admin_email"
		subject = "Nieuwe aanmelding ontvangen"
		recipient = data.AdminEmail
	} else {
		templateName = "aanmelding_email"
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

	msg := &EmailMessage{
		To:      recipient,
		Subject: subject,
		Body:    body.String(),
	}

	if !s.rateLimiter.AllowEmail("email_generic", "") {
		return fmt.Errorf("rate limit exceeded")
	}

	var err error
	if data.ToAdmin {
		err = s.smtpClient.Send(msg) // Gebruik standaard SMTP voor admin emails
	} else {
		err = s.smtpClient.SendRegistration(msg) // Gebruik registratie SMTP voor gebruiker emails
	}

	if err != nil {
		if s.metrics != nil {
			s.metrics.RecordEmailFailed("aanmelding_email")
		}
		s.prometheusMetrics.RecordEmailFailed("aanmelding_email", "smtp_error")
		return err
	}

	if s.metrics != nil {
		s.metrics.RecordEmailSent("aanmelding_email")
	}
	s.prometheusMetrics.RecordEmailSent("aanmelding_email", "success")
	return nil
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	if !s.rateLimiter.AllowEmail("email_generic", "") {
		return fmt.Errorf("rate limit exceeded")
	}

	msg := &EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	err := s.smtpClient.Send(msg)
	if err != nil {
		if s.metrics != nil {
			s.metrics.RecordEmailFailed("email_generic")
		}
		return err
	}

	if s.metrics != nil {
		s.metrics.RecordEmailSent("email_generic")
	}

	return nil
}

func (s *EmailService) GetTemplate(name string) *template.Template {
	return s.templates[name]
}

func ValidateTemplate(tmpl *template.Template, data interface{}) error {
	var buf bytes.Buffer
	return tmpl.Execute(&buf, data)
}

func (s *EmailService) SendEmail(to, subject, body string) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		s.prometheusMetrics.ObserveEmailLatency("email_generic", duration.Seconds())
	}()

	if !s.rateLimiter.AllowEmail("email_generic", "") {
		return fmt.Errorf("rate limit exceeded")
	}

	msg := &EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	err := s.smtpClient.Send(msg)
	if err != nil {
		if s.metrics != nil {
			s.metrics.RecordEmailFailed("email_generic")
		}
		s.prometheusMetrics.RecordEmailFailed("email", "smtp_error")
		return err
	}

	if s.metrics != nil {
		s.metrics.RecordEmailSent("email_generic")
	}
	s.prometheusMetrics.RecordEmailSent("email", "success")
	return nil
}

func NewTestEmailService(smtpClient SMTPClient) (*EmailService, error) {
	// Maak een eenvoudige template voor testen
	tmpl, err := template.New("test").Parse("<p>Test template</p>")
	if err != nil {
		return nil, err
	}

	templates := make(map[string]*template.Template)
	templates["contact_admin_email"] = tmpl
	templates["contact_email"] = tmpl
	templates["aanmelding_admin_email"] = tmpl
	templates["aanmelding_email"] = tmpl

	// Maak een metrics tracker voor testen
	metrics := NewEmailMetrics(24 * time.Hour)

	// Maak een test RateLimiter
	rateLimiter := NewRateLimiter(nil)
	rateLimiter.AddLimit("email_generic", 1000, time.Minute, false)

	return &EmailService{
		templates:         templates,
		rateLimiter:       rateLimiter,
		smtpClient:        smtpClient,
		metrics:           metrics,
		prometheusMetrics: nil,
	}, nil
}

func (s *EmailService) sendEmailWithTemplate(templateName, to, subject string, data interface{}) error {
	template := s.templates[templateName]
	if template == nil {
		logger.Error("Template not found", "template", templateName)
		return fmt.Errorf("template not found: %s", templateName)
	}

	var body bytes.Buffer
	if err := template.Execute(&body, data); err != nil {
		logger.Error("Failed to execute template", "error", err)
		return fmt.Errorf("failed to execute template: %v", err)
	}

	logger.Debug("Successfully generated email body for template", "template", templateName)
	return s.sendEmail(to, subject, body.String())
}

// SendTemplateEmail verzendt een email met template (voor batcher)
func (s *EmailService) SendTemplateEmail(recipient, subject, templateName string, templateData map[string]interface{}) error {
	template := s.templates[templateName]
	if template == nil {
		logger.Error("Template not found", "template", templateName)
		return fmt.Errorf("template not found: %s", templateName)
	}

	var body bytes.Buffer
	if err := template.Execute(&body, templateData); err != nil {
		logger.Error("Template rendering fout", "error", err, "template", templateName)
		return err
	}

	// Email verzenden
	err := s.smtpClient.SendEmail(recipient, subject, body.String())
	if err != nil {
		s.metrics.RecordEmailFailed(templateName)
		return err
	}

	s.metrics.RecordEmailSent(templateName)
	return nil
}

// SetMetrics stelt een nieuwe metrics tracker in (voor testen)
func (s *EmailService) SetMetrics(metrics *EmailMetrics) {
	s.metrics = metrics
}

// SetRateLimiter stelt een nieuwe rate limiter in (voor testen)
func (s *EmailService) SetRateLimiter(limiter RateLimiterInterface) {
	s.rateLimiter = limiter
}
