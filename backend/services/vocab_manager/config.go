package main

import (
	"time"

	"github.com/xolra0d/alias-online/shared/pkg/config"
)

type ServerConfig struct {
	// DATABASES
	PostgresUrl              string        // Env name: `POSTGRES_URL`. PostgreSQL connection string. Default: none, will panic, if not set.
	LoadVocabsTimeout        time.Duration // Env name: `LOAD_VOCABS_TIMEOUT`. Max wait time for loading vocab_manager in seconds. Default: 10.
	ClosePostgresConnTimeout time.Duration // Env name: `CLOSE_POSTGRES_CONN_TIMEOUT`. Max wait time for closing postgres connection in seconds. Default: 10.

	// HTTP
	RunningAddr     string        // Env name: `RUNNING_ADDR`. Address to run web on. Default: `:8070`.
	ShutdownTimeout time.Duration // Env name: `SHUTDOWN_TIMEOUT`. Time for transport server to shut down in seconds. Default: 10.

	// Vocabs
	VocabsPollInterval time.Duration // Env name: `VOCABS_POLL_INTERVAL`. Wait time between checks for vocab_manager update in seconds. Default: 10.
}

func LoadServerConfig() *ServerConfig {
	postgresUrl := config.GetEnvOrExit("POSTGRES_URL")
	loadVocabsTimeout := config.StringToSeconds("LOAD_VOCABS_TIMEOUT", config.GetEnvOrFallback("LOAD_VOCABS_TIMEOUT", "10"))
	closePostgresConnTimeout := config.StringToSeconds("CLOSE_POSTGRES_CONN_TIMEOUT", config.GetEnvOrFallback("CLOSE_POSTGRES_CONN_TIMEOUT", "10"))

	runningAddr := config.GetEnvOrFallback("RUNNING_ADDR", ":8070")
	shutdownTimeout := config.StringToSeconds("SHUTDOWN_TIMEOUT", config.GetEnvOrFallback("SHUTDOWN_TIMEOUT", "10"))

	vocabsPollInterval := config.StringToSeconds("VOCABS_POLL_INTERVAL", config.GetEnvOrFallback("VOCABS_POLL_INTERVAL", "10"))

	return &ServerConfig{
		postgresUrl,
		loadVocabsTimeout,
		closePostgresConnTimeout,

		runningAddr,
		shutdownTimeout,

		vocabsPollInterval,
	}
}
