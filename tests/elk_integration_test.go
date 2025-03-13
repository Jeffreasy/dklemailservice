package tests

import (
	"dklautomationgo/logger"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// startMockElkServer maakt een mock ELK server voor testen
func startMockElkServer(t *testing.T) (*httptest.Server, chan []map[string]interface{}) {
	receivedLogs := make(chan []map[string]interface{}, 10)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Controleer HTTP methode
		if r.Method != "POST" {
			t.Errorf("Verwacht POST request, kreeg %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Lees request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error bij het lezen van request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Parse JSON logs
		var logs []map[string]interface{}
		if err := json.Unmarshal(body, &logs); err != nil {
			t.Errorf("Error bij het parsen van JSON logs: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		t.Logf("Ontvangen data: %s", string(body))

		// Stuur logs naar het channel
		receivedLogs <- logs

		// Stuur een succesvolle response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))

	return server, receivedLogs
}

func TestLoggerWithELK(t *testing.T) {
	// Skip deze test als SKIP_ELK_TEST omgevingsvariabele is ingesteld
	if os.Getenv("SKIP_ELK_TEST") == "true" {
		t.Skip("ELK test overgeslagen vanwege SKIP_ELK_TEST=true")
	}

	// Maak een mock ELK server
	mockServer, receivedLogs := startMockElkServer(t)
	defer mockServer.Close() // Zorg ervoor dat de server altijd wordt afgesloten

	// Configureer de logger om naar onze mock server te sturen
	elkURL := mockServer.URL
	os.Setenv("ELK_URL", elkURL)
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("APP_NAME", "test-app")
	os.Setenv("APP_ENV", "test")

	// Initialiseer de logger
	logger.Setup("debug")
	defer logger.Shutdown() // Belangrijke toevoeging: zorg dat de logger wordt afgesloten!

	// Test logging
	logger.Info("Test info bericht", "test_field", "test_value")
	logger.Error("Test error bericht", "error_code", 12345)

	// Geef wat tijd voor asynchrone verwerking
	time.Sleep(100 * time.Millisecond)

	// Sluit de ELK writer expliciet af
	logger.Shutdown()

	// Geef wat tijd om de laatste logs te verzenden
	time.Sleep(500 * time.Millisecond)

	// Controleer de logs
	select {
	case logs := <-receivedLogs:
		if len(logs) < 2 {
			t.Errorf("Verwacht minstens 2 log berichten, maar kreeg %d", len(logs))
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout bij het wachten op logs")
	}

	// Test is nu voltooid, server wordt automatisch afgesloten door defer
}
