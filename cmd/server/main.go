package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/azezaz/HostMK-Security-Framework/internal/api"
	"github.com/azezaz/HostMK-Security-Framework/internal/api/handlers"
	"github.com/azezaz/HostMK-Security-Framework/internal/auth"
	"github.com/azezaz/HostMK-Security-Framework/internal/config"
	"github.com/azezaz/HostMK-Security-Framework/internal/crypto"
	"github.com/azezaz/HostMK-Security-Framework/internal/database"
	"github.com/azezaz/HostMK-Security-Framework/internal/observability"
	"github.com/azezaz/HostMK-Security-Framework/internal/persistence"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	logger := observability.NewLogger(cfg.App.Environment)
	ctx := context.Background()
	shutdownTracer, err := observability.InitTracer(ctx, cfg.Telemetry.ServiceName, cfg.Telemetry.OTLPEndpoint)
	if err != nil {
		log.Fatalf("telemetry init failed: %v", err)
	}
	defer func() { _ = shutdownTracer(context.Background()) }()

	db, err := database.New(database.Config{DatabaseURL: cfg.DB.URL})
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Migrate(ctx); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	aead, err := crypto.NewAEAD(cfg.Crypto.EncryptionKey)
	if err != nil {
		log.Fatalf("encryption init failed: %v", err)
	}

	appAPI := &handlers.API{
		Config: cfg,
		Store:  persistence.New(db.DB()),
		Tokens: auth.Tokens{
			Issuer:     cfg.JWT.Issuer,
			Secret:     []byte(cfg.JWT.Secret),
			AccessTTL:  cfg.JWT.AccessTTL,
			RefreshTTL: cfg.JWT.RefreshTTL,
		},
		AEAD: aead,
		Log:  logger,
	}

	server := &http.Server{
		Addr:              cfg.Address(),
		Handler:           api.NewRouter(appAPI),
		ReadTimeout:       cfg.App.ReadTimeout,
		WriteTimeout:      cfg.App.WriteTimeout,
		IdleTimeout:       cfg.App.IdleTimeout,
		ReadHeaderTimeout: cfg.App.ReadTimeout,
	}

	go func() {
		logger.Info("server_start", "addr", cfg.Address(), "name", cfg.App.Name, "version", cfg.App.Version)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.App.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
