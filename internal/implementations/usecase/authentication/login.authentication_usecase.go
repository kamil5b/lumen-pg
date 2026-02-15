package authentication

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

func (u *AuthenticationUseCaseImplementation) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	// Validate the login form
	validationErrors, err := u.ValidateLoginForm(ctx, req)
	if err != nil {
		return &domain.LoginResponse{Success: false, Message: err.Error()}, err
	}

	if len(validationErrors) > 0 {
		return &domain.LoginResponse{Success: false, Message: "validation failed"}, nil
	}

	// Probe the connection with the provided credentials
	success, err := u.ProbeConnection(ctx, req.Username, req.Password)
	if err != nil || !success {
		return &domain.LoginResponse{Success: false, Message: "invalid credentials"}, err
	}

	// Check if the user has accessible databases
	databases, err := u.rbacRepo.GetAccessibleDatabases(ctx, req.Username)
	if err != nil {
		return &domain.LoginResponse{Success: false, Message: err.Error()}, err
	}

	if len(databases) == 0 {
		return &domain.LoginResponse{Success: false, Message: "no accessible databases"}, nil
	}

	return &domain.LoginResponse{
		Success:  true,
		Message:  "login successful",
		Username: req.Username,
	}, nil
}
