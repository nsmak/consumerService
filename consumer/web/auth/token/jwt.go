package token

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

// Claims - custom jwt token claims
type Claims struct {
	ID int
	*jwt.StandardClaims
}

// Valid claims
func (c Claims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}
	if c.ID == 0 {
		return errors.New("bad user id")
	}
	return nil
}
