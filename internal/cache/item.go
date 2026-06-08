package cache

import (
	"context"
	"encoding/json"
	"errors"

	backErr "github.com/MeidoNoHitsuji/kolikosoft-test/internal/errors"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/model"
	"github.com/redis/go-redis/v9"
)

const (
	itemsKey = "skinport:items"
	// Отдельный ключ, чтобы если в момент перезаписи мы обратились к предмету, то он у нас был
	itemsTempKey = "skinport:items:tmp"
	hsetBatchLen = 500
)

type ItemCache struct {
	client *redis.Client
}

func NewItemCache(client *redis.Client) *ItemCache {
	return &ItemCache{
		client: client,
	}
}

func (c *ItemCache) SetItems(ctx context.Context, items []model.Item) error {
	if len(items) == 0 {
		return c.client.Del(ctx, itemsKey, itemsTempKey).Err()
	}

	pipe := c.client.Pipeline()
	pipe.Del(ctx, itemsTempKey)

	written := 0
	for _, item := range items {
		name := item.MarketHashName()
		if name == "" {
			continue
		}

		payload, err := json.Marshal(item)
		if err != nil {
			return err
		}

		pipe.HSet(ctx, itemsTempKey, name, payload)
		written++
		// Записываем предметы батчами по 500шт
		if written%hsetBatchLen == 0 {
			if _, err = pipe.Exec(ctx); err != nil {
				return err
			}
			pipe = c.client.Pipeline()
		}
	}

	if written == 0 {
		return c.client.Del(ctx, itemsKey, itemsTempKey).Err()
	}

	pipe.Rename(ctx, itemsTempKey, itemsKey)
	_, err := pipe.Exec(ctx)
	return err
}

func (c *ItemCache) GetItem(ctx context.Context, name string) (*model.Item, error) {
	payload, err := c.client.HGet(ctx, itemsKey, name).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, backErr.ErrItemNotFound
		}
		return nil, err
	}

	var item model.Item
	if err = json.Unmarshal(payload, &item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (c *ItemCache) GetItems(ctx context.Context) ([]model.Item, error) {
	payloads, err := c.client.HGetAll(ctx, itemsKey).Result()
	if err != nil {
		return nil, err
	}

	items := make([]model.Item, 0, len(payloads))
	for _, payload := range payloads {
		var item model.Item
		if err = json.Unmarshal([]byte(payload), &item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
