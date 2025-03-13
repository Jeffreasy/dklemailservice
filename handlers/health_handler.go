package handlers

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"dklautomationgo/services"

	"gopkg.in/gomail.v2"
)

// ServiceStatus representeert de verschillende statussen die de service kan hebben
type ServiceStatus string

const (
	StatusHealthy   ServiceStatus = "healthy"
	StatusDegraded  ServiceStatus = "degraded"
	StatusUnhealthy ServiceStatus = "unhealthy"
)

// HealthResponse bevat uitgebreide informatie over de status van de service
type HealthResponse struct {
	Status      ServiceStatus   `json:"status"`
	Version     string          `json:"version"`
	Timestamp   time.Time       `json:"timestamp"`
	Uptime      string          `json:"uptime"`
	Environment string          `json:"environment"`
	Memory      MemoryStats     `json:"memory"`
	System      SystemStats     `json:"system"`
	Checks      ComponentChecks `json:"checks"`
}

// MemoryStats bevat geheugen gerelateerde metrics
type MemoryStats struct {
	Alloc      uint64 `json:"alloc"`       // Bytes in gebruik
	TotalAlloc uint64 `json:"total_alloc"` // Totaal gealloceerde bytes
	HeapAlloc  uint64 `json:"heap_alloc"`  // Bytes in heap gebruik
	NumGC      uint32 `json:"num_gc"`      // Aantal garbage collections
}

// SystemStats bevat systeem gerelateerde metrics
type SystemStats struct {
	NumGoroutines int    `json:"num_goroutines"`
	NumCPU        int    `json:"num_cpu"`
	GoVersion     string `json:"go_version"`
}

// ComponentChecks bevat de status van verschillende service componenten
type ComponentChecks struct {
	SMTP struct {
		Default      bool   `json:"default"`      // Standaard SMTP verbinding
		Registration bool   `json:"registration"` // Registratie SMTP verbinding
		LastError    string `json:"last_error,omitempty"`
	} `json:"smtp"`
	RateLimiter struct {
		Status bool             `json:"status"`
		Limits map[string]Limit `json:"limits"`
	} `json:"rate_limiter"`
	Templates struct {
		Status    bool              `json:"status"`
		Available []string          `json:"available"`
		Errors    map[string]string `json:"errors,omitempty"`
	} `json:"templates"`
}

// Limit representeert een rate limit configuratie
type Limit struct {
	Count  int           `json:"count"`
	Window time.Duration `json:"window"`
	PerIP  bool          `json:"per_ip"`
}

// Exporteer de constanten zodat ze vanuit main.go bereikbaar zijn
var (
	StartTime   = time.Now()
	Version     = "1.1.0" // Deze zou uit buildinfo moeten komen
	Environment = "development"
	rateLimiter services.RateLimiterInterface // Wordt gezet via SetRateLimiter
)

// SetRateLimiter stelt de rate limiter in voor health checks
func SetRateLimiter(rl services.RateLimiterInterface) {
	rateLimiter = rl
}

// HealthHandler biedt een uitgebreide health check endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Verzamel memory stats
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Voer component checks uit
	checks := checkComponents()

	// Bepaal overall status
	status := determineOverallStatus(checks)

	response := HealthResponse{
		Status:      status,
		Version:     Version,
		Timestamp:   time.Now(),
		Uptime:      time.Since(StartTime).String(),
		Environment: Environment,
		Memory: MemoryStats{
			Alloc:      memStats.Alloc,
			TotalAlloc: memStats.TotalAlloc,
			HeapAlloc:  memStats.HeapAlloc,
			NumGC:      memStats.NumGC,
		},
		System: SystemStats{
			NumGoroutines: runtime.NumGoroutine(),
			NumCPU:        runtime.NumCPU(),
			GoVersion:     runtime.Version(),
		},
		Checks: checks,
	}

	w.Header().Set("Content-Type", "application/json")

	// Set HTTP status based on service status
	switch status {
	case StatusHealthy:
		w.WriteHeader(http.StatusOK)
	case StatusDegraded:
		w.WriteHeader(http.StatusOK) // Still 200 but with degraded status
	case StatusUnhealthy:
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// checkComponents voert health checks uit op verschillende service componenten
func checkComponents() ComponentChecks {
	var checks ComponentChecks

	// Check SMTP connections
	defaultSMTP, regSMTP, lastError := checkSMTPConnections()
	checks.SMTP.Default = defaultSMTP
	checks.SMTP.Registration = regSMTP
	checks.SMTP.LastError = lastError

	// Check rate limiter
	checks.RateLimiter.Status = rateLimiter != nil
	checks.RateLimiter.Limits = getRateLimits()

	// Check templates
	templateStatus, available, errors := checkTemplates()
	checks.Templates.Status = templateStatus
	checks.Templates.Available = available
	checks.Templates.Errors = errors

	return checks
}

// getRateLimits haalt de actuele rate limits op
func getRateLimits() map[string]Limit {
	if rateLimiter == nil {
		// Fallback naar default configuratie
		return map[string]Limit{
			"contact_email_global":    {Count: 100, Window: time.Hour, PerIP: false},
			"contact_email_ip":        {Count: 5, Window: time.Hour, PerIP: true},
			"aanmelding_email_global": {Count: 200, Window: time.Hour, PerIP: false},
			"aanmelding_email_ip":     {Count: 10, Window: time.Hour, PerIP: true},
		}
	}

	// Haal limieten op van de rate limiter
	serviceLimits := rateLimiter.GetLimits()
	limits := make(map[string]Limit)

	// Converteer de service limieten naar het health check formaat
	for name, limit := range serviceLimits {
		// Voor elke limiet maken we een globale en IP-specifieke variant
		if limit.PerIP {
			limits[name+"_ip"] = Limit{
				Count:  limit.Count,
				Window: limit.Period,
				PerIP:  true,
			}
		} else {
			limits[name+"_global"] = Limit{
				Count:  limit.Count,
				Window: limit.Period,
				PerIP:  false,
			}
		}
	}

	return limits
}

// checkTemplates controleert of alle benodigde templates aanwezig en geldig zijn
func checkTemplates() (bool, []string, map[string]string) {
	templateFiles := []string{
		"contact_admin_email",
		"contact_email",
		"aanmelding_admin_email",
		"aanmelding_email",
	}

	available := make([]string, 0)
	errors := make(map[string]string)

	for _, name := range templateFiles {
		templatePath := filepath.Join("templates", name+".html")

		// Check of template bestaat
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			errors[name] = "Template file not found"
			continue
		}

		// Probeer template te parsen
		_, err := template.ParseFiles(templatePath)
		if err != nil {
			errors[name] = "Template parse error: " + err.Error()
			continue
		}

		available = append(available, name)
	}

	return len(errors) == 0, available, errors
}

// determineOverallStatus bepaalt de algemene status op basis van component checks
func determineOverallStatus(checks ComponentChecks) ServiceStatus {
	// Als beide SMTP verbindingen down zijn, is de service unhealthy
	if !checks.SMTP.Default && !checks.SMTP.Registration {
		return StatusUnhealthy
	}

	// Als één van de SMTP verbindingen down is, is de service degraded
	if !checks.SMTP.Default || !checks.SMTP.Registration {
		return StatusDegraded
	}

	// Als templates of rate limiter issues hebben, is de service degraded
	if !checks.Templates.Status || !checks.RateLimiter.Status {
		return StatusDegraded
	}

	return StatusHealthy
}

// checkSMTPConnections test beide SMTP verbindingen
func checkSMTPConnections() (defaultOK bool, registrationOK bool, lastError string) {
	// Check default SMTP
	defaultOK = checkSMTPConnection(
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
	)

	// Check registration SMTP
	registrationOK = checkSMTPConnection(
		os.Getenv("REGISTRATION_SMTP_HOST"),
		os.Getenv("REGISTRATION_SMTP_PORT"),
		os.Getenv("REGISTRATION_SMTP_USER"),
		os.Getenv("REGISTRATION_SMTP_PASSWORD"),
	)

	if !defaultOK && !registrationOK {
		lastError = "Both SMTP connections failed"
	} else if !defaultOK {
		lastError = "Default SMTP connection failed"
	} else if !registrationOK {
		lastError = "Registration SMTP connection failed"
	}

	return
}

// checkSMTPConnection test een enkele SMTP verbinding
func checkSMTPConnection(host, portStr, user, password string) bool {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 587 // Default SMTP port
	}

	d := gomail.NewDialer(host, port, user, password)

	// Gebruik een context met timeout voor de health check
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Maak een channel voor het resultaat
	done := make(chan error, 1)

	go func() {
		_, err := d.Dial()
		done <- err
	}()

	// Wacht op het resultaat of timeout
	select {
	case err := <-done:
		return err == nil
	case <-ctx.Done():
		return false
	}
}
