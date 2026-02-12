package domain_test

import (
	"testing"
	"time"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// UC-S5-09: Transaction Start
func TestTransaction_Start(t *testing.T) {
	txn := domain.NewTransaction()
	err := txn.Start()
	require.NoError(t, err)
	assert.True(t, txn.Active)
	assert.False(t, txn.StartedAt.IsZero())
	assert.False(t, txn.ExpiresAt.IsZero())
	assert.Equal(t, 1*time.Minute, txn.Duration)
	assert.Empty(t, txn.Operations)
}

// UC-S5-10: Transaction Already Active Error
func TestTransaction_Start_AlreadyActive(t *testing.T) {
	txn := domain.NewTransaction()
	err := txn.Start()
	require.NoError(t, err)

	err = txn.Start()
	assert.ErrorIs(t, err, domain.ErrTransactionAlreadyActive)
}

// UC-S5-11: Cell Edit Buffering
func TestTransaction_AddOperation(t *testing.T) {
	txn := domain.NewTransaction()
	err := txn.Start()
	require.NoError(t, err)

	op := domain.BufferedOperation{
		Type:    domain.OpUpdate,
		Table:   "users",
		Schema:  "public",
		RowData: map[string]interface{}{"name": "John"},
	}
	err = txn.AddOperation(op)
	require.NoError(t, err)
	assert.Len(t, txn.Operations, 1)
	assert.Equal(t, domain.OpUpdate, txn.Operations[0].Type)
}

// UC-S5-11: Cell Edit Buffering - No Active Transaction
func TestTransaction_AddOperation_NoActiveTransaction(t *testing.T) {
	txn := domain.NewTransaction()
	op := domain.BufferedOperation{
		Type:  domain.OpUpdate,
		Table: "users",
	}
	err := txn.AddOperation(op)
	assert.ErrorIs(t, err, domain.ErrNoActiveTransaction)
}

// UC-S5-12: Transaction Commit
func TestTransaction_Commit(t *testing.T) {
	txn := domain.NewTransaction()
	err := txn.Start()
	require.NoError(t, err)

	op1 := domain.BufferedOperation{Type: domain.OpInsert, Table: "users"}
	op2 := domain.BufferedOperation{Type: domain.OpUpdate, Table: "users"}
	err = txn.AddOperation(op1)
	require.NoError(t, err)
	err = txn.AddOperation(op2)
	require.NoError(t, err)

	ops, err := txn.Commit()
	require.NoError(t, err)
	assert.Len(t, ops, 2)
	assert.False(t, txn.Active)
}

// UC-S5-12: Transaction Commit - No Active Transaction
func TestTransaction_Commit_NoActiveTransaction(t *testing.T) {
	txn := domain.NewTransaction()
	ops, err := txn.Commit()
	assert.ErrorIs(t, err, domain.ErrNoActiveTransaction)
	assert.Nil(t, ops)
}

// UC-S5-13: Transaction Rollback
func TestTransaction_Rollback(t *testing.T) {
	txn := domain.NewTransaction()
	err := txn.Start()
	require.NoError(t, err)

	op := domain.BufferedOperation{Type: domain.OpDelete, Table: "users"}
	err = txn.AddOperation(op)
	require.NoError(t, err)

	err = txn.Rollback()
	require.NoError(t, err)
	assert.False(t, txn.Active)
	assert.Nil(t, txn.Operations)
}

// UC-S5-13: Transaction Rollback - No Active Transaction
func TestTransaction_Rollback_NoActiveTransaction(t *testing.T) {
	txn := domain.NewTransaction()
	err := txn.Rollback()
	assert.ErrorIs(t, err, domain.ErrNoActiveTransaction)
}

// UC-S5-14: Transaction Timer Expiration
func TestTransaction_IsExpired(t *testing.T) {
	txn := domain.NewTransaction()
	txn.Duration = 1 * time.Millisecond
	err := txn.Start()
	require.NoError(t, err)

	// Wait for expiration
	time.Sleep(5 * time.Millisecond)
	assert.True(t, txn.IsExpired())
}

// UC-S5-14: Transaction Timer Not Yet Expired
func TestTransaction_IsNotExpired(t *testing.T) {
	txn := domain.NewTransaction()
	err := txn.Start()
	require.NoError(t, err)
	assert.False(t, txn.IsExpired())
}

// UC-S5-14: Transaction AddOperation After Expiry
func TestTransaction_AddOperation_Expired(t *testing.T) {
	txn := domain.NewTransaction()
	txn.Duration = 1 * time.Millisecond
	err := txn.Start()
	require.NoError(t, err)

	time.Sleep(5 * time.Millisecond)

	op := domain.BufferedOperation{Type: domain.OpUpdate, Table: "users"}
	err = txn.AddOperation(op)
	assert.ErrorIs(t, err, domain.ErrTransactionExpired)
}

// UC-S5-14: Transaction Commit After Expiry
func TestTransaction_Commit_Expired(t *testing.T) {
	txn := domain.NewTransaction()
	txn.Duration = 1 * time.Millisecond
	err := txn.Start()
	require.NoError(t, err)

	time.Sleep(5 * time.Millisecond)

	ops, err := txn.Commit()
	assert.ErrorIs(t, err, domain.ErrTransactionExpired)
	assert.Nil(t, ops)
}

// UC-S5-15: Row Deletion Buffering
func TestTransaction_DeleteOperation(t *testing.T) {
	txn := domain.NewTransaction()
	err := txn.Start()
	require.NoError(t, err)

	op := domain.BufferedOperation{
		Type:      domain.OpDelete,
		Table:     "users",
		Schema:    "public",
		WhereData: map[string]interface{}{"id": 1},
	}
	err = txn.AddOperation(op)
	require.NoError(t, err)
	assert.Equal(t, domain.OpDelete, txn.Operations[0].Type)
}

// UC-S5-16: Row Insertion Buffering
func TestTransaction_InsertOperation(t *testing.T) {
	txn := domain.NewTransaction()
	err := txn.Start()
	require.NoError(t, err)

	op := domain.BufferedOperation{
		Type:    domain.OpInsert,
		Table:   "users",
		Schema:  "public",
		RowData: map[string]interface{}{"username": "newuser", "email": "new@example.com"},
	}
	err = txn.AddOperation(op)
	require.NoError(t, err)
	assert.Equal(t, domain.OpInsert, txn.Operations[0].Type)
}
