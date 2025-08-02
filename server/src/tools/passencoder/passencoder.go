package passencoder

import (
	"golang.org/x/crypto/bcrypt"
)

func CreateHashPassword(password string) (string, error) {
	salt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(salt), nil
}

func CheckHashPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
