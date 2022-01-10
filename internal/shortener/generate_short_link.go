package shortener

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

// GenerateShortLink from full link
func (s *Shortener) GenerateShortLink(originalURL string) (string, error) {
	id := s.generateRandomString(5)
	s.linksMap[id] = originalURL
	u, err := url.Parse(s.addr)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}
	u.Path = id
	return u.String(), nil
}

func (s *Shortener) generateRandomString(n int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := ""
	randSrc := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(randSrc)
	for i := 0; i < n; i++ {
		result += string(letters[rnd.Intn(len(letters))])
	}
	return result
}
