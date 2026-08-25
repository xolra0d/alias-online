package transport

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/xolra0d/alias-online/internal/config"
	"github.com/xolra0d/alias-online/internal/database"
	"github.com/xolra0d/alias-online/internal/room"
)

// Handles holds transport handles and, realistically, the state of the program.
type Handles struct {
	postgres *database.Postgres
	rooms    *Rooms
	vocabs   *room.Vocabularies
	logger   *config.Logger

	createUserTimeout time.Duration
	createRoomTimeout time.Duration
}

func NewHandles(
	postgres *database.Postgres,
	rooms *Rooms,
	vocabs *room.Vocabularies,
	logger *config.Logger,
	createUserTimeout time.Duration,
	createRoomTimeout time.Duration,
) *Handles {
	return &Handles{
		postgres:          postgres,
		rooms:             rooms,
		vocabs:            vocabs,
		logger:            logger,
		createUserTimeout: createUserTimeout,
		createRoomTimeout: createRoomTimeout,
	}
}

// Healthy handles /ok requests
func (h *Handles) Healthy(w http.ResponseWriter, _ *http.Request) {
	const op = "transport.Healthy"

	err := WriteJSON(w, http.StatusOK, P{"ok": true})
	if err != nil {
		h.logger.Error(op, "could not write response", "err", err)
		return
	}
}

// AvailableLanguages returns loaded vocabs names.
func (h *Handles) AvailableLanguages(w http.ResponseWriter, _ *http.Request) {
	const op = "transport.AvailableLanguages"

	languages := h.vocabs.Languages()
	slices.SortFunc(languages, func(a, b string) int {
		aEnglish := strings.EqualFold(a, "English")
		bEnglish := strings.EqualFold(b, "English")
		switch {
		case aEnglish && !bEnglish:
			return -1
		case !aEnglish && bEnglish:
			return 1
		}

		aLower := strings.ToLower(a)
		bLower := strings.ToLower(b)
		if cmp := strings.Compare(aLower, bLower); cmp != 0 {
			return cmp
		}
		return strings.Compare(a, b)
	})

	err := WriteJSON(w, http.StatusOK, P{"languages": languages})
	if err != nil {
		h.logger.Error(op, "could not write response", "err", err)
		return
	}
}

// CreateUser generates random login, name, secret for user and returns as Credentials.
func (h *Handles) CreateUser(w http.ResponseWriter, r *http.Request) {
	const op = "transport.CreateUser"

	ctx, cancel := context.WithTimeout(r.Context(), h.createUserTimeout)
	defer cancel()
	credentials, err := h.postgres.CreateUser(ctx)
	if err != nil {
		h.logger.Error(op, "could not create user", "err", err)
		err = WriteJSON(w, http.StatusInternalServerError, P{"err": err})
		if err != nil {
			h.logger.Error(op, "could not write response", "err", err)
		}
		return
	}
	err = WriteJSON(w, http.StatusOK, P{"credentials": credentials})
	if err != nil {
		h.logger.Error(op, "could not write response", "err", err)
	}
}

func normalizeAdditionalVocabulary(raw string, maxWordLength int, maxWords int) ([]string, error) {
	seen := map[string]struct{}{}

	for part := range strings.SplitSeq(raw, ",") {
		word := strings.TrimSpace(part)
		if word == "" {
			continue
		}
		if utf8.RuneCountInString(word) > maxWordLength {
			return nil, fmt.Errorf("additional vocabulary word %q is too long (max %d chars)", word, maxWordLength)
		}
		key := strings.ToLower(word)
		seen[key] = struct{}{}
		if len(seen) > maxWords {
			return nil, fmt.Errorf("additional vocabulary exceeds %d words", maxWords)
		}
	}

	if len(seen) == 0 {
		return []string{}, nil
	}

	return slices.Collect(maps.Keys(seen)), nil
}

// CreateRoom validates room config from form and inserts it to database, returning roomId.
func (h *Handles) CreateRoom(w http.ResponseWriter, r *http.Request) {
	const op = "transport.CreateRoom"

	adminId := uuid.MustParse(r.Header.Get("User-Id")) // verified at `Auth`

	language := r.PostFormValue("language")
	ok := h.vocabs.Contains(language)
	if !ok {
		err := WriteJSON(w, http.StatusBadRequest, P{"err": "invalid language"})
		if err != nil {
			h.logger.Error(op, "could not write response", "err", err)
		}
		return
	}

	rudeWords, err := strconv.ParseBool(r.PostFormValue("rude-words"))
	if err != nil {
		err = WriteJSON(w, http.StatusBadRequest, P{"err": "invalid rude-words"})
		if err != nil {
			h.logger.Error(op, "could not write response", "err", err)
		}
		return
	}

	additional, err := normalizeAdditionalVocabulary(
		r.PostFormValue("additional-vocabulary"),
		h.rooms.maxAdditionalWordLength,
		h.rooms.maxAdditionalVocabularyWords,
	)

	if err != nil {
		err := WriteJSON(w, http.StatusBadRequest, P{"err": err})
		if err != nil {
			h.logger.Error(op, "could not write response", "err", err)
		}
		return
	}

	clock, err := strconv.ParseUint(r.PostFormValue("clock"), 10, 64)
	if err != nil {
		err = WriteJSON(w, http.StatusBadRequest, P{"err": "invalid clock"})
		if err != nil {
			h.logger.Error(op, "could not write response", "err", err)
		}
		return
	}
	if int(clock) < h.rooms.minClock || int(clock) > h.rooms.maxClock {
		err = WriteJSON(w, http.StatusBadRequest, P{"err": fmt.Sprintf("invalid clock: must be between %d and %d", h.rooms.minClock, h.rooms.maxClock)})
		if err != nil {
			h.logger.Error(op, "could not write response", "err", err)
		}
		return
	}

	cfg := room.RoomConfig{
		Language:             language,
		RudeWords:            rudeWords,
		AdditionalVocabulary: additional,
		Clock:                int(clock),
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.createRoomTimeout)
	defer cancel()
	roomId, err := h.postgres.AddRoom(ctx, adminId, cfg)
	if err != nil {
		h.logger.Error("could not add room", "err", err)
		err = WriteJSON(w, http.StatusInternalServerError, P{"err": err})
		if err != nil {
			h.logger.Error(op, "could not write response", "err", err)
		}
		return
	}
	err = WriteJSON(w, http.StatusOK, P{"room": roomId})
	if err != nil {
		h.logger.Error(op, "could not write response", "err", err)
	}
}

// InitWS validates user credentials and tries to update HTTP to Websocket connection.
func (h *Handles) InitWS(w http.ResponseWriter, r *http.Request) {
	const op = "transport.InitWS"

	roomId := r.PathValue("roomId")
	if roomId == "" {
		err := WriteJSON(w, http.StatusBadRequest, P{"err": "missing room id"})
		if err != nil {
			h.logger.Error(op, "could not write response", "err", err)
		}
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		err := WriteJSON(w, http.StatusBadRequest, P{"err": "missing user id"})
		if err != nil {
			h.logger.Error(op, "could not write response", "err", err)
		}
		return
	}
	userId, err := uuid.Parse(id)
	if err != nil {
		err := WriteJSON(w, http.StatusBadRequest, P{"err": "invalid user id"})
		if err != nil {
			h.logger.Error(op, "could not write response", "err", err)
		}
		return
	}
	secret := r.URL.Query().Get("secret")
	if secret == "" {
		err := WriteJSON(w, http.StatusBadRequest, P{"err": "missing secret"})
		if err != nil {
			h.logger.Error(op, "could not write response", "err", err)
		}
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		err := WriteJSON(w, http.StatusBadRequest, P{"err": "missing name"})
		if err != nil {
			h.logger.Error(op, "could not write response", "err", err)
		}
		return
	}
	if ok := h.postgres.ValidateUser(r.Context(), database.UserCredentials{Id: userId, Secret: secret}); !ok {
		err := WriteJSON(w, http.StatusUnauthorized, P{"err": "invalid credentials"})
		if err != nil {
			h.logger.Error(op, "could not write response", "err", err)
		}
		return
	}

	err = h.rooms.ServeWS(w, r, userId, name, roomId, h.postgres, h.vocabs)
	if err != nil {
		h.logger.Error(op, "room exit error", "err", err)
	}
}
