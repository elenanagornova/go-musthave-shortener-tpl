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
	query := `select short_link, original_link, user_uid from shortener.links where short_link = $1`
	var links UserLinks
	result := D.conn.QueryRow(context.Background(), query, shortLink)
	if err := result.Scan(&links.ShortURL, &links.OriginalURL); err != nil {
		return "", err
	}
	return links.OriginalURL, nil
}

func (D *DBRepo) SaveLinks(shortLink string, originalLink string, userUID string) error {
	_, err := D.conn.Exec(context.Background(), "insert into shortener.links (short_link, original_link, user_uid) values ($1, $2, $3)", shortLink, originalLink, userUID)
	return err
}

func (D *DBRepo) CreateUser(userUID string) error {
	return nil
}

func (D *DBRepo) GetLinksByuserUID(userUID string) []UserLinks {
	query := `select short_link, original_link, user_uid from shortener.links where user_uid = $1`
	var result []UserLinks
	rows, err := D.conn.Query(context.Background(), query, userUID)
	if err != nil {
		return nil
	}
	for rows.Next() {
		var link UserLinks
		err = rows.Scan(&link.ShortURL, &link.OriginalURL)
		if err != nil {
			return nil
		}
		result = append(result, link)
	}
	err = rows.Err()
	if err != nil {
		return nil
	}
	return result
}

func (D *DBRepo) Ping() error {
	return D.conn.Ping(context.Background())
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
	pgRepo.conn.Exec(context.Background(), "insert into shortener.links (short_link, original_link, user_uid) values ($1, $2, $3)", "shortLink", "originalLink", "userUID")

	return &pgRepo, nil
}
