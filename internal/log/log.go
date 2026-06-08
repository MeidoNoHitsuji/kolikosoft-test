package log

import (
	"os"
	"sync"
	"time"

	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/config"
	"github.com/rs/zerolog"
	zl "github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var once sync.Once
var Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Timestamp().Logger()

func InitLogger(cfg *config.Config) {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano
		zerolog.TimestampFieldName = "log_time"
		zerolog.MessageFieldName = "msg"

		logCtx := Logger.With()
		level, err := zerolog.ParseLevel(cfg.Log.Level)
		if err != nil {
			zl.Error().Err(err).Msg("уровень логирования не определен, установлен уровень trace")
			level = zerolog.TraceLevel
		}

		Logger = logCtx.Logger().Level(level) //.Hook(NewHook(cfg))
	})
}
