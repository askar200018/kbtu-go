package inmemory

import (
	"context"
	"fmt"
	"hw-6/internal/models"
	"hw-6/internal/store"
	"sync"
)

type DB struct {
	data     map[int]*models.Transaction
	accounts map[int]*models.Account

	lastAccountId     int
	lastTransactionId int

	mu *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		data:     make(map[int]*models.Transaction),
		accounts: make(map[int]*models.Account),

		lastAccountId:     0,
		lastTransactionId: 0,

		mu: new(sync.RWMutex),
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

func (db *DB) CreateAccount(ctx context.Context, name string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	account := &models.Account{
		ID:            db.lastAccountId,
		Name:          name,
		Transactions:  make([]models.Transaction, 0),
		CurrentAmount: 0,
	}
	db.accounts[db.lastAccountId] = account
	db.lastAccountId++

	return nil
}

func (db *DB) GetAccounts(ctx context.Context) ([]*models.Account, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	accounts := make([]*models.Account, 0, len(db.data))
	for _, account := range db.accounts {
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (db *DB) GetAccount(ctx context.Context, id int) (*models.Account, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	account, ok := db.accounts[id]
	if !ok {
		return nil, fmt.Errorf("no transaction with id %d", id)
	}

	return account, nil
}

func (db *DB) DeleteAccount(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.accounts, id)
	return nil
}

func (db *DB) GetCurrentAmountOfAccount(ctx context.Context, accountId int) (int, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	amount := db.accounts[accountId].CurrentAmount

	return amount, nil
}

func (db *DB) CreateTransaction(
	ctx context.Context,
	transaction *models.Transaction,
	accountId int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	newTransaction := models.Transaction{
		ID:     db.lastTransactionId,
		Date:   transaction.Date,
		Amount: transaction.Amount,
		Note:   transaction.Note,
		Type:   transaction.Type,
	}
	db.lastTransactionId++

	if newTransaction.Type == "income" {
		db.accounts[accountId].CurrentAmount += newTransaction.Amount
	} else if newTransaction.Type == "expense" {
		db.accounts[accountId].CurrentAmount -= newTransaction.Amount
	}

	db.accounts[accountId].Transactions = append(db.accounts[accountId].Transactions, newTransaction)

	return nil
}

func (db *DB) GetTransactionsByAccount(ctx context.Context, accountId int) ([]*models.Transaction, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	transactions := db.accounts[accountId].Transactions
	resultTransactions := make([]*models.Transaction, 0, len(db.data))
	for _, transaction := range transactions {
		resultTransactions = append(resultTransactions, &transaction)
	}

	return resultTransactions, nil
}

func (db *DB) UpdateTransaction(
	ctx context.Context,
	transaction *models.Transaction,
	accountId int) error {

	db.mu.Lock()
	defer db.mu.Unlock()

	transactions := db.accounts[accountId].Transactions
	for index, t := range transactions {
		if t.ID == transaction.ID {
			fmt.Println("OK")
			db.accounts[accountId].CurrentAmount -= t.Amount
			db.accounts[accountId].CurrentAmount += transaction.Amount
			db.accounts[accountId].Transactions[index] = *transaction
		}
	}
	return nil
}
