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

	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/entity"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/test"
)

func TestUserInfo(t *testing.T) {
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

func (u *userInfoSuite) TestList() {
	u.sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_infos` WHERE `user_infos`.`deleted_at` IS NULL LIMIT 10")).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "login_id", "role", "phone", "email"}).
				AddRow(test.UserIDCorrect, test.UserLoginIDCorrect, test.UserRoleCorrect, test.UserPhoneCorrect, test.UserEmailCorrect).
				AddRow(test.UserIDCorrect2, test.UserLoginIDCorrect2, test.UserRoleCorrect, test.UserPhoneCorrect, test.UserEmailCorrect),
		)

	userInfos, err := u.repo.List(context.Background(), 0, 10)
	require.NoError(u.T(), err)
	require.Equal(u.T(), test.UserIDCorrect, userInfos[0].ID)
	require.Equal(u.T(), test.UserLoginIDCorrect, userInfos[0].LoginID)
	require.Equal(u.T(), test.UserRoleCorrect, userInfos[0].Role)
	require.Equal(u.T(), test.UserPhoneCorrect, userInfos[0].Phone)
	require.Equal(u.T(), test.UserEmailCorrect, userInfos[0].Email)
	require.Equal(u.T(), test.UserIDCorrect2, userInfos[1].ID)
	require.Equal(u.T(), test.UserLoginIDCorrect2, userInfos[1].LoginID)
	require.Equal(u.T(), test.UserRoleCorrect, userInfos[1].Role)
	require.Equal(u.T(), test.UserPhoneCorrect, userInfos[1].Phone)
	require.Equal(u.T(), test.UserEmailCorrect, userInfos[1].Email)
}

func (u *userInfoSuite) TestCreate() {
	u.sqlMock.ExpectBegin()
	u.sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_infos` (`id`,`created_at`,`updated_at`,`deleted_at`,`login_id`,`role`,`phone`,`email`) VALUES (?,?,?,?,?,?,?,?)")).
		WithArgs(test.UserIDCorrect, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), test.UserLoginIDCorrect, test.UserRoleCorrect, test.UserPhoneCorrect, test.UserEmailCorrect).
		WillReturnResult(sqlmock.NewResult(1, 1))
	u.sqlMock.ExpectCommit()

	err := u.repo.Create(context.Background(), &entity.UserInfo{
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

func (u *userInfoSuite) TestGetByLoginID() {
	u.sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_infos` WHERE login_id = ? AND `user_infos`.`deleted_at` IS NULL ORDER BY `user_infos`.`id` LIMIT 1")).
		WithArgs(test.UserLoginIDCorrect).
		WillReturnRows(sqlmock.NewRows([]string{"id", "login_id", "role", "phone", "email"}).
			AddRow(test.UserIDCorrect, test.UserLoginIDCorrect, test.UserRoleCorrect, test.UserPhoneCorrect, test.UserEmailCorrect))

	userInfo, err := u.repo.GetByLoginID(context.Background(), test.UserLoginIDCorrect)
	require.NoError(u.T(), err)
	require.Equal(u.T(), test.UserIDCorrect, userInfo.ID)
	require.Equal(u.T(), test.UserLoginIDCorrect, userInfo.LoginID)
	require.Equal(u.T(), test.UserRoleCorrect, userInfo.Role)
	require.Equal(u.T(), test.UserPhoneCorrect, userInfo.Phone)
	require.Equal(u.T(), test.UserEmailCorrect, userInfo.Email)
}

func (u *userInfoSuite) TestCreateAndGetWithTx() {
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
	err := u.repo.WithTx(tx).Create(context.Background(), &entity.UserInfo{
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

func (u *userInfoSuite) TestUpdate() {
	u.sqlMock.ExpectBegin()
	u.sqlMock.ExpectExec(regexp.QuoteMeta("UPDATE `user_infos` SET `updated_at`=?,`login_id`=?,`role`=?,`phone`=?,`email`=? WHERE `id` = ?")).
		WithArgs(sqlmock.AnyArg(), test.UserLoginIDCorrect, test.UserRoleCorrect, test.UserPhoneCorrect, test.UserEmailCorrect, test.UserIDCorrect).
		WillReturnResult(sqlmock.NewResult(1, 1))
	u.sqlMock.ExpectCommit()

	err := u.repo.Update(context.Background(), &entity.UserInfo{
		ID:      test.UserIDCorrect,
		LoginID: test.UserLoginIDCorrect,
		Role:    test.UserRoleCorrect,
		Phone:   test.UserPhoneCorrect,
		Email:   test.UserEmailCorrect,
	})
	require.NoError(u.T(), err)
}

func (u *userInfoSuite) TestDelete() {
	u.sqlMock.ExpectBegin()
	u.sqlMock.ExpectExec(regexp.QuoteMeta("UPDATE `user_infos` SET `deleted_at`=? WHERE id = ? AND `user_infos`.`deleted_at` IS NULL")).
		WithArgs(sqlmock.AnyArg(), test.UserIDCorrect).
		WillReturnResult(sqlmock.NewResult(1, 1))
	u.sqlMock.ExpectCommit()

	err := u.repo.Delete(context.Background(), test.UserIDCorrect)
	require.NoError(u.T(), err)
}
