package e2e_integration

import (
	"net/http"
	"testing"
)

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
//       router := setupRouter() // Your actual router with all middleware
//       e2e.RunAllE2ETests(t, router)
//   }

// RunAllE2ETests executes all E2E test runners for all stories
func RunAllE2ETests(t *testing.T, router http.Handler) {
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
func RunAuthenticationE2ETests(t *testing.T, router http.Handler) {
	t.Helper()
	Story2AuthE2ERunner(t, router)
}

// RunDataInteractionE2ETests runs data interaction E2E tests (ERD, Query Editor, Main View)
func RunDataInteractionE2ETests(t *testing.T, router http.Handler) {
	t.Helper()

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
func RunSecurityE2ETests(t *testing.T, router http.Handler) {
	t.Helper()

	t.Run("Story 6: Isolation", func(t *testing.T) {
		Story6IsolationE2ERunner(t, router)
	})

	t.Run("Story 7: Security & Best Practices", func(t *testing.T) {
		Story7SecurityE2ERunner(t, router)
	})
}

// RunStory runs a specific story's E2E tests
func RunStory(t *testing.T, router http.Handler, storyNumber int) {
	t.Helper()

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
