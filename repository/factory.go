package repository

import (
	"gorm.io/gorm"
)

// Repository is een overkoepelende struct die alle specifieke repositories bevat
type Repository struct {
	Contact            ContactRepository
	ContactAntwoord    ContactAntwoordRepository
	Aanmelding         AanmeldingRepository
	AanmeldingAntwoord AanmeldingAntwoordRepository
	Gebruiker          GebruikerRepository
	VerzondEmail       VerzondEmailRepository
	EmailTemplate      EmailTemplateRepository
	Migratie           MigratieRepository
	IncomingEmail      IncomingEmailRepository
	Notification       NotificationRepository
}

// NewRepository maakt een nieuwe Repository met concrete implementaties
func NewRepository(db *gorm.DB) *Repository {
	baseRepo := NewPostgresRepository(db)

	repo := &Repository{
		Contact:            NewPostgresContactRepository(baseRepo),
		ContactAntwoord:    NewPostgresContactAntwoordRepository(baseRepo),
		Aanmelding:         NewPostgresAanmeldingRepository(baseRepo),
		AanmeldingAntwoord: NewPostgresAanmeldingAntwoordRepository(baseRepo),
		Gebruiker:          NewPostgresGebruikerRepository(baseRepo),
		VerzondEmail:       NewPostgresVerzondEmailRepository(baseRepo),
		EmailTemplate:      NewPostgresEmailTemplateRepository(baseRepo),
		Migratie:           NewPostgresMigratieRepository(baseRepo),
		IncomingEmail:      NewPostgresIncomingEmailRepository(db),
		Notification:       NewPostgresNotificationRepository(baseRepo),
	}

	return repo
}
