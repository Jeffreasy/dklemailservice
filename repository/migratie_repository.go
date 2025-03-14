package repository

import (
	"context"
	"dklautomationgo/models"
)

// PostgresMigratieRepository implementeert MigratieRepository met PostgreSQL
type PostgresMigratieRepository struct {
	*PostgresRepository
}

// NewPostgresMigratieRepository maakt een nieuwe PostgreSQL migratie repository
func NewPostgresMigratieRepository(base *PostgresRepository) *PostgresMigratieRepository {
	return &PostgresMigratieRepository{
		PostgresRepository: base,
	}
}

// Create slaat een nieuwe migratie op
func (r *PostgresMigratieRepository) Create(ctx context.Context, migratie *models.Migratie) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(migratie)
	return r.handleError("Create", result.Error)
}

// GetByVersie haalt een migratie op basis van versie
func (r *PostgresMigratieRepository) GetByVersie(ctx context.Context, versie string) (*models.Migratie, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var migratie models.Migratie
	result := r.DB().WithContext(ctx).Where("versie = ?", versie).First(&migratie)
	if err := r.handleError("GetByVersie", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &migratie, nil
}

// List haalt een lijst van migraties op
func (r *PostgresMigratieRepository) List(ctx context.Context) ([]*models.Migratie, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var migraties []*models.Migratie
	result := r.DB().WithContext(ctx).
		Order("id ASC").
		Find(&migraties)

	if err := r.handleError("List", result.Error); err != nil {
		return nil, err
	}

	return migraties, nil
}

// GetLatest haalt de laatste migratie op
func (r *PostgresMigratieRepository) GetLatest(ctx context.Context) (*models.Migratie, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var migratie models.Migratie
	result := r.DB().WithContext(ctx).
		Order("id DESC").
		First(&migratie)

	if err := r.handleError("GetLatest", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &migratie, nil
}
