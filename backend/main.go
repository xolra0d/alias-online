package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rs/cors"
)

func main() {
	serverConfig := loadServerConfig()

	baseLogger := NewBaseLogger(serverConfig.LogMessageMaxQueue)
	go baseLogger.StartLogging()
	defer baseLogger.EndLogging()

	startupLogger := baseLogger.WithPrefix("STARTUP")

	pgPool, err := InitPool(serverConfig.PostgresUrl)
	if err != nil {
		startupLogger.Error(err.Error())
		return
	}
	defer pgPool.Close()
	postgres := &Postgres{
		pgPool,
		&Secrets{
			baseLogger.WithPrefix("SECRETS"),
			serverConfig.Argon2idTime,
			serverConfig.Argon2idMemory,
			serverConfig.Argon2idThreads,
			serverConfig.Argon2idOutLen,
		},
		baseLogger.WithPrefix("POSTGRES"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), serverConfig.LoadVocabsTimeout)
	vocabs, err := postgres.LoadVocabs(ctx)
	cancel()
	if err != nil {
		startupLogger.Error(err.Error())
		return
	}

	rooms := &Rooms{
		rooms:  map[string]*Room{},
		logger: baseLogger.WithPrefix("ROOMS"),

		MinClock:                     serverConfig.MinClock,
		MaxClock:                     serverConfig.MaxClock,
		MaxAdditionalVocabularyWords: serverConfig.MaxAdditionalVocabularyWords,
		MaxAdditionalWordLength:      serverConfig.MaxAdditionalWordLength,
		WSOriginPatterns:             strings.Split(serverConfig.WSOriginPatterns, ","),
		MaxMessagesPerSecond:         serverConfig.MaxMessagesPerSecond,
		PingTimeout:                  serverConfig.PingTimeout,
		WSWriteTimeout:               serverConfig.WSWriteTimeout,
		WSReadTimeout:                serverConfig.WSReadTimeout,
		LoadRoomTimeout:              serverConfig.LoadRoomTimeout,
		SaveRoomTimeout:              serverConfig.SaveRoomTimeout,
	}
	handles := &Handles{
		postgres,
		rooms,
		&Vocabularies{vocabulary: vocabs},
		baseLogger.WithPrefix("HANDLES"),

		serverConfig.CreateUserTimeout,
		serverConfig.CreateRoomTimeout,
	}
	corsRules := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(serverConfig.AllowedOrigins, ","),
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Length", "Content-Type", "User-Id", "User-Secret"},
		AllowCredentials: true,
	})

	createUserLimiter := NewRateLimiter(serverConfig.CreateUserLimitPerWindow, serverConfig.LimiterWindow, serverConfig.LimiterCleanupEvery)
	createRoomLimiter := NewRateLimiter(serverConfig.CreateRoomLimitPerWindow, serverConfig.LimiterWindow, serverConfig.LimiterCleanupEvery)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/ok", handles.Healthy)
	mux.HandleFunc("GET /api/available-vocabs", handles.AvailableLanguages)
	mux.Handle("POST /api/create-user", Chain(
		http.HandlerFunc(handles.CreateUser),
		handles.IpRateLimiter(createUserLimiter),
	))
	mux.HandleFunc("GET /api/ws/{roomId}", handles.InitWS)

	mux.Handle("GET /api/protected/ok", Chain(
		http.HandlerFunc(handles.Healthy),
		handles.Auth(),
	))
	mux.Handle("POST /api/protected/create-room", Chain(
		http.HandlerFunc(handles.CreateRoom),
		handles.Auth(),
		handles.UserIdRateLimiter(createRoomLimiter),
	))

	httpLogger := baseLogger.WithPrefix("HTTP")
	server := &http.Server{
		Addr: serverConfig.RunningAddr,
		Handler: Chain(
			mux,
			corsRules.Handler,
			Logging(httpLogger),
		),
		ReadTimeout:  serverConfig.ReadTimeout,
		WriteTimeout: serverConfig.WriteTimeout,
		IdleTimeout:  serverConfig.IdleTimeout,
	}

	go func() {
		startupLogger.Info("server started", "addr", serverConfig.RunningAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			startupLogger.Error("listen error", "err", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown
	startupLogger.Info("shutdown initiated")

	ctx, cancel = context.WithTimeout(context.Background(), serverConfig.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		startupLogger.Error("shutdown failed", "err", err)
	}

	startupLogger.Info("server stopped")
}
