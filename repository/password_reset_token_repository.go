package repository

import (
	"context"
	"dklautomationgo/models"
	"time"
)

// PasswordResetTokenRepository definieert de interface voor password reset token operaties
type PasswordResetTokenRepository interface {
	Create(ctx context.Context, token *models.PasswordResetToken) error
	GetByToken(ctx context.Context, token string) (*models.PasswordResetToken, error)
	GetByEmail(ctx context.Context, email string) ([]*models.PasswordResetToken, error)
	MarkAsUsed(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context) error
}

// PostgresPasswordResetTokenRepository implementeert PasswordResetTokenRepository met PostgreSQL
type PostgresPasswordResetTokenRepository struct {
	*PostgresRepository
}

// NewPostgresPasswordResetTokenRepository maakt een nieuwe PostgreSQL password reset token repository
func NewPostgresPasswordResetTokenRepository(base *PostgresRepository) *PostgresPasswordResetTokenRepository {
	return &PostgresPasswordResetTokenRepository{
		PostgresRepository: base,
	}
}

// Create slaat een nieuwe password reset token op
func (r *PostgresPasswordResetTokenRepository) Create(ctx context.Context, token *models.PasswordResetToken) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(token)
	return r.handleError("Create", result.Error)
}

// GetByToken haalt een password reset token op basis van token string
func (r *PostgresPasswordResetTokenRepository) GetByToken(ctx context.Context, token string) (*models.PasswordResetToken, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var resetToken models.PasswordResetToken
	result := r.DB().WithContext(ctx).
		Where("token = ? AND is_used = ? AND expires_at > ?", token, false, time.Now()).
		First(&resetToken)

	if err := r.handleError("GetByToken", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &resetToken, nil
}

// GetByEmail haalt alle actieve password reset tokens op voor een email adres
func (r *PostgresPasswordResetTokenRepository) GetByEmail(ctx context.Context, email string) ([]*models.PasswordResetToken, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var tokens []*models.PasswordResetToken
	result := r.DB().WithContext(ctx).
		Where("email = ? AND is_used = ? AND expires_at > ?", email, false, time.Now()).
		Order("created_at DESC").
		Find(&tokens)

	if err := r.handleError("GetByEmail", result.Error); err != nil {
		return nil, err
	}

	return tokens, nil
}

// MarkAsUsed markeert een password reset token als gebruikt
func (r *PostgresPasswordResetTokenRepository) MarkAsUsed(ctx context.Context, token string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	now := time.Now()
	result := r.DB().WithContext(ctx).
		Model(&models.PasswordResetToken{}).
		Where("token = ?", token).
		Updates(map[string]interface{}{
			"is_used": true,
			"used_at": now,
		})

	return r.handleError("MarkAsUsed", result.Error)
}

// DeleteExpired verwijdert verlopen password reset tokens
func (r *PostgresPasswordResetTokenRepository) DeleteExpired(ctx context.Context) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.PasswordResetToken{})

	return r.handleError("DeleteExpired", result.Error)
}
