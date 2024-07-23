package storage

type Storage interface {
	CreateURL(origin string) (url string)
	GetURL(url string) (origin string)
}
