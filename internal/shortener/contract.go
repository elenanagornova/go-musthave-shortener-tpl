package shortener

import "go-musthave-shortener-tpl/internal/repository"

type Storager interface {
	FindOriginLinkByShortLink(shortLink string) (string, error)

	SaveLinks(shortLink string, originalLink string, userUID string) error

	CreateUser(userUid string) error

	GetLinksByUserUID(userUid string) []repository.UserLinks

	Ping() error

	Close()

	FinalSave() error
}
