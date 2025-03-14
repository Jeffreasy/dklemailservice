package database

import (
	"context"
	"dklautomationgo/database/migrations"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// MigrationManager beheert database migraties
type MigrationManager struct {
	db       *gorm.DB
	migrRepo repository.MigratieRepository
}

// NewMigrationManager maakt een nieuwe migratie manager
func NewMigrationManager(db *gorm.DB, migrRepo repository.MigratieRepository) *MigrationManager {
	return &MigrationManager{
		db:       db,
		migrRepo: migrRepo,
	}
}

// MigrateDatabase voert alle migraties uit
func (m *MigrationManager) MigrateDatabase() error {
	logger.Info("Database migratie gestart")

	// Maak de migraties tabel aan als deze nog niet bestaat
	if err := m.db.AutoMigrate(&models.Migratie{}); err != nil {
		return fmt.Errorf("fout bij aanmaken migraties tabel: %w", err)
	}

	// Voer SQL migraties uit
	if err := migrations.RunSQLMigrations(m.db); err != nil {
		return fmt.Errorf("fout bij uitvoeren SQL migraties: %w", err)
	}

	// Voer alle migraties uit
	if err := m.createTables(); err != nil {
		return err
	}

	logger.Info("Database migratie voltooid")
	return nil
}

// createTables maakt alle tabellen aan
func (m *MigrationManager) createTables() error {
	// Controleer of de migratie al is uitgevoerd
	ctx := context.Background()
	migratie, err := m.migrRepo.GetByVersie(ctx, "1.0.0")
	if err != nil {
		return fmt.Errorf("fout bij controleren migratie: %w", err)
	}

	if migratie != nil {
		logger.Info("Migratie 1.0.0 is al uitgevoerd", "toegepast", migratie.Toegepast)
		return nil
	}

	// Voer de migratie uit in een transactie
	err = m.db.Transaction(func(tx *gorm.DB) error {
		// Maak alle tabellen aan
		if err := tx.AutoMigrate(
			&models.ContactFormulier{},
			&models.ContactAntwoord{},
			&models.Aanmelding{},
			&models.AanmeldingAntwoord{},
			&models.EmailTemplate{},
			&models.VerzondEmail{},
			&models.Gebruiker{},
		); err != nil {
			return fmt.Errorf("fout bij aanmaken tabellen: %w", err)
		}

		// Registreer de migratie
		migratie := &models.Migratie{
			Versie:    "1.0.0",
			Naam:      "Initiële database setup",
			Toegepast: time.Now(),
		}

		if err := m.migrRepo.Create(ctx, migratie); err != nil {
			return fmt.Errorf("fout bij registreren migratie: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("fout bij uitvoeren migratie: %w", err)
	}

	logger.Info("Migratie 1.0.0 succesvol uitgevoerd")
	return nil
}

// SeedDatabase vult de database met initiële data
func (m *MigrationManager) SeedDatabase() error {
	logger.Info("Database seeding gestart")

	// Controleer of er al een admin gebruiker is
	var count int64
	if err := m.db.Model(&models.Gebruiker{}).Where("rol = ?", "admin").Count(&count).Error; err != nil {
		return fmt.Errorf("fout bij controleren admin gebruiker: %w", err)
	}

	if count > 0 {
		logger.Info("Admin gebruiker bestaat al, seeding overgeslagen")
		return nil
	}

	// Maak een admin gebruiker aan
	adminUser := &models.Gebruiker{
		Naam:           "Admin",
		Email:          "admin@dekoninklijkeloop.nl",
		WachtwoordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // password: admin
		Rol:            "admin",
		IsActief:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := m.db.Create(adminUser).Error; err != nil {
		return fmt.Errorf("fout bij aanmaken admin gebruiker: %w", err)
	}

	logger.Info("Admin gebruiker aangemaakt", "email", adminUser.Email)

	// Maak standaard email templates aan
	templates := []models.EmailTemplate{
		{
			Naam:         "contact_admin_email",
			Onderwerp:    "Nieuw contactformulier",
			Inhoud:       "<p>Er is een nieuw contactformulier ingevuld door {{.Contact.Naam}}.</p><p>Email: {{.Contact.Email}}</p><p>Bericht: {{.Contact.Bericht}}</p>",
			Beschrijving: "Email die naar de admin wordt gestuurd bij een nieuw contactformulier",
			IsActief:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			CreatedBy:    adminUser.ID,
		},
		{
			Naam:         "contact_email",
			Onderwerp:    "Bedankt voor je bericht",
			Inhoud:       "<p>Beste {{.Contact.Naam}},</p><p>Bedankt voor je bericht. We nemen zo snel mogelijk contact met je op.</p>",
			Beschrijving: "Bevestigingsemail die naar de gebruiker wordt gestuurd bij een contactformulier",
			IsActief:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			CreatedBy:    adminUser.ID,
		},
		{
			Naam:         "aanmelding_admin_email",
			Onderwerp:    "Nieuwe aanmelding ontvangen",
			Inhoud:       "<p>Er is een nieuwe aanmelding ontvangen van {{.Aanmelding.Naam}}.</p><p>Email: {{.Aanmelding.Email}}</p>",
			Beschrijving: "Email die naar de admin wordt gestuurd bij een nieuwe aanmelding",
			IsActief:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			CreatedBy:    adminUser.ID,
		},
		{
			Naam:         "aanmelding_email",
			Onderwerp:    "Bedankt voor je aanmelding",
			Inhoud:       "<p>Beste {{.Aanmelding.Naam}},</p><p>Bedankt voor je aanmelding. We hebben je aanmelding ontvangen en zullen deze zo snel mogelijk verwerken.</p>",
			Beschrijving: "Bevestigingsemail die naar de gebruiker wordt gestuurd bij een aanmelding",
			IsActief:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			CreatedBy:    adminUser.ID,
		},
	}

	for _, template := range templates {
		if err := m.db.Create(&template).Error; err != nil {
			return fmt.Errorf("fout bij aanmaken email template %s: %w", template.Naam, err)
		}
		logger.Info("Email template aangemaakt", "naam", template.Naam)
	}

	logger.Info("Database seeding voltooid")
	return nil
}
