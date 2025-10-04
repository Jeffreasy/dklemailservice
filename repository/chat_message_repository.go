package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// Force gorm import
var _ = gorm.ErrRecordNotFound

// PostgresChatMessageRepository implements the repository for ChatMessage
type PostgresChatMessageRepository struct {
	*PostgresRepository
}

// NewPostgresChatMessageRepository creates a new instance
func NewPostgresChatMessageRepository(base *PostgresRepository) *PostgresChatMessageRepository {
	return &PostgresChatMessageRepository{PostgresRepository: base}
}

// Create a new chat message
func (r *PostgresChatMessageRepository) Create(ctx context.Context, message *models.ChatMessage) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	return r.handleError("CreateChatMessage", r.DB().WithContext(ctx).Create(message).Error)
}

// GetByID retrieves a chat message by ID
func (r *PostgresChatMessageRepository) GetByID(ctx context.Context, id string) (*models.ChatMessage, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	var message models.ChatMessage
	err := r.DB().WithContext(ctx).First(&message, "id = ?", id).Error
	if err != nil {
		return nil, r.handleError("GetChatMessageByID", err)
	}
	return &message, nil
}

// List retrieves a list of chat messages
func (r *PostgresChatMessageRepository) List(ctx context.Context, limit, offset int) ([]*models.ChatMessage, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	var messages []*models.ChatMessage
	err := r.DB().WithContext(ctx).Limit(limit).Offset(offset).Find(&messages).Error
	if err != nil {
		return nil, r.handleError("ListChatMessages", err)
	}
	return messages, nil
}

// Update a chat message
func (r *PostgresChatMessageRepository) Update(ctx context.Context, message *models.ChatMessage) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	return r.handleError("UpdateChatMessage", r.DB().WithContext(ctx).Save(message).Error)
}

// Delete a chat message
func (r *PostgresChatMessageRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	return r.handleError("DeleteChatMessage", r.DB().WithContext(ctx).Delete(&models.ChatMessage{}, "id = ?", id).Error)
}

// ListByChannelID retrieves messages by channel ID with pagination
func (r *PostgresChatMessageRepository) ListByChannelID(ctx context.Context, channelID string, limit, offset int) ([]*models.ChatMessage, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	var messages []*models.ChatMessage
	err := r.DB().WithContext(ctx).Where("channel_id = ?", channelID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&messages).Error
	if err != nil {
		return nil, r.handleError("ListByChannelID", err)
	}
	return messages, nil
}
