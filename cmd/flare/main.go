package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/grpcreflect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/tommyo/flare"
	"github.com/tommyo/flare/proto/protoconnect"
)

func main() {
	config := flare.NewConfig()
	config.RegisterDefault("server.addr", "a", "localhost:8080", "host to bind to")
	config.RegisterDefault("server.deadline", "", "10s", "time to wait for the server to shutdown")

	server := flare.NewServer(config)

	if err := config.Parse(); err != nil {
		config.Usage()
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle(protoconnect.NewSparkConnectServiceHandler(server))

	reflector := grpcreflect.NewStaticReflector(protoconnect.SparkConnectServiceName)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	public := &http.Server{
		Addr:    config.String("addr"),
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

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
