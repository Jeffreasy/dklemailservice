package database

import (
	"dklautomationgo/database/migrations"
	"dklautomationgo/logger"
	"dklautomationgo/repository"
	"fmt"

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

	// [GEMINI] De AutoMigrate(&models.Migratie{}) call is verwijderd.
	// Dit was onderdeel van het oude, handmatige migratiesysteem.

	// Voer SQL migraties uit (V01 t/m V24)
	if err := migrations.RunSQLMigrations(m.db); err != nil {
		return fmt.Errorf("fout bij uitvoeren SQL migraties: %w", err)
	}

	// [GEMINI] De aanroep naar createTables() is verwijderd.
	// Deze functie bevatte de GORM AutoMigrate() die alle conflicten
	// (TIMESTAMP vs TIMESTAMPTZ, VARCHAR vs TEXT, en view locks) veroorzaakte.
	// Onze SQL-bestanden hebben dit werk al correct gedaan.

	logger.Info("Database migratie voltooid")
	return nil
}

// [GEMINI] De volledige 'createTables' functie is verwijderd.

// [GEMINI] De volledige 'SeedDatabase' functie is verwijderd.
// De logica hiervan zit nu in V02__seed_data.sql en V03__add_test_data.sql.
