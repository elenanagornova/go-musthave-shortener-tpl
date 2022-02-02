package shortener

import "fmt"

var (
	ErrLinkNotFound = fmt.Errorf("link not found")
)

// GetLink returns full link by short link
func (s *Shortener) GetLink(url string, userId string) (string, error) {
	link, ok := s.userLinks[userId]
	if !ok {
		return "", fmt.Errorf("link not found")
	}
	for _, links := range link {
		if links.ShortURL == url {
			return links.OriginalURL, nil
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
