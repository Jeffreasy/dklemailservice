package tests

import (
	"dklautomationgo/models"

	"github.com/stretchr/testify/mock"
)

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
