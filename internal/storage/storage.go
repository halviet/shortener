package storage

import "fmt"

type Storage interface {
	SaveURL(shortURL, origin string)
	GetOrigin(shortURL string) (origin string)
}

type ShortURL struct {
	Origin string
	Short  string
}

func New() *Store {
	return &Store{}
}

type Store struct {
	urls []ShortURL
}

func (s *Store) SaveURL(url ShortURL) {
	s.urls = append(s.urls, url)
}

func (s *Store) GetOrigin(shortURL string) (origin string, err error) {
	for _, url := range s.urls {
		if url.Short == shortURL {
			return url.Origin, nil
		}
	}

	return "", fmt.Errorf("url not found")
}
