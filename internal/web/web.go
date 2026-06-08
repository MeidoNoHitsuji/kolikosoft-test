package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/config"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/controller"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/shutdowner"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Serve struct {
	handler   *gin.Engine
	server    *http.Server
	log       *zerolog.Logger
	sd        *shutdowner.Shutdowner
	cfg       config.Config
	ctrHodler *controller.ControllerHolder
}

func NewServe(ctrHodler *controller.ControllerHolder, log *zerolog.Logger, cfg config.Config, sd *shutdowner.Shutdowner) *Serve {
	srv := &Serve{
		handler: gin.Default(),
		log:     log,
		sd:      sd,
		cfg:     cfg,
	}

	srv.registerRoutes(ctrHodler)
	return srv
}

func (srv *Serve) registerRoutes(ctrHodler *controller.ControllerHolder) {
	api := srv.handler.Group("/api")
	v1 := api.Group("/v1")

	for _, controllerInterface := range ctrHodler.Controllers {
		controllerInterface.RegisterRoutes(v1)
	}
}

func (srv *Serve) Run() {
	srv.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", srv.cfg.App.WebPort),
		Handler: srv.handler,
	}

	go func() {
		err := srv.server.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				srv.log.Info().Msg("web-сервер отключен")
			} else {
				srv.log.Error().Err(err).Msg("ошибка работы web-сервера")
			}
		}
	}()

	srv.log.Info().Msg("web-сервер запущен")
}

func (srv *Serve) Shutdown(ctx context.Context) {
	srv.sd.AddShutdownOption(func() error {
		return srv.server.Shutdown(ctx)
	}, shutdowner.PriorityLayerWeb)
}
