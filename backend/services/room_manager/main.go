package main

import (
	"log/slog"

	"github.com/xolra0d/alias-online/shared/pkg/logger"
)

func main() {
	l := slog.New(logger.NewHandler(nil))

	l.Info("loading configuration")
	cfg := LoadServerConfig()

	l.Info("trying to connect to redis")
	db := NewDatabase(cfg.RedisAddr, cfg.RedisUsername, cfg.RedisPassword, cfg.RedisDB)

	l.Info("initializing room manager")
	manager := NewManager(db, l, cfg.PollInterval, cfg.WorkerExpiry, cfg.RetrieveActiveWorkersTimeout)

	l.Info("starting scan for new workers")
	go manager.ScanForNewWorkers()

	l.Info("starting server")
	RunGrpcServer(
		manager,
		l,
		cfg.RunningAddr,
		cfg.ShutdownTimeout,
	)

	l.Info("stopping scan for new workers")
	manager.StopScanForNewWorkers()
}
