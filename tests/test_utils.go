package tests

import (
	"context"
	"dklautomationgo/models"

	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of the AuthService interface
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) VerifyToken(token string) (*models.Gebruiker, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Gebruiker), args.Error(1)
}

func (m *MockAuthService) GenerateToken(gebruiker *models.Gebruiker) (string, error) {
	args := m.Called(gebruiker)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) IsAdmin(gebruiker *models.Gebruiker) bool {
	args := m.Called(gebruiker)
	return args.Bool(0)
}

// GetUserFromToken implements the AuthService interface
func (m *MockAuthService) GetUserFromToken(ctx context.Context, token string) (*models.Gebruiker, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Gebruiker), args.Error(1)
}

// Login implements the AuthService interface
func (m *MockAuthService) Login(ctx context.Context, email, wachtwoord string) (string, error) {
	args := m.Called(ctx, email, wachtwoord)
	return args.String(0), args.Error(1)
}

// ValidateToken implements the AuthService interface
func (m *MockAuthService) ValidateToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

// HashPassword implements the AuthService interface
func (m *MockAuthService) HashPassword(wachtwoord string) (string, error) {
	args := m.Called(wachtwoord)
	return args.String(0), args.Error(1)
}

// VerifyPassword implements the AuthService interface
func (m *MockAuthService) VerifyPassword(hash, wachtwoord string) bool {
	args := m.Called(hash, wachtwoord)
	return args.Bool(0)
}

// ResetPassword implements the AuthService interface
func (m *MockAuthService) ResetPassword(ctx context.Context, email, nieuwWachtwoord string) error {
	args := m.Called(ctx, email, nieuwWachtwoord)
	return args.Error(0)
}

// MockEmailService is a mock implementation of the EmailService interface
type MockEmailService struct {
	mock.Mock
}

// SendEmail implements the EmailService interface
func (m *MockEmailService) SendEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

// SendContactEmail implements the EmailService interface
func (m *MockEmailService) SendContactEmail(data *models.ContactEmailData) error {
	args := m.Called(data)
	return args.Error(0)
}

// SendAanmeldingEmail implements the EmailService interface
func (m *MockEmailService) SendAanmeldingEmail(data *models.AanmeldingEmailData) error {
	args := m.Called(data)
	return args.Error(0)
}

// SendBatchEmail implements the EmailService interface
func (m *MockEmailService) SendBatchEmail(batchKey string, recipients []string, subject, body string) error {
	args := m.Called(batchKey, recipients, subject, body)
	return args.Error(0)
}

// NoOpEmailService is a mock implementation that does nothing
type NoOpEmailService struct{}

func (m *NoOpEmailService) SendEmail(to, subject, body string) error {
	return nil
}

func (m *NoOpEmailService) SendContactEmail(data *models.ContactEmailData) error {
	return nil
}

func (m *NoOpEmailService) SendAanmeldingEmail(data *models.AanmeldingEmailData) error {
	return nil
}

func (m *NoOpEmailService) SendBatchEmail(batchKey string, recipients []string, subject, body string) error {
	return nil
}
