package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// Force gorm import
var _ = gorm.ErrRecordNotFound

// PostgresChatUserPresenceRepository implements the repository for ChatUserPresence
type PostgresChatUserPresenceRepository struct {
	*PostgresRepository
}

// NewPostgresChatUserPresenceRepository creates a new instance
func NewPostgresChatUserPresenceRepository(base *PostgresRepository) *PostgresChatUserPresenceRepository {
	return &PostgresChatUserPresenceRepository{PostgresRepository: base}
}

// Create or update user presence (since primary key is user_id, use Upsert)
func (r *PostgresChatUserPresenceRepository) Upsert(ctx context.Context, presence *models.ChatUserPresence) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	return r.handleError("UpsertChatUserPresence", r.DB().WithContext(ctx).Save(presence).Error)
}

// GetByUserID retrieves user presence by user ID
func (r *PostgresChatUserPresenceRepository) GetByUserID(ctx context.Context, userID string) (*models.ChatUserPresence, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	var presence models.ChatUserPresence
	err := r.DB().WithContext(ctx).First(&presence, "user_id = ?", userID).Error
	if err != nil {
		return nil, r.handleError("GetChatUserPresenceByUserID", err)
	}
	return &presence, nil
}

// Delete user presence
func (r *PostgresChatUserPresenceRepository) Delete(ctx context.Context, userID string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	return r.handleError("DeleteChatUserPresence", r.DB().WithContext(ctx).Delete(&models.ChatUserPresence{}, "user_id = ?", userID).Error)
}

// ListOnlineUserIDs lists user IDs of online users
func (r *PostgresChatUserPresenceRepository) ListOnlineUsers(ctx context.Context) ([]*models.OnlineUser, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()
	var users []*models.OnlineUser
	err := r.DB().WithContext(ctx).Raw(`
		SELECT g.id, g.naam AS name
		FROM chat_user_presence p
		JOIN gebruikers g ON p.user_id = g.id
		WHERE p.status = 'online'
	`).Scan(&users).Error
	if err != nil {
		return nil, r.handleError("ListOnlineUsers", err)
	}
	return users, nil
}
