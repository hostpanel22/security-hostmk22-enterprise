package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/azezaz/HostMK-Security-Framework/internal/api"
	"github.com/azezaz/HostMK-Security-Framework/internal/config"
	"github.com/azezaz/HostMK-Security-Framework/internal/database"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	cfg := config.Load()

	if err := cfg.Validate(); err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	log.Printf(
		"Starting %s v%s (%s)",
		cfg.App.Name,
		cfg.App.Version,
		cfg.App.Environment,
	)

	var db *database.Postgres

	if cfg.DB.URL != "" {

		var err error

		db, err = database.New(database.Config{
			DatabaseURL: cfg.DB.URL,
		})

		if err != nil {
			log.Fatalf(
				"database connection failed: %v",
				err,
			)
		}

		defer func() {
			if err := db.Close(); err != nil {
				log.Printf(
					"database close error: %v",
					err,
				)
			}
		}()

		log.Println(
			"PostgreSQL connection established",
		)

	} else {

		log.Println(
			"DATABASE_URL not configured, running without database",
		)
	}

	router := api.NewRouter()

	server := &http.Server{
		Addr:              cfg.Address(),
		Handler:           router,
		ReadTimeout:       cfg.App.ReadTimeout,
		WriteTimeout:      cfg.App.WriteTimeout,
		IdleTimeout:       cfg.App.IdleTimeout,
		ReadHeaderTimeout: cfg.App.ReadTimeout,
	}

	go func() {

		log.Printf(
			"%s listening on %s",
			cfg.App.Name,
			cfg.Address(),
		)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			log.Fatalf(
				"server error: %v",
				err,
			)
		}
	}()

	waitForShutdown(
		server,
		cfg.App.ShutdownTimeout,
	)
}

func waitForShutdown(
	server *http.Server,
	timeout time.Duration,
) {

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	log.Println(
		"shutdown signal received",
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		timeout,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {

		log.Printf(
			"shutdown error: %v",
			err,
		)
	}

	log.Println(
		"server stopped",
	)
}
