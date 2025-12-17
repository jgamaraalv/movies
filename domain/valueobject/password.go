package valueobject

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const MinPasswordLength = 6

var (
	ErrEmptyPassword    = errors.New("password is required")
	ErrPasswordTooShort = errors.New("password must be at least 6 characters")
	ErrPasswordMismatch = errors.New("invalid password")
)

type Password struct {
	plainText  string
	hashedText string
}

func NewPassword(password string) (Password, error) {
	if password == "" {
		return Password{}, ErrEmptyPassword
	}

	if len(password) < MinPasswordLength {
		return Password{}, ErrPasswordTooShort
	}

	return Password{plainText: password}, nil
}

func NewHashedPassword(hashedPassword string) Password {
	return Password{hashedText: hashedPassword}
}

func (p Password) Hash() (string, error) {
	if p.hashedText != "" {
		return p.hashedText, nil
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(p.plainText), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func (p Password) Verify(plainPassword string) error {
	if p.hashedText == "" {
		return errors.New("cannot verify against non-hashed password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(p.hashedText), []byte(plainPassword))
	if err != nil {
		return ErrPasswordMismatch
	}

	return nil
}

func (p Password) IsHashed() bool {
	return p.hashedText != ""
}
