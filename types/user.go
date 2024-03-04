package types

import (
	"fmt"
	"net/mail"
	"unicode/utf8"
)

type User struct {
	ID       int     `json:"id" db:"id"`
	Name     *string `json:"name" db:"name"`
	Email    string  `json:"email" db:"email"`
	Password string  `json:"-" db:"password"`
}

type CreateUser struct {
	Email    string `json:"email" `
	Password string `json:"password"`
}

func (createUser *CreateUser) Validate() error {
	if _, err := mail.ParseAddress(createUser.Email); err != nil {
		return fmt.Errorf("invalid email")
	}

	if utf8.RuneCountInString(createUser.Password) < 8 {
		return fmt.Errorf("password must contain at least 8 symbols")
	}

	return nil
}
