package domain_test

import (
	"testing"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/stretchr/testify/assert"
)

// UC-S7-03: Password Encryption in Cookie
func TestEncryptPassword_Stub(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	_, err := domain.EncryptPassword("mysecretpassword", key)
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
}

// UC-S7-04: Password Decryption from Cookie
func TestDecryptPassword_Stub(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	_, err := domain.DecryptPassword("encryptedvalue", key)
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
}
