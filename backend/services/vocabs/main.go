package main

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/xolra0d/alias-online/shared/pkg/logger"
	"github.com/xolra0d/alias-online/shared/pkg/middleware"
)

func main() {
	const op = "main.main"

	cfg := LoadServerConfig()
	l := slog.New(logger.NewHandler(nil))

	postgres, err := NewPostgres(cfg.PostgresUrl, l)
	if err != nil {
		l.Error("failed to connect to database", "op", op, "error", err)
		return
	}
	defer postgres.Close()
	vocabs, ok := NewVocabManager(postgres, l, cfg.LoadVocabsTimeout, cfg.PollInterval)
	if !ok {
		l.Error("failed to initiate vocab manager", "op", op, "error", err)
		return
	}
	go vocabs.StartObservation()

	handles := NewHTTPHandles(vocabs, l)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/ok", handles.Healthy)
	mux.HandleFunc("GET /api/available-vocabs", handles.AvailableVocabs)
	mux.HandleFunc("GET /api/vocab", handles.Vocab)

	cors := middleware.NewCors(
		strings.Split(cfg.AllowedOrigins, ","),
		[]string{"GET", "OPTIONS"},
		[]string{"Origin", "Content-Length", "Content-Type"},
		true,
	)

	RunServer(mux, cors, l, cfg.RunningAddr, cfg.ShutdownTimeout)
	vocabs.StopObservation()
}
