package repository

import (
	"context"
	"dklautomationgo/models"
)

// PostgresParticipantRoleRepository implementeert ParticipantRoleRepository met PostgreSQL
type PostgresParticipantRoleRepository struct {
	*PostgresRepository
}

// NewPostgresParticipantRoleRepository maakt een nieuwe PostgreSQL participant role repository
func NewPostgresParticipantRoleRepository(base *PostgresRepository) *PostgresParticipantRoleRepository {
	return &PostgresParticipantRoleRepository{
		PostgresRepository: base,
	}
}

// List haalt alle participant rollen op
func (r *PostgresParticipantRoleRepository) List(ctx context.Context) ([]*models.ParticipantRole, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var roles []*models.ParticipantRole
	result := r.DB().WithContext(ctx).
		Where("is_active = true").
		Order("name ASC").
		Find(&roles)

	if err := r.handleError("List", result.Error); err != nil {
		return nil, err
	}
	return roles, nil
}

// Create slaat een nieuwe participant rol op
func (r *PostgresParticipantRoleRepository) Create(ctx context.Context, role *models.ParticipantRole) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(role)
	return r.handleError("Create", result.Error)
}

// GetByName haalt een participant rol op basis van naam
func (r *PostgresParticipantRoleRepository) GetByName(ctx context.Context, name string) (*models.ParticipantRole, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var role models.ParticipantRole
	result := r.DB().WithContext(ctx).
		Where("name = ? AND is_active = true", name).
		First(&role)

	if err := r.handleError("GetByName", result.Error); err != nil {
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &role, nil
}

// Update werkt een bestaande participant rol bij
func (r *PostgresParticipantRoleRepository) Update(ctx context.Context, role *models.ParticipantRole) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Save(role)
	return r.handleError("Update", result.Error)
}

// Delete verwijdert een participant rol (soft delete door is_active = false)
func (r *PostgresParticipantRoleRepository) Delete(ctx context.Context, name string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).
		Model(&models.ParticipantRole{}).
		Where("name = ?", name).
		Update("is_active", false)
	return r.handleError("Delete", result.Error)
}

// GetByID haalt een participant rol op basis van ID
func (r *PostgresParticipantRoleRepository) GetByID(ctx context.Context, id string) (*models.ParticipantRole, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var role models.ParticipantRole
	result := r.DB().WithContext(ctx).
		Where("id = ? AND is_active = true", id).
		First(&role)

	if err := r.handleError("GetByID", result.Error); err != nil {
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &role, nil
}

// ListActive haalt alleen actieve participant rollen op (zelfde als List)
func (r *PostgresParticipantRoleRepository) ListActive(ctx context.Context) ([]*models.ParticipantRole, error) {
	return r.List(ctx)
}
