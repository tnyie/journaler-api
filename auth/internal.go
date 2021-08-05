package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(password string, hash []byte) bool {
	return bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil
}

func HashPassword(password string) []byte {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return hashed
}
