package migrations

import (
	"dklautomationgo/logger"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"gorm.io/gorm"
)

//go:embed *.sql
var sqlMigrations embed.FS

// RunSQLMigrations voert alle SQL migratie scripts uit
func RunSQLMigrations(db *gorm.DB) error {
	logger.Info("SQL migraties worden uitgevoerd")

	// Lees alle SQL bestanden uit de embedded FS
	files, err := fs.ReadDir(sqlMigrations, ".")
	if err != nil {
		return fmt.Errorf("fout bij lezen embedded migrations: %w", err)
	}

	// Filter SQL bestanden en sorteer ze op naam
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	// Voer elke migratie uit
	for _, file := range sqlFiles {
		logger.Info("Migratie wordt uitgevoerd", "file", file)

		// Lees de inhoud van het bestand uit de embedded FS
		content, err := sqlMigrations.ReadFile(file)
		if err != nil {
			return fmt.Errorf("fout bij lezen migratie bestand %s: %w", file, err)
		}

		// Voer de SQL uit
		if err := db.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("fout bij uitvoeren migratie %s: %w", file, err)
		}

		logger.Info("Migratie succesvol uitgevoerd", "file", file)
	}

	logger.Info("Alle SQL migraties zijn succesvol uitgevoerd")
	return nil
}
