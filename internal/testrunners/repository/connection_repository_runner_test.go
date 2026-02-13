package repository

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// ConnectionRepoConstructor creates a connection repository with database connection
type ConnectionRepoConstructor func(db *sql.DB) repository.ConnectionRepository

// ConnectionRepositoryRunner runs integration tests for connection repository (Story 1)
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

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	repo := constructor(db)

	t.Run("IT-S1-01: Parse Valid Connection String", func(t *testing.T) {
		validConnStr := "postgres://user:password@localhost:5432/testdb?sslmode=disable"

		conn, err := repo.ParseConnectionString(validConnStr)

		require.NoError(t, err)
		assert.NotNil(t, conn)
		assert.Equal(t, "localhost", conn.Host)
		assert.Equal(t, 5432, conn.Port)
		assert.Equal(t, "testdb", conn.Database)
		assert.Equal(t, "user", conn.User)
		assert.Equal(t, "password", conn.Password)
		assert.Equal(t, "disable", conn.SSLMode)
	})

	t.Run("IT-S1-02: Parse Connection String - Alternative Format", func(t *testing.T) {
		validConnStr := "postgres://admin:secret@db.example.com:5433/mydb?sslmode=require"

		conn, err := repo.ParseConnectionString(validConnStr)

		require.NoError(t, err)
		assert.NotNil(t, conn)
		assert.Equal(t, "db.example.com", conn.Host)
		assert.Equal(t, 5433, conn.Port)
		assert.Equal(t, "mydb", conn.Database)
		assert.Equal(t, "admin", conn.User)
		assert.Equal(t, "secret", conn.Password)
		assert.Equal(t, "require", conn.SSLMode)
	})

	t.Run("IT-S1-03: Parse Invalid Connection String", func(t *testing.T) {
		invalidConnStr := "invalid connection string"

		conn, err := repo.ParseConnectionString(invalidConnStr)

		require.Error(t, err)
		assert.Nil(t, conn)
	})

	t.Run("IT-S1-04: Validate Real Connection - Success", func(t *testing.T) {
		conn := &domain.Connection{
			Host:     "localhost",
			Port:     5432,
			Database: "testdb",
			User:     "postgres",
			Password: "postgres",
			SSLMode:  "disable",
		}

		err := repo.ValidateConnection(ctx, conn)

		require.NoError(t, err)
	})

	t.Run("IT-S1-05: Validate Connection - Wrong Password", func(t *testing.T) {
		conn := &domain.Connection{
			Host:     "localhost",
			Port:     5432,
			Database: "testdb",
			User:     "postgres",
			Password: "wrongpassword",
			SSLMode:  "disable",
		}

		err := repo.ValidateConnection(ctx, conn)

		require.Error(t, err)
	})

	t.Run("IT-S1-06: Validate Connection - Non-existent Database", func(t *testing.T) {
		conn := &domain.Connection{
			Host:     "localhost",
			Port:     5432,
			Database: "nonexistent",
			User:     "postgres",
			Password: "postgres",
			SSLMode:  "disable",
		}

		err := repo.ValidateConnection(ctx, conn)

		require.Error(t, err)
	})

	t.Run("IT-S1-07: Connect to PostgreSQL - Success", func(t *testing.T) {
		conn := &domain.Connection{
			Host:     "localhost",
			Port:     5432,
			Database: "testdb",
			User:     "postgres",
			Password: "postgres",
			SSLMode:  "disable",
		}

		dbConn, err := repo.Connect(ctx, conn)

		require.NoError(t, err)
		assert.NotNil(t, dbConn)
	})

	t.Run("IT-S1-08: Connection String with Special Characters", func(t *testing.T) {
		validConnStr := "postgres://user:p%40ssw0rd%23@localhost:5432/testdb?sslmode=disable"

		conn, err := repo.ParseConnectionString(validConnStr)

		require.NoError(t, err)
		assert.NotNil(t, conn)
		// Password should be URL-decoded
		assert.Contains(t, conn.Password, "@") // Special char decoded
	})

	t.Run("IT-S1-09: Connection String Missing Port", func(t *testing.T) {
		validConnStr := "postgres://user:password@localhost/testdb?sslmode=disable"

		conn, err := repo.ParseConnectionString(validConnStr)

		require.NoError(t, err)
		assert.NotNil(t, conn)
		// Should use default PostgreSQL port
		assert.Equal(t, 5432, conn.Port)
	})

	t.Run("IT-S1-10: Parse Multiple Connections Independently", func(t *testing.T) {
		connStr1 := "postgres://user1:pass1@host1:5432/db1?sslmode=disable"
		connStr2 := "postgres://user2:pass2@host2:5433/db2?sslmode=require"

		conn1, err1 := repo.ParseConnectionString(connStr1)
		conn2, err2 := repo.ParseConnectionString(connStr2)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, conn1.User, conn2.User)
		assert.NotEqual(t, conn1.Host, conn2.Host)
		assert.NotEqual(t, conn1.Port, conn2.Port)
	})
}
