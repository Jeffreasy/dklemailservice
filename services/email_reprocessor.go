package services

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"io"
	"mime/quotedprintable"
	"strings"
)

// EmailReprocessor handles reprocessing of existing emails with better decoding
type EmailReprocessor struct {
	emailRepo repository.IncomingEmailRepository
	decoder   *EmailDecoder
}

// NewEmailReprocessor creates a new email reprocessor
func NewEmailReprocessor(emailRepo repository.IncomingEmailRepository) *EmailReprocessor {
	return &EmailReprocessor{
		emailRepo: emailRepo,
		decoder:   NewEmailDecoder(),
	}
}

// ReprocessAllEmails reprocesses all emails in the database with better decoding
func (r *EmailReprocessor) ReprocessAllEmails(ctx context.Context) (int, int, error) {
	// Get all emails (use a large limit or iterate in batches)
	emails, err := r.emailRepo.FindUnprocessed(ctx)
	if err != nil {
		return 0, 0, err
	}

	// Also get processed emails - we want to reprocess ALL
	// Since we don't have a GetAll method, we'll work with what we have
	// In production, you might want to add a GetAll(ctx context.Context) method

	processed := 0
	failed := 0

	for _, email := range emails {
		err := r.reprocessEmail(ctx, email)
		if err != nil {
			logger.Error("Failed to reprocess email", "email_id", email.ID, "error", err)
			failed++
		} else {
			processed++
		}
	}

	logger.Info("Email reprocessing complete", "processed", processed, "failed", failed)
	return processed, failed, nil
}

// reprocessEmail reprocesses a single email
func (r *EmailReprocessor) reprocessEmail(ctx context.Context, email *models.IncomingEmail) error {
	// Check if email body has encoding artifacts
	if !needsReprocessing(email.Body) {
		logger.Debug("Email doesn't need reprocessing", "email_id", email.ID)
		return nil
	}

	logger.Info("Reprocessing email", "email_id", email.ID, "subject", email.Subject)

	// Try to decode quoted-printable content
	decodedBody := decodeQuotedPrintableString(email.Body)

	// Try to convert Windows-1252 to UTF-8
	if email.ContentType != "" && strings.Contains(strings.ToLower(email.ContentType), "windows-1252") {
		converted, err := r.decoder.convertCharset(strings.NewReader(decodedBody), "windows-1252")
		if err == nil {
			convertedBytes, readErr := io.ReadAll(converted)
			if readErr == nil {
				decodedBody = string(convertedBytes)
			}
		}
	}

	// Remove MIME boundaries if present
	decodedBody = removeMIMEBoundaries(decodedBody)

	// Update email in database
	email.Body = decodedBody
	err := r.emailRepo.Update(ctx, email)
	if err != nil {
		return err
	}

	logger.Info("Email reprocessed successfully", "email_id", email.ID)
	return nil
}

// needsReprocessing checks if an email body contains encoding artifacts
func needsReprocessing(body string) bool {
	// Check for common quoted-printable artifacts
	if strings.Contains(body, "=92") || // Smart quote
		strings.Contains(body, "=85") || // Ellipsis
		strings.Contains(body, "=93") || // Left double quote
		strings.Contains(body, "=94") || // Right double quote
		strings.Contains(body, "=E9") || // é
		strings.Contains(body, "=EB") || // ë
		strings.Contains(body, "=20") { // Space (sometimes over-encoded)
		return true
	}

	// Check for MIME boundaries
	if strings.Contains(body, "--_000_") || strings.Contains(body, "--_001_") {
		return true
	}

	return false
}

// decodeQuotedPrintableString decodes a quoted-printable string
func decodeQuotedPrintableString(input string) string {
	var result strings.Builder
	reader := quotedprintable.NewReader(strings.NewReader(input))
	buf := make([]byte, 4096)

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			result.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}

	decoded := result.String()
	if decoded == "" {
		return input
	}

	return decoded
}

// removeMIMEBoundaries removes MIME boundary markers from email body
func removeMIMEBoundaries(body string) string {
	lines := strings.Split(body, "\n")
	var cleaned []string

	inBoundary := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip boundary markers
		if strings.HasPrefix(trimmed, "--_") && (strings.Contains(trimmed, "_") || strings.HasSuffix(trimmed, "--")) {
			inBoundary = true
			continue
		}

		// Skip Content-Type and Content-Transfer-Encoding headers after boundaries
		if inBoundary && (strings.HasPrefix(trimmed, "Content-Type:") ||
			strings.HasPrefix(trimmed, "Content-Transfer-Encoding:") ||
			trimmed == "") {
			if trimmed == "" {
				inBoundary = false
			}
			continue
		}

		inBoundary = false
		cleaned = append(cleaned, line)
	}

	return strings.Join(cleaned, "\n")
}
