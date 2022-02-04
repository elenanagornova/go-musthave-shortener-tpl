package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"go-musthave-shortener-tpl/internal/entity"
)

type DBRepo struct {
	conn *pgx.Conn
}

func (D *DBRepo) BatchSaveLinks(links []entity.DBBatchShortenerLinks) ([]entity.DBBatchShortenerLinks, error) {
	tx, err := D.conn.Begin(context.Background())
	if err != nil {
		return []entity.DBBatchShortenerLinks{}, err
	}
	defer tx.Rollback(context.Background())

	query := "insert into shortener.links(short_link, original_link, user_uid) values ($1, $2, $3)"

	for _, value := range links {
		_, err = tx.Exec(context.Background(), query, value.ShortURL, value.OriginalURL, value.UserUID)
		if err != nil {
			return []entity.DBBatchShortenerLinks{}, err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return []entity.DBBatchShortenerLinks{}, err
	}
	return links, nil
}

func (D *DBRepo) FinalSave() error {
	return D.conn.Close(context.Background())
}

func (D *DBRepo) FindOriginLinkByShortLink(shortLink string) (string, error) {
	query := `select short_link, original_link, user_uid from shortener.links where short_link = $1`
	var links entity.UserLinks
	result := D.conn.QueryRow(context.Background(), query, shortLink)
	if err := result.Scan(&links.ShortURL, &links.OriginalURL, &links.UserUID); err != nil {
		return "", err
	}
	return links.OriginalURL, nil
}

func (D *DBRepo) SaveLinks(shortLink string, originalLink string, userUID string) error {
	_, err := D.conn.Exec(context.Background(), "insert into shortener.links(short_link, original_link, user_uid) values ($1, $2, $3)", shortLink, originalLink, userUID)
	return err
}

func (D *DBRepo) CreateUser(userUID string) error {
	return nil
}

func (D *DBRepo) GetLinksByuserUID(userUID string) []entity.UserLinks {
	query := `select short_link, original_link, user_uid from shortener.links where user_uid = $1`
	var result []entity.UserLinks
	rows, err := D.conn.Query(context.Background(), query, userUID)
	if err != nil {
		return nil
	}
	for rows.Next() {
		var link entity.UserLinks
		err = rows.Scan(&link.ShortURL, &link.OriginalURL, &link.UserUID)
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
	pgRepo.conn.Exec(context.Background(), "insert into shortener.links(short_link, original_link, user_uid) values ($1, $2, $3)", "shortLink", "originalLink", "userUID")

	return &pgRepo, nil
}
