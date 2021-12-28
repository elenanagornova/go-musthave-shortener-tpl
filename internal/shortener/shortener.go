package shortener

type Shortener struct {
	linksMap map[string]string
	addr     string
}

func New(addr string) *Shortener {
	return &Shortener{
		addr:     addr,
		linksMap: make(map[string]string),
	}
}
