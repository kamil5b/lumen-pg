package dataview

import (
	"context"
	"errors"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *DataViewUseCaseImplementation) GetChildTableReferences(ctx context.Context, username, database, schema, table string, pkValues map[string]interface{}) ([]domain.ChildTableReference, error) {
	return nil, errors.New("not implemented")
}
