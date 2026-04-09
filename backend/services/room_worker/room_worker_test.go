package main

import (
	"context"
	"maps"
	"testing"
	"time"
)

func TestPrepareStateWaitUntilOperational(t *testing.T) {
	state := NewPrepareState()
	done := make(chan error, 1)

	go func() {
		done <- state.WaitUntilOperational()
	}()

	time.Sleep(10 * time.Millisecond)
	state.SetOperational()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("WaitUntilOperational did not unblock")
	}
}

func TestPrepareStateWaitUntilErrored(t *testing.T) {
	state := NewPrepareState()
	done := make(chan error, 1)

	go func() {
		done <- state.WaitUntilOperational()
	}()

	time.Sleep(10 * time.Millisecond)
	state.SetErrored()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	case <-time.After(time.Second):
		t.Fatal("WaitUntilOperational did not unblock")
	}
}

func TestMapToRoomConfigSuccess(t *testing.T) {
	ctx := context.Background()
	input := map[string]any{
		"rude-words":            true,
		"additional-vocabulary": []any{"extra"},
		"clock":                 float64(30),
		"language":              "en",
	}

	getVocab := func(_ context.Context, name string) (Vocabulary, error) {
		if name != "en" {
			t.Fatalf("unexpected vocab name: %s", name)
		}
		return Vocabulary{
			PrimaryWords: []string{"p1", "p2"},
			RudeWords:    []string{"r1"},
		}, nil
	}

	cfg, err := mapToRoomConfig(ctx, input, 60, getVocab)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Language != "en" {
		t.Fatalf("expected language en, got %s", cfg.Language)
	}
	if !cfg.RudeWords {
		t.Fatal("expected rude words to be enabled")
	}
	if cfg.Clock != 30 {
		t.Fatalf("expected clock 30, got %d", cfg.Clock)
	}
	if len(cfg.AdditionalVocabulary) != 1 || cfg.AdditionalVocabulary[0] != "extra" {
		t.Fatalf("unexpected additional vocabulary: %#v", cfg.AdditionalVocabulary)
	}

	if len(cfg.AllWords) != 4 {
		t.Fatalf("expected 4 words, got %d", len(cfg.AllWords))
	}
	expected := map[string]int{"p1": 1, "p2": 1, "r1": 1, "extra": 1}
	actual := map[string]int{}
	for _, w := range cfg.AllWords {
		actual[w]++
	}
	if !maps.Equal(expected, actual) {
		t.Fatalf("all words mismatch, expected %#v, got %#v", expected, actual)
	}
}

func TestMapToRoomConfigInvalidClock(t *testing.T) {
	ctx := context.Background()
	input := map[string]any{
		"rude-words":            false,
		"additional-vocabulary": []any{},
		"clock":                 float64(999),
		"language":              "en",
	}

	getVocab := func(_ context.Context, _ string) (Vocabulary, error) {
		return Vocabulary{}, nil
	}

	_, err := mapToRoomConfig(ctx, input, 60, getVocab)
	if err == nil {
		t.Fatal("expected error for invalid clock")
	}
}

func TestMapToRoomConfigInvalidAdditionalVocabulary(t *testing.T) {
	ctx := context.Background()
	input := map[string]any{
		"rude-words":            false,
		"additional-vocabulary": []any{"ok", 123},
		"clock":                 float64(10),
		"language":              "en",
	}

	getVocab := func(_ context.Context, _ string) (Vocabulary, error) {
		return Vocabulary{}, nil
	}

	_, err := mapToRoomConfig(ctx, input, 60, getVocab)
	if err == nil {
		t.Fatal("expected error for invalid additional vocabulary item")
	}
}

func TestToClientMessageSuccess(t *testing.T) {
	input := map[string]any{
		"user_id": "u1",
		"type":    float64(GetState),
		"data":    map[string]any{"k": "v"},
	}

	msg, err := toClientMessage(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if msg.UserId != "u1" {
		t.Fatalf("expected user_id u1, got %s", msg.UserId)
	}
	if msg.MsgType != GetState {
		t.Fatalf("expected type %d, got %d", GetState, msg.MsgType)
	}
	if v := msg.MsgData["k"]; v != "v" {
		t.Fatalf("expected data[k] to be v, got %#v", v)
	}
}

func TestToClientMessageInvalidUserID(t *testing.T) {
	input := map[string]any{
		"user_id": 1,
		"type":    float64(GetState),
		"data":    map[string]any{},
	}

	_, err := toClientMessage(input)
	if err == nil {
		t.Fatal("expected error for invalid user_id")
	}
}

func TestRoomWordHelpers(t *testing.T) {
	r := &Room{
		Config: &RoomConfig{
			AllWords: []string{"one", "two"},
		},
		CurrentWordIndex: 0,
	}

	if got := r.CurrentWord(); got != "one" {
		t.Fatalf("expected current word one, got %s", got)
	}
	if got := r.NextWord(); got != "two" {
		t.Fatalf("expected next word two, got %s", got)
	}
	if got := r.getWordAt(-1); got != "" {
		t.Fatalf("expected empty word for out-of-range index, got %s", got)
	}
	if got := r.getWordAt(2); got != "" {
		t.Fatalf("expected empty word for out-of-range index, got %s", got)
	}
}
