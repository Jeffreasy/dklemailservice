package repository

import (
	"context"
	"dklautomationgo/models"
	"fmt"

	"gorm.io/gorm"
)

// AchievementRepository interface voor achievement operaties
type AchievementRepository interface {
	Create(ctx context.Context, achievement *models.ParticipantAchievement) error
	GetByParticipantID(ctx context.Context, participantID string) ([]models.AchievementWithBadge, error)
	GetByBadgeID(ctx context.Context, badgeID string) ([]models.ParticipantAchievement, error)
	Delete(ctx context.Context, participantID, badgeID string) error
	HasAchievement(ctx context.Context, participantID, badgeID string) (bool, error)
	GetParticipantSummary(ctx context.Context, participantID string) (*models.ParticipantAchievementsSummary, error)
}

// PostgresAchievementRepository implementatie voor PostgreSQL
type PostgresAchievementRepository struct {
	*PostgresRepository
	badgeRepo BadgeRepository
}

// NewAchievementRepository maakt een nieuwe achievement repository
func NewAchievementRepository(db *gorm.DB, badgeRepo BadgeRepository) AchievementRepository {
	return &PostgresAchievementRepository{
		PostgresRepository: NewPostgresRepository(db),
		badgeRepo:          badgeRepo,
	}
}

// Create maakt een nieuw achievement aan
func (r *PostgresAchievementRepository) Create(ctx context.Context, achievement *models.ParticipantAchievement) error {
	return r.DB().WithContext(ctx).Create(achievement).Error
}

// GetByParticipantID haalt alle achievements op voor een deelnemer
func (r *PostgresAchievementRepository) GetByParticipantID(ctx context.Context, participantID string) ([]models.AchievementWithBadge, error) {
	var results []struct {
		models.ParticipantAchievement
		Badge models.Badge `gorm:"embedded;embeddedPrefix:badge_"`
	}

	err := r.DB().WithContext(ctx).
		Table("participant_achievements pa").
		Select(`
			pa.id, pa.participant_id, pa.badge_id, pa.earned_at, pa.created_at,
			b.id as badge_id, b.name as badge_name, b.description as badge_description,
			b.icon_url as badge_icon_url, b.criteria as badge_criteria, b.points as badge_points,
			b.is_active as badge_is_active, b.display_order as badge_display_order,
			b.created_at as badge_created_at, b.updated_at as badge_updated_at
		`).
		Joins("INNER JOIN badges b ON pa.badge_id = b.id").
		Where("pa.participant_id = ?", participantID).
		Order("pa.earned_at DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	achievements := make([]models.AchievementWithBadge, len(results))
	for i, r := range results {
		achievements[i] = models.AchievementWithBadge{
			ID:            r.ID,
			ParticipantID: r.ParticipantID,
			EarnedAt:      r.EarnedAt,
			Badge:         r.Badge,
		}
	}

	return achievements, nil
}

// GetByBadgeID haalt alle achievements op voor een badge
func (r *PostgresAchievementRepository) GetByBadgeID(ctx context.Context, badgeID string) ([]models.ParticipantAchievement, error) {
	var achievements []models.ParticipantAchievement
	err := r.DB().WithContext(ctx).
		Where("badge_id = ?", badgeID).
		Order("earned_at DESC").
		Find(&achievements).Error
	return achievements, err
}

// Delete verwijdert een achievement
func (r *PostgresAchievementRepository) Delete(ctx context.Context, participantID, badgeID string) error {
	result := r.DB().WithContext(ctx).
		Where("participant_id = ? AND badge_id = ?", participantID, badgeID).
		Delete(&models.ParticipantAchievement{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("achievement niet gevonden")
	}

	return nil
}

// HasAchievement controleert of deelnemer een specifieke badge heeft
func (r *PostgresAchievementRepository) HasAchievement(ctx context.Context, participantID, badgeID string) (bool, error) {
	var count int64
	err := r.DB().WithContext(ctx).
		Model(&models.ParticipantAchievement{}).
		Where("participant_id = ? AND badge_id = ?", participantID, badgeID).
		Count(&count).Error
	return count > 0, err
}

// GetParticipantSummary haalt een volledig overzicht op van achievements voor een deelnemer
func (r *PostgresAchievementRepository) GetParticipantSummary(ctx context.Context, participantID string) (*models.ParticipantAchievementsSummary, error) {
	// Haal participant info op
	var participant struct {
		Naam string `gorm:"column:naam"`
	}

	err := r.DB().WithContext(ctx).
		Table("participants").
		Select("naam").
		Where("id = ?", participantID).
		Scan(&participant).Error

	if err != nil {
		return nil, fmt.Errorf("deelnemer niet gevonden: %w", err)
	}

	// Haal alle achievements op
	achievements, err := r.GetByParticipantID(ctx, participantID)
	if err != nil {
		return nil, err
	}

	// Bereken totale punten
	totalPoints := 0
	for _, a := range achievements {
		totalPoints += a.Badge.Points
	}

	// Haal beschikbare badges op met stats
	availableBadges, err := r.badgeRepo.GetBadgesWithStats(ctx, &participantID)
	if err != nil {
		return nil, err
	}

	return &models.ParticipantAchievementsSummary{
		ParticipantID:   participantID,
		ParticipantName: participant.Naam,
		TotalBadges:     len(achievements),
		TotalPoints:     totalPoints,
		Achievements:    achievements,
		AvailableBadges: availableBadges,
	}, nil
}
