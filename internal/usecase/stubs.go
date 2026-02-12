package usecase

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/repository"
)

// StubAuthUsecase is a stub implementation of AuthUsecase.
type StubAuthUsecase struct {
	ConnRepo     repository.ConnectionRepository
	MetadataRepo repository.MetadataRepository
}

func NewStubAuthUsecase(connRepo repository.ConnectionRepository, metadataRepo repository.MetadataRepository) *StubAuthUsecase {
	return &StubAuthUsecase{ConnRepo: connRepo, MetadataRepo: metadataRepo}
}

func (s *StubAuthUsecase) Login(ctx context.Context, username, password string) (*domain.Session, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubAuthUsecase) ValidateSession(session *domain.Session) error {
	return domain.ErrNotImplemented
}

func (s *StubAuthUsecase) Logout(ctx context.Context, username string) error {
	return domain.ErrNotImplemented
}

func (s *StubAuthUsecase) ReAuthenticate(ctx context.Context, username, encryptedPassword string) error {
	return domain.ErrNotImplemented
}

func (s *StubAuthUsecase) GetAccessibleResources(ctx context.Context, username string) (*domain.Role, error) {
	return nil, domain.ErrNotImplemented
}

// StubMetadataUsecase is a stub implementation of MetadataUsecase.
type StubMetadataUsecase struct {
	MetadataRepo repository.MetadataRepository
}

func NewStubMetadataUsecase(metadataRepo repository.MetadataRepository) *StubMetadataUsecase {
	return &StubMetadataUsecase{MetadataRepo: metadataRepo}
}

func (s *StubMetadataUsecase) InitializeMetadata(ctx context.Context, connStr string) error {
	return domain.ErrNotImplemented
}

func (s *StubMetadataUsecase) GetMetadata() *domain.Metadata {
	return nil
}

func (s *StubMetadataUsecase) RefreshMetadata(ctx context.Context) error {
	return domain.ErrNotImplemented
}

func (s *StubMetadataUsecase) GetRBACMapping() map[string]*domain.Role {
	return nil
}

func (s *StubMetadataUsecase) GetERDData(ctx context.Context, database, schema string) (*domain.ERDData, error) {
	return nil, domain.ErrNotImplemented
}

// StubDataExplorerUsecase is a stub implementation of DataExplorerUsecase.
type StubDataExplorerUsecase struct {
	QueryRepo repository.QueryRepository
}

func NewStubDataExplorerUsecase(queryRepo repository.QueryRepository) *StubDataExplorerUsecase {
	return &StubDataExplorerUsecase{QueryRepo: queryRepo}
}

func (s *StubDataExplorerUsecase) LoadTableData(ctx context.Context, database, schema, table string, cursor *domain.Cursor, orderBy string, orderDir string, whereClause string) (*domain.CursorPage, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubDataExplorerUsecase) GetTotalRowCount(ctx context.Context, database, schema, table string, whereClause string) (int64, error) {
	return 0, domain.ErrNotImplemented
}

func (s *StubDataExplorerUsecase) GetReferencingTables(ctx context.Context, database, schema, table, pkColumn string, pkValue interface{}) ([]domain.ReferencingTable, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubDataExplorerUsecase) NavigateToForeignKey(ctx context.Context, database, schema, table, column string, value interface{}) (*domain.CursorPage, error) {
	return nil, domain.ErrNotImplemented
}

// StubQueryUsecase is a stub implementation of QueryUsecase.
type StubQueryUsecase struct {
	QueryRepo repository.QueryRepository
}

func NewStubQueryUsecase(queryRepo repository.QueryRepository) *StubQueryUsecase {
	return &StubQueryUsecase{QueryRepo: queryRepo}
}

func (s *StubQueryUsecase) ExecuteQuery(ctx context.Context, database, query string) (*domain.QueryResult, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubQueryUsecase) ExecuteQueries(ctx context.Context, database, sql string) ([]domain.QueryResult, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubQueryUsecase) ExecuteQueryWithPagination(ctx context.Context, database, query string, offset, limit int) (*domain.QueryResult, error) {
	return nil, domain.ErrNotImplemented
}

// StubTransactionUsecase is a stub implementation of TransactionUsecase.
type StubTransactionUsecase struct {
	QueryRepo repository.QueryRepository
}

func NewStubTransactionUsecase(queryRepo repository.QueryRepository) *StubTransactionUsecase {
	return &StubTransactionUsecase{QueryRepo: queryRepo}
}

func (s *StubTransactionUsecase) StartTransaction(ctx context.Context, sessionID string) error {
	return domain.ErrNotImplemented
}

func (s *StubTransactionUsecase) AddOperation(ctx context.Context, sessionID string, op domain.BufferedOperation) error {
	return domain.ErrNotImplemented
}

func (s *StubTransactionUsecase) Commit(ctx context.Context, sessionID string) error {
	return domain.ErrNotImplemented
}

func (s *StubTransactionUsecase) Rollback(ctx context.Context, sessionID string) error {
	return domain.ErrNotImplemented
}

func (s *StubTransactionUsecase) GetTransaction(sessionID string) (*domain.Transaction, error) {
	return nil, domain.ErrNotImplemented
}
