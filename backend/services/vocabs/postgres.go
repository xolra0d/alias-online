package main

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPool(postgresUrl string) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), postgresUrl)
}

type Postgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewPostgres(db *pgxpool.Pool, logger *slog.Logger) *Postgres {
	return &Postgres{
		db:     db,
		logger: logger,
	}
}

func (p *Postgres) LoadVocabs(ctx context.Context) (map[string]Vocabulary, bool) {
	const op = "main.LoadVocabs"

	vocabs := map[string]Vocabulary{}
	query := "SELECT language, primary_words, rude_words FROM vocabularies WHERE available = TRUE ORDER BY language"
	rows, err := p.db.Query(ctx, query)
	if err != nil {
		p.logger.Error("failed to load vocabs", "op", op, "err", err)
		return vocabs, false
	}
	defer rows.Close()
	for rows.Next() {
		var language string
		var primaryWords []string
		var rudeWords []string
		if err := rows.Scan(&language, &primaryWords, &rudeWords); err != nil {
			p.logger.Error("failed to load vocabulary", op, "op", "err", err)
			return vocabs, false
		}
		vocabs[language] = Vocabulary{primaryWords, rudeWords}
	}
	return vocabs, true
}
