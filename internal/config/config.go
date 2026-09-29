package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	DefaultHost            = "0.0.0.0"
	DefaultPort            = "8080"
	DefaultEnvironment     = "development"
	DefaultReadTimeout     = 15 * time.Second
	DefaultWriteTimeout    = 15 * time.Second
	DefaultIdleTimeout     = 60 * time.Second
	DefaultShutdownTimeout = 10 * time.Second
	DefaultAccessTTL       = 15 * time.Minute
	DefaultRefreshTTL      = 720 * time.Hour
)

type Config struct {
	App       AppConfig
	DB        DatabaseConfig
	JWT       JWTConfig
	Crypto    CryptoConfig
	Telemetry TelemetryConfig
}

type AppConfig struct {
	Name            string
	Version         string
	Environment     string
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	URL string
}

type JWTConfig struct {
	Secret     string
	Issuer     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type CryptoConfig struct {
	EncryptionKey string
}

type TelemetryConfig struct {
	OTLPEndpoint string
	ServiceName  string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		App: AppConfig{
			Name:            getEnv("APP_NAME", "HostMK Security Framework"),
			Version:         getEnv("APP_VERSION", "0.1.0"),
			Environment:     getEnv("APP_ENV", DefaultEnvironment),
			Host:            getEnv("HOST", DefaultHost),
			Port:            getEnv("PORT", DefaultPort),
			ReadTimeout:     getDurationEnv("READ_TIMEOUT_SECONDS", DefaultReadTimeout),
			WriteTimeout:    getDurationEnv("WRITE_TIMEOUT_SECONDS", DefaultWriteTimeout),
			IdleTimeout:     getDurationEnv("IDLE_TIMEOUT_SECONDS", DefaultIdleTimeout),
			ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT_SECONDS", DefaultShutdownTimeout),
		},
		DB: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", ""),
			Issuer:     getEnv("JWT_ISSUER", "hostmk-security"),
			AccessTTL:  getDurationEnv("JWT_ACCESS_TTL_SECONDS", DefaultAccessTTL),
			RefreshTTL: getDurationEnv("JWT_REFRESH_TTL_SECONDS", DefaultRefreshTTL),
		},
		Crypto: CryptoConfig{
			EncryptionKey: getEnv("ENCRYPTION_KEY", ""),
		},
		Telemetry: TelemetryConfig{
			OTLPEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
			ServiceName:  getEnv("OTEL_SERVICE_NAME", "hostmk-security"),
		},
	}
}

func (c *Config) Validate() error {
	if c.App.Port == "" {
		return fmt.Errorf("PORT is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.DB.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.Crypto.EncryptionKey == "" {
		return fmt.Errorf("ENCRYPTION_KEY is required")
	}
	return nil
}

func (c *Config) Address() string {
	return fmt.Sprintf("%s:%s", c.App.Host, c.App.Port)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		return fallback
	}
	return time.Duration(seconds) * time.Second
}
