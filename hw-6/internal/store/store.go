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
	CreateAccount(ctx context.Context, name string) error
	GetAccounts(ctx context.Context) ([]*models.Account, error)
	GetAccount(ctx context.Context, id int) (*models.Account, error)
	DeleteAccount(ctx context.Context, id int) error
	CreateTransaction(ctx context.Context, transaction *models.Transaction, accountId int) error
	GetTransactionsByAccount(ctx context.Context, accountId int) ([]*models.Transaction, error)
	GetCurrentAmountOfAccount(ctx context.Context, accountId int) (int, error)
	UpdateTransaction(ctx context.Context, transaction *models.Transaction, accountId int) error
}
