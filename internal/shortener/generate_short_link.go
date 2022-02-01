package shortener

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

// GenerateShortLink from full link
func (s *Shortener) GenerateShortLink(originalURL string, userId string) (string, error) {
	id := GenerateRandomString(5)

	u, err := url.Parse(s.addr)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}
	u.Path = id
	if _, exist := s.userLinks[userId]; !exist {
		s.userLinks[userId] = []UserLinks{}
	}
	s.userLinks[userId] = append(s.userLinks[userId], UserLinks{ShortURL: id, OriginalURL: originalURL})
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

//// GenerateShortLink from full link
//func (s *Shortener) GenerateShortLink(OriginalURL string, userId string) (string, error) {
//	id := s.GenerateRandomString(5)
//	u, err := url.Parse(s.addr)
//	if err != nil {
//		return "", fmt.Errorf("failed to parse url: %w", err)
//	}
//	u.Path = id
//	if _, exist := s.userLinks[userId]; !exist {
//		s.userLinks[userId] = []UserLinks{}
//	}
//	s.userLinks[userId] = append(s.userLinks[userId], UserLinks{ShortURL: id, OriginalURL: OriginalURL})
//	return u.String(), nil
//}
