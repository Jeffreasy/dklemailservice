package repository

import (
	"context"
	"dklautomationgo/models"
	"fmt"

	"gorm.io/gorm"
)

// BadgeRepository interface voor badge operaties
type BadgeRepository interface {
	Create(ctx context.Context, badge *models.Badge) error
	GetByID(ctx context.Context, id string) (*models.Badge, error)
	GetByName(ctx context.Context, name string) (*models.Badge, error)
	GetAll(ctx context.Context, activeOnly bool) ([]models.Badge, error)
	Update(ctx context.Context, id string, badge *models.Badge) error
	Delete(ctx context.Context, id string) error
	GetBadgesWithStats(ctx context.Context, participantID *string) ([]models.BadgeWithStats, error)
}

// PostgresBadgeRepository implementatie voor PostgreSQL
type PostgresBadgeRepository struct {
	*PostgresRepository
}

// NewBadgeRepository maakt een nieuwe badge repository
func NewBadgeRepository(db *gorm.DB) BadgeRepository {
	return &PostgresBadgeRepository{
		PostgresRepository: NewPostgresRepository(db),
	}
}

// Create maakt een nieuwe badge aan
func (r *PostgresBadgeRepository) Create(ctx context.Context, badge *models.Badge) error {
	return r.DB().WithContext(ctx).Create(badge).Error
}

// GetByID haalt een badge op via ID
func (r *PostgresBadgeRepository) GetByID(ctx context.Context, id string) (*models.Badge, error) {
	var badge models.Badge
	err := r.DB().WithContext(ctx).Where("id = ?", id).First(&badge).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("badge niet gevonden")
	}
	return &badge, err
}

// GetByName haalt een badge op via naam
func (r *PostgresBadgeRepository) GetByName(ctx context.Context, name string) (*models.Badge, error) {
	var badge models.Badge
	err := r.DB().WithContext(ctx).Where("name = ?", name).First(&badge).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("badge niet gevonden")
	}
	return &badge, err
}

// GetAll haalt alle badges op
func (r *PostgresBadgeRepository) GetAll(ctx context.Context, activeOnly bool) ([]models.Badge, error) {
	var badges []models.Badge
	query := r.DB().WithContext(ctx)

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	err := query.Order("display_order, name").Find(&badges).Error
	return badges, err
}

// Update werkt een badge bij
func (r *PostgresBadgeRepository) Update(ctx context.Context, id string, badge *models.Badge) error {
	result := r.DB().WithContext(ctx).Model(&models.Badge{}).
		Where("id = ?", id).
		Updates(badge)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("badge niet gevonden")
	}

	return nil
}

// Delete verwijdert een badge
func (r *PostgresBadgeRepository) Delete(ctx context.Context, id string) error {
	result := r.DB().WithContext(ctx).Where("id = ?", id).Delete(&models.Badge{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("badge niet gevonden")
	}

	return nil
}

// GetBadgesWithStats haalt alle badges op met statistieken
func (r *PostgresBadgeRepository) GetBadgesWithStats(ctx context.Context, participantID *string) ([]models.BadgeWithStats, error) {
	var badges []models.BadgeWithStats

	query := `
		SELECT 
			b.id, b.name, b.description, b.icon_url, b.criteria, b.points, 
			b.is_active, b.display_order, b.created_at, b.updated_at,
			COUNT(DISTINCT pa.id) as earned_count,
			MAX(pa.earned_at) as last_earned_at,
			CASE 
				WHEN ? IS NOT NULL AND EXISTS(
					SELECT 1 FROM participant_achievements 
					WHERE badge_id = b.id AND participant_id = ?
				) THEN true 
				ELSE false 
			END as earned_by_current_user
		FROM badges b
		LEFT JOIN participant_achievements pa ON b.id = pa.badge_id
		WHERE b.is_active = true
		GROUP BY b.id, b.name, b.description, b.icon_url, b.criteria, b.points, 
		         b.is_active, b.display_order, b.created_at, b.updated_at
		ORDER BY b.display_order, b.name
	`

	err := r.DB().WithContext(ctx).Raw(query, participantID, participantID).Scan(&badges).Error
	return badges, err
}
