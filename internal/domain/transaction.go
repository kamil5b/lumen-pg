package domain

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrNotImplemented          = errors.New("Haven't implemented yet")
	ErrTransactionAlreadyActive = errors.New("transaction already active")
	ErrNoActiveTransaction     = errors.New("no active transaction")
	ErrTransactionExpired      = errors.New("transaction expired")
)

// OperationType represents the type of buffered operation.
type OperationType string

const (
	OpInsert OperationType = "INSERT"
	OpUpdate OperationType = "UPDATE"
	OpDelete OperationType = "DELETE"
)

// BufferedOperation represents an operation buffered within a transaction.
type BufferedOperation struct {
	Type      OperationType
	Table     string
	Schema    string
	Database  string
	RowData   map[string]interface{}
	WhereData map[string]interface{}
}

// Transaction represents a user's active transaction with a timer.
type Transaction struct {
	mu         sync.Mutex
	Active     bool
	Operations []BufferedOperation
	StartedAt  time.Time
	ExpiresAt  time.Time
	Duration   time.Duration
}

// NewTransaction creates a new transaction with a 1-minute timer.
func NewTransaction() *Transaction {
	return &Transaction{
		Duration: 1 * time.Minute,
	}
}

// Start begins the transaction.
func (t *Transaction) Start() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.Active {
		return ErrTransactionAlreadyActive
	}

	t.Active = true
	t.StartedAt = time.Now()
	t.ExpiresAt = t.StartedAt.Add(t.Duration)
	t.Operations = []BufferedOperation{}
	return nil
}

// IsExpired checks if the transaction timer has expired.
func (t *Transaction) IsExpired() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.Active && time.Now().After(t.ExpiresAt)
}

// AddOperation adds a buffered operation to the transaction.
func (t *Transaction) AddOperation(op BufferedOperation) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.Active {
		return ErrNoActiveTransaction
	}
	if time.Now().After(t.ExpiresAt) {
		t.Active = false
		t.Operations = nil
		return ErrTransactionExpired
	}

	t.Operations = append(t.Operations, op)
	return nil
}

// Commit returns the buffered operations and ends the transaction.
func (t *Transaction) Commit() ([]BufferedOperation, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.Active {
		return nil, ErrNoActiveTransaction
	}
	if time.Now().After(t.ExpiresAt) {
		t.Active = false
		t.Operations = nil
		return nil, ErrTransactionExpired
	}

	ops := t.Operations
	t.Active = false
	t.Operations = nil
	return ops, nil
}

// Rollback discards all buffered operations and ends the transaction.
func (t *Transaction) Rollback() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.Active {
		return ErrNoActiveTransaction
	}

	t.Active = false
	t.Operations = nil
	return nil
}
