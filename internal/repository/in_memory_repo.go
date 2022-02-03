package repository

import (
	"fmt"
)

var (
	ErrLinkNotFound = fmt.Errorf("link not found")
)

type MRepo struct {
	userLinks map[string][]UserLinks
}

func (M MRepo) FinalSave() error {
	// do nothing
	return nil
}

func (M MRepo) GetLinksByuserUID(userUID string) []UserLinks {
	links, ok := M.userLinks[userUID]
	if !ok {
		return []UserLinks{}
	}
	return links
}

func (M MRepo) Ping() error {
	return nil
}

func (M MRepo) FindOriginLinkByShortLink(shortLink string) (string, error) {
	link := M.userLinks

	for _, links := range link {
		for _, link := range links {
			if link.ShortURL == shortLink {
				return link.OriginalURL, nil
			}
		}
	}
	return "", ErrLinkNotFound
}

func (M MRepo) SaveLinks(shortLink string, originalLink string, userUID string) error {
	if _, exist := M.userLinks[userUID]; !exist {
		M.userLinks[userUID] = []UserLinks{}
	}

	M.userLinks[userUID] = append(M.userLinks[userUID], UserLinks{ShortURL: shortLink, OriginalURL: originalLink})
	return nil
}

func (M MRepo) CreateUser(userUID string) error {
	return nil
}

func (MRepo) Close() {
	// do nothing
}

func NewInMemoryConnect() *MRepo {
	return &MRepo{userLinks: map[string][]UserLinks{}}
}
