package main

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

// NewPostgres creates new instance of postgres client.
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

// Close closes postgres pool.
func (p *Postgres) Close() {
	p.db.Close()
}

// LoadVocabs loads vocabularies from postgres.
func (p *Postgres) LoadVocabs(ctx context.Context) (map[string]Vocabulary, bool) {
	vocabs := map[string]Vocabulary{}
	query := "SELECT language, primary_words, rude_words FROM vocabularies WHERE available = TRUE ORDER BY language"
	rows, err := p.db.Query(ctx, query)
	if err != nil {
		p.logger.Error("failed to load available languages", "err", err)
		return vocabs, false
	}
	defer rows.Close()
	for rows.Next() {
		var language string
		var primaryWords []string
		var rudeWords []string
		if err := rows.Scan(&language, &primaryWords, &rudeWords); err != nil {
			p.logger.Error("failed to load vocabulary", "err", err)
			return vocabs, false
		}
		vocabs[language] = Vocabulary{primaryWords, rudeWords}
	}
	return vocabs, true
}
