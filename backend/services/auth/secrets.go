package main

import (
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.pact.im/x/option"
	"go.pact.im/x/phcformat"
	"go.pact.im/x/phcformat/encode"
	"golang.org/x/crypto/argon2"
)

type Secrets struct {
	logger *slog.Logger

	argon2idTime    uint32
	argon2idMemory  uint32
	argon2idThreads uint8
	argon2idOutLen  uint32
	privateKey      *rsa.PrivateKey
}

// loads rsa private key from file.
func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := k.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("could not parse as rsa private key")
	}
	return key, nil
}

// NewSecrets creates new secrets manager.
func NewSecrets(
	logger *slog.Logger,
	Argon2idTime uint32,
	Argon2idMemory uint32,
	Argon2idThreads uint8,
	Argon2idOutLen uint32,
	RsaPrivateKeyFilename string,
) (*Secrets, error) {
	key, err := loadPrivateKey(RsaPrivateKeyFilename)
	if err != nil {
		return nil, err
	}
	return &Secrets{
		logger:          logger,
		argon2idTime:    Argon2idTime,
		argon2idMemory:  Argon2idMemory,
		argon2idThreads: Argon2idThreads,
		argon2idOutLen:  Argon2idOutLen,
		privateKey:      key,
	}, nil
}

const (
	hashArgon2idName = "argon2id"
	hashTimeKey      = "t"
	hashMemoryKey    = "m"
	hashProcessesKey = "p"
	hashLengthKey    = "l"
)

// hashSecret hashes using argon2id any password with random salt and returns result in phcformat string.
// More: https://github.com/P-H-C/phc-string-format/blob/master/phc-sf-spec.md
func (s *Secrets) hashSecret(secret string) string {
	salt := s.GenerateSecretBase32()
	hashed := argon2.IDKey([]byte(secret), []byte(salt), s.argon2idTime, s.argon2idMemory, s.argon2idThreads, s.argon2idOutLen)

	newParam := func(k string, v uint) encode.KV[encode.String, encode.Byte, encode.Uint] {
		return encode.NewKV(encode.NewByte('='), encode.NewString(k), encode.NewUint(v))
	}

	out := string(phcformat.Append(nil,
		encode.String(hashArgon2idName),
		option.Value(encode.NewUint(uint(argon2.Version))),
		option.Value(encode.NewList(
			encode.NewByte(','),
			newParam(hashTimeKey, uint(s.argon2idTime)),
			newParam(hashMemoryKey, uint(s.argon2idMemory)),
			newParam(hashProcessesKey, uint(s.argon2idThreads)),
			newParam(hashLengthKey, uint(s.argon2idOutLen)),
		)),
		option.Value(encode.NewBase64(salt)),
		option.Value(encode.NewBase64(hashed)),
	))

	return out
}

// verifyArgon2id checks that secret is the same as encoded one in phcformat hash.
func (s *Secrets) verifyArgon2id(secret string, hash phcformat.Hash) *Error {
	const op = "database.verifyArgon2id"

	paramsStr, ok := hash.Params.Unwrap()
	if !ok {
		s.logger.Error("params for hashing algorithm are not set", "op", op, "hash", hash.String())
		return NewError(ErrInternal, fmt.Errorf("params for hashing algorithm are not set"))
	}
	params := make(map[string]uint32, 4)
	for part := range strings.SplitSeq(paramsStr, ",") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			s.logger.Error("params for hashing algorithm are not in key-value format", "op", op, "hash", hash.String())
			return NewError(ErrInternal, fmt.Errorf("params for hashing algorithm are not in key-value format"))
		}

		v, err := strconv.ParseUint(kv[1], 10, 32)
		if err != nil {
			s.logger.Error("could not parse params for hashing algorithm", "op", op, "hash", hash.String(), "name", kv[0], "val", kv[1])
			return NewError(ErrInternal, fmt.Errorf("could not parse params for hashing algorithm"))
		}
		params[kv[0]] = uint32(v)
	}

	if !hasKeys(params, hashTimeKey, hashMemoryKey, hashProcessesKey, hashLengthKey) {
		s.logger.Error("missing one or more of critical params for hashing algorithm", "op", op, "hash", hash.String(), "params", params)
		return NewError(ErrInternal, fmt.Errorf("missing one or more of critical params for hashing algorithm"))
	}
	t := params[hashTimeKey]
	m := params[hashMemoryKey]
	p := params[hashProcessesKey]
	l := params[hashLengthKey]

	salt, ok := hash.Salt.Unwrap()
	if !ok {
		s.logger.Error("salt is missing", "op", op, "hash", hash.String())
		return NewError(ErrInternal, fmt.Errorf("salt is missing"))
	}
	saltDecoded, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		s.logger.Error("salt decode error", "op", op, "hash", hash.String(), "err", err)
		return NewError(ErrInternal, fmt.Errorf("salt decode error"))
	}
	expected, ok := hash.Output.Unwrap()
	if !ok {
		s.logger.Error("output is missing", "op", op, "hash", hash.String(), "err", err)
		return NewError(ErrInternal, fmt.Errorf("output is missing"))
	}
	received := argon2.IDKey([]byte(secret), saltDecoded, t, m, uint8(p), l)

	if expected != base64.RawStdEncoding.EncodeToString(received) {
		return NewError(ErrWrongCredentials, fmt.Errorf("wrong credentials"))
	}
	return nil
}

func hasKeys(a map[string]uint32, keys ...string) bool {
	for _, k := range keys {
		if _, ok := a[k]; !ok {
			return false
		}
	}
	return true
}

// VerifyPassword checks if secret is equal to hash's secret.
func (s *Secrets) VerifyPassword(secret, hash string) *Error {
	const op = "database.VerifyPassword"

	h, ok := phcformat.Parse(hash)
	if !ok {
		s.logger.Error("could not decode phcformat hash", "hash", hash, "op", op)
		return NewError(ErrInternal, fmt.Errorf("could not decode phcformat hash: %s", hash))
	}

	switch h.ID {
	case hashArgon2idName:
		return s.verifyArgon2id(secret, h)
	default:
		s.logger.Error("Invalid hash ID", "hash", hash)
		return NewError(ErrInternal, fmt.Errorf("invalid hash ID: %s", hash))
	}
}

// GenerateSecretBase32 creates secure base32 secret.
func (s *Secrets) GenerateSecretBase32() string {
	return crand.Text()
}

// ValidateName check if name is valid for name field.
func ValidateName(name string) error {
	if len(name) < 8 || len(name) > 20 {
		return fmt.Errorf("invalid name")
	}
	return nil
}

// ValidateLogin check if login is valid for login field.
func ValidateLogin(login string) error {
	if len(login) < 8 || len(login) > 20 {
		return fmt.Errorf("invalid login")
	}
	return nil
}

// ValidatePassword check if password is valid for password field.
func ValidatePassword(password string) error {
	if len(password) < 8 || len(password) > 20 {
		return fmt.Errorf("invalid password")
	}
	return nil
}

// ValidateForRegister checks if name, login and password are valid for according fields in db.
func ValidateForRegister(name, login, password string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	if err := ValidateLogin(login); err != nil {
		return err
	}
	if err := ValidatePassword(password); err != nil {
		return err
	}
	return nil
}

// ValidateForLogin checks if login and password are valid for according fields in db.
func ValidateForLogin(login, password string) error {
	if err := ValidateLogin(login); err != nil {
		return err
	}
	if err := ValidatePassword(password); err != nil {
		return err
	}
	return nil
}

// NewJWT Issues new JWT.
func (s *Secrets) NewJWT(login string, exp time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub": login,
		"iat": time.Now().Unix(),
		"exp": exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	t, err := token.SignedString(s.privateKey)
	if err != nil {
		s.logger.Error("error signing token", "claims", claims, "err", err)
		return "", fmt.Errorf("jwt error")
	}
	return t, nil
}
