package shortener

// Shortener service for shortener links
type Shortener struct {
	linksMap  map[string]string
	addr      string
	userLinks map[string][]UserLinks
}

// New shortener instance
func New(addr string) *Shortener {
	return &Shortener{
		addr:      addr,
		linksMap:  make(map[string]string),
		userLinks: map[string][]UserLinks{},
	}
}

type UserLinks struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
