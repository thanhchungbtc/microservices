package services

import "golang.org/x/crypto/bcrypt"

type Password struct {
}

func NewPassword() *Password {
	return &Password{}
}

func (Password) ToHash(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

func (Password) Compare(storedPassword string, suppliedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(suppliedPassword))
}
