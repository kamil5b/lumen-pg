package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/handler"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
	mockUsecase "github.com/kamil5b/lumen-pg/internal/testrunners/mocks/usecase"
)

// ERDViewerHandlerConstructor is a function type that creates an ERDViewerHandler
type ERDViewerHandlerConstructor func(
	erdUC usecase.ERDUseCase,
	authUC usecase.AuthenticationUseCase,
) handler.ERDViewerHandler

// ERDViewerHandlerRunner runs all ERD viewer E2E handler tests
// Maps to TEST_PLAN.md:
// - Story 3: ERD Viewer [E2E-S3-01~04]
func ERDViewerHandlerRunner(t *testing.T, constructor ERDViewerHandlerConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockERD := mockUsecase.NewMockERDUseCase(ctrl)
	mockAuth := mockUsecase.NewMockAuthenticationUseCase(ctrl)

	h := constructor(mockERD, mockAuth)

	// E2E-S3-01: ERD Viewer Page Access
	t.Run("E2E-S3-01: ERD Viewer Page Access", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockERD.EXPECT().
			GetAvailableSchemas(gomock.Any(), "testuser", "testdb").
			Return([]string{"public", "private"}, nil)

		mockERD.EXPECT().
			GenerateERD(gomock.Any(), "testuser", "testdb", "public").
			Return(&domain.ERDData{
				Tables: []domain.ERDTable{
					{
						Name: "users",
						Columns: []domain.ERDColumn{
							{Name: "id", DataType: "integer", IsPrimary: true},
							{Name: "name", DataType: "text", IsPrimary: false},
						},
					},
					{
						Name: "posts",
						Columns: []domain.ERDColumn{
							{Name: "id", DataType: "integer", IsPrimary: true},
							{Name: "user_id", DataType: "integer", IsPrimary: false},
						},
					},
				},
				Relationships: []domain.ERDRelationship{
					{
						FromTable:  "posts",
						FromColumn: "user_id",
						ToTable:    "users",
						ToColumn:   "id",
					},
				},
			}, nil)

		req := httptest.NewRequest(http.MethodGet, "/erd?database=testdb&schema=public", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleERDViewerPage(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify ERD page is rendered with diagram
		require.Contains(t, body, "users")
		require.Contains(t, body, "posts")
		require.Contains(t, body, "erd-diagram")
		require.Contains(t, body, "erd-canvas")
	})

	// E2E-S3-02: ERD Zoom Controls
	t.Run("E2E-S3-02: ERD Zoom Controls", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockERD.EXPECT().
			GetAvailableSchemas(gomock.Any(), "testuser", "testdb").
			Return([]string{"public"}, nil)

		mockERD.EXPECT().
			GenerateERD(gomock.Any(), "testuser", "testdb", "public").
			Return(&domain.ERDData{
				Tables:        []domain.ERDTable{},
				Relationships: []domain.ERDRelationship{},
			}, nil)

		req := httptest.NewRequest(http.MethodGet, "/erd?database=testdb&schema=public", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleERDViewerPage(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify zoom controls are present
		require.Contains(t, body, "zoom-in")
		require.Contains(t, body, "zoom-out")
		require.Contains(t, body, "zoom-reset")
	})

	// E2E-S3-03: ERD Pan
	t.Run("E2E-S3-03: ERD Pan", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockERD.EXPECT().
			GetAvailableSchemas(gomock.Any(), "testuser", "testdb").
			Return([]string{"public"}, nil)

		mockERD.EXPECT().
			GenerateERD(gomock.Any(), "testuser", "testdb", "public").
			Return(&domain.ERDData{
				Tables:        []domain.ERDTable{},
				Relationships: []domain.ERDRelationship{},
			}, nil)

		req := httptest.NewRequest(http.MethodGet, "/erd?database=testdb&schema=public", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleERDViewerPage(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify pan functionality is enabled (via draggable canvas)
		require.Contains(t, body, "draggable")
		require.Contains(t, body, "pan")
	})

	// E2E-S3-04: Table Click in ERD
	t.Run("E2E-S3-04: Table Click in ERD", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockERD.EXPECT().
			GetTableBoxData(gomock.Any(), "testuser", "testdb", "public", "users").
			Return(&domain.TableMetadata{
				Name: "users",
				Columns: []domain.ColumnMetadata{
					{Name: "id", DataType: "integer", IsNullable: false, IsPrimary: true},
					{Name: "name", DataType: "text", IsNullable: true, IsPrimary: false},
					{Name: "email", DataType: "varchar", IsNullable: true, IsPrimary: false},
				},
				PrimaryKeys: []string{"id"},
				ForeignKeys: []domain.ForeignKeyMetadata{},
			}, nil)

		req := httptest.NewRequest(http.MethodGet, "/erd/table?database=testdb&schema=public&table=users", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleTableClickInERD(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify table details are shown in side panel
		require.Contains(t, body, "users")
		require.Contains(t, body, "id")
		require.Contains(t, body, "name")
		require.Contains(t, body, "email")
		require.Contains(t, body, "integer")
		require.Contains(t, body, "text")
	})

	// Additional test: ERD with empty schema
	t.Run("ERD with Empty Schema", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockERD.EXPECT().
			GetAvailableSchemas(gomock.Any(), "testuser", "testdb").
			Return([]string{"empty_schema"}, nil)

		mockERD.EXPECT().
			IsSchemaEmpty(gomock.Any(), "testuser", "testdb", "empty_schema").
			Return(true, nil)

		req := httptest.NewRequest(http.MethodGet, "/erd?database=testdb&schema=empty_schema", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleERDViewerPage(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify empty schema message
		require.Contains(t, body, "No tables found")
	})

	// Additional test: ERD with complex relationships
	t.Run("ERD with Complex Relationships", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockERD.EXPECT().
			GetAvailableSchemas(gomock.Any(), "testuser", "testdb").
			Return([]string{"public"}, nil)

		mockERD.EXPECT().
			GenerateERD(gomock.Any(), "testuser", "testdb", "public").
			Return(&domain.ERDData{
				Tables: []domain.ERDTable{
					{Name: "users", Columns: []domain.ERDColumn{{Name: "id", DataType: "integer", IsPrimary: true}}},
					{Name: "posts", Columns: []domain.ERDColumn{{Name: "id", DataType: "integer", IsPrimary: true}, {Name: "user_id", DataType: "integer"}}},
					{Name: "comments", Columns: []domain.ERDColumn{{Name: "id", DataType: "integer", IsPrimary: true}, {Name: "post_id", DataType: "integer"}, {Name: "user_id", DataType: "integer"}}},
				},
				Relationships: []domain.ERDRelationship{
					{FromTable: "posts", FromColumn: "user_id", ToTable: "users", ToColumn: "id"},
					{FromTable: "comments", FromColumn: "post_id", ToTable: "posts", ToColumn: "id"},
					{FromTable: "comments", FromColumn: "user_id", ToTable: "users", ToColumn: "id"},
				},
			}, nil)

		req := httptest.NewRequest(http.MethodGet, "/erd?database=testdb&schema=public", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleERDViewerPage(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify all tables and relationships are rendered
		require.Contains(t, body, "users")
		require.Contains(t, body, "posts")
		require.Contains(t, body, "comments")
	})

	// Additional test: Schema selection dropdown
	t.Run("Schema Selection Dropdown", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockERD.EXPECT().
			GetAvailableSchemas(gomock.Any(), "testuser", "testdb").
			Return([]string{"public", "private", "staging"}, nil)

		mockERD.EXPECT().
			GenerateERD(gomock.Any(), "testuser", "testdb", "public").
			Return(&domain.ERDData{
				Tables:        []domain.ERDTable{},
				Relationships: []domain.ERDRelationship{},
			}, nil)

		req := httptest.NewRequest(http.MethodGet, "/erd?database=testdb&schema=public", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleERDViewerPage(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify schema selection dropdown contains all schemas
		require.Contains(t, body, "public")
		require.Contains(t, body, "private")
		require.Contains(t, body, "staging")
	})

	// Additional test: Unauthorized access
	t.Run("Unauthorized Access to ERD Viewer", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "").
			Return(nil, domain.ValidationError{Field: "session", Message: "No session"})

		req := httptest.NewRequest(http.MethodGet, "/erd?database=testdb&schema=public", nil)
		rec := httptest.NewRecorder()

		h.HandleERDViewerPage(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
