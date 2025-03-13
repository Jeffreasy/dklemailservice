package tests

import (
	"dklautomationgo/logger"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestELKWriter(t *testing.T) {
	// Setup een test HTTP server die logs ontvangt
	var receivedLogs []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Lees de body
		buf, err := io.ReadAll(r.Body)
		if err == nil {
			receivedLogs = buf
		}

		// Stuur 200 OK terug
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Maak een ELK writer die naar onze test server stuurt
	config := logger.ELKConfig{
		Endpoint:      server.URL,
		BatchSize:     2, // Klein voor de test
		FlushInterval: 100 * time.Millisecond,
		AppName:       "test-app",
		Environment:   "test",
	}

	writer := logger.NewELKWriter(config)

	// Log een paar berichten
	entry1 := map[string]interface{}{
		"level":   "info",
		"message": "Test bericht 1",
		"extra":   "waarde1",
	}

	entry2 := map[string]interface{}{
		"level":   "error",
		"message": "Test bericht 2",
		"extra":   "waarde2",
	}

	// Schrijf logs
	if err := writer.Write(entry1); err != nil {
		t.Fatalf("Fout bij schrijven log 1: %v", err)
	}

	// Wacht kort om de test betrouwbaarder te maken
	time.Sleep(50 * time.Millisecond)

	if err := writer.Write(entry2); err != nil {
		t.Fatalf("Fout bij schrijven log 2: %v", err)
	}

	// De buffer zou nu vol moeten zijn en moet automatisch flushen
	// Wacht even om te zorgen dat de flush voltooid is
	time.Sleep(150 * time.Millisecond)

	// Controleer of logs zijn ontvangen (niet perfect, maar geeft basis werking aan)
	if len(receivedLogs) == 0 {
		t.Error("Geen logs ontvangen op de server")
	}

	// Sluit de writer
	if err := writer.Close(); err != nil {
		t.Fatalf("Fout bij sluiten writer: %v", err)
	}
}
