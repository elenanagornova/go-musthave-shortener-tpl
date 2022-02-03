package shortener

// Shortener service for shortener links
type Shortener struct {
	Repo Storager
	Addr string
}

// New shortener instance
func New(addr string, storager Storager) *Shortener {
	return &Shortener{
		Addr: addr,
		Repo: storager,
	}
}
