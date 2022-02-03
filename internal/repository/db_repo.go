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
	return D.conn.Close(context.Background())
}

func (D *DBRepo) FindOriginLinkByShortLink(shortLink string) (string, error) {
	return "", nil
}

func (D *DBRepo) SaveLinks(shortLink string, originalLink string, userUID string) error {
	_, err := D.conn.Exec(context.Background(), "insert into shortener.links (short_link, original_link, user_uid) values ($1, $2, $3)", shortLink, originalLink, userUID)
	return err
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
	//migration := `
	//	CREATE SCHEMA IF NOT EXISTS shortener;
	//	-- DROP SCHEMA shortener CASCADE ;
	//	-- CREATE SCHEMA shortener;
	//	SET SEARCH_PATH TO shortener;
	//
	//	CREATE TABLE IF NOT EXISTS links(
	//		  id           serial primary key,
	//          short_link      varchar,
	//          original_link varchar,
	//          user_uid          varchar
	//	);
	//	ALTER TABLE links ALTER COLUMN created_at SET DEFAULT now();
	//	ALTER TABLE links ALTER COLUMN removed SET DEFAULT false;
	//	CREATE UNIQUE INDEX IF NOT EXISTS original_url_idx ON links USING btree (original_url);
	//	`

	m, err := RunMigration(databaseDSN)
	if err != nil && !m {
		return nil, err
	}
	return &pgRepo, nil
}
