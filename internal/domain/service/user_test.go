package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/model"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/repo"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/repo/mocks"
)

const (
	userUUIDCorrect   = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	userIDCorrect     = "test0000"
	userPasswdCorrect = "test0000"
	userPhoneCorrect  = "000-0000-0000"
	userEmailCorrect  = "test@test.com"
)

func TestInit(t *testing.T) {
	suite.Run(t, new(userSuite))
}

type userSuite struct {
	suite.Suite

	userInfoRepo   mocks.UserInfoRepo
	userSecretRepo mocks.UserSecretRepo

	userService UserService
}

func (u *userSuite) SetupTest() {
	// Init repo
	u.userInfoRepo = mocks.UserInfoRepo{}
	u.userSecretRepo = mocks.UserSecretRepo{}

	// Init service
	u.userService = NewUserServiceImp(&u.userInfoRepo, &u.userSecretRepo)
}

func (u *userSuite) TestGetUserSuccess() {
	// Set repo
	u.userInfoRepo.On("GetSecondary", context.Background(), userUUIDCorrect).Return(&model.UserInfo{
		UUID:  userUUIDCorrect,
		ID:    userIDCorrect,
		Phone: userPhoneCorrect,
		Email: userEmailCorrect,
	}, nil)

	// Test
	userInfo, err := u.userService.GetUser(context.Background(), userUUIDCorrect)
	require.Nil(u.T(), err, "Failed to get user info")
	require.Equal(u.T(), userUUIDCorrect, userInfo.UUID, "Not equal UUID")
	require.Equal(u.T(), userIDCorrect, userInfo.ID, "Not equal ID")
	require.Equal(u.T(), userPhoneCorrect, userInfo.Phone, "Not equal Phone")
	require.Equal(u.T(), userEmailCorrect, userInfo.Email, "Not equal Eamil")
}

func (u *userSuite) TestGetUserRepoNotFound() {
	// Set repo
	repoError := repo.ErrNotFound
	u.userInfoRepo.On("GetSecondary", context.Background(), userUUIDCorrect).Return(&model.UserInfo{}, repoError)

	// Test
	_, err := u.userService.GetUser(context.Background(), userUUIDCorrect)
	require.Equal(u.T(), ErrRepoNotFound, err)
}
