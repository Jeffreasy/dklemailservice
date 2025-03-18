// Package tests contains test utilities and mock implementations
package tests

import (
	"context"
	"dklautomationgo/models"
	"dklautomationgo/services"
	"errors"
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
func (m *mockRateLimiter) Allow(key string) bool {
	return !m.shouldLimit
}

//nolint:unused // These mocks are kept for future tests
func (m *mockRateLimiter) GetLimits() map[string]services.RateLimit {
	result := make(map[string]services.RateLimit)
	for op, limit := range m.limits {
		result[op] = services.RateLimit{
			Count:  limit.limit,
			Period: limit.period,
			PerIP:  limit.perEmail,
		}
	}
	return result
}

//nolint:unused // These mocks are kept for future tests
func (m *mockRateLimiter) GetCurrentValues() map[string]int {
	return m.currentCounts
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

// MockNotificationService implements the services.NotificationService interface for testing
type MockNotificationService struct {
	shouldFail    bool
	calls         map[string]int
	notifications []*models.Notification
}

// NewMockNotificationService creates a new mock notification service
func NewMockNotificationService() *MockNotificationService {
	return &MockNotificationService{
		shouldFail:    false,
		calls:         make(map[string]int),
		notifications: make([]*models.Notification, 0),
	}
}

// SendNotification mocks sending a notification
func (m *MockNotificationService) SendNotification(ctx context.Context, notification *models.Notification) error {
	m.calls["SendNotification"]++

	if m.shouldFail {
		return errors.New("mock notification send failure")
	}

	notification.Sent = true
	return nil
}

// CreateNotification mocks creating a notification
func (m *MockNotificationService) CreateNotification(
	ctx context.Context,
	notificationType models.NotificationType,
	priority models.NotificationPriority,
	title, message string,
) (*models.Notification, error) {
	m.calls["CreateNotification"]++

	if m.shouldFail {
		return nil, errors.New("mock notification creation failure")
	}

	notification := &models.Notification{
		Type:     notificationType,
		Priority: priority,
		Title:    title,
		Message:  message,
		Sent:     false,
	}

	m.notifications = append(m.notifications, notification)
	return notification, nil
}

// GetNotification mocks retrieving a notification by ID
func (m *MockNotificationService) GetNotification(ctx context.Context, id string) (*models.Notification, error) {
	m.calls["GetNotification"]++

	if m.shouldFail {
		return nil, errors.New("mock get notification failure")
	}

	// Simple mock implementation that doesn't actually check the ID
	if len(m.notifications) > 0 {
		return m.notifications[0], nil
	}

	return nil, nil
}

// ListUnsentNotifications mocks listing unsent notifications
func (m *MockNotificationService) ListUnsentNotifications(ctx context.Context) ([]*models.Notification, error) {
	m.calls["ListUnsentNotifications"]++

	if m.shouldFail {
		return nil, errors.New("mock list notifications failure")
	}

	// Filter for unsent notifications
	unsent := make([]*models.Notification, 0)
	for _, n := range m.notifications {
		if !n.Sent {
			unsent = append(unsent, n)
		}
	}

	return unsent, nil
}

// ProcessUnsentNotifications mocks processing unsent notifications
func (m *MockNotificationService) ProcessUnsentNotifications(ctx context.Context) error {
	m.calls["ProcessUnsentNotifications"]++

	if m.shouldFail {
		return errors.New("mock process notifications failure")
	}

	// Mark all as sent
	for _, n := range m.notifications {
		if !n.Sent {
			n.Sent = true
		}
	}

	return nil
}

// Start mocks starting the notification service
func (m *MockNotificationService) Start() {
	m.calls["Start"]++
}

// Stop mocks stopping the notification service
func (m *MockNotificationService) Stop() {
	m.calls["Stop"]++
}

// IsRunning mocks checking if the notification service is running
func (m *MockNotificationService) IsRunning() bool {
	m.calls["IsRunning"]++
	return true
}
