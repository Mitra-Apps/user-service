package tools

import "golang.org/x/crypto/bcrypt"

type Bcrypt struct{}

//go:generate mockgen -source=bcrypt.go -destination=mock/bcrypt.go -package=mock
type BcryptInterface interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword, password []byte) error
}

func New(b *Bcrypt) BcryptInterface {
	return b
}

func (h *Bcrypt) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (h *Bcrypt) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
