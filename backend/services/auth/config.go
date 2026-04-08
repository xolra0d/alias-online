package main

import (
	"time"

	"github.com/xolra0d/alias-online/shared/pkg/config"
)

type ServerConfig struct {
	// DATABASES
	PostgresUrl        string        // Env name: `POSTGRES_URL`. PostgreSQL connection string. Default: none, will exit, if not set.
	AddAccountTimeout  time.Duration // Env name: `ADD_ACCOUNT_TIMEOUT`. Max wait time for saving new account in seconds. Default: 10.
	FindAccountTimeout time.Duration // Env name: `FIND_ACCOUNT_TIMEOUT`. Max wait time for finding account in seconds. Default: 10.

	// GRPC SERVER
	RunningAddr     string        // Env name: `RUNNING_ADDR`. Address to run web on. Default: `:8090`.
	ShutdownTimeout time.Duration // Env name: `SHUTDOWN_TIMEOUT`. Time for transport server to shut down in seconds. Default: 10.

	// SECURITY
	Argon2idTime          uint32        // Env name: `ARGON2ID_TIME`. Number of iterations to perform. Default: 2.
	Argon2idMemory        uint32        // Env name: `ARGON2ID_MEMORY`. Amount of memory to use in bytes. Default: 65536.
	Argon2idThreads       uint8         // Env name: `ARGON2ID_THREADS`. Degree of parallelism. Default: 1.
	Argon2idOutLen        uint32        // Env name: `ARGON2ID_OUT_LEN`. Desired number of returned bytes. Default: 32.
	JwtPrivateKeyFilename string        // Env name: `JWT_PRIVATE_KEY_FILENAME`. Path to private key used for creating JWT tokens. Default: none, will exit, if not set.
	JWTCookieTimeout      time.Duration // Env name: `JWT_COOKIE_TIMEOUT`. Time for JWT to expire in seconds. Default: 3600.
}

func LoadServerConfig() *ServerConfig {
	postgresUrl := config.GetEnvOrExit("POSTGRES_URL")
	addAccountTimeout := config.StringToSeconds("ADD_ACCOUNT_TIMEOUT", config.GetEnvOrFallback("ADD_ACCOUNT_TIMEOUT", "10"))
	findAccountTimeout := config.StringToSeconds("FIND_ACCOUNT_TIMEOUT", config.GetEnvOrFallback("FIND_ACCOUNT_TIMEOUT", "10"))

	runningAddr := config.GetEnvOrFallback("RUNNING_ADDR", ":8090")
	shutdownTimeout := config.StringToSeconds("SHUTDOWN_TIMEOUT", config.GetEnvOrFallback("SHUTDOWN_TIMEOUT", "10"))

	agon2idTime := config.StringToUInt("ARGON2ID_TIME", config.GetEnvOrFallback("ARGON2ID_TIME", "2"))
	agon2idMemory := config.StringToUInt("ARGON2ID_MEMORY", config.GetEnvOrFallback("ARGON2ID_MEMORY", "65536"))
	agon2idThreads := config.StringToUInt("ARGON2ID_THREADS", config.GetEnvOrFallback("ARGON2ID_THREADS", "1"))
	argon2idOutLen := config.StringToUInt("ARGON2ID_OUT_LEN", config.GetEnvOrFallback("ARGON2ID_OUT_LEN", "32"))
	jwtPrivateKeyFilename := config.GetEnvOrExit("RSA_PRIVATE_KEY_FILENAME")
	jwtCookieTimeout := config.StringToSeconds("JWT_COOKIE_TIMEOUT", config.GetEnvOrFallback("JWT_COOKIE_TIMEOUT", "3600"))

	return &ServerConfig{
		PostgresUrl:        postgresUrl,
		AddAccountTimeout:  addAccountTimeout,
		FindAccountTimeout: findAccountTimeout,

		RunningAddr:     runningAddr,
		ShutdownTimeout: shutdownTimeout,

		Argon2idTime:          uint32(agon2idTime),
		Argon2idMemory:        uint32(agon2idMemory),
		Argon2idThreads:       uint8(agon2idThreads),
		Argon2idOutLen:        uint32(argon2idOutLen),
		JwtPrivateKeyFilename: jwtPrivateKeyFilename,
		JWTCookieTimeout:      jwtCookieTimeout,
	}
}
