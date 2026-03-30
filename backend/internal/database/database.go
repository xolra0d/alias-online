package database

import (
	"context"
	"errors"
	"fmt"
	mrand "math/rand"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xolra0d/alias-online/internal/config"
	"github.com/xolra0d/alias-online/internal/room"
)

// InitPool initializes PostgreSQL pool from `postgresUrl`.
func InitPool(postgresUrl string) (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), postgresUrl)
}

// Postgres helps with postgres-specific commands.
type Postgres struct {
	db      *pgxpool.Pool
	secrets *Secrets
	logger  *config.Logger
}

func NewPostgres(db *pgxpool.Pool, secrets *Secrets, logger *config.Logger) *Postgres {
	return &Postgres{
		db:      db,
		secrets: secrets,
		logger:  logger,
	}
}

// UserCredentials stores user credentials sent while creating user.
type UserCredentials struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Secret string    `json:"secret"`
}

// CreateUser creates random user and inserts into db.
func (p *Postgres) CreateUser(ctx context.Context) (UserCredentials, error) { // TODO: split into 2 function (generation and insertion).
	const op = "database.CreateUser"

	id, err := uuid.NewV7()
	if err != nil {
		p.logger.Error(op, "uuidV7 generation error", "err", err)
		return UserCredentials{}, err
	}
	secret := p.secrets.GenerateSecretBase32()
	secretHash := p.secrets.hashSecret(secret)
	name := p.secrets.GenerateName()

	query := "INSERT INTO users (id, name, secret_hash) VALUES ($1, $2, $3)"
	_, err = p.db.Exec(ctx, query, id, name, secretHash)
	if err != nil {
		p.logger.Error(op, "user insertion error", "err", err)
		return UserCredentials{}, err
	}
	return UserCredentials{id, name, secret}, nil
}

// ValidateUser tries to check if user with `credentials.Id` and `credentials.Secret` exists.
func (p *Postgres) ValidateUser(ctx context.Context, credentials UserCredentials) bool {
	const op = "database.ValidateUser"

	query := "SELECT secret_hash FROM users WHERE id=$1"
	row := p.db.QueryRow(ctx, query, credentials.Id)
	var storedHash string
	err := row.Scan(&storedHash)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			p.logger.Error(op, "hash retrieval error", "err", err)
		}
		return false
	}
	return p.secrets.VerifyPassword(credentials.Secret, storedHash)
}

// AddRoom generates seeds vocab and saves config to database.
func (p *Postgres) AddRoom(ctx context.Context, adminId uuid.UUID, cfg room.RoomConfig) (string, error) {
	const op = "database.AddRoom"

	query := "INSERT INTO rooms (id, admin, seed, current_word_index, current_player_id, game_state, language, rude_words, additional_vocabulary, clock) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	roomId := p.secrets.GenerateRoomId()
	seed := mrand.Int31()
	currentWordIndex := 0
	currentPlayerId := adminId
	gameState := int(room.RoundOver)
	_, err := p.db.Exec(ctx, query, roomId, adminId, seed, currentWordIndex, currentPlayerId, gameState, cfg.Language, cfg.RudeWords, cfg.AdditionalVocabulary, cfg.Clock)
	if err != nil {
		p.logger.Error(op, "could not save new room", "err", err)
		return "", err
	}
	return roomId, nil
}

// LoadRoom tries to load `Room` with its users and turn order. Sets all users to not ready.
func (p *Postgres) LoadRoom(ctx context.Context, roomId string, vocabs *room.Vocabularies) (*room.Room, error) {
	const op = "database.LoadRoom"

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
		p.logger.Error(op, "could not load room", "roomId", roomId, "err", err)
		return nil, err
	}

	wordsTotal, ok := vocabs.WordsInVocab(language, rudeWords)
	if !ok {
		p.logger.Error(op, "could not find vocab", "language", language)
		return nil, fmt.Errorf("could not find vocab: %s", language)
	}
	wordsTotal += len(additionalVocabulary)

	words := mrand.New(mrand.NewSource(int64(seed))).Perm(wordsTotal)

	cfg := &room.RoomConfig{
		Seed:                 seed,
		WordsPerm:            words,
		Language:             language,
		RudeWords:            rudeWords,
		AdditionalVocabulary: additionalVocabulary,
		Clock:                clock,
	}

	query = `SELECT rp.user_id, rp.words_tried, rp.words_guessed, u.name
		FROM room_participants rp
		JOIN users u ON rp.user_id = u.id
		WHERE rp.room_id = $1
		ORDER BY rp.turn_order ASC`
	rows, err := p.db.Query(ctx, query, roomId)
	if err != nil {
		p.logger.Error(op, "failed to load room participants", "roomId", roomId, "err", err)
		return nil, err
	}
	defer rows.Close()

	players := map[uuid.UUID]*room.Player{}
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
			p.logger.Error(op, "failed to scan room participant", "roomId", roomId, "err", err)
			return nil, err
		}

		id := uuid.MustParse(userId)
		players[id] = room.NewPlayer(
			id,
			name,
			wordsTried,
			wordsGuessed,
		)

		turnOrder = append(turnOrder, id)
		if currentPlayerId == id {
			currentPlayer = i
		}
		i++
	}

	r := room.NewRoom(id, admin, cfg, players, turnOrder, currentPlayer, currentWordIndex, room.GameState(gameState), p.logger)
	return r, nil
}

// LoadVocabs tries to load all vocabs from database.
func (p *Postgres) LoadVocabs(ctx context.Context) (map[string]*room.Vocabulary, error) {
	const op = "database.LoadVocabs"

	vocabs := map[string]*room.Vocabulary{}
	query := "SELECT language, primary_words, rude_words FROM vocabularies WHERE available = TRUE"
	rows, err := p.db.Query(ctx, query)
	if err != nil {
		p.logger.Error(op, "failed to load vocabs", "err", err)
		return vocabs, err
	}
	defer rows.Close()
	for rows.Next() {
		var language string
		var primaryWords []string
		var rudeWords []string
		if err := rows.Scan(&language, &primaryWords, &rudeWords); err != nil {
			p.logger.Error(op, "failed to load vocabulary", "err", err)
			return vocabs, err
		}
		vocabs[language] = &room.Vocabulary{PrimaryWords: primaryWords, RudeWords: rudeWords}
	}
	return vocabs, nil
}

func (p *Postgres) SaveRoomSnapshot(ctx context.Context, r *room.Room) (err error) {
	const op = "database.SaveRoomSnapshot"

	tx, err := p.db.Begin(ctx)
	if err != nil {
		p.logger.Error(op, "failed to begin transaction", "err", err, "roomId", r.Id)
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			p.logger.Error(op, "failed to rollback transaction", "err", rollbackErr, "roomId", r.Id)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			err = commitErr
			p.logger.Error(op, "failed to commit transaction", "err", err, "roomId", r.Id)
		}
	}()

	roomQuery := `
		INSERT INTO rooms (id, admin, seed, current_word_index, current_player_id, game_state, language, rude_words, additional_vocabulary, clock)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id)
		DO UPDATE SET
			current_word_index    = EXCLUDED.current_word_index,
			current_player_id     = EXCLUDED.current_player_id,
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
		p.logger.Error(op, "failed to upsert room", "err", err, "roomId", r.Id)
		return err
	}

	participantQuery := `
		INSERT INTO room_participants (room_id, user_id, words_tried, words_guessed, turn_order)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (room_id, user_id)
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
			p.logger.Error(op, "upsert participant error", "err", err, "roomId", r.Id, "playerId", player.Id)
			results.Close()
			return err
		}
	}
	if err = results.Close(); err != nil {
		p.logger.Error(op, "failed to close batch results", "err", err, "roomId", r.Id)
		return err
	}
	return nil
}
