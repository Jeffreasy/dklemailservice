package services

import (
	"dklautomationgo/logger"
	"sync"
	"sync/atomic"
	"time"
)

// EmailMetrics houdt statistieken bij over verzonden emails
type EmailMetrics struct {
	totalEmails   int64
	successEmails int64
	failedEmails  int64
	emailsByType  map[string]int64
	mutex         sync.RWMutex
	startTime     time.Time
	resetInterval time.Duration
	lastResetTime time.Time
}

// NewEmailMetrics creëert een nieuwe email metrics tracker
func NewEmailMetrics(resetInterval time.Duration) *EmailMetrics {
	now := time.Now()
	return &EmailMetrics{
		emailsByType:  make(map[string]int64),
		startTime:     now,
		resetInterval: resetInterval,
		lastResetTime: now,
	}
}

// RecordEmailSent registreert een succesvolle email
func (m *EmailMetrics) RecordEmailSent(emailType string) {
	atomic.AddInt64(&m.totalEmails, 1)
	atomic.AddInt64(&m.successEmails, 1)

	m.mutex.Lock()
	m.emailsByType[emailType]++
	m.mutex.Unlock()

	m.checkAndResetPeriodic()
}

// RecordEmailFailed registreert een mislukte email
func (m *EmailMetrics) RecordEmailFailed(emailType string) {
	atomic.AddInt64(&m.totalEmails, 1)
	atomic.AddInt64(&m.failedEmails, 1)

	m.checkAndResetPeriodic()
}

// GetTotalEmails geeft het totaal aantal verzonden emails
func (m *EmailMetrics) GetTotalEmails() int64 {
	return atomic.LoadInt64(&m.totalEmails)
}

// GetSuccessRate berekent het succespercentage
func (m *EmailMetrics) GetSuccessRate() float64 {
	total := atomic.LoadInt64(&m.totalEmails)
	if total == 0 {
		return 100.0
	}

	success := atomic.LoadInt64(&m.successEmails)
	return float64(success) / float64(total) * 100.0
}

// GetEmailsByType geeft een kopie van de emails-per-type map
func (m *EmailMetrics) GetEmailsByType() map[string]int64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[string]int64, len(m.emailsByType))
	for k, v := range m.emailsByType {
		result[k] = v
	}

	return result
}

// Reset stelt alle tellers terug naar 0
func (m *EmailMetrics) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.resetCountersLocked()
	m.lastResetTime = time.Now()
}

// LogMetrics logt de huidige metrics
func (m *EmailMetrics) LogMetrics() {
	total := m.GetTotalEmails()
	success := atomic.LoadInt64(&m.successEmails)
	failed := atomic.LoadInt64(&m.failedEmails)
	successRate := m.GetSuccessRate()

	logger.Info("Email metrics",
		"total", total,
		"success", success,
		"failed", failed,
		"success_rate", successRate,
		"since", m.lastResetTime.Format(time.RFC3339),
	)

	// Log details per email type
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for emailType, count := range m.emailsByType {
		logger.Info("Email type metrics",
			"type", emailType,
			"count", count,
		)
	}
}

// checkAndResetPeriodic reset metrics als de reset interval is verstreken
func (m *EmailMetrics) checkAndResetPeriodic() {
	if m.resetInterval <= 0 {
		return
	}

	m.mutex.RLock()
	shouldReset := time.Since(m.lastResetTime) >= m.resetInterval
	m.mutex.RUnlock()

	if shouldReset {
		m.LogMetrics()
		m.Reset()
	}
}

// CheckAndResetIfNeeded forceert een check of een reset nodig is
// Deze methode is vooral nuttig voor tests
func (m *EmailMetrics) CheckAndResetIfNeeded() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.checkAndResetIfNeeded()
}

// De interne methode blijft privé
func (m *EmailMetrics) checkAndResetIfNeeded() {
	elapsed := time.Since(m.lastResetTime)
	if elapsed >= m.resetInterval {
		m.resetCountersLocked()
		m.lastResetTime = time.Now()
	}
}

// resetCountersLocked reset alle tellers naar 0
// BELANGRIJK: deze methode gaat ervan uit dat de mutex al is vergrendeld
func (m *EmailMetrics) resetCountersLocked() {
	atomic.StoreInt64(&m.totalEmails, 0)
	atomic.StoreInt64(&m.successEmails, 0)
	atomic.StoreInt64(&m.failedEmails, 0)
	m.emailsByType = make(map[string]int64)
}
