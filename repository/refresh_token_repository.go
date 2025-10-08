package repository

import (
	"context"
	"dklautomationgo/models"
	"time"
)

// RefreshTokenRepository definieert de interface voor refresh token operaties
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *models.RefreshToken) error
	GetByToken(ctx context.Context, token string) (*models.RefreshToken, error)
	RevokeToken(ctx context.Context, token string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error
	DeleteExpired(ctx context.Context) error
}

// PostgresRefreshTokenRepository implementeert RefreshTokenRepository met PostgreSQL
type PostgresRefreshTokenRepository struct {
	*PostgresRepository
}

// NewPostgresRefreshTokenRepository maakt een nieuwe PostgreSQL refresh token repository
func NewPostgresRefreshTokenRepository(base *PostgresRepository) *PostgresRefreshTokenRepository {
	return &PostgresRefreshTokenRepository{
		PostgresRepository: base,
	}
}

// Create slaat een nieuwe refresh token op
func (r *PostgresRefreshTokenRepository) Create(ctx context.Context, token *models.RefreshToken) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(token)
	return r.handleError("Create", result.Error)
}

// GetByToken haalt een refresh token op basis van token string
func (r *PostgresRefreshTokenRepository) GetByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var refreshToken models.RefreshToken
	result := r.DB().WithContext(ctx).
		Where("token = ? AND is_revoked = ? AND expires_at > ?", token, false, time.Now()).
		First(&refreshToken)

	if err := r.handleError("GetByToken", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &refreshToken, nil
}

// RevokeToken markeert een refresh token als ingetrokken
func (r *PostgresRefreshTokenRepository) RevokeToken(ctx context.Context, token string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	now := time.Now()
	result := r.DB().WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("token = ?", token).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
		})

	return r.handleError("RevokeToken", result.Error)
}

// RevokeAllUserTokens trekt alle refresh tokens van een gebruiker in
func (r *PostgresRefreshTokenRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	now := time.Now()
	result := r.DB().WithContext(ctx).
		Model(&models.RefreshToken{}).
		Where("user_id = ? AND is_revoked = ?", userID, false).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
		})

	return r.handleError("RevokeAllUserTokens", result.Error)
}

// DeleteExpired verwijdert verlopen refresh tokens
func (r *PostgresRefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.RefreshToken{})

	return r.handleError("DeleteExpired", result.Error)
}
