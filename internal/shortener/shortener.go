package shortener

// Shortener service for shortener links
type Shortener struct {
	userLinks map[string][]UserLinks
	addr      string
}

func (s *Shortener) GetLinks(userID string) []string {
	links, ok := s.userLinks[userID]
	if !ok {
		return []string{}
	}
	originalLinks := make([]string, 0, len(links))
	for _, link := range links {
		originalLinks = append(originalLinks, link.OriginalURL)
	}
	return originalLinks
}

// New shortener instance
func New(addr string) *Shortener {
	return &Shortener{
		addr:      addr,
		userLinks: map[string][]UserLinks{},
	}
}

type UserLinks struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
