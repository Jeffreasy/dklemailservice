package services

import (
	"context"
	"dklautomationgo/logger"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
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
	redisClient       *redis.Client
	useRedis          bool
}

type counter struct {
	count     int
	resetTime time.Time
}

// NewRateLimiter creëert een nieuwe rate limiter
func NewRateLimiter(prometheusMetrics *PrometheusMetrics) *RateLimiter {
	return NewRateLimiterWithRedis(prometheusMetrics, nil)
}

// NewRateLimiterWithRedis creëert een nieuwe rate limiter met optionele Redis ondersteuning
func NewRateLimiterWithRedis(prometheusMetrics *PrometheusMetrics, redisClient *redis.Client) *RateLimiter {
	rl := &RateLimiter{
		globalCounts:      make(map[string]*counter),
		ipCounts:          make(map[string]map[string]*counter),
		limits:            make(map[string]RateLimit),
		done:              make(chan bool),
		prometheusMetrics: prometheusMetrics,
		redisClient:       redisClient,
		useRedis:          redisClient != nil,
	}

	// Test Redis verbinding indien beschikbaar
	if rl.useRedis {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if _, err := rl.redisClient.Ping(ctx).Result(); err != nil {
			logger.Warn("Redis connection failed, falling back to in-memory rate limiting", "error", err)
			rl.useRedis = false
		} else {
			logger.Info("Redis rate limiting enabled")
		}
	}

	// Start periodieke opschoning van verlopen limieten (alleen voor in-memory)
	if !rl.useRedis {
		rl.cleanupTicker = time.NewTicker(10 * time.Minute)
		go rl.cleanupRoutine()
	}

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
	// Controleer globale rate limiting
	if limit, exists := r.limits[emailType]; exists {
		// Bereken hoeveel requests in deze periode zijn toegestaan
		key := emailType
		if userID != "" && limit.PerIP {
			key = fmt.Sprintf("%s:%s", emailType, userID)
		}

		if r.useRedis {
			return r.allowWithRedis(key, limit, emailType, userID)
		} else {
			return r.allowWithMemory(key, limit, emailType, userID)
		}
	}

	// Geen rate limiting voor dit type
	return true
}

// allowWithRedis gebruikt Redis voor distributed rate limiting
func (r *RateLimiter) allowWithRedis(key string, limit RateLimit, emailType, userID string) bool {
	ctx := context.Background()

	// Gebruik Redis sorted sets voor sliding window rate limiting
	// Voeg huidige timestamp toe aan de set
	now := time.Now().Unix()
	windowStart := now - int64(limit.Period.Seconds())

	// Verwijder oude entries buiten het window
	r.redisClient.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))

	// Tel huidige entries in het window
	count, err := r.redisClient.ZCard(ctx, key).Result()
	if err != nil {
		logger.Error("Redis ZCard failed, falling back to memory", "error", err, "key", key)
		return r.allowWithMemory(key, limit, emailType, userID)
	}

	// Check tegen de limiet
	if count >= int64(limit.Count) {
		// Overschreden
		logger.Warn("Redis rate limit overschreden",
			"operation", emailType,
			"limit", limit.Count,
			"period", limit.Period,
			"current_count", count)

		// Prometheus metrics registreren
		limitType := "global"
		if userID != "" {
			limitType = "per_user"
		}
		if r.prometheusMetrics != nil {
			r.prometheusMetrics.RecordRateLimitExceeded(emailType, limitType)
		}

		return false
	}

	// Voeg nieuwe request toe aan het window
	err = r.redisClient.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: strconv.FormatInt(now, 10) + ":" + fmt.Sprintf("%d", time.Now().Nanosecond()),
	}).Err()

	if err != nil {
		logger.Error("Redis ZAdd failed, falling back to memory", "error", err, "key", key)
		return r.allowWithMemory(key, limit, emailType, userID)
	}

	// Stel TTL in op de limiet periode + buffer
	r.redisClient.Expire(ctx, key, limit.Period+time.Minute)

	return true
}

// allowWithMemory gebruikt in-memory rate limiting (originele implementatie)
func (r *RateLimiter) allowWithMemory(key string, limit RateLimit, emailType, userID string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()

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
		logger.Warn("Memory rate limit overschreden",
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
			time.Sleep(100 * time.Millisecond)
		}

		return false
	}

	// Voor tests: voeg progressieve vertraging toe als TEST_RATE_LIMITING=true
	if os.Getenv("TEST_RATE_LIMITING") == "true" && cnt.count > 0 {
		time.Sleep(time.Duration(cnt.count) * 100 * time.Millisecond)
	}

	// Toegestaan, incrementeer counter
	cnt.count++
	return true
}

// Allow controleert of een operatie is toegestaan op basis van de rate limit
func (r *RateLimiter) Allow(key string) bool {
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

	// Controleer globale rate limiting
	if limit, exists := r.limits[operationType]; exists {
		// Bereken hoeveel requests in deze periode zijn toegestaan
		limitKey := operationType
		if userID != "" && limit.PerIP {
			limitKey = fmt.Sprintf("%s:%s", operationType, userID)
		}

		if r.useRedis {
			return r.allowWithRedis(limitKey, limit, operationType, userID)
		} else {
			return r.allowWithMemory(limitKey, limit, operationType, userID)
		}
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
	limit, exists := rl.limits[operationType]
	if !exists {
		return 0
	}

	if rl.useRedis {
		ctx := context.Background()
		key := operationType
		if limit.PerIP && ipAddress != "" {
			key = fmt.Sprintf("%s:%s", operationType, ipAddress)
		}

		// Tel entries in het huidige window
		now := time.Now().Unix()
		windowStart := now - int64(limit.Period.Seconds())

		count, err := rl.redisClient.ZCount(ctx, key, strconv.FormatInt(windowStart, 10), strconv.FormatInt(now+1, 10)).Result()
		if err != nil {
			logger.Error("Redis ZCount failed", "error", err, "key", key)
			return 0
		}

		return int(count)
	} else {
		rl.mutex.Lock()
		defer rl.mutex.Unlock()

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
	if !rl.useRedis {
		rl.cleanupTicker.Stop()
		rl.done <- true
	}

	// Redis client wordt afgesloten door de caller
	logger.Info("Rate limiter shutdown complete", "use_redis", rl.useRedis)
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
	if rl.useRedis {
		// Voor Redis implementatie, geef een lege map terug voor nu
		// In productie zou je alle Redis keys kunnen scannen, maar dat is duur
		logger.Debug("GetCurrentValues not implemented for Redis mode")
		return make(map[string]int)
	} else {
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
}
