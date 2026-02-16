package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/handler"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
	mockUsecase "github.com/kamil5b/lumen-pg/internal/testrunners/mocks/usecase"
)

// TransactionHandlerConstructor is a function type that creates a TransactionHandler
type TransactionHandlerConstructor func(
	txnUC usecase.TransactionUseCase,
	authUC usecase.AuthenticationUseCase,
	rbacUC usecase.RBACUseCase,
) handler.TransactionHandler

// TransactionHandlerRunner runs all transaction E2E handler tests
// Maps to TEST_PLAN.md:
// - Story 5: Main View & Data Interaction [E2E-S5-06~13]
func TransactionHandlerRunner(t *testing.T, constructor TransactionHandlerConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockTxn := mockUsecase.NewMockTransactionUseCase(ctrl)
	mockAuth := mockUsecase.NewMockAuthenticationUseCase(ctrl)
	mockRBAC := mockUsecase.NewMockRBACUseCase(ctrl)

	h := constructor(mockTxn, mockAuth, mockRBAC)

	// E2E-S5-06: Start Transaction Button
	t.Run("E2E-S5-06: Start Transaction Button", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockRBAC.EXPECT().
			CheckUpdatePermission(gomock.Any(), "testuser", "testdb", "public", "users").
			Return(true, nil)

		mockTxn.EXPECT().
			CheckActiveTransaction(gomock.Any(), "testuser").
			Return(false, nil)

		mockTxn.EXPECT().
			StartTransaction(gomock.Any(), "testuser", "testdb", "public", "users").
			Return(&domain.TransactionState{
				ID:       "txn_123",
				Username: "testuser",
			}, nil)

		req := httptest.NewRequest(http.MethodPost, "/transaction/start?database=testdb&schema=public&table=users", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleStartTransaction(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify transaction started indicator
		require.Contains(t, body, "Transaction Active")
		require.Contains(t, body, "timer")
		require.Contains(t, body, "60")
	})

	// E2E-S5-07: Transaction Mode Cell Editing
	t.Run("E2E-S5-07: Transaction Mode Cell Editing", func(t *testing.T) {
		form := url.Values{}
		form.Add("row_index", "0")
		form.Add("column", "name")
		form.Add("value", "NewName")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockTxn.EXPECT().
			CheckActiveTransaction(gomock.Any(), "testuser").
			Return(true, nil)

		mockTxn.EXPECT().
			EditCell(gomock.Any(), "testuser", "testdb", "public", "users", 0, "name", "NewName").
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/transaction/edit-cell?database=testdb&schema=public&table=users", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleEditCell(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify cell edit response
		require.Contains(t, body, "success")
	})

	// E2E-S5-08: Transaction Mode Edit Buffer Display
	t.Run("E2E-S5-08: Transaction Mode Edit Buffer Display", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockTxn.EXPECT().
			GetTransactionEdits(gomock.Any(), "testuser").
			Return(map[int]domain.RowEdit{
				0: {
					RowIndex:   0,
					ColumnName: "name",
					OldValue:   "OldName",
					NewValue:   "NewName",
				},
				1: {
					RowIndex:   1,
					ColumnName: "email",
					OldValue:   "old@example.com",
					NewValue:   "new@example.com",
				},
			}, nil)

		mockTxn.EXPECT().
			GetTransactionDeletes(gomock.Any(), "testuser").
			Return([]int{}, nil)

		mockTxn.EXPECT().
			GetTransactionInserts(gomock.Any(), "testuser").
			Return([]domain.RowInsert{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/transaction/status", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleGetTransactionStatus(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify transaction status includes edit information
		require.NotEmpty(t, body)
	})

	// E2E-S5-09: Transaction Commit Button
	t.Run("E2E-S5-09: Transaction Commit Button", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockTxn.EXPECT().
			CheckActiveTransaction(gomock.Any(), "testuser").
			Return(true, nil)

		mockTxn.EXPECT().
			CommitTransaction(gomock.Any(), "testuser").
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/transaction/commit", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleCommitTransaction(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify commit success message
		require.Contains(t, body, "committed")
		require.Contains(t, body, "success")
	})

	// E2E-S5-10: Transaction Rollback Button
	t.Run("E2E-S5-10: Transaction Rollback Button", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockTxn.EXPECT().
			CheckActiveTransaction(gomock.Any(), "testuser").
			Return(true, nil)

		mockTxn.EXPECT().
			RollbackTransaction(gomock.Any(), "testuser").
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/transaction/rollback", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleRollbackTransaction(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify rollback success message
		require.Contains(t, body, "rolled back")
		require.Contains(t, body, "success")
	})

	// E2E-S5-11: Transaction Timer Countdown
	t.Run("E2E-S5-11: Transaction Timer Countdown", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockTxn.EXPECT().
			GetTransactionRemainingTime(gomock.Any(), "testuser").
			Return(int64(45), nil)

		req := httptest.NewRequest(http.MethodGet, "/transaction/status", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleGetTransactionStatus(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify transaction status includes remaining time information
		require.NotEmpty(t, body)
	})

	// E2E-S5-12: Transaction Row Delete Button
	t.Run("E2E-S5-12: Transaction Row Delete Button", func(t *testing.T) {
		form := url.Values{}
		form.Add("row_index", "2")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockTxn.EXPECT().
			CheckActiveTransaction(gomock.Any(), "testuser").
			Return(true, nil)

		mockTxn.EXPECT().
			DeleteRow(gomock.Any(), "testuser", "testdb", "public", "users", 2).
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/transaction/delete-row?database=testdb&schema=public&table=users", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleDeleteRow(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify row marked for deletion
		require.Contains(t, body, "deleted")
	})

	// E2E-S5-13: Transaction New Row Button
	t.Run("E2E-S5-13: Transaction New Row Button", func(t *testing.T) {
		form := url.Values{}
		form.Add("name", "NewUser")
		form.Add("email", "newuser@example.com")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockTxn.EXPECT().
			CheckActiveTransaction(gomock.Any(), "testuser").
			Return(true, nil)

		mockTxn.EXPECT().
			InsertRow(gomock.Any(), "testuser", "testdb", "public", "users", gomock.Any()).
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/transaction/insert-row?database=testdb&schema=public&table=users", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleInsertRow(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify new row added to buffer
		require.Contains(t, body, "inserted")
	})

	// Additional test: Transaction already active error
	t.Run("Transaction Already Active Error", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockRBAC.EXPECT().
			CheckUpdatePermission(gomock.Any(), "testuser", "testdb", "public", "users").
			Return(true, nil)

		mockTxn.EXPECT().
			CheckActiveTransaction(gomock.Any(), "testuser").
			Return(true, nil)

		req := httptest.NewRequest(http.MethodPost, "/transaction/start?database=testdb&schema=public&table=users", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleStartTransaction(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusConflict, rec.Code)
		body := rec.Body.String()

		// Verify error message
		require.Contains(t, body, "Transaction already active")
	})

	// Additional test: Transaction timeout enforcement
	t.Run("Transaction Timer Expiration", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockTxn.EXPECT().
			IsTransactionExpired(gomock.Any(), "testuser").
			Return(true, nil)

		req := httptest.NewRequest(http.MethodPost, "/transaction/commit", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleCommitTransaction(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusRequestTimeout, rec.Code)
		body := rec.Body.String()

		// Verify timeout message
		require.Contains(t, body, "Transaction expired")
	})

	// Additional test: Edit without active transaction
	t.Run("Edit Cell Without Active Transaction", func(t *testing.T) {
		form := url.Values{}
		form.Add("row_index", "0")
		form.Add("column", "name")
		form.Add("value", "NewName")

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockTxn.EXPECT().
			CheckActiveTransaction(gomock.Any(), "testuser").
			Return(false, nil)

		req := httptest.NewRequest(http.MethodPost, "/transaction/edit-cell?database=testdb&schema=public&table=users", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleEditCell(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusBadRequest, rec.Code)
		body := rec.Body.String()

		// Verify error message
		require.Contains(t, body, "No active transaction")
	})

	// Additional test: Read-only user cannot start transaction
	t.Run("Read-Only User Cannot Start Transaction", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "readonly_user",
			}, nil)

		mockRBAC.EXPECT().
			CheckUpdatePermission(gomock.Any(), "readonly_user", "testdb", "public", "users").
			Return(false, nil)

		req := httptest.NewRequest(http.MethodPost, "/transaction/start?database=testdb&schema=public&table=users", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleStartTransaction(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusForbidden, rec.Code)
		body := rec.Body.String()

		// Verify permission error
		require.Contains(t, body, "permission denied")
	})

	// Additional test: Unauthorized access
	t.Run("Unauthorized Access to Transaction", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "").
			Return(nil, domain.ValidationError{Field: "session", Message: "No session"})

		req := httptest.NewRequest(http.MethodPost, "/transaction/start?database=testdb&schema=public&table=users", nil)
		rec := httptest.NewRecorder()

		h.HandleStartTransaction(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	// Additional test: Get transaction status
	t.Run("Get Transaction Status", func(t *testing.T) {
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_123").
			Return(&domain.Session{
				ID:       "session_123",
				Username: "testuser",
			}, nil)

		mockTxn.EXPECT().
			GetActiveTransaction(gomock.Any(), "testuser").
			Return(&domain.TransactionState{
				ID:       "txn_123",
				Username: "testuser",
			}, nil)

		req := httptest.NewRequest(http.MethodGet, "/transaction/status", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_123",
		})
		rec := httptest.NewRecorder()

		h.HandleGetTransactionStatus(rec, req.WithContext(ctx))

		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		// Verify transaction status
		require.Contains(t, body, "active")
		require.Contains(t, body, "txn_123")
	})
}
