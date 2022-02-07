package shortener

import (
	"go-musthave-shortener-tpl/internal/entity"
)

// GetLink returns full link by short link
func (s *Shortener) GetLink(url string) (string, error) {
	return s.Repo.FindOriginLinkByShortLink(url)
}

func (s *Shortener) GetLinks(userUID string) []entity.UserLinks {
	return s.Repo.GetLinksByuserUID(userUID)
}

func (s *Shortener) GetOriginalByShort(shortLink string) (string, error) {
	return s.Repo.FindOriginLinkByShortLink(shortLink)
}
