package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// Force gorm import
var _ = gorm.ErrRecordNotFound

// PostgresChatMessageReactionRepository implements the repository for ChatMessageReaction
type PostgresChatMessageReactionRepository struct {
	*PostgresRepository
}

// NewPostgresChatMessageReactionRepository creates a new instance
func NewPostgresChatMessageReactionRepository(base *PostgresRepository) *PostgresChatMessageReactionRepository {
	return &PostgresChatMessageReactionRepository{PostgresRepository: base}
}

// Create a new chat message reaction
func (r *PostgresChatMessageReactionRepository) Create(ctx context.Context, reaction *models.ChatMessageReaction) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	return r.handleError("CreateChatMessageReaction", r.DB().WithContext(ctx).Create(reaction).Error)
}

// GetByID retrieves a chat message reaction by ID
func (r *PostgresChatMessageReactionRepository) GetByID(ctx context.Context, id string) (*models.ChatMessageReaction, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	var reaction models.ChatMessageReaction
	err := r.DB().WithContext(ctx).First(&reaction, "id = ?", id).Error
	if err != nil {
		return nil, r.handleError("GetChatMessageReactionByID", err)
	}
	return &reaction, nil
}

// List retrieves a list of chat message reactions
func (r *PostgresChatMessageReactionRepository) List(ctx context.Context, limit, offset int) ([]*models.ChatMessageReaction, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	var reactions []*models.ChatMessageReaction
	err := r.DB().WithContext(ctx).Limit(limit).Offset(offset).Find(&reactions).Error
	if err != nil {
		return nil, r.handleError("ListChatMessageReactions", err)
	}
	return reactions, nil
}

// Delete a chat message reaction
func (r *PostgresChatMessageReactionRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	return r.handleError("DeleteChatMessageReaction", r.DB().WithContext(ctx).Delete(&models.ChatMessageReaction{}, "id = ?", id).Error)
}

// ListByMessageID retrieves reactions by message ID
func (r *PostgresChatMessageReactionRepository) ListByMessageID(ctx context.Context, messageID string) ([]*models.ChatMessageReaction, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	var reactions []*models.ChatMessageReaction
	err := r.DB().WithContext(ctx).Where("message_id = ?", messageID).Find(&reactions).Error
	if err != nil {
		return nil, r.handleError("ListByMessageID", err)
	}
	return reactions, nil
}
