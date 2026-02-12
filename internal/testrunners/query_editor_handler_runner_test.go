package testrunners

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
	"github.com/kamil5b/lumen-pg/internal/implementations/mocks"
)

// QueryEditorHandlerConstructor creates a query editor handler with its dependencies
type QueryEditorHandlerConstructor func(queryUseCase usecase.QueryUseCase) usecase.QueryEditorHandler

// QueryEditorHandlerRunner runs test specs for query editor handler (Story 4 E2E)
func QueryEditorHandlerRunner(t *testing.T, constructor QueryEditorHandlerConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryUseCase := mocks.NewMockQueryUseCase(ctrl)
	handler := constructor(mockQueryUseCase)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	t.Run("E2E-S4-01: Query Editor Page Access", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/query-editor", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Header().Get("Content-Type"), "text/html")
	})

	t.Run("E2E-S4-02: Execute Single Query", func(t *testing.T) {
		queryReq := map[string]interface{}{
			"sql": "SELECT * FROM users",
		}

		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username"},
			Rows:       [][]interface{}{{1, "test"}},
			TotalRows:  1,
			LoadedRows: 1,
			Success:    true,
		}

		mockQueryUseCase.EXPECT().ExecuteQuery(gomock.Any(), "SELECT * FROM users").Return(expectedResult, nil)

		body, _ := json.Marshal(queryReq)
		req := httptest.NewRequest(http.MethodPost, "/query-editor/execute", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		
		var result domain.QueryResult
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Len(t, result.Rows, 1)
	})

	t.Run("E2E-S4-03: Execute Multiple Queries", func(t *testing.T) {
		queryReq := map[string]interface{}{
			"sql": "SELECT * FROM users; SELECT * FROM posts;",
		}

		expectedResults := []*domain.QueryResult{
			{
				Columns:    []string{"id", "username"},
				Rows:       [][]interface{}{{1, "test"}},
				TotalRows:  1,
				LoadedRows: 1,
				Success:    true,
			},
			{
				Columns:    []string{"id", "title"},
				Rows:       [][]interface{}{{1, "post1"}},
				TotalRows:  1,
				LoadedRows: 1,
				Success:    true,
			},
		}

		mockQueryUseCase.EXPECT().ExecuteMultipleQueries(gomock.Any(), "SELECT * FROM users; SELECT * FROM posts;").Return(expectedResults, nil)

		body, _ := json.Marshal(queryReq)
		req := httptest.NewRequest(http.MethodPost, "/query-editor/execute-multiple", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("E2E-S4-04: Query Error Display", func(t *testing.T) {
		queryReq := map[string]interface{}{
			"sql": "SELECT * FROM nonexistent",
		}

		expectedResult := &domain.QueryResult{
			Success:      false,
			ErrorMessage: "table does not exist",
		}

		mockQueryUseCase.EXPECT().ExecuteQuery(gomock.Any(), "SELECT * FROM nonexistent").Return(expectedResult, nil)

		body, _ := json.Marshal(queryReq)
		req := httptest.NewRequest(http.MethodPost, "/query-editor/execute", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		
		var result domain.QueryResult
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		require.NoError(t, err)
		assert.False(t, result.Success)
		assert.Contains(t, result.ErrorMessage, "does not exist")
	})

	t.Run("E2E-S4-05: Offset Pagination Results with Hard Limit", func(t *testing.T) {
		queryReq := map[string]interface{}{
			"sql": "SELECT * FROM large_table",
		}

		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "data"},
			Rows:       make([][]interface{}, 1000), // First 1000 rows
			TotalRows:  5000,                        // Total available
			LoadedRows: 1000,                        // Hard limit
			Success:    true,
		}

		mockQueryUseCase.EXPECT().ExecuteQuery(gomock.Any(), "SELECT * FROM large_table").Return(expectedResult, nil)

		body, _ := json.Marshal(queryReq)
		req := httptest.NewRequest(http.MethodPost, "/query-editor/execute", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		
		var result domain.QueryResult
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		require.NoError(t, err)
		assert.Equal(t, int64(5000), result.TotalRows) // Shows actual size
		assert.Equal(t, 1000, result.LoadedRows)       // Only 1000 loaded
	})
}
