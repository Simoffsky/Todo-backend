package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ConnectToDB(connStr string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
