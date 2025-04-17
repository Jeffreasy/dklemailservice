package tests

import (
	"dklautomationgo/handlers"
	"dklautomationgo/logger"
	"dklautomationgo/tests/mocks"
	"testing"
)

func TestLogger(t *testing.T) {
	// Setup de testlogger
	testLogger := logger.UseTestLogger()
	defer logger.RestoreDefaultLogger()

	// Log berichten op verschillende niveaus
	logger.Debug("Dit is een debug bericht", "key", "value")
	logger.Info("Dit is een info bericht", "user", "test@example.com")
	logger.Warn("Dit is een waarschuwing", "code", 123)
	logger.Error("Dit is een fout", "error", "something went wrong")

	// Haal alle logs op
	entries := testLogger.GetEntries()

	// Controleer of we het juiste aantal logs hebben
	if len(entries) != 4 {
		t.Errorf("Verwacht 4 logberichten, maar kreeg er %d", len(entries))
	}

	// Verifieer de log niveaus
	expectedLevels := []string{logger.DebugLevel, logger.InfoLevel, logger.WarnLevel, logger.ErrorLevel}
	for i, level := range expectedLevels {
		if i < len(entries) && entries[i].Level != level {
			t.Errorf("Verwacht logniveau %s, maar kreeg %s", level, entries[i].Level)
		}
	}

	// Verifieer de inhoud van een specifiek bericht
	foundError := false
	for _, entry := range entries {
		if entry.Level == logger.ErrorLevel && entry.Message == "Dit is een fout" {
			foundError = true
			errorVal, ok := entry.Fields["error"]
			if !ok || errorVal != "something went wrong" {
				t.Errorf("Error veld ontbreekt of heeft onjuiste waarde: %v", entry.Fields)
			}
		}
	}

	if !foundError {
		t.Error("Error bericht niet gevonden in logs")
	}
}

func TestHandlerLogging(t *testing.T) {
	// Setup de testlogger
	testLogger := logger.UseTestLogger()
	defer logger.RestoreDefaultLogger()

	// Maak een mock email service
	mockService := newMockEmailService()
	mockNotificationService := NewMockNotificationService()
	mockAanmeldingRepo := new(mocks.MockAanmeldingRepository)

	// Maak de email handler
	handler := handlers.NewEmailHandler(mockService, mockNotificationService, mockAanmeldingRepo)

	// Simuleer het afhandelen van een aanvraag (zonder daadwerkelijk HTTP te gebruiken)
	// We roepen hier alleen bepaalde functies aan die loggen
	handler.LogUserActivity("test@example.com", "contact formulier", "127.0.0.1")

	// Controleer of er een log is gemaakt
	entries := testLogger.GetEntries()

	// Verifieer dat we minimaal 1 log entry hebben
	if len(entries) < 1 {
		t.Fatal("Geen logberichten gegenereerd")
	}

	// Zoek naar een specifiek logbericht
	found := false
	for _, entry := range entries {
		if entry.Message == "Gebruikersactiviteit" {
			found = true
			// Controleer de velden
			if email, ok := entry.Fields["email"]; !ok || email != "test@example.com" {
				t.Errorf("Email veld niet gevonden of onjuist: %v", entry.Fields)
			}
			if ip, ok := entry.Fields["ip"]; !ok || ip != "127.0.0.1" {
				t.Errorf("IP veld niet gevonden of onjuist: %v", entry.Fields)
			}
		}
	}

	if !found {
		t.Error("Verwacht logbericht niet gevonden")
	}
}
