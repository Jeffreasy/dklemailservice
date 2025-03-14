package repository

import (
	"context"
	"dklautomationgo/models"
	"time"
)

// PostgresGebruikerRepository implementeert GebruikerRepository met PostgreSQL
type PostgresGebruikerRepository struct {
	*PostgresRepository
}

// NewPostgresGebruikerRepository maakt een nieuwe PostgreSQL gebruiker repository
func NewPostgresGebruikerRepository(base *PostgresRepository) *PostgresGebruikerRepository {
	return &PostgresGebruikerRepository{
		PostgresRepository: base,
	}
}

// Create slaat een nieuwe gebruiker op
func (r *PostgresGebruikerRepository) Create(ctx context.Context, gebruiker *models.Gebruiker) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(gebruiker)
	return r.handleError("Create", result.Error)
}

// GetByID haalt een gebruiker op basis van ID
func (r *PostgresGebruikerRepository) GetByID(ctx context.Context, id string) (*models.Gebruiker, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var gebruiker models.Gebruiker
	result := r.DB().WithContext(ctx).First(&gebruiker, "id = ?", id)
	if err := r.handleError("GetByID", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &gebruiker, nil
}

// GetByEmail haalt een gebruiker op basis van email
func (r *PostgresGebruikerRepository) GetByEmail(ctx context.Context, email string) (*models.Gebruiker, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var gebruiker models.Gebruiker
	result := r.DB().WithContext(ctx).Where("email = ?", email).First(&gebruiker)
	if err := r.handleError("GetByEmail", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &gebruiker, nil
}

// List haalt een lijst van gebruikers op
func (r *PostgresGebruikerRepository) List(ctx context.Context, limit, offset int) ([]*models.Gebruiker, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var gebruikers []*models.Gebruiker
	result := r.DB().WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("naam ASC").
		Find(&gebruikers)

	if err := r.handleError("List", result.Error); err != nil {
		return nil, err
	}

	return gebruikers, nil
}

// Update werkt een bestaande gebruiker bij
func (r *PostgresGebruikerRepository) Update(ctx context.Context, gebruiker *models.Gebruiker) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Save(gebruiker)
	return r.handleError("Update", result.Error)
}

// Delete verwijdert een gebruiker
func (r *PostgresGebruikerRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Delete(&models.Gebruiker{}, "id = ?", id)
	return r.handleError("Delete", result.Error)
}

// UpdateLastLogin werkt de laatste login tijd van een gebruiker bij
func (r *PostgresGebruikerRepository) UpdateLastLogin(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	now := time.Now()
	result := r.DB().WithContext(ctx).
		Model(&models.Gebruiker{}).
		Where("id = ?", id).
		Update("laatste_login", now)

	return r.handleError("UpdateLastLogin", result.Error)
}
