package main

import (
	"time"

	"github.com/xolra0d/alias-online/shared/pkg/config"
)

type ServerConfig struct {
	// HTTP
	RunningAddr     string        // Env name: `RUNNING_ADDR`. Address to run web on. Default: `:8060`.
	ShutdownTimeout time.Duration // Env name: `SHUTDOWN_TIMEOUT`. Time for transport server to shut down in seconds. Default: 10.

	// Workers
	PollInterval                 time.Duration // Env name: `WORKERS_POLL_INTERVAL`. Wait time between searches for new workers in seconds. Default: 10.
	WorkerExpiry                 time.Duration // Env name: `WORKER_EXPIRY`. Wait time before worker expires in seconds. Default: 30.
	RetrieveActiveWorkersTimeout time.Duration

	// DATABASE
	RedisAddr     string // Env name: `REDIS_ADDR`. Redis server address. Default: `localhost:6379`.
	RedisUsername string // Env name: `REDIS_USERNAME`. Redis auth username. Default: "".
	RedisPassword string // Env name: `REDIS_PASSWORD`. Redis auth password. Default: ".
	RedisDB       int    // Env name: `REDIS_DB`. Redis database index. Default: 0.
}

func LoadServerConfig() *ServerConfig {
	runningAddr := config.GetEnvOrFallback("RUNNING_ADDR", ":8060")
	shutdownTimeout := config.StringToSeconds("SHUTDOWN_TIMEOUT", config.GetEnvOrFallback("SHUTDOWN_TIMEOUT", "10"))

	pollInterval := config.StringToSeconds("WORKERS_POLL_INTERVAL", config.GetEnvOrFallback("WORKERS_POLL_INTERVAL", "10"))
	workerExpiry := config.StringToSeconds("WORKER_EXPIRY", config.GetEnvOrFallback("WORKER_EXPIRY", "30"))
	retrieveActiveWorkersTimeout := config.StringToSeconds("RETRIEVE_ACTIVE_WORKERS_TIMEOUT", config.GetEnvOrFallback("RETRIEVE_ACTIVE_WORKERS_TIMEOUT", "10"))

	redisAddr := config.GetEnvOrFallback("REDIS_ADDR", "localhost:6379")
	redisUsername := config.GetEnvOrFallback("REDIS_USERNAME", "")
	redisPassword := config.GetEnvOrFallback("REDIS_PASSWORD", "")
	redisDB := config.StringToUInt("REDIS_DB", config.GetEnvOrFallback("REDIS_DB", "0"))

	return &ServerConfig{
		RunningAddr:     runningAddr,
		ShutdownTimeout: shutdownTimeout,

		PollInterval:                 pollInterval,
		WorkerExpiry:                 workerExpiry,
		RetrieveActiveWorkersTimeout: retrieveActiveWorkersTimeout,

		RedisAddr:     redisAddr,
		RedisUsername: redisUsername,
		RedisPassword: redisPassword,
		RedisDB:       int(redisDB),
	}
}
