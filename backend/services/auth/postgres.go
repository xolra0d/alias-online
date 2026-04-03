package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewPostgres(postgresUrl string, logger *slog.Logger) (*Postgres, error) {
	db, err := pgxpool.New(context.Background(), postgresUrl)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		db:     db,
		logger: logger,
	}, nil
}

func (p *Postgres) Close() {
	p.db.Close()
}

func (p *Postgres) AddAccount(ctx context.Context, creds Credentials) error {
	const op = "main.AddAccount"

	query := "INSERT INTO users (login, name, secret_hash) VALUES ($1, $2, $3)"
	_, err := p.db.Exec(ctx, query, creds.Login, creds.Name, creds.Password)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" {
			return fmt.Errorf("user already exists")
		}
		p.logger.Error("user insertion error", "op", op, "err", err)
		return err
	}

	return nil
}

func (p *Postgres) FindAccount(ctx context.Context, login string) (hash string, ok bool) {
	query := "SELECT secret_hash FROM users WHERE login=$1"
	row := p.db.QueryRow(ctx, query, login)
	var storedHash string
	err := row.Scan(&storedHash)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			p.logger.Error("hash retrieval error", "err", err)
		}
		return "", false
	}
	return storedHash, true
}
