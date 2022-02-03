package repository

type Storager interface {
	FindOriginLinkByShortLink(shortLink string) (string, error)

	SaveLinks(shortLink string, originalLink string, userUID string) error

	CreateUser(userUID string) error

	GetLinksByuserUID(userUID string) []UserLinks

	Ping() error

	Close()

	FinalSave() error
}
