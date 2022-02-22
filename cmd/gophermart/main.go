package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/sbxb/loyalty/api"
	"github.com/sbxb/loyalty/internal/logger"
)

func main() {
	logger.SetLevel("DEBUG")

	router := api.NewRouter()
	server, _ := api.NewHTTPServer("localhost:8888", router)
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
