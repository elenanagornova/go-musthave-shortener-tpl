package shortener

import "fmt"

func (s *Shortener) GetLink(url string) (string, error) {
	link, ok := s.linksMap[url]
	if !ok {
		return "", fmt.Errorf("link not found")
	}
	return link, nil
}
