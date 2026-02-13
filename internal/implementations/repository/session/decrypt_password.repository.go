package session

import (
	"errors"
)

// DecryptPassword decrypts a password from cookie storage
func (s *SessionRepository) DecryptPassword(encrypted string) (string, error) {
	return "", errors.New("not implemented yet")
}
