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
	modeluuid "github.com/ssup2ket/ssup2ket-auth-service/pkg/model/uuid"
)

const (
	userLoginIDCorrect = "test0000"
	userPasswdCorrect  = "test0000"
	userPhoneCorrect   = "000-0000-0000"
	userEmailCorrect   = "test@test.com"
)

var (
	userIDCorrect = modeluuid.FromStringOrNil("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
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
	u.sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_infos` (`id`,`created_at`,`updated_at`,`deleted_at`,`login_id`,`phone`,`email`) VALUES (?,?,?,?,?,?,?)")).
		WithArgs(userIDCorrect, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), userLoginIDCorrect, userPhoneCorrect, userEmailCorrect).
		WillReturnResult(sqlmock.NewResult(1, 1))
	u.sqlMock.ExpectCommit()

	err := u.repo.Create(context.Background(), &model.UserInfo{
		ID:      userIDCorrect,
		LoginID: userLoginIDCorrect,
		Phone:   userPhoneCorrect,
		Email:   userEmailCorrect,
	})
	require.NoError(u.T(), err)
}

func (u *userInfoSuite) TestGet() {
	u.sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_infos` WHERE id = ? AND `user_infos`.`deleted_at` IS NULL ORDER BY `user_infos`.`id` LIMIT 1")).
		WithArgs(userIDCorrect).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login_id", "phone", "email"}).
			AddRow(userIDCorrect, userLoginIDCorrect, userPhoneCorrect, userEmailCorrect))

	userInfo, err := u.repo.Get(context.Background(), userIDCorrect)
	require.NoError(u.T(), err)
	require.Equal(u.T(), userIDCorrect, userInfo.ID)
	require.Equal(u.T(), userLoginIDCorrect, userInfo.LoginID)
	require.Equal(u.T(), userPhoneCorrect, userInfo.Phone)
	require.Equal(u.T(), userEmailCorrect, userInfo.Email)
}

func (u *userInfoSuite) TestCreateAndGet() {
	u.sqlMock.ExpectBegin()
	u.sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_infos` (`id`,`created_at`,`updated_at`,`deleted_at`,`login_id`,`phone`,`email`) VALUES (?,?,?,?,?,?,?)")).
		WithArgs(userIDCorrect, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), userLoginIDCorrect, userPhoneCorrect, userEmailCorrect).
		WillReturnResult(sqlmock.NewResult(1, 1))
	u.sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_infos` WHERE id = ? AND `user_infos`.`deleted_at` IS NULL ORDER BY `user_infos`.`id` LIMIT 1")).
		WithArgs(userIDCorrect).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login_id", "phone", "email"}).
			AddRow(userIDCorrect, userLoginIDCorrect, userPhoneCorrect, userEmailCorrect))
	u.sqlMock.ExpectCommit()

	tx := NewDBTx()
	tx.Begin()
	err := u.repo.WithTx(tx).Create(context.Background(), &model.UserInfo{
		ID:      userIDCorrect,
		LoginID: userLoginIDCorrect,
		Phone:   userPhoneCorrect,
		Email:   userEmailCorrect,
	})
	require.NoError(u.T(), err)

	userInfo, err := u.repo.WithTx(tx).Get(context.Background(), userIDCorrect)
	require.NoError(u.T(), err)
	require.Equal(u.T(), userIDCorrect, userInfo.ID)
	require.Equal(u.T(), userLoginIDCorrect, userInfo.LoginID)
	require.Equal(u.T(), userPhoneCorrect, userInfo.Phone)
	require.Equal(u.T(), userEmailCorrect, userInfo.Email)
	tx.Commit()
}
