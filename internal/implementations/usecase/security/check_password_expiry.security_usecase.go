package security

import (
	"context"
)

func (u *SecurityUseCaseImplementation) CheckPasswordExpiry(ctx context.Context, username string, encryptedPassword string) (bool, error) {
	// Check if the password has expired using the clock repository
	// Use a default expiration time of 24 hours from now
	expirationTime := int64(86400) // 24 hours in seconds

	isExpired := u.clockRepo.IsExpired(ctx, expirationTime)

	return isExpired, nil
}
