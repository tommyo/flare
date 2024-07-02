package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tommyo/flare"
)

func main() {
	config := flare.NewConfig()
	config.RegisterDefault("server.addr", "a", "localhost:8080", "host to bind to")
	config.RegisterDefault("server.deadline", "", "10s", "time to wait for the server to shutdown")

	server := flare.NewServer(config)

	sessionStore := flare.NewSessionStore(config)

	if err := config.Parse(); err != nil {
		config.Usage()
		os.Exit(1)
	}

	sessionStore.Init()

	public := server.Build(sessionStore)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(config.Duration("server.deadline")))
	defer cancel()

	go func() {
		public.ListenAndServe()
	}()
	defer func() {
		if err := public.Shutdown(ctx); err != nil {
			fmt.Println("error when shutting down the main server: ", err)
		}
	}()

	<-sigs
}
