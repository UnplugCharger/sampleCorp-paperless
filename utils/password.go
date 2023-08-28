package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword  Returns the BCrypt hash of the password.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("can not hash password: %w", err)
	}
	return string(hash), nil
}

// CheckPasswordHash Returns true if the password and hash match.
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

}
