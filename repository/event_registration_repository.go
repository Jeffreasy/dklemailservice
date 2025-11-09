package repository

import (
	"context"
	"dklautomationgo/models"
)

// PostgresEventRegistrationRepository implementeert EventRegistrationRepository.
type PostgresEventRegistrationRepository struct {
	*PostgresRepository
}

// NewPostgresEventRegistrationRepository maakt een nieuwe repository.
func NewPostgresEventRegistrationRepository(base *PostgresRepository) *PostgresEventRegistrationRepository {
	return &PostgresEventRegistrationRepository{
		PostgresRepository: base,
	}
}

// Create slaat een nieuwe event registratie op
func (r *PostgresEventRegistrationRepository) Create(ctx context.Context, registration *models.EventRegistration) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(registration)
	return r.handleError("Create", result.Error)
}

// GetByID haalt één registratie op basis van ID.
func (r *PostgresEventRegistrationRepository) GetByID(ctx context.Context, id string) (*models.EventRegistration, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var registration models.EventRegistration
	// Preload (automatisch ophalen) van gekoppelde data
	result := r.DB().WithContext(ctx).
		Preload("Participant").     // Laadt de Participant-gegevens
		Preload("Event").           // Laadt de Event-gegevens
		Preload("Distance").        // Laadt de Afstand-details
		Preload("ParticipantRole"). // Laadt de Rol-details
		First(&registration, "id = ?", id)

	if err := r.handleError("GetByID", result.Error); err != nil {
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, nil // Niet gevonden
	}
	return &registration, nil
}

// ListByParticipantID haalt alle registraties voor een specifieke participant op.
func (r *PostgresEventRegistrationRepository) ListByParticipantID(ctx context.Context, participantID string) ([]*models.EventRegistration, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var registrations []*models.EventRegistration
	result := r.DB().WithContext(ctx).
		Preload("Event").
		Preload("Distance").
		Preload("ParticipantRole").
		Where("participant_id = ?", participantID).
		Order("registered_at DESC").
		Find(&registrations)

	if err := r.handleError("ListByParticipantID", result.Error); err != nil {
		return nil, err
	}
	return registrations, nil
}

// GetActiveRegistrationForParticipant haalt de meest relevante (vaak laatste) registratie op.
// Een "actief" event wordt gedefinieerd als een event dat nog niet is voltooid (status != 'completed' en status != 'cancelled')
func (r *PostgresEventRegistrationRepository) GetActiveRegistrationForParticipant(ctx context.Context, participantID string) (*models.EventRegistration, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var registration models.EventRegistration
	// Haal de laatste registratie voor een actief event (niet voltooid of geannuleerd)
	result := r.DB().WithContext(ctx).
		Joins("JOIN events ON event_registrations.event_id = events.id").
		Where("event_registrations.participant_id = ? AND events.status NOT IN ('completed', 'cancelled')", participantID).
		Preload("Event").
		Preload("Distance").
		Preload("ParticipantRole").
		Order("event_registrations.registered_at DESC").
		First(&registration)

	if err := r.handleError("GetActiveRegistrationForParticipant", result.Error); err != nil {
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, nil // Niet gevonden
	}
	return &registration, nil
}

// ListByRole haalt alle registraties op met een specifieke rol.
// Dit is de vervanging voor de oude GetAanmeldingenByRol.
func (r *PostgresEventRegistrationRepository) ListByRole(ctx context.Context, roleName string) ([]*models.EventRegistration, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var registrations []*models.EventRegistration
	result := r.DB().WithContext(ctx).
		Preload("Participant").
		Preload("Event").
		Where("participant_role_name = ?", roleName).
		Order("registered_at DESC").
		Find(&registrations)

	if err := r.handleError("ListByRole", result.Error); err != nil {
		return nil, err
	}
	return registrations, nil
}

// Update slaat wijzigingen in een registratie op.
func (r *PostgresEventRegistrationRepository) Update(ctx context.Context, registration *models.EventRegistration) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Save(registration)
	return r.handleError("Update", result.Error)
}

// List haalt een lijst van event registraties op
func (r *PostgresEventRegistrationRepository) List(ctx context.Context, limit, offset int) ([]*models.EventRegistration, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var registrations []*models.EventRegistration
	result := r.DB().WithContext(ctx).
		Preload("Event").
		Preload("Distance").
		Preload("ParticipantRole").
		Limit(limit).
		Offset(offset).
		Order("registered_at DESC").
		Find(&registrations)

	if err := r.handleError("List", result.Error); err != nil {
		return nil, err
	}
	return registrations, nil
}

// Delete verwijdert een event registratie
func (r *PostgresEventRegistrationRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Delete(&models.EventRegistration{}, "id = ?", id)
	return r.handleError("Delete", result.Error)
}

// ListByEventID haalt alle registraties voor een specifiek event op
func (r *PostgresEventRegistrationRepository) ListByEventID(ctx context.Context, eventID string) ([]*models.EventRegistration, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var registrations []*models.EventRegistration
	result := r.DB().WithContext(ctx).
		Preload("Participant").
		Preload("Distance").
		Preload("ParticipantRole").
		Where("event_id = ?", eventID).
		Order("registered_at DESC").
		Find(&registrations)

	if err := r.handleError("ListByEventID", result.Error); err != nil {
		return nil, err
	}
	return registrations, nil
}

// ListByStatus haalt alle registraties met een specifieke status op
func (r *PostgresEventRegistrationRepository) ListByStatus(ctx context.Context, status string) ([]*models.EventRegistration, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var registrations []*models.EventRegistration
	result := r.DB().WithContext(ctx).
		Preload("Participant").
		Preload("Event").
		Preload("Distance").
		Preload("ParticipantRole").
		Where("tracking_status = ?", status).
		Order("registered_at DESC").
		Find(&registrations)

	if err := r.handleError("ListByStatus", result.Error); err != nil {
		return nil, err
	}
	return registrations, nil
}

// GetByEventAndParticipant haalt een specifieke registratie op voor een event-participant combinatie
func (r *PostgresEventRegistrationRepository) GetByEventAndParticipant(ctx context.Context, eventID, participantID string) (*models.EventRegistration, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var registration models.EventRegistration
	result := r.DB().WithContext(ctx).
		Preload("Participant").
		Preload("Event").
		Preload("Distance").
		Preload("ParticipantRole").
		Where("event_id = ? AND participant_id = ?", eventID, participantID).
		First(&registration)

	if err := r.handleError("GetByEventAndParticipant", result.Error); err != nil {
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, nil // Niet gevonden
	}
	return &registration, nil
}

// UpdateStatus werkt alleen de status van een registratie bij
func (r *PostgresEventRegistrationRepository) UpdateStatus(ctx context.Context, id, status string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).
		Model(&models.EventRegistration{}).
		Where("id = ?", id).
		Update("tracking_status", status)
	return r.handleError("UpdateStatus", result.Error)
}
