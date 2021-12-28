package shortener

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateShortLink from full link
func (s *Shortener) GenerateShortLink(url string) string {
	id := s.generateRandomString(5)
	s.linksMap[id] = url
	return fmt.Sprintf("http://%s/%s", s.addr, id)
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
