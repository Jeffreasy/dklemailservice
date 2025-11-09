package repository

import (
	"context"
	"dklautomationgo/models"
	"time"

	"dklautomationgo/logger"

	"github.com/google/uuid"
)

// NotulenRepository defines the interface for notulen data operations
type NotulenRepository interface {
	Create(ctx context.Context, notulen *models.Notulen) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Notulen, error)
	Update(ctx context.Context, notulen *models.Notulen) error
	Finalize(ctx context.Context, id uuid.UUID, finalizedBy uuid.UUID) error
	Archive(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters *models.NotulenSearchFilters) ([]models.Notulen, int, error)
	Search(ctx context.Context, searchQuery string, filters *models.NotulenSearchFilters) ([]models.Notulen, int, error)
	GetVersions(ctx context.Context, notulenID uuid.UUID) ([]models.NotulenVersie, error)
	GetVersion(ctx context.Context, notulenID uuid.UUID, versie int) (*models.NotulenVersie, error)
}

// PostgresNotulenRepository implements NotulenRepository with PostgreSQL
type PostgresNotulenRepository struct {
	*PostgresRepository
}

// NewPostgresNotulenRepository creates a new PostgreSQL notulen repository
func NewPostgresNotulenRepository(base *PostgresRepository) *PostgresNotulenRepository {
	return &PostgresNotulenRepository{
		PostgresRepository: base,
	}
}

// Create saves a new notulen document
func (r *PostgresNotulenRepository) Create(ctx context.Context, notulen *models.Notulen) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(notulen)
	return r.handleError("Create", result.Error)
}

// GetByID retrieves a notulen by ID
func (r *PostgresNotulenRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Notulen, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var notulen models.Notulen
	// De GORM tags in het model (Fix 1) zorgen dat GORM de kolommen correct mapt.
	result := r.DB().WithContext(ctx).First(&notulen, "id = ?", id)
	if err := r.handleError("GetByID", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &notulen, nil
}

// Update modifies an existing notulen
func (r *PostgresNotulenRepository) Update(ctx context.Context, notulen *models.Notulen) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	// GEFIXTE KOLOMNAMEN: Gebruik de namen uit de Postgres HINTs
	updateData := map[string]interface{}{
		"titel":                 notulen.Titel,
		"vergadering_datum":     notulen.VergaderingDatum,
		"locatie":               notulen.Locatie,
		"voorzitter":            notulen.Voorzitter,
		"notulist":              notulen.Notulist,
		"aanwezigen":            notulen.Aanwezigen,
		"afwezigen":             notulen.Afwezigen,
		"aanwezigen_gebruikers": notulen.AanwezigenGebruikerIDs, // GEFIXT
		"afwezigen_gebruikers":  notulen.AfwezigenGebruikerIDs,  // GEFIXT
		"aanwezigen_gasten":     notulen.AanwezigenGasten,
		"afwezigen_gasten":      notulen.AfwezigenGasten,
		"agenda_items":          notulen.AgendaItems,
		"besluiten":             notulen.Besluiten,
		"actiepunten":           notulen.Actiepunten,
		"notities":              notulen.Notities,
		"status":                notulen.Status,
		"versie":                notulen.Versie,
		"created_by":            notulen.CreatedBy,
		"created_at":            notulen.CreatedAt,
		"updated_at":            notulen.UpdatedAt,
		"updated_by":            notulen.UpdatedByID, // GEFIXT
		"finalized_at":          notulen.FinalizedAt,
		"finalized_by":          notulen.FinalizedBy,
	}

	// De GORM tags in het model (Fix 1) zorgen dat GORM dit correct mapt.
	result := r.DB().WithContext(ctx).Model(notulen).Updates(updateData)
	return r.handleError("Update", result.Error)
}

// Finalize marks a notulen as finalized
func (r *PostgresNotulenRepository) Finalize(ctx context.Context, id uuid.UUID, finalizedBy uuid.UUID) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	now := time.Now()
	result := r.DB().WithContext(ctx).Model(&models.Notulen{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       "finalized",
			"finalized_at": now,
			"finalized_by": finalizedBy,
			"updated_at":   now,
		})

	return r.handleError("Finalize", result.Error)
}

// Archive marks a notulen as archived
func (r *PostgresNotulenRepository) Archive(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Model(&models.Notulen{}).
		Where("id = ?", id).
		Update("status", "archived")

	return r.handleError("Archive", result.Error)
}

// Delete performs a soft delete on a notulen
func (r *PostgresNotulenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Delete(&models.Notulen{}, "id = ?", id)
	return r.handleError("Delete", result.Error)
}

// List retrieves notulen with filtering and pagination
func (r *PostgresNotulenRepository) List(ctx context.Context, filters *models.NotulenSearchFilters) ([]models.Notulen, int, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	query := r.DB().WithContext(ctx).Model(&models.Notulen{})

	// Apply filters
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.CreatedBy != nil {
		query = query.Where("created_by = ?", *filters.CreatedBy)
	}
	if filters.DateFrom != nil {
		query = query.Where("vergadering_datum >= ?", filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("vergadering_datum <= ?", filters.DateTo)
	}

	// Get total count (before applying pagination)
	var total int64
	logger.Info("Executing count query for List")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, r.handleError("List count", err)
	}

	// Apply pagination with defaults
	limit := 50
	if filters.Limit > 0 {
		limit = filters.Limit
	}
	offset := 0
	if filters.Offset > 0 {
		offset = filters.Offset
	}

	// Execute query with ordering - SELECT only the columns we need to avoid scanning issues
	var notulen []models.Notulen
	// GEFIXTE KOLOMNAMEN: Gebruik de namen uit de Postgres HINTs
	result := query.
		Select("id, titel, vergadering_datum, locatie, voorzitter, notulist, aanwezigen, afwezigen, agenda_items, besluiten, actiepunten, notities, status, versie, created_by, created_at, updated_at, finalized_at, finalized_by, updated_by, aanwezigen_gebruikers, afwezigen_gebruikers, aanwezigen_gasten, afwezigen_gasten").
		Order("vergadering_datum DESC").
		Limit(limit).
		Offset(offset).
		Find(&notulen) // Gebruik .Find()

	// Debug: Log the actual query being executed
	logger.Info("Executing List query", "query", query.Statement.SQL.String(), "vars", query.Statement.Vars)

	if err := r.handleError("List", result.Error); err != nil {
		return nil, 0, err
	}

	return notulen, int(total), nil
}

// Search performs full-text search on notulen
func (r *PostgresNotulenRepository) Search(ctx context.Context, searchQuery string, filters *models.NotulenSearchFilters) ([]models.Notulen, int, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	query := r.DB().WithContext(ctx).Model(&models.Notulen{})

	// Apply full-text search if query provided
	if searchQuery != "" {
		query = query.Where("to_tsvector('dutch', titel || ' ' || COALESCE(notities, '')) @@ plainto_tsquery('dutch', ?)", searchQuery)
	}

	// Apply same filters as List
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.CreatedBy != nil {
		query = query.Where("created_by = ?", *filters.CreatedBy)
	}
	if filters.DateFrom != nil {
		query = query.Where("vergadering_datum >= ?", filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("vergadering_datum <= ?", filters.DateTo)
	}

	// Get total count (before applying pagination)
	var total int64
	logger.Info("Executing count query for Search")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, r.handleError("Search count", err)
	}

	// Apply pagination with defaults
	limit := 50
	if filters.Limit > 0 {
		limit = filters.Limit
	}
	offset := 0
	if filters.Offset > 0 {
		offset = filters.Offset
	}

	// Execute query with ordering - SELECT only the columns we need to avoid scanning issues
	var notulen []models.Notulen
	// GEFIXTE KOLOMNAMEN: Gebruik de namen uit de Postgres HINTs
	result := query.
		Select("id, titel, vergadering_datum, locatie, voorzitter, notulist, aanwezigen, afwezigen, agenda_items, besluiten, actiepunten, notities, status, versie, created_by, created_at, updated_at, finalized_at, finalized_by, updated_by, aanwezigen_gebruikers, afwezigen_gebruikers, aanwezigen_gasten, afwezigen_gasten").
		Order("vergadering_datum DESC").
		Limit(limit).
		Offset(offset).
		Find(&notulen) // Gebruik .Find()

	if err := r.handleError("Search", result.Error); err != nil {
		return nil, 0, err
	}

	return notulen, int(total), nil
}

// GetVersions retrieves all versions of a notulen
func (r *PostgresNotulenRepository) GetVersions(ctx context.Context, notulenID uuid.UUID) ([]models.NotulenVersie, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var versions []models.NotulenVersie
	result := r.DB().WithContext(ctx).
		Where("notulen_id = ?", notulenID).
		Order("versie DESC").
		Find(&versions)

	if err := r.handleError("GetVersions", result.Error); err != nil {
		return nil, err
	}

	return versions, nil
}

// GetVersion retrieves a specific version of a notulen
func (r *PostgresNotulenRepository) GetVersion(ctx context.Context, notulenID uuid.UUID, versie int) (*models.NotulenVersie, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var version models.NotulenVersie
	result := r.DB().WithContext(ctx).
		Where("notulen_id = ? AND versie = ?", notulenID, versie).
		First(&version)

	if err := r.handleError("GetVersion", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &version, nil
}
