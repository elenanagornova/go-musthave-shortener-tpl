package repository

type Storager interface {
	FindOriginLinkByShortLink(shortLink string) (string, error)

	SaveLinks(shortLink string, originalLink string, userUID string) error

	CreateUser(userUid string) error

	GetLinksByUserUID(userUid string) []UserLinks

	Ping() error

	Close()

	FinalSave() error
}
