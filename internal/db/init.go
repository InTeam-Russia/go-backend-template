package db

import (
	"context"
	"os"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func InitDb(dbUrl string, logger *zap.Logger) (*pgxpool.Pool, error) {
	pool, err := CreatePool(dbUrl, logger)

	createTableSql, err := os.ReadFile("./db/CreateTables.sql")
	if err != nil {
		return nil, err
	}

	_, err = pool.Exec(context.Background(), string(createTableSql))
	if err != nil {
		return nil, err
	}

	return pool, err
}

func DropDb(dbUrl string, logger *zap.Logger) (*pgxpool.Pool, error) {
	pool, err := CreatePool(dbUrl, logger)

	dropTableSql, err := os.ReadFile("./db/DropTables.sql")
	if err != nil {
		return nil, err
	}

	_, err = pool.Exec(context.Background(), string(dropTableSql))
	if err != nil {
		return nil, err
	}

	return pool, err
}

func CreatePool(dbUrl string, logger *zap.Logger) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		logger.Error("Unable to parse connection string")
		os.Exit(1)
	}

	poolConfig.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(
		context.Background(),
		poolConfig,
	)
	if err != nil {
		logger.Error("Unable to create connection pool")
		os.Exit(1)
	}

	return pool, err
}