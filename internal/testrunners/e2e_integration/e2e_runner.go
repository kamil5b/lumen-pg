package e2e_integration

import (
	"context"
	"database/sql"
	"net/http"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// RouterConstructor is a function type that creates an http.Handler (router) with all middleware
// given a database connection. This allows E2E tests to build the complete application stack.
type RouterConstructor func(db *sql.DB) http.Handler

// E2ETestRunner orchestrates all end-to-end test runners for the Lumen-PG application
// This runs complete route-level tests with full middleware stack
// Maps to TEST_PLAN.md Phase 6: E2E Tests [L233-281]
//
// Unlike handler test runners (which mock use cases) and integration tests
// (which test repository implementations), E2E tests verify:
// - Complete HTTP request/response flow through all middleware
// - Session management and cookie handling across requests
// - Authentication and authorization enforcement
// - Multi-user isolation and concurrent request handling
// - Security features (SQL injection prevention, cookie security, etc.)
//
// Usage:
//   func TestE2ERoutes(t *testing.T) {
//       e2e.RunAllE2ETests(t, func(db *sql.DB) http.Handler {
//           // Build your router with all middleware using the provided DB
//           return setupRouter(db)
//       })
//   }

// RunAllE2ETests executes all E2E test runners for all stories
// It automatically sets up a testcontainer PostgreSQL database and tears it down after tests
func RunAllE2ETests(t *testing.T, constructor RouterConstructor) {
	t.Helper()

	ctx := context.Background()

	// Setup PostgreSQL testcontainer
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

	// Seed E2E test data
	seedE2ETestData(t, ctx, db)

	// Build router with the test database
	router := constructor(db)

	// Run all E2E test suites
	t.Run("Story 2: Authentication & Identity", func(t *testing.T) {
		Story2AuthE2ERunner(t, router)
	})

	t.Run("Story 3: ERD Viewer", func(t *testing.T) {
		Story3ERDViewerE2ERunner(t, router)
	})

	t.Run("Story 4: Manual Query Editor", func(t *testing.T) {
		Story4QueryEditorE2ERunner(t, router)
	})

	t.Run("Story 5: Main View & Data Interaction", func(t *testing.T) {
		Story5MainViewE2ERunner(t, router)
	})

	t.Run("Story 6: Isolation", func(t *testing.T) {
		Story6IsolationE2ERunner(t, router)
	})

	t.Run("Story 7: Security & Best Practices", func(t *testing.T) {
		Story7SecurityE2ERunner(t, router)
	})
}

// RunAllE2ETestsWithRouter is a convenience function for tests that want to provide
// a pre-configured router (e.g., for testing against a specific database setup)
// This is the old signature maintained for backward compatibility
func RunAllE2ETestsWithRouter(t *testing.T, router http.Handler) {
	t.Helper()

	t.Run("Story 2: Authentication & Identity", func(t *testing.T) {
		Story2AuthE2ERunner(t, router)
	})

	t.Run("Story 3: ERD Viewer", func(t *testing.T) {
		Story3ERDViewerE2ERunner(t, router)
	})

	t.Run("Story 4: Manual Query Editor", func(t *testing.T) {
		Story4QueryEditorE2ERunner(t, router)
	})

	t.Run("Story 5: Main View & Data Interaction", func(t *testing.T) {
		Story5MainViewE2ERunner(t, router)
	})

	t.Run("Story 6: Isolation", func(t *testing.T) {
		Story6IsolationE2ERunner(t, router)
	})

	t.Run("Story 7: Security & Best Practices", func(t *testing.T) {
		Story7SecurityE2ERunner(t, router)
	})
}

// RunAuthenticationE2ETests runs only authentication-related E2E tests
func RunAuthenticationE2ETests(t *testing.T, constructor RouterConstructor) {
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

	seedE2ETestData(t, ctx, db)

	router := constructor(db)

	Story2AuthE2ERunner(t, router)
}

// RunDataInteractionE2ETests runs data interaction E2E tests (ERD, Query Editor, Main View)
func RunDataInteractionE2ETests(t *testing.T, constructor RouterConstructor) {
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

	seedE2ETestData(t, ctx, db)

	router := constructor(db)

	t.Run("Story 3: ERD Viewer", func(t *testing.T) {
		Story3ERDViewerE2ERunner(t, router)
	})

	t.Run("Story 4: Manual Query Editor", func(t *testing.T) {
		Story4QueryEditorE2ERunner(t, router)
	})

	t.Run("Story 5: Main View & Data Interaction", func(t *testing.T) {
		Story5MainViewE2ERunner(t, router)
	})
}

// RunSecurityE2ETests runs security-focused E2E tests (Isolation, Security)
func RunSecurityE2ETests(t *testing.T, constructor RouterConstructor) {
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

	seedE2ETestData(t, ctx, db)

	router := constructor(db)

	t.Run("Story 6: Isolation", func(t *testing.T) {
		Story6IsolationE2ERunner(t, router)
	})

	t.Run("Story 7: Security & Best Practices", func(t *testing.T) {
		Story7SecurityE2ERunner(t, router)
	})
}

// RunStory runs a specific story's E2E tests
func RunStory(t *testing.T, constructor RouterConstructor, storyNumber int) {
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

	seedE2ETestData(t, ctx, db)

	router := constructor(db)

	switch storyNumber {
	case 2:
		t.Run("Story 2: Authentication & Identity", func(t *testing.T) {
			Story2AuthE2ERunner(t, router)
		})
	case 3:
		t.Run("Story 3: ERD Viewer", func(t *testing.T) {
			Story3ERDViewerE2ERunner(t, router)
		})
	case 4:
		t.Run("Story 4: Manual Query Editor", func(t *testing.T) {
			Story4QueryEditorE2ERunner(t, router)
		})
	case 5:
		t.Run("Story 5: Main View & Data Interaction", func(t *testing.T) {
			Story5MainViewE2ERunner(t, router)
		})
	case 6:
		t.Run("Story 6: Isolation", func(t *testing.T) {
			Story6IsolationE2ERunner(t, router)
		})
	case 7:
		t.Run("Story 7: Security & Best Practices", func(t *testing.T) {
			Story7SecurityE2ERunner(t, router)
		})
	default:
		t.Fatalf("Unknown story number: %d", storyNumber)
	}
}

// seedE2ETestData creates the necessary database schema and test data for E2E tests
// This includes test users with different permissions, test tables, and sample data
func seedE2ETestData(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()

	// Create test tables for E2E scenarios
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(200) NOT NULL,
			content TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			price DECIMAL(10, 2) NOT NULL,
			description TEXT,
			stock INTEGER DEFAULT 0
		)
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
			quantity INTEGER NOT NULL,
			total_price DECIMAL(10, 2) NOT NULL,
			order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err)

	// Insert test users for different E2E scenarios
	_, err = db.ExecContext(ctx, `
		INSERT INTO users (id, name, email) VALUES
		(1, 'Alice', 'alice@example.com'),
		(2, 'Bob', 'bob@example.com'),
		(3, 'Charlie', 'charlie@example.com'),
		(4, 'David', 'david@example.com')
		ON CONFLICT (id) DO NOTHING
	`)
	require.NoError(t, err)

	// Insert test posts
	_, err = db.ExecContext(ctx, `
		INSERT INTO posts (id, user_id, title, content) VALUES
		(1, 1, 'First Post', 'This is Alice first post'),
		(2, 1, 'Second Post', 'This is Alice second post'),
		(3, 2, 'Bob Introduction', 'Hello, I am Bob'),
		(4, 3, 'Charlie Thoughts', 'Some thoughts from Charlie')
		ON CONFLICT (id) DO NOTHING
	`)
	require.NoError(t, err)

	// Insert test comments
	_, err = db.ExecContext(ctx, `
		INSERT INTO comments (id, post_id, user_id, content) VALUES
		(1, 1, 2, 'Nice post Alice!'),
		(2, 1, 3, 'Great content'),
		(3, 2, 2, 'Looking forward to more'),
		(4, 3, 1, 'Welcome Bob!')
		ON CONFLICT (id) DO NOTHING
	`)
	require.NoError(t, err)

	// Insert test products
	_, err = db.ExecContext(ctx, `
		INSERT INTO products (id, name, price, description, stock) VALUES
		(1, 'Laptop', 999.99, 'High-performance laptop', 10),
		(2, 'Mouse', 29.99, 'Wireless mouse', 50),
		(3, 'Keyboard', 79.99, 'Mechanical keyboard', 30),
		(4, 'Monitor', 299.99, '27-inch 4K monitor', 15)
		ON CONFLICT (id) DO NOTHING
	`)
	require.NoError(t, err)

	// Insert test orders
	_, err = db.ExecContext(ctx, `
		INSERT INTO orders (id, user_id, product_id, quantity, total_price) VALUES
		(1, 1, 1, 1, 999.99),
		(2, 1, 2, 2, 59.98),
		(3, 2, 3, 1, 79.99),
		(4, 3, 4, 1, 299.99)
		ON CONFLICT (id) DO NOTHING
	`)
	require.NoError(t, err)

	// Reset sequences to avoid conflicts
	_, err = db.ExecContext(ctx, `
		SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));
		SELECT setval('posts_id_seq', (SELECT MAX(id) FROM posts));
		SELECT setval('comments_id_seq', (SELECT MAX(id) FROM comments));
		SELECT setval('products_id_seq', (SELECT MAX(id) FROM products));
		SELECT setval('orders_id_seq', (SELECT MAX(id) FROM orders));
	`)
	require.NoError(t, err)
}
