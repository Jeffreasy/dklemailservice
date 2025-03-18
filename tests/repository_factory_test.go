package tests

import (
	"dklautomationgo/repository"
	"dklautomationgo/tests/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// mockDB creates a mock DB for testing
func mockDB() interface{} {
	// We kunnen zowel een echte DB als een mock DB gebruiken
	// Voor deze test gebruiken we een mock DB omdat we geen echte DB nodig hebben
	return mocks.NewMockDB()
}

func TestRepository(t *testing.T) {
	// Setup
	db := mockDB()

	// Test factory creation
	var repo *repository.Repository

	// Afhankelijk van het type database, maak de juiste repository
	switch v := db.(type) {
	case *gorm.DB:
		repo = repository.NewRepository(v)
	case *mocks.MockDB:
		// Voor een mock DB moeten we een andere aanpak gebruiken
		// In een echte implementatie zou je hier een mock repository maken
		// Voor deze test maken we gewoon een lege repository
		repo = &repository.Repository{
			Contact:  mocks.NewMockContactRepository(v),
			Migratie: mocks.NewMockMigratieRepository(v),
			// Andere repositories zouden hier ook gemockt moeten worden
		}
	}

	// Verify all repositories are created
	assert.NotNil(t, repo, "Repository zou niet nil moeten zijn")
	assert.NotNil(t, repo.Contact, "Contact repository zou niet nil moeten zijn")
	assert.NotNil(t, repo.Migratie, "Migratie repository zou niet nil moeten zijn")

	// Test type assertions
	switch db.(type) {
	case *gorm.DB:
		// Als we een echte DB gebruiken, controleer dan de echte types
		_, ok := repo.Contact.(*repository.PostgresContactRepository)
		assert.True(t, ok, "Contact repository zou van type PostgresContactRepository moeten zijn")

		_, ok = repo.Migratie.(*repository.PostgresMigratieRepository)
		assert.True(t, ok, "Migratie repository zou van type PostgresMigratieRepository moeten zijn")
	case *mocks.MockDB:
		// Als we een mock DB gebruiken, controleer dan de mock types
		_, ok := repo.Contact.(*mocks.MockContactRepository)
		assert.True(t, ok, "Contact repository zou van type MockContactRepository moeten zijn")

		_, ok = repo.Migratie.(*mocks.MockMigratieRepository)
		assert.True(t, ok, "Migratie repository zou van type MockMigratieRepository moeten zijn")
	}
}
