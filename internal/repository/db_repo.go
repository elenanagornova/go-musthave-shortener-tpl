package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

type DBRepo struct {
	conn *pgx.Conn
}

func (D *DBRepo) FinalSave() error {
	// do nothing
	return nil
}

func (D *DBRepo) FindOriginLinkByShortLink(shortLink string) (string, error) {
	panic("implement me")
}

func (D *DBRepo) SaveLinks(shortLink string, originalLink string, userUID string) error {
	panic("implement me")
}

func (D *DBRepo) CreateUser(userUid string) error {
	panic("implement me")
}

func (D *DBRepo) GetLinksByUserUID(userUid string) []UserLinks {
	panic("implement me")
}

func (D *DBRepo) Ping() error {
	panic("implement me")
}

func (D *DBRepo) Close() {
	D.conn.Close(context.Background())
}

func NewDBConnect(databaseDSN string) (*DBRepo, error) {
	conn, err := pgx.Connect(context.Background(), databaseDSN)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v\n", err))
		return nil, err
	}
	//defer conn.Close(context.Background())

	pgRepo := DBRepo{
		conn: conn,
	}
	m, err := RunMigration(databaseDSN)
	if err != nil && !m {
		return nil, err
	}
	return &pgRepo, nil
}
