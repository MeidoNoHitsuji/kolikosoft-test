package main

import (
	"context"
	"flag"
	"time"

	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/cache"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/client"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/controller"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/db"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/layer"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/log"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/rate"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/repository"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/service"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/shutdowner"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/web"
	"github.com/MeidoNoHitsuji/kolikosoft-test/internal/worker"
)

func main() {
	ctx := context.Background()

	cfg := layer.NewConfig(parseFlags())
	log.InitLogger(&cfg)

	sd := shutdowner.NewShutdowner(ctx, &log.Logger)

	database := db.NewDatabase(ctx, cfg, &log.Logger)
	db.MustMigrate(ctx, database, &log.Logger)
	db.ShutdownPostgres(sd, database)

	redis := db.NewRedis(ctx, cfg, &log.Logger)
	db.ShutdownRedis(sd, redis)
	itemCache := cache.NewItemCache(redis)

	limiter := rate.New()
	limiter.Configure(client.SkinportItemsRateKey, 8.0/(5*time.Minute).Seconds(), 8)

	skinportClient := client.NewSkinportClient(cfg, limiter)
	skinportWorker := worker.NewSkinportWorker(skinportClient, itemCache, &log.Logger)
	skinportWorker.Run(ctx, sd)

	rep := repository.New(database, &log.Logger, cfg)

	srvHolder := service.NewServiceHolder(rep, itemCache, &log.Logger)
	ctrHolder := controller.NewControllerHolder(srvHolder)

	serve := web.NewServe(ctrHolder, &log.Logger, cfg, sd)
	serve.Run()
	serve.Shutdown(ctx)

	sd.WaitShutdown()
}

func parseFlags() string {
	var configPath string
	flag.StringVar(&configPath, "config", "./config/application.json", "path to config file")
	flag.Parse()

	return configPath
}
