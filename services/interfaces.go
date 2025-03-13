package services

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

type TemplateRenderer interface {
	RenderTemplate(name string, data interface{}) (string, error)
}
