package repository

import (
	"context"
	"dklautomationgo/logger"
	"dklautomationgo/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostgresIncomingEmailRepository implementeert IncomingEmailRepository met PostgreSQL
type PostgresIncomingEmailRepository struct {
	db *gorm.DB
}

// NewPostgresIncomingEmailRepository maakt een nieuw PostgresIncomingEmailRepository
func NewPostgresIncomingEmailRepository(db *gorm.DB) *PostgresIncomingEmailRepository {
	return &PostgresIncomingEmailRepository{
		db: db,
	}
}

// Create slaat een nieuwe inkomende e-mail op
func (r *PostgresIncomingEmailRepository) Create(ctx context.Context, email *models.IncomingEmail) error {
	if email.ID == "" {
		email.ID = uuid.New().String()
	}

	err := r.db.WithContext(ctx).Create(email).Error
	if err != nil {
		logger.Error("Fout bij opslaan inkomende e-mail", "error", err)
		return err
	}

	return nil
}

// GetByID haalt een inkomende e-mail op basis van ID
func (r *PostgresIncomingEmailRepository) GetByID(ctx context.Context, id string) (*models.IncomingEmail, error) {
	var email models.IncomingEmail

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		logger.Error("Fout bij ophalen inkomende e-mail", "error", err, "id", id)
		return nil, err
	}

	return &email, nil
}

// List haalt een lijst van inkomende e-mails op
func (r *PostgresIncomingEmailRepository) List(ctx context.Context, limit, offset int) ([]*models.IncomingEmail, error) {
	var emails []*models.IncomingEmail

	err := r.db.WithContext(ctx).
		Order("received_at desc").
		Limit(limit).
		Offset(offset).
		Find(&emails).Error

	if err != nil {
		logger.Error("Fout bij ophalen inkomende e-mails", "error", err)
		return nil, err
	}

	return emails, nil
}

// Update werkt een bestaande inkomende e-mail bij
func (r *PostgresIncomingEmailRepository) Update(ctx context.Context, email *models.IncomingEmail) error {
	err := r.db.WithContext(ctx).Save(email).Error
	if err != nil {
		logger.Error("Fout bij bijwerken inkomende e-mail", "error", err, "id", email.ID)
		return err
	}

	return nil
}

// Delete verwijdert een inkomende e-mail
func (r *PostgresIncomingEmailRepository) Delete(ctx context.Context, id string) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.IncomingEmail{}).Error
	if err != nil {
		logger.Error("Fout bij verwijderen inkomende e-mail", "error", err, "id", id)
		return err
	}

	return nil
}

// FindByUID zoekt een inkomende e-mail op basis van UID
func (r *PostgresIncomingEmailRepository) FindByUID(ctx context.Context, uid string) (*models.IncomingEmail, error) {
	var email models.IncomingEmail

	err := r.db.WithContext(ctx).Where("uid = ?", uid).First(&email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		logger.Error("Fout bij zoeken inkomende e-mail op UID", "error", err, "uid", uid)
		return nil, err
	}

	return &email, nil
}

// FindUnprocessed haalt alle onverwerkte e-mails op
func (r *PostgresIncomingEmailRepository) FindUnprocessed(ctx context.Context) ([]*models.IncomingEmail, error) {
	var emails []*models.IncomingEmail

	err := r.db.WithContext(ctx).
		Where("is_processed = ?", false).
		Order("received_at asc").
		Find(&emails).Error

	if err != nil {
		logger.Error("Fout bij ophalen onverwerkte e-mails", "error", err)
		return nil, err
	}

	return emails, nil
}

// FindByAccountType zoekt inkomende e-mails op basis van account type
func (r *PostgresIncomingEmailRepository) FindByAccountType(ctx context.Context, accountType string) ([]*models.IncomingEmail, error) {
	var emails []*models.IncomingEmail

	err := r.db.WithContext(ctx).
		Where("account_type = ?", accountType).
		Order("received_at desc").
		Find(&emails).Error

	if err != nil {
		logger.Error("Fout bij zoeken inkomende e-mails op account type", "error", err, "account_type", accountType)
		return nil, err
	}

	return emails, nil
}

// ListByAccountTypePaginated haalt een gepagineerde lijst van inkomende e-mails op basis van account type en retourneert het totaal aantal
func (r *PostgresIncomingEmailRepository) ListByAccountTypePaginated(ctx context.Context, accountType string, limit, offset int) ([]*models.IncomingEmail, int64, error) {
	var emails []*models.IncomingEmail
	var totalCount int64

	// Begin een transactie
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		logger.Error("Fout bij starten transactie voor paginated email fetch", "error", tx.Error)
		return nil, 0, tx.Error
	}

	// Query om het totaal aantal te tellen
	countQuery := tx.Model(&models.IncomingEmail{}).Where("account_type = ?", accountType)
	if err := countQuery.Count(&totalCount).Error; err != nil {
		tx.Rollback() // Rollback bij fout
		logger.Error("Fout bij tellen inkomende e-mails op account type", "error", err, "account_type", accountType)
		return nil, 0, err
	}

	// Query om de gepagineerde data op te halen
	dataQuery := tx.Model(&models.IncomingEmail{}).
		Where("account_type = ?", accountType).
		Order("received_at desc").
		Limit(limit).
		Offset(offset)

	if err := dataQuery.Find(&emails).Error; err != nil {
		tx.Rollback() // Rollback bij fout
		logger.Error("Fout bij ophalen gepagineerde inkomende e-mails op account type", "error", err, "account_type", accountType)
		return nil, 0, err
	}

	// Commit de transactie
	if err := tx.Commit().Error; err != nil {
		logger.Error("Fout bij commiten transactie voor paginated email fetch", "error", err)
		// Data is al opgehaald, maar de transactie kon niet worden gecommit.
		// We retourneren de data toch maar loggen de fout.
		return emails, totalCount, nil
	}

	return emails, totalCount, nil
}
