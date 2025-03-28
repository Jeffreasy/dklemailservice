package services

import (
	"bytes"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var RetryDelayFactor = 100 // milliseconden

// SMTPClient interface definieert de methode voor het verzenden van emails
type SMTPClient interface {
	Send(msg *EmailMessage) error
	SendRegistration(msg *EmailMessage) error
	SendEmail(to, subject, body string) error
	SendWFC(msg *EmailMessage) error             // Nieuwe methode voor Whisky for Charity
	SendWFCEmail(to, subject, body string) error // Helper methode voor Whisky for Charity
}

// EmailService is verantwoordelijk voor het versturen van emails
type EmailService struct {
	smtpClient        SMTPClient
	templates         map[string]*template.Template
	rateLimiter       RateLimiterInterface
	metrics           *EmailMetrics
	prometheusMetrics PrometheusMetricsInterface
	excludedEmails    []string
	mu                sync.RWMutex
}

// EmailMessage representeert een te verzenden email
type EmailMessage struct {
	To       string
	Subject  string
	Body     string
	TestMode bool
}

// NewEmailService maakt een nieuwe EmailService met de opgegeven SMTP client
// Laadt templates uit de configureerbare template directory
// en configureert rate limiting op basis van omgevingsvariabelen
func NewEmailService(smtpClient SMTPClient, metrics *EmailMetrics, rateLimiter RateLimiterInterface, prometheusMetrics PrometheusMetricsInterface) *EmailService {
	return NewEmailServiceWithTemplatesDir(smtpClient, metrics, rateLimiter, prometheusMetrics, "templates")
}

// NewEmailServiceWithTemplatesDir maakt een nieuwe EmailService met de opgegeven SMTP client en templates directory
func NewEmailServiceWithTemplatesDir(smtpClient SMTPClient, metrics *EmailMetrics, rateLimiter RateLimiterInterface, prometheusMetrics PrometheusMetricsInterface, templatesDir string) *EmailService {
	// Definieer template functies
	templateFuncs := template.FuncMap{
		"multiply": func(a, b interface{}) float64 {
			// Convert interface values to float64
			var floatA, floatB float64
			switch v := a.(type) {
			case int:
				floatA = float64(v)
			case float64:
				floatA = v
			}
			switch v := b.(type) {
			case int:
				floatB = float64(v)
			case float64:
				floatB = v
			}
			return floatA * floatB
		},
		"currentYear": func() int {
			return time.Now().Year()
		},
	}

	// Laad alle templates bij initialisatie
	templates := make(map[string]*template.Template)
	templateFiles := []string{
		"contact_admin_email",
		"contact_email",
		"aanmelding_admin_email",
		"aanmelding_email",
		"wfc_order_confirmation",
		"wfc_order_admin",
	}

	for _, name := range templateFiles {
		templatePath := filepath.Join(templatesDir, name+".html")

		// Maak een nieuwe template met functies
		tmpl := template.New(name + ".html").Funcs(templateFuncs)

		// Parse het template bestand
		tmpl, err := tmpl.ParseFiles(templatePath)
		if err != nil {
			logger.Error("Failed to load template", "template", name, "error", err)
			continue
		}
		templates[name] = tmpl
		logger.Info("Template loaded successfully", "template", name)
	}

	// Laad uitgesloten email adressen
	excludedEmails := []string{}
	if excludedEmailsStr := os.Getenv("EXCLUDE_TEST_EMAILS"); excludedEmailsStr != "" {
		excludedEmails = strings.Split(excludedEmailsStr, ",")
		for i, email := range excludedEmails {
			excludedEmails[i] = strings.TrimSpace(strings.ToLower(email))
		}
		logger.Info("Uitgesloten test email adressen geladen", "count", len(excludedEmails))
	}

	return &EmailService{
		templates:         templates,
		smtpClient:        smtpClient,
		metrics:           metrics,
		rateLimiter:       rateLimiter,
		prometheusMetrics: prometheusMetrics,
		excludedEmails:    excludedEmails,
	}
}

// isExcludedEmail controleert of een email adres is uitgesloten van test emails
func (s *EmailService) isExcludedEmail(email string) bool {
	email = strings.TrimSpace(strings.ToLower(email))
	for _, excluded := range s.excludedEmails {
		if email == excluded {
			return true
		}
	}
	return false
}

func (s *EmailService) SendContactEmail(data *models.ContactEmailData) error {
	// Check if this is a test email to an excluded address
	if data.Contact.TestMode && !data.ToAdmin {
		if s.isExcludedEmail(data.Contact.Email) {
			logger.Info("Test email overgeslagen voor uitgesloten adres",
				"email", data.Contact.Email,
				"type", "contact")
			return nil
		}
	}

	// Log appropriately
	if data.ToAdmin {
		logger.Debug("Contact admin email wordt voorbereid",
			"naam", data.Contact.Naam,
			"email", data.Contact.Email,
			"test_mode", data.Contact.TestMode)
		return s.sendEmailWithTemplate("contact_admin_email", data.AdminEmail, "Nieuw contactformulier", data)
	}

	logger.Debug("Contact bevestigingsemail wordt voorbereid",
		"naam", data.Contact.Naam,
		"email", data.Contact.Email,
		"test_mode", data.Contact.TestMode)
	return s.sendEmailWithTemplate("contact_email", data.Contact.Email, "Bedankt voor je bericht", data)
}

func (s *EmailService) SendAanmeldingEmail(data *models.AanmeldingEmailData) error {
	// Check if this is a test email to an excluded address
	if data.Aanmelding.TestMode && !data.ToAdmin {
		if s.isExcludedEmail(data.Aanmelding.Email) {
			logger.Info("Test email overgeslagen voor uitgesloten adres",
				"email", data.Aanmelding.Email,
				"type", "aanmelding")
			return nil
		}
	}

	var templateName string
	var subject string
	var recipient string
	start := time.Now()

	if data.ToAdmin {
		templateName = "aanmelding_admin_email"
		subject = "Nieuwe aanmelding ontvangen"
		recipient = data.AdminEmail
	} else {
		templateName = "aanmelding_email"
		subject = "Bedankt voor je aanmelding"
		recipient = data.Aanmelding.Email
	}

	template := s.GetTemplate(templateName)
	if template == nil {
		return fmt.Errorf("template not found: %s", templateName)
	}

	// Genereer email body
	var body bytes.Buffer
	if err := template.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	// Controleer rate limits voordat we een poging wagen
	if !s.rateLimiter.AllowEmail("email_generic", "") {
		if s.metrics != nil {
			s.metrics.RecordEmailFailed("aanmelding_email")
		}
		s.prometheusMetrics.RecordEmailFailed("aanmelding_email", "rate_limited")
		return fmt.Errorf("rate limit exceeded")
	}

	// Bereid bericht voor
	msg := &EmailMessage{
		To:       recipient,
		Subject:  subject,
		Body:     body.String(),
		TestMode: data.Aanmelding.TestMode,
	}

	// Verzend met de juiste client op basis van type
	var err error
	if data.ToAdmin {
		err = s.smtpClient.Send(msg) // Gebruik standaard SMTP voor admin emails
	} else {
		err = s.smtpClient.SendRegistration(msg) // Gebruik registratie SMTP voor gebruiker emails
	}

	elapsedTime := time.Since(start)

	// Metrics bijwerken op basis van resultaat
	if err != nil {
		if s.metrics != nil {
			s.metrics.RecordEmailFailed("aanmelding_email")
		}
		s.prometheusMetrics.RecordEmailFailed("aanmelding_email", "smtp_error")
		s.prometheusMetrics.ObserveEmailLatency("aanmelding_email", elapsedTime.Seconds())
		return err
	}

	// Succesvolle verzending
	if s.metrics != nil {
		s.metrics.RecordEmailSent("aanmelding_email")
	}
	s.prometheusMetrics.RecordEmailSent("aanmelding_email", "success")
	s.prometheusMetrics.ObserveEmailLatency("aanmelding_email", elapsedTime.Seconds())
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

// GetTemplate geeft een template terug op basis van de naam
func (s *EmailService) GetTemplate(name string) *template.Template {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tmpl, exists := s.templates[name]
	if !exists {
		return nil
	}
	return tmpl
}

// ValidateTemplate valideert of een template correct kan worden uitgevoerd met de gegeven data
func ValidateTemplate(tmpl *template.Template, data interface{}) error {
	if tmpl == nil {
		return fmt.Errorf("template is nil")
	}

	// Render de template naar een buffer om te controleren of het werkt
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("template execution error: %w", err)
	}

	return nil
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

// SendWhiskyForCharityEmail verzendt een email via de WFC configuratie
func (s *EmailService) SendWhiskyForCharityEmail(to, subject, body string) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		if s.prometheusMetrics != nil {
			s.prometheusMetrics.ObserveEmailLatency("wfc_email", duration.Seconds())
		}
	}()

	if !s.rateLimiter.AllowEmail("email_generic", "") {
		if s.metrics != nil {
			s.metrics.RecordEmailFailed("wfc_email")
		}
		if s.prometheusMetrics != nil {
			s.prometheusMetrics.RecordEmailFailed("wfc_email", "rate_limited")
		}
		return fmt.Errorf("rate limit exceeded")
	}

	msg := &EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	err := s.smtpClient.SendWFC(msg)
	if err != nil {
		if s.metrics != nil {
			s.metrics.RecordEmailFailed("wfc_email")
		}
		if s.prometheusMetrics != nil {
			s.prometheusMetrics.RecordEmailFailed("wfc_email", "smtp_error")
		}
		return err
	}

	if s.metrics != nil {
		s.metrics.RecordEmailSent("wfc_email")
	}
	if s.prometheusMetrics != nil {
		s.prometheusMetrics.RecordEmailSent("wfc_email", "success")
	}
	return nil
}

// SendWFCOrderEmail sends order emails for Whisky for Charity
func (s *EmailService) SendWFCOrderEmail(data *models.WFCOrderEmailData) error {
	start := time.Now()

	var templateName string
	var subject string
	var recipient string

	if data.ToAdmin {
		templateName = "wfc_order_admin"
		subject = "New Order Received - Whisky for Charity"
		recipient = data.AdminEmail
	} else {
		templateName = "wfc_order_confirmation"
		subject = "Order Confirmation - Whisky for Charity"
		recipient = data.Order.CustomerEmail
	}

	// Get the template
	template := s.GetTemplate(templateName)
	if template == nil {
		logger.Error("Template not found", "template", templateName)
		return fmt.Errorf("template not found: %s", templateName)
	}

	// Generate email body (functies zijn nu al geregistreerd tijdens initialisatie)
	var body bytes.Buffer
	if err := template.Execute(&body, data); err != nil {
		logger.Error("Failed to execute template", "template", templateName, "error", err)
		return fmt.Errorf("failed to execute template: %v", err)
	}

	// Check rate limits
	if !s.rateLimiter.AllowEmail("wfc_email", "") {
		if s.metrics != nil {
			s.metrics.RecordEmailFailed("wfc_email")
		}
		if s.prometheusMetrics != nil {
			s.prometheusMetrics.RecordEmailFailed("wfc_email", "rate_limited")
		}
		return fmt.Errorf("rate limit exceeded")
	}

	// Prepare message
	msg := &EmailMessage{
		To:       recipient,
		Subject:  subject,
		Body:     body.String(),
		TestMode: false,
	}

	// Send via WFC SMTP
	err := s.smtpClient.SendWFC(msg)

	elapsedTime := time.Since(start)

	// Update metrics
	if err != nil {
		if s.metrics != nil {
			s.metrics.RecordEmailFailed("wfc_email")
		}
		if s.prometheusMetrics != nil {
			s.prometheusMetrics.RecordEmailFailed("wfc_email", "smtp_error")
			s.prometheusMetrics.ObserveEmailLatency("wfc_email", elapsedTime.Seconds())
		}
		logger.Error("Failed to send WFC order email", "recipient", recipient, "error", err)
		return err
	}

	// Success
	if s.metrics != nil {
		s.metrics.RecordEmailSent("wfc_email")
	}
	if s.prometheusMetrics != nil {
		s.prometheusMetrics.RecordEmailSent("wfc_email", "success")
		s.prometheusMetrics.ObserveEmailLatency("wfc_email", elapsedTime.Seconds())
	}

	logger.Info("WFC order email sent successfully",
		"template", templateName,
		"recipient", recipient,
		"order_id", data.Order.ID)

	return nil
}
