package shortener

import "github.com/jackc/pgx/v4"

// Shortener service for shortener links
type Shortener struct {
	DbConn    *pgx.Conn
	linksMap  map[string]string
	addr      string
	userLinks map[string][]UserLinks
}

// New shortener instance
func New(addr string, dbConn *pgx.Conn) *Shortener {
	return &Shortener{
		addr:      addr,
		DbConn:    dbConn,
		linksMap:  make(map[string]string),
		userLinks: map[string][]UserLinks{},
	}
}

type UserLinks struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
