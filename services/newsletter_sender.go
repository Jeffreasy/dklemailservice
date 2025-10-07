package services

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"fmt"
	"time"
)

type NewsletterSender struct {
	emailSvc      *EmailService
	batcher       *EmailBatcher
	gebruikerRepo repository.GebruikerRepository
	nlRepo        repository.NewsletterRepository
	notifSvc      NotificationService
}

func NewNewsletterSender(es *EmailService, eb *EmailBatcher, gr repository.GebruikerRepository,
	nr repository.NewsletterRepository, ns NotificationService) *NewsletterSender {
	return &NewsletterSender{emailSvc: es, batcher: eb, gebruikerRepo: gr, nlRepo: nr, notifSvc: ns}
}

func (s *NewsletterSender) Send(ctx context.Context, content, subject string) error {
	subs, err := s.gebruikerRepo.GetNewsletterSubscribers(ctx)
	if err != nil {
		return err
	}
	if len(subs) == 0 {
		logger.Info("Geen subscribers voor nieuwsbrief")
		return nil
	}

	// Save newsletter record
	nl := &models.Newsletter{Subject: subject, Content: content}
	if err := s.nlRepo.Create(ctx, nl); err != nil {
		return err
	}

	batchKey := "newsletter_daily"
	data := map[string]interface{}{"Summary": "", "Items": []models.NewsItem{}}
	data["Content"] = content

	// Queue in batcher with specific from address
	newsletterFromAddress := "nieuwsbrief@dekoninklijkeloop.nl"
	for _, sub := range subs {
		s.batcher.AddToBatch(batchKey, sub.Email, subject, "newsletter", data, newsletterFromAddress)
	}

	// Force immediate sending for small newsletter batches
	if len(subs) < s.batcher.batchSize {
		logger.Info("SendManual: Small batch detected, forcing immediate send", "subscriber_count", len(subs), "batch_size", s.batcher.batchSize)
		// Force immediate flush of this specific batch
		s.batcher.FlushBatch(batchKey)
	}

	// Mark newsletter as sent
	sentAt := time.Now()
	if err := s.nlRepo.MarkSent(ctx, nl.ID, sentAt); err != nil {
		logger.Error("Fout bij markeren nieuwsbrief als verzonden", "error", err, "newsletter_id", nl.ID)
		return err
	}

	logger.Info("Nieuwsbrief verzonden", "newsletter_id", nl.ID, "recipients", len(subs), "batch_id", nl.BatchID, "sent_at", sentAt)
	return nil
}

func (s *NewsletterSender) SendManual(ctx context.Context, newsletterID string) error {
	logger.Info("SendManual: Getting newsletter", "newsletter_id", newsletterID)

	// Get the newsletter
	nl, err := s.nlRepo.GetByID(ctx, newsletterID)
	if err != nil {
		logger.Error("SendManual: Error getting newsletter", "error", err, "newsletter_id", newsletterID)
		return err
	}
	if nl == nil {
		logger.Warn("SendManual: Newsletter not found", "newsletter_id", newsletterID)
		return fmt.Errorf("newsletter not found: %s", newsletterID)
	}

	logger.Info("SendManual: Newsletter found", "newsletter_id", newsletterID, "subject", nl.Subject)

	// Check if already sent
	if nl.SentAt != nil {
		return fmt.Errorf("newsletter already sent")
	}

	// Get subscribers
	logger.Info("SendManual: Getting newsletter subscribers")
	subs, err := s.gebruikerRepo.GetNewsletterSubscribers(ctx)
	if err != nil {
		logger.Error("SendManual: Error getting subscribers", "error", err)
		return err
	}
	if len(subs) == 0 {
		logger.Info("SendManual: Geen subscribers voor nieuwsbrief")
		return nil
	}

	logger.Info("SendManual: Found subscribers", "count", len(subs))

	// Update batch ID
	batchKey := "newsletter_manual_" + newsletterID
	if err := s.nlRepo.UpdateBatchID(ctx, newsletterID, batchKey); err != nil {
		return err
	}

	// Queue in batcher
	data := map[string]interface{}{
		"Content": nl.Content,
	}

	// Queue in batcher with specific from address
	newsletterFromAddress := "nieuwsbrief@dekoninklijkeloop.nl"
	logger.Info("SendManual: Queueing emails in batcher", "batch_key", batchKey, "from_address", newsletterFromAddress)

	for _, sub := range subs {
		s.batcher.AddToBatch(batchKey, sub.Email, nl.Subject, "newsletter", data, newsletterFromAddress)
	}

	// Force immediate sending for small newsletter batches
	if len(subs) < s.batcher.batchSize {
		logger.Info("Send: Small batch detected, forcing immediate send", "subscriber_count", len(subs), "batch_size", s.batcher.batchSize)
		// Force immediate flush of this specific batch
		s.batcher.FlushBatch(batchKey)
	}

	// Mark newsletter as sent
	sentAt := time.Now()
	if err := s.nlRepo.MarkSent(ctx, newsletterID, sentAt); err != nil {
		logger.Error("Fout bij markeren nieuwsbrief als verzonden", "error", err, "newsletter_id", newsletterID)
		return err
	}

	logger.Info("Manual nieuwsbrief verzonden", "newsletter_id", newsletterID, "recipients", len(subs), "batch_id", batchKey, "sent_at", sentAt)
	return nil
}
