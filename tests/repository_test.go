package tests

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"dklautomationgo/tests/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// testDB stelt een database op voor tests
// Als CGO is ingeschakeld, gebruikt het SQLite, anders een in-memory mock
func testDB(t *testing.T) interface{} {
	// Probeer SQLite te gebruiken als CGO is ingeschakeld
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		// Als SQLite niet beschikbaar is, gebruik dan de mock database
		logger.Info("SQLite niet beschikbaar, gebruik mock database", "error", err)
		return mocks.NewMockDB()
	}

	// Voor SQLite: gebruik een custom UUID generator omdat gen_random_uuid() niet wordt ondersteund
	// We voeren handmatig de CREATE TABLE statements uit met SQLite-compatibele UUID defaults
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS contact_formulieren (
			id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))),
			created_at DATETIME,
			updated_at DATETIME,
			naam TEXT NOT NULL,
			email TEXT NOT NULL,
			bericht TEXT NOT NULL,
			email_verzonden BOOLEAN DEFAULT FALSE,
			email_verzonden_op DATETIME,
			privacy_akkoord BOOLEAN NOT NULL,
			status TEXT DEFAULT 'nieuw',
			behandeld_door TEXT,
			behandeld_op DATETIME,
			notities TEXT,
			beantwoord BOOLEAN DEFAULT FALSE,
			antwoord_tekst TEXT,
			antwoord_datum DATETIME,
			antwoord_door TEXT,
			test_mode BOOLEAN NOT NULL DEFAULT FALSE
		)
	`).Error
	require.NoError(t, err, "Kon contact_formulieren tabel niet aanmaken")

	return db
}

// getContactRepository maakt een contact repository op basis van de database
func getContactRepository(t *testing.T, db interface{}) repository.ContactRepository {
	switch v := db.(type) {
	case *gorm.DB:
		// Gebruik de echte PostgreSQL repository met SQLite
		baseRepo := repository.NewPostgresRepository(v)
		return repository.NewPostgresContactRepository(baseRepo)
	case *mocks.MockDB:
		// Gebruik de mock repository
		return mocks.NewMockContactRepository(v)
	default:
		t.Fatalf("Onbekend database type: %T", db)
		return nil
	}
}

// init initializes the test environment
func init() {
	// Zet de logger op error-only tijdens tests
	logger.Setup(logger.ErrorLevel)
}

func TestPostgresContactRepository_Create(t *testing.T) {
	// Setup
	db := testDB(t)
	repo := getContactRepository(t, db)
	ctx := context.Background()

	// Test data
	contact := &models.ContactFormulier{
		ID:             uuid.New().String(),
		Naam:           "Test Persoon",
		Email:          "test@example.com",
		Bericht:        "Dit is een testbericht",
		Status:         "nieuw",
		PrivacyAkkoord: true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Test
	err := repo.Create(ctx, contact)
	assert.NoError(t, err, "Kon contactformulier niet aanmaken")

	// Verify
	saved, err := repo.GetByID(ctx, contact.ID)
	assert.NoError(t, err, "Kon contactformulier niet ophalen")
	assert.NotNil(t, saved, "Contactformulier niet gevonden")
	assert.Equal(t, contact.Naam, saved.Naam, "Namen komen niet overeen")
	assert.Equal(t, contact.Email, saved.Email, "Emails komen niet overeen")
}

func TestPostgresContactRepository_GetByID(t *testing.T) {
	// Setup
	db := testDB(t)
	repo := getContactRepository(t, db)
	ctx := context.Background()

	// Test data
	contact := &models.ContactFormulier{
		ID:             uuid.New().String(),
		Naam:           "Test Persoon",
		Email:          "test@example.com",
		Bericht:        "Dit is een testbericht",
		Status:         "nieuw",
		PrivacyAkkoord: true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Voeg toe aan database
	err := repo.Create(ctx, contact)
	require.NoError(t, err, "Kon contactformulier niet aanmaken")

	// Test: Bestaand ID
	t.Run("Bestaand ID", func(t *testing.T) {
		found, err := repo.GetByID(ctx, contact.ID)
		assert.NoError(t, err, "Fout bij ophalen contactformulier")
		assert.NotNil(t, found, "Contactformulier niet gevonden")
		assert.Equal(t, contact.ID, found.ID, "IDs komen niet overeen")
	})

	// Test: Niet-bestaand ID
	t.Run("Niet-bestaand ID", func(t *testing.T) {
		found, err := repo.GetByID(ctx, "niet-bestaand-id")
		assert.NoError(t, err, "Onverwachte fout bij ophalen niet-bestaand contactformulier")
		assert.Nil(t, found, "Niet-bestaand contactformulier gevonden")
	})
}

func TestPostgresContactRepository_List(t *testing.T) {
	// Setup
	db := testDB(t)
	repo := getContactRepository(t, db)
	ctx := context.Background()

	// Voor SQLite: schoon de tabel eerst om conflicten met andere tests te voorkomen
	if gormDB, ok := db.(*gorm.DB); ok {
		gormDB.Exec("DELETE FROM contact_formulieren")
	}

	// Voeg meerdere contactformulieren toe
	for i := 0; i < 5; i++ {
		contact := &models.ContactFormulier{
			ID:             uuid.New().String(),
			Naam:           "Test Persoon",
			Email:          "test@example.com",
			Bericht:        "Dit is een testbericht",
			Status:         "nieuw",
			PrivacyAkkoord: true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		err := repo.Create(ctx, contact)
		require.NoError(t, err, "Kon contactformulier niet aanmaken")
	}

	// Test: Lijst met limiet
	t.Run("Lijst met limiet", func(t *testing.T) {
		contacts, err := repo.List(ctx, 3, 0)
		assert.NoError(t, err, "Fout bij ophalen contactformulieren")
		assert.Len(t, contacts, 3, "Verkeerd aantal contactformulieren opgehaald")
	})

	// Test: Lijst met offset
	t.Run("Lijst met offset", func(t *testing.T) {
		contacts, err := repo.List(ctx, 5, 2)
		assert.NoError(t, err, "Fout bij ophalen contactformulieren met offset")
		assert.Len(t, contacts, 3, "Verkeerd aantal contactformulieren opgehaald met offset")
	})
}

func TestPostgresContactRepository_Update(t *testing.T) {
	// Setup
	db := testDB(t)
	repo := getContactRepository(t, db)
	ctx := context.Background()

	// Test data
	contact := &models.ContactFormulier{
		ID:             uuid.New().String(),
		Naam:           "Test Persoon",
		Email:          "test@example.com",
		Bericht:        "Dit is een testbericht",
		Status:         "nieuw",
		PrivacyAkkoord: true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Voeg toe aan database
	err := repo.Create(ctx, contact)
	require.NoError(t, err, "Kon contactformulier niet aanmaken")

	// Update contactformulier
	contact.Naam = "Gewijzigde Naam"
	contact.Status = "in behandeling"
	err = repo.Update(ctx, contact)
	assert.NoError(t, err, "Kon contactformulier niet bijwerken")

	// Verify
	updated, err := repo.GetByID(ctx, contact.ID)
	assert.NoError(t, err, "Kon bijgewerkt contactformulier niet ophalen")
	assert.Equal(t, "Gewijzigde Naam", updated.Naam, "Naam niet bijgewerkt")
	assert.Equal(t, "in behandeling", updated.Status, "Status niet bijgewerkt")
}

func TestPostgresContactRepository_Delete(t *testing.T) {
	// Setup
	db := testDB(t)
	repo := getContactRepository(t, db)
	ctx := context.Background()

	// Test data
	contact := &models.ContactFormulier{
		ID:             uuid.New().String(),
		Naam:           "Test Persoon",
		Email:          "test@example.com",
		Bericht:        "Dit is een testbericht",
		Status:         "nieuw",
		PrivacyAkkoord: true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Voeg toe aan database
	err := repo.Create(ctx, contact)
	require.NoError(t, err, "Kon contactformulier niet aanmaken")

	// Verwijder contactformulier
	err = repo.Delete(ctx, contact.ID)
	assert.NoError(t, err, "Kon contactformulier niet verwijderen")

	// Verify
	deleted, err := repo.GetByID(ctx, contact.ID)
	assert.NoError(t, err, "Onverwachte fout bij ophalen verwijderd contactformulier")
	assert.Nil(t, deleted, "Verwijderd contactformulier nog steeds gevonden")
}

func TestPostgresContactRepository_FindByEmail(t *testing.T) {
	// Setup
	db := testDB(t)
	repo := getContactRepository(t, db)
	ctx := context.Background()

	// Voeg contactformulieren toe met verschillende emails
	emails := []string{"test1@example.com", "test2@example.com", "test1@example.com"}
	for i, email := range emails {
		contact := &models.ContactFormulier{
			ID:             uuid.New().String(),
			Naam:           "Test Persoon",
			Email:          email,
			Bericht:        "Dit is een testbericht",
			Status:         "nieuw",
			PrivacyAkkoord: true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		err := repo.Create(ctx, contact)
		require.NoError(t, err, "Kon contactformulier %d niet aanmaken", i)
	}

	// Test: Zoek op bestaande email
	t.Run("Zoek op bestaande email", func(t *testing.T) {
		contacts, err := repo.FindByEmail(ctx, "test1@example.com")
		assert.NoError(t, err, "Fout bij zoeken op email")
		assert.Len(t, contacts, 2, "Verkeerd aantal contactformulieren gevonden")
	})

	// Test: Zoek op niet-bestaande email
	t.Run("Zoek op niet-bestaande email", func(t *testing.T) {
		contacts, err := repo.FindByEmail(ctx, "niet-bestaand@example.com")
		assert.NoError(t, err, "Fout bij zoeken op niet-bestaande email")
		assert.Len(t, contacts, 0, "Contactformulieren gevonden voor niet-bestaande email")
	})
}

func TestPostgresContactRepository_FindByStatus(t *testing.T) {
	// Setup
	db := testDB(t)
	repo := getContactRepository(t, db)
	ctx := context.Background()

	// Voor SQLite: schoon de tabel eerst om conflicten met andere tests te voorkomen
	if gormDB, ok := db.(*gorm.DB); ok {
		gormDB.Exec("DELETE FROM contact_formulieren")
	}

	// Voeg contactformulieren toe met verschillende statussen
	statussen := []string{"nieuw", "in behandeling", "afgehandeld", "nieuw"}
	for i, status := range statussen {
		contact := &models.ContactFormulier{
			ID:             uuid.New().String(),
			Naam:           "Test Persoon",
			Email:          "test@example.com",
			Bericht:        "Dit is een testbericht",
			Status:         status,
			PrivacyAkkoord: true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		err := repo.Create(ctx, contact)
		require.NoError(t, err, "Kon contactformulier %d niet aanmaken", i)
	}

	// Test: Zoek op bestaande status
	t.Run("Zoek op bestaande status", func(t *testing.T) {
		contacts, err := repo.FindByStatus(ctx, "nieuw")
		assert.NoError(t, err, "Fout bij zoeken op status")
		assert.Len(t, contacts, 2, "Verkeerd aantal contactformulieren gevonden")
	})

	// Test: Zoek op niet-bestaande status
	t.Run("Zoek op niet-bestaande status", func(t *testing.T) {
		contacts, err := repo.FindByStatus(ctx, "niet-bestaand")
		assert.NoError(t, err, "Fout bij zoeken op niet-bestaande status")
		assert.Len(t, contacts, 0, "Contactformulieren gevonden voor niet-bestaande status")
	})
}
