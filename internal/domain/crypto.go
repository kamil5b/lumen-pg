package domain

// EncryptPassword encrypts a password for cookie storage.
func EncryptPassword(password string, key []byte) (string, error) {
	return "", ErrNotImplemented
}

// DecryptPassword decrypts a password from cookie storage.
func DecryptPassword(encrypted string, key []byte) (string, error) {
	return "", ErrNotImplemented
}
