package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

func EncodePassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(hashedPassword), nil
}

func CheckPassword(password, encrypted string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password))
	return err == nil
}
