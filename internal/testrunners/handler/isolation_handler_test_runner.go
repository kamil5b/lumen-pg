package handler

import (
	"context"
	"io"
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

// IsolationHandlerConstructor is a function type that creates handlers for isolation testing
type IsolationHandlerConstructor struct {
	LoginHandler        handler.LoginHandler
	MainViewHandler     handler.MainViewHandler
	TransactionHandler  handler.TransactionHandler
	DataExplorerHandler handler.DataExplorerHandler
}

// IsolationHandlerConstructorFunc creates all necessary handlers for isolation testing
type IsolationHandlerConstructorFunc func(
	authUC usecase.AuthenticationUseCase,
	dataViewUC usecase.DataViewUseCase,
	txnUC usecase.TransactionUseCase,
	rbacUC usecase.RBACUseCase,
	setupUC usecase.SetupUseCase,
) IsolationHandlerConstructor

// IsolationHandlerRunner runs all isolation handler tests
// Maps to TEST_PLAN.md:
// - Story 6: Isolation [UC-S6-01~03, E2E-S6-01~03]
//
// NOTE: Cookie and session isolation (UC-S6-03) is primarily a MIDDLEWARE concern
// NOTE: Handlers assume middleware has already isolated sessions per request
// NOTE: Handlers focus on business logic: multi-user scenarios, permission isolation, transaction isolation
func IsolationHandlerRunner(t *testing.T, constructor IsolationHandlerConstructorFunc) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockAuth := mockUsecase.NewMockAuthenticationUseCase(ctrl)
	mockDataView := mockUsecase.NewMockDataViewUseCase(ctrl)
	mockTxn := mockUsecase.NewMockTransactionUseCase(ctrl)
	mockRBAC := mockUsecase.NewMockRBACUseCase(ctrl)
	mockSetup := mockUsecase.NewMockSetupUseCase(ctrl)

	handlers := constructor(mockAuth, mockDataView, mockTxn, mockRBAC, mockSetup)

	// E2E-S6-01: Simultaneous Users Different Permissions
	t.Run("E2E-S6-01: Simultaneous Users Different Permissions", func(t *testing.T) {
		// User 1 - Admin with full access
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_admin").
			Return(&domain.Session{
				ID:       "session_admin",
				Username: "admin_user",
			}, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "admin_user").
			Return(&domain.RoleMetadata{
				Name:                "admin_user",
				AccessibleDatabases: []string{"testdb1", "testdb2", "testdb3"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb1", Schema: "public", Name: "users", HasSelect: true, HasInsert: true, HasUpdate: true, HasDelete: true},
					{Database: "testdb1", Schema: "public", Name: "posts", HasSelect: true, HasInsert: true, HasUpdate: true, HasDelete: true},
					{Database: "testdb1", Schema: "public", Name: "comments", HasSelect: true, HasInsert: true, HasUpdate: true, HasDelete: true},
				},
			}, nil)

		// User 2 - Limited user with read-only access
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_readonly").
			Return(&domain.Session{
				ID:       "session_readonly",
				Username: "readonly_user",
			}, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "readonly_user").
			Return(&domain.RoleMetadata{
				Name:                "readonly_user",
				AccessibleDatabases: []string{"testdb1"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb1", Schema: "public", Name: "users", HasSelect: true, HasInsert: false, HasUpdate: false, HasDelete: false},
				},
			}, nil)

		// Request from admin user - load data explorer
		reqAdmin := httptest.NewRequest(http.MethodGet, "/api/data-explorer", nil)
		reqAdmin.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_admin",
		})
		recAdmin := httptest.NewRecorder()

		handlers.DataExplorerHandler.HandleLoadDataExplorer(recAdmin, reqAdmin.WithContext(ctx))

		require.Equal(t, http.StatusOK, recAdmin.Code)
		bodyAdmin := recAdmin.Body.String()

		// Verify admin sees all resources
		require.Contains(t, bodyAdmin, "testdb1")
		require.Contains(t, bodyAdmin, "testdb2")
		require.Contains(t, bodyAdmin, "testdb3")
		require.Contains(t, bodyAdmin, "users")
		require.Contains(t, bodyAdmin, "posts")
		require.Contains(t, bodyAdmin, "comments")

		// Request from readonly user - load data explorer
		reqReadonly := httptest.NewRequest(http.MethodGet, "/api/data-explorer", nil)
		reqReadonly.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "session_readonly",
		})
		recReadonly := httptest.NewRecorder()

		handlers.DataExplorerHandler.HandleLoadDataExplorer(recReadonly, reqReadonly.WithContext(ctx))

		require.Equal(t, http.StatusOK, recReadonly.Code)
		bodyReadonly := recReadonly.Body.String()

		// Verify readonly user sees only permitted resources
		require.Contains(t, bodyReadonly, "testdb1")
		require.NotContains(t, bodyReadonly, "testdb2")
		require.NotContains(t, bodyReadonly, "testdb3")
		require.Contains(t, bodyReadonly, "users")
		require.NotContains(t, bodyReadonly, "posts")
		require.NotContains(t, bodyReadonly, "comments")
	})

	// E2E-S6-02: Simultaneous Transactions
	t.Run("E2E-S6-02: Simultaneous Transactions", func(t *testing.T) {
		// User 1 starts transaction
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_user1").
			Return(&domain.Session{
				ID:       "session_user1",
				Username: "user1",
			}, nil)

		mockRBAC.EXPECT().
			CheckUpdatePermission(gomock.Any(), "user1", "testdb", "public", "users").
			Return(true, nil)

		mockTxn.EXPECT().
			CheckActiveTransaction(gomock.Any(), "user1").
			Return(false, nil).Times(1)

		mockTxn.EXPECT().
			StartTransaction(gomock.Any(), "user1", "testdb", "public", "users").
			Return(&domain.TransactionState{
				ID:       "txn_user1",
				Username: "user1",
			}, nil)

		mockTxn.EXPECT().
			EditCell(gomock.Any(), "user1", "testdb", "public", "users", 0, "name", "User1Edit").
			Return(nil)

		mockTxn.EXPECT().
			GetTransactionEdits(gomock.Any(), "user1").
			Return(map[int]domain.RowEdit{
				0: {RowIndex: 0, ColumnName: "name", OldValue: "Original", NewValue: "User1Edit"},
			}, nil)

		// User 2 starts transaction
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_user2").
			Return(&domain.Session{
				ID:       "session_user2",
				Username: "user2",
			}, nil)

		mockRBAC.EXPECT().
			CheckUpdatePermission(gomock.Any(), "user2", "testdb", "public", "users").
			Return(true, nil)

		mockTxn.EXPECT().
			CheckActiveTransaction(gomock.Any(), "user2").
			Return(false, nil).Times(1)

		mockTxn.EXPECT().
			StartTransaction(gomock.Any(), "user2", "testdb", "public", "users").
			Return(&domain.TransactionState{
				ID:       "txn_user2",
				Username: "user2",
			}, nil)

		mockTxn.EXPECT().
			EditCell(gomock.Any(), "user2", "testdb", "public", "users", 0, "email", "user2@example.com").
			Return(nil)

		mockTxn.EXPECT().
			GetTransactionEdits(gomock.Any(), "user2").
			Return(map[int]domain.RowEdit{
				0: {RowIndex: 0, ColumnName: "email", OldValue: "old@example.com", NewValue: "user2@example.com"},
			}, nil)

		mockTxn.EXPECT().
			GetTransactionDeletes(gomock.Any(), "user1").
			Return([]int{}, nil)

		mockTxn.EXPECT().
			GetTransactionInserts(gomock.Any(), "user1").
			Return([]domain.RowInsert{}, nil)

		mockTxn.EXPECT().
			GetTransactionDeletes(gomock.Any(), "user2").
			Return([]int{}, nil)

		mockTxn.EXPECT().
			GetTransactionInserts(gomock.Any(), "user2").
			Return([]domain.RowInsert{}, nil)

		// User 1 starts transaction
		reqUser1Start := httptest.NewRequest(http.MethodPost, "/transaction/start?database=testdb&schema=public&table=users", nil)
		reqUser1Start.AddCookie(&http.Cookie{Name: "session_id", Value: "session_user1"})
		recUser1Start := httptest.NewRecorder()
		handlers.TransactionHandler.HandleStartTransaction(recUser1Start, reqUser1Start.WithContext(ctx))
		require.Equal(t, http.StatusOK, recUser1Start.Code)

		// User 2 starts transaction
		reqUser2Start := httptest.NewRequest(http.MethodPost, "/transaction/start?database=testdb&schema=public&table=users", nil)
		reqUser2Start.AddCookie(&http.Cookie{Name: "session_id", Value: "session_user2"})
		recUser2Start := httptest.NewRecorder()
		handlers.TransactionHandler.HandleStartTransaction(recUser2Start, reqUser2Start.WithContext(ctx))
		require.Equal(t, http.StatusOK, recUser2Start.Code)

		// User 1 makes edit
		mockTxn.EXPECT().CheckActiveTransaction(gomock.Any(), "user1").Return(true, nil)
		formUser1Edit := url.Values{}
		formUser1Edit.Add("row_index", "0")
		formUser1Edit.Add("column", "name")
		formUser1Edit.Add("value", "User1Edit")
		reqUser1Edit := httptest.NewRequest(http.MethodPost, "/transaction/edit-cell?database=testdb&schema=public&table=users", strings.NewReader(formUser1Edit.Encode()))
		reqUser1Edit.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		reqUser1Edit.AddCookie(&http.Cookie{Name: "session_id", Value: "session_user1"})
		recUser1Edit := httptest.NewRecorder()
		handlers.TransactionHandler.HandleEditCell(recUser1Edit, reqUser1Edit.WithContext(ctx))
		require.Equal(t, http.StatusOK, recUser1Edit.Code)

		// User 2 makes edit
		mockTxn.EXPECT().CheckActiveTransaction(gomock.Any(), "user2").Return(true, nil)
		formUser2Edit := url.Values{}
		formUser2Edit.Add("row_index", "0")
		formUser2Edit.Add("column", "email")
		formUser2Edit.Add("value", "user2@example.com")
		reqUser2Edit := httptest.NewRequest(http.MethodPost, "/transaction/edit-cell?database=testdb&schema=public&table=users", strings.NewReader(formUser2Edit.Encode()))
		reqUser2Edit.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		reqUser2Edit.AddCookie(&http.Cookie{Name: "session_id", Value: "session_user2"})
		recUser2Edit := httptest.NewRecorder()
		handlers.TransactionHandler.HandleEditCell(recUser2Edit, reqUser2Edit.WithContext(ctx))
		require.Equal(t, http.StatusOK, recUser2Edit.Code)

		// Verify User 1's edits are isolated
		reqUser1Buffer := httptest.NewRequest(http.MethodGet, "/transaction/buffer", nil)
		reqUser1Buffer.AddCookie(&http.Cookie{Name: "session_id", Value: "session_user1"})
		recUser1Buffer := httptest.NewRecorder()
		handlers.TransactionHandler.ServeHTTP(recUser1Buffer, reqUser1Buffer.WithContext(ctx))
		require.Equal(t, http.StatusOK, recUser1Buffer.Code)
		bodyUser1 := recUser1Buffer.Body.String()
		require.Contains(t, bodyUser1, "User1Edit")
		require.NotContains(t, bodyUser1, "user2@example.com")

		// Verify User 2's edits are isolated
		reqUser2Buffer := httptest.NewRequest(http.MethodGet, "/transaction/buffer", nil)
		reqUser2Buffer.AddCookie(&http.Cookie{Name: "session_id", Value: "session_user2"})
		recUser2Buffer := httptest.NewRecorder()
		handlers.TransactionHandler.ServeHTTP(recUser2Buffer, reqUser2Buffer.WithContext(ctx))
		require.Equal(t, http.StatusOK, recUser2Buffer.Code)
		bodyUser2 := recUser2Buffer.Body.String()
		require.Contains(t, bodyUser2, "user2@example.com")
		require.NotContains(t, bodyUser2, "User1Edit")
	})

	// E2E-S6-03: One User Cannot See Another's Session
	t.Run("E2E-S6-03: One User Cannot See Another's Session", func(t *testing.T) {
		// User A logs in
		formUserA := url.Values{}
		formUserA.Add("username", "userA")
		formUserA.Add("password", "passwordA")

		mockAuth.EXPECT().
			ValidateLoginForm(gomock.Any(), gomock.Any()).
			Return([]domain.ValidationError{}, nil)

		mockAuth.EXPECT().
			ProbeConnection(gomock.Any(), "userA", "passwordA").
			Return(true, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "userA").
			Return(&domain.RoleMetadata{
				Name:                "userA",
				AccessibleDatabases: []string{"testdb"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb", Schema: "public", Name: "users", HasSelect: true},
				},
			}, nil)

		mockAuth.EXPECT().
			GetFirstAccessibleDatabase(gomock.Any(), "userA").
			Return("testdb", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleSchema(gomock.Any(), "userA", "testdb").
			Return("public", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleTable(gomock.Any(), "userA", "testdb", "public").
			Return("users", nil)

		mockAuth.EXPECT().
			CreateSession(gomock.Any(), "userA", "passwordA", "testdb", "public", "users").
			Return(&domain.Session{
				ID:       "session_userA",
				Username: "userA",
			}, nil)

		reqLoginA := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formUserA.Encode()))
		reqLoginA.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		recLoginA := httptest.NewRecorder()

		handlers.LoginHandler.HandleLogin(recLoginA, reqLoginA.WithContext(ctx))

		require.Equal(t, http.StatusFound, recLoginA.Code)

		// User B logs in
		formUserB := url.Values{}
		formUserB.Add("username", "userB")
		formUserB.Add("password", "passwordB")

		mockAuth.EXPECT().
			ValidateLoginForm(gomock.Any(), gomock.Any()).
			Return([]domain.ValidationError{}, nil)

		mockAuth.EXPECT().
			ProbeConnection(gomock.Any(), "userB", "passwordB").
			Return(true, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "userB").
			Return(&domain.RoleMetadata{
				Name:                "userB",
				AccessibleDatabases: []string{"testdb"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb", Schema: "public", Name: "users", HasSelect: true},
				},
			}, nil)

		mockAuth.EXPECT().
			GetFirstAccessibleDatabase(gomock.Any(), "userB").
			Return("testdb", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleSchema(gomock.Any(), "userB", "testdb").
			Return("public", nil)

		mockAuth.EXPECT().
			GetFirstAccessibleTable(gomock.Any(), "userB", "testdb", "public").
			Return("users", nil)

		mockAuth.EXPECT().
			CreateSession(gomock.Any(), "userB", "passwordB", "testdb", "public", "users").
			Return(&domain.Session{
				ID:       "session_userB",
				Username: "userB",
			}, nil)

		reqLoginB := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formUserB.Encode()))
		reqLoginB.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		recLoginB := httptest.NewRecorder()

		handlers.LoginHandler.HandleLogin(recLoginB, reqLoginB.WithContext(ctx))

		require.Equal(t, http.StatusFound, recLoginB.Code)

		// Verify session cookies are different
		cookieA := recLoginA.Result().Cookies()[0]
		cookieB := recLoginB.Result().Cookies()[0]
		require.NotEqual(t, cookieA.Value, cookieB.Value)

		// User A loads main view
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_userA").
			Return(&domain.Session{
				ID:       "session_userA",
				Username: "userA",
			}, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "userA").
			Return(&domain.RoleMetadata{
				Name:                "userA",
				AccessibleDatabases: []string{"testdb"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb", Schema: "public", Name: "users", HasSelect: true},
				},
			}, nil)

		mockDataView.EXPECT().
			LoadTableData(gomock.Any(), "userA", gomock.Any()).
			Return(&domain.QueryResult{
				Columns:  []string{"id", "name"},
				Rows:     []map[string]interface{}{{"id": 1, "name": "UserA Data"}},
				RowCount: 1,
			}, nil)

		reqMainA := httptest.NewRequest(http.MethodGet, "/main", nil)
		reqMainA.AddCookie(&http.Cookie{Name: "session_id", Value: "session_userA"})
		recMainA := httptest.NewRecorder()

		handlers.MainViewHandler.HandleMainViewPage(recMainA, reqMainA.WithContext(ctx))

		require.Equal(t, http.StatusOK, recMainA.Code)
		require.Contains(t, recMainA.Body.String(), "userA")
		require.NotContains(t, recMainA.Body.String(), "userB")

		// User B loads main view
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_userB").
			Return(&domain.Session{
				ID:       "session_userB",
				Username: "userB",
			}, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "userB").
			Return(&domain.RoleMetadata{
				Name:                "userB",
				AccessibleDatabases: []string{"testdb"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb", Schema: "public", Name: "users", HasSelect: true},
				},
			}, nil)

		mockDataView.EXPECT().
			LoadTableData(gomock.Any(), "userB", gomock.Any()).
			Return(&domain.QueryResult{
				Columns:  []string{"id", "name"},
				Rows:     []map[string]interface{}{{"id": 2, "name": "UserB Data"}},
				RowCount: 1,
			}, nil)

		reqMainB := httptest.NewRequest(http.MethodGet, "/main", nil)
		reqMainB.AddCookie(&http.Cookie{Name: "session_id", Value: "session_userB"})
		recMainB := httptest.NewRecorder()

		handlers.MainViewHandler.HandleMainViewPage(recMainB, reqMainB.WithContext(ctx))

		require.Equal(t, http.StatusOK, recMainB.Code)
		require.Contains(t, recMainB.Body.String(), "userB")
		require.NotContains(t, recMainB.Body.String(), "userA")
	})

	// Additional test: Session cookie isolation
	t.Run("Session Cookie Isolation", func(t *testing.T) {
		// Two different sessions with different users
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_user1").
			Return(&domain.Session{
				ID:       "session_user1",
				Username: "user1",
			}, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "user1").
			Return(&domain.RoleMetadata{
				Name:                "user1",
				AccessibleDatabases: []string{"testdb"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb", Schema: "public", Name: "users", HasSelect: true},
				},
			}, nil)

		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_user2").
			Return(&domain.Session{
				ID:       "session_user2",
				Username: "user2",
			}, nil)

		mockAuth.EXPECT().
			GetUserAccessibleResources(gomock.Any(), "user2").
			Return(&domain.RoleMetadata{
				Name:                "user2",
				AccessibleDatabases: []string{"testdb"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables: []domain.AccessibleTable{
					{Database: "testdb", Schema: "public", Name: "posts", HasSelect: true},
				},
			}, nil)

		// User 1 request
		req1 := httptest.NewRequest(http.MethodGet, "/api/data-explorer", nil)
		req1.AddCookie(&http.Cookie{Name: "session_id", Value: "session_user1"})
		rec1 := httptest.NewRecorder()

		handlers.DataExplorerHandler.HandleLoadDataExplorer(rec1, req1.WithContext(ctx))
		require.Equal(t, http.StatusOK, rec1.Code)
		require.Contains(t, rec1.Body.String(), "users")

		// User 2 request
		req2 := httptest.NewRequest(http.MethodGet, "/api/data-explorer", nil)
		req2.AddCookie(&http.Cookie{Name: "session_id", Value: "session_user2"})
		rec2 := httptest.NewRecorder()

		handlers.DataExplorerHandler.HandleLoadDataExplorer(rec2, req2.WithContext(ctx))
		require.Equal(t, http.StatusOK, rec2.Code)
		require.Contains(t, rec2.Body.String(), "posts")
		require.NotContains(t, rec2.Body.String(), "users")
	})

	// Additional test: Concurrent data access isolation
	t.Run("Concurrent Data Access Isolation", func(t *testing.T) {
		// User 1 loads data
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_user1").
			Return(&domain.Session{ID: "session_user1", Username: "user1"}, nil)

		mockRBAC.EXPECT().
			CheckTableAccess(gomock.Any(), "user1", "testdb", "public", "users").
			Return(true, nil)

		mockDataView.EXPECT().
			LoadTableData(gomock.Any(), "user1", gomock.Any()).
			Return(&domain.QueryResult{
				Columns:  []string{"id", "name"},
				Rows:     []map[string]interface{}{{"id": 1, "name": "User1Data"}},
				RowCount: 1,
			}, nil)

		// User 2 loads data
		mockAuth.EXPECT().
			ValidateSession(gomock.Any(), "session_user2").
			Return(&domain.Session{ID: "session_user2", Username: "user2"}, nil)

		mockRBAC.EXPECT().
			CheckTableAccess(gomock.Any(), "user2", "testdb", "public", "users").
			Return(true, nil)

		mockDataView.EXPECT().
			LoadTableData(gomock.Any(), "user2", gomock.Any()).
			Return(&domain.QueryResult{
				Columns:  []string{"id", "name"},
				Rows:     []map[string]interface{}{{"id": 2, "name": "User2Data"}},
				RowCount: 1,
			}, nil)

		// Both users load data concurrently
		form1 := url.Values{}
		form1.Add("database", "testdb")
		form1.Add("schema", "public")
		form1.Add("table", "users")
		req1 := httptest.NewRequest(http.MethodPost, "/main/load-data", strings.NewReader(form1.Encode()))
		req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req1.AddCookie(&http.Cookie{Name: "session_id", Value: "session_user1"})
		rec1 := httptest.NewRecorder()

		form2 := url.Values{}
		form2.Add("database", "testdb")
		form2.Add("schema", "public")
		form2.Add("table", "users")
		req2 := httptest.NewRequest(http.MethodPost, "/main/load-data", io.NopCloser(strings.NewReader(form2.Encode())))
		req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req2.AddCookie(&http.Cookie{Name: "session_id", Value: "session_user2"})
		rec2 := httptest.NewRecorder()

		handlers.MainViewHandler.HandleLoadTableData(rec1, req1.WithContext(ctx))
		handlers.MainViewHandler.HandleLoadTableData(rec2, req2.WithContext(ctx))

		// Verify each user sees only their data
		require.Equal(t, http.StatusOK, rec1.Code)
		require.Contains(t, rec1.Body.String(), "User1Data")
		require.NotContains(t, rec1.Body.String(), "User2Data")

		require.Equal(t, http.StatusOK, rec2.Code)
		require.Contains(t, rec2.Body.String(), "User2Data")
		require.NotContains(t, rec2.Body.String(), "User1Data")
	})
}
