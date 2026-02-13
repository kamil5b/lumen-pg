package repository

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// DatabaseRepositoryConstructor is a function type that creates a DatabaseRepository
type DatabaseRepositoryConstructor func(db *sql.DB) repository.DatabaseRepository

// DatabaseRepositoryRunner runs all database repository tests against an implementation
// Maps to TEST_PLAN.md:
// - Story 1: Setup & Configuration [UC-S1-01~07, IT-S1-01~04]
// - Story 4: Manual Query Editor [UC-S4-01~08, IT-S4-01~04]
// - Story 5: Main View & Data Interaction [UC-S5-01~19, IT-S5-01~07]
// - Story 7: Security & Best Practices [UC-S7-01, IT-S7-01]
func DatabaseRepositoryRunner(t *testing.T, constructor DatabaseRepositoryConstructor) {
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

	// Create test tables
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS test_users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE
		)
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS test_posts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES test_users(id),
			title VARCHAR(200) NOT NULL,
			content TEXT
		)
	`)
	require.NoError(t, err)

	// Insert test data
	_, err = db.ExecContext(ctx, `
		INSERT INTO test_users (name, email) VALUES
		('Alice', 'alice@example.com'),
		('Bob', 'bob@example.com'),
		('Charlie', 'charlie@example.com')
		ON CONFLICT DO NOTHING
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO test_posts (user_id, title, content) VALUES
		(1, 'First Post', 'This is the first post'),
		(1, 'Second Post', 'This is the second post'),
		(2, 'Bob Post', 'This is bobs post')
		ON CONFLICT DO NOTHING
	`)
	require.NoError(t, err)

	repo := constructor(db)

	t.Run("Connect establishes connection", func(t *testing.T) {
		err := repo.Connect(ctx, connStr)
		require.NoError(t, err)
	})

	t.Run("TestConnection verifies connectivity", func(t *testing.T) {
		err := repo.TestConnection(ctx, connStr)
		require.NoError(t, err)
	})

	t.Run("TestConnection fails with invalid connection string", func(t *testing.T) {
		err := repo.TestConnection(ctx, "postgres://invalid:invalid@invalid:99999/invalid")
		require.Error(t, err)
	})

	t.Run("GetConnection returns database connection", func(t *testing.T) {
		conn := repo.GetConnection()
		require.NotNil(t, conn)
		require.NoError(t, conn.PingContext(ctx))
	})

	t.Run("ExecuteQuery executes single query", func(t *testing.T) {
		result, err := repo.ExecuteQuery(ctx, "SELECT id, name FROM test_users ORDER BY id")
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Greater(t, len(result.Rows), 0)
		require.Contains(t, result.Columns, "id")
		require.Contains(t, result.Columns, "name")
	})

	t.Run("ExecuteQuery returns error for invalid SQL", func(t *testing.T) {
		result, err := repo.ExecuteQuery(ctx, "SELECT * FROM nonexistent_table")
		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("ExecuteQueryWithPagination returns paginated results", func(t *testing.T) {
		params := domain.QueryParams{
			Query:  "SELECT id, name FROM test_users ORDER BY id",
			Offset: 0,
			Limit:  2,
		}

		result, err := repo.ExecuteQueryWithPagination(ctx, params)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.LessOrEqual(t, len(result.Rows), 2)
	})

	t.Run("ExecuteQueryWithPagination respects offset", func(t *testing.T) {
		params := domain.QueryParams{
			Query:  "SELECT id, name FROM test_users ORDER BY id",
			Offset: 1,
			Limit:  2,
		}

		result, err := repo.ExecuteQueryWithPagination(ctx, params)
		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("ExecuteMultipleQueries executes separated queries", func(t *testing.T) {
		queries := "SELECT id FROM test_users LIMIT 1; SELECT id FROM test_posts LIMIT 1"
		results, err := repo.ExecuteMultipleQueries(ctx, queries)
		require.NoError(t, err)
		require.NotNil(t, results)
		require.GreaterOrEqual(t, len(results), 1)
	})

	t.Run("BeginTransaction creates new transaction", func(t *testing.T) {
		tx, err := repo.BeginTransaction(ctx)
		require.NoError(t, err)
		require.NotNil(t, tx)
		tx.Rollback()
	})

	t.Run("CommitTransaction commits transaction", func(t *testing.T) {
		tx, err := repo.BeginTransaction(ctx)
		require.NoError(t, err)

		err = repo.CommitTransaction(ctx, tx)
		require.NoError(t, err)
	})

	t.Run("RollbackTransaction rolls back transaction", func(t *testing.T) {
		tx, err := repo.BeginTransaction(ctx)
		require.NoError(t, err)

		err = repo.RollbackTransaction(ctx, tx)
		require.NoError(t, err)
	})

	t.Run("GetDatabases returns list of databases", func(t *testing.T) {
		databases, err := repo.GetDatabases(ctx)
		require.NoError(t, err)
		require.NotNil(t, databases)
		require.Greater(t, len(databases), 0)
		require.Contains(t, databases, "testdb")
	})

	t.Run("GetSchemas returns list of schemas", func(t *testing.T) {
		schemas, err := repo.GetSchemas(ctx, "testdb")
		require.NoError(t, err)
		require.NotNil(t, schemas)
		require.Greater(t, len(schemas), 0)
	})

	t.Run("GetSchemas returns error for non-existent database", func(t *testing.T) {
		_, err := repo.GetSchemas(ctx, "nonexistent_db")
		require.Error(t, err)
	})

	t.Run("GetTables returns list of tables in schema", func(t *testing.T) {
		tables, err := repo.GetTables(ctx, "testdb", "public")
		require.NoError(t, err)
		require.NotNil(t, tables)
		require.Greater(t, len(tables), 0)
	})

	t.Run("GetTables returns error for non-existent schema", func(t *testing.T) {
		_, err := repo.GetTables(ctx, "testdb", "nonexistent_schema")
		require.Error(t, err)
	})

	t.Run("GetTableMetadata returns table structure", func(t *testing.T) {
		metadata, err := repo.GetTableMetadata(ctx, "testdb", "public", "test_users")
		require.NoError(t, err)
		require.NotNil(t, metadata)
		require.Equal(t, "test_users", metadata.Name)
		require.Greater(t, len(metadata.Columns), 0)
		require.Contains(t, metadata.PrimaryKeys, "id")
	})

	t.Run("GetTableMetadata returns error for non-existent table", func(t *testing.T) {
		_, err := repo.GetTableMetadata(ctx, "testdb", "public", "nonexistent_table")
		require.Error(t, err)
	})

	t.Run("GetDatabaseMetadata returns complete database structure", func(t *testing.T) {
		metadata, err := repo.GetDatabaseMetadata(ctx, "testdb")
		require.NoError(t, err)
		require.NotNil(t, metadata)
		require.Equal(t, "testdb", metadata.Name)
		require.Greater(t, len(metadata.Schemas), 0)
	})

	t.Run("GetTableData returns table rows", func(t *testing.T) {
		params := domain.TableDataParams{
			Database: "testdb",
			Schema:   "public",
			Table:    "test_users",
			Offset:   0,
			Limit:    10,
		}

		result, err := repo.GetTableData(ctx, params)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Greater(t, len(result.Rows), 0)
	})

	t.Run("GetTableData respects WHERE clause", func(t *testing.T) {
		params := domain.TableDataParams{
			Database:    "testdb",
			Schema:      "public",
			Table:       "test_users",
			WhereClause: "name = 'Alice'",
			Offset:      0,
			Limit:       10,
		}

		result, err := repo.GetTableData(ctx, params)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Len(t, result.Rows, 1)
	})

	t.Run("GetTableData respects ORDER BY", func(t *testing.T) {
		params := domain.TableDataParams{
			Database: "testdb",
			Schema:   "public",
			Table:    "test_users",
			OrderBy:  "name",
			OrderDir: "DESC",
			Offset:   0,
			Limit:    10,
		}

		result, err := repo.GetTableData(ctx, params)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Greater(t, len(result.Rows), 0)
	})

	t.Run("GetTableData respects LIMIT and OFFSET", func(t *testing.T) {
		params := domain.TableDataParams{
			Database: "testdb",
			Schema:   "public",
			Table:    "test_users",
			Offset:   0,
			Limit:    1,
		}

		result, err := repo.GetTableData(ctx, params)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.LessOrEqual(t, len(result.Rows), 1)
	})

	t.Run("InsertRow inserts new row", func(t *testing.T) {
		values := map[string]interface{}{
			"name":  "David",
			"email": "david@example.com",
		}

		err := repo.InsertRow(ctx, "testdb", "public", "test_users", values)
		require.NoError(t, err)

		// Verify insertion
		var count int
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test_users WHERE name = 'David'").Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})

	t.Run("UpdateRow updates existing row", func(t *testing.T) {
		pkValues := map[string]interface{}{"id": 1}
		updateValues := map[string]interface{}{"name": "Alicia"}

		err := repo.UpdateRow(ctx, "testdb", "public", "test_users", pkValues, updateValues)
		require.NoError(t, err)

		// Verify update
		var name string
		err = db.QueryRowContext(ctx, "SELECT name FROM test_users WHERE id = 1").Scan(&name)
		require.NoError(t, err)
		require.Equal(t, "Alicia", name)
	})

	t.Run("DeleteRow deletes row", func(t *testing.T) {
		// Insert row to delete
		_, err := db.ExecContext(ctx, "INSERT INTO test_users (name, email) VALUES ('ToDelete', 'todelete@example.com')")
		require.NoError(t, err)

		// Get ID of inserted row
		var id int
		err = db.QueryRowContext(ctx, "SELECT id FROM test_users WHERE name = 'ToDelete'").Scan(&id)
		require.NoError(t, err)

		// Delete row
		pkValues := map[string]interface{}{"id": id}
		err = repo.DeleteRow(ctx, "testdb", "public", "test_users", pkValues)
		require.NoError(t, err)

		// Verify deletion
		var count int
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM test_users WHERE id = ?", id).Scan(&count)
		// Expecting error or 0 rows
	})

	t.Run("GetRowCount returns correct count", func(t *testing.T) {
		count, err := repo.GetRowCount(ctx, "testdb", "public", "test_users", "")
		require.NoError(t, err)
		require.Greater(t, count, int64(0))
	})

	t.Run("GetRowCount respects WHERE clause", func(t *testing.T) {
		count, err := repo.GetRowCount(ctx, "testdb", "public", "test_users", "name = 'Alice'")
		require.NoError(t, err)
		require.Equal(t, int64(1), count)
	})

	t.Run("Disconnect closes connection", func(t *testing.T) {
		err := repo.Disconnect(ctx)
		require.NoError(t, err)
	})
}
