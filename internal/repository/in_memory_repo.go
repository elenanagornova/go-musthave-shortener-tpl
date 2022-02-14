package repository

import (
	"context"
	"fmt"
	"go-musthave-shortener-tpl/internal/deleter"
	"go-musthave-shortener-tpl/internal/entity"
)

var (
	ErrLinkNotFound = fmt.Errorf("link not found")
)

type MRepo struct {
	userLinks map[string][]entity.UserLinks
}

func (M MRepo) BatchUpdateLinks(ctx context.Context, task deleter.DeleteTask) error {
	return nil
}

func (M MRepo) FindShortLinkByOriginLink(originLink string) (string, error) {
	return "", nil
}

func (M MRepo) BatchSaveLinks(links []entity.DBBatchShortenerLinks) ([]entity.DBBatchShortenerLinks, error) {
	return []entity.DBBatchShortenerLinks{}, nil
}

func (M MRepo) FinalSave() error {
	// do nothing
	return nil
}

func (M MRepo) GetLinksByuserUID(userUID string) []entity.UserLinks {
	links, ok := M.userLinks[userUID]
	if !ok {
		return []entity.UserLinks{}
	}
	return links
}

func (M MRepo) Ping() error {
	return nil
}

func (M MRepo) FindOriginLinkByShortLink(shortLink string) (entity.UserLinks, error) {
	link := M.userLinks

	for _, links := range link {
		for _, link := range links {
			if link.ShortURL == shortLink {

				return entity.UserLinks{ShortURL: shortLink, OriginalURL: link.OriginalURL}, nil
			}
		}
	}
	return entity.UserLinks{}, ErrLinkNotFound
}

func (M MRepo) SaveLinks(shortLink string, originalLink string, userUID string) error {
	if _, exist := M.userLinks[userUID]; !exist {
		M.userLinks[userUID] = []entity.UserLinks{}
	}

	M.userLinks[userUID] = append(M.userLinks[userUID], entity.UserLinks{ShortURL: shortLink, OriginalURL: originalLink})
	return nil
}

func (M MRepo) CreateUser(userUID string) error {
	return nil
}

func (MRepo) Close() {
	// do nothing
}

func NewInMemoryConnect() *MRepo {
	return &MRepo{userLinks: map[string][]entity.UserLinks{}}
}
