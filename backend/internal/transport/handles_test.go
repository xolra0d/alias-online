package transport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/xolra0d/alias-online/internal/room"
)

func TestAvailableLanguagesSortsEnglishFirstThenAscending(t *testing.T) {
	h := &Handles{
		vocabs: room.NewVocabularies(map[string]*room.Vocabulary{
			"Zulu":    {},
			"French":  {},
			"English": {},
			"Arabic":  {},
		}),
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/available-vocabs", nil)

	h.AvailableLanguages(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var payload struct {
		Languages []string `json:"languages"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	want := []string{"English", "Arabic", "French", "Zulu"}
	if !reflect.DeepEqual(payload.Languages, want) {
		t.Fatalf("expected languages %v, got %v", want, payload.Languages)
	}
}
