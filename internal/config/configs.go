package config

import (
	"os"
)

// Config
const (
	// Server URL
	EnvServerURL = "SERVER_URL"

	// Deploy env
	EnvDeployEnv = "DEPLOY_ENV"

	// MySQL
	EnvMySQLPrimaryIP       = "MYSQL_PRIMARY_IP"
	EnvMySQLPrimaryPort     = "MYSQL_PRIMARY_PORT"
	EnvMySQLPrimaryUser     = "MYSQL_PRIMARY_USER"
	EnvMySQLPrimaryPassword = "MYSQL_PRIMARY_PASSWORD"
	EnvMySQLPrimaryDatabase = "MYSQL_PRIMARY_DATABASE"

	EnvMySQLSecondaryIP       = "MYSQL_SECONDARY_IP"
	EnvMySQLSecondaryPort     = "MYSQL_SECONDARY_PORT"
	EnvMySQLSecondaryUser     = "MYSQL_SECONDARY_USER"
	EnvMySQLSecondaryPassword = "MYSQL_SECONDARY_PASSWORD"
	EnvMySQLSeconaaryDatabase = "MYSQL_SECONDARY_DATABASE"
)

type Configs struct {
	// Server URL
	ServerURL string

	// Deploy env
	DeployEnv DeployEnv

	// MySQL
	// Primary
	MySQLPrimaryIP       string
	MySQLPrimaryPort     string
	MySQLPrimaryUser     string
	MySQLPrimaryPassword string
	MySQLPrimaryDatabase string

	// Secondary
	MySQLSecondaryIP       string
	MySQLSecondaryPort     string
	MySQLSecondaryUser     string
	MySQLSecondaryPassword string
	MySQLSecondaryDatabase string
}

func GetConfigs() *Configs {
	return &Configs{
		ServerURL: os.Getenv(EnvServerURL),

		DeployEnv: DeployEnv(os.Getenv(EnvDeployEnv)),

		MySQLPrimaryIP:       os.Getenv(EnvMySQLPrimaryIP),
		MySQLPrimaryPort:     os.Getenv(EnvMySQLPrimaryPort),
		MySQLPrimaryUser:     os.Getenv(EnvMySQLPrimaryUser),
		MySQLPrimaryPassword: os.Getenv(EnvMySQLPrimaryPassword),
		MySQLPrimaryDatabase: os.Getenv(EnvMySQLPrimaryDatabase),

		MySQLSecondaryIP:       os.Getenv(EnvMySQLPrimaryIP),
		MySQLSecondaryPort:     os.Getenv(EnvMySQLPrimaryPort),
		MySQLSecondaryUser:     os.Getenv(EnvMySQLPrimaryUser),
		MySQLSecondaryPassword: os.Getenv(EnvMySQLPrimaryPassword),
		MySQLSecondaryDatabase: os.Getenv(EnvMySQLPrimaryDatabase),
	}
}

// Deploy env
type DeployEnv string

const (
	DeployEnvLocal DeployEnv = "local"
	DeployEnvDev   DeployEnv = "dev"
	DeployEnvProd  DeployEnv = "prod"
)
