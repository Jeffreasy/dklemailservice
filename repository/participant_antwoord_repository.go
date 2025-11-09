package repository

import (
	"context"
	"dklautomationgo/models"
)

// PostgresParticipantAntwoordRepository implementeert ParticipantAntwoordRepository met PostgreSQL
type PostgresParticipantAntwoordRepository struct {
	*PostgresRepository
}

// NewPostgresParticipantAntwoordRepository maakt een nieuwe PostgreSQL participant antwoord repository
func NewPostgresParticipantAntwoordRepository(base *PostgresRepository) *PostgresParticipantAntwoordRepository {
	return &PostgresParticipantAntwoordRepository{
		PostgresRepository: base,
	}
}

// Create slaat een nieuw participant antwoord op
func (r *PostgresParticipantAntwoordRepository) Create(ctx context.Context, antwoord *models.ParticipantAntwoord) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(antwoord)
	return r.handleError("Create", result.Error)
}

// GetByID haalt een participant antwoord op basis van ID
func (r *PostgresParticipantAntwoordRepository) GetByID(ctx context.Context, id string) (*models.ParticipantAntwoord, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var antwoord models.ParticipantAntwoord
	result := r.DB().WithContext(ctx).First(&antwoord, "id = ?", id)
	if err := r.handleError("GetByID", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &antwoord, nil
}

// ListByParticipantID haalt alle antwoorden voor een participant op
func (r *PostgresParticipantAntwoordRepository) ListByParticipantID(ctx context.Context, participantID string) ([]*models.ParticipantAntwoord, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var antwoorden []*models.ParticipantAntwoord
	result := r.DB().WithContext(ctx).
		Where("participant_id = ?", participantID).
		Order("verzonden_op DESC").
		Find(&antwoorden)

	if err := r.handleError("ListByParticipantID", result.Error); err != nil {
		return nil, err
	}

	return antwoorden, nil
}

// Update werkt een bestaand participant antwoord bij
func (r *PostgresParticipantAntwoordRepository) Update(ctx context.Context, antwoord *models.ParticipantAntwoord) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Save(antwoord)
	return r.handleError("Update", result.Error)
}

// Delete verwijdert een participant antwoord
func (r *PostgresParticipantAntwoordRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Delete(&models.ParticipantAntwoord{}, "id = ?", id)
	return r.handleError("Delete", result.Error)
}
