package inmemory

import (
	"context"
	"fmt"
	"hw-6/internal/models"
	"hw-6/internal/store"
	"sync"
)

type DB struct {
	data map[int]*models.Transaction

	mu *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		data: make(map[int]*models.Transaction),
		mu:   new(sync.RWMutex),
	}
}

func (db *DB) Create(ctx context.Context, transaction *models.Transaction) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[transaction.ID] = transaction
	return nil
}

func (db *DB) All(ctx context.Context) ([]*models.Transaction, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	transactions := make([]*models.Transaction, 0, len(db.data))
	for _, transaction := range db.data {
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (db *DB) ByID(ctx context.Context, id int) (*models.Transaction, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	transaction, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("no transaction with id %d", id)
	}

	return transaction, nil
}

func (db *DB) Update(ctx context.Context, transaction *models.Transaction) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[transaction.ID] = transaction
	return nil
}

func (db *DB) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}
