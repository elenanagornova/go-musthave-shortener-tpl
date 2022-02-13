package entity

type UserLinks struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserUID     string `json:"uid,omitempty"`
	Removed     bool   `json:"-"`
}

type DBBatchShortenerLinks struct {
	ShortURL      string
	OriginalURL   string
	UserUID       string
	CorrelationID string
}
