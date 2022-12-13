package db

import (
	"context"
	"ehdw/smartiko-test/src/config"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	lg "github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func ConnectPGDatabase(ctx context.Context, withLogger bool) *pgxpool.Pool {
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
