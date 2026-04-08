package main

import (
	"time"

	"github.com/xolra0d/alias-online/shared/pkg/config"
)

type ServerConfig struct {
	// DATABASES
	PostgresUrl              string        // Env name: `POSTGRES_URL`. PostgreSQL connection string. Default: none, will panic, if not set.
	LoadVocabTimeout         time.Duration // Env name: `LOAD_VOCAB_TIMEOUT`. Max wait time for loading vocab in seconds. Default: 10.
	ClosePostgresConnTimeout time.Duration // Env name: `CLOSE_POSTGRES_CONN_TIMEOUT`. Max wait time for closing postgres connection in seconds. Default: 10.

	// HTTP
	AllowedOrigins   string        // Env name: `ALLOWED_ORIGINS`. Origins to respond to (e.g., http://website.com:12), separated by comma. Default: none, will exit, if not set.
	RunningAddr      string        // Env name: `RUNNING_ADDR`. Address to bind HTTP server to. Default: `:8050`.
	WorkerPublicAddr string        // Env name: `WORKER_PUBLIC_ADDR`. Public worker address sent to room_manager and clients. Default: value of `RUNNING_ADDR`.
	ShutdownTimeout  time.Duration // Env name: `SHUTDOWN_TIMEOUT`. Time for transport server to shut down in seconds. Default: 10.

	// WORKERS
	WorkerPollInterval time.Duration // Env name: `WORKER_POLL_INTERVAL`. Wait time for worker pings in seconds. Default: 10.

	WsOriginPatterns     string        // Env name: `WS_ORIGIN_PATTERNS`. Defines allowed origins for ws connections, separated by comma. Default: none, will exit, if not set.
	MaxMessagesPerSecond int           // Env name: `MAX_MESSAGES_PER_SECOND`. Maximum messages sent per second, before connection is closed. Default: 50.
	WsPingTimeout        time.Duration // Env name: `PING_TIMEOUT`. Max wait time for ping request in seconds. Default: 5.
	WsWriteTimeout       time.Duration // Env name: `WS_WRITE_TIMEOUT`. Max wait time for writing response in seconds. Default: 5.
	LoadRoomTimeout      time.Duration // Env name: `LOAD_ROOM_TIMEOUT`. Max wait time for loading room_worker in seconds. Default: 10.
	SaveRoomTimeout      time.Duration // Env name: `SAVE_ROOM_TIMEOUT`. Max wait time for saving room_worker in seconds. Default: 10.
	MaxClockValue        int           // Env name: `MAX_CLOCK_VALUE`. Max clock value used for room state. Default 36000.
	JwtPublicKeyPath     string        // Env name: `JWT_PUBLIC_KEY_PATH`. Path to the JWT public key file. Default: none, will exit, if not set.
	RoomManagerUrl       string        // Env name: `ROOM_MANAGER_URL`. Address of the room manager service. Default: `localhost:8060`.
	VocabManagerUrl      string        // Env name: `VOCAB_MANAGER_URL`. Address of the vocab manager service. Default: `localhost:8070`.
}

func LoadServerConfig() *ServerConfig {
	postgresUrl := config.GetEnvOrExit("POSTGRES_URL")
	loadVocabTimeout := config.StringToSeconds("LOAD_VOCAB_TIMEOUT", config.GetEnvOrFallback("LOAD_VOCAB_TIMEOUT", "10"))
	closePostgresConnTimeout := config.StringToSeconds("CLOSE_POSTGRES_CONN_TIMEOUT", config.GetEnvOrFallback("CLOSE_POSTGRES_CONN_TIMEOUT", "10"))

	allowedOrigins := config.GetEnvOrExit("ALLOWED_ORIGINS")
	runningAddr := config.GetEnvOrFallback("RUNNING_ADDR", ":8050")
	workerPublicAddr := config.GetEnvOrFallback("WORKER_PUBLIC_ADDR", runningAddr)
	shutdownTimeout := config.StringToSeconds("SHUTDOWN_TIMEOUT", config.GetEnvOrFallback("SHUTDOWN_TIMEOUT", "10"))

	pollInterval := config.StringToSeconds("WORKER_POLL_INTERVAL", config.GetEnvOrFallback("WORKER_POLL_INTERVAL", "10"))
	wsOriginPatterns := config.GetEnvOrExit("WS_ORIGIN_PATTERNS")
	maxMessagesPerSecond := config.StringToUInt("MAX_MESSAGES_PER_SECOND", config.GetEnvOrFallback("MAX_MESSAGES_PER_SECOND", "50"))
	pingTimeout := config.StringToSeconds("PING_TIMEOUT", config.GetEnvOrFallback("PING_TIMEOUT", "5"))
	wsWriteTimeout := config.StringToSeconds("WS_WRITE_TIMEOUT", config.GetEnvOrFallback("WS_WRITE_TIMEOUT", "5"))
	loadRoomTimeout := config.StringToSeconds("LOAD_ROOM_TIMEOUT", config.GetEnvOrFallback("LOAD_ROOM_TIMEOUT", "10"))
	saveRoomTimeout := config.StringToSeconds("SAVE_ROOM_TIMEOUT", config.GetEnvOrFallback("SAVE_ROOM_TIMEOUT", "10"))

	jwtPublicKeyPath := config.GetEnvOrExit("JWT_PUBLIC_KEY_PATH")
	maxClockValue := config.StringToUInt("MAX_CLOCK_VALUE", config.GetEnvOrFallback("MAX_CLOCK_VALUE", "36000"))

	roomManagerUrl := config.GetEnvOrFallback("ROOM_MANAGER_URL", "localhost:8060")
	vocabManagerUrl := config.GetEnvOrFallback("VOCAB_MANAGER_URL", "localhost:8070")

	return &ServerConfig{
		PostgresUrl:              postgresUrl,
		LoadVocabTimeout:         loadVocabTimeout,
		ClosePostgresConnTimeout: closePostgresConnTimeout,

		AllowedOrigins:   allowedOrigins,
		RunningAddr:      runningAddr,
		WorkerPublicAddr: workerPublicAddr,
		ShutdownTimeout:  shutdownTimeout,

		WorkerPollInterval: pollInterval,

		WsOriginPatterns:     wsOriginPatterns,
		MaxMessagesPerSecond: int(maxMessagesPerSecond),
		WsPingTimeout:        pingTimeout,
		WsWriteTimeout:       wsWriteTimeout,
		LoadRoomTimeout:      loadRoomTimeout,
		SaveRoomTimeout:      saveRoomTimeout,
		JwtPublicKeyPath:     jwtPublicKeyPath,
		MaxClockValue:        int(maxClockValue),
		RoomManagerUrl:       roomManagerUrl,
		VocabManagerUrl:      vocabManagerUrl,
	}
}
