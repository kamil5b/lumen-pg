package security

import (
	"context"
)

func (u *SecurityUseCaseImplementation) IsHTTPSEnabled(ctx context.Context) (bool, error) {
	// For now, return true as HTTPS should be enabled in production
	// This can be extended to check environment variables or configuration
	return true, nil
}
