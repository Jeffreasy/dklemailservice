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
	aanmeldingen         map[string]*models.Aanmelding
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
		aanmeldingen:         make(map[string]*models.Aanmelding),
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

// MockAanmeldingRepository is een mock implementatie van AanmeldingRepository
type MockAanmeldingRepository struct {
	db *MockDB
}

// NewMockAanmeldingRepository maakt een nieuwe mock aanmelding repository
func NewMockAanmeldingRepository(db *MockDB) *MockAanmeldingRepository {
	return &MockAanmeldingRepository{
		db: db,
	}
}

// Create slaat een nieuwe aanmelding op
func (r *MockAanmeldingRepository) Create(ctx context.Context, aanmelding *models.Aanmelding) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if aanmelding.ID == "" {
		return errors.New("aanmelding ID is vereist")
	}

	r.db.aanmeldingen[aanmelding.ID] = aanmelding
	return nil
}

// GetByID haalt een aanmelding op basis van ID
func (r *MockAanmeldingRepository) GetByID(ctx context.Context, id string) (*models.Aanmelding, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	aanmelding, exists := r.db.aanmeldingen[id]
	if !exists {
		return nil, nil
	}

	return aanmelding, nil
}

// List haalt een lijst van aanmeldingen op
func (r *MockAanmeldingRepository) List(ctx context.Context, limit, offset int) ([]*models.Aanmelding, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var result []*models.Aanmelding
	for _, aanmelding := range r.db.aanmeldingen {
		result = append(result, aanmelding)
	}

	// Pas limit en offset toe
	if offset >= len(result) {
		return []*models.Aanmelding{}, nil
	}

	end := offset + limit
	if end > len(result) {
		end = len(result)
	}

	return result[offset:end], nil
}

// Update werkt een aanmelding bij
func (r *MockAanmeldingRepository) Update(ctx context.Context, aanmelding *models.Aanmelding) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, exists := r.db.aanmeldingen[aanmelding.ID]; !exists {
		return errors.New("aanmelding niet gevonden")
	}

	aanmelding.UpdatedAt = time.Now()
	r.db.aanmeldingen[aanmelding.ID] = aanmelding
	return nil
}

// Delete verwijdert een aanmelding
func (r *MockAanmeldingRepository) Delete(ctx context.Context, id string) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, exists := r.db.aanmeldingen[id]; !exists {
		return nil // Geen fout als de aanmelding niet bestaat
	}

	delete(r.db.aanmeldingen, id)
	return nil
}

// FindByEmail zoekt aanmeldingen op basis van email
func (r *MockAanmeldingRepository) FindByEmail(ctx context.Context, email string) ([]*models.Aanmelding, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var result []*models.Aanmelding
	for _, aanmelding := range r.db.aanmeldingen {
		if aanmelding.Email == email {
			result = append(result, aanmelding)
		}
	}

	return result, nil
}

// FindByStatus zoekt aanmeldingen op basis van status
func (r *MockAanmeldingRepository) FindByStatus(ctx context.Context, status string) ([]*models.Aanmelding, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	// Bekende statussen
	knownStatuses := map[string]bool{
		"nieuw":          true,
		"in_behandeling": true,
		"beantwoord":     true,
		"gesloten":       true,
	}

	// Als de status niet bekend is, retourneer alle aanmeldingen
	// Dit is nodig voor de GetAanmeldingenByRol methode die FindByStatus gebruikt met rol als parameter
	if !knownStatuses[status] {
		var result []*models.Aanmelding
		for _, aanmelding := range r.db.aanmeldingen {
			result = append(result, aanmelding)
		}
		return result, nil
	}

	// Anders filter op status
	var result []*models.Aanmelding
	for _, aanmelding := range r.db.aanmeldingen {
		if aanmelding.Status == status {
			result = append(result, aanmelding)
		}
	}

	return result, nil
}

// MockAanmeldingAntwoordRepository is een mock implementatie van AanmeldingAntwoordRepository
type MockAanmeldingAntwoordRepository struct {
	db *MockDB
}

// NewMockAanmeldingAntwoordRepository maakt een nieuwe mock aanmelding antwoord repository
func NewMockAanmeldingAntwoordRepository(db *MockDB) *MockAanmeldingAntwoordRepository {
	return &MockAanmeldingAntwoordRepository{
		db: db,
	}
}

// Create slaat een nieuw aanmeldingantwoord op
func (r *MockAanmeldingAntwoordRepository) Create(ctx context.Context, antwoord *models.AanmeldingAntwoord) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if antwoord.ID == "" {
		return errors.New("antwoord ID is vereist")
	}

	r.db.aanmeldingAntwoorden[antwoord.ID] = antwoord
	return nil
}

// GetByID haalt een aanmeldingantwoord op basis van ID
func (r *MockAanmeldingAntwoordRepository) GetByID(ctx context.Context, id string) (*models.AanmeldingAntwoord, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	antwoord, exists := r.db.aanmeldingAntwoorden[id]
	if !exists {
		return nil, nil
	}

	return antwoord, nil
}

// ListByAanmeldingID haalt een lijst van aanmeldingantwoorden op basis van aanmeldingID
func (r *MockAanmeldingAntwoordRepository) ListByAanmeldingID(ctx context.Context, aanmeldingID string) ([]*models.AanmeldingAntwoord, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	var result []*models.AanmeldingAntwoord
	for _, antwoord := range r.db.aanmeldingAntwoorden {
		if antwoord.AanmeldingID == aanmeldingID {
			result = append(result, antwoord)
		}
	}

	return result, nil
}

// Update werkt een aanmeldingantwoord bij
func (r *MockAanmeldingAntwoordRepository) Update(ctx context.Context, antwoord *models.AanmeldingAntwoord) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, exists := r.db.aanmeldingAntwoorden[antwoord.ID]; !exists {
		return errors.New("antwoord niet gevonden")
	}

	r.db.aanmeldingAntwoorden[antwoord.ID] = antwoord
	return nil
}

// Delete verwijdert een aanmeldingantwoord
func (r *MockAanmeldingAntwoordRepository) Delete(ctx context.Context, id string) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, exists := r.db.aanmeldingAntwoorden[id]; !exists {
		return nil // Geen fout als het antwoord niet bestaat
	}

	delete(r.db.aanmeldingAntwoorden, id)
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
