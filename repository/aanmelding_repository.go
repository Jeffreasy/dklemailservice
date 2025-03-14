package repository

import (
	"context"
	"dklautomationgo/models"
)

// PostgresAanmeldingRepository implementeert AanmeldingRepository met PostgreSQL
type PostgresAanmeldingRepository struct {
	*PostgresRepository
}

// NewPostgresAanmeldingRepository maakt een nieuwe PostgreSQL aanmelding repository
func NewPostgresAanmeldingRepository(base *PostgresRepository) *PostgresAanmeldingRepository {
	return &PostgresAanmeldingRepository{
		PostgresRepository: base,
	}
}

// Create slaat een nieuwe aanmelding op
func (r *PostgresAanmeldingRepository) Create(ctx context.Context, aanmelding *models.Aanmelding) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(aanmelding)
	return r.handleError("Create", result.Error)
}

// GetByID haalt een aanmelding op basis van ID
func (r *PostgresAanmeldingRepository) GetByID(ctx context.Context, id string) (*models.Aanmelding, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var aanmelding models.Aanmelding
	result := r.DB().WithContext(ctx).First(&aanmelding, "id = ?", id)
	if err := r.handleError("GetByID", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &aanmelding, nil
}

// List haalt een lijst van aanmeldingen op
func (r *PostgresAanmeldingRepository) List(ctx context.Context, limit, offset int) ([]*models.Aanmelding, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var aanmeldingen []*models.Aanmelding
	result := r.DB().WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&aanmeldingen)

	if err := r.handleError("List", result.Error); err != nil {
		return nil, err
	}

	return aanmeldingen, nil
}

// Update werkt een bestaande aanmelding bij
func (r *PostgresAanmeldingRepository) Update(ctx context.Context, aanmelding *models.Aanmelding) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Save(aanmelding)
	return r.handleError("Update", result.Error)
}

// Delete verwijdert een aanmelding
func (r *PostgresAanmeldingRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Delete(&models.Aanmelding{}, "id = ?", id)
	return r.handleError("Delete", result.Error)
}

// FindByEmail zoekt aanmeldingen op basis van email
func (r *PostgresAanmeldingRepository) FindByEmail(ctx context.Context, email string) ([]*models.Aanmelding, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var aanmeldingen []*models.Aanmelding
	result := r.DB().WithContext(ctx).
		Where("email = ?", email).
		Order("created_at DESC").
		Find(&aanmeldingen)

	if err := r.handleError("FindByEmail", result.Error); err != nil {
		return nil, err
	}

	return aanmeldingen, nil
}

// FindByStatus zoekt aanmeldingen op basis van status
func (r *PostgresAanmeldingRepository) FindByStatus(ctx context.Context, status string) ([]*models.Aanmelding, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var aanmeldingen []*models.Aanmelding
	result := r.DB().WithContext(ctx).
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&aanmeldingen)

	if err := r.handleError("FindByStatus", result.Error); err != nil {
		return nil, err
	}

	return aanmeldingen, nil
}
