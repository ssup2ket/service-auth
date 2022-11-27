package repo

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ssup2ket/service-auth/internal/domain/entity"
	"github.com/ssup2ket/service-auth/internal/test"
	"github.com/ssup2ket/service-auth/pkg/auth/hashing"
	"github.com/ssup2ket/service-auth/pkg/auth/token"
)

func TestUserSecret(t *testing.T) {
	suite.Run(t, new(userSecretSuite))
}

type userSecretSuite struct {
	suite.Suite
	sqlMock sqlmock.Sqlmock

	tx   *DBTxImp
	repo UserSecretRepo

	passwdHash       []byte
	passwdSalt       []byte
	refreshTokenHash []byte
	refreshTokenSalt []byte
}

func (u *userSecretSuite) SetupTest() {
	var err error
	var db *sql.DB

	// Init sqlMock
	db, u.sqlMock, err = sqlmock.New()
	require.NoError(u.T(), err)

	// Init DB
	primaryMySQL, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}))
	require.NoError(u.T(), err)

	// Init transaction, repo
	u.tx = NewDBTxImp(primaryMySQL)
	u.repo = NewUserSecretRepoImp(primaryMySQL)

	// Get password and refresh token's hash and salt
	tokenInfo, _ := token.CreateRefreshToken(&token.AuthClaims{UserID: test.UserLoginIDCorrect, UserLoginID: test.UserLoginIDCorrect})
	u.passwdHash, u.passwdSalt, _ = hashing.GetStrHashAndSalt(test.UserPasswdCorrect)
	u.refreshTokenHash, u.refreshTokenSalt, _ = hashing.GetStrHashAndSalt(tokenInfo.Token)
}

func (u *userSecretSuite) AfterTest(_, _ string) {
	require.NoError(u.T(), u.sqlMock.ExpectationsWereMet())
}

func (u *userSecretSuite) TestCreateSuccess() {
	u.sqlMock.ExpectBegin()
	u.sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_secrets` (`id`,`created_at`,`updated_at`,`deleted_at`,`passwd_hash`,`passwd_salt`,`refresh_token_hash`,`refresh_token_salt`) VALUES (?,?,?,?,?,?,?,?)")).
		WithArgs(test.UserIDCorrect, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), u.passwdHash, u.passwdSalt, u.refreshTokenHash, u.refreshTokenSalt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	u.sqlMock.ExpectCommit()

	err := u.repo.Create(context.Background(), &entity.UserSecret{
		ID:               test.UserIDCorrect,
		PasswdHash:       u.passwdHash,
		PasswdSalt:       u.passwdSalt,
		RefreshTokenHash: u.refreshTokenHash,
		RefreshTokenSalt: u.refreshTokenSalt,
	})
	require.NoError(u.T(), err)
}

func (u *userSecretSuite) TestCreateError() {
	u.sqlMock.ExpectBegin()
	u.sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_secrets` (`id`,`created_at`,`updated_at`,`deleted_at`,`passwd_hash`,`passwd_salt`,`refresh_token_hash`,`refresh_token_salt`) VALUES (?,?,?,?,?,?,?,?)")).
		WithArgs(test.UserIDCorrect, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), u.passwdHash, u.passwdSalt, u.refreshTokenHash, u.refreshTokenSalt).
		WillReturnError(fmt.Errorf("error"))
	u.sqlMock.ExpectRollback()

	err := u.repo.Create(context.Background(), &entity.UserSecret{
		ID:               test.UserIDCorrect,
		PasswdHash:       u.passwdHash,
		PasswdSalt:       u.passwdSalt,
		RefreshTokenHash: u.refreshTokenHash,
		RefreshTokenSalt: u.refreshTokenSalt,
	})
	require.Error(u.T(), err)
}

func (u *userSecretSuite) TestGetSuccess() {
	u.sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_secrets` WHERE id = ? AND `user_secrets`.`deleted_at` IS NULL ORDER BY `user_secrets`.`id` LIMIT 1")).
		WithArgs(test.UserIDCorrect).
		WillReturnRows(sqlmock.NewRows([]string{"id", "passwd_hash", "passwd_salt", "refresh_token_hash", "refresh_token_salt"}).
			AddRow(test.UserIDCorrect, u.passwdHash, u.passwdSalt, u.refreshTokenHash, u.refreshTokenSalt))

	userSecret, err := u.repo.Get(context.Background(), test.UserIDCorrect)
	require.NoError(u.T(), err)
	require.Equal(u.T(), test.UserIDCorrect, userSecret.ID)
	require.Equal(u.T(), u.passwdHash, userSecret.PasswdHash)
	require.Equal(u.T(), u.passwdSalt, userSecret.PasswdSalt)
	require.Equal(u.T(), u.refreshTokenHash, userSecret.RefreshTokenHash)
	require.Equal(u.T(), u.refreshTokenSalt, userSecret.RefreshTokenSalt)
}

func (u *userSecretSuite) TestGetError() {
	u.sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_secrets` WHERE id = ? AND `user_secrets`.`deleted_at` IS NULL ORDER BY `user_secrets`.`id` LIMIT 1")).
		WithArgs(test.UserIDCorrect).
		WillReturnError(fmt.Errorf("error"))

	_, err := u.repo.Get(context.Background(), test.UserIDCorrect)
	require.Error(u.T(), err)
}

func (u *userSecretSuite) TestCreateAndGetWithTxSuccess() {
	u.sqlMock.ExpectBegin()
	u.sqlMock.ExpectExec(regexp.QuoteMeta("INSERT INTO `user_secrets` (`id`,`created_at`,`updated_at`,`deleted_at`,`passwd_hash`,`passwd_salt`,`refresh_token_hash`,`refresh_token_salt`) VALUES (?,?,?,?,?,?,?,?)")).
		WithArgs(test.UserIDCorrect, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), u.passwdHash, u.passwdSalt, u.refreshTokenHash, u.refreshTokenSalt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	u.sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `user_secrets` WHERE id = ? AND `user_secrets`.`deleted_at` IS NULL ORDER BY `user_secrets`.`id` LIMIT 1")).
		WithArgs(test.UserIDCorrect).
		WillReturnRows(sqlmock.NewRows([]string{"id", "passwd_hash", "passwd_salt", "refresh_token_hash", "refresh_token_salt"}).
			AddRow(test.UserIDCorrect, u.passwdHash, u.passwdSalt, u.refreshTokenHash, u.refreshTokenSalt))
	u.sqlMock.ExpectCommit()

	tx, _ := u.tx.Begin()
	err := u.repo.WithTx(tx).Create(context.Background(), &entity.UserSecret{
		ID:               test.UserIDCorrect,
		PasswdHash:       u.passwdHash,
		PasswdSalt:       u.passwdSalt,
		RefreshTokenHash: u.refreshTokenHash,
		RefreshTokenSalt: u.refreshTokenSalt,
	})
	require.NoError(u.T(), err)

	userSecret, err := u.repo.WithTx(tx).Get(context.Background(), test.UserIDCorrect)
	require.NoError(u.T(), err)
	require.Equal(u.T(), test.UserIDCorrect, userSecret.ID)
	require.Equal(u.T(), u.passwdHash, userSecret.PasswdHash)
	require.Equal(u.T(), u.passwdSalt, userSecret.PasswdSalt)
	require.Equal(u.T(), u.refreshTokenHash, userSecret.RefreshTokenHash)
	require.Equal(u.T(), u.refreshTokenSalt, userSecret.RefreshTokenSalt)
	tx.Commit()
}

func (u *userSecretSuite) TestUpdateSuccess() {
	u.sqlMock.ExpectBegin()
	u.sqlMock.ExpectExec(regexp.QuoteMeta("UPDATE `user_secrets` SET `updated_at`=?,`passwd_hash`=?,`passwd_salt`=?,`refresh_token_hash`=?,`refresh_token_salt`=? WHERE `id` = ?")).
		WithArgs(sqlmock.AnyArg(), u.passwdHash, u.passwdSalt, u.refreshTokenHash, u.refreshTokenSalt, test.UserIDCorrect).
		WillReturnResult(sqlmock.NewResult(1, 1))
	u.sqlMock.ExpectCommit()

	err := u.repo.Update(context.Background(), &entity.UserSecret{
		ID:               test.UserIDCorrect,
		PasswdHash:       u.passwdHash,
		PasswdSalt:       u.passwdSalt,
		RefreshTokenHash: u.refreshTokenHash,
		RefreshTokenSalt: u.refreshTokenSalt,
	})
	require.NoError(u.T(), err)
}

func (u *userSecretSuite) TestUpdateError() {
	u.sqlMock.ExpectBegin()
	u.sqlMock.ExpectExec(regexp.QuoteMeta("UPDATE `user_secrets` SET `updated_at`=?,`passwd_hash`=?,`passwd_salt`=?,`refresh_token_hash`=?,`refresh_token_salt`=? WHERE `id` = ?")).
		WithArgs(sqlmock.AnyArg(), u.passwdHash, u.passwdSalt, u.refreshTokenHash, u.refreshTokenSalt, test.UserIDCorrect).
		WillReturnError(fmt.Errorf("error"))
	u.sqlMock.ExpectRollback()

	err := u.repo.Update(context.Background(), &entity.UserSecret{
		ID:               test.UserIDCorrect,
		PasswdHash:       u.passwdHash,
		PasswdSalt:       u.passwdSalt,
		RefreshTokenHash: u.refreshTokenHash,
		RefreshTokenSalt: u.refreshTokenSalt,
	})
	require.Error(u.T(), err)
}

func (u *userSecretSuite) TestDeleteSuccess() {
	u.sqlMock.ExpectBegin()
	u.sqlMock.ExpectExec(regexp.QuoteMeta("UPDATE `user_secrets` SET `deleted_at`=? WHERE id = ? AND `user_secrets`.`deleted_at` IS NULL")).
		WithArgs(sqlmock.AnyArg(), test.UserIDCorrect).
		WillReturnResult(sqlmock.NewResult(1, 1))
	u.sqlMock.ExpectCommit()

	err := u.repo.Delete(context.Background(), test.UserIDCorrect)
	require.NoError(u.T(), err)
}

func (u *userSecretSuite) TestDeleteError() {
	u.sqlMock.ExpectBegin()
	u.sqlMock.ExpectExec(regexp.QuoteMeta("UPDATE `user_secrets` SET `deleted_at`=? WHERE id = ? AND `user_secrets`.`deleted_at` IS NULL")).
		WithArgs(sqlmock.AnyArg(), test.UserIDCorrect).
		WillReturnError(fmt.Errorf("error"))
	u.sqlMock.ExpectRollback()

	err := u.repo.Delete(context.Background(), test.UserIDCorrect)
	require.Error(u.T(), err)
}
