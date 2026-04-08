package main

import (
	"log/slog"

	"github.com/xolra0d/alias-online/shared/pkg/logger"
)

func main() {
	l := slog.New(logger.NewHandler(nil))

	l.Info("loading configuration")
	cfg := LoadServerConfig()

	l.Info("trying to connect to postgres")
	postgres, err := NewPostgres(cfg.PostgresUrl, l)
	if err != nil {
		l.Error("failed to connect to postgres", "error", err)
		return
	}
	defer postgres.Close()

	l.Info("initializing vocabularies manager")
	vocabs, ok := NewVocabManager(postgres, l, cfg.LoadVocabsTimeout, cfg.VocabsPollInterval, cfg.ClosePostgresConnTimeout)
	if !ok {
		l.Error("failed to initiate vocabularies manager", "error", err)
		return
	}

	l.Info("starting vocabularies search")
	go vocabs.StartObservation()

	l.Info("starting server")
	RunGrpcServer(vocabs, l, cfg.RunningAddr, cfg.ShutdownTimeout)
	l.Info("stopping vocabularies search")
	vocabs.StopObservation()
	l.Info("All done")
}
