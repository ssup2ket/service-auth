package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/model"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/repo"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/repo/mocks"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/test"
)

func TestInit(t *testing.T) {
	suite.Run(t, new(userSuite))
}

type userSuite struct {
	suite.Suite

	userInfoRepo   mocks.UserInfoRepo
	userSecretRepo mocks.UserSecretRepo
	userOutboxRepo mocks.UserOutboxRepo

	userService UserService
}

func (u *userSuite) SetupTest() {
	// Init repo
	u.userInfoRepo = mocks.UserInfoRepo{}
	u.userSecretRepo = mocks.UserSecretRepo{}
	u.userOutboxRepo = mocks.UserOutboxRepo{}

	// Init service
	u.userService = NewUserServiceImp(&u.userInfoRepo, &u.userInfoRepo, &u.userSecretRepo, &u.userSecretRepo, &u.userOutboxRepo)
}

func (u *userSuite) TestGetUserSuccess() {
	// Set repo
	u.userInfoRepo.On("Get", context.Background(), test.UserIDCorrect).Return(&model.UserInfo{
		ID:      test.UserIDCorrect,
		LoginID: test.UserLoginIDCorrect,
		Phone:   test.UserPhoneCorrect,
		Email:   test.UserEmailCorrect,
	}, nil)

	// Test
	userInfo, err := u.userService.GetUser(context.Background(), test.UserIDCorrect)
	require.Nil(u.T(), err, "Failed to get user info")
	require.Equal(u.T(), test.UserIDCorrect, userInfo.ID, "Not equal UUID")
	require.Equal(u.T(), test.UserLoginIDCorrect, userInfo.LoginID, "Not equal ID")
	require.Equal(u.T(), test.UserPhoneCorrect, userInfo.Phone, "Not equal Phone")
	require.Equal(u.T(), test.UserEmailCorrect, userInfo.Email, "Not equal Eamil")
}

func (u *userSuite) TestGetUserRepoNotFound() {
	// Set repo
	repoError := repo.ErrNotFound
	u.userInfoRepo.On("Get", context.Background(), test.UserIDCorrect).Return(&model.UserInfo{}, repoError)

	// Test
	_, err := u.userService.GetUser(context.Background(), test.UserIDCorrect)
	require.Equal(u.T(), ErrRepoNotFound, err)
}
