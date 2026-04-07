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

func (p *Postgres) AddAccount(ctx context.Context, name, login, hash string) *Error {
	const op = "main.AddAccount"

	const PostgresUniqueViolation = "23505" // err code for violating uniqueness requirement

	query := "INSERT INTO users (login, name, secret_hash) VALUES ($1, $2, $3)"
	_, err := p.db.Exec(ctx, query, login, name, hash)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == PostgresUniqueViolation {
			return NewError(ErrUserAlreadyExists, fmt.Errorf("user already exists"))
		}
		p.logger.Error("user insertion error", "op", op, "creds", map[string]string{
			"login": login,
			"name":  name,
			"hash":  hash,
		}, "err", err)
		return NewError(ErrInternal, err)
	}
	return nil
}

func (p *Postgres) FindAccount(ctx context.Context, login string) (string, *Error) {
	query := "SELECT secret_hash FROM users WHERE login=$1"
	row := p.db.QueryRow(ctx, query, login)
	var storedHash string
	err := row.Scan(&storedHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", NewError(ErrUserNotFound, fmt.Errorf("user not found"))
		}
		return "", NewError(ErrInternal, err)
	}
	return storedHash, nil
}
