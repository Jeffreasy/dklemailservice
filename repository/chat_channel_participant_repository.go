package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// Force gorm import
var _ = gorm.ErrRecordNotFound

// PostgresChatChannelParticipantRepository implements the repository for ChatChannelParticipant
type PostgresChatChannelParticipantRepository struct {
	*PostgresRepository
}

// NewPostgresChatChannelParticipantRepository creates a new instance
func NewPostgresChatChannelParticipantRepository(base *PostgresRepository) *PostgresChatChannelParticipantRepository {
	return &PostgresChatChannelParticipantRepository{PostgresRepository: base}
}

// Create a new chat channel participant
func (r *PostgresChatChannelParticipantRepository) Create(ctx context.Context, participant *models.ChatChannelParticipant) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	return r.handleError("CreateChatChannelParticipant", r.DB().WithContext(ctx).Create(participant).Error)
}

// GetByID retrieves a chat channel participant by ID
func (r *PostgresChatChannelParticipantRepository) GetByID(ctx context.Context, id string) (*models.ChatChannelParticipant, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	var participant models.ChatChannelParticipant
	err := r.DB().WithContext(ctx).First(&participant, "id = ?", id).Error
	if err != nil {
		return nil, r.handleError("GetChatChannelParticipantByID", err)
	}
	return &participant, nil
}

// List retrieves a list of chat channel participants
func (r *PostgresChatChannelParticipantRepository) List(ctx context.Context, limit, offset int) ([]*models.ChatChannelParticipant, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	var participants []*models.ChatChannelParticipant
	err := r.DB().WithContext(ctx).Limit(limit).Offset(offset).Find(&participants).Error
	if err != nil {
		return nil, r.handleError("ListChatChannelParticipants", err)
	}
	return participants, nil
}

// Update a chat channel participant
func (r *PostgresChatChannelParticipantRepository) Update(ctx context.Context, participant *models.ChatChannelParticipant) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	return r.handleError("UpdateChatChannelParticipant", r.DB().WithContext(ctx).Save(participant).Error)
}

// Delete a chat channel participant
func (r *PostgresChatChannelParticipantRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	return r.handleError("DeleteChatChannelParticipant", r.DB().WithContext(ctx).Delete(&models.ChatChannelParticipant{}, "id = ?", id).Error)
}

// Additional methods, e.g., ListByChannelID
func (r *PostgresChatChannelParticipantRepository) ListByChannelID(ctx context.Context, channelID string) ([]*models.ChatChannelParticipant, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	var participants []*models.ChatChannelParticipant
	err := r.DB().WithContext(ctx).Where("channel_id = ?", channelID).Find(&participants).Error
	if err != nil {
		return nil, r.handleError("ListByChannelID", err)
	}
	return participants, nil
}
