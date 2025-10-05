package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// PostgresChatChannelRepository implements the repository for ChatChannel
type PostgresChatChannelRepository struct {
	*PostgresRepository
}

// Force gorm import
type gormDB = *gorm.DB

var _ gormDB

var _ = gorm.ErrRecordNotFound

// NewPostgresChatChannelRepository creates a new instance
func NewPostgresChatChannelRepository(base *PostgresRepository) *PostgresChatChannelRepository {
	return &PostgresChatChannelRepository{PostgresRepository: base}
}

// Create a new chat channel
func (r *PostgresChatChannelRepository) Create(ctx context.Context, channel *models.ChatChannel) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	return r.handleError("CreateChatChannel", r.DB().WithContext(ctx).Create(channel).Error)
}

// GetByID retrieves a chat channel by ID
func (r *PostgresChatChannelRepository) GetByID(ctx context.Context, id string) (*models.ChatChannel, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	var channel models.ChatChannel
	err := r.DB().WithContext(ctx).First(&channel, "id = ?", id).Error
	if err != nil {
		return nil, r.handleError("GetChatChannelByID", err)
	}
	return &channel, nil
}

// List retrieves a list of chat channels
func (r *PostgresChatChannelRepository) List(ctx context.Context, limit, offset int) ([]*models.ChatChannel, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	var channels []*models.ChatChannel
	err := r.DB().WithContext(ctx).Limit(limit).Offset(offset).Find(&channels).Error
	if err != nil {
		return nil, r.handleError("ListChatChannels", err)
	}
	return channels, nil
}

// Update a chat channel
func (r *PostgresChatChannelRepository) Update(ctx context.Context, channel *models.ChatChannel) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	return r.handleError("UpdateChatChannel", r.DB().WithContext(ctx).Save(channel).Error)
}

// Delete a chat channel
func (r *PostgresChatChannelRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	return r.handleError("DeleteChatChannel", r.DB().WithContext(ctx).Delete(&models.ChatChannel{}, "id = ?", id).Error)
}

// ListByUserID lists channels for a specific user
func (r *PostgresChatChannelRepository) ListByUserID(ctx context.Context, userID string, limit, offset int) ([]*models.ChatChannel, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	var channels []*models.ChatChannel
	err := r.DB().WithContext(ctx).
		Joins("JOIN chat_channel_participants ON chat_channel_participants.channel_id = chat_channels.id").
		Where("chat_channel_participants.user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Find(&channels).Error
	if err != nil {
		return nil, r.handleError("ListChannelsByUserID", err)
	}
	return channels, nil
}

// ListPublicChannels lists all public channels
func (r *PostgresChatChannelRepository) ListPublicChannels(ctx context.Context) ([]*models.ChatChannel, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	var channels []*models.ChatChannel
	err := r.DB().WithContext(ctx).Where("is_public = ?", true).Find(&channels).Error
	if err != nil {
		return nil, r.handleError("ListPublicChannels", err)
	}
	return channels, nil
}
