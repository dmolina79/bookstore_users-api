package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	LogLevel      = "LOG_LEVEL"
	goEnvironment = "GO_ENVIRONMENT"
	production    = "production"

	// mysql DB env vars
	mysqlUsersUsername = "MYSQL_DB_USERNAME"
	mysqlUsersPwd      = "MYSQL_DB_PASSWORD"
	mysqlUsersHost     = "MYSQL_DB_HOST"
	mysqlUsersDb       = "MYSQL_DB_NAME"
)

var (
	logLevel string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}

	logLevel = os.Getenv(LogLevel)
}

func GetUsersDBDataSourceName() string {
	username := os.Getenv(mysqlUsersUsername)
	password := os.Getenv(mysqlUsersPwd)
	host := os.Getenv(mysqlUsersHost)
	dbName := os.Getenv(mysqlUsersDb)
	datasourceName :=  fmt.Sprintf("%s:%s@tcp(%s)/%s?charset?utf8",
		username,
		password,
		host,
		dbName,
	)

	return datasourceName
}

func GetLogLevel() string {
	return logLevel
}

func IsProduction() bool {
	return os.Getenv(goEnvironment) == production
}
