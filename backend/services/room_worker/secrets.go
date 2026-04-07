package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Secrets struct {
	jwtPublicToken *rsa.PublicKey

	logger *slog.Logger
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	k, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := k.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("could not parse as rsa private key")
	}
	return key, nil
}

func NewSecrets(jwtPublicTokenPath string, logger *slog.Logger) (*Secrets, error) {
	publicKey, err := loadPublicKey(jwtPublicTokenPath)
	if err != nil {
		return nil, err
	}
	return &Secrets{
		jwtPublicToken: publicKey,
		logger:         logger,
	}, nil
}

func (s *Secrets) CheckJwtToken(tokenString string) (string, error) {
	var ErrUnAuthorized = fmt.Errorf("unauthorized")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			s.logger.Error("unexpected signing method", "err", t.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return &s.jwtPublicToken, nil
	})
	if err != nil || token == nil || !token.Valid {
		s.logger.Warn("invalid token", "err", err, "token", tokenString)
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrUnAuthorized
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", ErrUnAuthorized
	}

	exp, ok := claims["exp"].(float64)
	if !ok || exp == 0 || int64(exp) < time.Now().Unix() {
		return "", ErrUnAuthorized
	}

	return sub, nil
}
