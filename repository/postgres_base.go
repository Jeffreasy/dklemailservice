package repository

import (
	"context"
	"dklautomationgo/logger"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// PostgresRepository bevat gemeenschappelijke functionaliteit voor alle PostgreSQL repositories
type PostgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository maakt een nieuwe PostgreSQL repository
func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// withTimeout voegt een timeout toe aan de context
func (r *PostgresRepository) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, 5*time.Second)
}

// handleError verwerkt database fouten op een consistente manier
func (r *PostgresRepository) handleError(op string, err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil // Geen record gevonden wordt als nil teruggegeven
	}

	logger.Error("Database fout", "operation", op, "error", err)
	return fmt.Errorf("%s: %w", op, err)
}

// DB geeft de database verbinding terug
func (r *PostgresRepository) DB() *gorm.DB {
	return r.db
}

// Commented out unused method - can be uncommented when needed
/*
// transaction voert een functie uit binnen een transactie
func (r *PostgresRepository) transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	ctx, cancel := r.withTimeout(ctx)
	defer cancel()

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
*/
