package repo

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/model"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/test"
)

func TestInit(t *testing.T) {
	suite.Run(t, new(userInfoSuite))
}

type userInfoSuite struct {
	suite.Suite
	sqlMock sqlmock.Sqlmock

	repo UserInfoRepo
}

func (u *userInfoSuite) SetupTest() {
	var err error
	var db *sql.DB

	// Init sqlMock
	db, u.sqlMock, err = sqlmock.New()
	require.NoError(u.T(), err)

	// Init DB
	primaryMySQL, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}))
	require.NoError(u.T(), err)

	// Init repo
	u.repo = NewUserInfoRepoImp(primaryMySQL)
}

func (u *userInfoSuite) AfterTest(_, _ string) {
	require.NoError(u.T(), u.sqlMock.ExpectationsWereMet())
}

func (u *userInfoSuite) TestCreate() {
	u.sqlMock.ExpectBegin()
	u.sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_infos` (`id`,`created_at`,`updated_at`,`deleted_at`,`login_id`,`role`,`phone`,`email`) VALUES (?,?,?,?,?,?,?,?)")).
		WithArgs(test.UserIDCorrect, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), test.UserLoginIDCorrect, test.UserRoleCorrect, test.UserPhoneCorrect, test.UserEmailCorrect).
		WillReturnResult(sqlmock.NewResult(1, 1))
	u.sqlMock.ExpectCommit()

	err := u.repo.Create(context.Background(), &model.UserInfo{
		ID:      test.UserIDCorrect,
		LoginID: test.UserLoginIDCorrect,
		Role:    test.UserRoleCorrect,
		Phone:   test.UserPhoneCorrect,
		Email:   test.UserEmailCorrect,
	})
	require.NoError(u.T(), err)
}

func (u *userInfoSuite) TestGet() {
	u.sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_infos` WHERE id = ? AND `user_infos`.`deleted_at` IS NULL ORDER BY `user_infos`.`id` LIMIT 1")).
		WithArgs(test.UserIDCorrect).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login_id", "role", "phone", "email"}).
			AddRow(test.UserIDCorrect, test.UserLoginIDCorrect, test.UserRoleCorrect, test.UserPhoneCorrect, test.UserEmailCorrect))

	userInfo, err := u.repo.Get(context.Background(), test.UserIDCorrect)
	require.NoError(u.T(), err)
	require.Equal(u.T(), test.UserIDCorrect, userInfo.ID)
	require.Equal(u.T(), test.UserLoginIDCorrect, userInfo.LoginID)
	require.Equal(u.T(), test.UserRoleCorrect, userInfo.Role)
	require.Equal(u.T(), test.UserPhoneCorrect, userInfo.Phone)
	require.Equal(u.T(), test.UserEmailCorrect, userInfo.Email)
}

func (u *userInfoSuite) TestCreateAndGet() {
	u.sqlMock.ExpectBegin()
	u.sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_infos` (`id`,`created_at`,`updated_at`,`deleted_at`,`login_id`,`role`,`phone`,`email`) VALUES (?,?,?,?,?,?,?,?)")).
		WithArgs(test.UserIDCorrect, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), test.UserLoginIDCorrect, test.UserRoleCorrect, test.UserPhoneCorrect, test.UserEmailCorrect).
		WillReturnResult(sqlmock.NewResult(1, 1))
	u.sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_infos` WHERE id = ? AND `user_infos`.`deleted_at` IS NULL ORDER BY `user_infos`.`id` LIMIT 1")).
		WithArgs(test.UserIDCorrect).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login_id", "role", "phone", "email"}).
			AddRow(test.UserIDCorrect, test.UserLoginIDCorrect, test.UserRoleCorrect, test.UserPhoneCorrect, test.UserEmailCorrect))
	u.sqlMock.ExpectCommit()

	tx := NewDBTx()
	tx.Begin()
	err := u.repo.WithTx(tx).Create(context.Background(), &model.UserInfo{
		ID:      test.UserIDCorrect,
		LoginID: test.UserLoginIDCorrect,
		Role:    test.UserRoleCorrect,
		Phone:   test.UserPhoneCorrect,
		Email:   test.UserEmailCorrect,
	})
	require.NoError(u.T(), err)

	userInfo, err := u.repo.WithTx(tx).Get(context.Background(), test.UserIDCorrect)
	require.NoError(u.T(), err)
	require.Equal(u.T(), test.UserIDCorrect, userInfo.ID)
	require.Equal(u.T(), test.UserLoginIDCorrect, userInfo.LoginID)
	require.Equal(u.T(), test.UserRoleCorrect, userInfo.Role)
	require.Equal(u.T(), test.UserPhoneCorrect, userInfo.Phone)
	require.Equal(u.T(), test.UserEmailCorrect, userInfo.Email)
	tx.Commit()
}
