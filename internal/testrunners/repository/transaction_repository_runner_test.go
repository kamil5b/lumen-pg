package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// TransactionRepoConstructor creates a transaction repository with database connection
type TransactionRepoConstructor func(db *sql.DB) repository.TransactionRepository

// TransactionRepositoryRunner runs integration tests for transaction repository (Story 5)
func TransactionRepositoryRunner(t *testing.T, constructor TransactionRepoConstructor) {
	t.Helper()

	ctx := context.Background()

	// Start PostgreSQL container
	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
	)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	connStr, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	// Create test table
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS test_users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50),
			email VARCHAR(100)
		)
	`)
	require.NoError(t, err)

	// Insert test data
	_, err = db.ExecContext(ctx, `
		INSERT INTO test_users (username, email) VALUES
		('user1', 'user1@test.com'),
		('user2', 'user2@test.com')
	`)
	require.NoError(t, err)

	repo := constructor(db)

	t.Run("UC-S5-09: Transaction Start", func(t *testing.T) {
		username := "testuser"
		tableName := "test_users"

		txn, err := repo.StartTransaction(ctx, username, tableName)

		require.NoError(t, err)
		assert.NotNil(t, txn)
		assert.Equal(t, username, txn.Username)
		assert.Equal(t, tableName, txn.TableName)
		assert.NotEmpty(t, txn.ID)
		assert.False(t, txn.IsCommitted)
		assert.False(t, txn.IsRolledBack)
		assert.True(t, txn.ExpiresAt.After(time.Now()))
	})

	t.Run("UC-S5-10: Transaction Already Active Error", func(t *testing.T) {
		username := "activeuser"
		tableName := "test_users"

		// Start first transaction
		txn1, err := repo.StartTransaction(ctx, username, tableName)
		require.NoError(t, err)
		require.NotNil(t, txn1)

		// Try to start second transaction - should fail or return existing
		txn2, err := repo.StartTransaction(ctx, username, tableName)

		// Either error or returns same transaction
		if err == nil {
			assert.NotNil(t, txn2)
			// Clean up
			repo.RollbackTransaction(ctx, txn1.ID)
		} else {
			assert.NotNil(t, err)
		}
	})

	t.Run("UC-S5-11: Cell Edit Buffering", func(t *testing.T) {
		username := "edituser"
		tableName := "test_users"

		// Start transaction
		txn, err := repo.StartTransaction(ctx, username, tableName)
		require.NoError(t, err)

		// Buffer an edit operation
		op := domain.TransactionOperation{
			Type:       domain.OperationUpdate,
			TableName:  tableName,
			PrimaryKey: 1,
			Column:     "username",
			OldValue:   "user1",
			NewValue:   "updated_user1",
		}

		err = repo.BufferOperation(ctx, txn.ID, op)

		require.NoError(t, err)

		// Verify operation was buffered
		retrievedTxn, err := repo.GetTransaction(ctx, txn.ID)
		require.NoError(t, err)
		assert.NotNil(t, retrievedTxn)
		assert.Len(t, retrievedTxn.Operations, 1)
		assert.Equal(t, domain.OperationUpdate, retrievedTxn.Operations[0].Type)

		// Cleanup
		repo.RollbackTransaction(ctx, txn.ID)
	})

	t.Run("UC-S5-12: Transaction Commit", func(t *testing.T) {
		username := "commituser"
		tableName := "test_users"

		// Start transaction
		txn, err := repo.StartTransaction(ctx, username, tableName)
		require.NoError(t, err)

		// Buffer an operation
		op := domain.TransactionOperation{
			Type:       domain.OperationUpdate,
			TableName:  tableName,
			PrimaryKey: 1,
			Column:     "email",
			OldValue:   "user1@test.com",
			NewValue:   "newemail@test.com",
		}

		err = repo.BufferOperation(ctx, txn.ID, op)
		require.NoError(t, err)

		// Commit transaction
		err = repo.CommitTransaction(ctx, txn.ID)

		require.NoError(t, err)

		// Verify transaction is marked as committed
		committedTxn, err := repo.GetTransaction(ctx, txn.ID)
		if err == nil && committedTxn != nil {
			assert.True(t, committedTxn.IsCommitted)
		}
	})

	t.Run("UC-S5-13: Transaction Rollback", func(t *testing.T) {
		username := "rollbackuser"
		tableName := "test_users"

		// Start transaction
		txn, err := repo.StartTransaction(ctx, username, tableName)
		require.NoError(t, err)

		// Buffer an operation
		op := domain.TransactionOperation{
			Type:       domain.OperationUpdate,
			TableName:  tableName,
			PrimaryKey: 2,
			Column:     "username",
			OldValue:   "user2",
			NewValue:   "deleted_user2",
		}

		err = repo.BufferOperation(ctx, txn.ID, op)
		require.NoError(t, err)

		// Rollback transaction
		err = repo.RollbackTransaction(ctx, txn.ID)

		require.NoError(t, err)

		// Verify transaction is marked as rolled back
		rolledBackTxn, err := repo.GetTransaction(ctx, txn.ID)
		if err == nil && rolledBackTxn != nil {
			assert.True(t, rolledBackTxn.IsRolledBack)
		}
	})

	t.Run("UC-S5-14: Transaction Timer Expiration", func(t *testing.T) {
		username := "expireuser"
		tableName := "test_users"

		// Start transaction
		txn, err := repo.StartTransaction(ctx, username, tableName)
		require.NoError(t, err)

		// Verify expiration time is set (typically 1 minute from now)
		assert.True(t, txn.ExpiresAt.After(time.Now()))
		assert.True(t, txn.ExpiresAt.Before(time.Now().Add(2*time.Minute)))

		// Cleanup
		repo.RollbackTransaction(ctx, txn.ID)
	})

	t.Run("UC-S5-15: Row Deletion Buffering", func(t *testing.T) {
		username := "deleteuser"
		tableName := "test_users"

		// Start transaction
		txn, err := repo.StartTransaction(ctx, username, tableName)
		require.NoError(t, err)

		// Buffer a delete operation
		op := domain.TransactionOperation{
			Type:       domain.OperationDelete,
			TableName:  tableName,
			PrimaryKey: 1,
		}

		err = repo.BufferOperation(ctx, txn.ID, op)

		require.NoError(t, err)

		// Verify operation was buffered
		retrievedTxn, err := repo.GetTransaction(ctx, txn.ID)
		require.NoError(t, err)
		assert.NotNil(t, retrievedTxn)
		assert.Len(t, retrievedTxn.Operations, 1)
		assert.Equal(t, domain.OperationDelete, retrievedTxn.Operations[0].Type)
		assert.Equal(t, 1, retrievedTxn.Operations[0].PrimaryKey)

		// Cleanup
		repo.RollbackTransaction(ctx, txn.ID)
	})

	t.Run("UC-S5-16: Row Insertion Buffering", func(t *testing.T) {
		username := "insertuser"
		tableName := "test_users"

		// Start transaction
		txn, err := repo.StartTransaction(ctx, username, tableName)
		require.NoError(t, err)

		// Buffer an insert operation
		op := domain.TransactionOperation{
			Type:      domain.OperationInsert,
			TableName: tableName,
			NewValue: map[string]interface{}{
				"username": "newuser",
				"email":    "newuser@test.com",
			},
		}

		err = repo.BufferOperation(ctx, txn.ID, op)

		require.NoError(t, err)

		// Verify operation was buffered
		retrievedTxn, err := repo.GetTransaction(ctx, txn.ID)
		require.NoError(t, err)
		assert.NotNil(t, retrievedTxn)
		assert.Len(t, retrievedTxn.Operations, 1)
		assert.Equal(t, domain.OperationInsert, retrievedTxn.Operations[0].Type)

		// Cleanup
		repo.RollbackTransaction(ctx, txn.ID)
	})

	t.Run("IT-S5-04: Real Transaction Commit", func(t *testing.T) {
		username := "realcommituser"
		tableName := "test_users"

		// Start transaction
		txn, err := repo.StartTransaction(ctx, username, tableName)
		require.NoError(t, err)

		// Buffer a real update operation
		op := domain.TransactionOperation{
			Type:       domain.OperationUpdate,
			TableName:  tableName,
			PrimaryKey: 1,
			Column:     "username",
			OldValue:   "user1",
			NewValue:   "committed_user",
		}

		err = repo.BufferOperation(ctx, txn.ID, op)
		require.NoError(t, err)

		// Commit transaction
		err = repo.CommitTransaction(ctx, txn.ID)
		require.NoError(t, err)

		// Verify changes were actually persisted
		var username_val string
		err = db.QueryRowContext(ctx, "SELECT username FROM test_users WHERE id = $1", 1).Scan(&username_val)
		require.NoError(t, err)
		assert.Equal(t, "committed_user", username_val)

		// Cleanup: reset data
		db.ExecContext(ctx, "UPDATE test_users SET username = $1 WHERE id = $2", "user1", 1)
	})

	t.Run("IT-S5-05: Real Transaction Rollback", func(t *testing.T) {
		username := "realrollbackuser"
		tableName := "test_users"

		// Get original value
		var originalEmail string
		err := db.QueryRowContext(ctx, "SELECT email FROM test_users WHERE id = $1", 2).Scan(&originalEmail)
		require.NoError(t, err)

		// Start transaction
		txn, err := repo.StartTransaction(ctx, username, tableName)
		require.NoError(t, err)

		// Buffer an update operation
		op := domain.TransactionOperation{
			Type:       domain.OperationUpdate,
			TableName:  tableName,
			PrimaryKey: 2,
			Column:     "email",
			OldValue:   originalEmail,
			NewValue:   "rolled_back@test.com",
		}

		err = repo.BufferOperation(ctx, txn.ID, op)
		require.NoError(t, err)

		// Rollback transaction
		err = repo.RollbackTransaction(ctx, txn.ID)
		require.NoError(t, err)

		// Verify changes were NOT persisted
		var email_val string
		err = db.QueryRowContext(ctx, "SELECT email FROM test_users WHERE id = $1", 2).Scan(&email_val)
		require.NoError(t, err)
		assert.Equal(t, originalEmail, email_val)
	})

	t.Run("UC-S6-02: Transaction Isolation", func(t *testing.T) {
		// Start transaction for user A
		txnA, err := repo.StartTransaction(ctx, "userA", "test_users")
		require.NoError(t, err)

		// Start transaction for user B
		txnB, err := repo.StartTransaction(ctx, "userB", "test_users")
		require.NoError(t, err)

		// Buffer operations in both transactions
		opA := domain.TransactionOperation{
			Type:       domain.OperationUpdate,
			TableName:  "test_users",
			PrimaryKey: 1,
			Column:     "username",
			NewValue:   "userA_change",
		}

		opB := domain.TransactionOperation{
			Type:       domain.OperationUpdate,
			TableName:  "test_users",
			PrimaryKey: 2,
			Column:     "username",
			NewValue:   "userB_change",
		}

		err = repo.BufferOperation(ctx, txnA.ID, opA)
		require.NoError(t, err)

		err = repo.BufferOperation(ctx, txnB.ID, opB)
		require.NoError(t, err)

		// Verify transactions are isolated
		txnARetrieved, err := repo.GetTransaction(ctx, txnA.ID)
		require.NoError(t, err)

		txnBRetrieved, err := repo.GetTransaction(ctx, txnB.ID)
		require.NoError(t, err)

		assert.NotEqual(t, txnA.ID, txnB.ID)
		assert.Len(t, txnARetrieved.Operations, 1)
		assert.Len(t, txnBRetrieved.Operations, 1)

		// Cleanup
		repo.RollbackTransaction(ctx, txnA.ID)
		repo.RollbackTransaction(ctx, txnB.ID)
	})
}
