package domain

import "time"

// Transaction represents an active transaction
type Transaction struct {
	ID           string
	Username     string
	TableName    string
	Operations   []TransactionOperation
	StartedAt    time.Time
	ExpiresAt    time.Time
	IsCommitted  bool
	IsRolledBack bool
}

// TransactionOperation represents a buffered operation
type TransactionOperation struct {
	Type       OperationType
	TableName  string
	PrimaryKey interface{}
	Column     string
	OldValue   interface{}
	NewValue   interface{}
}

// OperationType represents the type of transaction operation
type OperationType string

const (
	OperationUpdate OperationType = "UPDATE"
	OperationDelete OperationType = "DELETE"
	OperationInsert OperationType = "INSERT"
)
