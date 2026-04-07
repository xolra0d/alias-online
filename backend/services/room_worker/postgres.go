package main

import (
	"context"
	"fmt"
	"log/slog"
	mrand "math/rand"

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

//func (p *Postgres) LoadRoomConfig(ctx context.Context, roomId string, getVocab func(string) (Vocabulary, error)) (*RoomConfig, error) {
//	var id string
//	var admin uuid.UUID
//	var seed int
//	var currentWordIndex int
//	var currentPlayerId uuid.UUID
//	var gameState int
//	var language string
//	var rudeWords bool
//	var additionalVocabulary []string
//	var clock int
//	query := "SELECT id, admin, seed, current_word_index, current_player_id, game_state, language, rude_words, additional_vocabulary, clock FROM rooms WHERE id=$1"
//	err := p.db.QueryRow(ctx, query, roomId).Scan(
//		&id,
//		&admin,
//		&seed,
//		&currentWordIndex,
//		&currentPlayerId,
//		&gameState,
//		&language,
//		&rudeWords,
//		&additionalVocabulary,
//		&clock,
//	)
//	if err != nil {
//		p.logger.Error("could not load room_worker", "roomId", roomId, "err", err)
//		return nil, err
//	}
//
//	vocab, err := getVocab(language)
//	if err != nil {
//		p.logger.Error("could not find vocab", "language", language)
//		return nil, fmt.Errorf("could not find vocab: %s", language)
//
//	}
//
//	allWords := vocab.PrimaryWords // it's our copy, so it's fine to mutate
//	if rudeWords {
//		allWords = append(allWords, vocab.RudeWords...)
//	}
//
//	rng := mrand.New(mrand.NewSource(int64(seed)))
//	rng.Shuffle(len(allWords), func(i, j int) {
//		allWords[i], allWords[j] = allWords[j], allWords[i]
//	})
//
//	cfg := &RoomConfig{
//		Seed:                 seed,
//		AllWords:             allWords,
//		Language:             language,
//		RudeWords:            rudeWords,
//		AdditionalVocabulary: additionalVocabulary,
//		Clock:                clock,
//	}
//	return cfg, nil
//}
//
//// LoadRoom tries to load `Room` with its users and turn order. Sets all users to not ready.
//func (p *Postgres) LoadRoomPlayers(ctx context.Context, roomId string, currentPlayerId uuid.UUID) (map[uuid.UUID]*Player, []uuid.UUID, int, error) {
//	const op = "database.LoadRoom"
//
//	query := `SELECT rp.user_id, rp.words_tried, rp.words_guessed, u.name
//		FROM room_participants rp
//		JOIN users u ON rp.user_id = u.id
//		WHERE rp.room_id = $1
//		ORDER BY rp.turn_order ASC`
//	rows, err := p.db.Query(ctx, query, roomId)
//	if err != nil {
//		p.logger.Error(op, "failed to load room_worker participants", "roomId", roomId, "err", err)
//		return nil, nil, 0, err
//	}
//	defer rows.Close()
//
//	players := map[uuid.UUID]*Player{}
//	//goland:noinspection GoPreferNilSlice
//	turnOrder := []uuid.UUID{}
//	currentPlayer := 0
//
//	i := 0
//	for rows.Next() {
//		var userId string
//		var wordsTried int
//		var wordsGuessed int
//		var name string
//		err = rows.Scan(&userId, &wordsTried, &wordsGuessed, &name)
//		if err != nil {
//			p.logger.Error(op, "failed to scan room_worker participant", "roomId", roomId, "err", err)
//			return nil, nil, 0, err
//		}
//
//		id := uuid.MustParse(userId)
//		players[id] = NewPlayer(
//			id,
//			name,
//			wordsTried,
//			wordsGuessed,
//		)
//
//		turnOrder = append(turnOrder, id)
//		if currentPlayerId == id {
//			currentPlayer = i
//		}
//		i++
//	}
//
//	return players, turnOrder, currentPlayer, nil
//}

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
	query := "SELECT id, admin, seed, current_word_index, current_player_id, game_state, language, rude_words, additional_vocabulary, clock FROM room_worker WHERE id=$1"
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
		p.logger.Error(op, "could not load room_worker", "roomId", roomId, "err", err)
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
		p.logger.Error(op, "failed to load room_worker participants", "roomId", roomId, "err", err)
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
			p.logger.Error(op, "failed to scan room_worker participant", "roomId", roomId, "err", err)
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

func (p *Postgres) SaveRoom(ctx context.Context, room *Room) error {
	return nil
}
