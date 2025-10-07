package repository

import (
	"context"
	"dklautomationgo/models"
	"time"
)

// PostgresNewsletterRepository implementeert NewsletterRepository met PostgreSQL
type PostgresNewsletterRepository struct {
	*PostgresRepository
}

// NewPostgresNewsletterRepository maakt een nieuwe PostgreSQL newsletter repository
func NewPostgresNewsletterRepository(base *PostgresRepository) *PostgresNewsletterRepository {
	return &PostgresNewsletterRepository{PostgresRepository: base}
}

func (r *PostgresNewsletterRepository) Create(ctx context.Context, nl *models.Newsletter) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	result := r.DB().WithContext(ctx).Create(nl)
	return r.handleError("Create", result.Error)
}

func (r *PostgresNewsletterRepository) GetByID(ctx context.Context, id string) (*models.Newsletter, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var nl models.Newsletter
	result := r.DB().WithContext(ctx).First(&nl, "id = ?", id)
	if err := r.handleError("GetByID", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &nl, nil
}

func (r *PostgresNewsletterRepository) List(ctx context.Context, limit, offset int) ([]*models.Newsletter, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var newsletters []*models.Newsletter
	result := r.DB().WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&newsletters)

	if err := r.handleError("List", result.Error); err != nil {
		return nil, err
	}

	return newsletters, nil
}

func (r *PostgresNewsletterRepository) Update(ctx context.Context, nl *models.Newsletter) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	result := r.DB().WithContext(ctx).Save(nl)
	return r.handleError("Update", result.Error)
}

func (r *PostgresNewsletterRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	result := r.DB().WithContext(ctx).Delete(&models.Newsletter{}, "id = ?", id)
	return r.handleError("Delete", result.Error)
}

func (r *PostgresNewsletterRepository) UpdateBatchID(ctx context.Context, id, batchID string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	result := r.DB().WithContext(ctx).Model(&models.Newsletter{}).Where("id = ?", id).Update("batch_id", batchID)
	return r.handleError("UpdateBatchID", result.Error)
}

func (r *PostgresNewsletterRepository) MarkSent(ctx context.Context, id string, sentAt time.Time) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	result := r.DB().WithContext(ctx).Model(&models.Newsletter{}).Where("id = ?", id).Update("sent_at", sentAt)
	return r.handleError("MarkSent", result.Error)
}
