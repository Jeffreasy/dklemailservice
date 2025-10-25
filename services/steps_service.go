package services

import (
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"fmt"

	"gorm.io/gorm"
)

// StepsService bevat business logic voor stappen tracking
type StepsService struct {
	db             *gorm.DB
	aanmeldingRepo repository.AanmeldingRepository
	routeFundRepo  repository.RouteFundRepository
}

// NewStepsService maakt een nieuwe steps service
func NewStepsService(db *gorm.DB, aanmeldingRepo repository.AanmeldingRepository, routeFundRepo repository.RouteFundRepository) *StepsService {
	return &StepsService{
		db:             db,
		aanmeldingRepo: aanmeldingRepo,
		routeFundRepo:  routeFundRepo,
	}
}

// UpdateSteps werkt stappen bij voor een deelnemer (delta toevoegen)
func (s *StepsService) UpdateSteps(participantID string, deltaSteps int) (*models.Aanmelding, error) {
	// Haal deelnemer op
	participant, err := s.aanmeldingRepo.GetByID(nil, participantID)
	if err != nil {
		return nil, fmt.Errorf("deelnemer niet gevonden: %w", err)
	}
	if participant == nil {
		return nil, fmt.Errorf("deelnemer niet gevonden")
	}

	// Update stappen (voorkom negatieve stappen)
	newSteps := participant.Steps + deltaSteps
	if newSteps < 0 {
		newSteps = 0
	}
	participant.Steps = newSteps

	// Sla wijzigingen op
	if err := s.aanmeldingRepo.Update(nil, participant); err != nil {
		return nil, fmt.Errorf("kon stappen niet bijwerken: %w", err)
	}

	return participant, nil
}

// UpdateStepsByUserID werkt stappen bij voor een deelnemer via gebruiker ID
func (s *StepsService) UpdateStepsByUserID(userID string, deltaSteps int) (*models.Aanmelding, error) {
	// Haal deelnemer op via gebruiker_id
	var participant models.Aanmelding
	err := s.db.Where("gebruiker_id = ?", userID).First(&participant).Error
	if err != nil {
		return nil, fmt.Errorf("deelnemer niet gevonden: %w", err)
	}

	// Update stappen (voorkom negatieve stappen)
	newSteps := participant.Steps + deltaSteps
	if newSteps < 0 {
		newSteps = 0
	}
	participant.Steps = newSteps

	// Sla wijzigingen op
	if err := s.aanmeldingRepo.Update(nil, &participant); err != nil {
		return nil, fmt.Errorf("kon stappen niet bijwerken: %w", err)
	}

	return &participant, nil
}

// GetParticipantDashboard haalt dashboard data op voor een deelnemer
func (s *StepsService) GetParticipantDashboard(participantID string) (*models.Aanmelding, int, error) {
	// Haal deelnemer op
	participant, err := s.aanmeldingRepo.GetByID(nil, participantID)
	if err != nil {
		return nil, 0, fmt.Errorf("deelnemer niet gevonden: %w", err)
	}
	if participant == nil {
		return nil, 0, fmt.Errorf("deelnemer niet gevonden")
	}

	// Bereken allocated funds gebaseerd op afstand
	allocatedFunds := s.CalculateAllocatedFunds(participant.Afstand)

	return participant, allocatedFunds, nil
}

// GetParticipantDashboardByUserID haalt dashboard data op voor een deelnemer via gebruiker ID
func (s *StepsService) GetParticipantDashboardByUserID(userID string) (*models.Aanmelding, int, error) {
	// Haal deelnemer op via gebruiker_id
	var participant models.Aanmelding
	err := s.db.Where("gebruiker_id = ?", userID).First(&participant).Error
	if err != nil {
		return nil, 0, fmt.Errorf("deelnemer niet gevonden: %w", err)
	}

	// Bereken allocated funds gebaseerd op afstand
	allocatedFunds := s.CalculateAllocatedFunds(participant.Afstand)

	return &participant, allocatedFunds, nil
}

// CalculateAllocatedFunds berekent toegewezen fondsen gebaseerd op afstand
func (s *StepsService) CalculateAllocatedFunds(route string) int {
	// Haal fondsallocatie op uit database
	routeFund, err := s.routeFundRepo.GetByRoute(nil, route)
	if err != nil {
		// Fallback naar standaard waarden als route niet gevonden wordt
		switch route {
		case "6 KM":
			return 50
		case "10 KM":
			return 75
		case "15 KM":
			return 100
		case "20 KM":
			return 125
		default:
			return 50 // Standaard bedrag
		}
	}
	return routeFund.Amount
}

// GetTotalSteps haalt totaal aantal stappen op voor een jaar
func (s *StepsService) GetTotalSteps(year int) (int, error) {
	var total int
	err := s.db.Model(&models.Aanmelding{}).
		Where("EXTRACT(YEAR FROM created_at) = ?", year).
		Select("SUM(steps)").Scan(&total).Error
	if err != nil {
		return 0, fmt.Errorf("kon totaal stappen niet ophalen: %w", err)
	}
	return total, nil
}

// GetFundsDistribution haalt fondsverdeling op over routes
func (s *StepsService) GetFundsDistribution() (map[string]int, error) {
	// Totaal bedrag (kan uit config komen)
	totalFunds := 10000

	// Routes
	routes := []string{"6 KM", "10 KM", "15 KM", "20 KM"}
	distribution := make(map[string]int)

	// Verdeel gelijk over routes (kan proportioneel gemaakt worden)
	fundsPerRoute := totalFunds / len(routes)

	for _, route := range routes {
		distribution[route] = fundsPerRoute
	}

	return distribution, nil
}

// GetFundsDistributionProportional haalt proportionele fondsverdeling op gebaseerd op aantal deelnemers
func (s *StepsService) GetFundsDistributionProportional() (map[string]int, int, error) {
	// Haal alle route funds op
	routeFunds, err := s.routeFundRepo.GetAll(nil)
	if err != nil {
		return nil, 0, fmt.Errorf("kon route funds niet ophalen: %w", err)
	}

	// Bereken totaal bedrag
	totalFunds := 0
	for _, rf := range routeFunds {
		totalFunds += rf.Amount
	}

	distribution := make(map[string]int)

	// Tel aantal deelnemers per route
	totalParticipants := 0
	for _, rf := range routeFunds {
		var count int64
		s.db.Model(&models.Aanmelding{}).Where("afstand = ?", rf.Route).Count(&count)
		totalParticipants += int(count)
		distribution[rf.Route] = int(count)
	}

	// Verdeel proportioneel gebaseerd op aantal deelnemers
	if totalParticipants > 0 {
		for route, count := range distribution {
			// Zoek het fondsbedrag voor deze route
			for _, rf := range routeFunds {
				if rf.Route == route {
					distribution[route] = (rf.Amount * count) / totalParticipants
					break
				}
			}
		}
	}

	return distribution, totalFunds, nil
}

// GetRouteFunds haalt alle route fondsallocaties op
func (s *StepsService) GetRouteFunds() ([]*models.RouteFund, error) {
	return s.routeFundRepo.GetAll(nil)
}

// UpdateRouteFund werkt een route fondsallocatie bij
func (s *StepsService) UpdateRouteFund(route string, amount int) (*models.RouteFund, error) {
	// Controleer of route bestaat
	existing, err := s.routeFundRepo.GetByRoute(nil, route)
	if err != nil {
		return nil, fmt.Errorf("route niet gevonden: %w", err)
	}

	existing.Amount = amount
	if err := s.routeFundRepo.Update(nil, existing); err != nil {
		return nil, fmt.Errorf("kon route fund niet bijwerken: %w", err)
	}

	return existing, nil
}

// CreateRouteFund maakt een nieuwe route fondsallocatie aan
func (s *StepsService) CreateRouteFund(route string, amount int) (*models.RouteFund, error) {
	routeFund := &models.RouteFund{
		Route:  route,
		Amount: amount,
	}

	if err := s.routeFundRepo.Create(nil, routeFund); err != nil {
		return nil, fmt.Errorf("kon route fund niet aanmaken: %w", err)
	}

	return routeFund, nil
}

// DeleteRouteFund verwijdert een route fondsallocatie
func (s *StepsService) DeleteRouteFund(route string) error {
	return s.routeFundRepo.Delete(nil, route)
}
