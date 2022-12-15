package db

import (
	"context"
	"database/sql"
	"ehdw/smartiko-test/src/config"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/source/github"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	lg "github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func ConnectPGDatabase(ctx context.Context, withLogger bool) *pgxpool.Pool {
	createDatabase(ctx)
	migrateDb(ctx)
	conf := config.Config().Database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.User, conf.Password,
		conf.Host, conf.Port, conf.Db)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		lg.WithError(err).Fatal("Unable to parse database config")
	}
	if withLogger {
		pglg := lg.New()
		config.ConnConfig.LogLevel = pgx.LogLevelDebug
		config.ConnConfig.Logger = logrusadapter.NewLogger(pglg)
	}
	repo, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		lg.Fatal("Failed to initialise connection to Postgres Database", zap.Error(err))
	}
	return repo
}

func createDatabase(ctx context.Context) {
	conf := config.Config().Database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable",
		conf.User, conf.Password,
		conf.Host, conf.Port)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		lg.WithError(err).Fatal("Unable to parse database config")
	}
	repo, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		lg.Fatal("Failed to initialise connection to Postgres Database", zap.Error(err))
	}
	_, err = repo.Exec(ctx, "CREATE DATABASE \""+conf.Db+"\"")
	if err != nil && !strings.Contains(err.Error(), "(SQLSTATE 42P04)") {
		lg.Fatal("Failed to create a database: " + err.Error())
	}
}

func migrateDb(ctx context.Context) {
	conf := config.Config().Database
	lg.Println("Starting migrations")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.User, conf.Password,
		conf.Host, conf.Port, conf.Db)
	dbCon, err := sql.Open("postgres", dsn)
	if err != nil {
		lg.Fatal("failed to connect to db :: ", err)
	}

	driver, err := postgres.WithInstance(dbCon, &postgres.Config{})
	if err != nil {
		lg.Fatal("failed to create postgres abstraction :: ", err)
	}

	srcUrl := "file://migrations"
	migrator, err := migrate.NewWithDatabaseInstance(srcUrl, conf.Db, driver)
	if err != nil {
		lg.Fatal("failed to create migrator instance :: ", err)
	}

	err = migrator.Up()
	if err != nil && err.Error() != "no change" {
		lg.Fatal("failed to apply migrations :: ", err)
	}
	lg.Println("Migration finished successfully")
}
