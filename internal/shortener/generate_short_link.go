package shortener

import (
	"fmt"
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
	if err := s.Repo.SaveLinks(id, originalURL, userUID); err != nil {
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
