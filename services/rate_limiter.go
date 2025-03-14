package services

import (
	"dklautomationgo/logger"
	"fmt"
	"os"
	"sync"
	"time"
)

// RateLimiterInterface definieert de interface voor rate limiting
type RateLimiterInterface interface {
	AllowEmail(operation, key string) bool
	Allow(key string) bool
	GetLimits() map[string]RateLimit
	GetCurrentCount(operationType string, key string) int
	GetCurrentValues() map[string]int
}

// RateLimit definieert beperkingen voor email verzending
type RateLimit struct {
	Count  int           // Maximum aantal in de periode
	Period time.Duration // Periode waarbinnen de limiet geldt
	PerIP  bool          // Of limiet per IP geldt of globaal
}

// RateLimiter beheert het aantal verzendingen binnen tijdslimieten
type RateLimiter struct {
	mutex             sync.Mutex
	globalCounts      map[string]*counter
	ipCounts          map[string]map[string]*counter
	limits            map[string]RateLimit
	cleanupTicker     *time.Ticker
	done              chan bool
	prometheusMetrics *PrometheusMetrics
}

type counter struct {
	count     int
	resetTime time.Time
}

// NewRateLimiter creëert een nieuwe rate limiter
func NewRateLimiter(prometheusMetrics *PrometheusMetrics) *RateLimiter {
	rl := &RateLimiter{
		globalCounts:      make(map[string]*counter),
		ipCounts:          make(map[string]map[string]*counter),
		limits:            make(map[string]RateLimit),
		done:              make(chan bool),
		prometheusMetrics: prometheusMetrics,
	}

	// Start periodieke opschoning van verlopen limieten
	rl.cleanupTicker = time.NewTicker(10 * time.Minute)
	go rl.cleanupRoutine()

	return rl
}

// AddLimit voegt een nieuwe limiet toe voor een type operatie
func (rl *RateLimiter) AddLimit(operationType string, count int, period time.Duration, perIP bool) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	rl.limits[operationType] = RateLimit{
		Count:  count,
		Period: period,
		PerIP:  perIP,
	}
}

// AllowEmail controleert of een email mag worden verzonden
func (r *RateLimiter) AllowEmail(emailType, userID string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()

	// Controleer globale rate limiting
	if limit, exists := r.limits[emailType]; exists {
		// Bereken hoeveel requests in deze periode zijn toegestaan
		key := emailType
		if userID != "" && limit.PerIP {
			key = fmt.Sprintf("%s:%s", emailType, userID)
		}

		cnt, exists := r.globalCounts[key]
		if !exists {
			// Maak een nieuwe counter aan
			cnt = &counter{
				count:     0,
				resetTime: now,
			}
			r.globalCounts[key] = cnt
		}

		// Reset counter als de periode is verstreken
		elapsed := now.Sub(cnt.resetTime)
		if elapsed > limit.Period {
			cnt.count = 0
			cnt.resetTime = now
		}

		// Check tegen de limiet
		if cnt.count >= limit.Count {
			// Overschreden
			logger.Warn("Globale rate limit overschreden",
				"operation", emailType,
				"limit", limit.Count,
				"period", limit.Period)

			// Prometheus metrics registreren
			limitType := "global"
			if userID != "" {
				limitType = "per_user"
			}
			if r.prometheusMetrics != nil {
				r.prometheusMetrics.RecordRateLimitExceeded(emailType, limitType)
			}

			// Voor tests: voeg een korte vertraging toe als TEST_RATE_LIMITING=true
			if os.Getenv("TEST_RATE_LIMITING") == "true" {
				// In test mode gebruiken we een korte vertraging
				time.Sleep(100 * time.Millisecond)
			}

			return false
		}

		// Voor tests: voeg progressieve vertraging toe als TEST_RATE_LIMITING=true
		// en we al bij meer dan 1 request zijn
		if os.Getenv("TEST_RATE_LIMITING") == "true" && cnt.count > 0 {
			// In test mode gebruiken we een korte vertraging afhankelijk van het aantal requests
			time.Sleep(time.Duration(cnt.count) * 100 * time.Millisecond)
		}

		// Toegestaan, incrementeer counter
		cnt.count++
		return true
	}

	// Geen rate limiting voor dit type
	return true
}

// Allow controleert of een operatie is toegestaan op basis van de rate limit
func (r *RateLimiter) Allow(key string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Splits de key in operationType en userID (indien aanwezig)
	parts := make([]string, 2)
	if idx := indexOf(key, ":"); idx != -1 {
		parts[0] = key[:idx]
		if idx+1 < len(key) {
			parts[1] = key[idx+1:]
		}
	} else {
		parts[0] = key
	}

	operationType := parts[0]
	userID := parts[1]

	now := time.Now()

	// Controleer globale rate limiting
	if limit, exists := r.limits[operationType]; exists {
		// Bereken hoeveel requests in deze periode zijn toegestaan
		limitKey := operationType
		if userID != "" && limit.PerIP {
			limitKey = fmt.Sprintf("%s:%s", operationType, userID)
		}

		cnt, exists := r.globalCounts[limitKey]
		if !exists {
			// Maak een nieuwe counter aan
			cnt = &counter{
				count:     0,
				resetTime: now,
			}
			r.globalCounts[limitKey] = cnt
		}

		// Reset counter als de periode is verstreken
		elapsed := now.Sub(cnt.resetTime)
		if elapsed > limit.Period {
			cnt.count = 0
			cnt.resetTime = now
		}

		// Check tegen de limiet
		if cnt.count >= limit.Count {
			// Overschreden
			logger.Warn("Rate limit overschreden",
				"operation", operationType,
				"limit", limit.Count,
				"period", limit.Period)

			// Prometheus metrics registreren
			limitType := "global"
			if userID != "" {
				limitType = "per_user"
			}
			if r.prometheusMetrics != nil {
				r.prometheusMetrics.RecordRateLimitExceeded(operationType, limitType)
			}

			// Voor tests: voeg een korte vertraging toe als TEST_RATE_LIMITING=true
			if os.Getenv("TEST_RATE_LIMITING") == "true" {
				time.Sleep(100 * time.Millisecond)
			}

			return false
		}

		// Verhoog de counter
		cnt.count++

		return true
	}

	// Als er geen limiet is ingesteld, sta toe
	return true
}

// indexOf vindt de index van een karakter in een string
func indexOf(s string, char string) int {
	for i, c := range s {
		if string(c) == char {
			return i
		}
	}
	return -1
}

// GetCurrentCount geeft het huidige aantal voor een operatietype
func (rl *RateLimiter) GetCurrentCount(operationType string, ipAddress string) int {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	limit, exists := rl.limits[operationType]
	if !exists {
		return 0
	}

	now := time.Now()

	if limit.PerIP {
		ipMap, exists := rl.ipCounts[operationType]
		if !exists {
			return 0
		}

		ipCounter, exists := ipMap[ipAddress]
		if !exists || now.After(ipCounter.resetTime) {
			return 0
		}

		return ipCounter.count
	}

	// Globale teller
	globalCounter, exists := rl.globalCounts[operationType]
	if !exists || now.After(globalCounter.resetTime) {
		return 0
	}

	return globalCounter.count
}

// cleanupRoutine verwijdert verlopen tellers
func (rl *RateLimiter) cleanupRoutine() {
	for {
		select {
		case <-rl.cleanupTicker.C:
			rl.cleanup()
		case <-rl.done:
			return
		}
	}
}

// cleanup verwijdert verlopen tellers
func (rl *RateLimiter) cleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	// Opschonen van globale tellers
	for opType, counter := range rl.globalCounts {
		if now.After(counter.resetTime) {
			delete(rl.globalCounts, opType)
		}
	}

	// Opschonen van IP-specifieke tellers
	for opType, ipMap := range rl.ipCounts {
		for ip, counter := range ipMap {
			if now.After(counter.resetTime) {
				delete(ipMap, ip)
			}
		}

		// Verwijder lege maps
		if len(ipMap) == 0 {
			delete(rl.ipCounts, opType)
		}
	}
}

// Shutdown stopt de rate limiter en opschoning
func (rl *RateLimiter) Shutdown() {
	rl.cleanupTicker.Stop()
	rl.done <- true
}

// GetLimits geeft alle geconfigureerde limieten terug
func (rl *RateLimiter) GetLimits() map[string]RateLimit {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Maak een kopie van de limieten map
	limits := make(map[string]RateLimit, len(rl.limits))
	for k, v := range rl.limits {
		limits[k] = v
	}

	return limits
}

// GetCurrentValues haalt de huidige waarden op
func (rl *RateLimiter) GetCurrentValues() map[string]int {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	result := make(map[string]int)
	now := time.Now()

	// Verzamel alle actieve tellers
	for key, counter := range rl.globalCounts {
		// Controleer of de teller nog geldig is
		if !now.After(counter.resetTime) {
			result[key] = counter.count
		}
	}

	return result
}
