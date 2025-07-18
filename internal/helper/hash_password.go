package helper

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err == errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password") {
		return errors.New("Password is wrong")
	} else if err != nil {
		return err
	}
	return nil
}