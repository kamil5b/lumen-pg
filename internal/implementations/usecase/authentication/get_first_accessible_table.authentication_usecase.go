package authentication

import (
	"context"
	"fmt"
)

func (u *AuthenticationUseCaseImplementation) GetFirstAccessibleTable(ctx context.Context, username, database, schema string) (string, error) {
	metadata, err := u.metadataRepo.GetRoleMetadata(ctx, username)
	if err != nil {
		return "", fmt.Errorf("failed to get role metadata: %w", err)
	}

	if metadata == nil {
		return "", fmt.Errorf("no metadata found for user: %s", username)
	}

	// Filter tables by database and schema
	for _, table := range metadata.AccessibleTables {
		if table.Database == database && table.Schema == schema {
			return table.Name, nil
		}
	}

	return "", fmt.Errorf("no accessible tables found for user %s in database %s, schema %s", username, database, schema)
}
