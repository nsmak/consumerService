package domain

import (
	"time"

	"github.com/nsmak/consumerService/consumer"
	"github.com/nsmak/consumerService/consumer/models"
	"github.com/nsmak/consumerService/consumer/web/auth"
	"github.com/nsmak/consumerService/consumer/web/auth/token"
	"github.com/nsmak/consumerService/consumer/web/store"

	"github.com/dgrijalva/jwt-go"
)

type domainError struct {
	IsUserError bool   `json:"-"`
	Message     string `json:"message"`
	Err         error  `json:"err,omitempty"`
}

func (e *domainError) Error() string {
	if e.Err != nil {
		e.Message = e.Message + " --> " + e.Err.Error()
	}
	return e.Message
}

func (e *domainError) UserError() bool {
	if e.Err == nil {
		return e.IsUserError
	}
	err, ok := e.Err.(consumer.Error)
	if !ok {
		return e.IsUserError
	}
	return err.UserError()
}

func (e *domainError) Unwrap() error {
	return e.Err
}

// Service - domain service for work with user's account data
type Service struct {
	store store.DataStore
	auth  *auth.Service
}

// NewService returns new instance of domain service
func NewService(s store.DataStore, a *auth.Service) *Service {
	return &Service{store: s, auth: a}
}

// CreateConsumer creates new consumer and returns his data model or error
func (s *Service) CreateConsumer(form models.RegFrom) (models.Consumer, error) {
	if err := form.Validate(); err != nil {
		return models.Consumer{}, err
	}
	isExist, err := s.store.ConsumerIsExist(form.Email)
	if err != nil {
		return models.Consumer{}, &domainError{Message: "can't get user", Err: err}
	}
	if isExist {
		return models.Consumer{}, &domainError{IsUserError: true, Message: "user is already exist"}
	}
	c, err := s.store.CreateConsumer(s.auth.MakeConsumerModel(form))
	if err != nil {
		return models.Consumer{}, &domainError{Message: "can't create user", Err: err}
	}
	return c, nil
}

// Consumer finds consumer model by email address and returns his model or error
func (s *Service) Consumer(email string) (models.Consumer, error) {
	c, err := s.store.Consumer(email)
	if err != nil {
		return models.Consumer{}, err
	}
	return c, nil
}

// AuthenticateConsumer creates authentication token and returns it or returns error
func (s *Service) AuthenticateConsumer(form models.AuthForm) (tokenString string, err error) {
	err = form.Validate()
	if err != nil {
		err = &domainError{IsUserError: true, Message: "can't authenticate user", Err: err}
		return
	}
	c, err := s.store.Consumer(form.Email)
	if err != nil {
		err = &domainError{IsUserError: false, Message: "can't get user", Err: err}
		return
	}
	err = s.auth.MatchPasswordHash(c.PassHash, form)
	if err != nil {
		err = &domainError{IsUserError: true, Message: "can't authenticate user", Err: err}
		return
	}
	claims := token.Claims{
		ID: c.ID,
		StandardClaims: &jwt.StandardClaims{
			Subject:   "access",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(24 * time.Hour * 365).Unix(),
			Issuer:    "consumer_service",
		},
	}
	tokenString, err = s.auth.CreateJWT(claims)

	return
}

// ValidateToken - validate authentication token and return error if token is not valid
func (s *Service) ValidateToken(t string) error {
	claims, err := s.auth.ParseToken(t)
	if err != nil {
		return err
	}
	if err := claims.Valid(); err != nil {
		return err
	}
	return nil
}
