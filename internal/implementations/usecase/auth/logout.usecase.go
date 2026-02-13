package auth

import (
	"context"
	"errors"
)

// Logout performs logout operation
func (a *AuthUseCase) Logout(ctx context.Context, sessionToken string) error {
	return errors.New("not implemented yet")
}
