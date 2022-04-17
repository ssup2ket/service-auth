package repo

import (
	"errors"
	"fmt"

	gomysql "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ssup2ket/ssup2ket-auth-service/internal/config"
	"github.com/ssup2ket/ssup2ket-auth-service/internal/domain/entity"
)

// Init
func New(c *config.Configs) (DBTx, *gorm.DB, *gorm.DB, error) {
	var err error

	// Set config
	gormConfig := &gorm.Config{}
	if c.DeployEnv != config.DeployEnvLocal {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	// Connect to primary MySQL
	primaryDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.MySQLPrimaryUser, c.MySQLPrimaryPassword, c.MySQLPrimaryIP, c.MySQLPrimaryPort, c.MySQLDatabase)
	primaryMySQL, err := gorm.Open(mysql.Open(primaryDSN), gormConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to primary MySQL")
		return nil, nil, nil, err
	}

	// Connect to secondary MySQL
	secondaryDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.MySQLSecondaryUser, c.MySQLSecondaryPassword, c.MySQLSecondaryIP, c.MySQLSecondaryPort, c.MySQLDatabase)
	secondaryMySQL, err := gorm.Open(mysql.Open(secondaryDSN), gormConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to secondary MySQL")
		return nil, nil, nil, err
	}

	// Init schemas
	if err = primaryMySQL.AutoMigrate(
		&entity.UserInfo{},
		&entity.UserSecret{},
		&entity.Outbox{},
	); err != nil {
		log.Error().Err(err).Msg("Failed to init schemas")
		return nil, nil, nil, err
	}

	return NewDBTxImp(primaryMySQL), primaryMySQL, secondaryMySQL, nil
}

// DB transaction
type DBTx interface {
	GetTx() *gorm.DB

	Begin() (DBTx, error)
	Commit() error
	Rollback() error
}

type DBTxImp struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewDBTxImp(d *gorm.DB) *DBTxImp {
	return &DBTxImp{
		db: d,
	}
}

func (d *DBTxImp) GetTx() *gorm.DB {
	return d.tx
}

func (d *DBTxImp) Begin() (DBTx, error) {
	t := d.db.Begin()
	return &DBTxImp{
		tx: t,
	}, t.Error
}

func (d *DBTxImp) Commit() error {
	return d.tx.Commit().Error
}

func (d *DBTxImp) Rollback() error {
	return d.tx.Rollback().Error
}

// Error
var (
	ErrNotFound    error = fmt.Errorf("resource not found")
	ErrConflict    error = fmt.Errorf("conflict")
	ErrServerError error = fmt.Errorf("server error")
)

func getReturnErr(err error) error {
	var mysqlErr *gomysql.MySQLError
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	} else if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return ErrConflict
	}
	return ErrServerError
}
