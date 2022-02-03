package repository

import (
	"encoding/json"
	"os"
)

type FRepo struct {
	filepath  string
	userLinks map[string][]UserLinks
}

func (F FRepo) FinalSave() error {
	fp, err := os.Create(F.filepath)
	if err != nil {
		return err
	}
	defer fp.Close()
	return json.NewEncoder(fp).Encode(F.userLinks)
}

func (F FRepo) GetLinksByUserUID(userUid string) []UserLinks {
	links, ok := F.userLinks[userUid]
	if !ok {
		return []UserLinks{}
	}
	return links
}

func (F FRepo) Ping() error {
	panic("implement me")
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
		F.userLinks[userUID] = []UserLinks{}
	}

	F.userLinks[userUID] = append(F.userLinks[userUID], UserLinks{ShortURL: shortLink, OriginalURL: originalLink})
	return nil
}

func (F FRepo) CreateUser(userUid string) error {
	panic("implement me")
}
func (FRepo) Close() {
	// do nothing
}

func NewFileConnect(fileStoragePath string) *FRepo {
	return &FRepo{
		filepath:  fileStoragePath,
		userLinks: map[string][]UserLinks{},
	}
}
