package repository

import (
	"context"
	"dklautomationgo/models"
)

// PostgresVerzondEmailRepository implementeert VerzondEmailRepository met PostgreSQL
type PostgresVerzondEmailRepository struct {
	*PostgresRepository
}

// NewPostgresVerzondEmailRepository maakt een nieuwe PostgreSQL verzonden email repository
func NewPostgresVerzondEmailRepository(base *PostgresRepository) *PostgresVerzondEmailRepository {
	return &PostgresVerzondEmailRepository{
		PostgresRepository: base,
	}
}

// Create slaat een nieuwe verzonden email op
func (r *PostgresVerzondEmailRepository) Create(ctx context.Context, email *models.VerzondEmail) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(email)
	return r.handleError("Create", result.Error)
}

// GetByID haalt een verzonden email op basis van ID
func (r *PostgresVerzondEmailRepository) GetByID(ctx context.Context, id string) (*models.VerzondEmail, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var email models.VerzondEmail
	result := r.DB().WithContext(ctx).First(&email, "id = ?", id)
	if err := r.handleError("GetByID", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &email, nil
}

// List haalt een lijst van verzonden emails op
func (r *PostgresVerzondEmailRepository) List(ctx context.Context, limit, offset int) ([]*models.VerzondEmail, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var emails []*models.VerzondEmail
	result := r.DB().WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("verzonden_op DESC").
		Find(&emails)

	if err := r.handleError("List", result.Error); err != nil {
		return nil, err
	}

	return emails, nil
}

// Update werkt een bestaande verzonden email bij
func (r *PostgresVerzondEmailRepository) Update(ctx context.Context, email *models.VerzondEmail) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Save(email)
	return r.handleError("Update", result.Error)
}

// FindByContactID haalt verzonden emails op basis van contact ID
func (r *PostgresVerzondEmailRepository) FindByContactID(ctx context.Context, contactID string) ([]*models.VerzondEmail, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var emails []*models.VerzondEmail
	result := r.DB().WithContext(ctx).
		Where("contact_id = ?", contactID).
		Order("verzonden_op DESC").
		Find(&emails)

	if err := r.handleError("FindByContactID", result.Error); err != nil {
		return nil, err
	}

	return emails, nil
}

// FindByAanmeldingID haalt verzonden emails op basis van aanmelding ID
func (r *PostgresVerzondEmailRepository) FindByAanmeldingID(ctx context.Context, aanmeldingID string) ([]*models.VerzondEmail, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var emails []*models.VerzondEmail
	result := r.DB().WithContext(ctx).
		Where("aanmelding_id = ?", aanmeldingID).
		Order("verzonden_op DESC").
		Find(&emails)

	if err := r.handleError("FindByAanmeldingID", result.Error); err != nil {
		return nil, err
	}

	return emails, nil
}

// FindByOntvanger haalt verzonden emails op basis van ontvanger
func (r *PostgresVerzondEmailRepository) FindByOntvanger(ctx context.Context, ontvanger string) ([]*models.VerzondEmail, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var emails []*models.VerzondEmail
	result := r.DB().WithContext(ctx).
		Where("ontvanger = ?", ontvanger).
		Order("verzonden_op DESC").
		Find(&emails)

	if err := r.handleError("FindByOntvanger", result.Error); err != nil {
		return nil, err
	}

	return emails, nil
}
