package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	mrand "math/rand"
	"os"
	"strconv"

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
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
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
	num := mrand.Uint32()
	return strconv.FormatUint(uint64(num), 36)
}

func (p *Postgres) AddRoom(ctx context.Context, adminId uuid.UUID, language string, rudeWords bool, additionalVocabulary []string, clock int) (string, error) {
	query := "INSERT INTO rooms (id, admin, language, rude_words, additional_vocabulary, clock) VALUES ($1, $2, $3, $4, $5, $6)"
	roomId := generateRoomId()
	_, err := p.db.Exec(ctx, query, roomId, adminId, language, rudeWords, additionalVocabulary, clock)
	if err != nil {
		return "", err
	}
	return roomId, nil
}
