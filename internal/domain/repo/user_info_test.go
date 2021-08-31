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
	"github.com/ssup2ket/ssup2ket-auth-service/pkg/uuid"
)

const (
	userIDCorrect     = "test0000"
	userPasswdCorrect = "test0000"
	userPhoneCorrect  = "000-0000-0000"
	userEmailCorrect  = "test@test.com"
)

var (
	userUUIDCorrect = uuid.FromStringOrNil("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
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
	u.sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_infos` (`uuid`,`created_at`,`updated_at`,`deleted_at`,`id`,`phone`,`email`) VALUES (?,?,?,?,?,?,?)")).
		WithArgs(userUUIDCorrect, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), userIDCorrect, userPhoneCorrect, userEmailCorrect).
		WillReturnResult(sqlmock.NewResult(1, 1))
	u.sqlMock.ExpectCommit()

	err := u.repo.Create(context.Background(), &model.UserInfo{
		UUID:  userUUIDCorrect,
		ID:    userIDCorrect,
		Phone: userPhoneCorrect,
		Email: userEmailCorrect,
	})
	require.NoError(u.T(), err)
}

func (u *userInfoSuite) TestGetPrimary() {
	u.sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_infos` WHERE uuid = ? AND `user_infos`.`deleted_at` IS NULL ORDER BY `user_infos`.`uuid` LIMIT 1")).
		WithArgs(userUUIDCorrect).
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "id", "phone", "email"}).
			AddRow(userUUIDCorrect, userIDCorrect, userPhoneCorrect, userEmailCorrect))

	userInfo, err := u.repo.Get(context.Background(), userUUIDCorrect)
	require.NoError(u.T(), err)
	require.Equal(u.T(), userUUIDCorrect, userInfo.UUID)
	require.Equal(u.T(), userIDCorrect, userInfo.ID)
	require.Equal(u.T(), userPhoneCorrect, userInfo.Phone)
	require.Equal(u.T(), userEmailCorrect, userInfo.Email)
}

func (u *userInfoSuite) TestCreateAndGetPrimary() {
	u.sqlMock.ExpectBegin()
	u.sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_infos` (`uuid`,`created_at`,`updated_at`,`deleted_at`,`id`,`phone`,`email`) VALUES (?,?,?,?,?,?,?)")).
		WithArgs(userUUIDCorrect, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), userIDCorrect, userPhoneCorrect, userEmailCorrect).
		WillReturnResult(sqlmock.NewResult(1, 1))
	u.sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_infos` WHERE uuid = ? AND `user_infos`.`deleted_at` IS NULL ORDER BY `user_infos`.`uuid` LIMIT 1")).
		WithArgs(userUUIDCorrect).
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "id", "phone", "email"}).
			AddRow(userUUIDCorrect, userIDCorrect, userPhoneCorrect, userEmailCorrect))
	u.sqlMock.ExpectCommit()

	tx := NewDBTx()
	tx.Begin()
	err := u.repo.WithTx(tx).Create(context.Background(), &model.UserInfo{
		UUID:  userUUIDCorrect,
		ID:    userIDCorrect,
		Phone: userPhoneCorrect,
		Email: userEmailCorrect,
	})
	require.NoError(u.T(), err)

	userInfo, err := u.repo.WithTx(tx).Get(context.Background(), userUUIDCorrect)
	require.NoError(u.T(), err)
	require.Equal(u.T(), userUUIDCorrect, userInfo.UUID)
	require.Equal(u.T(), userIDCorrect, userInfo.ID)
	require.Equal(u.T(), userPhoneCorrect, userInfo.Phone)
	require.Equal(u.T(), userEmailCorrect, userInfo.Email)
	tx.Commit()
}
