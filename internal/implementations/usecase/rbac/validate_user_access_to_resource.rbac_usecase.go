package rbac

import (
	"context"
	"fmt"
)

func (u *RBACUseCaseImplementation) ValidateUserAccessToResource(ctx context.Context, username, resourceType, database, schema, table string) (bool, error) {
	// Validate if user has access to the specified resource
	accessible, err := u.metadataRepo.IsTableAccessible(ctx, username, database, schema, table)
	if err != nil {
		return false, fmt.Errorf("failed to validate user access to resource: %w", err)
	}

	return accessible, nil
}
