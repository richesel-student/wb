package usecase

import (
	"log"
	"shortener/internal/domain"
	"shortener/internal/repository"
	"shortener/pkg/shortener"
)

type Cache interface {
	Get(key string) (string, error)
	Set(key, value string)
}

type ShortenerUseCase struct {
	urlRepo   repository.URLRepository
	clickRepo repository.ClickRepository
	cache     Cache
}

func NewShortenerUseCase(u repository.URLRepository, c repository.ClickRepository, cache Cache) *ShortenerUseCase {
	return &ShortenerUseCase{u, c, cache}
}

func (uc *ShortenerUseCase) Create(original, custom string) (string, error) {
	short := custom
	if short == "" {
		short = shortener.Generate(6)
	}

	err := uc.urlRepo.Save(domain.URL{
		Original: original,
		Short:    short,
	})

	return short, err
}

func (uc *ShortenerUseCase) Redirect(short, ua string) (string, error) {
	// cache
	if uc.cache != nil {
		val, err := uc.cache.Get(short)
		if err == nil {
			return val, nil
		}
	}

	original, err := uc.urlRepo.Get(short)
	if err != nil {
		return "", err
	}

	if err := uc.clickRepo.SaveClick(domain.Click{
		Short:     short,
		UserAgent: ua,
	}); err != nil {
		log.Printf("failed to save click: %v", err)
	}

	if uc.cache != nil {
		uc.cache.Set(short, original)
	}

	return original, nil
}

func (uc *ShortenerUseCase) Analytics(short string) (map[string]interface{}, error) {
	return uc.clickRepo.GetStats(short)
}
