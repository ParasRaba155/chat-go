package main

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"

	"app/config"
	"app/logger"
)

var migrationPath = "migrations"

//go:embed migrations/*.sql
var fs embed.FS

func main() {
	appConfig, err := config.LoadConfig(".env.dev")
	if err != nil {
		panic("could not load the config ")
	}

	logger := logger.New(appConfig.LogFileLoc)
	dbConfig := appConfig.Database

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DefaultDB, dbConfig.SSLMODE)

	driver, err := iofs.New(fs, migrationPath)
	if err != nil {
		logger.Fatal("could not create migration driver", zap.Error(err))
	}

	m, err := migrate.NewWithSourceInstance("iofs", driver, dbURL)
	if err != nil {
		logger.Fatal("could not create migration source instance", zap.Error(err))
	}

	logger.Sugar().Infof("m = %+v", m)

	err = m.Up()
	if err != nil {
		logger.Fatal("could not migrate up", zap.Error(err))
	}
}
