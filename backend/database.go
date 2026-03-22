package main

import (
	"context"
	"errors"
	"fmt"
	mrand "math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitPool initializes PostgreSQL pool from `postgresUrl`.
func InitPool(postgresUrl string) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), postgresUrl)
}

// Postgres helps with postgres-specific commands.
type Postgres struct {
	db      *pgxpool.Pool
	secrets *Secrets
	logger  *PrefixLogger
}

// UserCredentials stores user credentials sent while creating user.
type UserCredentials struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Secret string    `json:"secret"`
}

// CreateUser creates random user and inserts into db.
func (p *Postgres) CreateUser(ctx context.Context) (UserCredentials, error) { // TODO: split into 2 function (generation and insertion).
	id, err := uuid.NewV7()
	if err != nil {
		p.logger.Error("uuidV7 generation error", "err", err)
		return UserCredentials{}, err
	}
	secret := p.secrets.GenerateSecretBase32()
	secretHash := p.secrets.hashSecret(secret)
	name := p.secrets.GenerateName()

	query := "INSERT INTO users (id, name, secret_hash) VALUES ($1, $2, $3)"
	_, err = p.db.Exec(ctx, query, id, name, secretHash)
	if err != nil {
		p.logger.Error("user insertion error", "err", err)
		return UserCredentials{}, err
	}
	return UserCredentials{id, name, secret}, nil
}

// ValidateUser tries to check if user with `credentials.Id` and `credentials.Secret` exists.
func (p *Postgres) ValidateUser(ctx context.Context, credentials UserCredentials) bool {
	query := "SELECT secret_hash FROM users WHERE id=$1"
	row := p.db.QueryRow(ctx, query, credentials.Id)
	var storedHash string
	err := row.Scan(&storedHash)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			p.logger.Error("hash retrieval error", "err", err)
		}
		return false
	}
	return p.secrets.VerifyPassword(credentials.Secret, storedHash)
}

// AddRoom generates seeds vocab and saves config to database.
func (p *Postgres) AddRoom(ctx context.Context, adminId uuid.UUID, cfg RoomConfig) (string, error) {
	query := "INSERT INTO rooms (id, admin, seed, current_word_index, current_player_id, game_state, language, rude_words, additional_vocabulary, clock) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	roomId := p.secrets.GenerateRoomId()
	seed := mrand.Int31()
	currentWordIndex := 0
	currentPlayerId := adminId
	gameState := int(RoundOver)
	_, err := p.db.Exec(ctx, query, roomId, adminId, seed, currentWordIndex, currentPlayerId, gameState, cfg.Language, cfg.RudeWords, cfg.AdditionalVocabulary, cfg.Clock)
	if err != nil {
		p.logger.Error("could not save new room", "err", err)
		return "", err
	}
	return roomId, nil
}

// LoadRoom tries to load `Room` with its users and turn order. Sets all users to not ready.
func (p *Postgres) LoadRoom(ctx context.Context, roomId string, vocabs *Vocabularies) (*Room, error) {
	var id string
	var admin uuid.UUID
	var seed int
	var currentWordIndex int
	var currentPlayerId uuid.UUID
	var gameState int
	var language string
	var rudeWords bool
	var additionalVocabulary []string
	var clock int
	query := "SELECT id, admin, seed, current_word_index, current_player_id, game_state, language, rude_words, additional_vocabulary, clock FROM rooms WHERE id=$1"
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
		p.logger.Error("could not load room", "roomId", roomId, "err", err)
		return nil, err
	}

	vocabs.lock.RLock()
	wordsTotal := len(vocabs.vocabulary[language].PrimaryWords)
	if rudeWords {
		wordsTotal += len(vocabs.vocabulary[language].RudeWords)
	}
	wordsTotal += len(additionalVocabulary)
	vocabs.lock.RUnlock()

	words := mrand.New(mrand.NewSource(int64(seed))).Perm(wordsTotal)

	cfg := &RoomConfig{
		seed,
		words,
		language,
		rudeWords,
		additionalVocabulary,
		clock,
	}

	query = `SELECT rp.user_id, rp.words_tried, rp.words_guessed, u.name
		FROM room_participants rp
		JOIN users u ON rp.user_id = u.id
		WHERE rp.room_id = $1
		ORDER BY rp.turn_order ASC`
	rows, err := p.db.Query(ctx, query, roomId)
	if err != nil {
		p.logger.Error("failed to load room participants", "roomId", roomId, "err", err)
		return nil, err
	}
	defer rows.Close()

	players := map[uuid.UUID]*Player{}
	//goland:noinspection GoPreferNilSlice
	turnOrder := []uuid.UUID{}
	currentPlayer := 0

	i := 0
	for rows.Next() {
		var userId string
		var wordsTried int
		var wordsGuessed int
		var name string
		err = rows.Scan(&userId, &wordsTried, &wordsGuessed, &name)
		if err != nil {
			p.logger.Error("failed to scan room participant", "roomId", roomId, "err", err)
			return nil, err
		}

		id := uuid.MustParse(userId)
		players[id] = &Player{
			id,
			name,
			make(chan []byte, 10),
			false,
			wordsTried,
			wordsGuessed,
		}

		turnOrder = append(turnOrder, id)
		if currentPlayerId == id {
			currentPlayer = i
		}
		i++
	}

	room := &Room{
		Id:     id,
		Admin:  admin,
		Config: cfg,

		Players:    players,
		ingest:     make(chan *ClientMessage, 50),
		readyCount: 0,
		join:       make(chan *Player, 5),
		leave:      make(chan uuid.UUID, 5),

		turnOrder:     turnOrder,
		currentPlayer: currentPlayer,
		currentWord:   currentWordIndex,
		State:         GameState(gameState),

		// only state loaded from db is RoundOver, so no need for ticker
		ticker:        &time.Ticker{},
		RemainingTime: cfg.Clock,
		logger:        p.logger.CopyWithPrefix(fmt.Sprintf("ROOM-%s", roomId)),
	}
	return room, nil
}

// LoadVocabs tries to load all vocabs from database.
func (p *Postgres) LoadVocabs(ctx context.Context) (map[string]*Vocabulary, error) {
	vocabs := map[string]*Vocabulary{}
	query := "SELECT language, primary_words, rude_words FROM vocabularies WHERE available = TRUE"
	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return vocabs, err
	}
	defer rows.Close()
	for rows.Next() {
		var language string
		var primaryWords []string
		var rudeWords []string
		if err := rows.Scan(&language, &primaryWords, &rudeWords); err != nil {
			p.logger.Error("failed to load vocabulary", "err", err)
			return vocabs, err
		}
		vocabs[language] = &Vocabulary{primaryWords, rudeWords}
	}
	return vocabs, nil
}

// UpdateRoomState updates room state and state of each player in who were in room.
func (p *Postgres) UpdateRoomState(ctx context.Context, r *Room) error {
	query := "UPDATE rooms SET current_word_index=$1, current_player_id=$2, game_state=$3 WHERE id=$4"

	r.State = RoundOver

	_, err := p.db.Exec(ctx, query, r.currentWord, r.turnOrder[r.currentPlayer], r.State, r.Id)
	if err != nil {
		p.logger.Error("failed to update room state", "err", err, "roomId", r.Id)
		return err
	}

	query = `INSERT INTO room_participants (room_id, user_id, words_tried, words_guessed, turn_order)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (room_id, user_id)
        DO UPDATE SET
            words_tried    = EXCLUDED.words_tried,
            words_guessed  = EXCLUDED.words_guessed,
            turn_order     = EXCLUDED.turn_order;`

	batch := &pgx.Batch{}
	for i, id := range r.turnOrder {
		p := r.Players[id]
		batch.Queue(query, r.Id, id, p.WordsTried, p.WordsGuessed, i)
	}

	results := p.db.SendBatch(ctx, batch)
	defer func(results pgx.BatchResults) {
		err := results.Close()
		if err != nil {
			p.logger.Error("failed to close batch results", "err", err)
		}
	}(results)

	var resErr error
	for _, id := range r.turnOrder {
		if _, err := results.Exec(); err != nil {
			p.logger.Error("upsert participant error", "err", err, "roomId", r.Id, "playerId", id)
			resErr = err
		}
	}
	return resErr
}
