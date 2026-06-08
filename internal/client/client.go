package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/config"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/model"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/rate"
	"github.com/andybalholm/brotli"
	"github.com/go-resty/resty/v2"
)

const SkinportItemsRateKey = "skinport:items"

type SkinportClient struct {
	client  *resty.Client
	limiter *rate.Limiter
}

type SkinportItem = model.Item

type SkinportTradableItem = model.Item

func NewSkinportClient(cfg config.Config, limiter *rate.Limiter) *SkinportClient {
	return &SkinportClient{
		client: resty.New().
			SetBaseURL(cfg.App.SkinportUrl).
			SetHeader("Accept-Encoding", "br").
			SetTimeout(30 * time.Second),
		limiter: limiter,
	}
}

func (c *SkinportClient) GetAllItems(ctx context.Context) ([]SkinportItem, error) {
	notTradable := false
	body, err := c.requestItems(ctx, &notTradable)
	if err != nil {
		return nil, err
	}

	var items []SkinportItem
	if err = json.Unmarshal(body, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (c *SkinportClient) GetTradableItems(ctx context.Context) ([]SkinportTradableItem, error) {
	body, err := c.requestItems(ctx, nil)
	if err != nil {
		return nil, err
	}

	var items []SkinportTradableItem
	if err = json.Unmarshal(body, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (c *SkinportClient) requestItems(ctx context.Context, tradable *bool) ([]byte, error) {
	// Технически лимитёр нам тут не нужен, ибо я ограничил воркер таймером в 5 минут.
	// Однако если нужно будет больше, то данный лимитёр не даст нам выйти за пределы rate-limit'а
	c.limiter.Wait(SkinportItemsRateKey)

	req := c.client.R().
		SetContext(ctx).
		SetDoNotParseResponse(true)
	if tradable != nil {
		req.SetQueryParam("tradable", strconv.FormatBool(*tradable))
	}

	resp, err := req.Get("/items")
	if err != nil {
		return nil, err
	}
	defer resp.RawBody().Close()

	if resp.IsError() {
		c.handleRateLimit(resp)
		return nil, fmt.Errorf("skinport вернул ошибочный статус запроса: %s", resp.Status())
	}

	body, err := readBody(resp)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *SkinportClient) handleRateLimit(resp *resty.Response) {
	if resp.StatusCode() != http.StatusTooManyRequests {
		return
	}

	retryAfter := resp.Header().Get("Retry-After")
	if retryAfter == "" {
		return
	}

	seconds, err := strconv.Atoi(retryAfter)
	if err != nil {
		return
	}

	c.limiter.RateLimited(SkinportItemsRateKey, time.Duration(seconds)*time.Second)
}

func readBody(resp *resty.Response) ([]byte, error) {
	var reader io.Reader = resp.RawBody()
	if strings.EqualFold(resp.Header().Get("Content-Encoding"), "br") {
		reader = brotli.NewReader(reader)
	}

	return io.ReadAll(reader)
}
