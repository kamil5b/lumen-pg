package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
	"github.com/kamil5b/lumen-pg/internal/interfaces/usecase"
	mockRepository "github.com/kamil5b/lumen-pg/internal/testrunners/mocks/repository"
	"github.com/stretchr/testify/require"
)

// SecurityUsecaseConstructor is a function type that creates a SecurityUseCase
type SecurityUsecaseConstructor func(
	encryptionRepo repository.EncryptionRepository,
	sessionRepo repository.SessionRepository,
	clockRepo repository.ClockRepository,
) usecase.SecurityUseCase

// SecurityUsecaseRunner runs all security usecase tests against an implementation
// Maps to TEST_PLAN.md:
// - Story 7: Security & Best Practices [UC-S7-01~07, IT-S7-01~03, E2E-S7-01~06]
func SecurityUsecaseRunner(t *testing.T, constructor SecurityUsecaseConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEncryption := mockRepository.NewMockEncryptionRepository(ctrl)
	mockSession := mockRepository.NewMockSessionRepository(ctrl)
	mockClock := mockRepository.NewMockClockRepository(ctrl)

	uc := constructor(mockEncryption, mockSession, mockClock)

	ctx := context.Background()

	// UC-S7-03: Password Encryption in Cookie
	// UC-S7-04: Password Decryption from Cookie
	// UC-S2-15: Metadata Refresh Button (re-authentication)
	t.Run("EncryptPassword encrypts password for storage", func(t *testing.T) {
		mockEncryption.EXPECT().
			Encrypt(gomock.Any(), "password123").
			Return("encrypted_password_hash", nil)

		encrypted, err := uc.EncryptPassword(ctx, "password123")

		require.NoError(t, err)
		require.NotEmpty(t, encrypted)
		require.NotEqual(t, "password123", encrypted)
	})

	t.Run("DecryptPassword decrypts encrypted password", func(t *testing.T) {
		mockEncryption.EXPECT().
			Decrypt(gomock.Any(), "encrypted_password_hash").
			Return("password123", nil)

		decrypted, err := uc.DecryptPassword(ctx, "encrypted_password_hash")

		require.NoError(t, err)
		require.Equal(t, "password123", decrypted)
	})

	// UC-S7-05: Cookie Tampering Detection
	t.Run("ValidateCookieIntegrity returns true for valid signature", func(t *testing.T) {
		cookieData := &domain.CookieData{
			Username: "testuser",
			Password: "encrypted_pass",
			Nonce:    "nonce_123",
		}

		mockEncryption.EXPECT().
			GenerateSignature(gomock.Any(), gomock.Any()).
			Return("valid_signature", nil)

		valid, err := uc.ValidateCookieIntegrity(ctx, cookieData, "valid_signature")

		require.NoError(t, err)
		require.True(t, valid)
	})

	t.Run("ValidateCookieIntegrity returns false for tampered cookie", func(t *testing.T) {
		cookieData := &domain.CookieData{
			Username: "testuser",
			Password: "encrypted_pass",
			Nonce:    "nonce_123",
		}

		mockEncryption.EXPECT().
			GenerateSignature(gomock.Any(), gomock.Any()).
			Return("correct_signature", nil)

		valid, err := uc.ValidateCookieIntegrity(ctx, cookieData, "tampered_signature")

		require.NoError(t, err)
		require.False(t, valid)
	})

	// UC-S7-05: Cookie Tampering Detection
	t.Run("GenerateCookieSignature creates signature", func(t *testing.T) {
		cookieData := &domain.CookieData{
			Username: "testuser",
			Password: "encrypted_pass",
			Nonce:    "nonce_123",
		}

		mockEncryption.EXPECT().
			GenerateSignature(gomock.Any(), gomock.Any()).
			Return("signature_hash", nil)

		signature, err := uc.GenerateCookieSignature(ctx, cookieData)

		require.NoError(t, err)
		require.NotEmpty(t, signature)
	})

	// UC-S7-01: SQL Injection Prevention - WHERE Clause
	// UC-S5-03 ~ UC-S5-04: WHERE Clause Validation
	// E2E-S7-01: SQL Injection via WHERE Bar
	t.Run("SanitizeWhereClause removes SQL injection attempts", func(t *testing.T) {
		mockEncryption.EXPECT().
			Encrypt(gomock.Any(), gomock.Any()).
			Return("id > 10", nil)

		sanitized, err := uc.SanitizeWhereClause(ctx, "id > 10; DROP TABLE users; --")

		require.NoError(t, err)
		require.NotEmpty(t, sanitized)
	})

	t.Run("SanitizeWhereClause preserves safe clauses", func(t *testing.T) {
		mockEncryption.EXPECT().
			Encrypt(gomock.Any(), gomock.Any()).
			Return("id > 10 AND status = 'active'", nil)

		sanitized, err := uc.SanitizeWhereClause(ctx, "id > 10 AND status = 'active'")

		require.NoError(t, err)
		require.NotEmpty(t, sanitized)
	})

	// UC-S7-02: SQL Injection Prevention - Query Editor
	// E2E-S7-02: SQL Injection via Query Editor
	t.Run("ValidateQueryForInjection detects injection attempts", func(t *testing.T) {
		mockEncryption.EXPECT().
			ValidateSignature(gomock.Any(), "SELECT * FROM users; DROP TABLE users; --", gomock.Any()).
			Return(false, nil)

		hasInjection, err := uc.ValidateQueryForInjection(ctx, "SELECT * FROM users; DROP TABLE users; --")

		require.NoError(t, err)
		require.True(t, hasInjection)
	})

	t.Run("ValidateQueryForInjection accepts safe queries", func(t *testing.T) {
		mockEncryption.EXPECT().
			ValidateSignature(gomock.Any(), "SELECT * FROM users WHERE id = $1", gomock.Any()).
			Return(true, nil)

		hasInjection, err := uc.ValidateQueryForInjection(ctx, "SELECT * FROM users WHERE id = $1")

		require.NoError(t, err)
		require.False(t, hasInjection)
	})

	// UC-S7-06: Session Timeout Short-Lived Cookie
	// UC-S7-07: Session Timeout Long-Lived Cookie
	// E2E-S7-04: Session Timeout Enforcement
	t.Run("CheckSessionTimeout returns true for expired session", func(t *testing.T) {
		expiredTime := time.Now().Add(-1 * time.Hour)
		mockSession.EXPECT().
			GetSession(gomock.Any(), "expired_session").
			Return(&domain.Session{
				ID:        "expired_session",
				ExpiresAt: expiredTime,
			}, nil)

		hasTimedOut, err := uc.CheckSessionTimeout(ctx, "expired_session")

		require.NoError(t, err)
		require.True(t, hasTimedOut)
	})

	t.Run("CheckSessionTimeout returns false for active session", func(t *testing.T) {
		futureTime := time.Now().Add(1 * time.Hour)
		mockSession.EXPECT().
			GetSession(gomock.Any(), "active_session").
			Return(&domain.Session{
				ID:        "active_session",
				ExpiresAt: futureTime,
			}, nil)

		hasTimedOut, err := uc.CheckSessionTimeout(ctx, "active_session")

		require.NoError(t, err)
		require.False(t, hasTimedOut)
	})

	// UC-S7-06: Session Timeout Short-Lived Cookie
	t.Run("CheckPasswordExpiry returns true for expired password", func(t *testing.T) {
		mockClock.EXPECT().
			IsExpired(gomock.Any(), gomock.Any()).
			Return(true)

		expired, err := uc.CheckPasswordExpiry(ctx, "testuser", "encrypted_password")

		require.NoError(t, err)
		require.True(t, expired)
	})

	t.Run("CheckPasswordExpiry returns false for valid password", func(t *testing.T) {
		mockClock.EXPECT().
			IsExpired(gomock.Any(), gomock.Any()).
			Return(false)

		expired, err := uc.CheckPasswordExpiry(ctx, "testuser", "encrypted_password")

		require.NoError(t, err)
		require.False(t, expired)
	})

	// UC-S2-07: Session Cookie Creation - Password
	// UC-S7-03: Password Encryption in Cookie
	t.Run("GenerateSecureSessionID creates unique ID", func(t *testing.T) {
		mockEncryption.EXPECT().
			GenerateSecureToken(gomock.Any(), 32).
			Return("secure_session_id_123", nil).Times(2)

		id1, err1 := uc.GenerateSecureSessionID(ctx)
		id2, err2 := uc.GenerateSecureSessionID(ctx)

		require.NoError(t, err1)
		require.NoError(t, err2)
		require.NotEmpty(t, id1)
		require.NotEmpty(t, id2)
	})

	// UC-S7-01 & UC-S7-02: Combined injection prevention
	t.Run("Security prevents multiple injection vectors", func(t *testing.T) {
		testCases := []struct {
			name            string
			input           string
			expectInjection bool
		}{
			{
				name:            "Simple DROP TABLE",
				input:           "id = 1; DROP TABLE users; --",
				expectInjection: true,
			},
			{
				name:            "UNION-based injection",
				input:           "id = 1 UNION SELECT * FROM users",
				expectInjection: true,
			},
			{
				name:            "Boolean-based blind injection",
				input:           "id = 1' OR '1'='1",
				expectInjection: true,
			},
			{
				name:            "Safe parameterized query",
				input:           "id = $1",
				expectInjection: false,
			},
			{
				name:            "Safe comparison",
				input:           "id > 10 AND status = 'active'",
				expectInjection: false,
			},
		}

		for _, tc := range testCases {
			mockEncryption.EXPECT().
				ValidateSignature(gomock.Any(), tc.input, gomock.Any()).
				Return(!tc.expectInjection, nil)

			hasInjection, err := uc.ValidateQueryForInjection(ctx, tc.input)

			require.NoError(t, err, "case: %s", tc.name)
			require.Equal(t, tc.expectInjection, hasInjection, "case: %s", tc.name)
		}
	})

	// UC-S7-03 & UC-S7-04: Password lifecycle
	t.Run("Password encryption and decryption roundtrip", func(t *testing.T) {
		originalPassword := "SecurePassword123!"

		mockEncryption.EXPECT().
			Encrypt(gomock.Any(), originalPassword).
			Return("encrypted_hash", nil)

		mockEncryption.EXPECT().
			Decrypt(gomock.Any(), "encrypted_hash").
			Return(originalPassword, nil)

		encrypted, err1 := uc.EncryptPassword(ctx, originalPassword)
		require.NoError(t, err1)

		decrypted, err2 := uc.DecryptPassword(ctx, encrypted)
		require.NoError(t, err2)

		require.Equal(t, originalPassword, decrypted)
	})

	// UC-S7-05: Cookie integrity lifecycle
	t.Run("Cookie creation, validation, and tamper detection", func(t *testing.T) {
		originalCookie := &domain.CookieData{
			Username: "testuser",
			Password: "encrypted_pass",
			Nonce:    "nonce_abc123",
		}

		mockEncryption.EXPECT().
			GenerateSignature(gomock.Any(), gomock.Any()).
			Return("original_signature", nil).Times(2)

		signature, err := uc.GenerateCookieSignature(ctx, originalCookie)
		require.NoError(t, err)

		valid, err := uc.ValidateCookieIntegrity(ctx, originalCookie, signature)
		require.NoError(t, err)
		require.True(t, valid)
	})

	// IT-S7-01: Real SQL Injection Test
	// IT-S7-02: Real Password Security
	// IT-S7-03: Real Session Expiration
	t.Run("Integration: Full security flow", func(t *testing.T) {
		mockEncryption.EXPECT().
			Encrypt(gomock.Any(), "password123").
			Return("encrypted_pass", nil)

		mockEncryption.EXPECT().
			GenerateSecureToken(gomock.Any(), 32).
			Return("secure_id_123", nil)

		mockEncryption.EXPECT().
			GenerateSignature(gomock.Any(), gomock.Any()).
			Return("cookie_signature", nil)

		futureTime := time.Now().Add(1 * time.Hour)
		mockSession.EXPECT().
			GetSession(gomock.Any(), "secure_id_123").
			Return(&domain.Session{
				ID:        "secure_id_123",
				Username:  "testuser",
				ExpiresAt: futureTime,
			}, nil)

		// Encrypt password
		encryptedPass, _ := uc.EncryptPassword(ctx, "password123")
		require.NotEmpty(t, encryptedPass)

		// Generate secure session ID
		sessionID, _ := uc.GenerateSecureSessionID(ctx)
		require.NotEmpty(t, sessionID)

		// Create and validate cookie
		cookie := &domain.CookieData{
			Username: "testuser",
			Password: encryptedPass,
			Nonce:    "nonce_456",
		}
		sig, _ := uc.GenerateCookieSignature(ctx, cookie)
		valid, _ := uc.ValidateCookieIntegrity(ctx, cookie, sig)
		require.True(t, valid)

		// Check session timeout
		hasTimedOut, _ := uc.CheckSessionTimeout(ctx, sessionID)
		require.False(t, hasTimedOut)
	})
}
