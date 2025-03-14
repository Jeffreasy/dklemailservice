// Package tests contains test utilities and mock implementations
package tests

import (
	"dklautomationgo/models"
	"dklautomationgo/services"
	"fmt"
	"sync"
	"time"
)

// mockEmailService implementeert de handlers.EmailServiceInterface
//
//nolint:unused // These mocks are kept for future tests
type mockEmailService struct {
	mu                    sync.Mutex
	contactEmailCalled    bool
	aanmeldingEmailCalled bool
	shouldFail            bool
	sentEmails            []services.EmailMessage
}

//nolint:unused // These mocks are kept for future tests
func newMockEmailService() *mockEmailService {
	return &mockEmailService{
		sentEmails: make([]services.EmailMessage, 0),
	}
}

//nolint:unused // These mocks are kept for future tests
func (m *mockEmailService) SendContactEmail(data *models.ContactEmailData) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.contactEmailCalled = true
	if m.shouldFail {
		return fmt.Errorf("mock email error")
	}

	recipient := data.Contact.Email
	if data.ToAdmin {
		recipient = data.AdminEmail
	}

	m.sentEmails = append(m.sentEmails, services.EmailMessage{
		To:      recipient,
		Subject: "Contact Formulier",
		Body:    fmt.Sprintf("Contact van %s", data.Contact.Naam),
	})
	return nil
}

//nolint:unused // These mocks are kept for future tests
func (m *mockEmailService) SendAanmeldingEmail(data *models.AanmeldingEmailData) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.aanmeldingEmailCalled = true
	if m.shouldFail {
		return fmt.Errorf("mock email error")
	}

	recipient := data.Aanmelding.Email
	if data.ToAdmin {
		recipient = data.AdminEmail
	}

	m.sentEmails = append(m.sentEmails, services.EmailMessage{
		To:      recipient,
		Subject: "Aanmelding DKL",
		Body:    fmt.Sprintf("Aanmelding van %s voor %s", data.Aanmelding.Naam, data.Aanmelding.Rol),
	})
	return nil
}

//nolint:unused // These mocks are kept for future tests
func (m *mockEmailService) SendEmail(recipient, subject, body string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldFail {
		return fmt.Errorf("mock email error")
	}

	m.sentEmails = append(m.sentEmails, services.EmailMessage{
		To:      recipient,
		Subject: subject,
		Body:    body,
	})
	return nil
}

//nolint:unused // These mocks are kept for future tests
func (m *mockEmailService) SendBatchEmail(batchKey string, recipients []string, subject, body string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldFail {
		return fmt.Errorf("mock email error")
	}

	for _, recipient := range recipients {
		m.sentEmails = append(m.sentEmails, services.EmailMessage{
			To:      recipient,
			Subject: subject,
			Body:    body,
		})
	}
	return nil
}

// mockSMTP implementeert zowel SMTPClient als SMTPDialer interfaces
//
//nolint:unused // These mocks are kept for future tests
type mockSMTP struct {
	mu         sync.Mutex
	sentEmails []*services.EmailMessage
	shouldFail bool
	failFirst  bool
	firstSent  bool
}

//nolint:unused // These mocks are kept for future tests
func (m *mockSMTP) Send(msg *services.EmailMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if msg.To == "" {
		return fmt.Errorf("invalid recipient")
	}

	if m.failFirst && !m.firstSent {
		m.firstSent = true
		return fmt.Errorf("mock smtp error")
	}

	if m.shouldFail {
		return fmt.Errorf("mock email error")
	}

	m.sentEmails = append(m.sentEmails, msg)
	return nil
}

//nolint:unused // These mocks are kept for future tests
func (m *mockSMTP) SendRegistration(msg *services.EmailMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if msg.To == "" {
		return fmt.Errorf("invalid recipient")
	}

	if m.shouldFail {
		return fmt.Errorf("mock email error")
	}

	m.sentEmails = append(m.sentEmails, msg)
	return nil
}

//nolint:unused // These mocks are kept for future tests
func (m *mockSMTP) Dial() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldFail {
		return fmt.Errorf("mock dialer error")
	}
	return nil
}

//nolint:unused // These mocks are kept for future tests
func (m *mockSMTP) SetShouldFail(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = fail
}

//nolint:unused // These mocks are kept for future tests
func (m *mockSMTP) SetFailFirst(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failFirst = fail
	m.firstSent = false
}

//nolint:unused // These mocks are kept for future tests
func (m *mockSMTP) GetLastEmail() *services.EmailMessage {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.sentEmails) == 0 {
		return nil
	}
	return m.sentEmails[len(m.sentEmails)-1]
}

//nolint:unused // These mocks are kept for future tests
func (m *mockSMTP) GetSentEmails() []*services.EmailMessage {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sentEmails
}

//nolint:unused // These mocks are kept for future tests
func newMockSMTP() *mockSMTP {
	return &mockSMTP{
		sentEmails: make([]*services.EmailMessage, 0),
	}
}

// SendEmail is een helper functie voor backwards compatibility
//
//nolint:unused // These mocks are kept for future tests
func (m *mockSMTP) SendEmail(to, subject, body string) error {
	msg := &services.EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	}
	return m.Send(msg)
}

// MockSMTPClient is een mock implementatie van de SMTP client voor tests
type MockSMTPClient struct {
	SendEmailCallCount int
	LastTo             string
	LastSubject        string
	LastBody           string
	ErrorToReturn      error
}

// SendEmail implementeert de SendEmail methode van de SMTP client interface
func (m *MockSMTPClient) SendEmail(to, subject, body string) error {
	m.SendEmailCallCount++
	m.LastTo = to
	m.LastSubject = subject
	m.LastBody = body
	return m.ErrorToReturn
}

// Send implementeert de Send methode van de SMTP client interface
func (m *MockSMTPClient) Send(msg *services.EmailMessage) error {
	return m.SendEmail(msg.To, msg.Subject, msg.Body)
}

// SendRegistration implementeert de SendRegistration methode van de SMTP client interface
func (m *MockSMTPClient) SendRegistration(msg *services.EmailMessage) error {
	return m.Send(msg)
}

// mockRateLimiter is een mock implementatie van de RateLimiterInterface
//
//nolint:unused // These mocks are kept for future tests
type mockRateLimiter struct {
	shouldLimit bool
	limits      map[string]struct {
		limit    int
		period   time.Duration
		perEmail bool
	}
	currentCounts map[string]int
}

//nolint:unused // These mocks are kept for future tests
func newMockRateLimiter() *mockRateLimiter {
	return &mockRateLimiter{
		shouldLimit: false,
		limits: make(map[string]struct {
			limit    int
			period   time.Duration
			perEmail bool
		}),
		currentCounts: make(map[string]int),
	}
}

//nolint:unused // These mocks are kept for future tests
func (m *mockRateLimiter) AllowEmail(operation, email string) bool {
	return !m.shouldLimit
}

//nolint:unused // These mocks are kept for future tests
func (m *mockRateLimiter) GetLimits() map[string]services.RateLimit {
	limits := make(map[string]services.RateLimit)
	for op, l := range m.limits {
		limits[op] = services.RateLimit{
			Count:  l.limit,
			Period: l.period,
			PerIP:  l.perEmail,
		}
	}
	return limits
}

//nolint:unused // These mocks are kept for future tests
func (m *mockRateLimiter) GetCurrentCount(operationType string, key string) int {
	return m.currentCounts[operationType]
}

//nolint:unused // These mocks are kept for future tests
func (m *mockRateLimiter) SetShouldLimit(limit bool) {
	m.shouldLimit = limit
}

//nolint:unused // These mocks are kept for future tests
func (m *mockRateLimiter) AddLimit(operation string, limit int, period time.Duration, perEmail bool) {
	m.limits[operation] = struct {
		limit    int
		period   time.Duration
		perEmail bool
	}{
		limit:    limit,
		period:   period,
		perEmail: perEmail,
	}
}

// SetCurrentCount is een helper functie voor tests
//
//nolint:unused // These mocks are kept for future tests
func (m *mockRateLimiter) SetCurrentCount(operationType string, count int) {
	m.currentCounts[operationType] = count
}

// mockPrometheusMetrics implementeert de PrometheusMetrics interface voor tests
//
//nolint:unused // These mocks are kept for future tests
type mockPrometheusMetrics struct {
	mu           sync.Mutex
	emailsSent   int
	emailsFailed int
	latencies    map[string][]float64
}

//nolint:unused // These mocks are kept for future tests
func newMockPrometheusMetrics() *mockPrometheusMetrics {
	return &mockPrometheusMetrics{
		latencies: make(map[string][]float64),
	}
}

//nolint:unused // These mocks are kept for future tests
func (m *mockPrometheusMetrics) RecordEmailSent(emailType, status string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.emailsSent++
}

//nolint:unused // These mocks are kept for future tests
func (m *mockPrometheusMetrics) RecordEmailFailed(emailType, reason string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.emailsFailed++
}

//nolint:unused // These mocks are kept for future tests
func (m *mockPrometheusMetrics) ObserveEmailLatency(emailType string, duration float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.latencies[emailType] == nil {
		m.latencies[emailType] = make([]float64, 0)
	}
	m.latencies[emailType] = append(m.latencies[emailType], duration)
}

//nolint:unused // These mocks are kept for future tests
func (m *mockPrometheusMetrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.emailsSent = 0
	m.emailsFailed = 0
	m.latencies = make(map[string][]float64)
}
