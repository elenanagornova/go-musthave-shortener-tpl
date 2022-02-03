package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type DBRepo struct {
	conn *pgx.Conn
}

func (D *DBRepo) FinalSave() error {
	return D.conn.Close(context.Background())
}

func (D *DBRepo) FindOriginLinkByShortLink(shortLink string) (string, error) {
	return "", nil
}

func (D *DBRepo) SaveLinks(shortLink string, originalLink string, userUID string) error {
	_, err := D.conn.Exec(context.Background(), "insert into shortener.links (short_link, original_link, user_uid) values ($1, $2, $3)", shortLink, originalLink, userUID)
	return err
}

func (D *DBRepo) CreateUser(userUID string) error {
	panic("implement me")
}

func (D *DBRepo) GetLinksByuserUID(userUID string) []UserLinks {
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
		return nil, err
	}

	pgRepo := DBRepo{
		conn: conn,
	}

	m, err := RunMigration(databaseDSN)
	if err != nil && !m {
		return nil, err
	}
	return &pgRepo, nil
}
