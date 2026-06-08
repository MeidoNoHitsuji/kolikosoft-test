package db

import (
	"context"

	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/config"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/shutdowner"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func NewRedis(ctx context.Context, cfg config.Config, log *zerolog.Logger) *redis.Client {
	opt, err := redis.ParseURL(cfg.Redis.URL())
	if err != nil {
		log.Fatal().Err(err).Msg("не удалось разобрать redis url")
	}

	client := redis.NewClient(opt)
	if err = client.Ping(ctx).Err(); err != nil {
		log.Fatal().Err(err).Msg("не удалось подключиться к redis")
	}

	return client
}

func ShutdownRedis(sd *shutdowner.Shutdowner, client *redis.Client) {
	sd.AddShutdownOption(client.Close, shutdowner.PriorityLayerStorage)
}
