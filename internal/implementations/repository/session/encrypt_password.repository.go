package session

import (
	"errors"
)

// EncryptPassword encrypts a password for cookie storage
func (s *SessionRepository) EncryptPassword(password string) (string, error) {
	return "", errors.New("not implemented yet")
}
