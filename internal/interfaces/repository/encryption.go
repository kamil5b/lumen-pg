package repository

import (
	"context"
)

// EncryptionRepository defines operations for encryption and decryption
type EncryptionRepository interface {
	// Encrypt encrypts data using a secure cipher
	Encrypt(ctx context.Context, plaintext string) (string, error)

	// Decrypt decrypts encrypted data
	Decrypt(ctx context.Context, ciphertext string) (string, error)

	// GenerateNonce generates a random nonce for encryption
	GenerateNonce(ctx context.Context) (string, error)

	// ValidateSignature validates a signature for data integrity
	ValidateSignature(ctx context.Context, data string, signature string) (bool, error)

	// GenerateSignature generates a signature for data
	GenerateSignature(ctx context.Context, data string) (string, error)

	// HashPassword generates a hash of a password
	HashPassword(ctx context.Context, password string) (string, error)

	// ComparePasswordHash compares a password with a hash
	ComparePasswordHash(ctx context.Context, password, hash string) (bool, error)

	// GenerateSecureToken generates a cryptographically secure random token
	GenerateSecureToken(ctx context.Context, length int) (string, error)
}

// LoggerRepository defines operations for logging
type LoggerRepository interface {
	// LogInfo logs an informational message
	LogInfo(ctx context.Context, message string, fields map[string]interface{}) error

	// LogWarn logs a warning message
	LogWarn(ctx context.Context, message string, fields map[string]interface{}) error

	// LogError logs an error message
	LogError(ctx context.Context, message string, err error, fields map[string]interface{}) error

	// LogDebug logs a debug message
	LogDebug(ctx context.Context, message string, fields map[string]interface{}) error

	// LogSecurityEvent logs a security-related event
	LogSecurityEvent(ctx context.Context, eventType string, username string, details map[string]interface{}) error

	// LogQueryExecution logs a query execution
	LogQueryExecution(ctx context.Context, username string, query string, executionTimeMs int64, success bool, error error) error

	// LogTransactionEvent logs a transaction event
	LogTransactionEvent(ctx context.Context, username string, eventType string, details map[string]interface{}) error
}

// CacheRepository defines operations for caching
type CacheRepository interface {
	// Set stores a value in the cache with a key
	Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error

	// Get retrieves a value from the cache by key
	Get(ctx context.Context, key string) (interface{}, error)

	// Delete removes a value from the cache
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists in the cache
	Exists(ctx context.Context, key string) (bool, error)

	// Clear clears all values from the cache
	Clear(ctx context.Context) error

	// SetWithExpiration sets a value with an expiration time
	SetWithExpiration(ctx context.Context, key string, value interface{}, expirationTime int64) error

	// GetAndDelete retrieves a value and deletes it atomically
	GetAndDelete(ctx context.Context, key string) (interface{}, error)
}

// ClockRepository defines operations for time-related functionality
type ClockRepository interface {
	// Now returns the current time in Unix seconds
	Now(ctx context.Context) int64

	// NowNano returns the current time in nanoseconds
	NowNano(ctx context.Context) int64

	// AddSeconds adds seconds to the current time and returns Unix timestamp
	AddSeconds(ctx context.Context, seconds int64) int64

	// IsExpired checks if a given Unix timestamp is in the past
	IsExpired(ctx context.Context, expirationTime int64) bool

	// TimeUntilExpiration returns seconds until a given expiration time
	TimeUntilExpiration(ctx context.Context, expirationTime int64) int64
}
