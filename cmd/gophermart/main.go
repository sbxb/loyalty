package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/sbxb/loyalty/api"
	"github.com/sbxb/loyalty/config"
	"github.com/sbxb/loyalty/internal/logger"
	"github.com/sbxb/loyalty/storage"
	"github.com/sbxb/loyalty/storage/inmemory"
	"github.com/sbxb/loyalty/storage/psql"
)

func main() {
	logger.SetLevel("DEBUG")

	cfg, err := config.New()
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Info("Config parsed")

	var store storage.Storage
	if cfg.DatabaseDSN != "" {
		store, err = psql.NewDBStorage(cfg.DatabaseDSN)
	} else {
		store, err = inmemory.NewMapStorage()
	}
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Info("Storage created")
	defer store.Close()

	router := api.NewRouter(store, cfg)
	server, _ := api.NewHTTPServer(cfg.ServerAddress, router)
	defer server.Close()

	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGTERM, syscall.SIGINT,
	)
	defer stop()

	go func() {
		err := server.Start(ctx)
		if err != nil {
			// Server either failed to start or exited unexpectedly
			// Cancel the context in order to let main() finish its work
			stop()
		}
	}()

	<-ctx.Done()
	server.Close()
}
