package repository

import (
	"context"
	"dklautomationgo/models"
)

// PostgresEmailTemplateRepository implementeert EmailTemplateRepository met PostgreSQL
type PostgresEmailTemplateRepository struct {
	*PostgresRepository
}

// NewPostgresEmailTemplateRepository maakt een nieuwe PostgreSQL email template repository
func NewPostgresEmailTemplateRepository(base *PostgresRepository) *PostgresEmailTemplateRepository {
	return &PostgresEmailTemplateRepository{
		PostgresRepository: base,
	}
}

// Create slaat een nieuwe email template op
func (r *PostgresEmailTemplateRepository) Create(ctx context.Context, template *models.EmailTemplate) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(template)
	return r.handleError("Create", result.Error)
}

// GetByID haalt een email template op basis van ID
func (r *PostgresEmailTemplateRepository) GetByID(ctx context.Context, id string) (*models.EmailTemplate, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var template models.EmailTemplate
	result := r.DB().WithContext(ctx).First(&template, "id = ?", id)
	if err := r.handleError("GetByID", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &template, nil
}

// GetByNaam haalt een email template op basis van naam
func (r *PostgresEmailTemplateRepository) GetByNaam(ctx context.Context, naam string) (*models.EmailTemplate, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var template models.EmailTemplate
	result := r.DB().WithContext(ctx).Where("naam = ?", naam).First(&template)
	if err := r.handleError("GetByNaam", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &template, nil
}

// List haalt een lijst van email templates op
func (r *PostgresEmailTemplateRepository) List(ctx context.Context, limit, offset int) ([]*models.EmailTemplate, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var templates []*models.EmailTemplate
	result := r.DB().WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("naam ASC").
		Find(&templates)

	if err := r.handleError("List", result.Error); err != nil {
		return nil, err
	}

	return templates, nil
}

// Update werkt een bestaande email template bij
func (r *PostgresEmailTemplateRepository) Update(ctx context.Context, template *models.EmailTemplate) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Save(template)
	return r.handleError("Update", result.Error)
}

// Delete verwijdert een email template
func (r *PostgresEmailTemplateRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Delete(&models.EmailTemplate{}, "id = ?", id)
	return r.handleError("Delete", result.Error)
}

// FindActive haalt alle actieve email templates op
func (r *PostgresEmailTemplateRepository) FindActive(ctx context.Context) ([]*models.EmailTemplate, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var templates []*models.EmailTemplate
	result := r.DB().WithContext(ctx).
		Where("is_actief = ?", true).
		Order("naam ASC").
		Find(&templates)

	if err := r.handleError("FindActive", result.Error); err != nil {
		return nil, err
	}

	return templates, nil
}
