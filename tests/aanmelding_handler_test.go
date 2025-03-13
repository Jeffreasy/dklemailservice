package tests

import (
	"bytes"
	"dklautomationgo/handlers"
	"dklautomationgo/models"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestAanmeldingHandler(t *testing.T) {
	tests := []struct {
		name          string
		input         models.AanmeldingFormulier
		expectedCode  int
		errorContains string
		setupMock     func(*mockEmailService)
	}{
		{
			name: "Succesvolle aanmelding",
			input: models.AanmeldingFormulier{
				Naam:           "Test Runner",
				Email:          "runner@example.com",
				Telefoon:       "0612345678",
				Rol:            "loper",
				Afstand:        "10km",
				Ondersteuning:  "geen",
				Bijzonderheden: "geen",
				Terms:          true,
			},
			expectedCode: fiber.StatusOK,
			setupMock: func(m *mockEmailService) {
				m.shouldFail = false
			},
		},
		{
			name: "Aanmelding zonder terms akkoord",
			input: models.AanmeldingFormulier{
				Naam:           "Test Runner",
				Email:          "runner@example.com",
				Telefoon:       "0612345678",
				Rol:            "loper",
				Afstand:        "10km",
				Ondersteuning:  "geen",
				Bijzonderheden: "geen",
				Terms:          false,
			},
			expectedCode:  fiber.StatusBadRequest,
			errorContains: "Je moet akkoord gaan met de voorwaarden",
			setupMock: func(m *mockEmailService) {
				m.shouldFail = false
			},
		},
		{
			name: "Ongeldige email",
			input: models.AanmeldingFormulier{
				Naam:           "Test Runner",
				Email:          "invalid-email",
				Telefoon:       "0612345678",
				Rol:            "loper",
				Afstand:        "10km",
				Ondersteuning:  "geen",
				Bijzonderheden: "geen",
				Terms:          true,
			},
			expectedCode:  fiber.StatusBadRequest,
			errorContains: "Ongeldig email adres",
			setupMock: func(m *mockEmailService) {
				m.shouldFail = false
			},
		},
		{
			name: "Ontbrekende verplichte velden",
			input: models.AanmeldingFormulier{
				Naam:           "",
				Email:          "runner@example.com",
				Telefoon:       "0612345678",
				Rol:            "loper",
				Afstand:        "10km",
				Ondersteuning:  "geen",
				Bijzonderheden: "geen",
				Terms:          true,
			},
			expectedCode:  fiber.StatusBadRequest,
			errorContains: "Naam is verplicht",
			setupMock: func(m *mockEmailService) {
				m.shouldFail = false
			},
		},
		{
			name: "Email verzending mislukt",
			input: models.AanmeldingFormulier{
				Naam:           "Test Runner",
				Email:          "runner@example.com",
				Telefoon:       "0612345678",
				Rol:            "loper",
				Afstand:        "10km",
				Ondersteuning:  "geen",
				Bijzonderheden: "geen",
				Terms:          true,
			},
			expectedCode:  fiber.StatusInternalServerError,
			errorContains: "Fout bij het verzenden van de email",
			setupMock: func(m *mockEmailService) {
				m.shouldFail = true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			mockService := newMockEmailService()

			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			handler := handlers.NewEmailHandler(mockService)
			app.Post("/aanmelding-email", handler.HandleAanmeldingEmail)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest("POST", "/aanmelding-email", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err, "Failed to test request")
			assert.Equal(t, tt.expectedCode, resp.StatusCode, "Unexpected status code")

			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			assert.NoError(t, err, "Failed to decode response")

			if tt.errorContains != "" {
				assert.Contains(t, result["error"], tt.errorContains, "Error message does not contain expected text")
			} else {
				assert.Equal(t, true, result["success"], "Expected success to be true")
			}
		})
	}
}
