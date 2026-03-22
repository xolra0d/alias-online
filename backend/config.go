package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ServerConfig holds all runtime configuration loaded from ENV variables.
type ServerConfig struct {
	// APP
	LogMessageMaxQueue int           // Env name: `LOG_MESSAGE_MAX_QUEUE`. Max queue length before client logger needs to wait. Default: 100.
	LoadVocabsTimeout  time.Duration // Env name: `LOAD_VOCABS_TIMEOUT`. Vocabularies load timeout in seconds. Default: 5.

	// DATABASES
	PostgresUrl string // Env name: `POSTGRES_URL`. PostgreSQL connection string. Default: none, will panic, if not set.

	// SECURITY
	Argon2idTime    uint32 // Env name: `ARGON2ID_TIME`. Number of iterations to perform. Default: 2.
	Argon2idMemory  uint32 // Env name: `ARGON2ID_MEMORY`. Amount of memory to use in bytes. Default: 65536.
	Argon2idThreads uint8  // Env name: `ARGON2ID_THREADS`. Degree of parallelism. Default: 1.
	Argon2idOutLen  uint32 // Env name: `ARGON2ID_OUT_LEN`. Desired number of returned bytes. Default: 32.

	// HTTP
	AllowedOrigins           string        // Env name: `ALLOWED_ORIGINS`. Origins to respond to (e.g., http://website.com:12), separetad by comma. Default: none, will exit, if not set.
	RunningAddr              string        // Env name: `RUNNING_ADDR`. Address to run web on. Default: `:8080`.
	ReadTimeout              time.Duration // Env name: `READ_TIMEOUT`. Docs: net/http.Server.ReadTimeout in seconds. Default: 5.
	WriteTimeout             time.Duration // Env name: `WRITE_TIMEOUT`. Docs: net/http.Server.WriteTimeout in seconds. Default: 5.
	IdleTimeout              time.Duration // Env name: `IDLE_TIMEOUT`. Docs: net/http.Server.IdleTimeout in seconds. Default: 30.
	ShutdownTimeout          time.Duration // Env name: `SHUTDOWN_TIMEOUT`. Time for http server to shut down in seconds. Default: 10.
	CreateUserLimitPerWindow int           // Env name: `CREATE_USER_LIMIT_PER_WINDOW`. Number of users allowed to be created for `LimiterWindow` time. Default: 30.
	CreateRoomLimitPerWindow int           // Env name: `CREATE_ROOM_LIMIT_PER_WINDOW`. Number of rooms allowed to be created for `LimiterWindow` time. Default: 30.
	LimiterCleanupEvery      int           // Env name: `LIMITER_CLEANUP_EVERY`. Removes outdated entries after `LimiterCleanupEvery` requests handled. Default: 100.
	LimiterWindow            time.Duration // Env name: `LIMITER_WINDOW`. Defines window for limiters in seconds. Default: 60.
	CreateUserTimeout        time.Duration // Env name: `CREATE_USER_TIMEOUT`. Defines timeout for user creation in seconds. Default: 5.
	CreateRoomTimeout        time.Duration // Env name: `CREATE_ROOM_TIMEOUT`. Defines timeout for room creation in seconds. Default: 5.

	// ROOMS
	MinClock                     int           // Env name: `MIN_CLOCK`. Min number of seconds for clock in round. Default: 1.
	MaxClock                     int           // Env name: `MAX_CLOCK`. Max number of seconds for clock in round. Default: 36000. (10 hours)
	MaxAdditionalVocabularyWords int           // Env name: `MAX_ADDITIONAL_VOCABULARY_WORDS`. Max number of words in additional vocabulary. Default: 1000.
	MaxAdditionalWordLength      int           // Env name: `MAX_ADDITIONAL_WORD_LENGTH`. Max number of runes (UTF-8 chars) in word in additional vocabulary. Default 64.
	LoadRoomTimeout              time.Duration // Env name: `LOAD_ROOM_TIMEOUT`. Max wait time for loading room in seconds. Default: 10.
	SaveRoomTimeout              time.Duration // Env name: `SAVE_ROOM_TIMEOUT`. Max wait time for saving room in seconds. Default: 10.

	// WS
	WSOriginPatterns     string        // Env name: `WS_ORIGIN_PATTERNS`. Defines allowed origins for ws connections, separetad by comma. Default: none, will exit, if not set.
	MaxMessagesPerSecond int           // Env name: `MAX_MESSAGES_PER_SECOND`. Maximum messages sent per second, before connection is closed. Default: 50.
	PingTimeout          time.Duration // Env name: `PING_TIMEOUT`. Max wait time for ping request in seconds. Default: 5.
	WSWriteTimeout       time.Duration // Env name: `WS_WRITE_TIMEOUT`. Max wait time for writing response in seconds. Default: 5.
}

// Loads server config from ENV variables.
func loadServerConfig() ServerConfig {
	logMessageMaxQueue := stringToUInt("LOG_MESSAGE_MAX_QUEUE", getEnvOrFallback("LOG_MESSAGE_MAX_QUEUE", "100"))
	loadVocabsTimeout := stringToSeconds("LOAD_VOCABS_TIMEOUT", getEnvOrFallback("LOAD_VOCABS_TIMEOUT", "5"))

	postgresUrl := getEnvOrExit("POSTGRES_URL")

	agon2idTime := stringToUInt("ARGON2ID_TIME", getEnvOrFallback("ARGON2ID_TIME", "2"))
	agon2idMemory := stringToUInt("ARGON2ID_MEMORY", getEnvOrFallback("ARGON2ID_MEMORY", "65536"))
	agon2idThreads := stringToUInt("ARGON2ID_THREADS", getEnvOrFallback("ARGON2ID_THREADS", "1"))
	argon2idOutLen := stringToUInt("ARGON2ID_OUT_LEN", getEnvOrFallback("ARGON2ID_OUT_LEN", "32"))

	allowedOrigins := getEnvOrExit("ALLOWED_ORIGINS")
	runningAddr := getEnvOrFallback("RUNNING_ADDR", ":8080")
	readTimeout := stringToSeconds("READ_TIMEOUT", getEnvOrFallback("READ_TIMEOUT", "5"))
	writeTimeout := stringToSeconds("WRITE_TIMEOUT", getEnvOrFallback("WRITE_TIMEOUT", "5"))
	idleTimeout := stringToSeconds("IDLE_TIMEOUT", getEnvOrFallback("IDLE_TIMEOUT", "30"))
	shutdownTimeout := stringToSeconds("SHUTDOWN_TIMEOUT", getEnvOrFallback("SHUTDOWN_TIMEOUT", "10"))
	createUserLimitPerWindow := stringToUInt("CREATE_USER_LIMIT_PER_WINDOW", getEnvOrFallback("CREATE_USER_LIMIT_PER_WINDOW", "30"))
	createRoomLimitPerWindow := stringToUInt("CREATE_ROOM_LIMIT_PER_WINDOW", getEnvOrFallback("CREATE_ROOM_LIMIT_PER_WINDOW", "30"))
	limiterCleanupEvery := stringToUInt("LIMITER_CLEANUP_EVERY", getEnvOrFallback("LIMITER_CLEANUP_EVERY", "100"))
	limiterWindow := stringToSeconds("LIMITER_WINDOW", getEnvOrFallback("LIMITER_WINDOW", "60"))
	createUserTimeout := stringToSeconds("CREATE_USER_TIMEOUT", getEnvOrFallback("CREATE_USER_TIMEOUT", "5"))
	createRoomTimeout := stringToSeconds("CREATE_ROOM_TIMEOUT", getEnvOrFallback("CREATE_ROOM_TIMEOUT", "5"))

	minClock := stringToUInt("MIN_CLOCK", getEnvOrFallback("MIN_CLOCK", "1"))
	maxClock := stringToUInt("MAX_CLOCK", getEnvOrFallback("MAX_CLOCK", "36000"))
	maxAdditionalVocabularyWords := stringToUInt("MAX_ADDITIONAL_VOCABULARY_WORDS", getEnvOrFallback("MAX_ADDITIONAL_VOCABULARY_WORDS", "1000"))
	maxAdditionalWordLength := stringToUInt("MAX_ADDITIONAL_WORD_LENGTH", getEnvOrFallback("MAX_ADDITIONAL_WORD_LENGTH", "64"))
	loadRoomTimeout := stringToSeconds("LOAD_ROOM_TIMEOUT", getEnvOrFallback("LOAD_ROOM_TIMEOUT", "10"))
	saveRoomTimeout := stringToSeconds("SAVE_ROOM_TIMEOUT", getEnvOrFallback("SAVE_ROOM_TIMEOUT", "10"))

	wsOriginPatterns := getEnvOrExit("WS_ORIGIN_PATTERNS")
	maxMessagesPerSecond := stringToUInt("MAX_MESSAGES_PER_SECOND", getEnvOrFallback("MAX_MESSAGES_PER_SECOND", "50"))
	pingTimeout := stringToSeconds("PING_TIMEOUT", getEnvOrFallback("PING_TIMEOUT", "5"))
	wsWriteTimeout := stringToSeconds("WS_WRITE_TIMEOUT", getEnvOrFallback("WS_WRITE_TIMEOUT", "5"))

	return ServerConfig{
		LogMessageMaxQueue: int(logMessageMaxQueue),
		LoadVocabsTimeout:  loadVocabsTimeout,

		PostgresUrl: postgresUrl,

		Argon2idTime:    uint32(agon2idTime),
		Argon2idMemory:  uint32(agon2idMemory),
		Argon2idThreads: uint8(agon2idThreads),
		Argon2idOutLen:  uint32(argon2idOutLen),

		AllowedOrigins:           allowedOrigins,
		RunningAddr:              runningAddr,
		ReadTimeout:              readTimeout,
		WriteTimeout:             writeTimeout,
		IdleTimeout:              idleTimeout,
		ShutdownTimeout:          shutdownTimeout * time.Second,
		CreateUserLimitPerWindow: int(createUserLimitPerWindow),
		CreateRoomLimitPerWindow: int(createRoomLimitPerWindow),
		LimiterCleanupEvery:      int(limiterCleanupEvery),
		LimiterWindow:            limiterWindow,
		CreateUserTimeout:        createUserTimeout,
		CreateRoomTimeout:        createRoomTimeout,

		MinClock:                     int(minClock),
		MaxClock:                     int(maxClock),
		MaxAdditionalVocabularyWords: int(maxAdditionalVocabularyWords),
		MaxAdditionalWordLength:      int(maxAdditionalWordLength),
		LoadRoomTimeout:              loadRoomTimeout,
		SaveRoomTimeout:              saveRoomTimeout,

		WSOriginPatterns:     wsOriginPatterns,
		MaxMessagesPerSecond: int(maxMessagesPerSecond),
		PingTimeout:          pingTimeout,
		WSWriteTimeout:       wsWriteTimeout,
	}
}

// getEnvOrExit tries to get `name` env var. If not set - exits.
func getEnvOrExit(name string) string {
	out := strings.TrimSpace(os.Getenv(name))
	if out == "" {
		fmt.Printf("ENV: `%s` is not specified!\n", name)
		os.Exit(1)
	}
	return out
}

// getEnvOrFallback tries to get `name` env var and return it. If not set - returns `fallback`.
func getEnvOrFallback(name, fallback string) string {
	out := strings.TrimSpace(os.Getenv(name))
	if out == "" {
		return fallback
	}
	return out
}

// stringToSeconds parses `val` (unsigned num) into num seconds. Uses `name` for panic.
func stringToSeconds(name, val string) time.Duration {
	out, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		panic(fmt.Sprintf("ENV: `%s`=%s is invalid: %s!", name, val, err))
	}

	return time.Duration(out) * time.Second
}

// stringToSeconds parses `val` (unsigned num) into uint64. Uses `name` for panic.
func stringToUInt(name, val string) uint64 {
	out, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("ENV: `%s`=%s is invalid: %s!", name, val, err))
	}

	return out
}
