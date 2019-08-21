package users

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns hash by password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// IsMatch returns true, if hash matches password
func IsMatch(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
