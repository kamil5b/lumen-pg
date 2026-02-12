package repository

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// StubConnectionRepository is a stub implementation of ConnectionRepository.
type StubConnectionRepository struct{}

func NewStubConnectionRepository() *StubConnectionRepository {
	return &StubConnectionRepository{}
}

func (s *StubConnectionRepository) ValidateConnectionString(connStr string) error {
	return domain.ErrNotImplemented
}

func (s *StubConnectionRepository) ParseConnectionString(connStr string) (*domain.ConnectionConfig, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubConnectionRepository) TestConnection(ctx context.Context, connStr string) error {
	return domain.ErrNotImplemented
}

func (s *StubConnectionRepository) ProbeConnection(ctx context.Context, username, password, host string, port int, sslMode string, accessibleDBs []string) (*domain.Session, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubConnectionRepository) Connect(ctx context.Context, config *domain.ConnectionConfig) (interface{}, error) {
	return nil, domain.ErrNotImplemented
}

// StubMetadataRepository is a stub implementation of MetadataRepository.
type StubMetadataRepository struct{}

func NewStubMetadataRepository() *StubMetadataRepository {
	return &StubMetadataRepository{}
}

func (s *StubMetadataRepository) LoadDatabases(ctx context.Context) ([]domain.Database, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubMetadataRepository) LoadSchemas(ctx context.Context, database string) ([]domain.Schema, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubMetadataRepository) LoadTables(ctx context.Context, database, schema string) ([]domain.Table, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubMetadataRepository) LoadColumns(ctx context.Context, database, schema, table string) ([]domain.Column, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubMetadataRepository) LoadForeignKeys(ctx context.Context, database, schema, table string) ([]domain.ForeignKey, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubMetadataRepository) LoadRoles(ctx context.Context) ([]domain.Role, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubMetadataRepository) LoadAllMetadata(ctx context.Context) (*domain.Metadata, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubMetadataRepository) GetAccessibleResources(ctx context.Context, roleName string) (*domain.Role, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubMetadataRepository) GenerateERDData(ctx context.Context, database, schema string) (*domain.ERDData, error) {
	return nil, domain.ErrNotImplemented
}

// StubQueryRepository is a stub implementation of QueryRepository.
type StubQueryRepository struct{}

func NewStubQueryRepository() *StubQueryRepository {
	return &StubQueryRepository{}
}

func (s *StubQueryRepository) ExecuteQuery(ctx context.Context, database, query string, params ...interface{}) (*domain.QueryResult, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubQueryRepository) ExecuteQueries(ctx context.Context, database string, queries []string) ([]domain.QueryResult, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubQueryRepository) LoadTableData(ctx context.Context, database, schema, table string, cursor *domain.Cursor, orderBy string, orderDir string, whereClause string) (*domain.CursorPage, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubQueryRepository) GetTotalRowCount(ctx context.Context, database, schema, table string, whereClause string) (int64, error) {
	return 0, domain.ErrNotImplemented
}

func (s *StubQueryRepository) GetReferencingTables(ctx context.Context, database, schema, table, pkColumn string, pkValue interface{}) ([]domain.ReferencingTable, error) {
	return nil, domain.ErrNotImplemented
}

func (s *StubQueryRepository) ExecuteTransaction(ctx context.Context, database string, operations []domain.BufferedOperation) error {
	return domain.ErrNotImplemented
}
