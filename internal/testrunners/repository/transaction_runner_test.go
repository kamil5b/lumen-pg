package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// TransactionRepositoryConstructor is a function type that creates a TransactionRepository
type TransactionRepositoryConstructor func(db *sql.DB) repository.TransactionRepository

// TransactionRepositoryRunner runs all transaction repository tests against an implementation
func TransactionRepositoryRunner(t *testing.T, constructor TransactionRepositoryConstructor) {
	t.Helper()

	ctx := context.Background()

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
	)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	connStr, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	err = db.PingContext(ctx)
	require.NoError(t, err)

	repo := constructor(db)

	t.Run("CreateTransaction and GetTransaction", func(t *testing.T) {
		now := time.Now()
		txn := &domain.TransactionState{
			ID:        "txn_123",
			Username:  "testuser",
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, txn)
		require.NoError(t, err)

		retrieved, err := repo.GetTransaction(ctx, "txn_123")
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		require.Equal(t, txn.ID, retrieved.ID)
		require.Equal(t, txn.Username, retrieved.Username)
	})

	t.Run("GetTransaction returns error for non-existent transaction", func(t *testing.T) {
		_, err := repo.GetTransaction(ctx, "nonexistent_txn")
		require.Error(t, err)
	})

	t.Run("UpdateTransaction modifies existing transaction", func(t *testing.T) {
		now := time.Now()
		txn := &domain.TransactionState{
			ID:        "update_txn",
			Username:  "testuser",
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, txn)
		require.NoError(t, err)

		newExpiry := now.Add(60 * time.Minute)
		txn.ExpiresAt = newExpiry

		err = repo.UpdateTransaction(ctx, txn)
		require.NoError(t, err)

		retrieved, err := repo.GetTransaction(ctx, "update_txn")
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		require.True(t, retrieved.ExpiresAt.After(newExpiry.Add(-1*time.Second)))
	})

	t.Run("DeleteTransaction removes transaction", func(t *testing.T) {
		now := time.Now()
		txn := &domain.TransactionState{
			ID:        "delete_txn",
			Username:  "testuser",
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, txn)
		require.NoError(t, err)

		err = repo.DeleteTransaction(ctx, "delete_txn")
		require.NoError(t, err)

		_, err = repo.GetTransaction(ctx, "delete_txn")
		require.Error(t, err)
	})

	t.Run("GetUserTransaction retrieves active transaction for user", func(t *testing.T) {
		now := time.Now()
		username := "user_txn"
		txn := &domain.TransactionState{
			ID:        "user_txn_1",
			Username:  username,
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, txn)
		require.NoError(t, err)

		retrieved, err := repo.GetUserTransaction(ctx, username)
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		require.Equal(t, username, retrieved.Username)
	})

	t.Run("GetUserTransaction returns error for user with no transaction", func(t *testing.T) {
		_, err := repo.GetUserTransaction(ctx, "user_no_txn")
		require.Error(t, err)
	})

	t.Run("AddRowEdit buffers cell edit", func(t *testing.T) {
		now := time.Now()
		txn := &domain.TransactionState{
			ID:        "edit_txn",
			Username:  "testuser",
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, txn)
		require.NoError(t, err)

		edit := domain.RowEdit{
			RowIndex:   0,
			ColumnName: "name",
			OldValue:   "oldname",
			NewValue:   "newname",
		}

		err = repo.AddRowEdit(ctx, "edit_txn", edit)
		require.NoError(t, err)

		edits, err := repo.GetRowEdits(ctx, "edit_txn")
		require.NoError(t, err)
		require.NotNil(t, edits)
		require.Len(t, edits, 1)
	})

	t.Run("AddRowDelete buffers row deletion", func(t *testing.T) {
		now := time.Now()
		txn := &domain.TransactionState{
			ID:        "delete_row_txn",
			Username:  "testuser",
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, txn)
		require.NoError(t, err)

		err = repo.AddRowDelete(ctx, "delete_row_txn", 1)
		require.NoError(t, err)

		deletes, err := repo.GetRowDeletes(ctx, "delete_row_txn")
		require.NoError(t, err)
		require.NotNil(t, deletes)
		require.Len(t, deletes, 1)
		require.Contains(t, deletes, 1)
	})

	t.Run("AddRowInsert buffers row insertion", func(t *testing.T) {
		now := time.Now()
		txn := &domain.TransactionState{
			ID:        "insert_row_txn",
			Username:  "testuser",
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, txn)
		require.NoError(t, err)

		insert := domain.RowInsert{
			Values: map[string]interface{}{
				"name": "newrow",
				"age":  25,
			},
		}

		err = repo.AddRowInsert(ctx, "insert_row_txn", insert)
		require.NoError(t, err)

		inserts, err := repo.GetRowInserts(ctx, "insert_row_txn")
		require.NoError(t, err)
		require.NotNil(t, inserts)
		require.Len(t, inserts, 1)
	})

	t.Run("GetRowEdits returns all buffered edits", func(t *testing.T) {
		now := time.Now()
		txn := &domain.TransactionState{
			ID:        "get_edits_txn",
			Username:  "testuser",
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, txn)
		require.NoError(t, err)

		edit1 := domain.RowEdit{
			RowIndex:   0,
			ColumnName: "col1",
			OldValue:   "val1",
			NewValue:   "newval1",
		}

		edit2 := domain.RowEdit{
			RowIndex:   1,
			ColumnName: "col2",
			OldValue:   "val2",
			NewValue:   "newval2",
		}

		err = repo.AddRowEdit(ctx, "get_edits_txn", edit1)
		require.NoError(t, err)

		err = repo.AddRowEdit(ctx, "get_edits_txn", edit2)
		require.NoError(t, err)

		edits, err := repo.GetRowEdits(ctx, "get_edits_txn")
		require.NoError(t, err)
		require.Len(t, edits, 2)
	})

	t.Run("GetRowDeletes returns all buffered deletions", func(t *testing.T) {
		now := time.Now()
		txn := &domain.TransactionState{
			ID:        "get_deletes_txn",
			Username:  "testuser",
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, txn)
		require.NoError(t, err)

		err = repo.AddRowDelete(ctx, "get_deletes_txn", 0)
		require.NoError(t, err)

		err = repo.AddRowDelete(ctx, "get_deletes_txn", 2)
		require.NoError(t, err)

		deletes, err := repo.GetRowDeletes(ctx, "get_deletes_txn")
		require.NoError(t, err)
		require.Len(t, deletes, 2)
	})

	t.Run("GetRowInserts returns all buffered insertions", func(t *testing.T) {
		now := time.Now()
		txn := &domain.TransactionState{
			ID:        "get_inserts_txn",
			Username:  "testuser",
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, txn)
		require.NoError(t, err)

		insert1 := domain.RowInsert{
			Values: map[string]interface{}{"col1": "val1"},
		}

		insert2 := domain.RowInsert{
			Values: map[string]interface{}{"col1": "val2"},
		}

		err = repo.AddRowInsert(ctx, "get_inserts_txn", insert1)
		require.NoError(t, err)

		err = repo.AddRowInsert(ctx, "get_inserts_txn", insert2)
		require.NoError(t, err)

		inserts, err := repo.GetRowInserts(ctx, "get_inserts_txn")
		require.NoError(t, err)
		require.Len(t, inserts, 2)
	})

	t.Run("ClearRowEdits removes all edits", func(t *testing.T) {
		now := time.Now()
		txn := &domain.TransactionState{
			ID:        "clear_edits_txn",
			Username:  "testuser",
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, txn)
		require.NoError(t, err)

		edit := domain.RowEdit{
			RowIndex:   0,
			ColumnName: "col",
			OldValue:   "old",
			NewValue:   "new",
		}

		err = repo.AddRowEdit(ctx, "clear_edits_txn", edit)
		require.NoError(t, err)

		err = repo.ClearRowEdits(ctx, "clear_edits_txn")
		require.NoError(t, err)

		edits, err := repo.GetRowEdits(ctx, "clear_edits_txn")
		require.NoError(t, err)
		require.Len(t, edits, 0)
	})

	t.Run("ClearRowDeletes removes all deletions", func(t *testing.T) {
		now := time.Now()
		txn := &domain.TransactionState{
			ID:        "clear_deletes_txn",
			Username:  "testuser",
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, txn)
		require.NoError(t, err)

		err = repo.AddRowDelete(ctx, "clear_deletes_txn", 1)
		require.NoError(t, err)

		err = repo.ClearRowDeletes(ctx, "clear_deletes_txn")
		require.NoError(t, err)

		deletes, err := repo.GetRowDeletes(ctx, "clear_deletes_txn")
		require.NoError(t, err)
		require.Len(t, deletes, 0)
	})

	t.Run("ClearRowInserts removes all insertions", func(t *testing.T) {
		now := time.Now()
		txn := &domain.TransactionState{
			ID:        "clear_inserts_txn",
			Username:  "testuser",
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, txn)
		require.NoError(t, err)

		insert := domain.RowInsert{
			Values: map[string]interface{}{"col": "val"},
		}

		err = repo.AddRowInsert(ctx, "clear_inserts_txn", insert)
		require.NoError(t, err)

		err = repo.ClearRowInserts(ctx, "clear_inserts_txn")
		require.NoError(t, err)

		inserts, err := repo.GetRowInserts(ctx, "clear_inserts_txn")
		require.NoError(t, err)
		require.Len(t, inserts, 0)
	})

	t.Run("TransactionExists returns true for existing transaction", func(t *testing.T) {
		now := time.Now()
		txn := &domain.TransactionState{
			ID:        "exists_txn",
			Username:  "testuser",
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, txn)
		require.NoError(t, err)

		exists, err := repo.TransactionExists(ctx, "exists_txn")
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("TransactionExists returns false for non-existent transaction", func(t *testing.T) {
		exists, err := repo.TransactionExists(ctx, "nonexistent_txn")
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("TransactionExists returns false for expired transaction", func(t *testing.T) {
		now := time.Now()
		txn := &domain.TransactionState{
			ID:        "expired_txn",
			Username:  "testuser",
			StartedAt: now.Add(-1 * time.Hour),
			ExpiresAt: now.Add(-10 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, txn)
		require.NoError(t, err)

		exists, err := repo.TransactionExists(ctx, "expired_txn")
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("InvalidateExpiredTransactions removes expired transactions", func(t *testing.T) {
		now := time.Now()

		expiredTxn := &domain.TransactionState{
			ID:        "old_txn",
			Username:  "user1",
			StartedAt: now.Add(-1 * time.Hour),
			ExpiresAt: now.Add(-10 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		activeTxn := &domain.TransactionState{
			ID:        "current_txn",
			Username:  "user2",
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, expiredTxn)
		require.NoError(t, err)

		err = repo.CreateTransaction(ctx, activeTxn)
		require.NoError(t, err)

		err = repo.InvalidateExpiredTransactions(ctx)
		require.NoError(t, err)

		_, err = repo.GetTransaction(ctx, "old_txn")
		require.Error(t, err)

		retrieved, err := repo.GetTransaction(ctx, "current_txn")
		require.NoError(t, err)
		require.NotNil(t, retrieved)
	})

	t.Run("Mixed transaction operations", func(t *testing.T) {
		now := time.Now()
		txn := &domain.TransactionState{
			ID:        "mixed_txn",
			Username:  "testuser",
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, txn)
		require.NoError(t, err)

		edit := domain.RowEdit{
			RowIndex:   0,
			ColumnName: "col1",
			OldValue:   "old1",
			NewValue:   "new1",
		}

		err = repo.AddRowEdit(ctx, "mixed_txn", edit)
		require.NoError(t, err)

		err = repo.AddRowDelete(ctx, "mixed_txn", 2)
		require.NoError(t, err)

		insert := domain.RowInsert{
			Values: map[string]interface{}{"col": "val"},
		}

		err = repo.AddRowInsert(ctx, "mixed_txn", insert)
		require.NoError(t, err)

		edits, _ := repo.GetRowEdits(ctx, "mixed_txn")
		deletes, _ := repo.GetRowDeletes(ctx, "mixed_txn")
		inserts, _ := repo.GetRowInserts(ctx, "mixed_txn")

		require.Len(t, edits, 1)
		require.Len(t, deletes, 1)
		require.Len(t, inserts, 1)
	})

	t.Run("Multiple transactions for different users", func(t *testing.T) {
		now := time.Now()

		txn1 := &domain.TransactionState{
			ID:        "user1_txn",
			Username:  "user1",
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		txn2 := &domain.TransactionState{
			ID:        "user2_txn",
			Username:  "user2",
			StartedAt: now,
			ExpiresAt: now.Add(30 * time.Minute),
			Edits:     make(map[int]domain.RowEdit),
			Deletes:   []int{},
			Inserts:   []domain.RowInsert{},
		}

		err := repo.CreateTransaction(ctx, txn1)
		require.NoError(t, err)

		err = repo.CreateTransaction(ctx, txn2)
		require.NoError(t, err)

		userTxn1, err := repo.GetUserTransaction(ctx, "user1")
		require.NoError(t, err)
		require.Equal(t, "user1", userTxn1.Username)

		userTxn2, err := repo.GetUserTransaction(ctx, "user2")
		require.NoError(t, err)
		require.Equal(t, "user2", userTxn2.Username)
	})
}
