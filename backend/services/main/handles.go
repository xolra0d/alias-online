package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/xolra0d/alias-online/shared/pkg/api"
	"github.com/xolra0d/alias-online/shared/pkg/middleware"
	pbAuth "github.com/xolra0d/alias-online/shared/proto/auth"
	pbRoomManager "github.com/xolra0d/alias-online/shared/proto/room_manager"
	pbVocabManager "github.com/xolra0d/alias-online/shared/proto/vocab_manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewRoomManagerClient(roomManagerUrl string, logger *slog.Logger) (pbRoomManager.RoomManagerServiceClient, func() error, error) {
	conn, err := grpc.NewClient(roomManagerUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to connect to room manager", "roomManagerUrl", roomManagerUrl, "err", err)
		return nil, nil, err
	}
	//defer conn.Close()
	return pbRoomManager.NewRoomManagerServiceClient(conn), conn.Close, nil
}

func NewVocabManagerClient(vocabManagerUrl string, logger *slog.Logger) (pbVocabManager.VocabManagerServiceClient, func() error, error) {
	conn, err := grpc.NewClient(vocabManagerUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to connect to vocab manager", "vocabManagerUrl", vocabManagerUrl, "err", err)
		return nil, nil, err
	}
	//defer conn.Close()
	return pbVocabManager.NewVocabManagerServiceClient(conn), conn.Close, nil
}

func NewAuthClient(authUrl string, logger *slog.Logger) (pbAuth.AuthServiceClient, func() error, error) {
	conn, err := grpc.NewClient(authUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to connect to vocab manager", "authUrl", authUrl, "err", err)
		return nil, nil, err
	}
	//defer conn.Close()
	return pbAuth.NewAuthServiceClient(conn), conn.Close, nil
}

type Handles struct {
	authClient         pbAuth.AuthServiceClient
	roomManagerClient  pbRoomManager.RoomManagerServiceClient
	vocabManagerClient pbVocabManager.VocabManagerServiceClient
	logger             *slog.Logger

	AddAccountTimeout  time.Duration
	FindAccountTimeout time.Duration
	JWTCookieTimeout   time.Duration
	JWTCookiePath      string
	JWTCookieSecure    bool
	JWTCookieHTTPOnly  bool
	JWTCookieDomain    string
}

func NewHTTPHandles(
	authClient pbAuth.AuthServiceClient,
	vocabManagerClient pbVocabManager.VocabManagerServiceClient,
	roomManagerClient pbRoomManager.RoomManagerServiceClient,
	logger *slog.Logger,
	addAccountTimeout, findAccountTimeout, JWTCookieTimeout time.Duration,
	JWTCookiePath string,
	JWTCookieSecure bool,
	JWTCookieHTTPOnly bool,
	JWTCookieDomain string,
) *Handles {
	return &Handles{
		authClient,
		roomManagerClient,
		vocabManagerClient,
		logger,

		addAccountTimeout,
		findAccountTimeout,
		JWTCookieTimeout,
		JWTCookiePath,
		JWTCookieSecure,
		JWTCookieHTTPOnly,
		JWTCookieDomain,
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

// AvailableVocabs handles /ok requests
func (h *Handles) AvailableVocabs(w http.ResponseWriter, r *http.Request) {
	const op = "main.AvailableVocabs"

	v, err := h.vocabManagerClient.GetAvailableVocabs(r.Context(), &emptypb.Empty{})
	if err != nil {
		h.logger.Error("could not get available vocabs", "op", op, "err", err)
		writeErrorAndLogWriteError(w, http.StatusInternalServerError, "could not get available vocabs", h.logger)
		return
	}

	if err = api.WriteJSON(w, http.StatusOK, map[string]any{"vocabs": v.Names}); err != nil {
		h.logger.Error("could not write response", "op", op, "err", err)
	}
}

func (h *Handles) Register(w http.ResponseWriter, r *http.Request) {
	const op = "main.Register"

	var regCreds struct {
		Name     string `json:"name"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&regCreds)
	if err != nil {
		h.logger.Error("could not decode register request", "op", op, "err", err, "body", r.Body)
		writeErrorAndLogWriteError(w, http.StatusBadRequest, "could not decode credentials", h.logger)
		return
	}

	resp, err := h.authClient.Register(r.Context(), &pbAuth.RegisterRequest{Name: regCreds.Name, Login: regCreds.Login, Password: regCreds.Password})
	if err != nil {
		h.logger.Error("could not register user", "op", op, "err", err)
		fmt.Println(err.Error())
		switch status.Code(err) {
		case codes.AlreadyExists:
			writeErrorAndLogWriteError(w, http.StatusConflict, "user already exists", h.logger)
		case codes.InvalidArgument:
			writeErrorAndLogWriteError(w, http.StatusBadRequest, "invalid credentials", h.logger)
		case codes.Internal:
			writeErrorAndLogWriteError(w, http.StatusInternalServerError, "internal server error", h.logger)
		default:
			h.logger.Error("invalid error type", "op", op, "err", err, "type", status.Code(err))
			writeErrorAndLogWriteError(w, http.StatusInternalServerError, "could not register user", h.logger)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     middleware.LoginCookieName,
		Value:    resp.Token,
		Path:     h.JWTCookiePath,
		MaxAge:   int(time.Until(time.Unix(resp.Exp, 0)).Seconds()),
		Secure:   h.JWTCookieSecure,
		HttpOnly: h.JWTCookieHTTPOnly,
		SameSite: http.SameSiteLaxMode,
		Domain:   h.JWTCookieDomain,
	})

	err = api.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
	if err != nil {
		h.logger.Error("could not write response", "err", err)
	}
}

func (h *Handles) Login(w http.ResponseWriter, r *http.Request) {
	const op = "main.Login"

	var regCreds struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&regCreds)
	if err != nil {
		h.logger.Error("could not decode login request", "op", op, "err", err, "body", r.Body)
		writeErrorAndLogWriteError(w, http.StatusBadRequest, "could not decode credentials", h.logger)
		return
	}

	resp, err := h.authClient.Login(r.Context(), &pbAuth.LoginRequest{Login: regCreds.Login, Password: regCreds.Password})
	if err != nil {
		h.logger.Error("could not login user", "op", op, "err", err)
		fmt.Println(err.Error())
		switch status.Code(err) {
		case codes.NotFound:
			writeErrorAndLogWriteError(w, http.StatusConflict, "user not found", h.logger)
		case codes.InvalidArgument:
			writeErrorAndLogWriteError(w, http.StatusBadRequest, "invalid credentials", h.logger)
		case codes.Unauthenticated:
			writeErrorAndLogWriteError(w, http.StatusBadRequest, "wrong credentials", h.logger)
		case codes.Internal:
			writeErrorAndLogWriteError(w, http.StatusInternalServerError, "internal server error", h.logger)
		default:
			h.logger.Error("invalid error type", "op", op, "err", err, "type", status.Code(err))
			writeErrorAndLogWriteError(w, http.StatusInternalServerError, "could not register user", h.logger)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     middleware.LoginCookieName,
		Value:    resp.Token,
		Path:     h.JWTCookiePath,
		MaxAge:   int(time.Until(time.Unix(resp.Exp, 0)).Seconds()),
		Secure:   h.JWTCookieSecure,
		HttpOnly: h.JWTCookieHTTPOnly,
		SameSite: http.SameSiteLaxMode,
		Domain:   h.JWTCookieDomain,
	})

	err = api.WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
	if err != nil {
		h.logger.Error("could not write response", "err", err)
	}
}

func (h *Handles) Play(w http.ResponseWriter, r *http.Request) {
	roomId := r.PathValue("roomId")
	if roomId == "" {
		writeErrorAndLogWriteError(w, http.StatusBadRequest, "missing roomId", h.logger)
		return
	}

	worker, err := h.roomManagerClient.GetRoomWorker(r.Context(), &pbRoomManager.GetRoomWorkerRequest{RoomId: roomId})
	if err != nil {
		h.logger.Error("could not get best room worker", "roomId", roomId, "err", err)
		writeErrorAndLogWriteError(w, http.StatusInternalServerError, "internal error", h.logger)
		return
	}
	err = api.WriteJSON(w, http.StatusOK, map[string]any{"worker": worker.GetWorker()})
	if err != nil {
		h.logger.Error("could not write response", "err", err)
	}
}

// shorthand for writing error
func writeErrorAndLogWriteError(w http.ResponseWriter, status int, err string, logger *slog.Logger) {
	writeErr := api.WriteJSON(w, status, map[string]any{"err": err})
	if writeErr != nil {
		logger.Error("could not write response", "err", writeErr)
	}
}
