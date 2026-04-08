package main

import (
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/xolra0d/alias-online/shared/pkg/logger"
)

func main() {
	l := slog.New(logger.NewHandler(nil))

	l.Info("loading configuration")
	cfg := LoadServerConfig()

	l.Info("trying to connect to postgres")
	postgres, err := NewPostgres(cfg.PostgresUrl, l)
	if err != nil {
		l.Error("failed to connect to database", "error", err)
		return
	}

	l.Info("initializing secrets")
	secrets, err := NewSecrets(cfg.JwtPublicKeyPath, l)
	if err != nil {
		l.Error("failed to load secrets", "error", err)
		return
	}

	roomManagerClient, closeRoomManagerFunc, err := NewRoomManagerClient(cfg.RoomManagerUrl, l)
	if err != nil {
		l.Error("failed to start room manager client", "error", err)
		return
	}
	defer closeRoomManagerFunc()
	vocabManagerClient, closeVocabManagerFunc, err := NewVocabManagerClient(cfg.VocabManagerUrl, l)
	if err != nil {
		l.Error("failed to start vocab manager client", "error", err)
		return
	}
	defer closeVocabManagerFunc()

	rooms := NewRooms(
		postgres,
		l,
		roomManagerClient,
		vocabManagerClient,

		cfg.WorkerPublicAddr,
		strings.Split(cfg.WsOriginPatterns, ","),
		cfg.LoadRoomTimeout,
		cfg.SaveRoomTimeout,
		cfg.WsWriteTimeout,
		cfg.WsPingTimeout,
		cfg.MaxMessagesPerSecond,
		cfg.MaxClockValue,
		cfg.LoadVocabTimeout,
	)
	handles := NewHandles(secrets, l, rooms)

	shouldStop := make(chan struct{})
	done := make(chan struct{})

	go rooms.RunPinger(l, cfg.WorkerPollInterval, cfg.WorkerPublicAddr, shouldStop, done)
	go RunHttpClient(handles, secrets, l, cfg.RunningAddr, cfg.ShutdownTimeout, shouldStop, done)

	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown
	l.Info("shutdown initiated")

	shouldStop <- struct{}{}
	shouldStop <- struct{}{}

	go func() {
		time.Sleep(cfg.ShutdownTimeout)
		l.Warn("timeout shutting down")
		os.Exit(1)
	}()

	<-done
	<-done
	l.Info("All done")
}
