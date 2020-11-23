package main

import (
	"context"
	"log"
	"time"

	"github.com/bitburst/burstconsumer/cmd/burst/burstcfg"
	"github.com/bitburst/burstconsumer/cmd/burst/models/user"
	"github.com/bitburst/burstconsumer/cmd/burst/router"
	"github.com/bitburst/burstconsumer/pkg/gracefulshutdown"
	"github.com/bitburst/burstconsumer/pkg/server"
	"github.com/joho/godotenv"
)

func main() {
	// load env vars from .env file if one is present
	godotenv.Load()

	// get program wide configuration that will have access to DB, logger and
	// generic configuration
	cfg, err := burstcfg.New()
	if err != nil {
		log.Fatalf("could not start program due to: %s", err)
	}

	// configure server instance
	srv := server.
		Get().
		WithAddr(cfg.GetAPIPort()).
		WithRouter(router.Get(cfg)).
		WithErrLogger(cfg.Errlog)

	// start server in separate goroutine so that it doesn't block graceful
	// shutdown handler
	go func() {
		cfg.Infolog.Printf("starting server at %s", cfg.GetAPIPort())
		if err := srv.Start(); err != nil {
			cfg.Errlog.Printf("starting server: %s", err)
		}
	}()

	// start sweeper in separate goroutine so that it doesn't block graceful
	// shutdown handler and will be able to remove old records from DB
	go func() {
		for {
			time.Sleep(1 * time.Second)
			t := time.Now().Add(-30 * time.Second).Format("2006-01-02 15:04:05")
			if err := user.RemoveOlderThan(context.Background(), cfg.DB.Client, t); err != nil {
				cfg.Errlog.Printf("performing sweeping: %s\n", err)
			}
		}
	}()

	// initiate graceful shutdown handler that will listen to api crash signals
	// and perform cleanup
	gracefulshutdown.Init(cfg.Errlog, func() {
		if err := srv.Close(); err != nil {
			cfg.Errlog.Printf("closing server: %s", err)
		}

		if err := cfg.DB.Close(); err != nil {
			cfg.Errlog.Printf("closing DB: %s", err)
		}
	})
}
