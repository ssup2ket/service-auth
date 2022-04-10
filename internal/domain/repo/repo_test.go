package repo

import (
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestRepo(t *testing.T) {
	suite.Run(t, new(repoSuite))
}

type repoSuite struct {
	suite.Suite
}

func (r *repoSuite) TestGetRetrunError() {
	errNotFound := getReturnErr(gorm.ErrRecordNotFound)
	require.Equal(r.T(), ErrNotFound, errNotFound)

	errConflict := getReturnErr(&gomysql.MySQLError{Number: 1062})
	require.Equal(r.T(), ErrConflict, errConflict)

	errNullError := getReturnErr(nil)
	require.Equal(r.T(), ErrServerError, errNullError)
}
