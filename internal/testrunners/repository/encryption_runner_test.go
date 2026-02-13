package repository

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// EncryptionRepositoryConstructor is a function type that creates an EncryptionRepository
type EncryptionRepositoryConstructor func(db *sql.DB) repository.EncryptionRepository

// EncryptionRepositoryRunner runs all encryption repository tests against an implementation
// Maps to TEST_PLAN.md:
// - Story 2: Authentication & Identity [UC-S2-06, UC-S2-07: Session Cookie Creation with Password]
// - Story 7: Security & Best Practices [UC-S7-03~07, IT-S7-02, E2E-S7-01~06]
func EncryptionRepositoryRunner(t *testing.T, constructor EncryptionRepositoryConstructor) {
	t.Helper()

	ctx := context.Background()

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
	)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	connStr, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	err = db.PingContext(ctx)
	require.NoError(t, err)

	repo := constructor(db)

	// UC-S7-03: Password Encryption in Cookie
	// IT-S7-02: Real Password Security
	t.Run("Encrypt encrypts plaintext", func(t *testing.T) {
		plaintext := "mysecretpassword"
		ciphertext, err := repo.Encrypt(ctx, plaintext)
		require.NoError(t, err)
		require.NotEmpty(t, ciphertext)
		require.NotEqual(t, plaintext, ciphertext)
	})

	// UC-S7-04: Password Decryption from Cookie
	t.Run("Decrypt decrypts ciphertext", func(t *testing.T) {
		plaintext := "mysecretpassword"
		ciphertext, err := repo.Encrypt(ctx, plaintext)
		require.NoError(t, err)

		decrypted, err := repo.Decrypt(ctx, ciphertext)
		require.NoError(t, err)
		require.Equal(t, plaintext, decrypted)
	})

	// UC-S7-05: Cookie Tampering Detection
	// E2E-S7-03: Cookie Tampering Prevention
	t.Run("Decrypt with invalid ciphertext returns error", func(t *testing.T) {
		_, err := repo.Decrypt(ctx, "invalid_ciphertext_data")
		require.Error(t, err)
	})

	t.Run("Encrypt and decrypt multiple values", func(t *testing.T) {
		values := []string{"password1", "password2", "password3"}

		ciphertexts := make([]string, len(values))
		for i, val := range values {
			ct, err := repo.Encrypt(ctx, val)
			require.NoError(t, err)
			ciphertexts[i] = ct
		}

		for i, ct := range ciphertexts {
			decrypted, err := repo.Decrypt(ctx, ct)
			require.NoError(t, err)
			require.Equal(t, values[i], decrypted)
		}
	})

	t.Run("GenerateNonce generates unique nonces", func(t *testing.T) {
		nonce1, err := repo.GenerateNonce(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, nonce1)

		nonce2, err := repo.GenerateNonce(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, nonce2)

		require.NotEqual(t, nonce1, nonce2)
	})

	t.Run("GenerateNonce generates valid nonces", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			nonce, err := repo.GenerateNonce(ctx)
			require.NoError(t, err)
			require.NotEmpty(t, nonce)
		}
	})

	// UC-S7-05: Cookie Tampering Detection
	t.Run("GenerateSignature generates signature for data", func(t *testing.T) {
		data := "important data"
		signature, err := repo.GenerateSignature(ctx, data)
		require.NoError(t, err)
		require.NotEmpty(t, signature)
	})

	// UC-S7-05: Cookie Tampering Detection
	// E2E-S7-03: Cookie Tampering Prevention
	t.Run("ValidateSignature validates correct signature", func(t *testing.T) {
		data := "important data"
		signature, err := repo.GenerateSignature(ctx, data)
		require.NoError(t, err)

		valid, err := repo.ValidateSignature(ctx, data, signature)
		require.NoError(t, err)
		require.True(t, valid)
	})

	// UC-S7-05: Cookie Tampering Detection
	t.Run("ValidateSignature rejects invalid signature", func(t *testing.T) {
		data := "important data"
		invalidSignature := "invalid_signature_data"

		valid, err := repo.ValidateSignature(ctx, data, invalidSignature)
		require.NoError(t, err)
		require.False(t, valid)
	})

	// UC-S7-05: Cookie Tampering Detection
	// E2E-S7-03: Cookie Tampering Prevention
	t.Run("ValidateSignature rejects tampered data", func(t *testing.T) {
		data := "important data"
		signature, err := repo.GenerateSignature(ctx, data)
		require.NoError(t, err)

		tamperedData := "tampered data"
		valid, err := repo.ValidateSignature(ctx, tamperedData, signature)
		require.NoError(t, err)
		require.False(t, valid)
	})

	// UC-S7-03: Password Encryption in Cookie
	// IT-S7-02: Real Password Security
	t.Run("HashPassword generates password hash", func(t *testing.T) {
		password := "mypassword123"
		hash, err := repo.HashPassword(ctx, password)
		require.NoError(t, err)
		require.NotEmpty(t, hash)
		require.NotEqual(t, password, hash)
	})

	// UC-S7-04: Password Decryption from Cookie
	// IT-S7-02: Real Password Security
	t.Run("ComparePasswordHash verifies correct password", func(t *testing.T) {
		password := "mypassword123"
		hash, err := repo.HashPassword(ctx, password)
		require.NoError(t, err)

		matches, err := repo.ComparePasswordHash(ctx, password, hash)
		require.NoError(t, err)
		require.True(t, matches)
	})

	// IT-S7-02: Real Password Security
	t.Run("ComparePasswordHash rejects incorrect password", func(t *testing.T) {
		password := "mypassword123"
		hash, err := repo.HashPassword(ctx, password)
		require.NoError(t, err)

		wrongPassword := "wrongpassword"
		matches, err := repo.ComparePasswordHash(ctx, wrongPassword, hash)
		require.NoError(t, err)
		require.False(t, matches)
	})

	// IT-S7-02: Real Password Security
	t.Run("ComparePasswordHash rejects invalid hash", func(t *testing.T) {
		password := "mypassword123"
		invalidHash := "invalid_hash_data"

		matches, err := repo.ComparePasswordHash(ctx, password, invalidHash)
		require.Error(t, err)
		require.False(t, matches)
	})

	// UC-S2-06: Session Cookie Creation - Username
	// UC-S2-07: Session Cookie Creation - Password
	t.Run("GenerateSecureToken generates random tokens", func(t *testing.T) {
		token1, err := repo.GenerateSecureToken(ctx, 32)
		require.NoError(t, err)
		require.NotEmpty(t, token1)
		require.Len(t, token1, 32)

		token2, err := repo.GenerateSecureToken(ctx, 32)
		require.NoError(t, err)
		require.NotEmpty(t, token2)

		require.NotEqual(t, token1, token2)
	})

	t.Run("GenerateSecureToken generates tokens of different lengths", func(t *testing.T) {
		token16, err := repo.GenerateSecureToken(ctx, 16)
		require.NoError(t, err)
		require.Len(t, token16, 16)

		token64, err := repo.GenerateSecureToken(ctx, 64)
		require.NoError(t, err)
		require.Len(t, token64, 64)
	})

	t.Run("GenerateSecureToken generates cryptographically random tokens", func(t *testing.T) {
		tokens := make(map[string]bool)
		for i := 0; i < 10; i++ {
			token, err := repo.GenerateSecureToken(ctx, 32)
			require.NoError(t, err)
			require.NotContains(t, tokens, token)
			tokens[token] = true
		}
	})

	// UC-S7-03: Password Encryption in Cookie
	// E2E-S7-01: SQL Injection via WHERE Bar (encryption security)
	t.Run("Multiple encryptions of same plaintext produce different ciphertexts", func(t *testing.T) {
		plaintext := "test_password"
		ct1, err := repo.Encrypt(ctx, plaintext)
		require.NoError(t, err)

		ct2, err := repo.Encrypt(ctx, plaintext)
		require.NoError(t, err)

		// Different encryptions should produce different ciphertexts due to IV/nonce
		require.NotEqual(t, ct1, ct2)

		// But both should decrypt to same value
		pt1, _ := repo.Decrypt(ctx, ct1)
		pt2, _ := repo.Decrypt(ctx, ct2)
		require.Equal(t, pt1, pt2)
		require.Equal(t, plaintext, pt1)
	})

	// IT-S7-02: Real Password Security
	t.Run("Password hashes are unique for same password", func(t *testing.T) {
		password := "test_password"
		hash1, err := repo.HashPassword(ctx, password)
		require.NoError(t, err)

		hash2, err := repo.HashPassword(ctx, password)
		require.NoError(t, err)

		// Hashes should be different but both match password
		require.NotEqual(t, hash1, hash2)

		match1, _ := repo.ComparePasswordHash(ctx, password, hash1)
		match2, _ := repo.ComparePasswordHash(ctx, password, hash2)

		require.True(t, match1)
		require.True(t, match2)
	})

	// UC-S7-05: Cookie Tampering Detection
	t.Run("Signature generation is deterministic", func(t *testing.T) {
		data := "consistent_data"
		sig1, err := repo.GenerateSignature(ctx, data)
		require.NoError(t, err)

		sig2, err := repo.GenerateSignature(ctx, data)
		require.NoError(t, err)

		require.Equal(t, sig1, sig2)
	})

	// UC-S7-03: Password Encryption in Cookie
	t.Run("Empty plaintext can be encrypted", func(t *testing.T) {
		ciphertext, err := repo.Encrypt(ctx, "")
		require.NoError(t, err)
		require.NotEmpty(t, ciphertext)

		decrypted, err := repo.Decrypt(ctx, ciphertext)
		require.NoError(t, err)
		require.Equal(t, "", decrypted)
	})

	// UC-S7-03: Password Encryption in Cookie
	t.Run("Large plaintext can be encrypted", func(t *testing.T) {
		plaintext := ""
		for i := 0; i < 1000; i++ {
			plaintext += "a"
		}

		ciphertext, err := repo.Encrypt(ctx, plaintext)
		require.NoError(t, err)

		decrypted, err := repo.Decrypt(ctx, ciphertext)
		require.NoError(t, err)
		require.Equal(t, plaintext, decrypted)
	})
}
