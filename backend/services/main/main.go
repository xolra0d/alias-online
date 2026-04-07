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

	//conn, err := grpc.NewClient(cfg.VocabsURLGateway, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	l.Error("could not connect to vocab manager", "err", err)
	//	return
	//}
	//defer conn.Close()
	//vocabManager := pbVocab.NewVocabManagerServiceClient(conn)

	//conn, err = grpc.NewClient(cfg.AuthURLGateway, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	l.Error("could not connect to auth", "err", err)
	//	return
	//}
	//defer conn.Close()
	//auth := pbAuth.NewAuthServiceClient(conn)
	vocabManager, closeVocabManager, err := NewVocabManagerClient("localhost:8070", l)
	if err != nil {
		l.Error("could not connect to vocab manager", "err", err)
		return
	}
	defer closeVocabManager()
	roomManager, closeRoomManager, err := NewRoomManagerClient("localhost:8060", l)
	if err != nil {
		l.Error("could not connect to room manager", "err", err)
		return
	}
	defer closeRoomManager()
	auth, closeAuth, err := NewAuthClient("localhost:8090", l)
	if err != nil {
		l.Error("could not connect to auth", "err", err)
		return
	}
	defer closeAuth()

	l.Info("initializing secrets")
	secrets, err := NewSecrets(cfg.JWTPublicKeyPath, l)
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
	mux.Handle("GET /api/play/{roomId}", middleware.Chain(
		http.HandlerFunc(handles.Play),
		middleware.AuthJWT(l, secrets.CheckJwtToken),
	))

	cors := middleware.NewCors(
		strings.Split(cfg.AllowedOrigins, ","),
		[]string{"GET", "POST", "OPTIONS"},
		[]string{"Origin", "Content-Length", "Content-Type"},
		true,
	)
	csrf := middleware.NewCSRF(strings.Split(cfg.AllowedOrigins, ","))

	RunServer(mux, csrf, cors, l, cfg.RunningAddr, cfg.ShutdownTimeout)
}
