package repository

import (
	"context"
	"dklautomationgo/models"
	"time"
)

// PostgresAccessTokenRepository implementeert AccessTokenRepository met PostgreSQL
type PostgresAccessTokenRepository struct {
	*PostgresRepository
}

// NewPostgresAccessTokenRepository maakt een nieuwe PostgreSQL access token repository
func NewPostgresAccessTokenRepository(base *PostgresRepository) *PostgresAccessTokenRepository {
	return &PostgresAccessTokenRepository{
		PostgresRepository: base,
	}
}

// Create slaat een nieuwe access token op
func (r *PostgresAccessTokenRepository) Create(ctx context.Context, token *models.AccessToken) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(token)
	return r.handleError("Create", result.Error)
}

// GetByToken haalt een access token op basis van token string
func (r *PostgresAccessTokenRepository) GetByToken(ctx context.Context, token string) (*models.AccessToken, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var accessToken models.AccessToken
	result := r.DB().WithContext(ctx).
		Where("token = ? AND is_revoked = ? AND expires_at > ?", token, false, time.Now()).
		First(&accessToken)

	if err := r.handleError("GetByToken", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &accessToken, nil
}

// RevokeToken markeert een access token als ingetrokken
func (r *PostgresAccessTokenRepository) RevokeToken(ctx context.Context, token string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	now := time.Now()
	result := r.DB().WithContext(ctx).
		Model(&models.AccessToken{}).
		Where("token = ?", token).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
		})

	return r.handleError("RevokeToken", result.Error)
}

// RevokeAllUserTokens trekt alle access tokens van een gebruiker in
func (r *PostgresAccessTokenRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	now := time.Now()
	result := r.DB().WithContext(ctx).
		Model(&models.AccessToken{}).
		Where("owner_id = ? AND is_revoked = ?", userID, false).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
		})

	return r.handleError("RevokeAllUserTokens", result.Error)
}

// UpdateSessionID werkt de session ID bij voor een access token
func (r *PostgresAccessTokenRepository) UpdateSessionID(ctx context.Context, token, sessionID string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).
		Model(&models.AccessToken{}).
		Where("token = ?", token).
		Update("session_id", sessionID)

	return r.handleError("UpdateSessionID", result.Error)
}

// DeleteExpired verwijdert verlopen access tokens
func (r *PostgresAccessTokenRepository) DeleteExpired(ctx context.Context) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.AccessToken{})

	return r.handleError("DeleteExpired", result.Error)
}
