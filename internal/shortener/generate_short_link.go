package shortener

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

// GenerateShortLink from full link
func (s *Shortener) GenerateShortLink(originalURL string, userUId string) (string, error) {
	id := GenerateRandomString(5)

	u, err := url.Parse(s.addr)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}
	u.Path = id
	if _, exist := s.userLinks[userUId]; !exist {
		s.userLinks[userUId] = []UserLinks{}
	}
	s.userLinks[userUId] = append(s.userLinks[userUId], UserLinks{ShortURL: id, OriginalURL: originalURL})
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
