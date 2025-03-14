// Package tests contains test utilities and mock implementations for the migration repository
package tests

import (
	"context"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/tests/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// testMigratieDB stelt een database op voor tests
func testMigratieDB() *mocks.MockDB {
	return mocks.NewMockDB()
}

// getMigratieRepository maakt een migratie repository op basis van de database
func getMigratieRepositoryForTest(db *mocks.MockDB) repository.MigratieRepository {
	return mocks.NewMockMigratieRepository(db)
}

// setupTestData voegt testdata toe aan de repository
func setupTestData(repo repository.MigratieRepository) {
	ctx := context.Background()

	// Voeg testmigraties toe
	repo.Create(ctx, &models.Migratie{
		Versie:    "001_initial_schema",
		Naam:      "Initial Schema",
		Toegepast: time.Now().Add(-2 * time.Hour),
	})

	repo.Create(ctx, &models.Migratie{
		Versie:    "002_seed_data",
		Naam:      "Seed Data",
		Toegepast: time.Now().Add(-1 * time.Hour),
	})
}

func TestPostgresMigratieRepository_Create(t *testing.T) {
	// Setup
	db := testMigratieDB()
	repo := getMigratieRepositoryForTest(db)
	ctx := context.Background()

	// Test data
	migratie := &models.Migratie{
		Versie:    "001_initial_schema",
		Naam:      "Initial Schema",
		Toegepast: time.Now(),
	}

	// Test
	err := repo.Create(ctx, migratie)
	assert.NoError(t, err, "Kon migratie niet aanmaken")

	// Verify
	saved, err := repo.GetByVersie(ctx, migratie.Versie)
	assert.NoError(t, err, "Kon migratie niet ophalen")
	assert.NotNil(t, saved, "Migratie niet gevonden")
	assert.Equal(t, migratie.Versie, saved.Versie, "Versies komen niet overeen")
	assert.Equal(t, migratie.Naam, saved.Naam, "Namen komen niet overeen")
}

func TestPostgresMigratieRepository_GetByVersie(t *testing.T) {
	// Setup
	db := testMigratieDB()
	repo := getMigratieRepositoryForTest(db)
	setupTestData(repo)
	ctx := context.Background()

	// Test: Bestaande versie
	t.Run("Bestaande versie", func(t *testing.T) {
		found, err := repo.GetByVersie(ctx, "001_initial_schema")
		assert.NoError(t, err, "Fout bij ophalen migratie")
		assert.NotNil(t, found, "Migratie niet gevonden")
		assert.Equal(t, "001_initial_schema", found.Versie, "Versies komen niet overeen")
	})

	// Test: Niet-bestaande versie
	t.Run("Niet-bestaande versie", func(t *testing.T) {
		found, err := repo.GetByVersie(ctx, "niet-bestaande-versie")
		assert.NoError(t, err, "Onverwachte fout bij ophalen niet-bestaande migratie")
		assert.Nil(t, found, "Niet-bestaande migratie gevonden")
	})
}

func TestPostgresMigratieRepository_List(t *testing.T) {
	// Setup
	db := testMigratieDB()
	repo := getMigratieRepositoryForTest(db)
	setupTestData(repo)
	ctx := context.Background()

	// Test: Lijst alle migraties
	migraties, err := repo.List(ctx)
	assert.NoError(t, err, "Fout bij ophalen migraties")
	assert.Len(t, migraties, 2, "Verkeerd aantal migraties opgehaald")

	// Controleer of alle migraties aanwezig zijn
	versieMap := make(map[string]bool)
	for _, m := range migraties {
		versieMap[m.Versie] = true
	}

	assert.True(t, versieMap["001_initial_schema"], "Migratie 001_initial_schema niet gevonden in lijst")
	assert.True(t, versieMap["002_seed_data"], "Migratie 002_seed_data niet gevonden in lijst")
}

func TestPostgresMigratieRepository_GetLatest(t *testing.T) {
	// Setup
	db := testMigratieDB()
	repo := getMigratieRepositoryForTest(db)
	setupTestData(repo)
	ctx := context.Background()

	// Test: Haal laatste migratie op
	latest, err := repo.GetLatest(ctx)
	assert.NoError(t, err, "Fout bij ophalen laatste migratie")
	assert.NotNil(t, latest, "Geen laatste migratie gevonden")

	// De laatste migratie zou de meest recente moeten zijn (002_seed_data)
	assert.Equal(t, "002_seed_data", latest.Versie, "Verkeerde laatste migratie")
}

func TestPostgresMigratieRepository_EmptyDatabase(t *testing.T) {
	// Setup
	db := testMigratieDB()
	repo := getMigratieRepositoryForTest(db)
	ctx := context.Background()

	// Test: Lijst in lege database
	t.Run("Lijst in lege database", func(t *testing.T) {
		migraties, err := repo.List(ctx)
		assert.NoError(t, err, "Fout bij ophalen migraties uit lege database")
		assert.Len(t, migraties, 0, "Migraties gevonden in lege database")
	})

	// Test: GetLatest in lege database
	t.Run("GetLatest in lege database", func(t *testing.T) {
		latest, err := repo.GetLatest(ctx)
		assert.NoError(t, err, "Fout bij ophalen laatste migratie uit lege database")
		assert.Nil(t, latest, "Migratie gevonden in lege database")
	})
}
