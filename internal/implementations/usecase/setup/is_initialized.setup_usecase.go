package setup

import (
	"context"
)

func (u *SetupUseCaseImplementation) IsInitialized(ctx context.Context) (bool, error) {
	// Check if metadata exists by trying to retrieve it
	// We need a database name, but since we don't have one, we'll use empty string
	// and let the repository determine if it has been initialized
	metadata, err := u.metadataRepo.GetMetadata(ctx, "")
	if err != nil {
		// If there's an error, the system is not initialized
		return false, nil
	}

	// If metadata exists and is not nil, the system is initialized
	if metadata != nil {
		return true, nil
	}

	return false, nil
}
