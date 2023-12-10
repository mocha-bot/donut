package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"buf.build/gen/go/mocha/remcall/connectrpc/go/donut/v1/donutv1connect"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	cfg, err := Get()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get config")
	}

	db, err := NewDatabaseInstance(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database instance")
	}

	mux := http.NewServeMux()

	// Create instances
	repo := NewDonutRepository(db)
	donut := NewDonutCall(repo)
	handler := NewHandler(donut)

	mmPath, mmHandler := donutv1connect.NewMatchMakerServiceHandler(handler)
	pPath, pHandler := donutv1connect.NewPeopleServiceHandler(handler)

	mux.Handle(mmPath, mmHandler)
	mux.Handle(pPath, pHandler)

	server := &http.Server{
		Addr:    cfg.ApplicationConfig.Address(),
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	// Run the server in a goroutine so that it doesn't block
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to serve")
		}

		log.Info().Msgf("server is listening on %s", cfg.ApplicationConfig.Address())
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info().Msg("server is shutting down")

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to shutdown server gracefully")
	}
}
