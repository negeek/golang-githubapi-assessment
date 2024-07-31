package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

var PostgreSQLDB *pgxpool.Pool

func Connect(url string) error {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return err
	}

	PostgreSQLDB, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return err
	}

	return nil
}
