package repository

import (
	"context"
	"dklautomationgo/models"
)

// EventRepository interface voor event operaties
type EventRepository interface {
	Create(ctx context.Context, event *models.Event) error
	GetByID(ctx context.Context, id string) (*models.Event, error)
	GetActiveEvent(ctx context.Context) (*models.Event, error)
	List(ctx context.Context, limit, offset int) ([]*models.Event, error)
	ListActive(ctx context.Context) ([]*models.Event, error)
	Update(ctx context.Context, event *models.Event) error
	Delete(ctx context.Context, id string) error

	// Event Participant methods
	RegisterParticipant(ctx context.Context, eventID, participantID string) (*models.EventParticipant, error)
	GetEventParticipant(ctx context.Context, eventID, participantID string) (*models.EventParticipant, error)
	UpdateEventParticipant(ctx context.Context, ep *models.EventParticipant) error
	GetEventParticipants(ctx context.Context, eventID string) ([]*models.EventParticipant, error)
	GetParticipantEvents(ctx context.Context, participantID string) ([]*models.EventParticipant, error)
}

// PostgresEventRepository implementeert EventRepository met PostgreSQL
type PostgresEventRepository struct {
	*PostgresRepository
}

// NewPostgresEventRepository maakt een nieuwe PostgreSQL event repository
func NewPostgresEventRepository(base *PostgresRepository) *PostgresEventRepository {
	return &PostgresEventRepository{
		PostgresRepository: base,
	}
}

// Create slaat een nieuw event op
func (r *PostgresEventRepository) Create(ctx context.Context, event *models.Event) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(event)
	return r.handleError("Create", result.Error)
}

// GetByID haalt een event op basis van ID
func (r *PostgresEventRepository) GetByID(ctx context.Context, id string) (*models.Event, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var event models.Event
	result := r.DB().WithContext(ctx).First(&event, "id = ?", id)
	if err := r.handleError("GetByID", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &event, nil
}

// GetActiveEvent haalt het eerstvolgende actieve event op
func (r *PostgresEventRepository) GetActiveEvent(ctx context.Context) (*models.Event, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var event models.Event
	result := r.DB().WithContext(ctx).
		Where("is_active = ? AND status IN (?)", true, []string{models.EventStatusUpcoming, models.EventStatusActive}).
		Order("start_time ASC").
		First(&event)

	if err := r.handleError("GetActiveEvent", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &event, nil
}

// List haalt een lijst van events op
func (r *PostgresEventRepository) List(ctx context.Context, limit, offset int) ([]*models.Event, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var events []*models.Event
	result := r.DB().WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("start_time DESC").
		Find(&events)

	if err := r.handleError("List", result.Error); err != nil {
		return nil, err
	}

	return events, nil
}

// ListActive haalt alleen actieve events op
func (r *PostgresEventRepository) ListActive(ctx context.Context) ([]*models.Event, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var events []*models.Event
	result := r.DB().WithContext(ctx).
		Where("is_active = ?", true).
		Order("start_time DESC").
		Find(&events)

	if err := r.handleError("ListActive", result.Error); err != nil {
		return nil, err
	}

	return events, nil
}

// Update werkt een bestaand event bij
func (r *PostgresEventRepository) Update(ctx context.Context, event *models.Event) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Save(event)
	return r.handleError("Update", result.Error)
}

// Delete verwijdert een event
func (r *PostgresEventRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Delete(&models.Event{}, "id = ?", id)
	return r.handleError("Delete", result.Error)
}

// RegisterParticipant registreert een participant voor een event
func (r *PostgresEventRepository) RegisterParticipant(ctx context.Context, eventID, participantID string) (*models.EventParticipant, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	ep := &models.EventParticipant{
		EventID:        eventID,
		ParticipantID:  participantID,
		TrackingStatus: models.TrackingStatusRegistered,
	}

	result := r.DB().WithContext(ctx).Create(ep)
	if err := r.handleError("RegisterParticipant", result.Error); err != nil {
		return nil, err
	}

	return ep, nil
}

// GetEventParticipant haalt een specifieke event participant op
func (r *PostgresEventRepository) GetEventParticipant(ctx context.Context, eventID, participantID string) (*models.EventParticipant, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var ep models.EventParticipant
	result := r.DB().WithContext(ctx).
		Preload("Event").
		Preload("Participant").
		Where("event_id = ? AND participant_id = ?", eventID, participantID).
		First(&ep)

	if err := r.handleError("GetEventParticipant", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &ep, nil
}

// UpdateEventParticipant werkt event participant gegevens bij
func (r *PostgresEventRepository) UpdateEventParticipant(ctx context.Context, ep *models.EventParticipant) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Save(ep)
	return r.handleError("UpdateEventParticipant", result.Error)
}

// GetEventParticipants haalt alle participants voor een event op
func (r *PostgresEventRepository) GetEventParticipants(ctx context.Context, eventID string) ([]*models.EventParticipant, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var participants []*models.EventParticipant
	result := r.DB().WithContext(ctx).
		Preload("Participant").
		Where("event_id = ?", eventID).
		Order("registered_at DESC").
		Find(&participants)

	if err := r.handleError("GetEventParticipants", result.Error); err != nil {
		return nil, err
	}

	return participants, nil
}

// GetParticipantEvents haalt alle events voor een participant op
func (r *PostgresEventRepository) GetParticipantEvents(ctx context.Context, participantID string) ([]*models.EventParticipant, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var events []*models.EventParticipant
	result := r.DB().WithContext(ctx).
		Preload("Event").
		Where("participant_id = ?", participantID).
		Order("registered_at DESC").
		Find(&events)

	if err := r.handleError("GetParticipantEvents", result.Error); err != nil {
		return nil, err
	}

	return events, nil
}
