package service

import (
	"context"

	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/cache"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/model"
	"github.com/rs/zerolog"
)

type ItemService struct {
	cache *cache.ItemCache
	log   *zerolog.Logger
}

func NewItemService(cache *cache.ItemCache, log *zerolog.Logger) *ItemService {
	return &ItemService{
		cache: cache,
		log:   log,
	}
}

func (s *ItemService) GetItems(ctx context.Context) ([]model.Item, error) {
	items, err := s.cache.GetItems(ctx)
	if err != nil {
		s.log.Err(err).Msg("не удалось получить предметы из кеша")
		return nil, err
	}

	return items, nil
}
