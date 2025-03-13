package tests

import (
	"dklautomationgo/services"
	"testing"
	"time"
)

func TestEmailMetrics_RecordEmailSent(t *testing.T) {
	metrics := services.NewEmailMetrics(time.Hour) // Reset na een uur

	// Test het registreren van verzonden emails
	metrics.RecordEmailSent("contact_email")
	metrics.RecordEmailSent("contact_email")
	metrics.RecordEmailSent("aanmelding_email")

	// Verifieer totalen
	if total := metrics.GetTotalEmails(); total != 3 {
		t.Errorf("Verwacht totaal aantal emails: 3, maar kreeg: %d", total)
	}

	// Verifieer per type
	emailsByType := metrics.GetEmailsByType()
	if count, exists := emailsByType["contact_email"]; !exists || count != 2 {
		t.Errorf("Verwacht 2 contact emails, maar kreeg: %d", count)
	}
	if count, exists := emailsByType["aanmelding_email"]; !exists || count != 1 {
		t.Errorf("Verwacht 1 aanmelding email, maar kreeg: %d", count)
	}

	// Verifieer success rate (zou 100% moeten zijn)
	if rate := metrics.GetSuccessRate(); rate != 100.0 {
		t.Errorf("Verwacht 100%% success rate, maar kreeg: %.2f%%", rate)
	}
}

func TestEmailMetrics_RecordEmailFailed(t *testing.T) {
	metrics := services.NewEmailMetrics(time.Hour)

	// Test het registreren van zowel success als failures
	metrics.RecordEmailSent("contact_email")
	metrics.RecordEmailFailed("contact_email")
	metrics.RecordEmailFailed("aanmelding_email")

	// Verifieer totalen
	if total := metrics.GetTotalEmails(); total != 3 {
		t.Errorf("Verwacht totaal aantal emails: 3, maar kreeg: %d", total)
	}

	// Verifieer success rate (1 van de 3 succesvol = 33.33%)
	expectedRate := 33.33 // Ongeveer 33.33%
	rate := metrics.GetSuccessRate()
	if rate < expectedRate-1 || rate > expectedRate+1 {
		t.Errorf("Verwacht ongeveer %.2f%% success rate, maar kreeg: %.2f%%", expectedRate, rate)
	}
}

func TestEmailMetrics_Reset(t *testing.T) {
	metrics := services.NewEmailMetrics(time.Hour)

	// Voeg wat data toe
	metrics.RecordEmailSent("contact_email")
	metrics.RecordEmailSent("aanmelding_email")
	metrics.RecordEmailFailed("contact_email")

	// Verifieer dat we data hebben
	if total := metrics.GetTotalEmails(); total != 3 {
		t.Errorf("Verwacht totaal aantal emails: 3, maar kreeg: %d", total)
	}

	// Reset de metrics
	metrics.Reset()

	// Verifieer dat alles op nul staat
	if total := metrics.GetTotalEmails(); total != 0 {
		t.Errorf("Verwacht totaal aantal emails na reset: 0, maar kreeg: %d", total)
	}
	if len(metrics.GetEmailsByType()) != 0 {
		t.Errorf("Verwacht lege map na reset, maar heeft nog entries: %v", metrics.GetEmailsByType())
	}
}

func TestEmailMetrics_AutoReset(t *testing.T) {
	// Gebruik een ZEER korte reset interval voor de test
	resetInterval := 10 * time.Millisecond // Extreem kort voor test
	metrics := services.NewEmailMetrics(resetInterval)

	// Voeg wat data toe
	metrics.RecordEmailSent("contact_email")
	metrics.RecordEmailSent("aanmelding_email")

	// Verifieer dat we data hebben
	if total := metrics.GetTotalEmails(); total != 2 {
		t.Errorf("Verwacht totaal aantal emails: 2, maar kreeg: %d", total)
	}

	// Wacht veel langer dan de reset interval om zeker te zijn dat reset gebeurt
	time.Sleep(200 * time.Millisecond) // 20x de reset interval

	// Forceer de reset via de publieke CheckAndResetIfNeeded methode
	metrics.CheckAndResetIfNeeded()

	// Voeg nieuwe data toe
	metrics.RecordEmailSent("new_email")

	// Na de reset en nieuwe registratie moet er 1 email zijn
	if total := metrics.GetTotalEmails(); total != 1 {
		t.Errorf("Verwacht 1 email na auto-reset, maar kreeg: %d", total)
		// Log voor debug informatie
		t.Logf("Reset interval: %v, Sleep time: %v", resetInterval, 200*time.Millisecond)
	}
}
