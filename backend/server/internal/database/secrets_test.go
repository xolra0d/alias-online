package database

import (
	"io"
	"strings"
	"testing"

	"github.com/xolra0d/alias-online/internal/config"
)

func testSecrets() *Secrets {
	return NewSecrets(config.NewLogger(64, io.Discard), 2, 64*1024, 1, 32)
}

func TestHashAndVerifyPassword(t *testing.T) {
	s := testSecrets()
	plain := "correct-password"

	hash := s.hashSecret(plain)
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if !s.VerifyPassword(plain, hash) {
		t.Fatal("expected password verification to pass")
	}
	if s.VerifyPassword("wrong-password", hash) {
		t.Fatal("expected password verification to fail for wrong password")
	}
}

func TestVerifyPasswordInvalidHash(t *testing.T) {
	s := testSecrets()
	if s.VerifyPassword("password", "not-a-valid-phc-hash") {
		t.Fatal("expected invalid hash verification to fail")
	}
}

func TestGenerateRoomIDUsesBase32Alphabet(t *testing.T) {
	s := testSecrets()
	roomID := s.GenerateRoomId()

	if len(roomID) != 8 {
		t.Fatalf("expected room id length 8, got %d (%q)", len(roomID), roomID)
	}
	if strings.Trim(roomID, "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567") != "" {
		t.Fatalf("room id contains invalid chars: %q", roomID)
	}
}
