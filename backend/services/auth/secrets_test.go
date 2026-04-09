package main

import (
	"crypto/rand"
	"crypto/rsa"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func testSecrets(t *testing.T) *Secrets {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate rsa key: %v", err)
	}

	return &Secrets{
		logger:          testLogger(),
		argon2idTime:    2,
		argon2idMemory:  64 * 1024,
		argon2idThreads: 1,
		argon2idOutLen:  32,
		privateKey:      key,
	}
}

func TestErrorError(t *testing.T) {
	var nilErr *Error
	if got := nilErr.Error(); got != "" {
		t.Fatalf("nil error should render empty string, got %q", got)
	}

	err := NewError(ErrUserNotFound, io.EOF)
	if got := err.Error(); got != string(ErrUserNotFound) {
		t.Fatalf("unexpected service error string: got %q", got)
	}
}

func TestHasKeys(t *testing.T) {
	m := map[string]uint32{"a": 1, "b": 2}
	if !hasKeys(m, "a", "b") {
		t.Fatal("expected keys to be present")
	}
	if hasKeys(m, "a", "missing") {
		t.Fatal("expected missing key check to fail")
	}
}

func TestValidateFieldsBoundaries(t *testing.T) {
	valid8 := "12345678"
	valid20 := "12345678901234567890"
	invalidShort := "1234567"
	invalidLong := "123456789012345678901"

	cases := []struct {
		name  string
		fn    func(string) error
		value string
		want  bool
	}{
		{"name valid min", ValidateName, valid8, false},
		{"name valid max", ValidateName, valid20, false},
		{"name invalid short", ValidateName, invalidShort, true},
		{"name invalid long", ValidateName, invalidLong, true},
		{"login valid min", ValidateLogin, valid8, false},
		{"login valid max", ValidateLogin, valid20, false},
		{"login invalid short", ValidateLogin, invalidShort, true},
		{"login invalid long", ValidateLogin, invalidLong, true},
		{"password valid min", ValidatePassword, valid8, false},
		{"password valid max", ValidatePassword, valid20, false},
		{"password invalid short", ValidatePassword, invalidShort, true},
		{"password invalid long", ValidatePassword, invalidLong, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.fn(tc.value)
			if tc.want && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.want && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}

func TestValidateForRegisterAndLogin(t *testing.T) {
	if err := ValidateForRegister("username1", "loginname", "password"); err != nil {
		t.Fatalf("expected valid register payload, got %v", err)
	}
	if err := ValidateForRegister("short", "loginname", "password"); err == nil {
		t.Fatal("expected invalid register payload to fail")
	}

	if err := ValidateForLogin("loginname", "password"); err != nil {
		t.Fatalf("expected valid login payload, got %v", err)
	}
	if err := ValidateForLogin("short", "password"); err == nil {
		t.Fatal("expected invalid login payload to fail")
	}
}

func TestHashAndVerifyPassword(t *testing.T) {
	s := testSecrets(t)
	secret := "correct-password"

	hash := s.hashSecret(secret)
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}

	if err := s.VerifyPassword(secret, hash); err != nil {
		t.Fatalf("expected password verification to pass, got %v", err)
	}

	err := s.VerifyPassword("wrong-password", hash)
	if err == nil || err.SvcError != ErrWrongCredentials {
		t.Fatalf("expected wrong credentials error, got %#v", err)
	}
}

func TestVerifyPasswordInvalidHash(t *testing.T) {
	s := testSecrets(t)

	err := s.VerifyPassword("password", "not-a-phc-hash")
	if err == nil || err.SvcError != ErrInternal {
		t.Fatalf("expected internal error for invalid hash format, got %#v", err)
	}
}

func TestNewJWT(t *testing.T) {
	s := testSecrets(t)
	exp := time.Now().Add(10 * time.Minute).Round(0)

	token, err := s.NewJWT("loginname", exp)
	if err != nil {
		t.Fatalf("failed to create jwt: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return &s.privateKey.PublicKey, nil
	})
	if err != nil {
		t.Fatalf("failed to parse signed jwt: %v", err)
	}
	if !parsed.Valid {
		t.Fatal("expected valid jwt")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("unexpected claims type: %T", parsed.Claims)
	}
	if claims["sub"] != "loginname" {
		t.Fatalf("expected sub claim to equal loginname, got %v", claims["sub"])
	}

	gotExp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatalf("expected numeric exp claim, got %T", claims["exp"])
	}
	if int64(gotExp) != exp.Unix() {
		t.Fatalf("unexpected exp claim: got %d want %d", int64(gotExp), exp.Unix())
	}
}
