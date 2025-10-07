package services

import (
	"dklautomationgo/models"
)

type NewsletterProcessor struct{}

func NewNewsletterProcessor() *NewsletterProcessor { return &NewsletterProcessor{} }

// Process converteert en filtert ruwe feed items naar ons model
func (p *NewsletterProcessor) Process(items []models.NewsItem) models.ProcessedNews {
	processed := models.ProcessedNews{Items: make([]models.NewsItem, 0), Summary: ""}
	for _, it := range items {
		if it.Title == "" || it.Link == "" {
			continue
		}
		processed.Items = append(processed.Items, it)
		if len(processed.Items) >= 20 {
			break
		}
	}
	if len(processed.Items) > 0 {
		processed.Summary = "Laatste updates en artikelen"
	}
	return processed
}
