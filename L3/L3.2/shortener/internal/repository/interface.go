package repository

import "shortener/internal/domain"

type URLRepository interface {
	Save(u domain.URL) error
	Get(short string) (string, error)
}

type ClickRepository interface {
	SaveClick(c domain.Click) error
	GetStats(short string) (map[string]interface{}, error)
}
