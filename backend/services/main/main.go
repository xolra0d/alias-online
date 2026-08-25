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

	//postgres, err := NewPostgres(cfg.PostgresUrl, l)
	//if err != nil {
	//	l.Error("failed to connect to database", "op", op, "error", err)
	//	return
	//}
	//defer postgres.Close()

	//secrets, err := NewSecrets(l, cfg.Argon2idTime, cfg.Argon2idMemory, cfg.Argon2idThreads, cfg.Argon2idOutLen, cfg.RSAPrivateKeyFilename)
	//if err != nil {
	//	l.Error("failed to start secrets", "op", op, "error", err)
	//	return
	//}
	handles := NewHTTPHandles(
		//secrets,
		//postgres,
		l,
		cfg.AddAccountTimeout,
		cfg.FindAccountTimeout,
		cfg.JWTCookieTimeout,
		//cfg.JWTCookiePath,
		//cfg.JWTCookieSecure,
		//cfg.JWTCookieHTTPOnly,
		//cfg.JWTCookieDomain,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/ok", handles.Healthy)
	mux.HandleFunc("GET /api/available-vocabs", handles.AvailableVocabs)

	//mux.HandleFunc("GET /api/.well-known/jwks", handles.PublicKeys)
	//mux.Handle("POST /api/register", middleware.Chain(
	//	http.HandlerFunc(handles.Register),
	//	IPRateLimiter(30, time.Minute, 128, l),
	//))
	//mux.Handle("POST /api/login", middleware.Chain(
	//	http.HandlerFunc(handles.Login),
	//	IPRateLimiter(30, time.Minute, 128, l),
	//))
	//mux.Handle("GET /api/protected/ok", middleware.Chain(
	//	http.HandlerFunc(handles.Healthy),
	//	middleware.AuthJWT(l, secrets.ValidateJWT),
	//))

	cors := middleware.NewCors(
		strings.Split(cfg.AllowedOrigins, ","),
		[]string{"GET", "POST", "OPTIONS"},
		[]string{"Origin", "Content-Length", "Content-Type"},
		true,
	)
	csrf := middleware.NewCSRF(strings.Split(cfg.AllowedOrigins, ","))

	RunServer(mux, csrf, cors, l, cfg.RunningAddr, cfg.ShutdownTimeout)
}
