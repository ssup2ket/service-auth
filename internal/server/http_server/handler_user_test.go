package http_server

import (
	"testing"

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

type handlerUserSuite struct {
	suite.Suite
}

func TestInit(t *testing.T) {
	suite.Run(t, new(handlerUserSuite))
}
