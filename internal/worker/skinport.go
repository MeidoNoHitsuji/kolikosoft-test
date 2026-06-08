package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/cache"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/client"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/model"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/shutdowner"
	"github.com/rs/zerolog"
)

const skinportRefreshInterval = 5 * time.Minute

type SkinportWorker struct {
	client *client.SkinportClient
	cache  *cache.ItemCache
	log    *zerolog.Logger
}

func NewSkinportWorker(client *client.SkinportClient, cache *cache.ItemCache, log *zerolog.Logger) *SkinportWorker {
	return &SkinportWorker{
		client: client,
		cache:  cache,
		log:    log,
	}
}

func (w *SkinportWorker) Run(ctx context.Context, sd *shutdowner.Shutdowner) {
	workerCtx, cancel := context.WithCancel(ctx)
	sd.AddShutdownOption(func() error {
		cancel()
		return nil
	}, shutdowner.PriorityLayerWeb)

	go w.run(workerCtx)
}

func (w *SkinportWorker) run(ctx context.Context) {
	w.refresh(ctx)

	ticker := time.NewTicker(skinportRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.refresh(ctx)
		}
	}
}

func (w *SkinportWorker) refresh(ctx context.Context) {
	items, err := w.loadItems(ctx)
	if err != nil {
		w.log.Err(err).Msg("не удалось обновить кеш предметов skinport")
		return
	}

	if err = w.cache.SetItems(ctx, items); err != nil {
		w.log.Err(err).Msg("не удалось сохранить предметы skinport в кеш")
		return
	}

	w.log.Info().Int("count", len(items)).Msg("кеш предметов skinport обновлен")
}

func (w *SkinportWorker) loadItems(ctx context.Context) ([]model.Item, error) {
	var (
		wg            sync.WaitGroup
		tradableItems []client.SkinportTradableItem
		allItems      []client.SkinportItem
		tradableErr   error
		allItemsErr   error
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		tradableItems, tradableErr = w.client.GetTradableItems(ctx)
	}()
	go func() {
		defer wg.Done()
		allItems, allItemsErr = w.client.GetAllItems(ctx)
	}()
	wg.Wait()

	if tradableErr != nil {
		return nil, fmt.Errorf("ошибка при получении tradable предметов: %w", tradableErr)
	}
	if allItemsErr != nil {
		return nil, fmt.Errorf("ошибка при получении всех предметов: %w", allItemsErr)
	}

	return mergeItems(allItems, tradableItems), nil
}

func mergeItems(allItems []client.SkinportItem, tradableItems []client.SkinportTradableItem) []model.Item {
	itemsByName := make(map[string]model.Item, len(allItems))
	itemNames := make([]string, 0, len(allItems))

	for _, item := range allItems {
		name := item.MarketHashName()
		if name == "" {
			continue
		}
		if _, ok := itemsByName[name]; ok {
			continue
		}

		item[model.TradableMinPriceKey] = nullRawMessage()
		itemsByName[name] = item
		itemNames = append(itemNames, name)
	}

	for _, item := range tradableItems {
		name := item.MarketHashName()
		if name == "" {
			continue
		}

		price := rawOrNull(item[model.MinPriceKey])
		cachedItem, ok := itemsByName[name]
		if ok {
			cachedItem[model.TradableMinPriceKey] = price
			continue
		}

		item[model.TradableMinPriceKey] = price
		item[model.MinPriceKey] = nullRawMessage()
		itemsByName[name] = item
		itemNames = append(itemNames, name)
	}

	result := make([]model.Item, 0, len(itemsByName))
	for _, name := range itemNames {
		result = append(result, itemsByName[name])
	}

	return result
}

func rawOrNull(value json.RawMessage) json.RawMessage {
	if len(value) == 0 {
		return nullRawMessage()
	}

	return append(json.RawMessage(nil), value...)
}

func nullRawMessage() json.RawMessage {
	return json.RawMessage("null")
}
