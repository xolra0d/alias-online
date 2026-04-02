package main

import (
	"time"

	"github.com/xolra0d/alias-online/shared/pkg/config"
)

type ServerConfig struct {
	// DATABASES
	PostgresUrl       string        // Env name: `POSTGRES_URL`. PostgreSQL connection string. Default: none, will panic, if not set.
	LoadVocabsTimeout time.Duration // Env name: `LOAD_VOCABS_TIMEOUT`. Max wait time for loading vocabs in seconds. Default: 10.

	// HTTP
	AllowedOrigins  string        // Env name: `ALLOWED_ORIGINS`. Origins to respond to (e.g., http://website.com:12), separated by comma. Default: none, will exit, if not set.
	RunningAddr     string        // Env name: `RUNNING_ADDR`. Address to run web on. Default: `:8080`.
	ShutdownTimeout time.Duration // Env name: `SHUTDOWN_TIMEOUT`. Time for transport server to shut down in seconds. Default: 10.

	// Vocabs
	PollInterval time.Duration // Env name: `VOCABS_POLL_INTERVAL`. Wait time between checks for vocabs update in seconds. Default: 10.
}

func LoadServerConfig() *ServerConfig {
	postgresUrl := config.GetEnvOrExit("POSTGRES_URL")
	loadVocabsTimeout := config.StringToSeconds("LOAD_VOCABS_TIMEOUT", config.GetEnvOrFallback("LOAD_VOCABS_TIMEOUT", "10"))

	allowedOrigins := config.GetEnvOrExit("ALLOWED_ORIGINS")
	runningAddr := config.GetEnvOrFallback("RUNNING_ADDR", ":8080")
	shutdownTimeout := config.StringToSeconds("SHUTDOWN_TIMEOUT", config.GetEnvOrFallback("SHUTDOWN_TIMEOUT", "10"))

	pollInterval := config.StringToSeconds("POLL_INTERVAL", config.GetEnvOrFallback("POLL_INTERVAL", "10"))

	return &ServerConfig{
		postgresUrl,
		loadVocabsTimeout,

		allowedOrigins,
		runningAddr,
		shutdownTimeout,

		pollInterval,
	}
}
