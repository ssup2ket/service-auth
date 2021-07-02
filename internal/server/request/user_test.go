package request

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	userUUIDCorrect   = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	userIDCorrect     = "test0000"
	userPasswdCorrect = "test0000"
	userPhoneCorrect  = "000-0000-0000"
	userEmailCorrect  = "test@test.com"

	userUUIDWrongFormat  = "aaaa-aaaa"
	userIDShort          = "test0"
	userIDLong           = "testtesttesttesttesttest"
	userPasswdShort      = "test0"
	userPasswdLong       = "testtesttesttesttesttest"
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
	err := ValidateUserCreate(userIDCorrect, userPasswdCorrect, userPhoneCorrect, userEmailCorrect)
	require.NoError(h.T(), err)
}

func (h *userSuite) TestBindUserCreateIDWrong() {
	wrongIDs := []string{userIDShort, userIDLong}
	for _, wrongID := range wrongIDs {
		err := ValidateUserCreate(wrongID, userPasswdCorrect, userPhoneCorrect, userEmailCorrect)
		require.Error(h.T(), err)
	}
}

func (h *userSuite) TestBindUserCreatePasswdWrong() {
	wrongPasswds := []string{userPasswdShort, userPasswdLong}
	for _, wrongPasswd := range wrongPasswds {
		err := ValidateUserCreate(userIDCorrect, wrongPasswd, userPhoneCorrect, userEmailCorrect)
		require.Error(h.T(), err)
	}
}

func (h *userSuite) TestBindUserCreatePhoneWrong() {
	err := ValidateUserCreate(userIDCorrect, userPasswdCorrect, userPhoneWrongFormat, userEmailCorrect)
	require.Error(h.T(), err)
}

func (h *userSuite) TestBindUserCreateEmailWrong() {
	err := ValidateUserCreate(userIDCorrect, userPasswdCorrect, userPhoneCorrect, userEmailWrongFormat)
	require.Error(h.T(), err)
}

// UserUpdate
func (h *userSuite) TestBindUserUpdateCorrect() {
	err := ValidateUserUpdate(userUUIDCorrect, userPasswdCorrect, userPhoneCorrect, userEmailCorrect)
	require.NoError(h.T(), err)
}

func (h *userSuite) TestBindUserUpdateUUIDWrong() {
	err := ValidateUserUpdate(userUUIDWrongFormat, userPasswdCorrect, userPhoneCorrect, userEmailCorrect)
	require.Error(h.T(), err)
}

func (h *userSuite) TestBindUserUpdatePasswdWrong() {
	wrongPasswds := []string{userPasswdShort, userPasswdLong}
	for _, wrongPasswd := range wrongPasswds {
		err := ValidateUserUpdate(userUUIDCorrect, wrongPasswd, userPhoneCorrect, userEmailCorrect)
		require.Error(h.T(), err)
	}
}

func (h *userSuite) TestBindUserUpdatePhoneWrong() {
	err := ValidateUserUpdate(userUUIDCorrect, userPasswdCorrect, userPhoneWrongFormat, userEmailCorrect)
	require.Error(h.T(), err)
}

func (h *userSuite) TestBindUserUpdateEmailWrong() {
	err := ValidateUserUpdate(userUUIDCorrect, userPasswdCorrect, userPhoneCorrect, userEmailWrongFormat)
	require.Error(h.T(), err)
}
