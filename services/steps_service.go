package services

import (
	"context"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// StepsService bevat business logic voor stappen tracking
type StepsService struct {
	db              *gorm.DB
	participantRepo repository.ParticipantRepository
	distanceRepo    repository.DistanceRepository
	stepsHub        *StepsHub // WebSocket hub voor real-time updates
}

// NewStepsService maakt een nieuwe steps service
func NewStepsService(db *gorm.DB, participantRepo repository.ParticipantRepository, distanceRepo repository.DistanceRepository) *StepsService {
	return &StepsService{
		db:              db,
		participantRepo: participantRepo,
		distanceRepo:    distanceRepo,
		stepsHub:        nil, // Wordt later gezet via SetStepsHub
	}
}

// SetStepsHub sets de StepsHub voor real-time updates
func (s *StepsService) SetStepsHub(hub *StepsHub) {
	s.stepsHub = hub
}

// GetDB returns the database connection
func (s *StepsService) GetDB() *gorm.DB {
	return s.db
}

// UpdateSteps werkt stappen bij voor een deelnemer (delta toevoegen)
func (s *StepsService) UpdateSteps(participantID string, deltaSteps int) (*models.Participant, error) {
	ctx := context.Background()

	// Haal deelnemer op
	participant, err := s.participantRepo.GetByID(ctx, participantID)
	if err != nil {
		return nil, fmt.Errorf("deelnemer niet gevonden: %w", err)
	}
	if participant == nil {
		return nil, fmt.Errorf("deelnemer niet gevonden")
	}

	// Haal event registration op voor stappen bijwerken
	var eventReg models.EventRegistration
	if err := s.db.Where("participant_id = ?", participantID).First(&eventReg).Error; err != nil {
		return nil, fmt.Errorf("event registration niet gevonden: %w", err)
	}

	// Update stappen (voorkom negatieve stappen)
	newSteps := eventReg.Steps + deltaSteps
	if newSteps < 0 {
		newSteps = 0
	}
	eventReg.Steps = newSteps

	// Sla event registration op
	if err := s.db.Save(&eventReg).Error; err != nil {
		return nil, fmt.Errorf("kon stappen niet bijwerken: %w", err)
	}

	// ✨ Broadcast WebSocket update
	s.broadcastStepUpdate(participant, deltaSteps)

	return participant, nil
}

// UpdateStepsByUserID werkt stappen bij voor een deelnemer via gebruiker ID
func (s *StepsService) UpdateStepsByUserID(userID string, deltaSteps int) (*models.Participant, error) {
	// Haal deelnemer op via gebruiker_id
	var participant models.Participant
	err := s.db.Where("gebruiker_id = ?", userID).First(&participant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("geen deelnemersregistratie gevonden voor gebruiker %s - gebruiker is mogelijk geen deelnemer", userID)
		}
		return nil, fmt.Errorf("fout bij ophalen deelnemer: %w", err)
	}

	// Haal event registration op voor stappen bijwerken
	var eventReg models.EventRegistration
	if err := s.db.Where("participant_id = ?", userID).First(&eventReg).Error; err != nil {
		return nil, fmt.Errorf("event registration niet gevonden: %w", err)
	}

	// Update stappen (voorkom negatieve stappen)
	newSteps := eventReg.Steps + deltaSteps
	if newSteps < 0 {
		newSteps = 0
	}
	eventReg.Steps = newSteps

	// Sla event registration op
	if err := s.db.Save(&eventReg).Error; err != nil {
		return nil, fmt.Errorf("kon stappen niet bijwerken: %w", err)
	}

	// Sla wijzigingen op
	ctx := context.Background()
	if err := s.participantRepo.Update(ctx, &participant); err != nil {
		return nil, fmt.Errorf("kon stappen niet bijwerken: %w", err)
	}

	// ✨ Broadcast WebSocket update
	s.broadcastStepUpdate(&participant, deltaSteps)

	return &participant, nil
}

// GetParticipantDashboard haalt dashboard data op voor een deelnemer
func (s *StepsService) GetParticipantDashboard(participantID string) (*models.Participant, int, error) {
	ctx := context.Background()

	// Haal deelnemer op
	participant, err := s.participantRepo.GetByID(ctx, participantID)
	if err != nil {
		return nil, 0, fmt.Errorf("deelnemer niet gevonden: %w", err)
	}
	if participant == nil {
		return nil, 0, fmt.Errorf("deelnemer niet gevonden")
	}

	// Haal event registration op voor afstand
	var eventReg models.EventRegistration
	if err := s.db.Where("participant_id = ?", participantID).First(&eventReg).Error; err != nil {
		return nil, 0, fmt.Errorf("event registration niet gevonden: %w", err)
	}

	// Bereken allocated funds gebaseerd op afstand
	route := ""
	if eventReg.DistanceRoute != nil {
		route = *eventReg.DistanceRoute
	}
	allocatedFunds := s.CalculateAllocatedFunds(route)

	return participant, allocatedFunds, nil
}

// GetParticipantDashboardByUserID haalt dashboard data op voor een deelnemer via gebruiker ID
func (s *StepsService) GetParticipantDashboardByUserID(userID string) (*models.Participant, int, error) {
	// Haal deelnemer op via gebruiker_id
	var participant models.Participant
	err := s.db.Where("gebruiker_id = ?", userID).First(&participant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, fmt.Errorf("geen deelnemersregistratie gevonden voor gebruiker %s - gebruiker is mogelijk geen deelnemer", userID)
		}
		return nil, 0, fmt.Errorf("fout bij ophalen deelnemer: %w", err)
	}

	// Haal event registration op voor afstand
	var eventReg models.EventRegistration
	if err := s.db.Where("participant_id = ?", userID).First(&eventReg).Error; err != nil {
		return nil, 0, fmt.Errorf("event registration niet gevonden: %w", err)
	}

	// Bereken allocated funds gebaseerd op afstand
	route := ""
	if eventReg.DistanceRoute != nil {
		route = *eventReg.DistanceRoute
	}
	allocatedFunds := s.CalculateAllocatedFunds(route)

	return &participant, allocatedFunds, nil
}

// CalculateAllocatedFunds berekent toegewezen fondsen gebaseerd op afstand
func (s *StepsService) CalculateAllocatedFunds(route string) int {
	ctx := context.Background()

	// Haal fondsallocatie op uit database
	distance, err := s.distanceRepo.GetByRoute(ctx, route)
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
	return distance.RegistrationFee
}

// GetTotalSteps haalt totaal aantal stappen op
// Als year = 0, tel dan alle stappen op (ongeacht jaar van aanmelding)
// Anders filter op jaar van aanmelding
func (s *StepsService) GetTotalSteps(year int) (int, error) {
	var total int
	query := s.db.Model(&models.Participant{})

	// Als year > 0, filter op jaar van aanmelding
	if year > 0 {
		query = query.Where("EXTRACT(YEAR FROM created_at) = ?", year)
	}
	// Anders: tel ALLE stappen op van ALLE deelnemers

	err := query.Select("COALESCE(SUM(steps), 0)").Scan(&total).Error
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
	ctx := context.Background()

	// Haal alle distances op
	distances, err := s.distanceRepo.GetAll(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("kon distances niet ophalen: %w", err)
	}

	// Bereken totaal bedrag
	totalFunds := 0
	for _, d := range distances {
		totalFunds += d.RegistrationFee
	}

	distribution := make(map[string]int)

	// Tel aantal deelnemers per route
	totalParticipants := 0
	for _, d := range distances {
		var count int64
		s.db.Model(&models.Participant{}).Where("afstand = ?", d.Route).Count(&count)
		totalParticipants += int(count)
		distribution[d.Route] = int(count)
	}

	// Verdeel proportioneel gebaseerd op aantal deelnemers
	if totalParticipants > 0 {
		for route, count := range distribution {
			// Zoek het fondsbedrag voor deze route
			for _, d := range distances {
				if d.Route == route {
					distribution[route] = (d.RegistrationFee * count) / totalParticipants
					break
				}
			}
		}
	}

	return distribution, totalFunds, nil
}

// GetRouteFunds haalt alle route fondsallocaties op
func (s *StepsService) GetRouteFunds() ([]*models.RouteFund, error) {
	ctx := context.Background()
	distances, err := s.distanceRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Convert Distance models to RouteFund models for backward compatibility
	routeFunds := make([]*models.RouteFund, len(distances))
	for i, d := range distances {
		routeFunds[i] = &models.RouteFund{
			Route:  d.Route,
			Amount: d.RegistrationFee,
		}
	}
	return routeFunds, nil
}

// UpdateRouteFund werkt een route fondsallocatie bij
func (s *StepsService) UpdateRouteFund(route string, amount int) (*models.RouteFund, error) {
	ctx := context.Background()

	// Controleer of route bestaat
	existing, err := s.distanceRepo.GetByRoute(ctx, route)
	if err != nil {
		return nil, fmt.Errorf("route niet gevonden: %w", err)
	}

	existing.RegistrationFee = amount
	if err := s.distanceRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("kon distance niet bijwerken: %w", err)
	}

	// Return RouteFund for backward compatibility
	return &models.RouteFund{
		Route:  existing.Route,
		Amount: existing.RegistrationFee,
	}, nil
}

// CreateRouteFund maakt een nieuwe route fondsallocatie aan
func (s *StepsService) CreateRouteFund(route string, amount int) (*models.RouteFund, error) {
	distance := &models.Distance{
		Route:           route,
		RegistrationFee: amount,
	}

	if err := s.distanceRepo.Create(context.TODO(), distance); err != nil {
		return nil, fmt.Errorf("kon distance niet aanmaken: %w", err)
	}

	// Return RouteFund for backward compatibility
	return &models.RouteFund{
		Route:  distance.Route,
		Amount: distance.RegistrationFee,
	}, nil
}

// DeleteRouteFund verwijdert een route fondsallocatie
func (s *StepsService) DeleteRouteFund(route string) error {
	ctx := context.Background()
	return s.distanceRepo.Delete(ctx, route)
}

// broadcastStepUpdate broadcast een stappen update via WebSocket
func (s *StepsService) broadcastStepUpdate(participant *models.Participant, delta int) {
	if s.stepsHub == nil {
		return // WebSocket niet geïnitialiseerd
	}

	// Haal event registration op voor afstand
	var eventReg models.EventRegistration
	if err := s.db.Where("participant_id = ?", participant.ID).First(&eventReg).Error; err != nil {
		return
	}

	route := ""
	if eventReg.DistanceRoute != nil {
		route = *eventReg.DistanceRoute
	}
	allocatedFunds := s.CalculateAllocatedFunds(route)

	// Broadcast step update
	s.stepsHub.StepUpdate <- &StepUpdateMessage{
		Type:           MessageTypeStepUpdate,
		ParticipantID:  participant.ID,
		Naam:           participant.Naam,
		Steps:          eventReg.Steps,
		Delta:          delta,
		Route:          route,
		AllocatedFunds: allocatedFunds,
		Timestamp:      time.Now().Unix(),
	}

	// Update totaal stappen in background
	go s.broadcastTotalSteps()

	// Update leaderboard in background
	go s.broadcastLeaderboard()
}

// broadcastTotalSteps broadcast totaal stappen update
func (s *StepsService) broadcastTotalSteps() {
	if s.stepsHub == nil {
		return
	}

	totalSteps, err := s.GetTotalSteps(0)
	if err != nil {
		return
	}

	s.stepsHub.TotalUpdate <- &TotalUpdateMessage{
		Type:       MessageTypeTotalUpdate,
		TotalSteps: totalSteps,
		Year:       0,
		Timestamp:  time.Now().Unix(),
	}
}

// broadcastLeaderboard broadcast leaderboard update (top 10)
func (s *StepsService) broadcastLeaderboard() {
	if s.stepsHub == nil {
		return
	}

	var entries []LeaderboardEntry
	err := s.db.Table("leaderboard_view").
		Limit(10).
		Find(&entries).Error

	if err != nil {
		return
	}

	s.stepsHub.LeaderboardUpdate <- &LeaderboardUpdateMessage{
		Type:      MessageTypeLeaderboardUpdate,
		TopN:      10,
		Entries:   entries,
		Timestamp: time.Now().Unix(),
	}
}
