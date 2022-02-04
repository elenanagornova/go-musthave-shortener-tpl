package entity

type BatchRequest []BatchShortenerRequest

type BatchShortenerRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	links []BatchShortenerResponse
}
type BatchShortenerResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
