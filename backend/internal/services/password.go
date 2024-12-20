package services

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}

type DefaultPasswordHasher struct{}

func NewDefaultPasswordHasher() *DefaultPasswordHasher {
	return &DefaultPasswordHasher{}
}

func (h *DefaultPasswordHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func (h *DefaultPasswordHasher) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
