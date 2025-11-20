package repository

import (
	"context"
	"dklautomationgo/models"
	"time"
)

// PostgresSessionRepository implementeert SessionRepository voor PostgreSQL
type PostgresSessionRepository struct {
	*PostgresRepository
}

// NewPostgresSessionRepository maakt een nieuwe PostgresSessionRepository
func NewPostgresSessionRepository(base *PostgresRepository) *PostgresSessionRepository {
	return &PostgresSessionRepository{
		PostgresRepository: base,
	}
}

// Create slaat een nieuwe sessie op
func (r *PostgresSessionRepository) Create(ctx context.Context, session *models.Session) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(session)
	return r.handleError("Create", result.Error)
}

// GetByID haalt een sessie op basis van ID
func (r *PostgresSessionRepository) GetByID(ctx context.Context, id string) (*models.Session, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var session models.Session
	result := r.DB().WithContext(ctx).Where("id = ?", id).First(&session)

	if err := r.handleError("GetByID", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &session, nil
}

// GetByAccessToken haalt een sessie op basis van access token
func (r *PostgresSessionRepository) GetByAccessToken(ctx context.Context, accessToken string) (*models.Session, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var session models.Session
	result := r.DB().WithContext(ctx).Where("access_token = ?", accessToken).First(&session)

	if err := r.handleError("GetByAccessToken", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &session, nil
}

// ListByOwnerID haalt alle actieve sessies op voor een gebruiker
func (r *PostgresSessionRepository) ListByOwnerID(ctx context.Context, ownerID string) ([]*models.Session, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var sessions []*models.Session
	result := r.DB().WithContext(ctx).
		Where("owner_id = ? AND is_active = true AND expires_at > ?", ownerID, time.Now()).
		Order("last_activity DESC").
		Find(&sessions)

	return sessions, r.handleError("ListByOwnerID", result.Error)
}

// Update werkt een bestaande sessie bij
func (r *PostgresSessionRepository) Update(ctx context.Context, session *models.Session) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Save(session)
	return r.handleError("Update", result.Error)
}

// Delete verwijdert een sessie
func (r *PostgresSessionRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Delete(&models.Session{}, "id = ?", id)
	return r.handleError("Delete", result.Error)
}

// RevokeByAccessToken trekt een sessie in op basis van access token
func (r *PostgresSessionRepository) RevokeByAccessToken(ctx context.Context, accessToken string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	now := time.Now()
	result := r.DB().WithContext(ctx).
		Model(&models.Session{}).
		Where("access_token = ?", accessToken).
		Updates(map[string]interface{}{
			"is_active":     false,
			"last_activity": now,
		})

	return r.handleError("RevokeByAccessToken", result.Error)
}

// RevokeAllUserSessions trekt alle sessies in voor een gebruiker behalve de huidige sessie
func (r *PostgresSessionRepository) RevokeAllUserSessions(ctx context.Context, ownerID, currentSessionID string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	now := time.Now()
	result := r.DB().WithContext(ctx).Model(&models.Session{}).
		Where("owner_id = ? AND id != ? AND is_active = true", ownerID, currentSessionID).
		Updates(map[string]interface{}{
			"is_active":     false,
			"last_activity": now,
		})

	return r.handleError("RevokeAllUserSessions", result.Error)
}

// RevokeAllUserSessionsComplete trekt alle sessies in voor een gebruiker
func (r *PostgresSessionRepository) RevokeAllUserSessionsComplete(ctx context.Context, ownerID string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	now := time.Now()
	result := r.DB().WithContext(ctx).Model(&models.Session{}).
		Where("owner_id = ? AND is_active = true", ownerID).
		Updates(map[string]interface{}{
			"is_active":     false,
			"last_activity": now,
		})

	return r.handleError("RevokeAllUserSessionsComplete", result.Error)
}

// MarkCurrentSession markeert een sessie als de huidige sessie
func (r *PostgresSessionRepository) MarkCurrentSession(ctx context.Context, sessionID string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	now := time.Now()

	// Eerst alle sessies voor deze gebruiker als niet-huidig markeren
	if err := r.DB().WithContext(ctx).Model(&models.Session{}).
		Where("id = ?", sessionID).
		Update("is_current", false).Error; err != nil {
		return r.handleError("MarkCurrentSession - reset", err)
	}

	// Dan deze sessie als huidig markeren
	result := r.DB().WithContext(ctx).Model(&models.Session{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"is_current":    true,
			"last_activity": now,
		})

	return r.handleError("MarkCurrentSession - set", result.Error)
}

// CleanupExpiredSessions verwijdert verlopen sessies
func (r *PostgresSessionRepository) CleanupExpiredSessions(ctx context.Context) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Unscoped().
		Where("expires_at < ? OR (is_active = false AND last_activity < ?)", time.Now(), time.Now().Add(-24*time.Hour)).
		Delete(&models.Session{})

	return r.handleError("CleanupExpiredSessions", result.Error)
}
