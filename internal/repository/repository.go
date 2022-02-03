package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
)

func CreateDBConnect(databaseDSN string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), databaseDSN)
	if err != nil {
		return nil, err
	}
	return conn, nil

}
