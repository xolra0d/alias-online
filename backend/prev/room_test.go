package room_worker

//
//import (
//	"testing"
//
//	"github.com/google/uuid"
//)
//
//func newRoomForMessageTests() (*Room, *Vocabularies, uuid.UUID, uuid.UUID) {
//	explainer := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
//	guesser := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
//
//	room_worker := &Room{
//		Id:    "TESTROOM",
//		Admin: explainer,
//		Config: &RoomConfig{
//			Language:             "English",
//			RudeWords:            false,
//			AdditionalVocabulary: []string{},
//			Clock:                60,
//			wordsPerm:            []int{0, 1},
//		},
//		Players: map[uuid.UUID]*Player{
//			explainer: {Id: explainer, Name: "Explainer", toSend: make(chan []byte, 32), Ready: true},
//			guesser:   {Id: guesser, Name: "Guesser", toSend: make(chan []byte, 32), Ready: true},
//		},
//		TurnOrder:        []uuid.UUID{explainer, guesser},
//		currentPlayer:    0,
//		CurrentWordIndex: 0,
//		wordShown:        false,
//		State:            Explaining,
//		logger:           DummyPrefixLogger(),
//	}
//
//	vocabs := &Vocabularies{
//		vocabulary: map[string]*Vocabulary{
//			"English": {
//				PrimaryWords: []string{"cat", "dog"},
//			},
//		},
//	}
//
//	return room_worker, vocabs, explainer, guesser
//}
//
//func TestHandleMessageTryGuessIncrementsScore(t *testing.T) {
//	room_worker, vocabs, explainer, guesser := newRoomForMessageTests()
//
//	room_worker.handleMessage(&ClientMessage{
//		UserId:  guesser,
//		MsgType: TryGuess,
//		MsgData: map[string]any{"guess": "cat"},
//	}, vocabs)
//
//	if room_worker.Players[explainer].WordsGuessed != 1 {
//		t.Fatalf("expected explainer guessed score to be 1, got %d", room_worker.Players[explainer].WordsGuessed)
//	}
//	if room_worker.CurrentWordIndex != 1 {
//		t.Fatalf("expected current word index to advance to 1, got %d", room_worker.CurrentWordIndex)
//	}
//	if room_worker.wordShown {
//		t.Fatal("expected wordShown to reset to false after successful guess")
//	}
//}
//
//func TestHandleMessageGetWordCountsTriedOncePerWord(t *testing.T) {
//	room_worker, vocabs, explainer, _ := newRoomForMessageTests()
//
//	room_worker.handleMessage(&ClientMessage{
//		UserId:  explainer,
//		MsgType: GetWord,
//		MsgData: map[string]any{},
//	}, vocabs)
//	room_worker.handleMessage(&ClientMessage{
//		UserId:  explainer,
//		MsgType: GetWord,
//		MsgData: map[string]any{},
//	}, vocabs)
//
//	if room_worker.Players[explainer].WordsTried != 1 {
//		t.Fatalf("expected words tried to increase once for same word, got %d", room_worker.Players[explainer].WordsTried)
//	}
//}
//
//func TestHandleMessageStartRoundOnlyCurrentPlayer(t *testing.T) {
//	room_worker, vocabs, explainer, guesser := newRoomForMessageTests()
//	room_worker.State = RoundOver
//
//	room_worker.handleMessage(&ClientMessage{
//		UserId:  guesser,
//		MsgType: StartRound,
//		MsgData: map[string]any{},
//	}, vocabs)
//	if room_worker.State != RoundOver {
//		t.Fatalf("expected non-current player to be ignored, got state %d", room_worker.State)
//	}
//
//	room_worker.handleMessage(&ClientMessage{
//		UserId:  explainer,
//		MsgType: StartRound,
//		MsgData: map[string]any{},
//	}, vocabs)
//	if room_worker.State != Explaining {
//		t.Fatalf("expected current player to start round, got state %d", room_worker.State)
//	}
//	if room_worker.ticker == nil {
//		t.Fatal("expected ticker to be initialized after round start")
//	}
//	room_worker.ticker.Stop()
//}
