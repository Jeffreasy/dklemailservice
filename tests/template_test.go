package tests

import (
	"dklautomationgo/models"
	"dklautomationgo/services"
	"testing"
	"time"
)

func TestEmailTemplateRendering(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     interface{}
		wantErr  bool
	}{
		{
			name:     "Contact admin template - alle velden",
			template: "contact_admin",
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
			template: "contact_admin",
			data: &models.ContactEmailData{
				Contact: &models.ContactFormulier{
					Naam:  "Minimaal",
					Email: "min@example.com",
				},
			},
		},
		{
			name:     "Aanmelding template - optionele velden",
			template: "aanmelding_user",
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

	smtp := newMockSMTP()
	service, err := services.NewTestEmailService(smtp)
	if err != nil {
		t.Fatalf("Failed to create email service: %v", err)
	}

	for _, tt := range tests {
		tt := tt // Maak een kopie van de loopvariabele
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run tests in parallel

			tmpl := service.GetTemplate(tt.template)
			if tmpl == nil {
				if !tt.wantErr {
					t.Errorf("Template %s not found", tt.template)
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
