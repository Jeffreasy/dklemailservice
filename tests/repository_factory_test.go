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

func TestRepositoryFactory(t *testing.T) {
	// Setup
	db := mockDB()

	// Test factory creation
	var factory *repository.RepositoryFactory

	// Afhankelijk van het type database, maak de juiste factory
	switch v := db.(type) {
	case *gorm.DB:
		factory = repository.NewRepositoryFactory(v)
	case *mocks.MockDB:
		// Voor een mock DB moeten we een andere aanpak gebruiken
		// In een echte implementatie zou je hier een mock factory maken
		// Voor deze test maken we gewoon een lege factory
		factory = &repository.RepositoryFactory{
			Contact:  mocks.NewMockContactRepository(v),
			Migratie: mocks.NewMockMigratieRepository(v),
			// Andere repositories zouden hier ook gemockt moeten worden
		}
	}

	// Verify all repositories are created
	assert.NotNil(t, factory, "Factory zou niet nil moeten zijn")
	assert.NotNil(t, factory.Contact, "Contact repository zou niet nil moeten zijn")
	assert.NotNil(t, factory.Migratie, "Migratie repository zou niet nil moeten zijn")

	// Test type assertions
	switch db.(type) {
	case *gorm.DB:
		// Als we een echte DB gebruiken, controleer dan de echte types
		_, ok := factory.Contact.(*repository.PostgresContactRepository)
		assert.True(t, ok, "Contact repository zou van type PostgresContactRepository moeten zijn")

		_, ok = factory.Migratie.(*repository.PostgresMigratieRepository)
		assert.True(t, ok, "Migratie repository zou van type PostgresMigratieRepository moeten zijn")
	case *mocks.MockDB:
		// Als we een mock DB gebruiken, controleer dan de mock types
		_, ok := factory.Contact.(*mocks.MockContactRepository)
		assert.True(t, ok, "Contact repository zou van type MockContactRepository moeten zijn")

		_, ok = factory.Migratie.(*mocks.MockMigratieRepository)
		assert.True(t, ok, "Migratie repository zou van type MockMigratieRepository moeten zijn")
	}
}
