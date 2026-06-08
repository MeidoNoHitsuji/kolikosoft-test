package db

import (
	"context"

	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/config"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/shutdowner"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

func NewDatabase(_ context.Context, cfg config.Config, log *zerolog.Logger) *sqlx.DB {
	db, err := sqlx.Connect("pgx", cfg.DB.URL())
	if err != nil {
		log.Fatal().Err(err).Msg("не удалось подключиться к базе данных")
	}
	return db
}

func Shutdown(sd *shutdowner.Shutdowner, db *sqlx.DB) {
	sd.AddShutdownOption(db.Close, shutdowner.PriorityLayerStorage)
}
