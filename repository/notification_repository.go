package repository

import (
	"context"
	"dklautomationgo/models"
)

// PostgresNotificationRepository implementeert NotificationRepository met PostgreSQL
type PostgresNotificationRepository struct {
	*PostgresRepository
}

// NewPostgresNotificationRepository maakt een nieuwe PostgreSQL notification repository
func NewPostgresNotificationRepository(base *PostgresRepository) *PostgresNotificationRepository {
	return &PostgresNotificationRepository{
		PostgresRepository: base,
	}
}

// Create slaat een nieuwe notificatie op
func (r *PostgresNotificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Create(notification)
	return r.handleError("Create", result.Error)
}

// GetByID haalt een notificatie op basis van ID
func (r *PostgresNotificationRepository) GetByID(ctx context.Context, id string) (*models.Notification, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var notification models.Notification
	result := r.DB().WithContext(ctx).First(&notification, "id = ?", id)
	if err := r.handleError("GetByID", result.Error); err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &notification, nil
}

// Update werkt een bestaande notificatie bij
func (r *PostgresNotificationRepository) Update(ctx context.Context, notification *models.Notification) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Save(notification)
	return r.handleError("Update", result.Error)
}

// Delete verwijdert een notificatie
func (r *PostgresNotificationRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	result := r.DB().WithContext(ctx).Delete(&models.Notification{}, "id = ?", id)
	return r.handleError("Delete", result.Error)
}

// ListUnsent haalt alle niet verzonden notificaties op
func (r *PostgresNotificationRepository) ListUnsent(ctx context.Context) ([]*models.Notification, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var notifications []*models.Notification
	result := r.DB().WithContext(ctx).
		Where("sent = ?", false).
		Order("priority DESC, created_at ASC").
		Find(&notifications)

	if err := r.handleError("ListUnsent", result.Error); err != nil {
		return nil, err
	}

	return notifications, nil
}

// ListByType haalt alle notificaties op van een bepaald type
func (r *PostgresNotificationRepository) ListByType(ctx context.Context, notificationType models.NotificationType) ([]*models.Notification, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var notifications []*models.Notification
	result := r.DB().WithContext(ctx).
		Where("type = ?", notificationType).
		Order("created_at DESC").
		Find(&notifications)

	if err := r.handleError("ListByType", result.Error); err != nil {
		return nil, err
	}

	return notifications, nil
}

// ListByPriority haalt alle notificaties op met een bepaalde prioriteit
func (r *PostgresNotificationRepository) ListByPriority(ctx context.Context, priority models.NotificationPriority) ([]*models.Notification, error) {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	var notifications []*models.Notification
	result := r.DB().WithContext(ctx).
		Where("priority = ?", priority).
		Order("created_at DESC").
		Find(&notifications)

	if err := r.handleError("ListByPriority", result.Error); err != nil {
		return nil, err
	}

	return notifications, nil
}
