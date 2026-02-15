package authentication

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *AuthenticationUseCaseImplementation) ValidateLoginForm(ctx context.Context, req domain.LoginRequest) ([]domain.ValidationError, error) {
	errors := []domain.ValidationError{}

	if req.Username == "" {
		errors = append(errors, domain.ValidationError{
			Field:   "username",
			Message: "username is required",
		})
	}

	if req.Password == "" {
		errors = append(errors, domain.ValidationError{
			Field:   "password",
			Message: "password is required",
		})
	}

	return errors, nil
}
