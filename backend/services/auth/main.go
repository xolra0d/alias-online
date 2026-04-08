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

	l.Info("initializing secrets")
	secrets, err := NewSecrets(
		l,
		cfg.Argon2idTime,
		cfg.Argon2idMemory,
		cfg.Argon2idThreads,
		cfg.Argon2idOutLen,
		cfg.JwtPrivateKeyFilename,
	)
	if err != nil {
		l.Error("failed to start secrets", "error", err)
		return
	}

	l.Info("starting server")
	RunGrpcServer(
		secrets,
		postgres,
		l,
		cfg.AddAccountTimeout,
		cfg.FindAccountTimeout,
		cfg.JWTCookieTimeout,
		cfg.RunningAddr,
		cfg.ShutdownTimeout,
	)
	l.Info("All done")
}
