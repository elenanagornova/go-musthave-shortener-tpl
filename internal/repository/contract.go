package repository

import "go-musthave-shortener-tpl/internal/entity"

type Storager interface {
	FindOriginLinkByShortLink(shortLink string) (string, error)

	FindShortLinkByOriginLink(originLink string) (string, error)

	SaveLinks(shortLink string, originalLink string, userUID string) error

	BatchSaveLinks(links []entity.DBBatchShortenerLinks) ([]entity.DBBatchShortenerLinks, error)

	CreateUser(userUID string) error

	GetLinksByuserUID(userUID string) []entity.UserLinks

	Ping() error

	Close()

	FinalSave() error
}
