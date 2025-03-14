package repository

import (
	"context"
	"dklautomationgo/models"
)

// PostgresAanmeldingAntwoordRepository implementeert AanmeldingAntwoordRepository met PostgreSQL
type PostgresAanmeldingAntwoordRepository struct {
	*PostgresRepository
}

// NewPostgresAanmeldingAntwoordRepository maakt een nieuwe PostgreSQL aanmelding antwoord repository
func NewPostgresAanmeldingAntwoordRepository(base *PostgresRepository) *PostgresAanmeldingAntwoordRepository {
	return &PostgresAanmeldingAntwoordRepository{
		PostgresRepository: base,
	}
}

// Create slaat een nieuw aanmelding antwoord op
func (r *PostgresAanmeldingAntwoordRepository) Create(ctx context.Context, antwoord *models.AanmeldingAntwoord) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(antwoord)
	return r.handleError("Create", result.Error)
}

// GetByID haalt een aanmelding antwoord op basis van ID
func (r *PostgresAanmeldingAntwoordRepository) GetByID(ctx context.Context, id string) (*models.AanmeldingAntwoord, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var antwoord models.AanmeldingAntwoord
	result := r.DB().WithContext(ctx).First(&antwoord, "id = ?", id)
	if err := r.handleError("GetByID", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &antwoord, nil
}

// ListByAanmeldingID haalt alle antwoorden voor een aanmelding op
func (r *PostgresAanmeldingAntwoordRepository) ListByAanmeldingID(ctx context.Context, aanmeldingID string) ([]*models.AanmeldingAntwoord, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var antwoorden []*models.AanmeldingAntwoord
	result := r.DB().WithContext(ctx).
		Where("aanmelding_id = ?", aanmeldingID).
		Order("verzonden_op DESC").
		Find(&antwoorden)

	if err := r.handleError("ListByAanmeldingID", result.Error); err != nil {
		return nil, err
	}

	return antwoorden, nil
}

// Update werkt een bestaand aanmelding antwoord bij
func (r *PostgresAanmeldingAntwoordRepository) Update(ctx context.Context, antwoord *models.AanmeldingAntwoord) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Save(antwoord)
	return r.handleError("Update", result.Error)
}

// Delete verwijdert een aanmelding antwoord
func (r *PostgresAanmeldingAntwoordRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Delete(&models.AanmeldingAntwoord{}, "id = ?", id)
	return r.handleError("Delete", result.Error)
}
