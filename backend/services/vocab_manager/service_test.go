package main

import (
	"context"
	"reflect"
	"sort"
	"testing"

	pb "github.com/xolra0d/alias-online/shared/proto/vocab_manager"
	"google.golang.org/protobuf/types/known/emptypb"
)

func newTestVocabManager(vocabs map[string]Vocabulary) *VocabManager {
	return &VocabManager{vocabs: vocabs}
}

func TestVocabManagerAvailableVocabs(t *testing.T) {
	vm := newTestVocabManager(map[string]Vocabulary{
		"en": {},
		"ua": {},
		"de": {},
	})

	got := vm.AvailableVocabs()
	sort.Strings(got)
	want := []string{"de", "en", "ua"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected available vocabs: got=%v want=%v", got, want)
	}
}

func TestVocabManagerVocabReturnsCopy(t *testing.T) {
	vm := newTestVocabManager(map[string]Vocabulary{
		"en": {
			PrimaryWords: []string{"apple", "banana"},
			RudeWords:    []string{"foo"},
		},
	})

	got := vm.Vocab("en")
	got.PrimaryWords[0] = "changed"
	got.RudeWords[0] = "changed"

	original := vm.vocabs["en"]
	if original.PrimaryWords[0] != "apple" {
		t.Fatalf("primary words were mutated in manager: %v", original.PrimaryWords)
	}
	if original.RudeWords[0] != "foo" {
		t.Fatalf("rude words were mutated in manager: %v", original.RudeWords)
	}
}

func TestVocabManagerVocabReturnsEmptyForUnknown(t *testing.T) {
	vm := newTestVocabManager(map[string]Vocabulary{
		"en": {
			PrimaryWords: []string{"apple"},
		},
	})

	got := vm.Vocab("missing")
	if len(got.PrimaryWords) != 0 || len(got.RudeWords) != 0 {
		t.Fatalf("expected empty vocabulary for unknown name, got=%+v", got)
	}
}

func TestServerPing(t *testing.T) {
	s := &server{vocabs: newTestVocabManager(map[string]Vocabulary{})}
	got, err := s.Ping(context.Background(), &emptypb.Empty{})
	if err != nil {
		t.Fatalf("Ping returned error: %v", err)
	}
	if !got.Ok {
		t.Fatal("Ping response should have Ok=true")
	}
}

func TestServerGetAvailableVocabs(t *testing.T) {
	s := &server{vocabs: newTestVocabManager(map[string]Vocabulary{
		"en": {},
		"ua": {},
	})}

	got, err := s.GetAvailableVocabs(context.Background(), &emptypb.Empty{})
	if err != nil {
		t.Fatalf("GetAvailableVocabs returned error: %v", err)
	}

	sort.Strings(got.Names)
	want := []string{"en", "ua"}
	if !reflect.DeepEqual(got.Names, want) {
		t.Fatalf("unexpected names: got=%v want=%v", got.Names, want)
	}
}

func TestServerGetVocab(t *testing.T) {
	s := &server{vocabs: newTestVocabManager(map[string]Vocabulary{
		"en": {
			PrimaryWords: []string{"apple", "banana"},
			RudeWords:    []string{"foo"},
		},
	})}

	got, err := s.GetVocab(context.Background(), &pb.GetVocabRequest{Name: "en"})
	if err != nil {
		t.Fatalf("GetVocab returned error: %v", err)
	}

	wantPrimary := []string{"apple", "banana"}
	wantRude := []string{"foo"}
	if !reflect.DeepEqual(got.PrimaryWords, wantPrimary) {
		t.Fatalf("unexpected primary words: got=%v want=%v", got.PrimaryWords, wantPrimary)
	}
	if !reflect.DeepEqual(got.RudeWords, wantRude) {
		t.Fatalf("unexpected rude words: got=%v want=%v", got.RudeWords, wantRude)
	}
}
