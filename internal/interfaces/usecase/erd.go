package usecase

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// ERDUseCase defines operations for Entity-Relationship Diagram generation and viewing
type ERDUseCase interface {
	// GenerateERD generates ERD data for a database schema
	GenerateERD(ctx context.Context, username, database, schema string) (interface{}, error)

	// GetERDTables returns all tables to display in an ERD
	GetERDTables(ctx context.Context, username, database, schema string) ([]domain.TableMetadata, error)

	// GetERDRelationships returns all relationships to display in an ERD
	GetERDRelationships(ctx context.Context, username, database, schema string) ([]domain.ForeignKeyMetadata, error)

	// GetTableBoxData returns box representation data for a table in ERD
	GetTableBoxData(ctx context.Context, username, database, schema, table string) (*domain.TableMetadata, error)

	// GetRelationshipLines returns relationship line data for ERD visualization
	GetRelationshipLines(ctx context.Context, username, database, schema string) ([]map[string]interface{}, error)

	// IsSchemaEmpty checks if a schema has no tables
	IsSchemaEmpty(ctx context.Context, username, database, schema string) (bool, error)

	// GetAvailableSchemas returns all schemas accessible by a user in a database
	GetAvailableSchemas(ctx context.Context, username, database string) ([]string, error)
}
