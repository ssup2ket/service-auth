package request

import (
	"testing"

	"github.com/ssup2ket/service-auth/internal/test"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type userSuite struct {
	suite.Suite
}

func TestInit(t *testing.T) {
	suite.Run(t, new(userSuite))
}

// UserUUID
func (u *userSuite) TestBindUserUUIDCorrect() {
	err := ValidateUserUUID(test.UserIDCorrect.String())
	require.NoError(u.T(), err)
}

func (u *userSuite) TestBindUserUUIDWrong() {
	err := ValidateUserUUID(test.UserIDWrongFormat)
	require.Error(u.T(), err)
}

// UserCreate
func (u *userSuite) TestBindUserCreateCorrect() {
	err := ValidateUserCreate(test.UserLoginIDCorrect, test.UserPasswdCorrect, string(test.UserRoleCorrect), test.UserPhoneCorrect, test.UserEmailCorrect)
	require.NoError(u.T(), err)
}

func (u *userSuite) TestBindUserCreateLoginIDWrong() {
	wrongLoginIDs := []string{test.UserLoginIDShort, test.UserLoginIDLong}
	for _, wrongID := range wrongLoginIDs {
		err := ValidateUserCreate(wrongID, test.UserPasswdCorrect, string(test.UserRoleCorrect), test.UserPhoneCorrect, test.UserEmailCorrect)
		require.Error(u.T(), err)
	}
}

func (u *userSuite) TestBindUserCreatePasswdWrong() {
	wrongPasswds := []string{test.UserPasswdShort, test.UserPasswdLong}
	for _, wrongPasswd := range wrongPasswds {
		err := ValidateUserCreate(test.UserLoginIDCorrect, wrongPasswd, string(test.UserRoleCorrect), test.UserPhoneCorrect, test.UserEmailCorrect)
		require.Error(u.T(), err)
	}
}

func (u *userSuite) TestBindUserCreateRoleWrong() {
	err := ValidateUserCreate(test.UserLoginIDCorrect, test.UserPasswdCorrect, test.UserRoleWrong, test.UserPhoneWrongFormat, test.UserEmailCorrect)
	require.Error(u.T(), err)
}

func (u *userSuite) TestBindUserCreatePhoneWrong() {
	err := ValidateUserCreate(test.UserLoginIDCorrect, test.UserPasswdCorrect, string(test.UserRoleCorrect), test.UserPhoneWrongFormat, test.UserEmailCorrect)
	require.Error(u.T(), err)
}

func (u *userSuite) TestBindUserCreateEmailWrong() {
	err := ValidateUserCreate(test.UserLoginIDCorrect, test.UserPasswdCorrect, string(test.UserRoleCorrect), test.UserPhoneCorrect, test.UserEmailWrongFormat)
	require.Error(u.T(), err)
}

// UserUpdate
func (u *userSuite) TestBindUserUpdateCorrect() {
	err := ValidateUserUpdate(test.UserIDCorrect.String(), test.UserPasswdCorrect, string(test.UserRoleCorrect), test.UserPhoneCorrect, test.UserEmailCorrect)
	require.NoError(u.T(), err)
}

func (u *userSuite) TestBindUserUpdateUUIDWrong() {
	err := ValidateUserUpdate(test.UserIDWrongFormat, test.UserPasswdCorrect, string(test.UserRoleCorrect), test.UserPhoneCorrect, test.UserEmailCorrect)
	require.Error(u.T(), err)
}

func (u *userSuite) TestBindUserUpdatePasswdWrong() {
	wrongPasswds := []string{test.UserPasswdShort, test.UserPasswdLong}
	for _, wrongPasswd := range wrongPasswds {
		err := ValidateUserUpdate(test.UserIDCorrect.String(), wrongPasswd, string(test.UserRoleCorrect), test.UserPhoneCorrect, test.UserEmailCorrect)
		require.Error(u.T(), err)
	}
}

func (u *userSuite) TestBindUserUpdateRoleWrong() {
	err := ValidateUserUpdate(test.UserIDCorrect.String(), test.UserPasswdCorrect, test.UserRoleWrong, test.UserPhoneWrongFormat, test.UserEmailCorrect)
	require.Error(u.T(), err)
}

func (u *userSuite) TestBindUserUpdatePhoneWrong() {
	err := ValidateUserUpdate(test.UserIDCorrect.String(), test.UserPasswdCorrect, string(test.UserRoleCorrect), test.UserPhoneWrongFormat, test.UserEmailCorrect)
	require.Error(u.T(), err)
}

func (u *userSuite) TestBindUserUpdateEmailWrong() {
	err := ValidateUserUpdate(test.UserIDCorrect.String(), test.UserPasswdCorrect, string(test.UserRoleCorrect), test.UserPhoneCorrect, test.UserEmailWrongFormat)
	require.Error(u.T(), err)
}
