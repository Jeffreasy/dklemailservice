// Package tests bevat tests voor de applicatie
package tests

import (
	"dklautomationgo/models"
	"dklautomationgo/services"
	"testing"
	"time"
)

// TestEmailTemplateRendering test het renderen van email templates
func TestEmailTemplateRendering(t *testing.T) {
	// Zoek de templates directory met onze helper
	templatesDir, err := GetTemplatesDir()
	if err != nil {
		t.Skip("Templates directory niet gevonden, test wordt overgeslagen")
	}

	// Initialiseer een mock SMTP client
	smtp := &mockSMTP{}

	// Maak een aangepaste EmailService met de juiste templates
	metrics := services.NewEmailMetrics(24 * time.Hour)
	rateLimiter := services.NewRateLimiter(nil)
	rateLimiter.AddLimit("email_generic", 1000, time.Minute, false)

	service := services.NewEmailServiceWithTemplatesDir(smtp, metrics, rateLimiter, nil, templatesDir)

	tests := []struct {
		name     string
		template string
		data     interface{}
		wantErr  bool
	}{
		{
			name:     "Contact admin template - alle velden",
			template: "contact_admin_email",
			data: &models.ContactEmailData{
				Contact: &models.ContactFormulier{
					Naam:           "Test Persoon",
					Email:          "test@example.com",
					Bericht:        "Test bericht",
					PrivacyAkkoord: true,
					CreatedAt:      time.Now(),
				},
			},
		},
		{
			name:     "Contact admin template - minimale velden",
			template: "contact_admin_email",
			data: &models.ContactEmailData{
				Contact: &models.ContactFormulier{
					Naam:  "Minimaal",
					Email: "min@example.com",
				},
			},
		},
		{
			name:     "Aanmelding template - optionele velden",
			template: "aanmelding_email",
			data: &models.AanmeldingEmailData{
				Aanmelding: &models.AanmeldingFormulier{
					Naam:  "Test Deelnemer",
					Email: "test@example.com",
					Rol:   "Deelnemer",
				},
			},
		},
		{
			name:     "Ongeldige template",
			template: "niet_bestaande_template",
			data:     &models.ContactEmailData{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		tt := tt // Maak een kopie van de loopvariabele
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run tests in parallel

			tmpl := service.GetTemplate(tt.template)
			if tmpl == nil {
				if !tt.wantErr {
					t.Errorf("Template %s not found, maar was wel verwacht", tt.template)
				}
				return
			}

			err := services.ValidateTemplate(tmpl, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
