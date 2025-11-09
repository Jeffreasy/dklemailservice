package services

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"
	"dklautomationgo/repository"
	"fmt"
	"time"
)

// GamificationService bevat de business logica voor gamification
type GamificationService struct {
	badgeRepo       repository.BadgeRepository
	achievementRepo repository.AchievementRepository
	leaderboardRepo repository.LeaderboardRepository
	aanmeldingRepo  repository.ParticipantRepository
	eventRegRepo    repository.EventRegistrationRepository
}

// NewGamificationService maakt een nieuwe gamification service
func NewGamificationService(
	badgeRepo repository.BadgeRepository,
	achievementRepo repository.AchievementRepository,
	leaderboardRepo repository.LeaderboardRepository,
	aanmeldingRepo repository.ParticipantRepository,
	eventRegRepo repository.EventRegistrationRepository,
) *GamificationService {
	return &GamificationService{
		badgeRepo:       badgeRepo,
		achievementRepo: achievementRepo,
		leaderboardRepo: leaderboardRepo,
		aanmeldingRepo:  aanmeldingRepo,
		eventRegRepo:    eventRegRepo,
	}
}

// ==================== BADGE MANAGEMENT ====================

// CreateBadge maakt een nieuwe badge aan
func (s *GamificationService) CreateBadge(ctx context.Context, req models.BadgeRequest) (*models.Badge, error) {
	// Check if badge with same name exists
	existingBadge, _ := s.badgeRepo.GetByName(ctx, req.Name)
	if existingBadge != nil {
		return nil, fmt.Errorf("badge met naam '%s' bestaat al", req.Name)
	}

	badge := &models.Badge{
		Name:         req.Name,
		Description:  req.Description,
		IconURL:      req.IconURL,
		Criteria:     req.Criteria,
		Points:       req.Points,
		IsActive:     true,
		DisplayOrder: 0,
	}

	if req.IsActive != nil {
		badge.IsActive = *req.IsActive
	}

	if req.DisplayOrder != nil {
		badge.DisplayOrder = *req.DisplayOrder
	}

	err := s.badgeRepo.Create(ctx, badge)
	if err != nil {
		logger.Error("Fout bij aanmaken badge", "error", err, "name", req.Name)
		return nil, fmt.Errorf("kon badge niet aanmaken: %w", err)
	}

	logger.Info("Badge aangemaakt", "id", badge.ID, "name", badge.Name)
	return badge, nil
}

// GetBadge haalt een badge op via ID
func (s *GamificationService) GetBadge(ctx context.Context, id string) (*models.Badge, error) {
	return s.badgeRepo.GetByID(ctx, id)
}

// GetAllBadges haalt alle badges op
func (s *GamificationService) GetAllBadges(ctx context.Context, activeOnly bool) ([]models.Badge, error) {
	return s.badgeRepo.GetAll(ctx, activeOnly)
}

// GetBadgesWithStats haalt badges op met statistieken
func (s *GamificationService) GetBadgesWithStats(ctx context.Context, participantID *string) ([]models.BadgeWithStats, error) {
	return s.badgeRepo.GetBadgesWithStats(ctx, participantID)
}

// UpdateBadge werkt een badge bij
func (s *GamificationService) UpdateBadge(ctx context.Context, id string, req models.BadgeRequest) (*models.Badge, error) {
	// Check if badge exists
	existingBadge, err := s.badgeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if name is being changed to existing name
	if req.Name != existingBadge.Name {
		nameCheck, _ := s.badgeRepo.GetByName(ctx, req.Name)
		if nameCheck != nil {
			return nil, fmt.Errorf("badge met naam '%s' bestaat al", req.Name)
		}
	}

	badge := &models.Badge{
		Name:         req.Name,
		Description:  req.Description,
		IconURL:      req.IconURL,
		Criteria:     req.Criteria,
		Points:       req.Points,
		IsActive:     existingBadge.IsActive,
		DisplayOrder: existingBadge.DisplayOrder,
	}

	if req.IsActive != nil {
		badge.IsActive = *req.IsActive
	}

	if req.DisplayOrder != nil {
		badge.DisplayOrder = *req.DisplayOrder
	}

	err = s.badgeRepo.Update(ctx, id, badge)
	if err != nil {
		logger.Error("Fout bij bijwerken badge", "error", err, "id", id)
		return nil, fmt.Errorf("kon badge niet bijwerken: %w", err)
	}

	logger.Info("Badge bijgewerkt", "id", id, "name", badge.Name)

	// Haal bijgewerkte badge op met alle data
	return s.badgeRepo.GetByID(ctx, id)
}

// DeleteBadge verwijdert een badge
func (s *GamificationService) DeleteBadge(ctx context.Context, id string) error {
	err := s.badgeRepo.Delete(ctx, id)
	if err != nil {
		logger.Error("Fout bij verwijderen badge", "error", err, "id", id)
		return fmt.Errorf("kon badge niet verwijderen: %w", err)
	}

	logger.Info("Badge verwijderd", "id", id)
	return nil
}

// ==================== ACHIEVEMENT MANAGEMENT ====================

// AwardBadge kent een badge toe aan een deelnemer
func (s *GamificationService) AwardBadge(ctx context.Context, participantID, badgeID string) (*models.AchievementWithBadge, error) {
	// Check if badge exists
	badge, err := s.badgeRepo.GetByID(ctx, badgeID)
	if err != nil {
		return nil, err
	}

	if !badge.IsActive {
		return nil, fmt.Errorf("badge is niet actief")
	}

	// Check if participant exists
	_, err = s.aanmeldingRepo.GetByID(ctx, participantID)
	if err != nil {
		return nil, fmt.Errorf("deelnemer niet gevonden")
	}

	// Check if already has achievement
	hasAchievement, err := s.achievementRepo.HasAchievement(ctx, participantID, badgeID)
	if err != nil {
		return nil, err
	}

	if hasAchievement {
		return nil, fmt.Errorf("deelnemer heeft deze badge al")
	}

	// Create achievement
	achievement := &models.ParticipantAchievement{
		ParticipantID: participantID,
		BadgeID:       badgeID,
		EarnedAt:      time.Now(),
	}

	err = s.achievementRepo.Create(ctx, achievement)
	if err != nil {
		logger.Error("Fout bij toekennen badge", "error", err, "participant_id", participantID, "badge_id", badgeID)
		return nil, fmt.Errorf("kon badge niet toekennen: %w", err)
	}

	logger.Info("Badge toegekend", "participant_id", participantID, "badge_id", badgeID, "badge_name", badge.Name)

	return &models.AchievementWithBadge{
		ID:            achievement.ID,
		ParticipantID: achievement.ParticipantID,
		EarnedAt:      achievement.EarnedAt,
		Badge:         *badge,
	}, nil
}

// RemoveBadge verwijdert een badge van een deelnemer
func (s *GamificationService) RemoveBadge(ctx context.Context, participantID, badgeID string) error {
	err := s.achievementRepo.Delete(ctx, participantID, badgeID)
	if err != nil {
		logger.Error("Fout bij verwijderen achievement", "error", err, "participant_id", participantID, "badge_id", badgeID)
		return fmt.Errorf("kon achievement niet verwijderen: %w", err)
	}

	logger.Info("Achievement verwijderd", "participant_id", participantID, "badge_id", badgeID)
	return nil
}

// GetParticipantAchievements haalt alle achievements op voor een deelnemer
func (s *GamificationService) GetParticipantAchievements(ctx context.Context, participantID string) ([]models.AchievementWithBadge, error) {
	return s.achievementRepo.GetByParticipantID(ctx, participantID)
}

// GetParticipantSummary haalt een volledig overzicht op van achievements
func (s *GamificationService) GetParticipantSummary(ctx context.Context, participantID string) (*models.ParticipantAchievementsSummary, error) {
	return s.achievementRepo.GetParticipantSummary(ctx, participantID)
}

// ==================== AUTO-AWARD LOGIC ====================

// CheckAndAwardBadges controleert of een deelnemer nieuwe badges heeft verdiend
func (s *GamificationService) CheckAndAwardBadges(ctx context.Context, participantID string) ([]models.Badge, error) {
	// Haal event registration data op voor badge criteria
	eventReg, err := s.GetEventRegistrationByParticipantID(ctx, participantID)
	if err != nil {
		logger.Error("Kon event registratie niet ophalen voor badge controle", "participant_id", participantID, "error", err)
		return []models.Badge{}, nil
	}
	if eventReg == nil {
		logger.Debug("Geen actieve event registratie gevonden voor badge controle", "participant_id", participantID)
		return []models.Badge{}, nil
	}

	var awardedBadges []models.Badge

	// Haal alle actieve badges op
	badges, err := s.badgeRepo.GetAll(ctx, true) // activeOnly = true
	if err != nil {
		logger.Error("Kon badges niet ophalen voor controle", "error", err)
		return []models.Badge{}, nil
	}

	// Controleer elke badge
	for _, badge := range badges {
		if s.meetsCriteria(eventReg, &badge) {
			// Controleer of deelnemer deze badge al heeft
			hasBadge, err := s.achievementRepo.HasAchievement(ctx, participantID, badge.ID)
			if err != nil {
				logger.Error("Kon achievement status niet controleren", "participant_id", participantID, "badge_id", badge.ID, "error", err)
				continue
			}

			if !hasBadge {
				// Ken badge toe
				_, err := s.AwardBadge(ctx, participantID, badge.ID)
				if err != nil {
					logger.Error("Kon badge niet toekennen", "participant_id", participantID, "badge_id", badge.ID, "error", err)
					continue
				}
				awardedBadges = append(awardedBadges, badge)
			}
		}
	}

	return awardedBadges, nil
}

// GetEventRegistrationByParticipantID haalt de actieve event registratie op voor een participant
func (s *GamificationService) GetEventRegistrationByParticipantID(ctx context.Context, participantID string) (*models.EventRegistration, error) {
	return s.eventRegRepo.GetActiveRegistrationForParticipant(ctx, participantID)
}

// meetsCriteria controleert of een deelnemer aan badge criteria voldoet
func (s *GamificationService) meetsCriteria(eventReg *models.EventRegistration, badge *models.Badge) bool {
	criteria := badge.Criteria

	// Check min_steps
	if criteria.MinSteps != nil && eventReg.Steps < *criteria.MinSteps {
		return false
	}

	// Check route
	if criteria.Route != nil && eventReg.DistanceRoute != nil && *eventReg.DistanceRoute != *criteria.Route {
		return false
	}

	// Note: criteria zoals consecutive_days, early_participant, has_team
	// vereisen extra logica/data die je later kunt implementeren

	return true
}

// ==================== LEADERBOARD ====================

// GetLeaderboard haalt de leaderboard op
func (s *GamificationService) GetLeaderboard(ctx context.Context, filters models.LeaderboardFilters) (*models.LeaderboardResponse, error) {
	return s.leaderboardRepo.GetLeaderboard(ctx, filters)
}

// GetParticipantRank haalt rank informatie op voor een deelnemer
func (s *GamificationService) GetParticipantRank(ctx context.Context, participantID string) (*models.ParticipantRankInfo, error) {
	return s.leaderboardRepo.GetParticipantRank(ctx, participantID)
}
