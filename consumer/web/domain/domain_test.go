package domain_test

import (
	"errors"
	"testing"

	"github.com/nsmak/consumerService/consumer"
	"github.com/nsmak/consumerService/consumer/models"
	"github.com/nsmak/consumerService/consumer/web/auth"
	"github.com/nsmak/consumerService/consumer/web/domain"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

const (
	validToken   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhY2Nlc3MiLCJpc3MiOiJjb25zdW1lcl9zZXJ2aWNlIiwiaWF0IjowLCJleHAiOjMyNTM1MTI5NjAwLCJpZCI6MX0.-U_HwBOfUrbNr8Ntif1ih8yKkaI65XM38ozQM2VFEbw"
	expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhY2Nlc3MiLCJpc3MiOiJjb25zdW1lcl9zZXJ2aWNlIiwiaWF0IjowLCJleHAiOjE1MDAsImlkIjoxfQ.U4bE3l7Mpz4DrlC9EmKTJ7fKXuYmrgiON6ntoD-RalU"
	invalidToken = "invalid token"
)

type DomainSuite struct {
	suite.Suite
	mockCtl   *gomock.Controller
	mockStore *MockDataStore
	auth      *auth.Service
	domain    *domain.Service
}

func (d *DomainSuite) SetupTest() {
	d.mockCtl = gomock.NewController(d.T())
	d.mockStore = NewMockDataStore(d.mockCtl)
	d.auth = auth.NewService(auth.Opts{SigningKey: []byte("1")})
	d.domain = domain.NewService(d.mockStore, d.auth)
}

func (d *DomainSuite) TeardownTest() {
	d.mockCtl.Finish()
}

func (d *DomainSuite) TestCreateConsumerSuccess() {
	reg := models.RegFrom{
		Email: "test@test.com",
		Pass1: "1234",
		Pass2: "1234",
	}

	newUser := d.auth.MakeConsumerModel(reg)

	d.mockStore.EXPECT().ConsumerIsExist(reg.Email).Return(false, nil)
	d.mockStore.EXPECT().CreateConsumer(newUser).Return(newUser, nil)
	u, err := d.domain.CreateConsumer(reg)

	d.Require().NoError(err)
	d.Require().Equal(u, newUser)
}

func (d *DomainSuite) TestCreateConsumerInvalidRegForm() {
	_, err := d.domain.CreateConsumer(models.RegFrom{})

	d.Require().Error(err)
}

func (d *DomainSuite) TestCreateConsumerUserIsExist() {
	reg := models.RegFrom{
		Email: "test@test.com",
		Pass1: "1234",
		Pass2: "1234",
	}

	d.mockStore.EXPECT().ConsumerIsExist(reg.Email).Return(true, nil)
	newUser, err := d.domain.CreateConsumer(reg)

	d.Require().Error(err)
	d.Require().NotEqual(newUser, d.auth.MakeConsumerModel(reg))
}

func (d *DomainSuite) TestCreateConsumerStoreIsExistErr() {
	reg := models.RegFrom{
		Email: "test@test.com",
		Pass1: "1234",
		Pass2: "1234",
	}
	dbErr := errors.New("db_error")

	d.mockStore.EXPECT().ConsumerIsExist(gomock.Any()).Return(false, dbErr)
	newUser, err := d.domain.CreateConsumer(reg)

	ok := errors.Is(err, dbErr)
	d.Require().True(ok)
	d.Require().NotEqual(newUser, d.auth.MakeConsumerModel(reg))
}

func (d *DomainSuite) TestCreateConsumerStoreCreateConsumerErr() {
	reg := models.RegFrom{
		Email: "test@test.com",
		Pass1: "1234",
		Pass2: "1234",
	}
	dbErr := errors.New("db_error")

	d.mockStore.EXPECT().ConsumerIsExist(reg.Email).Return(false, nil)
	d.mockStore.EXPECT().CreateConsumer(gomock.Any()).Return(models.Consumer{}, dbErr)
	newUser, err := d.domain.CreateConsumer(reg)

	ok := errors.Is(err, dbErr)
	d.Require().True(ok)
	d.Require().NotEqual(newUser, d.auth.MakeConsumerModel(reg))
}

func (d *DomainSuite) TestConsumerSuccess() {
	user := models.Consumer{
		ID:           0,
		Email:        "test@test.com",
		RegTimestamp: 0,
		PassHash:     "3",
	}

	d.mockStore.EXPECT().Consumer(user.Email).Return(user, nil)
	sUser, err := d.domain.Consumer(user.Email)

	d.Require().NoError(err)
	d.Require().Equal(sUser, user)
}

func (d *DomainSuite) TestConsumerFail() {
	emptyUser := models.Consumer{}
	dbErr := errors.New("db_error")

	d.mockStore.EXPECT().Consumer(gomock.Any()).Return(emptyUser, dbErr)
	user, err := d.domain.Consumer("test@test.com")

	d.Require().Error(err)
	d.Require().True(errors.Is(err, dbErr))
	d.Require().Equal(user, emptyUser)
}

func (d *DomainSuite) TestAuthenticateConsumerSuccess() {
	authForm := models.AuthForm{
		Email: "test@test.com",
		Pass:  "password",
	}
	user := models.Consumer{
		ID:           1,
		Email:        "test@test.com",
		RegTimestamp: 0,
		PassHash:     "test@test.compassword",
	}

	d.mockStore.EXPECT().Consumer(authForm.Email).Return(user, nil)
	token, err := d.domain.AuthenticateConsumer(authForm)

	d.Require().NoError(err)
	d.Require().NotEmpty(token)
}

func (d *DomainSuite) TestAuthenticateConsumerFormInvalid() {
	token, err := d.domain.AuthenticateConsumer(models.AuthForm{})
	_, ok := err.(consumer.Error)

	d.Require().Error(err)
	d.Require().Empty(token)
	d.Require().True(ok)
}

func (d *DomainSuite) TestAuthenticateConsumerGetConsumerErr() {
	authForm := models.AuthForm{
		Email: "test@test.com",
		Pass:  "password",
	}
	dbErr := errors.New("db_error")

	d.mockStore.EXPECT().Consumer(authForm.Email).Return(models.Consumer{}, dbErr)
	token, err := d.domain.AuthenticateConsumer(authForm)
	_, ok := err.(consumer.Error)

	d.Require().Error(err)
	d.Require().Empty(token)
	d.Require().True(ok)
	d.Require().True(errors.Is(err, dbErr))
}

func (d *DomainSuite) TestAuthenticateConsumerInvalidPassword() {
	authForm := models.AuthForm{
		Email: "test@test.com",
		Pass:  "password",
	}
	user := models.Consumer{
		ID:           1,
		Email:        "test@test.com",
		RegTimestamp: 0,
		PassHash:     "hash",
	}

	d.mockStore.EXPECT().Consumer(authForm.Email).Return(user, nil)
	token, err := d.domain.AuthenticateConsumer(authForm)
	_, ok := err.(consumer.Error)

	d.Require().Error(err)
	d.Require().Empty(token)
	d.Require().True(ok)
}

func (d *DomainSuite) TestValidateTokenSuccess() {
	err := d.domain.ValidateToken(validToken)

	d.Require().NoError(err)
}

func (d *DomainSuite) TestValidateTokenInvalidToken() {
	err := d.domain.ValidateToken(invalidToken)
	_, ok := err.(consumer.Error)

	d.Require().Error(err)
	d.Require().True(ok)
}

func (d *DomainSuite) TestValidateTokenExpiredToken() {
	err := d.domain.ValidateToken(expiredToken)

	d.Require().Error(err)
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(DomainSuite))
}
