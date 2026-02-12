package usecase_test

import (
	"context"
	"testing"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/mocks"
	"github.com/kamil5b/lumen-pg/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// Phase 3: Use Case Tests using gomock

// === AuthUsecase Tests ===

// UC-S2-03: Login Connection Probe
func TestAuthUsecase_Login_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnRepo := mocks.NewMockConnectionRepository(ctrl)
	mockMetadataRepo := mocks.NewMockMetadataRepository(ctrl)

	uc := usecase.NewStubAuthUsecase(mockConnRepo, mockMetadataRepo)
	session, err := uc.Login(context.Background(), "admin", "password")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Nil(t, session)
}

// UC-S2-04: Login Connection Probe Failure
func TestAuthUsecase_Login_NoAccessibleResources(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnRepo := mocks.NewMockConnectionRepository(ctrl)
	mockMetadataRepo := mocks.NewMockMetadataRepository(ctrl)

	uc := usecase.NewStubAuthUsecase(mockConnRepo, mockMetadataRepo)
	session, err := uc.Login(context.Background(), "limited_user", "password")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Nil(t, session)
}

// UC-S2-05: Login Success After Probe
func TestAuthUsecase_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnRepo := mocks.NewMockConnectionRepository(ctrl)
	mockMetadataRepo := mocks.NewMockMetadataRepository(ctrl)

	uc := usecase.NewStubAuthUsecase(mockConnRepo, mockMetadataRepo)
	session, err := uc.Login(context.Background(), "admin", "password")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Nil(t, session)
}

// UC-S2-08: Session Validation - Valid Session
func TestAuthUsecase_ValidateSession_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnRepo := mocks.NewMockConnectionRepository(ctrl)
	mockMetadataRepo := mocks.NewMockMetadataRepository(ctrl)

	uc := usecase.NewStubAuthUsecase(mockConnRepo, mockMetadataRepo)
	session := &domain.Session{Username: "admin"}
	err := uc.ValidateSession(session)
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
}

// UC-S2-10: Session Re-authentication
func TestAuthUsecase_ReAuthenticate_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnRepo := mocks.NewMockConnectionRepository(ctrl)
	mockMetadataRepo := mocks.NewMockMetadataRepository(ctrl)

	uc := usecase.NewStubAuthUsecase(mockConnRepo, mockMetadataRepo)
	err := uc.ReAuthenticate(context.Background(), "admin", "encrypted_password")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
}

// UC-S2-12: Logout
func TestAuthUsecase_Logout_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnRepo := mocks.NewMockConnectionRepository(ctrl)
	mockMetadataRepo := mocks.NewMockMetadataRepository(ctrl)

	uc := usecase.NewStubAuthUsecase(mockConnRepo, mockMetadataRepo)
	err := uc.Logout(context.Background(), "admin")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
}

// UC-S2-11: Get Accessible Resources
func TestAuthUsecase_GetAccessibleResources_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnRepo := mocks.NewMockConnectionRepository(ctrl)
	mockMetadataRepo := mocks.NewMockMetadataRepository(ctrl)

	uc := usecase.NewStubAuthUsecase(mockConnRepo, mockMetadataRepo)
	role, err := uc.GetAccessibleResources(context.Background(), "admin")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Nil(t, role)
}

// === MetadataUsecase Tests ===

// UC-S1-05: Metadata Initialization
func TestMetadataUsecase_InitializeMetadata_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetadataRepo := mocks.NewMockMetadataRepository(ctrl)

	uc := usecase.NewStubMetadataUsecase(mockMetadataRepo)
	err := uc.InitializeMetadata(context.Background(), "postgres://user:pass@host:5432/db")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
}

// UC-S1-06: In-Memory Metadata Storage
func TestMetadataUsecase_GetMetadata_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetadataRepo := mocks.NewMockMetadataRepository(ctrl)

	uc := usecase.NewStubMetadataUsecase(mockMetadataRepo)
	metadata := uc.GetMetadata()
	assert.Nil(t, metadata)
}

// UC-S1-07: RBAC Initialization
func TestMetadataUsecase_GetRBACMapping_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetadataRepo := mocks.NewMockMetadataRepository(ctrl)

	uc := usecase.NewStubMetadataUsecase(mockMetadataRepo)
	mapping := uc.GetRBACMapping()
	assert.Nil(t, mapping)
}

// UC-S2-15: Metadata Refresh
func TestMetadataUsecase_RefreshMetadata_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetadataRepo := mocks.NewMockMetadataRepository(ctrl)

	uc := usecase.NewStubMetadataUsecase(mockMetadataRepo)
	err := uc.RefreshMetadata(context.Background())
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
}

// UC-S3-01: ERD Data Generation
func TestMetadataUsecase_GetERDData_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMetadataRepo := mocks.NewMockMetadataRepository(ctrl)

	uc := usecase.NewStubMetadataUsecase(mockMetadataRepo)
	erd, err := uc.GetERDData(context.Background(), "testdb", "public")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Nil(t, erd)
}

// === DataExplorerUsecase Tests ===

// UC-S5-01: Table Data Loading
func TestDataExplorerUsecase_LoadTableData_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryRepo := mocks.NewMockQueryRepository(ctrl)

	uc := usecase.NewStubDataExplorerUsecase(mockQueryRepo)
	cursor := domain.NewCursor()
	page, err := uc.LoadTableData(context.Background(), "testdb", "public", "users", cursor, "id", "ASC", "")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Nil(t, page)
}

// UC-S5-07: Total Row Count
func TestDataExplorerUsecase_GetTotalRowCount_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryRepo := mocks.NewMockQueryRepository(ctrl)

	uc := usecase.NewStubDataExplorerUsecase(mockQueryRepo)
	count, err := uc.GetTotalRowCount(context.Background(), "testdb", "public", "users", "")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Equal(t, int64(0), count)
}

// UC-S5-18: Primary Key Navigation
func TestDataExplorerUsecase_GetReferencingTables_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryRepo := mocks.NewMockQueryRepository(ctrl)

	uc := usecase.NewStubDataExplorerUsecase(mockQueryRepo)
	tables, err := uc.GetReferencingTables(context.Background(), "testdb", "public", "users", "id", 1)
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Nil(t, tables)
}

// UC-S5-17: Foreign Key Navigation
func TestDataExplorerUsecase_NavigateToForeignKey_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryRepo := mocks.NewMockQueryRepository(ctrl)

	uc := usecase.NewStubDataExplorerUsecase(mockQueryRepo)
	page, err := uc.NavigateToForeignKey(context.Background(), "testdb", "public", "posts", "user_id", 1)
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Nil(t, page)
}

// === QueryUsecase Tests ===

// UC-S4-01: Single Query Execution
func TestQueryUsecase_ExecuteQuery_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryRepo := mocks.NewMockQueryRepository(ctrl)

	uc := usecase.NewStubQueryUsecase(mockQueryRepo)
	result, err := uc.ExecuteQuery(context.Background(), "testdb", "SELECT * FROM users")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Nil(t, result)
}

// UC-S4-02: Multiple Query Execution
func TestQueryUsecase_ExecuteQueries_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryRepo := mocks.NewMockQueryRepository(ctrl)

	uc := usecase.NewStubQueryUsecase(mockQueryRepo)
	results, err := uc.ExecuteQueries(context.Background(), "testdb", "SELECT 1; SELECT 2")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Nil(t, results)
}

// UC-S4-03: Query Result Offset Pagination
func TestQueryUsecase_ExecuteQueryWithPagination_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryRepo := mocks.NewMockQueryRepository(ctrl)

	uc := usecase.NewStubQueryUsecase(mockQueryRepo)
	result, err := uc.ExecuteQueryWithPagination(context.Background(), "testdb", "SELECT * FROM users", 0, 1000)
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Nil(t, result)
}

// === TransactionUsecase Tests ===

// UC-S5-09: Transaction Start
func TestTransactionUsecase_StartTransaction_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryRepo := mocks.NewMockQueryRepository(ctrl)

	uc := usecase.NewStubTransactionUsecase(mockQueryRepo)
	err := uc.StartTransaction(context.Background(), "session-123")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
}

// UC-S5-11: Add Operation
func TestTransactionUsecase_AddOperation_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryRepo := mocks.NewMockQueryRepository(ctrl)

	uc := usecase.NewStubTransactionUsecase(mockQueryRepo)
	op := domain.BufferedOperation{Type: domain.OpUpdate, Table: "users"}
	err := uc.AddOperation(context.Background(), "session-123", op)
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
}

// UC-S5-12: Transaction Commit
func TestTransactionUsecase_Commit_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryRepo := mocks.NewMockQueryRepository(ctrl)

	uc := usecase.NewStubTransactionUsecase(mockQueryRepo)
	err := uc.Commit(context.Background(), "session-123")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
}

// UC-S5-13: Transaction Rollback
func TestTransactionUsecase_Rollback_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryRepo := mocks.NewMockQueryRepository(ctrl)

	uc := usecase.NewStubTransactionUsecase(mockQueryRepo)
	err := uc.Rollback(context.Background(), "session-123")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
}

// UC-S5-09: Get Transaction
func TestTransactionUsecase_GetTransaction_Stub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryRepo := mocks.NewMockQueryRepository(ctrl)

	uc := usecase.NewStubTransactionUsecase(mockQueryRepo)
	txn, err := uc.GetTransaction("session-123")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Nil(t, txn)
}
