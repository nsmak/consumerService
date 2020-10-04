package auth

import (
	"errors"

	"consumerService/consumer"
	"consumerService/consumer/models"
	"consumerService/consumer/web/auth/token"

	"github.com/dgrijalva/jwt-go"
)

type authError struct {
	IsUserError bool   `json:"-"`
	Message     string `json:"message"`
	Err         error  `json:"err,omitempty"`
}

func (e *authError) Error() string {
	if e.Err != nil {
		e.Message = e.Message + " --> " + e.Err.Error()
	}
	return e.Message
}

func (e *authError) UserError() bool {
	if e.Err == nil {
		return e.IsUserError
	}
	err, ok := e.Err.(consumer.Error)
	if !ok {
		return e.IsUserError
	}
	return err.UserError()
}

func (e *authError) Unwrap() error {
	return e.Err
}

// Opts - auth service options stricture
type Opts struct {
	SigningKey []byte
}

// Service implements authentication and auth-helping methods
type Service struct {
	Opts
}

// NewService returns new instance of authentication service
func NewService(opts Opts) *Service {
	return &Service{Opts: opts}
}

// MakeConsumerModel creates domain model from registration form
func (s *Service) MakeConsumerModel(form models.RegFrom) models.Consumer {
	return models.Consumer{
		Email:    form.Email,
		PassHash: s.hashFrom(form.Email, form.Pass1),
	}
}

// MatchPasswordHash returns error if result of hashing email and password doesn't match with input hash
func (s *Service) MatchPasswordHash(hash string, form models.AuthForm) error {
	in := s.hashFrom(form.Email, form.Pass)
	if hash != in {
		return &authError{IsUserError: true, Message: "invalid password"}
	}
	return nil
}

func (s *Service) hashFrom(email, password string) string {
	// TODO: Not secure! Must be implemented hash algorithm
	return email + password
}

// CreateJWT returns string of jwt
func (s *Service) CreateJWT(claims token.Claims) (string, error) {
	if claims.ID == 0 {
		return "", &authError{IsUserError: true, Message: "invalid uid"}
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := t.SignedString(s.SigningKey)
	if err != nil {
		return "", &authError{Message: "can's signing token", Err: err}
	}
	return str, nil
}

// ParseToken claims from token string and returns its
func (s *Service) ParseToken(str string) (token.Claims, error) {
	parser := jwt.Parser{SkipClaimsValidation: true}
	t, err := parser.ParseWithClaims(str, &token.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("signing method is not valid")
		}
		return s.SigningKey, nil
	})
	if err != nil {
		return token.Claims{}, &authError{Message: "can't parse token", Err: err}
	}

	claims, ok := t.Claims.(*token.Claims)
	if !ok {
		return token.Claims{}, &authError{IsUserError: false, Message: "invalid token claims type"}
	}
	return *claims, nil
}
