package request

import (
	"testing"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/test"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type userSuite struct {
	suite.Suite
}

func TestInit(t *testing.T) {
	suite.Run(t, new(userSuite))
}

// UserCreate
func (h *userSuite) TestBindUserCreateCorrect() {
	err := ValidateUserCreate(test.UserLoginIDCorrect, test.UserPasswdCorrect, string(test.UserRoleCorrect), test.UserPhoneCorrect, test.UserEmailCorrect)
	require.NoError(h.T(), err)
}

func (h *userSuite) TestBindUserCreateIDWrong() {
	wrongIDs := []string{test.UserLoginIDShort, test.UserLoginIDLong}
	for _, wrongID := range wrongIDs {
		err := ValidateUserCreate(wrongID, test.UserPasswdCorrect, string(test.UserRoleCorrect), test.UserPhoneCorrect, test.UserEmailCorrect)
		require.Error(h.T(), err)
	}
}

func (h *userSuite) TestBindUserCreatePasswdWrong() {
	wrongPasswds := []string{test.UserPasswdShort, test.UserPasswdLong}
	for _, wrongPasswd := range wrongPasswds {
		err := ValidateUserCreate(test.UserLoginIDCorrect, wrongPasswd, string(test.UserRoleCorrect), test.UserPhoneCorrect, test.UserEmailCorrect)
		require.Error(h.T(), err)
	}
}

func (h *userSuite) TestBindUserCreateRoleWrong() {
	err := ValidateUserCreate(test.UserLoginIDCorrect, test.UserPasswdCorrect, test.UserRoleWrong, test.UserPhoneWrongFormat, test.UserEmailCorrect)
	require.Error(h.T(), err)
}

func (h *userSuite) TestBindUserCreatePhoneWrong() {
	err := ValidateUserCreate(test.UserLoginIDCorrect, test.UserPasswdCorrect, string(test.UserRoleCorrect), test.UserPhoneWrongFormat, test.UserEmailCorrect)
	require.Error(h.T(), err)
}

func (h *userSuite) TestBindUserCreateEmailWrong() {
	err := ValidateUserCreate(test.UserLoginIDCorrect, test.UserPasswdCorrect, string(test.UserRoleCorrect), test.UserPhoneCorrect, test.UserEmailWrongFormat)
	require.Error(h.T(), err)
}

// UserUpdate
func (h *userSuite) TestBindUserUpdateCorrect() {
	err := ValidateUserUpdate(test.UserIDCorrect.String(), test.UserPasswdCorrect, string(test.UserRoleCorrect), test.UserPhoneCorrect, test.UserEmailCorrect)
	require.NoError(h.T(), err)
}

func (h *userSuite) TestBindUserUpdateUUIDWrong() {
	err := ValidateUserUpdate(test.UserIDWrongFormat, test.UserPasswdCorrect, string(test.UserRoleCorrect), test.UserPhoneCorrect, test.UserEmailCorrect)
	require.Error(h.T(), err)
}

func (h *userSuite) TestBindUserUpdatePasswdWrong() {
	wrongPasswds := []string{test.UserPasswdShort, test.UserPasswdLong}
	for _, wrongPasswd := range wrongPasswds {
		err := ValidateUserUpdate(test.UserIDCorrect.String(), wrongPasswd, string(test.UserRoleCorrect), test.UserPhoneCorrect, test.UserEmailCorrect)
		require.Error(h.T(), err)
	}
}

func (h *userSuite) TestBindUserUpdateRoleWrong() {
	err := ValidateUserUpdate(test.UserIDCorrect.String(), test.UserPasswdCorrect, test.UserRoleWrong, test.UserPhoneWrongFormat, test.UserEmailCorrect)
	require.Error(h.T(), err)
}

func (h *userSuite) TestBindUserUpdatePhoneWrong() {
	err := ValidateUserUpdate(test.UserIDCorrect.String(), test.UserPasswdCorrect, string(test.UserRoleCorrect), test.UserPhoneWrongFormat, test.UserEmailCorrect)
	require.Error(h.T(), err)
}

func (h *userSuite) TestBindUserUpdateEmailWrong() {
	err := ValidateUserUpdate(test.UserIDCorrect.String(), test.UserPasswdCorrect, string(test.UserRoleCorrect), test.UserPhoneCorrect, test.UserEmailWrongFormat)
	require.Error(h.T(), err)
}
