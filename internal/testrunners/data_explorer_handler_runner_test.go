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

// DataExplorerHandlerConstructor creates a data explorer handler with its dependencies
type DataExplorerHandlerConstructor func(
	queryUseCase usecase.QueryUseCase,
	transactionUseCase usecase.TransactionUseCase,
	metadataUseCase usecase.MetadataUseCase,
) usecase.DataExplorerHandler

// DataExplorerHandlerRunner runs test specs for data explorer handler (Story 5 E2E)
func DataExplorerHandlerRunner(t *testing.T, constructor DataExplorerHandlerConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryUseCase := mocks.NewMockQueryUseCase(ctrl)
	mockTransactionUseCase := mocks.NewMockTransactionUseCase(ctrl)
	mockMetadataUseCase := mocks.NewMockMetadataUseCase(ctrl)
	
	handler := constructor(mockQueryUseCase, mockTransactionUseCase, mockMetadataUseCase)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	t.Run("E2E-S5-01: Main View Default Load", func(t *testing.T) {
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username"},
			Rows:       [][]interface{}{{1, "test"}},
			TotalRows:  1,
			LoadedRows: 1,
			Success:    true,
		}

		mockQueryUseCase.EXPECT().ExecuteQuery(gomock.Any(), gomock.Any()).Return(expectedResult, nil)

		req := httptest.NewRequest(http.MethodGet, "/data-explorer/table/users", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("E2E-S5-02: Table Selection from Sidebar", func(t *testing.T) {
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "title"},
			Rows:       [][]interface{}{{1, "post1"}},
			TotalRows:  1,
			LoadedRows: 1,
			Success:    true,
		}

		mockQueryUseCase.EXPECT().ExecuteQuery(gomock.Any(), gomock.Any()).Return(expectedResult, nil)

		req := httptest.NewRequest(http.MethodGet, "/data-explorer/table/posts", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("E2E-S5-03: WHERE Bar Filtering", func(t *testing.T) {
		whereReq := map[string]interface{}{
			"table":  "users",
			"where":  "id > 10",
			"limit":  50,
			"offset": 0,
		}

		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username"},
			Rows:       [][]interface{}{{11, "user11"}, {12, "user12"}},
			TotalRows:  2,
			LoadedRows: 2,
			Success:    true,
		}

		mockQueryUseCase.EXPECT().ExecuteQuery(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedResult, nil)

		body, _ := json.Marshal(whereReq)
		req := httptest.NewRequest(http.MethodPost, "/data-explorer/filter", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("E2E-S5-06: Start Transaction Button", func(t *testing.T) {
		txnReq := map[string]interface{}{
			"username":  "testuser",
			"tableName": "users",
		}

		expectedTxn := &domain.Transaction{
			ID:        "txn-123",
			Username:  "testuser",
			TableName: "users",
		}

		mockTransactionUseCase.EXPECT().StartTransaction(gomock.Any(), "testuser", "users").Return(expectedTxn, nil)

		body, _ := json.Marshal(txnReq)
		req := httptest.NewRequest(http.MethodPost, "/data-explorer/transaction/start", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		
		var txn domain.Transaction
		err := json.Unmarshal(rec.Body.Bytes(), &txn)
		require.NoError(t, err)
		assert.Equal(t, "txn-123", txn.ID)
	})

	t.Run("E2E-S5-09: Transaction Commit Button", func(t *testing.T) {
		commitReq := map[string]interface{}{
			"transactionId": "txn-123",
		}

		mockTransactionUseCase.EXPECT().CommitTransaction(gomock.Any(), "txn-123").Return(nil)

		body, _ := json.Marshal(commitReq)
		req := httptest.NewRequest(http.MethodPost, "/data-explorer/transaction/commit", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("E2E-S5-10: Transaction Rollback Button", func(t *testing.T) {
		rollbackReq := map[string]interface{}{
			"transactionId": "txn-123",
		}

		mockTransactionUseCase.EXPECT().RollbackTransaction(gomock.Any(), "txn-123").Return(nil)

		body, _ := json.Marshal(rollbackReq)
		req := httptest.NewRequest(http.MethodPost, "/data-explorer/transaction/rollback", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("E2E-S5-05: Cursor Pagination Infinite Scroll with Actual Size", func(t *testing.T) {
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "data"},
			Rows:       make([][]interface{}, 50), // First page: 50 rows
			TotalRows:  5000,                      // Total available
			LoadedRows: 50,
			Success:    true,
		}

		mockQueryUseCase.EXPECT().ExecuteQuery(gomock.Any(), gomock.Any()).Return(expectedResult, nil)

		req := httptest.NewRequest(http.MethodGet, "/data-explorer/table/large_table?limit=50&offset=0", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		
		var result domain.QueryResult
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		require.NoError(t, err)
		assert.Equal(t, int64(5000), result.TotalRows) // Shows total
		assert.Equal(t, 50, result.LoadedRows)         // Current page
	})
}
