package repository

import (
	"context"
	"dklautomationgo/models"
)

// PostgresContactRepository implementeert ContactRepository met PostgreSQL
type PostgresContactRepository struct {
	*PostgresRepository
}

// NewPostgresContactRepository maakt een nieuwe PostgreSQL contact repository
func NewPostgresContactRepository(base *PostgresRepository) *PostgresContactRepository {
	return &PostgresContactRepository{
		PostgresRepository: base,
	}
}

// Create slaat een nieuw contactformulier op
func (r *PostgresContactRepository) Create(ctx context.Context, contact *models.ContactFormulier) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(contact)
	return r.handleError("Create", result.Error)
}

// GetByID haalt een contactformulier op basis van ID
func (r *PostgresContactRepository) GetByID(ctx context.Context, id string) (*models.ContactFormulier, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var contact models.ContactFormulier
	result := r.DB().WithContext(ctx).First(&contact, "id = ?", id)
	if err := r.handleError("GetByID", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &contact, nil
}

// List haalt een lijst van contactformulieren op
func (r *PostgresContactRepository) List(ctx context.Context, limit, offset int) ([]*models.ContactFormulier, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var contacts []*models.ContactFormulier
	result := r.DB().WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&contacts)

	if err := r.handleError("List", result.Error); err != nil {
		return nil, err
	}

	return contacts, nil
}

// Update werkt een bestaand contactformulier bij
func (r *PostgresContactRepository) Update(ctx context.Context, contact *models.ContactFormulier) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Save(contact)
	return r.handleError("Update", result.Error)
}

// Delete verwijdert een contactformulier
func (r *PostgresContactRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Delete(&models.ContactFormulier{}, "id = ?", id)
	return r.handleError("Delete", result.Error)
}

// FindByEmail zoekt contactformulieren op basis van email
func (r *PostgresContactRepository) FindByEmail(ctx context.Context, email string) ([]*models.ContactFormulier, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var contacts []*models.ContactFormulier
	result := r.DB().WithContext(ctx).
		Where("email = ?", email).
		Order("created_at DESC").
		Find(&contacts)

	if err := r.handleError("FindByEmail", result.Error); err != nil {
		return nil, err
	}

	return contacts, nil
}

// FindByStatus zoekt contactformulieren op basis van status
func (r *PostgresContactRepository) FindByStatus(ctx context.Context, status string) ([]*models.ContactFormulier, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var contacts []*models.ContactFormulier
	result := r.DB().WithContext(ctx).
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&contacts)

	if err := r.handleError("FindByStatus", result.Error); err != nil {
		return nil, err
	}

	return contacts, nil
}
