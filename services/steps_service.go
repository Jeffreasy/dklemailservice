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
}

// NewStepsService maakt een nieuwe steps service
func NewStepsService(db *gorm.DB, aanmeldingRepo repository.AanmeldingRepository) *StepsService {
	return &StepsService{
		db:             db,
		aanmeldingRepo: aanmeldingRepo,
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

// CalculateAllocatedFunds berekent toegewezen fondsen gebaseerd op afstand
func (s *StepsService) CalculateAllocatedFunds(route string) int {
	// Voorbeeld logica: verschillende bedragen per afstand
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
	// Totaal bedrag
	totalFunds := 10000

	// Routes
	routes := []string{"6 KM", "10 KM", "15 KM", "20 KM"}
	distribution := make(map[string]int)

	// Tel aantal deelnemers per route
	totalParticipants := 0
	for _, route := range routes {
		var count int64
		s.db.Model(&models.Aanmelding{}).Where("afstand = ?", route).Count(&count)
		totalParticipants += int(count)
		distribution[route] = int(count)
	}

	// Verdeel proportioneel
	if totalParticipants > 0 {
		for route, count := range distribution {
			distribution[route] = (totalFunds * count) / totalParticipants
		}
	}

	return distribution, totalFunds, nil
}
