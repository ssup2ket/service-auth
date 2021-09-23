package request

import (
	"testing"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/model"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	userIDCorrect      = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	userLoginIDCorrect = "test0000"
	userRoleCorret     = model.UserRoleAdmin
	userPasswdCorrect  = "test0000"
	userPhoneCorrect   = "000-0000-0000"
	userEmailCorrect   = "test@test.com"

	userIDWrongFormat    = "aaaa-aaaa"
	userLoginIDShort     = "test0"
	userLoginIDLong      = "testtesttesttesttesttest"
	userPasswdShort      = "test0"
	userPasswdLong       = "testtesttesttesttesttest"
	userRoleWrong        = "tester"
	userPhoneWrongFormat = "00-000-00000"
	userEmailWrongFormat = "testtest.com"
)

type userSuite struct {
	suite.Suite
}

func TestInit(t *testing.T) {
	suite.Run(t, new(userSuite))
}

// UserCreate
func (h *userSuite) TestBindUserCreateCorrect() {
	err := ValidateUserCreate(userLoginIDCorrect, userPasswdCorrect, string(userRoleCorret), userPhoneCorrect, userEmailCorrect)
	require.NoError(h.T(), err)
}

func (h *userSuite) TestBindUserCreateIDWrong() {
	wrongIDs := []string{userLoginIDShort, userLoginIDLong}
	for _, wrongID := range wrongIDs {
		err := ValidateUserCreate(wrongID, userPasswdCorrect, string(userRoleCorret), userPhoneCorrect, userEmailCorrect)
		require.Error(h.T(), err)
	}
}

func (h *userSuite) TestBindUserCreatePasswdWrong() {
	wrongPasswds := []string{userPasswdShort, userPasswdLong}
	for _, wrongPasswd := range wrongPasswds {
		err := ValidateUserCreate(userLoginIDCorrect, wrongPasswd, string(userRoleCorret), userPhoneCorrect, userEmailCorrect)
		require.Error(h.T(), err)
	}
}

func (h *userSuite) TestBindUserCreateRoleWrong() {
	err := ValidateUserCreate(userLoginIDCorrect, userPasswdCorrect, userRoleWrong, userPhoneWrongFormat, userEmailCorrect)
	require.Error(h.T(), err)
}

func (h *userSuite) TestBindUserCreatePhoneWrong() {
	err := ValidateUserCreate(userLoginIDCorrect, userPasswdCorrect, string(userRoleCorret), userPhoneWrongFormat, userEmailCorrect)
	require.Error(h.T(), err)
}

func (h *userSuite) TestBindUserCreateEmailWrong() {
	err := ValidateUserCreate(userLoginIDCorrect, userPasswdCorrect, string(userRoleCorret), userPhoneCorrect, userEmailWrongFormat)
	require.Error(h.T(), err)
}

// UserUpdate
func (h *userSuite) TestBindUserUpdateCorrect() {
	err := ValidateUserUpdate(userIDCorrect, userPasswdCorrect, string(userRoleCorret), userPhoneCorrect, userEmailCorrect)
	require.NoError(h.T(), err)
}

func (h *userSuite) TestBindUserUpdateUUIDWrong() {
	err := ValidateUserUpdate(userIDWrongFormat, userPasswdCorrect, string(userRoleCorret), userPhoneCorrect, userEmailCorrect)
	require.Error(h.T(), err)
}

func (h *userSuite) TestBindUserUpdatePasswdWrong() {
	wrongPasswds := []string{userPasswdShort, userPasswdLong}
	for _, wrongPasswd := range wrongPasswds {
		err := ValidateUserUpdate(userIDCorrect, wrongPasswd, string(userRoleCorret), userPhoneCorrect, userEmailCorrect)
		require.Error(h.T(), err)
	}
}

func (h *userSuite) TestBindUserUpdateRoleWrong() {
	err := ValidateUserUpdate(userIDCorrect, userPasswdCorrect, userRoleWrong, userPhoneWrongFormat, userEmailCorrect)
	require.Error(h.T(), err)
}

func (h *userSuite) TestBindUserUpdatePhoneWrong() {
	err := ValidateUserUpdate(userIDCorrect, userPasswdCorrect, string(userRoleCorret), userPhoneWrongFormat, userEmailCorrect)
	require.Error(h.T(), err)
}

func (h *userSuite) TestBindUserUpdateEmailWrong() {
	err := ValidateUserUpdate(userIDCorrect, userPasswdCorrect, string(userRoleCorret), userPhoneCorrect, userEmailWrongFormat)
	require.Error(h.T(), err)
}
