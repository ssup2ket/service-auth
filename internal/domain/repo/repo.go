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

// Pkg variables
var cfg *config.Configs
var primaryMySQL *gorm.DB
var secondaryMySQL *gorm.DB

// Init
func Init(c *config.Configs) error {
	var err error

	// Set config
	cfg = c
	gormConfig := &gorm.Config{}
	if cfg.DeployEnv != config.DeployEnvLocal {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	// Connect to primary MySQL
	primaryDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.MySQLPrimaryUser, c.MySQLPrimaryPassword, c.MySQLPrimaryIP, c.MySQLPrimaryPort, c.MySQLDatabase)
	primaryMySQL, err = gorm.Open(mysql.Open(primaryDSN), gormConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to primary MySQL")
		return err
	}

	// Connect to secondary MySQL
	secondaryDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.MySQLSecondaryUser, c.MySQLSecondaryPassword, c.MySQLSecondaryIP, c.MySQLSecondaryPort, c.MySQLDatabase)
	secondaryMySQL, err = gorm.Open(mysql.Open(secondaryDSN), gormConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to secondary MySQL")
		return err
	}

	// Init schemas
	if err = primaryMySQL.AutoMigrate(
		&entity.UserInfo{},
		&entity.UserSecret{},
		&entity.Outbox{},
	); err != nil {
		log.Error().Err(err).Msg("Failed to init schemas")
		return err
	}

	return nil
}

func GetDBConns() (*gorm.DB, *gorm.DB) {
	return primaryMySQL, secondaryMySQL
}

// DB transaction
type DBTx struct {
	tx *gorm.DB
}

func NewDBTx() *DBTx {
	return &DBTx{}
}

func (d *DBTx) getTx() *gorm.DB {
	return d.tx
}

func (d *DBTx) Begin() error {
	d.tx = primaryMySQL.Begin()
	return d.tx.Error
}

func (d *DBTx) Commit() error {
	return d.tx.Commit().Error
}

func (d *DBTx) Rollback() error {
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
