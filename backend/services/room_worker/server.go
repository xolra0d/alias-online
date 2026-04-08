package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/xolra0d/alias-online/shared/pkg/api"
	"github.com/xolra0d/alias-online/shared/pkg/middleware"
)

type Handles struct {
	secrets *Secrets
	logger  *slog.Logger

	rooms *Rooms
}

func NewHandles(secrets *Secrets, logger *slog.Logger, rooms *Rooms) *Handles {
	return &Handles{
		secrets: secrets,
		logger:  logger,
		rooms:   rooms,
	}
}

func RunHttpClient(handles *Handles, secrets *Secrets, logger *slog.Logger, runningAddr string, shutdownTimeout time.Duration, shouldStop, done chan struct{}) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/ok", handles.Healthy)
	mux.Handle("GET /api/play/{roomId}", middleware.Chain( // GET for websockets
		http.HandlerFunc(handles.InitWS),
		middleware.AuthJWT(logger, secrets.CheckJwt),
	))

	server := &http.Server{
		Addr: runningAddr,
		Handler: middleware.Chain(
			mux,
			middleware.Logging(logger),
		),
	}

	go func() {
		logger.Info("starting HTTP server", "addr", runningAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("listen error", "err", err)
		}
	}()

	<-shouldStop
	logger.Info("HTTP server shutdown initiated")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("shutdown failed", "err", err)
	}

	logger.Info("HTTP server stopped")
	done <- struct{}{}
}

// InitWS validates user credentials and tries to update HTTP to Websocket connection.
func (h *Handles) InitWS(w http.ResponseWriter, r *http.Request) {
	const op = "main.InitWS"

	roomId := r.PathValue("roomId")
	if roomId == "" {
		err := api.WriteJSON(w, http.StatusBadRequest, map[string]any{"err": "missing room_worker id"})
		if err != nil {
			h.logger.Error("could not write response", "op", op, "err", err)
		}
		return
	}

	username := r.Context().Value(middleware.LoginContextKey).(string)
	name := r.URL.Query().Get("name")
	if name == "" {
		err := api.WriteJSON(w, http.StatusBadRequest, map[string]any{"err": "missing name"})
		if err != nil {
			h.logger.Error("could not write response", "op", op, "err", err)
		}
		return
	}
	err := h.rooms.RunWS(w, r, roomId, username, name)
	if err != nil {
		h.logger.Error("error while RunWS", "roomId", roomId, "username", username, "err", err)
		err = api.WriteJSON(w, http.StatusInternalServerError, map[string]any{"err": "internal error"})
		if err != nil {
			h.logger.Error("could not write response", "op", op, "err", err)
			return
		}
	}
}

// Healthy handles /ok requests
func (h *Handles) Healthy(w http.ResponseWriter, _ *http.Request) {
	const op = "main.Healthy"

	err := api.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
	if err != nil {
		h.logger.Error("could not write response", "op", op, "err", err)
		return
	}
}
