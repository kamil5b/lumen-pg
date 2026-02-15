package authentication

import (
	"context"
	"fmt"
)

func (u *AuthenticationUseCaseImplementation) GetFirstAccessibleDatabase(ctx context.Context, username string) (string, error) {
	roleMetadata, err := u.metadataRepo.GetRoleMetadata(ctx, username)
	if err != nil {
		return "", fmt.Errorf("failed to get role metadata: %w", err)
	}

	if roleMetadata == nil || len(roleMetadata.AccessibleDatabases) == 0 {
		return "", fmt.Errorf("no accessible databases found for user: %s", username)
	}

	return roleMetadata.AccessibleDatabases[0], nil
}
