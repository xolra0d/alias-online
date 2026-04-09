package main

import (
	"context"
	"fmt"
	"log/slog"
	mrand "math/rand"

	"github.com/jackc/pgx/v5"
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

// LoadRoom tries to load `Room` with its users and turn order. Sets all users to not ready.
func (p *Postgres) LoadRoom(ctx context.Context, roomId string, getVocab func(ctx context.Context, s string) (Vocabulary, error)) (*Room, error) {
	const op = "database.LoadRoom"

	var id string
	var admin string
	var seed int
	var currentWordIndex int
	var currentPlayerId string
	var gameState int
	var language string
	var rudeWords bool
	var additionalVocabulary []string
	var clock int
	query := "SELECT id, admin, seed, current_word_index, current_player_login, game_state, language, rude_words, additional_vocabulary, clock FROM rooms WHERE id=$1"
	err := p.db.QueryRow(ctx, query, roomId).Scan(
		&id,
		&admin,
		&seed,
		&currentWordIndex,
		&currentPlayerId,
		&gameState,
		&language,
		&rudeWords,
		&additionalVocabulary,
		&clock,
	)
	if err != nil {
		p.logger.Error("could not load rooms", "roomId", roomId, "err", err)
		return nil, err
	}

	vocab, err := getVocab(ctx, language)
	if err != nil {
		p.logger.Error("could not find vocab", "language", language)
		return nil, fmt.Errorf("could not find vocab: %s", language)

	}

	allWords := vocab.PrimaryWords // it's our copy, so it's fine to mutate
	if rudeWords {
		allWords = append(allWords, vocab.RudeWords...)
	}
	if len(additionalVocabulary) != 0 {
		allWords = append(allWords, additionalVocabulary...)
	}

	rng := mrand.New(mrand.NewSource(int64(seed)))
	rng.Shuffle(len(allWords), func(i, j int) {
		allWords[i], allWords[j] = allWords[j], allWords[i]
	})

	cfg := &RoomConfig{
		Seed:                 seed,
		AllWords:             allWords,
		Language:             language,
		RudeWords:            rudeWords,
		AdditionalVocabulary: additionalVocabulary,
		Clock:                clock,
	}

	query = `SELECT rp.user_login, rp.words_tried, rp.words_guessed, u.name
		FROM room_participants rp
		JOIN users u ON rp.user_login = u.login
		WHERE rp.room_id = $1
		ORDER BY rp.turn_order ASC`
	rows, err := p.db.Query(ctx, query, roomId)
	if err != nil {
		p.logger.Error("failed to load rooms participants", "roomId", roomId, "err", err)
		return nil, err
	}
	defer rows.Close()

	players := map[string]*Player{}
	//goland:noinspection GoPreferNilSlice
	turnOrder := []string{}
	currentPlayer := 0

	i := 0
	for rows.Next() {
		var userLogin string
		var wordsTried int
		var wordsGuessed int
		var name string
		err = rows.Scan(&userLogin, &wordsTried, &wordsGuessed, &name)
		if err != nil {
			p.logger.Error("failed to scan rooms participant", "roomId", roomId, "err", err)
			return nil, err
		}

		players[userLogin] = NewPlayer(
			userLogin,
			name,
			wordsTried,
			wordsGuessed,
		)

		turnOrder = append(turnOrder, userLogin)
		if currentPlayerId == userLogin {
			currentPlayer = i
		}
		i++
	}

	r := NewRoom(id, admin, cfg, players, turnOrder, currentPlayer, currentWordIndex, GameState(gameState), p.logger)
	return r, nil
}

// SaveRoom tries to save `Room` with its users and turn order.
func (p *Postgres) SaveRoom(ctx context.Context, r *Room) (err error) {
	const op = "database.SaveRoom"

	tx, err := p.db.Begin(ctx)
	if err != nil {
		p.logger.Error("failed to begin transaction", "err", err, "roomId", r.Id)
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			p.logger.Error("failed to rollback transaction", "err", rollbackErr, "roomId", r.Id)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = commitErr
			p.logger.Error("failed to commit transaction", "err", err, "roomId", r.Id)
		}
	}()

	roomQuery := `
		INSERT INTO rooms (id, admin, seed, current_word_index, current_player_login, game_state, language, rude_words, additional_vocabulary, clock)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id)
		DO UPDATE SET
			current_word_index    = EXCLUDED.current_word_index,
			current_player_login     = EXCLUDED.current_player_login,
			game_state            = EXCLUDED.game_state,
			language              = EXCLUDED.language,
			rude_words            = EXCLUDED.rude_words,
			additional_vocabulary = EXCLUDED.additional_vocabulary,
			clock                 = EXCLUDED.clock`

	_, err = tx.Exec(ctx, roomQuery,
		r.Id,
		r.Admin,
		r.Config.Seed,
		r.CurrentWordIndex,
		r.CurrentPlayer(),
		r.State,
		r.Config.Language,
		r.Config.RudeWords,
		r.Config.AdditionalVocabulary,
		r.Config.Clock,
	)
	if err != nil {
		p.logger.Error("failed to upsert rooms", "err", err, "roomId", r.Id)
		return err
	}

	participantQuery := `
		INSERT INTO room_participants (room_id, user_login, words_tried, words_guessed, turn_order)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (room_id, user_login)
		DO UPDATE SET
			words_tried   = EXCLUDED.words_tried,
			words_guessed = EXCLUDED.words_guessed,
			turn_order    = EXCLUDED.turn_order`

	batch := &pgx.Batch{}
	for i, playerId := range r.TurnOrder {
		player := r.Players[playerId]
		batch.Queue(participantQuery, r.Id, player.Id, player.WordsTried, player.WordsGuessed, i)
	}

	results := tx.SendBatch(ctx, batch)
	for _, player := range r.Players {
		if _, err = results.Exec(); err != nil {
			p.logger.Error("upsert participant error", "err", err, "roomId", r.Id, "playerId", player.Id)
			results.Close()
			return err
		}
	}
	if err = results.Close(); err != nil {
		p.logger.Error("failed to close batch results", "err", err, "roomId", r.Id)
		return err
	}
	return nil
}
