package testrunners

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// ConnectionRepoConstructor is a function that creates a ConnectionRepository
type ConnectionRepoConstructor func(superAdminConnStr string) interfaces.ConnectionRepository

// ConnectionRepositoryRunner tests ConnectionRepository implementations
func ConnectionRepositoryRunner(t *testing.T, constructor ConnectionRepoConstructor) {
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

	repo := constructor(connStr)

	t.Run("ValidateConnectionString - valid", func(t *testing.T) {
		err := repo.ValidateConnectionString(connStr)
		assert.NoError(t, err)
	})

	t.Run("ValidateConnectionString - invalid", func(t *testing.T) {
		err := repo.ValidateConnectionString("invalid://connection")
		assert.Error(t, err)
	})

	t.Run("TestConnection - success", func(t *testing.T) {
		err := repo.TestConnection(ctx, connStr)
		assert.NoError(t, err)
	})

	t.Run("TestConnection - failure", func(t *testing.T) {
		badConnStr := "postgres://baduser:badpass@localhost:5432/baddb?sslmode=disable"
		err := repo.TestConnection(ctx, badConnStr)
		assert.Error(t, err)
	})

	t.Run("GetConnection - success", func(t *testing.T) {
		db, err := repo.GetConnection(ctx, "postgres", "postgres", "testdb")
		require.NoError(t, err)
		assert.NotNil(t, db)
		defer db.Close()

		// Test the connection works
		err = db.Ping()
		assert.NoError(t, err)
	})
}

// MetadataRepoConstructor is a function that creates a MetadataRepository
type MetadataRepoConstructor func() interfaces.MetadataRepository

// MetadataRepositoryRunner tests MetadataRepository implementations
func MetadataRepositoryRunner(t *testing.T, constructor MetadataRepoConstructor) {
	t.Helper()

	ctx := context.Background()

	// Start PostgreSQL container with test data
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

	// Create test schema and tables
	_, err = db.Exec(`
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(100),
			created_at TIMESTAMP DEFAULT NOW()
		);
		CREATE TABLE posts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			title VARCHAR(200),
			content TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		);
	`)
	require.NoError(t, err)

	repo := constructor()

	t.Run("LoadGlobalMetadata", func(t *testing.T) {
		metadata, err := repo.LoadGlobalMetadata(ctx, db)
		require.NoError(t, err)
		assert.NotNil(t, metadata)
		assert.NotEmpty(t, metadata.Databases)
	})

	t.Run("LoadDatabaseMetadata", func(t *testing.T) {
		dbMeta, err := repo.LoadDatabaseMetadata(ctx, db, "testdb")
		require.NoError(t, err)
		assert.NotNil(t, dbMeta)
		assert.Equal(t, "testdb", dbMeta.Name)
		assert.NotEmpty(t, dbMeta.Tables)
	})

	t.Run("GetTableMetadata", func(t *testing.T) {
		tableMeta, err := repo.GetTableMetadata(ctx, db, "public", "users")
		require.NoError(t, err)
		assert.NotNil(t, tableMeta)
		assert.Equal(t, "users", tableMeta.Name)
		assert.NotEmpty(t, tableMeta.Columns)
	})

	t.Run("GetERDData", func(t *testing.T) {
		erdData, err := repo.GetERDData(ctx, db, "testdb", "public")
		require.NoError(t, err)
		assert.NotNil(t, erdData)
		assert.NotEmpty(t, erdData.Tables)
		assert.NotEmpty(t, erdData.Relationships)
	})
}

// QueryRepoConstructor is a function that creates a QueryRepository
type QueryRepoConstructor func() interfaces.QueryRepository

// QueryRepositoryRunner tests QueryRepository implementations
func QueryRepositoryRunner(t *testing.T, constructor QueryRepoConstructor) {
	t.Helper()

	ctx := context.Background()

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

	// Create test table with data
	_, err = db.Exec(`
		CREATE TABLE test_data (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100),
			value INTEGER
		);
		INSERT INTO test_data (name, value) VALUES 
			('item1', 10),
			('item2', 20),
			('item3', 30);
	`)
	require.NoError(t, err)

	repo := constructor()

	t.Run("ExecuteQuery - SELECT", func(t *testing.T) {
		req := domain.QueryRequest{
			SQL:    "SELECT * FROM test_data",
			Params: nil,
		}
		result, err := repo.ExecuteQuery(ctx, db, req)
		require.NoError(t, err)
		assert.True(t, result.IsSelect)
		assert.NotEmpty(t, result.Rows)
		assert.NotEmpty(t, result.Columns)
	})

	t.Run("ExecuteQuery - INSERT", func(t *testing.T) {
		req := domain.QueryRequest{
			SQL:    "INSERT INTO test_data (name, value) VALUES ($1, $2)",
			Params: []interface{}{"item4", 40},
		}
		result, err := repo.ExecuteQuery(ctx, db, req)
		require.NoError(t, err)
		assert.False(t, result.IsSelect)
		assert.Equal(t, int64(1), result.RowsAffected)
	})

	t.Run("ExecuteMultipleQueries", func(t *testing.T) {
		sql := "INSERT INTO test_data (name, value) VALUES ('item5', 50); SELECT * FROM test_data WHERE name = 'item5';"
		results, err := repo.ExecuteMultipleQueries(ctx, db, sql)
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("GetTableData - basic", func(t *testing.T) {
		req := domain.TableDataRequest{
			Schema: "public",
			Table:  "test_data",
			Limit:  50,
		}
		result, err := repo.GetTableData(ctx, db, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.Rows)
		assert.GreaterOrEqual(t, result.TotalCount, int64(3))
	})

	t.Run("GetTableData - with WHERE clause", func(t *testing.T) {
		req := domain.TableDataRequest{
			Schema:      "public",
			Table:       "test_data",
			WhereClause: "value > 20",
			Limit:       50,
		}
		result, err := repo.GetTableData(ctx, db, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("GetTableData - with sorting", func(t *testing.T) {
		req := domain.TableDataRequest{
			Schema:   "public",
			Table:    "test_data",
			OrderBy:  "value",
			OrderDir: "DESC",
			Limit:    50,
		}
		result, err := repo.GetTableData(ctx, db, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

// TransactionRepoConstructor is a function that creates a TransactionRepository
type TransactionRepoConstructor func() interfaces.TransactionRepository

// TransactionRepositoryRunner tests TransactionRepository implementations
func TransactionRepositoryRunner(t *testing.T, constructor TransactionRepoConstructor) {
	t.Helper()

	ctx := context.Background()

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
	_, err = db.Exec(`
		CREATE TABLE txn_test (
			id SERIAL PRIMARY KEY,
			data VARCHAR(100)
		);
	`)
	require.NoError(t, err)

	repo := constructor()

	t.Run("BeginTransaction - success", func(t *testing.T) {
		tx, err := repo.BeginTransaction(ctx, db)
		require.NoError(t, err)
		assert.NotNil(t, tx)
		tx.Rollback()
	})

	t.Run("CommitTransaction - success", func(t *testing.T) {
		tx, err := repo.BeginTransaction(ctx, db)
		require.NoError(t, err)

		_, err = tx.Exec("INSERT INTO txn_test (data) VALUES ('test')")
		require.NoError(t, err)

		err = repo.CommitTransaction(ctx, tx)
		assert.NoError(t, err)

		// Verify data was committed
		var count int
		db.QueryRow("SELECT COUNT(*) FROM txn_test WHERE data = 'test'").Scan(&count)
		assert.Equal(t, 1, count)
	})

	t.Run("RollbackTransaction - success", func(t *testing.T) {
		tx, err := repo.BeginTransaction(ctx, db)
		require.NoError(t, err)

		_, err = tx.Exec("INSERT INTO txn_test (data) VALUES ('rollback_test')")
		require.NoError(t, err)

		err = repo.RollbackTransaction(ctx, tx)
		assert.NoError(t, err)

		// Verify data was not committed
		var count int
		db.QueryRow("SELECT COUNT(*) FROM txn_test WHERE data = 'rollback_test'").Scan(&count)
		assert.Equal(t, 0, count)
	})

	t.Run("ExecuteOperations - buffered operations", func(t *testing.T) {
		tx, err := repo.BeginTransaction(ctx, db)
		require.NoError(t, err)

		operations := []domain.Operation{
			{
				Type:   domain.OperationInsert,
				Schema: "public",
				Table:  "txn_test",
				Data: map[string]interface{}{
					"data": "op1",
				},
			},
			{
				Type:   domain.OperationInsert,
				Schema: "public",
				Table:  "txn_test",
				Data: map[string]interface{}{
					"data": "op2",
				},
			},
		}

		err = repo.ExecuteOperations(ctx, tx, operations)
		require.NoError(t, err)

		err = repo.CommitTransaction(ctx, tx)
		assert.NoError(t, err)

		// Verify operations were executed
		var count int
		db.QueryRow("SELECT COUNT(*) FROM txn_test WHERE data IN ('op1', 'op2')").Scan(&count)
		assert.Equal(t, 2, count)
	})
}
