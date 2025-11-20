package repository

import (
	"context"
	"dklautomationgo/models"
	"time"
)

// EmailVerificationTokenRepository definieert de interface voor email verification token operaties
type EmailVerificationTokenRepository interface {
	Create(ctx context.Context, token *models.EmailVerificationToken) error
	GetByToken(ctx context.Context, token string) (*models.EmailVerificationToken, error)
	GetByEmail(ctx context.Context, email string) ([]*models.EmailVerificationToken, error)
	GetByUserID(ctx context.Context, userID string, userType string) ([]*models.EmailVerificationToken, error)
	MarkAsUsed(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
	InvalidateUserTokens(ctx context.Context, userID string, userType string) error
}

// PostgresEmailVerificationTokenRepository implementeert EmailVerificationTokenRepository met PostgreSQL
type PostgresEmailVerificationTokenRepository struct {
	*PostgresRepository
}

// NewPostgresEmailVerificationTokenRepository maakt een nieuwe PostgreSQL email verification token repository
func NewPostgresEmailVerificationTokenRepository(base *PostgresRepository) *PostgresEmailVerificationTokenRepository {
	return &PostgresEmailVerificationTokenRepository{
		PostgresRepository: base,
	}
}

// Create slaat een nieuwe email verification token op
func (r *PostgresEmailVerificationTokenRepository) Create(ctx context.Context, token *models.EmailVerificationToken) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(token)
	return r.handleError("Create", result.Error)
}

// GetByToken haalt een email verification token op basis van token string
func (r *PostgresEmailVerificationTokenRepository) GetByToken(ctx context.Context, token string) (*models.EmailVerificationToken, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var verificationToken models.EmailVerificationToken
	result := r.DB().WithContext(ctx).
		Where("token = ? AND is_used = ? AND expires_at > ?", token, false, time.Now()).
		First(&verificationToken)

	if err := r.handleError("GetByToken", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &verificationToken, nil
}

// GetByEmail haalt alle actieve email verification tokens op voor een email adres
func (r *PostgresEmailVerificationTokenRepository) GetByEmail(ctx context.Context, email string) ([]*models.EmailVerificationToken, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var tokens []*models.EmailVerificationToken
	result := r.DB().WithContext(ctx).
		Where("email = ? AND is_used = ? AND expires_at > ?", email, false, time.Now()).
		Order("created_at DESC").
		Find(&tokens)

	if err := r.handleError("GetByEmail", result.Error); err != nil {
		return nil, err
	}

	return tokens, nil
}

// GetByUserID haalt alle actieve email verification tokens op voor een specifieke gebruiker
func (r *PostgresEmailVerificationTokenRepository) GetByUserID(ctx context.Context, userID string, userType string) ([]*models.EmailVerificationToken, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var tokens []*models.EmailVerificationToken
	result := r.DB().WithContext(ctx).
		Where("user_id = ? AND user_type = ? AND is_used = ? AND expires_at > ?", userID, userType, false, time.Now()).
		Order("created_at DESC").
		Find(&tokens)

	if err := r.handleError("GetByUserID", result.Error); err != nil {
		return nil, err
	}

	return tokens, nil
}

// MarkAsUsed markeert een email verification token als gebruikt
func (r *PostgresEmailVerificationTokenRepository) MarkAsUsed(ctx context.Context, token string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	now := time.Now()
	result := r.DB().WithContext(ctx).
		Model(&models.EmailVerificationToken{}).
		Where("token = ?", token).
		Updates(map[string]interface{}{
			"is_used": true,
			"used_at": now,
		})

	return r.handleError("MarkAsUsed", result.Error)
}

// DeleteExpired verwijdert verlopen email verification tokens
func (r *PostgresEmailVerificationTokenRepository) DeleteExpired(ctx context.Context) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.EmailVerificationToken{})

	return r.handleError("DeleteExpired", result.Error)
}

// InvalidateUserTokens maakt alle tokens voor een gebruiker ongeldig
func (r *PostgresEmailVerificationTokenRepository) InvalidateUserTokens(ctx context.Context, userID string, userType string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	now := time.Now()
	result := r.DB().WithContext(ctx).
		Model(&models.EmailVerificationToken{}).
		Where("user_id = ? AND user_type = ? AND is_used = ?", userID, userType, false).
		Updates(map[string]interface{}{
			"is_used": true,
			"used_at": now,
		})

	return r.handleError("InvalidateUserTokens", result.Error)
}
