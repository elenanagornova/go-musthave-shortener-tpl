package shortener

import (
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"math/rand"
	"net/url"
	"time"
)

// GenerateShortLink from full link
func (s *Shortener) GenerateShortLink(originalURL string, userUID string) (string, error) {
	id := GenerateRandomString(5)

	u, err := url.Parse(s.Addr)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}
	u.Path = id
	if error := s.Repo.SaveLinks(id, originalURL, userUID); error != nil {
		var pgErr *pgconn.PgError
		if errors.As(error, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			path, err := s.Repo.FindShortLinkByOriginLink(originalURL)
			if err != nil {
				return "", err
			}
			return s.Addr + path, error
		}
		return "", err
	}
	return u.String(), nil
}

func GenerateRandomString(n int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := ""
	randSrc := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(randSrc)
	for i := 0; i < n; i++ {
		result += string(letters[rnd.Intn(len(letters))])
	}
	return result
}
