package test

import (
	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/model"
	modeluuid "github.com/ssup2ket/ssup2ket-auth-service/pkg/model/uuid"
)

const (
	UserLoginIDCorrect = "test0000"
	UserRoleCorrect    = model.UserRoleAdmin
	UserPasswdCorrect  = "test0000"
	UserPhoneCorrect   = "000-0000-0000"
	UserEmailCorrect   = "test@test.com"

	UserIDWrongFormat    = "aaaa-aaaa"
	UserLoginIDShort     = "test0"
	UserLoginIDLong      = "testtesttesttesttesttest"
	UserPasswdShort      = "test0"
	UserPasswdLong       = "testtesttesttesttesttest"
	UserRoleWrong        = "tester"
	UserPhoneWrongFormat = "00-000-00000"
	UserEmailWrongFormat = "testtest.com"
)

var (
	UserIDCorrect = modeluuid.FromStringOrNil("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
)
