package shortener

import (
	"go-musthave-shortener-tpl/internal/deleter"
	"go-musthave-shortener-tpl/internal/entity"
)

type Storager interface {
	FindOriginLinkByShortLink(shortLink string) (entity.UserLinks, error)

	FindShortLinkByOriginLink(originLink string) (string, error)

	SaveLinks(shortLink string, originalLink string, userUID string) error

	BatchSaveLinks(links []entity.DBBatchShortenerLinks) ([]entity.DBBatchShortenerLinks, error)

	BatchUpdateLinks(tsk deleter.DeleteTask) error

	CreateUser(userUID string) error

	GetLinksByuserUID(userUID string) []entity.UserLinks

	Ping() error

	Close()

	FinalSave() error
}
