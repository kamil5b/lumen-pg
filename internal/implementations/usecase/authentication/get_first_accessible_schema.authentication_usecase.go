package authentication

import (
	"context"
	"fmt"
)

func (u *AuthenticationUseCaseImplementation) GetFirstAccessibleSchema(ctx context.Context, username, database string) (string, error) {
	metadata, err := u.metadataRepo.GetRoleMetadata(ctx, username)
	if err != nil {
		return "", fmt.Errorf("failed to get role metadata: %w", err)
	}

	if metadata == nil || len(metadata.AccessibleSchemas) == 0 {
		return "", fmt.Errorf("no accessible schemas found")
	}

	return metadata.AccessibleSchemas[0], nil
}
