// Package tests contains test utilities and mock implementations
package tests

import (
	"context"
	"dklautomationgo/models"
	"dklautomationgo/services"
	"fmt"
	"sync"
	"time"

	"github.com/stretchr/testify/mock"
)

// mockEmailService implementeert de handlers.EmailServiceInterface
type mockEmailService struct {
	mu                    sync.Mutex
	contactEmailCalled    bool
	aanmeldingEmailCalled bool
	shouldFail            bool
	sentEmails            []services.EmailMessage
}

func newMockEmailService() *mockEmailService {
	return &mockEmailService{
		sentEmails: make([]services.EmailMessage, 0),
	}
}

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

func (m *mockEmailService) SendRegistrationEmail(data *models.RegistrationEmailData) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.aanmeldingEmailCalled = true
	if m.shouldFail {
		return fmt.Errorf("mock email error")
	}

	recipient := data.Participant.Email
	if data.ToAdmin {
		recipient = data.AdminEmail
	}

	m.sentEmails = append(m.sentEmails, services.EmailMessage{
		To:      recipient,
		Subject: "Aanmelding DKL",
		Body:    fmt.Sprintf("Aanmelding van %s", data.Participant.Naam),
	})
	return nil
}

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
type mockSMTP struct {
	mock.Mock
	mutex              sync.Mutex
	Dialer             services.SMTPDialer
	ShouldFail         bool
	SendCalled         bool
	SendRegCalled      bool
	SendWFCCalled      bool
	SendWithFromCalled bool
	LastFrom           string
	LastTo             string
	LastSubject        string
	LastBody           string
	SendError          error
}

func (m *mockSMTP) Send(msg *services.EmailMessage) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.SendCalled = true
	m.LastTo = msg.To
	m.LastSubject = msg.Subject
	m.LastBody = msg.Body
	args := m.Called(msg)
	if m.SendError != nil {
		return m.SendError
	}
	return args.Error(0)
}

func (m *mockSMTP) SendRegistration(msg *services.EmailMessage) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.SendRegCalled = true
	m.LastTo = msg.To
	m.LastSubject = msg.Subject
	m.LastBody = msg.Body
	args := m.Called(msg)
	if m.SendError != nil {
		return m.SendError
	}
	return args.Error(0)
}

func (m *mockSMTP) SendWFC(msg *services.EmailMessage) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.SendWFCCalled = true
	m.LastTo = msg.To
	m.LastSubject = msg.Subject
	m.LastBody = msg.Body
	args := m.Called(msg)
	if m.SendError != nil {
		return m.SendError
	}
	return args.Error(0)
}

// Implement SendWithFrom for the mock
func (m *mockSMTP) SendWithFrom(from string, msg *services.EmailMessage) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.SendWithFromCalled = true
	m.LastFrom = from // Store the specific from address used
	m.LastTo = msg.To
	m.LastSubject = msg.Subject
	m.LastBody = msg.Body
	args := m.Called(from, msg)
	if m.SendError != nil {
		return m.SendError
	}
	return args.Error(0)
}

// Update SendEmail signature for the mock
func (m *mockSMTP) SendEmail(to, subject, body string, fromAddress ...string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	finalFrom := ""
	if len(fromAddress) > 0 {
		finalFrom = fromAddress[0]
	}

	msg := &services.EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	if finalFrom == "" {
		return m.Send(msg)
	} else {
		return m.SendWithFrom(finalFrom, msg)
	}
}

func (m *mockSMTP) SendWFCEmail(to, subject, body string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	msg := &services.EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	}
	return m.SendWFC(msg)
}

// Reset clears the mock state
func (m *mockSMTP) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Mock = mock.Mock{}
	m.SendCalled = false
	m.SendRegCalled = false
	m.SendWFCCalled = false
	m.SendWithFromCalled = false
	m.LastTo = ""
	m.LastSubject = ""
	m.LastBody = ""
	m.LastFrom = ""
	m.SendError = nil
}

// Dial implements the SMTPDialer interface
func (m *mockSMTP) Dial(addr string) (services.SMTPClient, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	args := m.Called(addr)
	if m.ShouldFail {
		return nil, fmt.Errorf("mock SMTP dial error")
	}
	clientArg := args.Get(0)
	if clientArg == nil {
		return nil, args.Error(1)
	}
	if client, ok := clientArg.(services.SMTPClient); ok {
		return client, args.Error(1)
	}
	return nil, fmt.Errorf("mock Dial returned incorrect type: %T", clientArg)
}

// SetShouldFail configures the mock to return an error on the next Send/Dial
func (m *mockSMTP) SetShouldFail(shouldFail bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if shouldFail {
		m.SendError = fmt.Errorf("mock SMTP error")
	} else {
		m.SendError = nil
	}
}

// GetLastEmail returns the details of the last email sent
func (m *mockSMTP) GetLastEmail() (to, subject, body string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.LastTo, m.LastSubject, m.LastBody
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

// SendWFC implementeert de SendWFC methode van de SMTP client interface
func (m *MockSMTPClient) SendWFC(msg *services.EmailMessage) error {
	return m.Send(msg)
}

// SendWFCEmail implementeert de SendWFCEmail methode van de SMTP client interface
func (m *MockSMTPClient) SendWFCEmail(to, subject, body string) error {
	return m.SendEmail(to, subject, body)
}

// mockRateLimiter is een mock implementatie van de RateLimiterInterface
type mockRateLimiter struct {
	AllowFunc func(key string) bool
}

func newMockRateLimiter() *mockRateLimiter {
	return &mockRateLimiter{}
}

func (m *mockRateLimiter) Allow(key string) bool {
	if m.AllowFunc != nil {
		return m.AllowFunc(key)
	}
	return true // Default allow
}

func (m *mockRateLimiter) AllowEmail(emailType, recipient string) bool {
	// Simple mock: always allow or use AllowFunc if specific logic needed
	key := fmt.Sprintf("%s:%s", emailType, recipient)
	return m.Allow(key)
}

// Implement AddLimit for the mock (needed by some tests)
func (m *mockRateLimiter) AddLimit(operationType string, count int, period time.Duration, perIP bool) {
	// No-op for mock, or add specific test logic if needed
}

// Implement GetCurrentCount for the mock (to satisfy RateLimiterInterface)
func (m *mockRateLimiter) GetCurrentCount(operationType string, key string) int {
	// Return 0 or some specific test value if needed
	return 0
}

// Implement Shutdown for the mock
func (m *mockRateLimiter) Shutdown() {
	// No-op for mock
}

// Implement GetLimits for the mock
func (m *mockRateLimiter) GetLimits() map[string]services.RateLimit {
	// Return an empty map or some default test limits if needed
	return make(map[string]services.RateLimit)
}

// GetCurrentValues implements the RateLimiterInterface for the mock.
// Returns a dummy value for testing.
func (m *mockRateLimiter) GetCurrentValues() map[string]int {
	// Return a fixed dummy map or implement logic if needed for specific tests
	return map[string]int{"dummy_key": 0}
}

// mockPrometheusMetrics implementeert de PrometheusMetrics interface voor tests
type mockPrometheusMetrics struct {
	mock.Mock
	mutex        sync.Mutex
	latencies    map[string]float64
	emailsSent   map[string]int
	emailsFailed map[string]int
}

func newMockPrometheusMetrics() *mockPrometheusMetrics {
	return &mockPrometheusMetrics{
		latencies:    make(map[string]float64),
		emailsSent:   make(map[string]int),
		emailsFailed: make(map[string]int),
	}
}

func (m *mockPrometheusMetrics) RecordEmailSent(emailType, status string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.emailsSent == nil {
		m.emailsSent = make(map[string]int)
	}
	m.emailsSent[emailType+"_"+status]++
	m.Called(emailType, status)
}

func (m *mockPrometheusMetrics) RecordEmailFailed(emailType, reason string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.emailsFailed == nil {
		m.emailsFailed = make(map[string]int)
	}
	m.emailsFailed[emailType+"_"+reason]++
	m.Called(emailType, reason)
}

func (m *mockPrometheusMetrics) ObserveEmailLatency(emailType string, duration float64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.latencies == nil {
		m.latencies = make(map[string]float64)
	}
	m.latencies[emailType] = duration
	m.Called(emailType, duration)
}

// RecordRateLimitExceeded implements the PrometheusMetricsInterface
func (m *mockPrometheusMetrics) RecordRateLimitExceeded(operationType, limitType string) {
	m.Called(operationType, limitType)
}

func (m *mockPrometheusMetrics) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Mock = mock.Mock{}
	m.latencies = make(map[string]float64)
	m.emailsSent = make(map[string]int)
	m.emailsFailed = make(map[string]int)
}

// MockNotificationService implements the services.NotificationService interface for testing
type MockNotificationService struct {
	mock.Mock
}

// NewMockNotificationService creates a new mock notification service
func NewMockNotificationService() *MockNotificationService {
	return &MockNotificationService{}
}

// SendNotification mocks sending a notification
func (m *MockNotificationService) SendNotification(ctx context.Context, notification *models.Notification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

// CreateNotification mocks creating a notification
// V27: Updated to use strings instead of types
func (m *MockNotificationService) CreateNotification(
	ctx context.Context,
	notificationType string,
	priority string,
	title, message string,
) (*models.Notification, error) {
	args := m.Called(ctx, notificationType, priority, title, message)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Notification), args.Error(1)
}

// GetNotification mocks retrieving a notification by ID
func (m *MockNotificationService) GetNotification(ctx context.Context, id string) (*models.Notification, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Notification), args.Error(1)
}

// ListUnsentNotifications mocks listing unsent notifications
func (m *MockNotificationService) ListUnsentNotifications(ctx context.Context) ([]*models.Notification, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Notification), args.Error(1)
}

// ProcessUnsentNotifications mocks processing unsent notifications
func (m *MockNotificationService) ProcessUnsentNotifications(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Start mocks starting the notification service
func (m *MockNotificationService) Start() {
	m.Called()
}

// Stop mocks stopping the notification service
func (m *MockNotificationService) Stop() {
	m.Called()
}

// IsRunning mocks checking if the notification service is running
func (m *MockNotificationService) IsRunning() bool {
	args := m.Called()
	return args.Bool(0)
}

// --- Mock Email Sender (used for testing handlers) ---
type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendEmail(to, subject, body string, fromAddress ...string) error {
	args := m.Called(to, subject, body, fromAddress)
	return args.Error(0)
}

func (m *MockEmailSender) SendTemplateEmail(recipient, subject, templateName string, templateData map[string]interface{}, fromAddress ...string) error {
	args := m.Called(recipient, subject, templateName, templateData, fromAddress)
	return args.Error(0)
}

func (m *MockEmailSender) SendContactEmail(data *models.ContactEmailData) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockEmailSender) SendRegistrationEmail(data *models.RegistrationEmailData) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockEmailSender) SendWFCEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

func (m *MockEmailSender) SendWFCOrderEmail(data *models.WFCOrderEmailData) error {
	args := m.Called(data)
	return args.Error(0)
}

// --- WFCMockSMTPClient (als deze specifiek is) ---
// Als WFCMockSMTPClient bestaat en verschilt van mockSMTP,
// moet SendEmail en SendWithFrom hier ook aangepast/toegevoegd worden.
// Voorbeeld:
type WFCMockSMTPClient struct {
	mockSMTP // Embed mockSMTP of implementeer methoden opnieuw
}

// Dial implements the SMTPDialer interface
func (m *WFCMockSMTPClient) Dial() error {
	// Mock implementation, return nil or configured error
	return nil // Or return m.SendError if that's how errors are handled
}

// Zorg ervoor dat WFCMockSMTPClient ook voldoet aan de interface
// (Dit vereist mogelijk aanpassingen gebaseerd op de werkelijke implementatie)

// mockSMTPDialer implementeert de services.SMTPDialer interface voor tests
type mockSMTPDialer struct {
	ShouldFail bool
	DialCalled bool
	mutex      sync.Mutex
}

// Dial implementeert de Dial methode voor de mock dialer
func (m *mockSMTPDialer) Dial() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.DialCalled = true
	if m.ShouldFail {
		return fmt.Errorf("mock dialer error")
	}
	return nil
}

// Reset resets the mock dialer state
func (m *mockSMTPDialer) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.ShouldFail = false
	m.DialCalled = false
}

// MockAuthService is a mock implementation of services.AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(ctx context.Context, email, wachtwoord string) (string, string, error) {
	args := m.Called(ctx, email, wachtwoord)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockAuthService) ValidateToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) GetUserFromToken(ctx context.Context, token string) (*models.Gebruiker, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*models.Gebruiker), args.Error(1)
}

func (m *MockAuthService) HashPassword(wachtwoord string) (string, error) {
	args := m.Called(wachtwoord)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) VerifyPassword(hash, wachtwoord string) bool {
	args := m.Called(hash, wachtwoord)
	return args.Bool(0)
}

func (m *MockAuthService) ResetPassword(ctx context.Context, email, nieuwWachtwoord string) error {
	args := m.Called(ctx, email, nieuwWachtwoord)
	return args.Error(0)
}

func (m *MockAuthService) CreateUser(ctx context.Context, gebruiker *models.Gebruiker, password string) error {
	args := m.Called(ctx, gebruiker, password)
	return args.Error(0)
}

func (m *MockAuthService) ListUsers(ctx context.Context, limit, offset int) ([]*models.Gebruiker, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*models.Gebruiker), args.Error(1)
}

func (m *MockAuthService) GetUser(ctx context.Context, id string) (*models.Gebruiker, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Gebruiker), args.Error(1)
}

func (m *MockAuthService) UpdateUser(ctx context.Context, gebruiker *models.Gebruiker, password *string) error {
	args := m.Called(ctx, gebruiker, password)
	return args.Error(0)
}

func (m *MockAuthService) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAuthService) GenerateRefreshToken(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, string, error) {
	args := m.Called(ctx, refreshToken)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockAuthService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *MockAuthService) RevokeAllUserRefreshTokens(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthService) SearchUsers(ctx context.Context, query string, limit int) ([]*models.Gebruiker, error) {
	args := m.Called(ctx, query, limit)
	return args.Get(0).([]*models.Gebruiker), args.Error(1)
}

func (m *MockAuthService) GetParticipantByGebruikerID(ctx context.Context, gebruikerID string) (*models.Participant, error) {
	args := m.Called(ctx, gebruikerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Participant), args.Error(1)
}

func (m *MockAuthService) DeleteUserAccount(ctx context.Context, userID, password, reason string) error {
	args := m.Called(ctx, userID, password, reason)
	return args.Error(0)
}

func (m *MockAuthService) GetGebruikerByEmail(ctx context.Context, email string) (*models.Gebruiker, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Gebruiker), args.Error(1)
}

func (m *MockAuthService) GetParticipantByEmail(ctx context.Context, email string) ([]*models.Participant, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Participant), args.Error(1)
}

func (m *MockAuthService) IsEmailVerified(ctx context.Context, email, userType string) (bool, error) {
	args := m.Called(ctx, email, userType)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthService) RequestPasswordReset(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockAuthService) ResetPasswordWithToken(ctx context.Context, token, newPassword string) error {
	args := m.Called(ctx, token, newPassword)
	return args.Error(0)
}

func (m *MockAuthService) SendEmailVerification(ctx context.Context, email, userID, userType string) error {
	args := m.Called(ctx, email, userID, userType)
	return args.Error(0)
}

func (m *MockAuthService) VerifyEmailWithToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAuthService) ResendEmailVerification(ctx context.Context, email, userID, userType string) error {
	args := m.Called(ctx, email, userID, userType)
	return args.Error(0)
}

func (m *MockAuthService) RevokeAllUserAccessTokens(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthService) ValidateAccessToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAuthService) CreateSessionForLogin(ctx context.Context, userID, userType, accessToken, ipAddress, userAgent string) (*models.Session, error) {
	args := m.Called(ctx, userID, userType, accessToken, ipAddress, userAgent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Session), args.Error(1)
}

func (m *MockAuthService) ListUserSessions(ctx context.Context, userID string) ([]*models.Session, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Session), args.Error(1)
}

func (m *MockAuthService) RevokeSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockAuthService) RevokeAllUserSessions(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthService) RevokeOtherUserSessions(ctx context.Context, userID, currentSessionID string) error {
	args := m.Called(ctx, userID, currentSessionID)
	return args.Error(0)
}

func (m *MockAuthService) UpdateAccessTokenWithSession(ctx context.Context, token, sessionID string) error {
	args := m.Called(ctx, token, sessionID)
	return args.Error(0)
}

// MockPermissionService is a mock implementation of services.PermissionService
type MockPermissionService struct {
	mock.Mock
}

// NewMockPermissionService creates a new mock permission service
func NewMockPermissionService() *MockPermissionService {
	return &MockPermissionService{}
}

// HasPermission mocks checking if a user has a specific permission
func (m *MockPermissionService) HasPermission(ctx context.Context, userID, resource, action string) bool {
	args := m.Called(ctx, userID, resource, action)
	return args.Bool(0)
}

// GetUserPermissions mocks retrieving all permissions for a user
func (m *MockPermissionService) GetUserPermissions(ctx context.Context, userID string) ([]*models.UserPermission, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.UserPermission), args.Error(1)
}

// GetUserRoles mocks retrieving all active roles for a user
func (m *MockPermissionService) GetUserRoles(ctx context.Context, userID string) ([]*models.UserRole, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.UserRole), args.Error(1)
}

// AssignRole mocks assigning a role to a user
func (m *MockPermissionService) AssignRole(ctx context.Context, userID, roleID string, assignedBy *string) error {
	args := m.Called(ctx, userID, roleID, assignedBy)
	return args.Error(0)
}

// RevokeRole mocks revoking a role from a user
func (m *MockPermissionService) RevokeRole(ctx context.Context, userID, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

// CreateRole mocks creating a new role
func (m *MockPermissionService) CreateRole(ctx context.Context, role *models.RBACRole, createdBy *string) error {
	args := m.Called(ctx, role, createdBy)
	return args.Error(0)
}

// UpdateRole mocks updating a role
func (m *MockPermissionService) UpdateRole(ctx context.Context, role *models.RBACRole) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

// DeleteRole mocks deleting a role
func (m *MockPermissionService) DeleteRole(ctx context.Context, roleID string) error {
	args := m.Called(ctx, roleID)
	return args.Error(0)
}

// AssignPermissionToRole mocks assigning a permission to a role
func (m *MockPermissionService) AssignPermissionToRole(ctx context.Context, roleID, permissionID string, assignedBy *string) error {
	args := m.Called(ctx, roleID, permissionID, assignedBy)
	return args.Error(0)
}

// RevokePermissionFromRole mocks revoking a permission from a role
func (m *MockPermissionService) RevokePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	args := m.Called(ctx, roleID, permissionID)
	return args.Error(0)
}

// GetRoles mocks retrieving all roles
func (m *MockPermissionService) GetRoles(ctx context.Context, limit, offset int) ([]*models.RBACRole, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.RBACRole), args.Error(1)
}

// GetPermissions mocks retrieving all permissions
func (m *MockPermissionService) GetPermissions(ctx context.Context, limit, offset int) ([]*models.Permission, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Permission), args.Error(1)
}

// InvalidateUserCache mocks invalidating the cache for a user
func (m *MockPermissionService) InvalidateUserCache(userID string) {
	m.Called(userID)
}

// RefreshCache mocks refreshing all caches
func (m *MockPermissionService) RefreshCache(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// CanParticipantRegisterForEvent mocks checking if a participant can register for an event
func (m *MockPermissionService) CanParticipantRegisterForEvent(ctx context.Context, userID string) bool {
	args := m.Called(ctx, userID)
	return args.Bool(0)
}

// GetParticipantPermissionLevel mocks getting participant permission level
func (m *MockPermissionService) GetParticipantPermissionLevel(ctx context.Context, userID string) string {
	args := m.Called(ctx, userID)
	return args.String(0)
}

// HasParticipantAppAccess mocks checking if participant has app access
func (m *MockPermissionService) HasParticipantAppAccess(ctx context.Context, userID string) bool {
	args := m.Called(ctx, userID)
	return args.Bool(0)
}
