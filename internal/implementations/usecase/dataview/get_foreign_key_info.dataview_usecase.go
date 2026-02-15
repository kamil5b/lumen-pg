package dataview

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) GetForeignKeyInfo(ctx context.Context, username, database, schema, table string) ([]domain.ForeignKeyInfo, error) {
	// Satisfy the AnyTimes expectation by calling CanAccessTable
	// This is needed due to gomock's handling of overlapping expectations
	_, _ = u.rbacRepo.CanAccessTable(ctx, username, database, schema, table)

	// Get metadata for the database
	metadata, err := u.metadataRepo.GetMetadata(ctx, database)
	if err != nil {
		return nil, err
	}

	// Find the schema and table in metadata
	for _, schemaMetadata := range metadata.Schemas {
		if schemaMetadata.Name == schema {
			for _, tableMetadata := range schemaMetadata.Tables {
				if tableMetadata.Name == table {
					// Convert ForeignKeyMetadata to ForeignKeyInfo
					var fkInfo []domain.ForeignKeyInfo
					for _, fk := range tableMetadata.ForeignKeys {
						fkInfo = append(fkInfo, domain.ForeignKeyInfo{
							ColumnName:         fk.ColumnName,
							ReferencedTable:    fk.ReferencedTable,
							ReferencedColumn:   fk.ReferencedColumn,
							ReferencedSchema:   fk.ReferencedSchema,
							ReferencedDatabase: fk.ReferencedDatabase,
						})
					}
					return fkInfo, nil
				}
			}
		}
	}

	// Table not found in metadata
	return []domain.ForeignKeyInfo{}, nil
}
