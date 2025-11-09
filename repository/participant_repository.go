package repository

import (
	"context"
	"dklautomationgo/models"
)

// PostgresParticipantRepository implementeert ParticipantRepository met PostgreSQL
type PostgresParticipantRepository struct {
	*PostgresRepository
}

// NewPostgresParticipantRepository maakt een nieuwe PostgreSQL participant repository
func NewPostgresParticipantRepository(base *PostgresRepository) *PostgresParticipantRepository {
	return &PostgresParticipantRepository{
		PostgresRepository: base,
	}
}

// Standaard kolommen voor een Participant "Persoon"
// Alle event-specifieke data (steps, status, route) is hier bewust weggelaten.
var participantColumns = []string{
	"id",
	"created_at",
	"updated_at",
	"naam",
	"email",
	"telefoon",
	"terms",
	"gebruiker_id",
	"test_mode",
}

// Create slaat een nieuwe participant op
func (r *PostgresParticipantRepository) Create(ctx context.Context, participant *models.Participant) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(participant)
	return r.handleError("Create", result.Error)
}

// GetByID haalt een participant op basis van ID
func (r *PostgresParticipantRepository) GetByID(ctx context.Context, id string) (*models.Participant, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var participant models.Participant
	// GEFIXTE QUERY: Gebruik 'Select' om scan errors te voorkomen
	result := r.DB().WithContext(ctx).
		Select(participantColumns).
		First(&participant, "id = ?", id)

	if err := r.handleError("GetByID", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &participant, nil
}

// List haalt een lijst van participants op
func (r *PostgresParticipantRepository) List(ctx context.Context, limit, offset int) ([]*models.Participant, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var participants []*models.Participant
	// GEFIXTE QUERY: Gebruik 'Select' om scan errors te voorkomen
	result := r.DB().WithContext(ctx).
		Select(participantColumns).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&participants)

	if err := r.handleError("List", result.Error); err != nil {
		return nil, err
	}

	return participants, nil
}

// Update werkt een bestaande participant bij
func (r *PostgresParticipantRepository) Update(ctx context.Context, participant *models.Participant) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	// Update alleen de 'participant' velden
	result := r.DB().WithContext(ctx).Model(participant).
		Select(participantColumns). // Zorgt dat we geen event-data proberen te updaten
		Updates(participant)

	return r.handleError("Update", result.Error)
}

// Delete verwijdert een participant
func (r *PostgresParticipantRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Delete(&models.Participant{}, "id = ?", id)
	return r.handleError("Delete", result.Error)
}

// FindByEmail zoekt participants op basis van email
func (r *PostgresParticipantRepository) FindByEmail(ctx context.Context, email string) ([]*models.Participant, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var participants []*models.Participant
	// GEFIXTE QUERY: Gebruik 'Select' om scan errors te voorkomen
	result := r.DB().WithContext(ctx).
		Select(participantColumns).
		Where("email = ?", email).
		Order("created_at DESC").
		Find(&participants)

	if err := r.handleError("FindByEmail", result.Error); err != nil {
		return nil, err
	}

	return participants, nil
}
