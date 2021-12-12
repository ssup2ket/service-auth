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
	EnvMySQLDatabase = "MYSQL_DATABASE"

	EnvMySQLPrimaryIP       = "MYSQL_PRIMARY_IP"
	EnvMySQLPrimaryPort     = "MYSQL_PRIMARY_PORT"
	EnvMySQLPrimaryUser     = "MYSQL_PRIMARY_USER"
	EnvMySQLPrimaryPassword = "MYSQL_PRIMARY_PASSWORD"

	EnvMySQLSecondaryIP       = "MYSQL_SECONDARY_IP"
	EnvMySQLSecondaryPort     = "MYSQL_SECONDARY_PORT"
	EnvMySQLSecondaryUser     = "MYSQL_SECONDARY_USER"
	EnvMySQLSecondaryPassword = "MYSQL_SECONDARY_PASSWORD"

	// Jaeger
	EnvJaegerCollectorEndpoint = "JAEGER_COLLECTOR_ENDPOINT"
)

type Configs struct {
	// Server URL
	ServerURL string

	// Deploy env
	DeployEnv DeployEnv

	// MySQL
	MySQLDatabase string

	// Primary
	MySQLPrimaryIP       string
	MySQLPrimaryPort     string
	MySQLPrimaryUser     string
	MySQLPrimaryPassword string

	// Secondary
	MySQLSecondaryIP       string
	MySQLSecondaryPort     string
	MySQLSecondaryUser     string
	MySQLSecondaryPassword string

	// Jaeger
	JaegerCollectorEndpoint string
}

func GetConfigs() *Configs {
	return &Configs{
		ServerURL: os.Getenv(EnvServerURL),

		DeployEnv: DeployEnv(os.Getenv(EnvDeployEnv)),

		MySQLDatabase: os.Getenv(EnvMySQLDatabase),

		MySQLPrimaryIP:       os.Getenv(EnvMySQLPrimaryIP),
		MySQLPrimaryPort:     os.Getenv(EnvMySQLPrimaryPort),
		MySQLPrimaryUser:     os.Getenv(EnvMySQLPrimaryUser),
		MySQLPrimaryPassword: os.Getenv(EnvMySQLPrimaryPassword),

		MySQLSecondaryIP:       os.Getenv(EnvMySQLSecondaryIP),
		MySQLSecondaryPort:     os.Getenv(EnvMySQLSecondaryPort),
		MySQLSecondaryUser:     os.Getenv(EnvMySQLSecondaryUser),
		MySQLSecondaryPassword: os.Getenv(EnvMySQLSecondaryPassword),

		JaegerCollectorEndpoint: os.Getenv(EnvJaegerCollectorEndpoint),
	}
}

// Deploy env
type DeployEnv string

const (
	DeployEnvLocal DeployEnv = "local"
	DeployEnvDev   DeployEnv = "dev"
	DeployEnvStage DeployEnv = "stage"
	DeployEnvProd  DeployEnv = "prod"
)
