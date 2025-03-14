package repository

import (
	"context"
	"dklautomationgo/models"
)

// PostgresContactAntwoordRepository implementeert ContactAntwoordRepository met PostgreSQL
type PostgresContactAntwoordRepository struct {
	*PostgresRepository
}

// NewPostgresContactAntwoordRepository maakt een nieuwe PostgreSQL contact antwoord repository
func NewPostgresContactAntwoordRepository(base *PostgresRepository) *PostgresContactAntwoordRepository {
	return &PostgresContactAntwoordRepository{
		PostgresRepository: base,
	}
}

// Create slaat een nieuw contact antwoord op
func (r *PostgresContactAntwoordRepository) Create(ctx context.Context, antwoord *models.ContactAntwoord) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(antwoord)
	return r.handleError("Create", result.Error)
}

// GetByID haalt een contact antwoord op basis van ID
func (r *PostgresContactAntwoordRepository) GetByID(ctx context.Context, id string) (*models.ContactAntwoord, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var antwoord models.ContactAntwoord
	result := r.DB().WithContext(ctx).First(&antwoord, "id = ?", id)
	if err := r.handleError("GetByID", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &antwoord, nil
}

// ListByContactID haalt alle antwoorden voor een contact op
func (r *PostgresContactAntwoordRepository) ListByContactID(ctx context.Context, contactID string) ([]*models.ContactAntwoord, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var antwoorden []*models.ContactAntwoord
	result := r.DB().WithContext(ctx).
		Where("contact_id = ?", contactID).
		Order("verzonden_op DESC").
		Find(&antwoorden)

	if err := r.handleError("ListByContactID", result.Error); err != nil {
		return nil, err
	}

	return antwoorden, nil
}

// Update werkt een bestaand contact antwoord bij
func (r *PostgresContactAntwoordRepository) Update(ctx context.Context, antwoord *models.ContactAntwoord) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Save(antwoord)
	return r.handleError("Update", result.Error)
}

// Delete verwijdert een contact antwoord
func (r *PostgresContactAntwoordRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Delete(&models.ContactAntwoord{}, "id = ?", id)
	return r.handleError("Delete", result.Error)
}
