package repository

import (
	"context"
	"dklautomationgo/models"

	"gorm.io/gorm"
)

// PostgresProgramScheduleRepository implements ProgramScheduleRepository
type PostgresProgramScheduleRepository struct {
	db *gorm.DB
}

// NewPostgresProgramScheduleRepository creates a new program schedule repository
func NewPostgresProgramScheduleRepository(db *gorm.DB) *PostgresProgramScheduleRepository {
	return &PostgresProgramScheduleRepository{db: db}
}

// Create saves a new program schedule
func (r *PostgresProgramScheduleRepository) Create(ctx context.Context, schedule *models.ProgramSchedule) error {
	return r.db.WithContext(ctx).Create(schedule).Error
}

// GetByID retrieves a program schedule by ID
func (r *PostgresProgramScheduleRepository) GetByID(ctx context.Context, id string) (*models.ProgramSchedule, error) {
	var schedule models.ProgramSchedule
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

// List retrieves a paginated list of program schedules
func (r *PostgresProgramScheduleRepository) List(ctx context.Context, limit, offset int) ([]*models.ProgramSchedule, error) {
	var schedules []*models.ProgramSchedule
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("order_number ASC, created_at DESC").Find(&schedules).Error
	return schedules, err
}

// ListVisible retrieves only visible program schedules ordered by order_number
func (r *PostgresProgramScheduleRepository) ListVisible(ctx context.Context) ([]*models.ProgramSchedule, error) {
	var schedules []*models.ProgramSchedule
	err := r.db.WithContext(ctx).Where("visible = ?", true).Order("order_number ASC, created_at DESC").Find(&schedules).Error
	return schedules, err
}

// Update updates an existing program schedule
func (r *PostgresProgramScheduleRepository) Update(ctx context.Context, schedule *models.ProgramSchedule) error {
	return r.db.WithContext(ctx).Save(schedule).Error
}

// Delete removes a program schedule
func (r *PostgresProgramScheduleRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.ProgramSchedule{}, "id = ?", id).Error
}
