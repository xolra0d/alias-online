package main

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/xolra0d/alias-online/shared/pkg/logger"
	"github.com/xolra0d/alias-online/shared/pkg/middleware"
)

func main() {
	const op = "main.main"

	cfg := LoadServerConfig()
	l := slog.New(logger.NewHandler(nil))

	vocabManager, closeVocabManager, err := NewVocabManagerClient(cfg.VocabsURLGateway, l)
	if err != nil {
		l.Error("could not connect to vocab manager", "err", err)
		return
	}
	defer closeVocabManager()
	roomManager, closeRoomManager, err := NewRoomManagerClient(cfg.RoomManagerURLGateway, l)
	if err != nil {
		l.Error("could not connect to room manager", "err", err)
		return
	}
	defer closeRoomManager()
	auth, closeAuth, err := NewAuthClient(cfg.AuthURLGateway, l)
	if err != nil {
		l.Error("could not connect to auth", "err", err)
		return
	}
	defer closeAuth()

	l.Info("initializing secrets")
	secrets, err := NewSecrets(cfg.JwtPublicKeyPath, l)
	if err != nil {
		l.Error("failed to load secrets", "error", err)
		return
	}

	handles := NewHTTPHandles(
		auth,
		vocabManager,
		roomManager,
		l,
		cfg.AddAccountTimeout,
		cfg.FindAccountTimeout,
		cfg.JWTCookieTimeout,
		cfg.JWTCookiePath,
		cfg.JWTCookieSecure,
		cfg.JWTCookieHTTPOnly,
		cfg.JWTCookieDomain,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/ok", handles.Healthy)
	mux.HandleFunc("GET /api/available-vocabs", handles.AvailableVocabs)
	mux.Handle("POST /api/register", middleware.Chain(
		http.HandlerFunc(handles.Register),
		IPRateLimiter(30, time.Minute, 128, l),
	))
	mux.Handle("POST /api/login", middleware.Chain(
		http.HandlerFunc(handles.Login),
		IPRateLimiter(30, time.Minute, 128, l),
	))
	mux.Handle("GET /api/protected/ok", middleware.Chain(
		http.HandlerFunc(handles.Healthy),
		middleware.AuthJWT(l, secrets.CheckJwt),
	))
	mux.Handle("GET /api/protected/play/{roomId}", middleware.Chain(
		http.HandlerFunc(handles.Play),
		middleware.AuthJWT(l, secrets.CheckJwt),
	))

	cors := middleware.NewCors(
		strings.Split(cfg.AllowedOrigins, ","),
		[]string{"GET", "POST", "OPTIONS"},
		[]string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		true,
	)
	csrf := middleware.NewCSRF(strings.Split(cfg.AllowedOrigins, ","))

	RunServer(mux, csrf, cors, l, cfg.RunningAddr, cfg.ShutdownTimeout)
}
