package mocks

import (
	"context"
	"dklautomationgo/models"
	"errors"
	"sync"
	"time"
)

// MockDB is een in-memory database implementatie voor tests
type MockDB struct {
	mu                   sync.RWMutex
	contacts             map[string]*models.ContactFormulier
	contactAntwoorden    map[string]*models.ContactAntwoord
	aanmeldingen         map[string]*models.AanmeldingFormulier
	aanmeldingAntwoorden map[string]*models.AanmeldingAntwoord
	emailTemplates       map[string]*models.EmailTemplate
	verzondenEmails      map[string]*models.VerzondEmail
	gebruikers           map[string]*models.Gebruiker
	migraties            map[string]*models.Migratie
}

// NewMockDB maakt een nieuwe mock database
func NewMockDB() *MockDB {
	return &MockDB{
		contacts:             make(map[string]*models.ContactFormulier),
		contactAntwoorden:    make(map[string]*models.ContactAntwoord),
		aanmeldingen:         make(map[string]*models.AanmeldingFormulier),
		aanmeldingAntwoorden: make(map[string]*models.AanmeldingAntwoord),
		emailTemplates:       make(map[string]*models.EmailTemplate),
		verzondenEmails:      make(map[string]*models.VerzondEmail),
		gebruikers:           make(map[string]*models.Gebruiker),
		migraties:            make(map[string]*models.Migratie),
	}
}

// MockContactRepository is een mock implementatie van ContactRepository
type MockContactRepository struct {
	db *MockDB
}

// NewMockContactRepository maakt een nieuwe mock contact repository
func NewMockContactRepository(db *MockDB) *MockContactRepository {
	return &MockContactRepository{
		db: db,
	}
}

// Create slaat een nieuw contactformulier op
func (r *MockContactRepository) Create(ctx context.Context, contact *models.ContactFormulier) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if contact.ID == "" {
		return errors.New("contact ID is vereist")
	}

	r.db.contacts[contact.ID] = contact
	return nil
}

// GetByID haalt een contactformulier op basis van ID
func (r *MockContactRepository) GetByID(ctx context.Context, id string) (*models.ContactFormulier, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	contact, exists := r.db.contacts[id]
	if !exists {
		return nil, nil
	}

	return contact, nil
}

// List haalt een lijst van contactformulieren op
func (r *MockContactRepository) List(ctx context.Context, limit, offset int) ([]*models.ContactFormulier, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var result []*models.ContactFormulier
	for _, contact := range r.db.contacts {
		result = append(result, contact)
	}

	// Pas limit en offset toe
	if offset >= len(result) {
		return []*models.ContactFormulier{}, nil
	}

	end := offset + limit
	if end > len(result) {
		end = len(result)
	}

	return result[offset:end], nil
}

// Update werkt een contactformulier bij
func (r *MockContactRepository) Update(ctx context.Context, contact *models.ContactFormulier) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, exists := r.db.contacts[contact.ID]; !exists {
		return errors.New("contact niet gevonden")
	}

	contact.UpdatedAt = time.Now()
	r.db.contacts[contact.ID] = contact
	return nil
}

// Delete verwijdert een contactformulier
func (r *MockContactRepository) Delete(ctx context.Context, id string) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, exists := r.db.contacts[id]; !exists {
		return nil // Geen fout als het contact niet bestaat
	}

	delete(r.db.contacts, id)
	return nil
}

// FindByEmail zoekt contactformulieren op basis van email
func (r *MockContactRepository) FindByEmail(ctx context.Context, email string) ([]*models.ContactFormulier, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var result []*models.ContactFormulier
	for _, contact := range r.db.contacts {
		if contact.Email == email {
			result = append(result, contact)
		}
	}

	return result, nil
}

// FindByStatus zoekt contactformulieren op basis van status
func (r *MockContactRepository) FindByStatus(ctx context.Context, status string) ([]*models.ContactFormulier, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var result []*models.ContactFormulier
	for _, contact := range r.db.contacts {
		if contact.Status == status {
			result = append(result, contact)
		}
	}

	return result, nil
}

// MockMigratieRepository is een mock implementatie van MigratieRepository
type MockMigratieRepository struct {
	db *MockDB
}

// NewMockMigratieRepository maakt een nieuwe mock migratie repository
func NewMockMigratieRepository(db *MockDB) *MockMigratieRepository {
	return &MockMigratieRepository{
		db: db,
	}
}

// Create slaat een nieuwe migratie op
func (r *MockMigratieRepository) Create(ctx context.Context, migratie *models.Migratie) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if migratie.Versie == "" {
		return errors.New("migratie versie is vereist")
	}

	// Gebruik string versie als key in de map
	r.db.migraties[migratie.Versie] = migratie
	return nil
}

// GetByVersie haalt een migratie op basis van versie
func (r *MockMigratieRepository) GetByVersie(ctx context.Context, versie string) (*models.Migratie, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	migratie, exists := r.db.migraties[versie]
	if !exists {
		return nil, nil
	}

	return migratie, nil
}

// List haalt een lijst van migraties op
func (r *MockMigratieRepository) List(ctx context.Context) ([]*models.Migratie, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var result []*models.Migratie
	for _, migratie := range r.db.migraties {
		result = append(result, migratie)
	}

	return result, nil
}

// GetLatest haalt de meest recente migratie op
func (r *MockMigratieRepository) GetLatest(ctx context.Context) (*models.Migratie, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var latest *models.Migratie
	for _, migratie := range r.db.migraties {
		if latest == nil || migratie.Toegepast.After(latest.Toegepast) {
			latest = migratie
		}
	}

	return latest, nil
}
