package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// PostgresSocialEmbedRepository implements SocialEmbedRepository
type PostgresSocialEmbedRepository struct {
	db *gorm.DB
}

// NewPostgresSocialEmbedRepository creates a new social embed repository
func NewPostgresSocialEmbedRepository(db *gorm.DB) *PostgresSocialEmbedRepository {
	return &PostgresSocialEmbedRepository{db: db}
}

// Create saves a new social embed
func (r *PostgresSocialEmbedRepository) Create(ctx context.Context, embed *models.SocialEmbed) error {
	return r.db.WithContext(ctx).Create(embed).Error
}

// GetByID retrieves a social embed by ID
func (r *PostgresSocialEmbedRepository) GetByID(ctx context.Context, id string) (*models.SocialEmbed, error) {
	var embed models.SocialEmbed
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&embed).Error
	if err != nil {
		return nil, err
	}
	return &embed, nil
}

// List retrieves a paginated list of social embeds
func (r *PostgresSocialEmbedRepository) List(ctx context.Context, limit, offset int) ([]*models.SocialEmbed, error) {
	var embeds []*models.SocialEmbed
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("order_number ASC, created_at DESC").Find(&embeds).Error
	return embeds, err
}

// ListVisible retrieves only visible social embeds ordered by order_number
func (r *PostgresSocialEmbedRepository) ListVisible(ctx context.Context) ([]*models.SocialEmbed, error) {
	var embeds []*models.SocialEmbed
	err := r.db.WithContext(ctx).Where("visible = ?", true).Order("order_number ASC, created_at DESC").Find(&embeds).Error
	return embeds, err
}

// Update updates an existing social embed
func (r *PostgresSocialEmbedRepository) Update(ctx context.Context, embed *models.SocialEmbed) error {
	return r.db.WithContext(ctx).Save(embed).Error
}

// Delete removes a social embed
func (r *PostgresSocialEmbedRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.SocialEmbed{}, "id = ?", id).Error
}
