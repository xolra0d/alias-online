package main

import (
	"log/slog"
	"net/http"

	"github.com/xolra0d/alias-online/shared/pkg/logger"
)

func main() {
	const op = "main.main"

	cfg := LoadServerConfig()
	l := slog.New(logger.NewHandler(nil))

	pgPool, err := InitPool(cfg.PostgresUrl)
	if err != nil {
		l.Error("failed to connect to database", "op", op, "error", err)
		return
	}
	defer pgPool.Close()

	postgres := NewPostgres(pgPool, l)
	vocabs := NewVocabManager(postgres, l)
	go vocabs.StartObservation(cfg.LoadVocabsTimeout, cfg.PollInterval)

	handles := NewHTTPHandles(vocabs, l)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/ok", handles.Healthy)
	mux.HandleFunc("GET /api/available-vocabs", handles.AvailableVocabs)
	mux.HandleFunc("GET /api/vocab", handles.Vocab)

	RunServer(mux, l, cfg.RunningAddr, cfg.ShutdownTimeout)
	vocabs.StopObservation()
}
