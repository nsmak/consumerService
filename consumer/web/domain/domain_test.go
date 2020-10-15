package domain_test

import (
	"errors"
	"testing"

	"github.com/nsmak/consumerService/consumer/models"
	"github.com/nsmak/consumerService/consumer/web/auth"
	"github.com/nsmak/consumerService/consumer/web/domain"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type StoreSuite struct {
	suite.Suite
	mockCtl   *gomock.Controller
	mockStore *MockDataStore
	auth      *auth.Service
	domain    *domain.Service
}

func (s *StoreSuite) SetupTest() {
	s.mockCtl = gomock.NewController(s.T())
	s.mockStore = NewMockDataStore(s.mockCtl)
	s.auth = auth.NewService(auth.Opts{SigningKey: []byte("1")})
	s.domain = domain.NewService(s.mockStore, s.auth)
}

func (s *StoreSuite) TeardownTest() {
	s.mockCtl.Finish()
}

func (s *StoreSuite) TestCreateConsumerSuccess() {
	reg := models.RegFrom{
		Email: "test@test.com",
		Pass1: "1234",
		Pass2: "1234",
	}

	newUser := s.auth.MakeConsumerModel(reg)

	s.mockStore.EXPECT().ConsumerIsExist(reg.Email).Return(false, nil)
	s.mockStore.EXPECT().CreateConsumer(newUser).Return(newUser, nil)
	u, err := s.domain.CreateConsumer(reg)

	s.Require().NoError(err)
	s.Require().Equal(u, newUser)
}

func (s *StoreSuite) TestCreateConsumerInvalidRegForm() {
	_, err := s.domain.CreateConsumer(models.RegFrom{})

	s.Require().Error(err)
}

func (s *StoreSuite) TestCreateConsumerUserIsExist() {
	reg := models.RegFrom{
		Email: "test@test.com",
		Pass1: "1234",
		Pass2: "1234",
	}

	s.mockStore.EXPECT().ConsumerIsExist(reg.Email).Return(true, nil)
	newUser, err := s.domain.CreateConsumer(reg)

	s.Require().Error(err)
	s.Require().NotEqual(newUser, s.auth.MakeConsumerModel(reg))
}

func (s *StoreSuite) TestCreateConsumerStoreIsExistErr() {
	reg := models.RegFrom{
		Email: "test@test.com",
		Pass1: "1234",
		Pass2: "1234",
	}
	dbErr := errors.New("db_error")

	s.mockStore.EXPECT().ConsumerIsExist(gomock.Any()).Return(false, dbErr)
	newUser, err := s.domain.CreateConsumer(reg)

	ok := errors.Is(err, dbErr)
	s.Require().True(ok)
	s.Require().NotEqual(newUser, s.auth.MakeConsumerModel(reg))
}

func (s *StoreSuite) TestCreateConsumerStoreCreateConsumerErr() {
	reg := models.RegFrom{
		Email: "test@test.com",
		Pass1: "1234",
		Pass2: "1234",
	}
	dbErr := errors.New("db_error")

	s.mockStore.EXPECT().ConsumerIsExist(reg.Email).Return(false, nil)
	s.mockStore.EXPECT().CreateConsumer(gomock.Any()).Return(models.Consumer{}, dbErr)
	newUser, err := s.domain.CreateConsumer(reg)

	ok := errors.Is(err, dbErr)
	s.Require().True(ok)
	s.Require().NotEqual(newUser, s.auth.MakeConsumerModel(reg))
}

func (s *StoreSuite) TestConsumerSuccess() {
	user := models.Consumer{
		ID:           0,
		Email:        "test@test.com",
		RegTimestamp: 0,
		PassHash:     "3",
	}

	s.mockStore.EXPECT().Consumer(user.Email).Return(user, nil)
	sUser, err := s.domain.Consumer(user.Email)

	s.Require().NoError(err)
	s.Require().Equal(sUser, user)
}

func (s *StoreSuite) TestConsumerFail() {
	emptyUser := models.Consumer{}
	dbErr := errors.New("db_error")

	s.mockStore.EXPECT().Consumer(gomock.Any()).Return(emptyUser, dbErr)
	user, err := s.domain.Consumer("test@test.com")

	s.Require().Error(err)
	ok := errors.Is(err, dbErr)
	s.Require().True(ok)
	s.Require().Equal(user, emptyUser)
}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(StoreSuite))
}
