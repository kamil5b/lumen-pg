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
	"github.com/kamil5b/lumen-pg/internal/implementations/mocks"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
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

	t.Run("E2E-S5-04: Column Header Sorting", func(t *testing.T) {
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username"},
			Rows:       [][]interface{}{{1, "alice"}, {2, "bob"}},
			TotalRows:  2,
			LoadedRows: 2,
			Success:    true,
		}

		mockQueryUseCase.EXPECT().ExecuteQuery(gomock.Any(), gomock.Any()).Return(expectedResult, nil)

		req := httptest.NewRequest(http.MethodGet, "/data-explorer/table/users?sort=username&order=asc", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("E2E-S5-05a: Cursor Pagination Infinite Scroll Loading", func(t *testing.T) {
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "username"},
			Rows:       make([][]interface{}, 50),
			TotalRows:  1000,
			LoadedRows: 50,
			Success:    true,
		}

		mockQueryUseCase.EXPECT().ExecuteQuery(gomock.Any(), gomock.Any()).Return(expectedResult, nil)

		req := httptest.NewRequest(http.MethodGet, "/data-explorer/table/users?limit=50&offset=50", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var result domain.QueryResult
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		require.NoError(t, err)
		assert.Equal(t, 50, result.LoadedRows)
	})

	t.Run("E2E-S5-05b: Pagination Hard Limit Enforcement", func(t *testing.T) {
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "data"},
			Rows:       make([][]interface{}, 1000),
			TotalRows:  1000000,
			LoadedRows: 1000,
			Success:    true,
		}

		mockQueryUseCase.EXPECT().ExecuteQuery(gomock.Any(), gomock.Any()).Return(expectedResult, nil)

		req := httptest.NewRequest(http.MethodGet, "/data-explorer/table/huge_table", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var result domain.QueryResult
		err := json.Unmarshal(rec.Body.Bytes(), &result)
		require.NoError(t, err)
		assert.Equal(t, 1000, result.LoadedRows)
		assert.Greater(t, result.TotalRows, int64(1000))
	})

	t.Run("E2E-S5-07: Transaction Mode Cell Editing", func(t *testing.T) {
		editReq := map[string]interface{}{
			"transactionId": "txn-123",
			"rowId":         1,
			"column":        "username",
			"value":         "newusername",
		}

		mockTransactionUseCase.EXPECT().BufferEdit(gomock.Any(), "txn-123", gomock.Any()).Return(nil)

		body, _ := json.Marshal(editReq)
		req := httptest.NewRequest(http.MethodPost, "/data-explorer/transaction/edit", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("E2E-S5-08: Transaction Mode Edit Buffer Display", func(t *testing.T) {
		expectedTxn := &domain.Transaction{
			ID:        "txn-123",
			Username:  "testuser",
			TableName: "users",
			Operations: []domain.TransactionOperation{
				{
					Type:       domain.OperationUpdate,
					PrimaryKey: 1,
					Column:     "username",
					OldValue:   "oldname",
					NewValue:   "newname",
				},
			},
		}

		mockTransactionUseCase.EXPECT().GetTransaction(gomock.Any(), "txn-123").Return(expectedTxn, nil)

		req := httptest.NewRequest(http.MethodGet, "/data-explorer/transaction/txn-123", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("E2E-S5-11: Transaction Timer Countdown", func(t *testing.T) {
		expectedTxn := &domain.Transaction{
			ID:        "txn-123",
			Username:  "testuser",
			TableName: "users",
		}

		mockTransactionUseCase.EXPECT().GetTransaction(gomock.Any(), "txn-123").Return(expectedTxn, nil)

		req := httptest.NewRequest(http.MethodGet, "/data-explorer/transaction/txn-123", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("E2E-S5-12: Transaction Row Delete Button", func(t *testing.T) {
		deleteReq := map[string]interface{}{
			"transactionId": "txn-123",
			"rowId":         1,
		}

		mockTransactionUseCase.EXPECT().BufferEdit(gomock.Any(), "txn-123", gomock.Any()).Return(nil)

		body, _ := json.Marshal(deleteReq)
		req := httptest.NewRequest(http.MethodPost, "/data-explorer/transaction/delete-row", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("E2E-S5-13: Transaction New Row Button", func(t *testing.T) {
		insertReq := map[string]interface{}{
			"transactionId": "txn-123",
			"rowData": map[string]interface{}{
				"username": "newuser",
				"email":    "new@test.com",
			},
		}

		mockTransactionUseCase.EXPECT().BufferEdit(gomock.Any(), "txn-123", gomock.Any()).Return(nil)

		body, _ := json.Marshal(insertReq)
		req := httptest.NewRequest(http.MethodPost, "/data-explorer/transaction/insert-row", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("E2E-S5-14: FK Cell Navigation (Read-Only)", func(t *testing.T) {
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "user_id"},
			Rows:       [][]interface{}{{1, 1}},
			TotalRows:  1,
			LoadedRows: 1,
			Success:    true,
		}

		mockQueryUseCase.EXPECT().ExecuteQuery(gomock.Any(), gomock.Any()).Return(expectedResult, nil)

		req := httptest.NewRequest(http.MethodGet, "/data-explorer/table/posts/fk/1/users", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("E2E-S5-15: PK Cell Navigation (Read-Only)", func(t *testing.T) {
		expectedResult := &domain.QueryResult{
			Columns:    []string{"table_name", "row_count"},
			Rows:       [][]interface{}{{"posts", 5}, {"comments", 3}},
			TotalRows:  2,
			LoadedRows: 2,
			Success:    true,
		}

		mockMetadataUseCase.EXPECT().LoadGlobalMetadata(gomock.Any()).Return(&domain.GlobalMetadata{}, nil)
		mockQueryUseCase.EXPECT().ExecuteQuery(gomock.Any(), gomock.Any()).Return(expectedResult, nil)

		req := httptest.NewRequest(http.MethodGet, "/data-explorer/table/users/pk/1", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("E2E-S5-15a: PK Cell Navigation - Table Click", func(t *testing.T) {
		expectedResult := &domain.QueryResult{
			Columns:    []string{"id", "user_id", "title"},
			Rows:       [][]interface{}{{1, 1, "Post 1"}, {2, 1, "Post 2"}},
			TotalRows:  2,
			LoadedRows: 2,
			Success:    true,
		}

		mockQueryUseCase.EXPECT().ExecuteQuery(gomock.Any(), gomock.Any()).Return(expectedResult, nil)

		req := httptest.NewRequest(http.MethodGet, "/data-explorer/table/posts?user_id=1", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
