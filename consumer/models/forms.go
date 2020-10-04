package models

import (
	"errors"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Consumer model
type Consumer struct {
	ID           int
	Email        string
	RegTimestamp int64
	PassHash     string
}

// RegFrom is used in registration process
type RegFrom struct {
	Email string
	Pass1 string
	Pass2 string
}

// Validate registration form
func (f RegFrom) Validate() error {
	if f.Email == "" {
		return errors.New("empty email")
	}
	if !emailRegex.MatchString(f.Email) {
		return errors.New("invalid email")
	}
	if f.Pass1 == "" {
		return errors.New("empty password")
	}
	if f.Pass1 != f.Pass2 {
		return errors.New("passwords don't match")
	}
	return nil
}

// AuthForm is used in authentication process
type AuthForm struct {
	Email string
	Pass  string
}

// Validate authentication form
func (f AuthForm) Validate() error {
	if f.Email == "" {
		return errors.New("empty email")
	}
	if !emailRegex.MatchString(f.Email) {
		return errors.New("email is not valid")
	}
	if f.Pass == "" {
		return errors.New("empty password")
	}
	return nil
}
