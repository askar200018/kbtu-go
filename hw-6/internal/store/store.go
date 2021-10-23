package store

import (
	"context"
	"hw-6/internal/models"
)

type Store interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	All(ctx context.Context) ([]*models.Transaction, error)
	ByID(ctx context.Context, id int) (*models.Transaction, error)
	Update(ctx context.Context, transaction *models.Transaction) error
	Delete(ctx context.Context, id int) error
}
