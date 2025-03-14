package repository

import (
	"dklautomationgo/logger"

	"gorm.io/gorm"
)

// RepositoryFactory bevat alle repositories
type RepositoryFactory struct {
	Contact            ContactRepository
	ContactAntwoord    ContactAntwoordRepository
	Aanmelding         AanmeldingRepository
	AanmeldingAntwoord AanmeldingAntwoordRepository
	EmailTemplate      EmailTemplateRepository
	VerzondEmail       VerzondEmailRepository
	Gebruiker          GebruikerRepository
	Migratie           MigratieRepository
}

// NewRepositoryFactory maakt een nieuwe repository factory
func NewRepositoryFactory(db *gorm.DB) *RepositoryFactory {
	logger.Info("Initialiseren repository factory")

	// Maak de basis repository
	baseRepo := NewPostgresRepository(db)

	// Maak alle repositories
	return &RepositoryFactory{
		Contact:            NewPostgresContactRepository(baseRepo),
		ContactAntwoord:    NewPostgresContactAntwoordRepository(baseRepo),
		Aanmelding:         NewPostgresAanmeldingRepository(baseRepo),
		AanmeldingAntwoord: NewPostgresAanmeldingAntwoordRepository(baseRepo),
		EmailTemplate:      NewPostgresEmailTemplateRepository(baseRepo),
		VerzondEmail:       NewPostgresVerzondEmailRepository(baseRepo),
		Gebruiker:          NewPostgresGebruikerRepository(baseRepo),
		Migratie:           NewPostgresMigratieRepository(baseRepo),
	}
}
