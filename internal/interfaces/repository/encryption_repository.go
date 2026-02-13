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
