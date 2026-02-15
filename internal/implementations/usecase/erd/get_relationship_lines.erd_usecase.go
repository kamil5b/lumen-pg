package erd

import (
	"context"
	"errors"
)

func (u *ERDUseCaseImplementation) GetRelationshipLines(ctx context.Context, username, database, schema string) ([]map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}
