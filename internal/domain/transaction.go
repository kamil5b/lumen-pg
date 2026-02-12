package domain

import "time"

// TransactionState represents the state of a user's transaction
type TransactionState struct {
	Active     bool
	StartedAt  time.Time
	ExpiresAt  time.Time
	Operations []Operation
}

// Operation represents a buffered database operation
type Operation struct {
	Type   OperationType
	Schema string
	Table  string
	Data   map[string]interface{}
	Where  map[string]interface{} // For UPDATE/DELETE
}

// OperationType represents the type of database operation
type OperationType string

const (
	OperationInsert OperationType = "INSERT"
	OperationUpdate OperationType = "UPDATE"
	OperationDelete OperationType = "DELETE"
)

// TransactionTimeout is the transaction timeout duration (1 minute)
const TransactionTimeout = 1 * time.Minute
