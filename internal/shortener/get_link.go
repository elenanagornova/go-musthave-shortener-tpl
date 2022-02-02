package shortener

import "fmt"

var (
	ErrLinkNotFound = fmt.Errorf("link not found")
)

// GetLink returns full link by short link
func (s *Shortener) GetLink(url string) (string, error) {
	link := s.userLinks

	for _, links := range link {
		for _, link := range links {
			if link.ShortURL == url {
				return link.OriginalURL, nil
			}
		}
	}
	return "", ErrLinkNotFound
}

func (s *Shortener) GetLinks(userID string) []UserLinks {
	links, ok := s.userLinks[userID]
	if !ok {
		return []UserLinks{}
	}
	return links
}
