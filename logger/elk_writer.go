package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ELKConfig configuratie voor ELK/Graylog integratie
type ELKConfig struct {
	Endpoint      string        // URL endpoint voor Logstash of Graylog HTTP input
	BatchSize     int           // Aantal logs om te bufferen voor ze worden verzonden
	FlushInterval time.Duration // Hoe vaak logs worden verzonden
	AppName       string        // Naam van de applicatie (voor identificatie)
	Environment   string        // Omgeving (dev, test, prod, etc.)
}

// LogWriter interface voor verschillende log outputs
type LogWriter interface {
	Write(entry map[string]interface{}) error
	Flush() error
	Close() error
}

// ELKWriter implementeert verzenden van logs naar ELK of Graylog
type ELKWriter struct {
	config     ELKConfig
	httpClient *http.Client
	buffer     []map[string]interface{}
	ticker     *time.Ticker
	done       chan bool
}

// NewELKWriter maakt een nieuwe ELK/Graylog writer
func NewELKWriter(config ELKConfig) *ELKWriter {
	writer := &ELKWriter{
		config:     config,
		httpClient: &http.Client{Timeout: 5 * time.Second},
		buffer:     make([]map[string]interface{}, 0, config.BatchSize),
		done:       make(chan bool),
	}

	// Start een achtergrondproces dat logs periodiek verstuurt
	writer.ticker = time.NewTicker(config.FlushInterval)
	go writer.flushRoutine()

	return writer
}

// Write voegt een log entry toe aan de buffer
func (w *ELKWriter) Write(entry map[string]interface{}) error {
	// Voeg extra velden toe voor identificatie
	entry["application"] = w.config.AppName
	entry["environment"] = w.config.Environment
	entry["@timestamp"] = time.Now().UTC().Format(time.RFC3339)

	// Voeg toe aan buffer
	w.buffer = append(w.buffer, entry)

	// Als de buffer vol is, verstuur dan direct
	if len(w.buffer) >= w.config.BatchSize {
		return w.Flush()
	}

	return nil
}

// Flush verstuurt alle gebufferde logs
func (w *ELKWriter) Flush() error {
	if len(w.buffer) == 0 {
		return nil
	}

	// Converteer buffer naar JSON
	jsonData, err := json.Marshal(w.buffer)
	if err != nil {
		return err
	}

	fmt.Printf("ELK Writer verzendt %d logs naar %s\n", len(w.buffer), w.config.Endpoint)

	// Stuur naar ELK/Graylog
	resp, err := w.httpClient.Post(
		w.config.Endpoint,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Printf("ELK Writer fout: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	// Reset buffer na succesvol versturen
	w.buffer = w.buffer[:0]
	return nil
}

// flushRoutine verstuurt logs periodiek
func (w *ELKWriter) flushRoutine() {
	for {
		select {
		case <-w.ticker.C:
			_ = w.Flush() // Errors loggen we niet om oneindige loops te voorkomen
		case <-w.done:
			return
		}
	}
}

// Close stopt de writer en verstuurt resterende logs
func (w *ELKWriter) Close() error {
	w.ticker.Stop()
	w.done <- true
	return w.Flush()
}
