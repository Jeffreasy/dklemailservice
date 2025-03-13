package tests

import (
	"bytes"
	"dklautomationgo/handlers"
	"dklautomationgo/models"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestEmailHandlerFlow(t *testing.T) {
	tests := []struct {
		name       string
		input      interface{}
		wantErr    bool
		checkCalls func(*testing.T, *mockEmailService)
	}{
		{
			name: "Succesvol contactformulier",
			input: models.ContactFormulier{
				Naam:           "Test",
				Email:          "test@example.com",
				Bericht:        "Test bericht",
				PrivacyAkkoord: true,
			},
			wantErr: false,
			checkCalls: func(t *testing.T, m *mockEmailService) {
				if !m.contactEmailCalled {
					t.Error("Contact email method was not called")
				}
			},
		},
		{
			name: "Contact formulier zonder privacy akkoord",
			input: models.ContactFormulier{
				Naam:           "Test",
				Email:          "test@example.com",
				Bericht:        "Test bericht",
				PrivacyAkkoord: false,
			},
			wantErr: true,
			checkCalls: func(t *testing.T, m *mockEmailService) {
				if m.contactEmailCalled {
					t.Error("Contact email method should not be called")
				}
			},
		},
		{
			name: "Email verzending mislukt",
			input: models.ContactFormulier{
				Naam:           "Test",
				Email:          "test@example.com",
				Bericht:        "Test bericht",
				PrivacyAkkoord: true,
			},
			wantErr: true,
			checkCalls: func(t *testing.T, m *mockEmailService) {
				// Geen check nodig, we weten dat de call mislukt
				// omdat we shouldFail op true hebben gezet
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			// Gebruik de bestaande mock uit mocks.go
			mockService := newMockEmailService()

			// Stel shouldFail in als we een error willen simuleren
			if tt.name == "Email verzending mislukt" {
				mockService.shouldFail = true
			}

			handler := handlers.NewEmailHandler(mockService)

			app.Post("/contact-email", handler.HandleContactEmail)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/contact-email", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}

			if tt.wantErr {
				if resp.StatusCode != fiber.StatusBadRequest &&
					resp.StatusCode != fiber.StatusInternalServerError {
					t.Errorf("Expected error status code, got %d", resp.StatusCode)
				}
			} else {
				if resp.StatusCode != fiber.StatusOK {
					t.Errorf("Expected status OK, got %d", resp.StatusCode)
				}
			}

			if tt.checkCalls != nil {
				tt.checkCalls(t, mockService)
			}
		})
	}
}
