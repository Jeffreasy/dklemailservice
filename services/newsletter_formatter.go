package services

import (
	"bytes"
	"dklautomationgo/models"
	"fmt"
)

type NewsletterFormatter struct {
	emailSvc *EmailService
}

func NewNewsletterFormatter(es *EmailService) *NewsletterFormatter {
	return &NewsletterFormatter{emailSvc: es}
}

func (f *NewsletterFormatter) Format(processed *models.ProcessedNews, subject string) (string, error) {
	tmplName := "newsletter"
	data := map[string]interface{}{
		"Items":   processed.Items,
		"Summary": processed.Summary,
	}
	var body bytes.Buffer
	tmpl := f.emailSvc.GetTemplate(tmplName)
	if tmpl == nil {
		return "", fmt.Errorf("template newsletter not found")
	}
	if err := tmpl.Execute(&body, data); err != nil {
		return "", err
	}
	return body.String(), nil
}
