package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/repo"
)

func TestService(t *testing.T) {
	suite.Run(t, new(serviceSuite))
}

type serviceSuite struct {
	suite.Suite
}

func (s *serviceSuite) TestGetRetrunError() {
	errNotFound := getReturnErr(repo.ErrNotFound)
	require.Equal(s.T(), ErrRepoNotFound, errNotFound)

	errConflict := getReturnErr(repo.ErrConflict)
	require.Equal(s.T(), ErrRepoConflict, errConflict)

	errServerError := getReturnErr(repo.ErrServerError)
	require.Equal(s.T(), ErrRepoServerError, errServerError)

	errTestError := getReturnErr(nil)
	require.Equal(s.T(), ErrServerErr, errTestError)
}
