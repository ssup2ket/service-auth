package test

import (
	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/entity"
	"github.com/ssup2ket/ssup2ket-auth-service/pkg/entity/uuid"
)

const (
	UserLoginIDCorrect  = "test0000"
	UserLoginIDCorrect2 = "test1111"
	UserRoleCorrect     = entity.UserRoleAdmin
	UserPasswdCorrect   = "test0000"
	UserPhoneCorrect    = "000-0000-0000"
	UserEmailCorrect    = "test@test.com"

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
	UserIDCorrect  = uuid.FromStringOrNil("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	UserIDCorrect2 = uuid.FromStringOrNil("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
)
