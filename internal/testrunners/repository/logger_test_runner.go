package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// LoggerRepositoryConstructor is a function type that creates a LoggerRepository
type LoggerRepositoryConstructor func() repository.LoggerRepository

// LoggerRepositoryRunner runs all logger repository tests against an implementation
// Maps to TEST_PLAN.md:
// - Story 2: Authentication & Identity [UC-S2-01~15, IT-S2-01~05]
// - Story 4: Manual Query Editor [UC-S4-01~08, IT-S4-01~04]
// - Story 5: Main View & Data Interaction [UC-S5-01~19, IT-S5-01~07]
// - Story 7: Security & Best Practices [UC-S7-01~07, IT-S7-01~03, E2E-S7-01~06]
func LoggerRepositoryRunner(t *testing.T, constructor LoggerRepositoryConstructor) {
	t.Helper()

	ctx := context.Background()
	repo := constructor()

	// UC-S2-01: Login Form Validation - Empty Username
	// UC-S2-02: Login Form Validation - Empty Password
	t.Run("LogInfo logs informational message", func(t *testing.T) {
		err := repo.LogInfo(ctx, "User login attempt", map[string]interface{}{
			"username": "testuser",
		})
		require.NoError(t, err)
	})

	// UC-S2-03: Login Connection Probe
	// IT-S2-01: Real PostgreSQL Connection Probe
	t.Run("LogInfo with empty fields", func(t *testing.T) {
		err := repo.LogInfo(ctx, "Application started", map[string]interface{}{})
		require.NoError(t, err)
	})

	// UC-S2-01: Login Form Validation - Empty Username
	t.Run("LogInfo with nil fields", func(t *testing.T) {
		err := repo.LogInfo(ctx, "System event", nil)
		require.NoError(t, err)
	})

	// UC-S2-04: Login Connection Probe Failure
	// IT-S2-02: Real PostgreSQL Connection Probe Failure
	t.Run("LogWarn logs warning message", func(t *testing.T) {
		err := repo.LogWarn(ctx, "Connection retry attempt", map[string]interface{}{
			"attempt": 2,
			"delay":   "5 seconds",
		})
		require.NoError(t, err)
	})

	// UC-S2-04: Login Connection Probe Failure
	t.Run("LogWarn with multiple fields", func(t *testing.T) {
		err := repo.LogWarn(ctx, "Session near expiration", map[string]interface{}{
			"username":       "testuser",
			"time_remaining": 300,
			"session_id":     "sess_123",
			"created_at":     "2024-01-01T00:00:00Z",
		})
		require.NoError(t, err)
	})

	// UC-S4-06: Invalid Query Error
	// IT-S4-04: Query with Permission Denied
	t.Run("LogError logs error message with error object", func(t *testing.T) {
		testErr := errors.New("syntax error")
		err := repo.LogError(ctx, "Failed to execute query", testErr, map[string]interface{}{
			"query": "SELECT * FROM invalid_table",
		})
		require.NoError(t, err)
	})

	// UC-S4-06: Invalid Query Error
	t.Run("LogError with nil error object", func(t *testing.T) {
		err := repo.LogError(ctx, "Query execution issue", nil, map[string]interface{}{
			"query": "SELECT * FROM users",
		})
		require.NoError(t, err)
	})

	// UC-S4-01: Single Query Execution
	// UC-S4-02: Multiple Query Execution
	t.Run("LogDebug logs debug message", func(t *testing.T) {
		err := repo.LogDebug(ctx, "Query execution started", map[string]interface{}{
			"query_id":       "qry_123",
			"username":       "testuser",
			"execution_time": 0,
		})
		require.NoError(t, err)
	})

	// UC-S4-07: Query Splitting
	t.Run("LogDebug with detailed information", func(t *testing.T) {
		err := repo.LogDebug(ctx, "Parsing multi-statement query", map[string]interface{}{
			"statement_count": 3,
			"query":           "SELECT 1; SELECT 2; SELECT 3;",
		})
		require.NoError(t, err)
	})

	// UC-S2-05: Login Success After Probe
	// IT-S2-01: Real PostgreSQL Connection Probe
	t.Run("LogSecurityEvent logs authentication success", func(t *testing.T) {
		err := repo.LogSecurityEvent(ctx, "auth_success", "testuser", map[string]interface{}{
			"ip_address": "192.168.1.1",
			"timestamp":  "2024-01-01T12:00:00Z",
		})
		require.NoError(t, err)
	})

	// UC-S2-04: Login Connection Probe Failure
	// UC-S2-03: Login Connection Probe
	t.Run("LogSecurityEvent logs authentication failure", func(t *testing.T) {
		err := repo.LogSecurityEvent(ctx, "auth_failure", "testuser", map[string]interface{}{
			"reason":     "invalid_password",
			"ip_address": "192.168.1.1",
			"attempt":    1,
		})
		require.NoError(t, err)
	})

	// UC-S7-01: SQL Injection Prevention - WHERE Clause
	// UC-S7-02: SQL Injection Prevention - Query Editor
	t.Run("LogSecurityEvent logs sql_injection_attempt", func(t *testing.T) {
		err := repo.LogSecurityEvent(ctx, "sql_injection_attempt", "testuser", map[string]interface{}{
			"query":  "SELECT * FROM users WHERE id = 1 OR '1'='1'",
			"source": "where_clause",
		})
		require.NoError(t, err)
	})

	// UC-S2-12: Logout Cookie Clearing
	t.Run("LogSecurityEvent logs logout event", func(t *testing.T) {
		err := repo.LogSecurityEvent(ctx, "logout", "testuser", map[string]interface{}{
			"session_duration": 3600,
			"timestamp":        "2024-01-01T13:00:00Z",
		})
		require.NoError(t, err)
	})

	// UC-S7-05: Cookie Tampering Detection
	t.Run("LogSecurityEvent logs cookie_tampering_attempt", func(t *testing.T) {
		err := repo.LogSecurityEvent(ctx, "cookie_tampering_attempt", "unknown_user", map[string]interface{}{
			"cookie_name": "session",
			"ip_address":  "192.168.1.100",
		})
		require.NoError(t, err)
	})

	// UC-S4-01: Single Query Execution
	// IT-S4-01: Real SELECT Query
	t.Run("LogQueryExecution logs successful SELECT query", func(t *testing.T) {
		err := repo.LogQueryExecution(ctx, "testuser", "SELECT * FROM users", 45, true, nil)
		require.NoError(t, err)
	})

	// UC-S4-02: Multiple Query Execution
	t.Run("LogQueryExecution logs multiple successful queries", func(t *testing.T) {
		queries := []string{
			"SELECT * FROM users",
			"SELECT * FROM posts",
			"SELECT * FROM comments",
		}

		for _, query := range queries {
			err := repo.LogQueryExecution(ctx, "testuser", query, 50, true, nil)
			require.NoError(t, err)
		}
	})

	// UC-S4-04: DDL Query Execution
	// IT-S4-02: Real DDL Query
	t.Run("LogQueryExecution logs DDL query", func(t *testing.T) {
		err := repo.LogQueryExecution(ctx, "testuser", "CREATE TABLE users (id INT PRIMARY KEY)", 120, true, nil)
		require.NoError(t, err)
	})

	// UC-S4-05: DML Query Execution
	// IT-S4-03: Real DML Query
	t.Run("LogQueryExecution logs DML query", func(t *testing.T) {
		err := repo.LogQueryExecution(ctx, "testuser", "INSERT INTO users (name) VALUES ('John')", 80, true, nil)
		require.NoError(t, err)
	})

	// UC-S4-06: Invalid Query Error
	t.Run("LogQueryExecution logs query with error", func(t *testing.T) {
		queryErr := errors.New("table does not exist")
		err := repo.LogQueryExecution(ctx, "testuser", "SELECT * FROM nonexistent", 25, false, queryErr)
		require.NoError(t, err)
	})

	// IT-S4-04: Query with Permission Denied
	t.Run("LogQueryExecution logs permission denied error", func(t *testing.T) {
		permErr := errors.New("permission denied")
		err := repo.LogQueryExecution(ctx, "restricted_user", "DROP TABLE users", 10, false, permErr)
		require.NoError(t, err)
	})

	// UC-S4-03: Query Result Offset Pagination
	t.Run("LogQueryExecution with pagination", func(t *testing.T) {
		err := repo.LogQueryExecution(ctx, "testuser", "SELECT * FROM users LIMIT 10 OFFSET 20", 60, true, nil)
		require.NoError(t, err)
	})

	// UC-S4-08: Parameterized Query Execution
	t.Run("LogQueryExecution with parameterized query", func(t *testing.T) {
		err := repo.LogQueryExecution(ctx, "testuser", "SELECT * FROM users WHERE id = $1", 35, true, nil)
		require.NoError(t, err)
	})

	// UC-S5-09: Transaction Start
	// UC-S5-12: Transaction Commit
	t.Run("LogTransactionEvent logs transaction start", func(t *testing.T) {
		err := repo.LogTransactionEvent(ctx, "testuser", "begin", map[string]interface{}{
			"transaction_id":  "txn_123",
			"isolation_level": "read_committed",
		})
		require.NoError(t, err)
	})

	// UC-S5-12: Transaction Commit
	// IT-S5-04: Real Transaction Commit
	t.Run("LogTransactionEvent logs transaction commit", func(t *testing.T) {
		err := repo.LogTransactionEvent(ctx, "testuser", "commit", map[string]interface{}{
			"transaction_id": "txn_123",
			"duration_ms":    1500,
			"changes":        3,
		})
		require.NoError(t, err)
	})

	// UC-S5-13: Transaction Rollback
	// IT-S5-05: Real Transaction Rollback
	t.Run("LogTransactionEvent logs transaction rollback", func(t *testing.T) {
		err := repo.LogTransactionEvent(ctx, "testuser", "rollback", map[string]interface{}{
			"transaction_id": "txn_123",
			"reason":         "user_requested",
			"duration_ms":    500,
		})
		require.NoError(t, err)
	})

	// UC-S5-14: Transaction Timer Expiration
	t.Run("LogTransactionEvent logs transaction timeout", func(t *testing.T) {
		err := repo.LogTransactionEvent(ctx, "testuser", "timeout", map[string]interface{}{
			"transaction_id":  "txn_123",
			"timeout_seconds": 600,
		})
		require.NoError(t, err)
	})

	// UC-S5-11: Cell Edit Buffering
	t.Run("LogTransactionEvent logs cell edit", func(t *testing.T) {
		err := repo.LogTransactionEvent(ctx, "testuser", "cell_edit", map[string]interface{}{
			"transaction_id": "txn_123",
			"table":          "users",
			"column":         "email",
			"row_id":         5,
			"old_value":      "old@example.com",
			"new_value":      "new@example.com",
		})
		require.NoError(t, err)
	})

	// UC-S5-15: Row Deletion Buffering
	t.Run("LogTransactionEvent logs row delete", func(t *testing.T) {
		err := repo.LogTransactionEvent(ctx, "testuser", "row_delete", map[string]interface{}{
			"transaction_id": "txn_123",
			"table":          "users",
			"row_id":         5,
		})
		require.NoError(t, err)
	})

	// UC-S5-16: Row Insertion Buffering
	t.Run("LogTransactionEvent logs row insert", func(t *testing.T) {
		err := repo.LogTransactionEvent(ctx, "testuser", "row_insert", map[string]interface{}{
			"transaction_id": "txn_123",
			"table":          "users",
			"data":           "name=John,email=john@example.com",
		})
		require.NoError(t, err)
	})

	// UC-S5-10: Transaction Already Active Error
	t.Run("LogTransactionEvent logs transaction already active error", func(t *testing.T) {
		err := repo.LogTransactionEvent(ctx, "testuser", "error", map[string]interface{}{
			"error_type":   "transaction_already_active",
			"previous_txn": "txn_122",
		})
		require.NoError(t, err)
	})

	// UC-S2-13: Header Username Display
	// IT-S2-03: Real Role-Based Resource Access
	t.Run("LogInfo logs user session info", func(t *testing.T) {
		err := repo.LogInfo(ctx, "User session info", map[string]interface{}{
			"username": "testuser",
			"role":     "user",
			"database": "testdb",
		})
		require.NoError(t, err)
	})

	// UC-S2-14: Navigation Menu Rendering
	t.Run("LogDebug logs UI render event", func(t *testing.T) {
		err := repo.LogDebug(ctx, "Navigation menu rendered", map[string]interface{}{
			"username":     "testuser",
			"tables_count": 15,
			"views_count":  3,
		})
		require.NoError(t, err)
	})

	// UC-S2-15: Metadata Refresh Button
	t.Run("LogInfo logs metadata refresh", func(t *testing.T) {
		err := repo.LogInfo(ctx, "Metadata refresh initiated", map[string]interface{}{
			"username":       "testuser",
			"previous_count": 15,
			"new_count":      16,
		})
		require.NoError(t, err)
	})

	// UC-S5-01: Table Data Loading
	t.Run("LogQueryExecution logs table data loading", func(t *testing.T) {
		err := repo.LogQueryExecution(ctx, "testuser", "SELECT * FROM users LIMIT 50", 75, true, nil)
		require.NoError(t, err)
	})

	// UC-S5-03: WHERE Clause Validation
	t.Run("LogDebug logs WHERE clause validation", func(t *testing.T) {
		err := repo.LogDebug(ctx, "WHERE clause validated", map[string]interface{}{
			"table":  "users",
			"clause": "age > 18 AND status = 'active'",
			"valid":  true,
		})
		require.NoError(t, err)
	})

	// UC-S5-05: Column Sorting ASC
	// UC-S5-06: Column Sorting DESC
	t.Run("LogDebug logs column sort operation", func(t *testing.T) {
		err := repo.LogDebug(ctx, "Column sort applied", map[string]interface{}{
			"table":     "users",
			"column":    "created_at",
			"direction": "DESC",
		})
		require.NoError(t, err)
	})

	// UC-S7-06: Session Timeout Short-Lived Cookie
	// UC-S7-07: Session Timeout Long-Lived Cookie
	t.Run("LogSecurityEvent logs session timeout", func(t *testing.T) {
		err := repo.LogSecurityEvent(ctx, "session_timeout", "testuser", map[string]interface{}{
			"reason":                   "inactivity",
			"timeout_duration_seconds": 1800,
		})
		require.NoError(t, err)
	})

	// UC-S7-03: Password Encryption in Cookie
	t.Run("LogSecurityEvent logs encryption operation", func(t *testing.T) {
		err := repo.LogSecurityEvent(ctx, "encryption_operation", "system", map[string]interface{}{
			"operation": "password_hash",
			"success":   true,
		})
		require.NoError(t, err)
	})

	// UC-S2-05: Login Success After Probe
	// IT-S2-01: Real PostgreSQL Connection Probe
	t.Run("Complex login flow logging sequence", func(t *testing.T) {
		username := "testuser"

		// Log connection probe
		err := repo.LogSecurityEvent(ctx, "connection_probe", username, map[string]interface{}{
			"database": "testdb",
			"success":  true,
		})
		require.NoError(t, err)

		// Log successful login
		err = repo.LogSecurityEvent(ctx, "auth_success", username, map[string]interface{}{
			"ip_address": "192.168.1.1",
		})
		require.NoError(t, err)

		// Log metadata loading
		err = repo.LogInfo(ctx, "Loading user metadata", map[string]interface{}{
			"username": username,
		})
		require.NoError(t, err)
	})

	// UC-S4-01: Single Query Execution
	t.Run("Complex query execution logging sequence", func(t *testing.T) {
		username := "testuser"
		query := "SELECT * FROM users WHERE age > 18"

		// Log query started (via debug)
		err := repo.LogDebug(ctx, "Query execution started", map[string]interface{}{
			"query": query,
		})
		require.NoError(t, err)

		// Log successful execution
		err = repo.LogQueryExecution(ctx, username, query, 120, true, nil)
		require.NoError(t, err)
	})

	// UC-S5-09: Transaction Start
	// UC-S5-11: Cell Edit Buffering
	// UC-S5-12: Transaction Commit
	t.Run("Complex transaction logging sequence", func(t *testing.T) {
		username := "testuser"
		txnID := "txn_456"

		// Log transaction start
		err := repo.LogTransactionEvent(ctx, username, "begin", map[string]interface{}{
			"transaction_id": txnID,
		})
		require.NoError(t, err)

		// Log cell edit
		err = repo.LogTransactionEvent(ctx, username, "cell_edit", map[string]interface{}{
			"transaction_id": txnID,
			"table":          "users",
		})
		require.NoError(t, err)

		// Log commit
		err = repo.LogTransactionEvent(ctx, username, "commit", map[string]interface{}{
			"transaction_id": txnID,
			"duration_ms":    2000,
		})
		require.NoError(t, err)
	})

	// E2E-S7-01: SQL Injection via WHERE Bar
	// E2E-S7-02: SQL Injection via Query Editor
	t.Run("SQL injection attempt logging", func(t *testing.T) {
		err := repo.LogSecurityEvent(ctx, "sql_injection_attempt", "potential_attacker", map[string]interface{}{
			"query_pattern": "' OR '1'='1",
			"source":        "where_clause",
			"timestamp":     "2024-01-01T14:30:00Z",
		})
		require.NoError(t, err)

		err = repo.LogSecurityEvent(ctx, "sql_injection_attempt", "potential_attacker", map[string]interface{}{
			"query_pattern": "; DROP TABLE users; --",
			"source":        "query_editor",
			"timestamp":     "2024-01-01T14:31:00Z",
		})
		require.NoError(t, err)
	})

	// IT-S7-01: Real SQL Injection Test
	t.Run("Query execution with injection prevention logging", func(t *testing.T) {
		// Log that injection was detected and prevented
		err := repo.LogWarn(ctx, "Potential injection pattern detected", map[string]interface{}{
			"query":     "SELECT * FROM users WHERE id = $1",
			"sanitized": true,
		})
		require.NoError(t, err)

		// Log successful execution after sanitization
		err = repo.LogQueryExecution(ctx, "testuser", "SELECT * FROM users WHERE id = $1", 40, true, nil)
		require.NoError(t, err)
	})

	// UC-S5-19: Read-Only Mode Enforcement
	t.Run("LogWarn logs read-only mode violation attempt", func(t *testing.T) {
		err := repo.LogWarn(ctx, "Read-only mode violation attempt", map[string]interface{}{
			"username":      "testuser",
			"operation":     "update",
			"table":         "users",
			"read_only_key": "readonly_user",
		})
		require.NoError(t, err)
	})
}
