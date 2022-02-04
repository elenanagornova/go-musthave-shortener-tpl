package repository

import (
	"encoding/json"
	"go-musthave-shortener-tpl/internal/entity"
	"os"
)

type FRepo struct {
	filepath  string
	userLinks map[string][]entity.UserLinks
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

func (F FRepo) FindOriginLinkByShortLink(shortLink string) (string, error) {
	link := F.userLinks

	for _, links := range link {
		for _, link := range links {
			if link.ShortURL == shortLink {
				return link.OriginalURL, nil
			}
		}
	}
	return "", ErrLinkNotFound
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
