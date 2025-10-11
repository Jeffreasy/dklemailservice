package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// PostgresRadioRecordingRepository implements RadioRecordingRepository
type PostgresRadioRecordingRepository struct {
	db *gorm.DB
}

// NewPostgresRadioRecordingRepository creates a new radio recording repository
func NewPostgresRadioRecordingRepository(db *gorm.DB) *PostgresRadioRecordingRepository {
	return &PostgresRadioRecordingRepository{db: db}
}

// Create saves a new radio recording
func (r *PostgresRadioRecordingRepository) Create(ctx context.Context, recording *models.RadioRecording) error {
	return r.db.WithContext(ctx).Create(recording).Error
}

// GetByID retrieves a radio recording by ID
func (r *PostgresRadioRecordingRepository) GetByID(ctx context.Context, id string) (*models.RadioRecording, error) {
	var recording models.RadioRecording
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&recording).Error
	if err != nil {
		return nil, err
	}
	return &recording, nil
}

// List retrieves a paginated list of radio recordings
func (r *PostgresRadioRecordingRepository) List(ctx context.Context, limit, offset int) ([]*models.RadioRecording, error) {
	var recordings []*models.RadioRecording
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("order_number ASC, created_at DESC").Find(&recordings).Error
	return recordings, err
}

// ListVisible retrieves only visible radio recordings ordered by order_number
func (r *PostgresRadioRecordingRepository) ListVisible(ctx context.Context) ([]*models.RadioRecording, error) {
	var recordings []*models.RadioRecording
	err := r.db.WithContext(ctx).Where("visible = ?", true).Order("order_number ASC, created_at DESC").Find(&recordings).Error
	return recordings, err
}

// Update updates an existing radio recording
func (r *PostgresRadioRecordingRepository) Update(ctx context.Context, recording *models.RadioRecording) error {
	return r.db.WithContext(ctx).Save(recording).Error
}

// Delete removes a radio recording
func (r *PostgresRadioRecordingRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.RadioRecording{}, "id = ?", id).Error
}
