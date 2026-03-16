package main

import (
	"container/ring"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	mrand "math/rand"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPool() (*pgxpool.Pool, error) {
	postgresUrl := os.Getenv("POSTGRES_URL")
	if postgresUrl == "" {
		log.Fatal("POSTGRES_URL environment variable not set")
	}
	return pgxpool.New(context.Background(), postgresUrl)
}

type Postgres struct {
	db *pgxpool.Pool
}

func (p *Postgres) AvailableLanguages(ctx context.Context) ([]string, error) {
	query := "SELECT language FROM vocabularies"
	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return pgx.CollectRows(rows, pgx.RowTo[string])
}

type UserCredentials struct {
	Id     uuid.UUID `json:"id"`
	Secret string    `json:"secret"`
}

func generateSecret() (string, error) {
	b := [16]byte{}
	_, err := rand.Read(b[:])
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

func generateFunnyName() string {
	adjectives := []string{"Grumpy", "Sleepy", "Chaotic", "Spicy", "Wobbly", "Fluffy", "Sneaky"}
	nouns := []string{"Waffle", "Penguin", "Muffin", "Wizard", "Noodle", "Taco", "Biscuit"}
	return adjectives[mrand.Intn(len(adjectives))] + nouns[mrand.Intn(len(nouns))] + strconv.Itoa(mrand.Intn(100))
}

func (p *Postgres) CreateUser(ctx context.Context) (string, UserCredentials, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", UserCredentials{}, err
	}
	secret, err := generateSecret()
	if err != nil {
		return "", UserCredentials{}, err
	}
	hash := sha256.Sum256([]byte(secret))
	secretHash := hex.EncodeToString(hash[:])
	name := generateFunnyName()

	query := "INSERT INTO users (id, name, secret_hash) VALUES ($1, $2, $3)"
	_, err = p.db.Exec(ctx, query, id, name, secretHash)
	if err != nil {
		return "", UserCredentials{}, err
	}
	return name, UserCredentials{id, secret}, nil
}

func (p *Postgres) ValidateUser(ctx context.Context, credentials UserCredentials) error {
	hash := sha256.Sum256([]byte(credentials.Secret))
	hashedSecret := hex.EncodeToString(hash[:])

	query := "SELECT TRUE FROM users WHERE id=$1 AND secret_hash=$2"
	row := p.db.QueryRow(ctx, query, credentials.Id, hashedSecret)
	var exists bool
	err := row.Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.New("user not found")
	}
	return err
}

func generateRoomId() string {
	data := [5]byte{}
	_, _ = rand.Read(data[:]) // possible collision at ~1 million games.
	return base32.StdEncoding.EncodeToString(data[:])
}

func (p *Postgres) AddRoom(ctx context.Context, adminId uuid.UUID, language string, rudeWords bool, additionalVocabulary []string, clock int) (string, error) {
	fmt.Println(additionalVocabulary == nil)
	query := "INSERT INTO rooms (id, admin, seed, current_word_index, current_player_id, game_state, language, rude_words, additional_vocabulary, clock) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	roomId := generateRoomId()
	seed := mrand.Int31()
	currentWordIndex := 0
	currentPlayerId := adminId
	gameState := 0
	_, err := p.db.Exec(ctx, query, roomId, adminId, seed, currentWordIndex, currentPlayerId, gameState, language, rudeWords, additionalVocabulary, clock)
	if err != nil {
		return "", err
	}
	return roomId, nil
}

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
		return nil, err
	}

	vocabs.lock.RLock()
	wordsTotal := len(vocabs.vocabulary[language].PrimaryWords)
	if rudeWords {
		wordsTotal += len(vocabs.vocabulary[language].RudeWords)
	}
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

	query = "SELECT user_id, words_tried, words_guessed FROM room_participants WHERE room_id=$1 ORDER BY turn_order ASC"
	rows, err := p.db.Query(ctx, query, roomId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := map[uuid.UUID]*Player{}
	var queue *ring.Ring
	var currentPlayer *ring.Ring

	for rows.Next() {
		var userId string
		var wordsTried int
		var wordsGuessed int
		err = rows.Scan(&userId, &wordsTried, &wordsGuessed)
		if err != nil {
			return nil, err
		}

		id := uuid.MustParse(userId)
		players[id] = &Player{
			id,
			make(chan []byte, 10),
			false,
			wordsTried,
			wordsGuessed,
		}

		q := ring.New(1)
		q.Value = id
		if queue == nil {
			queue = q
		} else {
			queue.Prev().Link(q)
		}

		if id == currentPlayerId {
			currentPlayer = q
		}
	}

	if currentPlayer == nil {
		currentPlayer = queue
	}
	room := &Room{
		Id:     id,
		Admin:  admin,
		Config: cfg,

		Players: players,
		ingest:  make(chan *ClientMessage, 50),
		count:   0,
		join:    make(chan *Player, 5),
		leave:   make(chan uuid.UUID, 5),

		currentPlayer: currentPlayer,
		currentWord:   currentWordIndex,
		State:         GameState(gameState),

		// only state loaded from db is RoundOver, so no need for ticker
		ticker:        &time.Ticker{},
		RemainingTime: cfg.Clock,
	}
	return room, nil
}

func (p *Postgres) LoadVocabs(ctx context.Context) (map[string]Vocabulary, error) {
	vocabs := map[string]Vocabulary{}
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
			return vocabs, err
		}
		vocabs[language] = Vocabulary{primaryWords, rudeWords}
	}
	return vocabs, nil
}

func (p *Postgres) UpdateRoomState(ctx context.Context, r *Room) error {
	query := "UPDATE rooms SET current_word_index=$1, current_player_id=$2, game_state=$3 WHERE id=$4"

	playerId, ok := r.currentPlayer.Value.(uuid.UUID)
	if !ok {
		panic("invalid player type")
	}

	_, err := p.db.Exec(ctx, query, r.currentWord, playerId, r.State, r.Id)
	if err != nil {
		return err
	}

	turnOrder := make(map[uuid.UUID]int, r.currentPlayer.Len())
	head := r.currentPlayer
	for i := 0; i < head.Len(); i++ {
		if head.Value.(uuid.UUID) == r.Admin {
			break
		}
		head = head.Next()
	}
	for i := 0; i < head.Len(); i++ {
		id := head.Move(i).Value.(uuid.UUID)
		turnOrder[id] = i
	}

	query = `INSERT INTO room_participants (room_id, user_id, words_tried, words_guessed, turn_order)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (room_id, user_id)
        DO UPDATE SET
            words_tried    = EXCLUDED.words_tried,
            words_guessed  = EXCLUDED.words_guessed,
            turn_order     = EXCLUDED.turn_order;`

	batch := &pgx.Batch{}
	for _, pl := range r.Players {
		batch.Queue(query, r.Id, pl.Id, pl.WordsTried, pl.WordsGuessed, turnOrder[pl.Id])
	}

	results := p.db.SendBatch(ctx, batch)
	defer results.Close()

	for range r.Players {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("batch upsert participant: %w", err)
		}
	}
	return nil
}
