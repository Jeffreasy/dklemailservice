package mocks

import (
	"context"
	"dklautomationgo/models"
	"errors"
	"strings"
	"sync"
	"time"
)

// MockDB is een in-memory database implementatie voor tests
type MockDB struct {
	mu                    sync.RWMutex
	contacts              map[string]*models.ContactFormulier
	contactAntwoorden     map[string]*models.ContactAntwoord
	participants          map[string]*models.Participant
	participantAntwoorden map[string]*models.ParticipantAntwoord
	emailTemplates        map[string]*models.EmailTemplate
	verzondenEmails       map[string]*models.VerzondEmail
	gebruikers            map[string]*models.Gebruiker
	migraties             map[string]*models.Migratie
}

// NewMockDB maakt een nieuwe mock database
func NewMockDB() *MockDB {
	return &MockDB{
		contacts:              make(map[string]*models.ContactFormulier),
		contactAntwoorden:     make(map[string]*models.ContactAntwoord),
		participants:          make(map[string]*models.Participant),
		participantAntwoorden: make(map[string]*models.ParticipantAntwoord),
		emailTemplates:        make(map[string]*models.EmailTemplate),
		verzondenEmails:       make(map[string]*models.VerzondEmail),
		gebruikers:            make(map[string]*models.Gebruiker),
		migraties:             make(map[string]*models.Migratie),
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
		// V27: Direct field access (database column is 'status')
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

// MockContactAntwoordRepository is een mock implementatie van ContactAntwoordRepository
type MockContactAntwoordRepository struct {
	db *MockDB
}

// NewMockContactAntwoordRepository maakt een nieuwe mock contact antwoord repository
func NewMockContactAntwoordRepository(db *MockDB) *MockContactAntwoordRepository {
	return &MockContactAntwoordRepository{
		db: db,
	}
}

// Create slaat een nieuw contactantwoord op
func (r *MockContactAntwoordRepository) Create(ctx context.Context, antwoord *models.ContactAntwoord) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if antwoord.ID == "" {
		return errors.New("antwoord ID is vereist")
	}

	r.db.contactAntwoorden[antwoord.ID] = antwoord
	return nil
}

// GetByID haalt een contactantwoord op basis van ID
func (r *MockContactAntwoordRepository) GetByID(ctx context.Context, id string) (*models.ContactAntwoord, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	antwoord, exists := r.db.contactAntwoorden[id]
	if !exists {
		return nil, nil
	}

	return antwoord, nil
}

// ListByContactID haalt een lijst van contactantwoorden op basis van contactID
func (r *MockContactAntwoordRepository) ListByContactID(ctx context.Context, contactID string) ([]*models.ContactAntwoord, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var result []*models.ContactAntwoord
	for _, antwoord := range r.db.contactAntwoorden {
		if antwoord.ContactID == contactID {
			result = append(result, antwoord)
		}
	}

	return result, nil
}

// Update werkt een contactantwoord bij
func (r *MockContactAntwoordRepository) Update(ctx context.Context, antwoord *models.ContactAntwoord) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, exists := r.db.contactAntwoorden[antwoord.ID]; !exists {
		return errors.New("antwoord niet gevonden")
	}

	r.db.contactAntwoorden[antwoord.ID] = antwoord
	return nil
}

// Delete verwijdert een contactantwoord
func (r *MockContactAntwoordRepository) Delete(ctx context.Context, id string) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, exists := r.db.contactAntwoorden[id]; !exists {
		return nil // Geen fout als het antwoord niet bestaat
	}

	delete(r.db.contactAntwoorden, id)
	return nil
}

// MockParticipantRepository is een mock implementatie van ParticipantRepository
type MockParticipantRepository struct {
	db *MockDB
}

// NewMockParticipantRepository maakt een nieuwe mock participant repository
func NewMockParticipantRepository(db *MockDB) *MockParticipantRepository {
	return &MockParticipantRepository{
		db: db,
	}
}

// Create slaat een nieuwe participant op
func (r *MockParticipantRepository) Create(ctx context.Context, participant *models.Participant) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if participant.ID == "" {
		return errors.New("participant ID is vereist")
	}

	r.db.participants[participant.ID] = participant
	return nil
}

// GetByID haalt een participant op basis van ID
func (r *MockParticipantRepository) GetByID(ctx context.Context, id string) (*models.Participant, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	participant, exists := r.db.participants[id]
	if !exists {
		return nil, nil
	}

	return participant, nil
}

// List haalt een lijst van participants op
func (r *MockParticipantRepository) List(ctx context.Context, limit, offset int) ([]*models.Participant, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var result []*models.Participant
	for _, participant := range r.db.participants {
		result = append(result, participant)
	}

	// Pas limit en offset toe
	if offset >= len(result) {
		return []*models.Participant{}, nil
	}

	end := offset + limit
	if end > len(result) {
		end = len(result)
	}

	return result[offset:end], nil
}

// Update werkt een participant bij
func (r *MockParticipantRepository) Update(ctx context.Context, participant *models.Participant) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, exists := r.db.participants[participant.ID]; !exists {
		return errors.New("participant niet gevonden")
	}

	participant.UpdatedAt = time.Now()
	r.db.participants[participant.ID] = participant
	return nil
}

// Delete verwijdert een participant
func (r *MockParticipantRepository) Delete(ctx context.Context, id string) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, exists := r.db.participants[id]; !exists {
		return nil // Geen fout als de participant niet bestaat
	}

	delete(r.db.participants, id)
	return nil
}

// FindByEmail zoekt participants op basis van email
func (r *MockParticipantRepository) FindByEmail(ctx context.Context, email string) ([]*models.Participant, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var result []*models.Participant
	for _, participant := range r.db.participants {
		if participant.Email == email {
			result = append(result, participant)
		}
	}

	return result, nil
}

// FindByStatus zoekt participants op basis van status
// Status is niet meer deel van het Participant model, dus retourneer alle participants
func (r *MockParticipantRepository) FindByStatus(ctx context.Context, status string) ([]*models.Participant, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var result []*models.Participant
	for _, participant := range r.db.participants {
		result = append(result, participant)
	}

	return result, nil
}

// MockParticipantAntwoordRepository is een mock implementatie van ParticipantAntwoordRepository
type MockParticipantAntwoordRepository struct {
	db *MockDB
}

// NewMockParticipantAntwoordRepository maakt een nieuwe mock participant antwoord repository
func NewMockParticipantAntwoordRepository(db *MockDB) *MockParticipantAntwoordRepository {
	return &MockParticipantAntwoordRepository{
		db: db,
	}
}

// Create slaat een nieuw participantantwoord op
func (r *MockParticipantAntwoordRepository) Create(ctx context.Context, antwoord *models.ParticipantAntwoord) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if antwoord.ID == "" {
		return errors.New("antwoord ID is vereist")
	}

	r.db.participantAntwoorden[antwoord.ID] = antwoord
	return nil
}

// GetByID haalt een participantantwoord op basis van ID
func (r *MockParticipantAntwoordRepository) GetByID(ctx context.Context, id string) (*models.ParticipantAntwoord, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	antwoord, exists := r.db.participantAntwoorden[id]
	if !exists {
		return nil, nil
	}

	return antwoord, nil
}

// ListByParticipantID haalt een lijst van participantantwoorden op basis van participantID
func (r *MockParticipantAntwoordRepository) ListByParticipantID(ctx context.Context, participantID string) ([]*models.ParticipantAntwoord, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var result []*models.ParticipantAntwoord
	for _, antwoord := range r.db.participantAntwoorden {
		if antwoord.ParticipantID == participantID {
			result = append(result, antwoord)
		}
	}

	return result, nil
}

// Update werkt een participantantwoord bij
func (r *MockParticipantAntwoordRepository) Update(ctx context.Context, antwoord *models.ParticipantAntwoord) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, exists := r.db.participantAntwoorden[antwoord.ID]; !exists {
		return errors.New("antwoord niet gevonden")
	}

	r.db.participantAntwoorden[antwoord.ID] = antwoord
	return nil
}

// Delete verwijdert een participantantwoord
func (r *MockParticipantAntwoordRepository) Delete(ctx context.Context, id string) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, exists := r.db.participantAntwoorden[id]; !exists {
		return nil // Geen fout als het antwoord niet bestaat
	}

	delete(r.db.participantAntwoorden, id)
	return nil
}

// MockGebruikerRepository is een mock implementatie van GebruikerRepository
type MockGebruikerRepository struct {
	db *MockDB
}

// NewMockGebruikerRepository maakt een nieuwe mock gebruiker repository
func NewMockGebruikerRepository(db *MockDB) *MockGebruikerRepository {
	return &MockGebruikerRepository{
		db: db,
	}
}

// Create slaat een nieuwe gebruiker op
func (r *MockGebruikerRepository) Create(ctx context.Context, gebruiker *models.Gebruiker) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if gebruiker.ID == "" {
		return errors.New("gebruiker ID is vereist")
	}

	r.db.gebruikers[gebruiker.ID] = gebruiker
	return nil
}

// GetByID haalt een gebruiker op basis van ID
func (r *MockGebruikerRepository) GetByID(ctx context.Context, id string) (*models.Gebruiker, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	gebruiker, exists := r.db.gebruikers[id]
	if !exists {
		return nil, nil
	}

	return gebruiker, nil
}

// GetByEmail haalt een gebruiker op basis van email
func (r *MockGebruikerRepository) GetByEmail(ctx context.Context, email string) (*models.Gebruiker, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	for _, gebruiker := range r.db.gebruikers {
		if gebruiker.Email == email {
			return gebruiker, nil
		}
	}

	return nil, nil
}

// List haalt een lijst van gebruikers op
func (r *MockGebruikerRepository) List(ctx context.Context, limit, offset int) ([]*models.Gebruiker, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var result []*models.Gebruiker
	for _, gebruiker := range r.db.gebruikers {
		result = append(result, gebruiker)
	}

	// Pas limit en offset toe
	if offset >= len(result) {
		return []*models.Gebruiker{}, nil
	}

	end := offset + limit
	if end > len(result) {
		end = len(result)
	}

	return result[offset:end], nil
}

// Update werkt een gebruiker bij
func (r *MockGebruikerRepository) Update(ctx context.Context, gebruiker *models.Gebruiker) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, exists := r.db.gebruikers[gebruiker.ID]; !exists {
		return errors.New("gebruiker niet gevonden")
	}

	r.db.gebruikers[gebruiker.ID] = gebruiker
	return nil
}

// Delete verwijdert een gebruiker
func (r *MockGebruikerRepository) Delete(ctx context.Context, id string) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, exists := r.db.gebruikers[id]; !exists {
		return nil // Geen fout als de gebruiker niet bestaat
	}

	delete(r.db.gebruikers, id)
	return nil
}

// UpdateLastLogin werkt de laatste login tijd van een gebruiker bij
func (r *MockGebruikerRepository) UpdateLastLogin(ctx context.Context, id string) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	gebruiker, exists := r.db.gebruikers[id]
	if !exists {
		return errors.New("gebruiker niet gevonden")
	}

	now := time.Now()
	gebruiker.LaatsteLogin = &now
	r.db.gebruikers[id] = gebruiker
	return nil
}

// Search zoekt gebruikers op basis van naam of email
func (r *MockGebruikerRepository) Search(ctx context.Context, query string, limit int) ([]*models.Gebruiker, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var result []*models.Gebruiker
	for _, gebruiker := range r.db.gebruikers {
		// Simple search implementation - check if query matches name or email (case insensitive)
		queryLower := strings.ToLower(query)
		nameLower := strings.ToLower(gebruiker.Naam)
		emailLower := strings.ToLower(gebruiker.Email)

		if strings.Contains(nameLower, queryLower) || strings.Contains(emailLower, queryLower) {
			result = append(result, gebruiker)
			if limit > 0 && len(result) >= limit {
				break
			}
		}
	}

	return result, nil
}

// GetNewsletterSubscribers haalt actieve subscribers op
func (r *MockGebruikerRepository) GetNewsletterSubscribers(ctx context.Context) ([]*models.Gebruiker, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var result []*models.Gebruiker
	for _, gebruiker := range r.db.gebruikers {
		if gebruiker.IsActief && gebruiker.NewsletterSubscribed {
			result = append(result, gebruiker)
		}
	}

	return result, nil
}

// MockPermissionService is een mock implementatie van PermissionService
type MockPermissionService struct{}

// NewMockPermissionService maakt een nieuwe mock permission service
func NewMockPermissionService() *MockPermissionService {
	return &MockPermissionService{}
}

// HasPermission controleert of een gebruiker een specifieke permissie heeft (altijd true voor tests)
func (m *MockPermissionService) HasPermission(ctx context.Context, userID, resource, action string) bool {
	return true
}

// GetUserPermissions haalt alle permissies op voor een gebruiker
func (m *MockPermissionService) GetUserPermissions(ctx context.Context, userID string) ([]*models.UserPermission, error) {
	return []*models.UserPermission{}, nil
}

// GetUserRoles haalt alle actieve rollen op voor een gebruiker
func (m *MockPermissionService) GetUserRoles(ctx context.Context, userID string) ([]*models.UserRole, error) {
	return []*models.UserRole{}, nil
}

// AssignRole kent een rol toe aan een gebruiker
func (m *MockPermissionService) AssignRole(ctx context.Context, userID, roleID string, assignedBy *string) error {
	return nil
}

// RevokeRole verwijdert een rol van een gebruiker
func (m *MockPermissionService) RevokeRole(ctx context.Context, userID, roleID string) error {
	return nil
}

// CreateRole maakt een nieuwe rol aan
func (m *MockPermissionService) CreateRole(ctx context.Context, role *models.RBACRole, createdBy *string) error {
	return nil
}

// UpdateRole werkt een rol bij
func (m *MockPermissionService) UpdateRole(ctx context.Context, role *models.RBACRole) error {
	return nil
}

// DeleteRole verwijdert een rol
func (m *MockPermissionService) DeleteRole(ctx context.Context, roleID string) error {
	return nil
}

// AssignPermissionToRole kent een permissie toe aan een rol
func (m *MockPermissionService) AssignPermissionToRole(ctx context.Context, roleID, permissionID string, assignedBy *string) error {
	return nil
}

// RevokePermissionFromRole verwijdert een permissie van een rol
func (m *MockPermissionService) RevokePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	return nil
}

// GetRoles haalt alle rollen op
func (m *MockPermissionService) GetRoles(ctx context.Context, limit, offset int) ([]*models.RBACRole, error) {
	return []*models.RBACRole{}, nil
}

// GetPermissions haalt alle permissies op
func (m *MockPermissionService) GetPermissions(ctx context.Context, limit, offset int) ([]*models.Permission, error) {
	return []*models.Permission{}, nil
}

// InvalidateUserCache wist de cache voor een gebruiker
func (m *MockPermissionService) InvalidateUserCache(userID string) {
	// Do nothing in mock
}

// RefreshCache vernieuwt alle caches
func (m *MockPermissionService) RefreshCache(ctx context.Context) error {
	return nil
}

// MockEventRepository is een mock implementatie van EventRepository
type MockEventRepository struct {
	db *MockDB
}

// NewMockEventRepository maakt een nieuwe mock event repository
func NewMockEventRepository(db *MockDB) *MockEventRepository {
	return &MockEventRepository{
		db: db,
	}
}

// GetActiveEvent haalt het actieve event op
func (r *MockEventRepository) GetActiveEvent(ctx context.Context) (*models.Event, error) {
	// Return a mock active event for testing
	return &models.Event{
		ID:          "active-event-id",
		Name:        "Test Event",
		Description: "Test Event Description",
		StartTime:   time.Now().Add(24 * time.Hour),
		EndTime:     &[]time.Time{time.Now().Add(48 * time.Hour)}[0],
		Status:      "upcoming", // V27: Direct field (database column is 'status')
		IsActive:    true,
	}, nil
}

// Create slaat een nieuw event op
func (r *MockEventRepository) Create(ctx context.Context, event *models.Event) error {
	// Mock implementation - do nothing
	return nil
}

// GetByID haalt een event op basis van ID
func (r *MockEventRepository) GetByID(ctx context.Context, id string) (*models.Event, error) {
	// Return a mock event for testing
	return &models.Event{
		ID:          id,
		Name:        "Test Event",
		Description: "Test Event Description",
		StartTime:   time.Now().Add(24 * time.Hour),
		EndTime:     &[]time.Time{time.Now().Add(48 * time.Hour)}[0],
		Status:      "upcoming", // V27: Direct field (database column is 'status')
		IsActive:    true,
	}, nil
}

// List haalt een lijst van events op
func (r *MockEventRepository) List(ctx context.Context, limit, offset int) ([]*models.Event, error) {
	// Return empty list for testing
	return []*models.Event{}, nil
}

// Update werkt een event bij
func (r *MockEventRepository) Update(ctx context.Context, event *models.Event) error {
	// Mock implementation - do nothing
	return nil
}

// Delete verwijdert een event
func (r *MockEventRepository) Delete(ctx context.Context, id string) error {
	// Mock implementation - do nothing
	return nil
}

// GetEventParticipant haalt een event participant op
func (r *MockEventRepository) GetEventParticipant(ctx context.Context, eventID, participantID string) (*models.EventParticipant, error) {
	// Return a mock event participant for testing
	return &models.EventParticipant{
		ID:            "event-participant-id",
		EventID:       eventID,
		ParticipantID: participantID,
	}, nil
}

// GetEventParticipants haalt alle event participants op
func (r *MockEventRepository) GetEventParticipants(ctx context.Context, eventID string) ([]*models.EventParticipant, error) {
	// Return empty list for testing
	return []*models.EventParticipant{}, nil
}

// GetParticipantEvents haalt alle events voor een participant op
func (r *MockEventRepository) GetParticipantEvents(ctx context.Context, participantID string) ([]*models.EventParticipant, error) {
	// Return empty list for testing
	return []*models.EventParticipant{}, nil
}

// ListActive haalt alle actieve events op
func (r *MockEventRepository) ListActive(ctx context.Context) ([]*models.Event, error) {
	// Return empty list for testing
	return []*models.Event{}, nil
}

// RegisterParticipant registreert een participant voor een event
func (r *MockEventRepository) RegisterParticipant(ctx context.Context, eventID, participantID string) (*models.EventParticipant, error) {
	// Return a mock event participant for testing
	return &models.EventParticipant{
		ID:            "new-registration-id",
		EventID:       eventID,
		ParticipantID: participantID,
	}, nil
}

// UpdateEventParticipant werkt een event participant bij
func (r *MockEventRepository) UpdateEventParticipant(ctx context.Context, eventParticipant *models.EventParticipant) error {
	// Mock implementation - do nothing
	return nil
}

// MockEventRegistrationRepository is een mock implementatie van EventRegistrationRepository
type MockEventRegistrationRepository struct {
	db *MockDB
}

// NewMockEventRegistrationRepository maakt een nieuwe mock event registration repository
func NewMockEventRegistrationRepository(db *MockDB) *MockEventRegistrationRepository {
	return &MockEventRegistrationRepository{
		db: db,
	}
}

// Create slaat een nieuwe event registration op
func (r *MockEventRegistrationRepository) Create(ctx context.Context, registration *models.EventRegistration) error {
	// Mock implementation - do nothing
	return nil
}

// GetByID haalt een event registration op basis van ID
func (r *MockEventRegistrationRepository) GetByID(ctx context.Context, id string) (*models.EventRegistration, error) {
	// Return a mock registration for testing
	return &models.EventRegistration{
		ID:            id,
		EventID:       "event-id",
		ParticipantID: "participant-id",
		TestMode:      true,
	}, nil
}

// ListByEventID haalt een lijst van registrations op basis van event ID
func (r *MockEventRegistrationRepository) ListByEventID(ctx context.Context, eventID string) ([]*models.EventRegistration, error) {
	// Return empty list for testing
	return []*models.EventRegistration{}, nil
}

// ListByParticipantID haalt een lijst van registrations op basis van participant ID
func (r *MockEventRegistrationRepository) ListByParticipantID(ctx context.Context, participantID string) ([]*models.EventRegistration, error) {
	// Return empty list for testing
	return []*models.EventRegistration{}, nil
}

// Update werkt een event registration bij
func (r *MockEventRegistrationRepository) Update(ctx context.Context, registration *models.EventRegistration) error {
	// Mock implementation - do nothing
	return nil
}

// Delete verwijdert een event registration
func (r *MockEventRegistrationRepository) Delete(ctx context.Context, id string) error {
	// Mock implementation - do nothing
	return nil
}

// GetByEventAndParticipant haalt een registration op basis van event en participant ID
func (r *MockEventRegistrationRepository) GetByEventAndParticipant(ctx context.Context, eventID, participantID string) (*models.EventRegistration, error) {
	// Return a mock registration for testing
	return &models.EventRegistration{
		ID:            "registration-id",
		EventID:       eventID,
		ParticipantID: participantID,
		TestMode:      true,
	}, nil
}

// GetActiveRegistrationForParticipant haalt de actieve registration op voor een participant
func (r *MockEventRegistrationRepository) GetActiveRegistrationForParticipant(ctx context.Context, participantID string) (*models.EventRegistration, error) {
	// Return a mock active registration for testing
	return &models.EventRegistration{
		ID:            "active-registration-id",
		EventID:       "active-event-id",
		ParticipantID: participantID,
		TestMode:      true,
	}, nil
}

// ListByRole haalt een lijst van registrations op basis van rol
func (r *MockEventRegistrationRepository) ListByRole(ctx context.Context, roleID string) ([]*models.EventRegistration, error) {
	// Return empty list for testing
	return []*models.EventRegistration{}, nil
}

// List haalt een lijst van event registrations op
func (r *MockEventRegistrationRepository) List(ctx context.Context, limit, offset int) ([]*models.EventRegistration, error) {
	// Return empty list for testing
	return []*models.EventRegistration{}, nil
}

// ListByStatus haalt alle registraties met een specifieke status op
func (r *MockEventRegistrationRepository) ListByStatus(ctx context.Context, status string) ([]*models.EventRegistration, error) {
	// Return empty list for testing
	return []*models.EventRegistration{}, nil
}

// UpdateStatus werkt alleen de status van een registratie bij
func (r *MockEventRegistrationRepository) UpdateStatus(ctx context.Context, id, status string) error {
	// Mock implementation - do nothing
	return nil
}
