package shortener

// Shortener service for shortener links
type Shortener struct {
	linksMap map[string]string
	addr     string
}

// New shortener instance
func New(addr string) *Shortener {
	return &Shortener{
		addr:     addr,
		linksMap: make(map[string]string),
	}
}
