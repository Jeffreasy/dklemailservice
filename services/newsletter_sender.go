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

	// Queue in batcher
	for _, sub := range subs {
		s.batcher.AddToBatch(batchKey, sub.Email, subject, "newsletter", data)
	}

	// Force flush if small batch
	if len(subs) < s.batcher.batchSize {
		time.Sleep(s.batcher.batchWindow)
	}

	logger.Info("Nieuwsbrief gequeued", "recipients", len(subs), "batch_id", nl.BatchID)
	return nil
}

func (s *NewsletterSender) SendManual(ctx context.Context, newsletterID string) error {
	// Get the newsletter
	nl, err := s.nlRepo.GetByID(ctx, newsletterID)
	if err != nil {
		return err
	}
	if nl == nil {
		return fmt.Errorf("newsletter not found: %s", newsletterID)
	}

	// Check if already sent
	if nl.SentAt != nil {
		return fmt.Errorf("newsletter already sent")
	}

	// Get subscribers
	subs, err := s.gebruikerRepo.GetNewsletterSubscribers(ctx)
	if err != nil {
		return err
	}
	if len(subs) == 0 {
		logger.Info("Geen subscribers voor nieuwsbrief")
		return nil
	}

	// Update batch ID
	batchKey := "newsletter_manual_" + newsletterID
	if err := s.nlRepo.UpdateBatchID(ctx, newsletterID, batchKey); err != nil {
		return err
	}

	// Queue in batcher
	data := map[string]interface{}{
		"Content": nl.Content,
	}

	for _, sub := range subs {
		s.batcher.AddToBatch(batchKey, sub.Email, nl.Subject, "newsletter", data)
	}

	// Force flush if small batch
	if len(subs) < s.batcher.batchSize {
		time.Sleep(s.batcher.batchWindow)
	}

	logger.Info("Manual nieuwsbrief gequeued", "newsletter_id", newsletterID, "recipients", len(subs), "batch_id", batchKey)
	return nil
}
