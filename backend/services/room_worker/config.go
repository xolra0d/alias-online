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
	AllowedOrigins  string        // Env name: `ALLOWED_ORIGINS`. Origins to respond to (e.g., http://website.com:12), separated by comma. Default: none, will exit, if not set.
	RunningAddr     string        // Env name: `RUNNING_ADDR`. Address to run web on. Default: `localhost:8050`.
	ShutdownTimeout time.Duration // Env name: `SHUTDOWN_TIMEOUT`. Time for transport server to shut down in seconds. Default: 10.

	// Vocabs
	PollInterval time.Duration // Env name: `VOCABS_POLL_INTERVAL`. Wait time between checks for vocab_manager update in seconds. Default: 10.

	WsOriginPatterns     string        // Env name: `WS_ORIGIN_PATTERNS`. Defines allowed origins for ws connections, separetad by comma. Default: none, will exit, if not set.
	MaxMessagesPerSecond int           // Env name: `MAX_MESSAGES_PER_SECOND`. Maximum messages sent per second, before connection is closed. Default: 50.
	WsPingTimeout        time.Duration // Env name: `PING_TIMEOUT`. Max wait time for ping request in seconds. Default: 5.
	WsWriteTimeout       time.Duration // Env name: `WS_WRITE_TIMEOUT`. Max wait time for writing response in seconds. Default: 5.
	LoadRoomTimeout      time.Duration // Env name: `LOAD_ROOM_TIMEOUT`. Max wait time for loading room_worker in seconds. Default: 10.
	SaveRoomTimeout      time.Duration // Env name: `SAVE_ROOM_TIMEOUT`. Max wait time for saving room_worker in seconds. Default: 10.

	JwtPublicKeyPath string
}

func LoadServerConfig() *ServerConfig {
	postgresUrl := config.GetEnvOrExit("POSTGRES_URL")
	loadVocabsTimeout := config.StringToSeconds("LOAD_VOCABS_TIMEOUT", config.GetEnvOrFallback("LOAD_VOCABS_TIMEOUT", "10"))
	closePostgresConnTimeout := config.StringToSeconds("CLOSE_POSTGRES_CONN_TIMEOUT", config.GetEnvOrFallback("CLOSE_POSTGRES_CONN_TIMEOUT", "10"))

	allowedOrigins := config.GetEnvOrExit("ALLOWED_ORIGINS")
	runningAddr := config.GetEnvOrFallback("RUNNING_ADDR", "localhost:8050")
	shutdownTimeout := config.StringToSeconds("SHUTDOWN_TIMEOUT", config.GetEnvOrFallback("SHUTDOWN_TIMEOUT", "10"))

	pollInterval := config.StringToSeconds("VOCABS_POLL_INTERVAL", config.GetEnvOrFallback("VOCABS_POLL_INTERVAL", "10"))

	wsOriginPatterns := config.GetEnvOrExit("WS_ORIGIN_PATTERNS")
	maxMessagesPerSecond := config.StringToUInt("MAX_MESSAGES_PER_SECOND", config.GetEnvOrFallback("MAX_MESSAGES_PER_SECOND", "50"))
	pingTimeout := config.StringToSeconds("PING_TIMEOUT", config.GetEnvOrFallback("PING_TIMEOUT", "5"))
	wsWriteTimeout := config.StringToSeconds("WS_WRITE_TIMEOUT", config.GetEnvOrFallback("WS_WRITE_TIMEOUT", "5"))
	loadRoomTimeout := config.StringToSeconds("LOAD_ROOM_TIMEOUT", config.GetEnvOrFallback("LOAD_ROOM_TIMEOUT", "10"))
	saveRoomTimeout := config.StringToSeconds("SAVE_ROOM_TIMEOUT", config.GetEnvOrFallback("SAVE_ROOM_TIMEOUT", "10"))

	jwtPublicKeyPath := config.GetEnvOrExit("JWT_PUBLIC_KEY_PATH")

	return &ServerConfig{
		postgresUrl,
		loadVocabsTimeout,
		closePostgresConnTimeout,

		allowedOrigins,
		runningAddr,
		shutdownTimeout,

		pollInterval,

		wsOriginPatterns,
		int(maxMessagesPerSecond),
		pingTimeout,
		wsWriteTimeout,
		loadRoomTimeout,
		saveRoomTimeout,

		jwtPublicKeyPath,
	}
}
