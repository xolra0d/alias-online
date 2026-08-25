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
	"github.com/xolra0d/alias-online/internal/config"
	"github.com/xolra0d/alias-online/internal/database"
	"github.com/xolra0d/alias-online/internal/room"
	"github.com/xolra0d/alias-online/internal/transport"
)

func main() {
	const op = "server.main"

	serverConfig := config.LoadServerConfig()

	logger := config.NewLogger(serverConfig.LogMessageMaxQueue, os.Stdout)
	go logger.StartLogging()
	defer logger.EndLogging()

	pgPool, err := database.InitPool(serverConfig.PostgresUrl)
	if err != nil {
		logger.Error(op, "Failed to connect to database", "error", err)
		return
	}
	defer pgPool.Close()
	secrets := database.NewSecrets(
		logger,
		serverConfig.Argon2idTime,
		serverConfig.Argon2idMemory,
		serverConfig.Argon2idThreads,
		serverConfig.Argon2idOutLen,
	)
	postgres := database.NewPostgres(pgPool, secrets, logger)

	ctx, cancel := context.WithTimeout(context.Background(), serverConfig.LoadVocabsTimeout)
	vocabs, err := postgres.LoadVocabs(ctx)
	cancel()
	if err != nil {
		logger.Error(op, "could not load vocabs", "error", err)
		return
	}

	rooms := transport.NewRooms(
		map[string]*room.Room{},
		logger,
		serverConfig.MinClock,
		serverConfig.MaxClock,
		serverConfig.MaxAdditionalVocabularyWords,
		serverConfig.MaxAdditionalWordLength,
		strings.Split(serverConfig.WSOriginPatterns, ","),
		serverConfig.MaxMessagesPerSecond,
		serverConfig.PingTimeout,
		serverConfig.WSWriteTimeout,
		serverConfig.LoadRoomTimeout,
		serverConfig.SaveRoomTimeout,
	)

	handles := transport.NewHandles(
		postgres,
		rooms,
		room.NewVocabularies(vocabs),
		logger,
		serverConfig.CreateUserTimeout,
		serverConfig.CreateRoomTimeout,
	)
	corsRules := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(serverConfig.AllowedOrigins, ","),
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Length", "Content-Type", "User-Id", "User-Secret"},
		AllowCredentials: true,
	})

	createUserLimiter := transport.NewRateLimiter(serverConfig.CreateUserLimitPerWindow, serverConfig.LimiterWindow, serverConfig.LimiterCleanupEvery)
	createRoomLimiter := transport.NewRateLimiter(serverConfig.CreateRoomLimitPerWindow, serverConfig.LimiterWindow, serverConfig.LimiterCleanupEvery)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/ok", handles.Healthy)
	mux.HandleFunc("GET /api/available-vocabs", handles.AvailableLanguages)
	mux.Handle("POST /api/create-user", transport.Chain(
		http.HandlerFunc(handles.CreateUser),
		handles.IpRateLimiter(createUserLimiter),
	))
	mux.HandleFunc("GET /api/ws/{roomId}", handles.InitWS)

	mux.Handle("GET /api/protected/ok", transport.Chain(
		http.HandlerFunc(handles.Healthy),
		handles.Auth(),
	))
	mux.Handle("POST /api/protected/create-room", transport.Chain(
		http.HandlerFunc(handles.CreateRoom),
		handles.Auth(),
		handles.UserIdRateLimiter(createRoomLimiter),
	))

	server := &http.Server{
		Addr: serverConfig.RunningAddr,
		Handler: transport.Chain(
			mux,
			corsRules.Handler,
			transport.Logging(logger),
		),
		ReadTimeout:  serverConfig.ReadTimeout,
		WriteTimeout: serverConfig.WriteTimeout,
		IdleTimeout:  serverConfig.IdleTimeout,
	}

	go func() {
		logger.Info(op, "server started", "addr", serverConfig.RunningAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(op, "listen error", "err", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown
	logger.Info(op, "shutdown initiated")

	ctx, cancel = context.WithTimeout(context.Background(), serverConfig.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(op, "shutdown failed", "err", err)
	}

	logger.Info(op, "server stopped")
}
