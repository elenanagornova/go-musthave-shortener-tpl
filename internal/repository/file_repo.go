package repository

import (
	"context"
	"encoding/json"
	"go-musthave-shortener-tpl/internal/deleter"
	"go-musthave-shortener-tpl/internal/entity"
	"os"
)

type FRepo struct {
	filepath  string
	userLinks map[string][]entity.UserLinks
}

func (F FRepo) BatchUpdateLinks(ctx context.Context, task deleter.DeleteTask) error {
	return nil
}

func (F FRepo) FindShortLinkByOriginLink(originLink string) (string, error) {
	return "", nil
}

func (F FRepo) BatchSaveLinks(links []entity.DBBatchShortenerLinks) ([]entity.DBBatchShortenerLinks, error) {
	return []entity.DBBatchShortenerLinks{}, nil
}

func (F FRepo) FinalSave() error {
	fp, err := os.Create(F.filepath)
	if err != nil {
		return err
	}
	defer fp.Close()
	return json.NewEncoder(fp).Encode(F.userLinks)
}

func (F FRepo) GetLinksByuserUID(userUID string) []entity.UserLinks {
	links, ok := F.userLinks[userUID]
	if !ok {
		return []entity.UserLinks{}
	}
	return links
}

func (F FRepo) Ping() error {
	return nil
}

func (F FRepo) FindOriginLinkByShortLink(shortLink string) (entity.UserLinks, error) {
	userLinks := F.userLinks

	for _, links := range userLinks {
		for _, link := range links {
			if link.ShortURL == shortLink {
				return entity.UserLinks{ShortURL: shortLink, OriginalURL: link.OriginalURL}, nil
			}
		}
	}
	return entity.UserLinks{}, ErrLinkNotFound
}

func (F FRepo) SaveLinks(shortLink string, originalLink string, userUID string) error {
	if _, exist := F.userLinks[userUID]; !exist {
		F.userLinks[userUID] = []entity.UserLinks{}
	}

	F.userLinks[userUID] = append(F.userLinks[userUID], entity.UserLinks{ShortURL: shortLink, OriginalURL: originalLink})
	return nil
}

func (F FRepo) CreateUser(userUID string) error {
	return nil
}
func (FRepo) Close() {
	// do nothing
}

func NewFileConnect(fileStoragePath string) *FRepo {
	return &FRepo{
		filepath:  fileStoragePath,
		userLinks: map[string][]entity.UserLinks{},
	}
}
