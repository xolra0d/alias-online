package main

import (
	"log/slog"
	"net/http"

	"github.com/xolra0d/alias-online/shared/pkg/api"
)

type Handles struct {
	vocabs *VocabManager
	logger *slog.Logger
}

func NewHTTPHandles(vocabs *VocabManager, logger *slog.Logger) *Handles {
	return &Handles{
		vocabs,
		logger,
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

// AvailableVocabs retrieves names of available vocabs
func (h *Handles) AvailableVocabs(w http.ResponseWriter, _ *http.Request) {
	const op = "main.AvailableVocabs"

	err := api.WriteJSON(w, http.StatusOK, map[string]any{"vocabs": h.vocabs.AvailableVocabs()})
	if err != nil {
		h.logger.Error("could not write response", "op", op, "err", err)
		return
	}
}

// Vocab retrieves names of available vocabs
func (h *Handles) Vocab(w http.ResponseWriter, r *http.Request) {
	primary, rude := h.vocabs.Vocab(r.URL.Query().Get("vocab"))
	err := api.WriteJSON(w, http.StatusOK, map[string]any{"primary_words": primary, "rude_words": rude})
	if err != nil {
		h.logger.Error("could not write response", "err", err)
		return
	}
}
