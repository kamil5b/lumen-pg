package usecase

import (
	"context"

	"github.com/kamil5b/lumen-pg/internal/domain"
)

// SecurityUseCase defines operations for security-related functionality
type SecurityUseCase interface {
	// EncryptPassword encrypts a password for secure storage in cookies
	EncryptPassword(ctx context.Context, password string) (string, error)

	// DecryptPassword decrypts a password from a cookie
	DecryptPassword(ctx context.Context, encryptedPassword string) (string, error)

	// ValidateCookieIntegrity validates that a cookie has not been tampered with
	ValidateCookieIntegrity(ctx context.Context, cookieData *domain.CookieData, signature string) (bool, error)

	// GenerateCookieSignature generates a signature for cookie data
	GenerateCookieSignature(ctx context.Context, cookieData *domain.CookieData) (string, error)

	// SanitizeWhereClause sanitizes a WHERE clause to prevent SQL injection
	SanitizeWhereClause(ctx context.Context, whereClause string) (string, error)

	// ValidateQueryForInjection checks if a query contains potential SQL injection
	ValidateQueryForInjection(ctx context.Context, query string) (bool, error)

	// CheckSessionTimeout checks if a session has timed out
	CheckSessionTimeout(ctx context.Context, sessionID string) (bool, error)

	// CheckPasswordExpiry checks if a stored password is still valid
	CheckPasswordExpiry(ctx context.Context, username string, encryptedPassword string) (bool, error)

	// GenerateSecureSessionID generates a secure session ID
	GenerateSecureSessionID(ctx context.Context) (string, error)

	// IsHTTPSEnabled checks if HTTPS is enabled for the application
	IsHTTPSEnabled(ctx context.Context) (bool, error)
}
