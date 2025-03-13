package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthResponse bevat informatie over de status van de service
type HealthResponse struct {
	Status    string    `json:"status"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
}

// Exporteer de constanten zodat ze vanuit main.go bereikbaar zijn
var StartTime = time.Now()
var Version = "1.0.0" // Deze zou uit buildinfo moeten komen

// HealthHandler biedt een health check endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "ok",
		Version:   Version,
		Timestamp: time.Now(),
		Uptime:    time.Since(StartTime).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
