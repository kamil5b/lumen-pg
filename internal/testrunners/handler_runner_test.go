package testrunners

import (
	"testing"

	"github.com/kamil5b/lumen-pg/internal/interfaces"
	"go.uber.org/mock/gomock"
)

// AuthHandlerConstructor is a function that creates an AuthHandler
type AuthHandlerConstructor func(authService interfaces.AuthService) interfaces.AuthHandler

// AuthHandlerRunner tests AuthHandler implementations
func AuthHandlerRunner(t *testing.T, constructor AuthHandlerConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("POST /login - success returns cookies and redirects", func(t *testing.T) {
		t.Skip("Requires mock implementation - run mockgen to generate mocks")
		// Test successful login with cookie setting
	})

	t.Run("POST /login - failure returns error message", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test failed login
	})

	t.Run("POST /login - no accessible resources", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test login with no accessible resources
	})

	t.Run("POST /logout - clears cookies and redirects", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test logout flow
	})

	t.Run("GET /login - displays login form", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test login page display
	})
}

// MainViewHandlerConstructor is a function that creates a MainViewHandler
type MainViewHandlerConstructor func(
	dataExplorerService interfaces.DataExplorerService,
	metadataService interfaces.MetadataService,
	transactionService interfaces.TransactionService,
) interfaces.MainViewHandler

// MainViewHandlerRunner tests MainViewHandler implementations
func MainViewHandlerRunner(t *testing.T, constructor MainViewHandlerConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("GET / - displays first accessible table", func(t *testing.T) {
		t.Skip("Requires mock implementation - run mockgen to generate mocks")
		// Test main view with first table loaded
	})

	t.Run("GET /table/:schema/:name - loads specific table", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test table selection
	})

	t.Run("POST /filter - applies WHERE clause", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test filtering
	})

	t.Run("GET /table/:schema/:name/sort - sorts by column", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test sorting
	})

	t.Run("GET /table/:schema/:name/page - infinite scroll pagination", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test pagination
	})

	t.Run("GET /pk-references - shows referencing tables", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test PK navigation
	})

	t.Run("GET /fk-navigate - navigates to parent table", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test FK navigation
	})
}

// QueryEditorHandlerConstructor is a function that creates a QueryEditorHandler
type QueryEditorHandlerConstructor func(queryService interfaces.QueryService) interfaces.QueryEditorHandler

// QueryEditorHandlerRunner tests QueryEditorHandler implementations
func QueryEditorHandlerRunner(t *testing.T, constructor QueryEditorHandlerConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("GET /query - displays query editor", func(t *testing.T) {
		t.Skip("Requires mock implementation - run mockgen to generate mocks")
		// Test query editor page
	})

	t.Run("POST /query/execute - executes single query", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test single query execution
	})

	t.Run("POST /query/execute - executes multiple queries", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test multiple query execution
	})

	t.Run("POST /query/execute - returns error for invalid query", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test error handling
	})

	t.Run("POST /query/execute - SELECT returns paginated results", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test SELECT with pagination (hard limit 1000 rows)
	})

	t.Run("POST /query/execute - DDL returns success message", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test DDL execution
	})

	t.Run("POST /query/execute - DML returns affected rows", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test DML execution
	})
}

// ERDViewerHandlerConstructor is a function that creates an ERDViewerHandler
type ERDViewerHandlerConstructor func(metadataService interfaces.MetadataService) interfaces.ERDViewerHandler

// ERDViewerHandlerRunner tests ERDViewerHandler implementations
func ERDViewerHandlerRunner(t *testing.T, constructor ERDViewerHandlerConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("GET /erd - displays ERD viewer page", func(t *testing.T) {
		t.Skip("Requires mock implementation - run mockgen to generate mocks")
		// Test ERD viewer page
	})

	t.Run("GET /erd/data - returns ERD JSON data", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test ERD data retrieval
	})

	t.Run("GET /erd/table/:schema/:name - returns table details", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test table details panel
	})
}

// TransactionHandlerConstructor is a function that creates a TransactionHandler
type TransactionHandlerConstructor func(transactionService interfaces.TransactionService) interfaces.TransactionHandler

// TransactionHandlerRunner tests TransactionHandler implementations
func TransactionHandlerRunner(t *testing.T, constructor TransactionHandlerConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("POST /transaction/start - starts transaction", func(t *testing.T) {
		t.Skip("Requires mock implementation - run mockgen to generate mocks")
		// Test transaction start
	})

	t.Run("POST /transaction/buffer - buffers operation", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test operation buffering
	})

	t.Run("POST /transaction/commit - commits operations", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test commit
	})

	t.Run("POST /transaction/rollback - discards operations", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test rollback
	})

	t.Run("GET /transaction/state - returns current state", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test state retrieval
	})

	t.Run("POST /transaction/edit-cell - buffers cell edit", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test cell editing in transaction mode
	})

	t.Run("POST /transaction/delete-row - buffers row delete", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test row deletion in transaction mode
	})

	t.Run("POST /transaction/insert-row - buffers row insert", func(t *testing.T) {
		t.Skip("Requires mock implementation")
		// Test row insertion in transaction mode
	})
}
