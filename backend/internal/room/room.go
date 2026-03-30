package room

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
	"github.com/xolra0d/alias-online/internal/config"
)

type Vocabulary struct {
	PrimaryWords []string
	RudeWords    []string
}

type Vocabularies struct {
	vocabulary map[string]*Vocabulary
	lock       sync.RWMutex
}

func NewVocabularies(vocabs map[string]*Vocabulary) *Vocabularies {
	return &Vocabularies{
		vocabulary: vocabs,
	}
}

func (v *Vocabularies) Contains(name string) bool {
	v.lock.RLock()
	_, ok := v.vocabulary[name]
	v.lock.RUnlock()
	return ok
}

func (v *Vocabularies) WordsInVocab(name string, rudeWords bool) (int, bool) {
	v.lock.RLock()
	defer v.lock.RUnlock()
	vocab, ok := v.vocabulary[name]
	if !ok {
		return 0, false
	}
	result := len(vocab.PrimaryWords)
	if rudeWords {
		result += len(vocab.RudeWords)
	}
	return result, true
}

func (v *Vocabularies) Languages() []string {
	v.lock.RLock()
	languages := make([]string, 0, len(v.vocabulary))
	for name := range v.vocabulary {
		languages = append(languages, name)
	}
	v.lock.RUnlock()
	return languages
}

type ClientMessageType int

const (
	GetState ClientMessageType = iota
	StartRound
	GetWord
	TryGuess
	FinishGame
	GetNewWord
)

type ServerMessageType int

const (
	NewUpdate ServerMessageType = iota
	CurrentState
	YourWord
	WordGuessed
	RightGuess
	WrongGuess
)

type GameState int

const (
	RoundOver GameState = iota
	Explaining
	Finished
)

type ClientMessage struct {
	UserId  uuid.UUID         `json:"user_id"`
	MsgType ClientMessageType `json:"type"`
	MsgData map[string]any    `json:"data"`
}

type ServerMessage struct {
	MsgType ServerMessageType `json:"msg_type"`
	MsgData map[string]any    `json:"msg_data"`
}

type Player struct {
	Id     uuid.UUID   `json:"id"`
	Name   string      `json:"name"`
	ToSend chan []byte `json:"-"`

	Ready        bool `json:"ready"`
	WordsTried   int  `json:"words_tried"`
	WordsGuessed int  `json:"words_guessed"`
}

func NewPlayer(id uuid.UUID, name string, wordsTried, wordsGuessed int) *Player {
	return &Player{
		id,
		name,
		make(chan []byte, 10),
		false,
		wordsTried,
		wordsGuessed,
	}
}

// RoomConfig holds specific room configuration
type RoomConfig struct {
	Seed                 int      `form:"-" json:"-"`
	WordsPerm            []int    `form:"-" json:"-"` // words permutation, unique for every room, dependent on Seed
	Language             string   `form:"language" json:"language"`
	RudeWords            bool     `form:"rude-words" json:"rude-words"`
	AdditionalVocabulary []string `form:"additional-vocabulary" json:"additional-vocabulary"`
	Clock                int      `form:"clock" json:"clock"`
}

type Room struct {
	Id     string      `json:"id"`
	Admin  uuid.UUID   `json:"admin"`
	Config *RoomConfig `json:"config"`

	Players    map[uuid.UUID]*Player `json:"players"`
	ingest     chan *ClientMessage
	readyCount uint
	join       chan *Player
	leave      chan uuid.UUID

	CurrentWordIndex int
	wordShown        bool
	TurnOrder        []uuid.UUID // circular queue
	currentPlayer    int         // points into TurnOrder. When RoundOver, points to next player
	State            GameState   `json:"game_state"`
	ticker           *time.Ticker
	RemainingTime    int `json:"remaining_time"`

	logger *config.Logger
}

func NewRoom(
	id string,
	admin uuid.UUID,
	cfg *RoomConfig,
	players map[uuid.UUID]*Player,
	turnOrder []uuid.UUID,
	currentPlayer int,
	currentWordIndex int,
	gameState GameState,
	logger *config.Logger,
) *Room {
	return &Room{
		Id:     id,
		Admin:  admin,
		Config: cfg,

		Players:    players,
		ingest:     make(chan *ClientMessage, 50),
		readyCount: 0,
		join:       make(chan *Player, 5),
		leave:      make(chan uuid.UUID, 5),

		TurnOrder:        turnOrder,
		currentPlayer:    currentPlayer,
		CurrentWordIndex: currentWordIndex,
		State:            gameState,

		// only state loaded from db is RoundOver, so no need for ticker
		ticker:        &time.Ticker{},
		RemainingTime: cfg.Clock,
		logger:        logger,
	}
}

func (r *Room) Join(player *Player) {
	r.join <- player
}

func (r *Room) Leave(player uuid.UUID) {
	r.leave <- player
}

func (r *Room) Ingest(msg *ClientMessage) {
	r.ingest <- msg
}

func (r *Room) IncCurrentPlayer() {
	if r.currentPlayer == len(r.TurnOrder)-1 {
		r.currentPlayer = 0
	} else {
		r.currentPlayer++
	}
}

func (r *Room) CurrentPlayer() uuid.UUID {
	return r.TurnOrder[r.currentPlayer]
}

func (r *Room) ToMap() map[string]any {
	return map[string]any{
		"admin":          r.Admin,
		"config":         r.Config,
		"players":        r.Players,
		"game_state":     r.State,
		"remaining_time": r.RemainingTime,
		"current_player": r.TurnOrder[r.currentPlayer],
	}
}

func (r *Room) Run(vocabs *Vocabularies, onEmpty func(room *Room) error) {
	op := "room.Run." + r.Id

	for {
		select {
		case <-r.ticker.C:
			if r.RemainingTime != 0 {
				r.RemainingTime--
				continue
			}

			r.logger.Info(op, "round ended", "currentPlayer", r.TurnOrder[r.currentPlayer])

			r.ticker.Stop()
			r.State = RoundOver
			r.RemainingTime = r.Config.Clock

			for range r.TurnOrder {
				r.IncCurrentPlayer()
				if r.Players[r.TurnOrder[r.currentPlayer]].Ready {
					break
				}
			}
			r.CurrentWordIndex += 1
			r.wordShown = false
			r.ReportUpdate()
		case msg := <-r.ingest:
			r.handleMessage(msg, vocabs)
		case player := <-r.join:
			r.handleJoin(player)

			for range r.TurnOrder {
				if r.Players[r.TurnOrder[r.currentPlayer]].Ready {
					break
				}
				r.IncCurrentPlayer()
			}
			r.logger.Info(op, "player joined", "playerId", player.Id)
		case id := <-r.leave:
			r.handleLeave(id)
			r.logger.Info(op, "player left", "playerId", id)
			if r.readyCount == 0 {
				err := onEmpty(r)
				if err != nil {
					r.logger.Error(op, "failed to save state", "error", err)
				}

				r.logger.Info(op, "room is removed")
				return
			}
		}
	}
}

func (r *Room) handleMessage(msg *ClientMessage, vocabs *Vocabularies) {
	op := "room.handleMessage." + r.Id

	switch r.State {
	case RoundOver:
		switch msg.MsgType {
		case GetState:
			r.sendState(msg.UserId)
		case StartRound:
			if msg.UserId != r.TurnOrder[r.currentPlayer] {
				return
			}
			r.State = Explaining
			r.ticker = time.NewTicker(time.Second) // update Room.RemainingTime every second.
			r.wordShown = false
			r.ReportUpdate()
			r.logger.Info(op, "round started", "currentPlayer", r.TurnOrder[r.currentPlayer])
		case FinishGame:
			if msg.UserId != r.Admin {
				return
			}
			r.State = Finished
			r.ReportUpdate()
			r.logger.Info(op, "game is finished")
		default:
			r.logger.Warn(op, "unknown msg type", "msg", msg)
		}
	case Explaining:
		switch msg.MsgType {
		case GetState:
			r.sendState(msg.UserId)
		case GetWord:
			if msg.UserId != r.TurnOrder[r.currentPlayer] {
				return
			}
			word := r.CurrentWord(vocabs)
			m := ServerMessage{
				MsgType: YourWord,
				MsgData: map[string]any{
					"word": word,
				},
			}
			b, err := json.Marshal(&m)
			if err != nil {
				r.logger.Error(op, "unexpected panic, could not marshal", "msg", m, "error", err)
				return
			}
			r.Players[msg.UserId].ToSend <- b

			if !r.wordShown {
				r.wordShown = true

				r.logger.Info(op, "new word", "word", word, "currentPlayer", msg.UserId)

				r.Players[msg.UserId].WordsTried++
				r.ReportUpdate()
			}
		case GetNewWord:
			if msg.UserId != r.TurnOrder[r.currentPlayer] {
				return
			}
			r.CurrentWordIndex++
			r.wordShown = true
			r.Players[msg.UserId].WordsTried++

			word := r.CurrentWord(vocabs)
			m := ServerMessage{
				MsgType: YourWord,
				MsgData: map[string]any{
					"word": word,
				},
			}

			r.logger.Info(op, "new word", "word", word, "currentPlayer", msg.UserId)

			b, err := json.Marshal(&m)
			if err != nil {
				r.logger.Error(op, "unexpected panic, could not marshal", "msg", m, "error", err)
				return
			}
			r.Players[msg.UserId].ToSend <- b
			r.ReportUpdate()
		case TryGuess:
			if msg.UserId == r.TurnOrder[r.currentPlayer] {
				return
			}
			g, ok := msg.MsgData["guess"]
			if !ok {
				return
			}
			guess, ok := g.(string)
			if !ok {
				return
			}

			word := r.CurrentWord(vocabs)
			if guess == word {
				id := r.TurnOrder[r.currentPlayer]
				r.Players[id].WordsGuessed++

				r.logger.Info(op, "Word guessed", "word", word)

				r.CurrentWordIndex++
				r.wordShown = false

				r.WordGuessed(msg.UserId)
				r.ReportUpdate()
			} else {
				m := ServerMessage{
					MsgType: WrongGuess,
					MsgData: map[string]any{},
				}
				b, err := json.Marshal(&m)
				if err != nil {
					r.logger.Error(op, "unexpected panic, could not marshal", "msg", m, "error", err)
					return
				}
				r.Players[msg.UserId].ToSend <- b
			}
		default:
			r.logger.Warn(op, "unknown msg type", "msg", msg)
		}
	case Finished:
		if msg.MsgType == GetState {
			r.sendState(msg.UserId)
		}
	}
}

func (r *Room) WordGuessed(guesser uuid.UUID) {
	op := "room.WordGuessed." + r.Id

	for _, player := range r.Players {
		if !player.Ready {
			continue
		}

		var msgType ServerMessageType
		if player.Id == guesser {
			msgType = RightGuess
		} else {
			msgType = WordGuessed
		}

		m := ServerMessage{
			MsgType: msgType,
			MsgData: map[string]any{
				"guesser": guesser,
			},
		}
		b, err := json.Marshal(&m)
		if err != nil {
			r.logger.Error(op, "unexpected panic, could not marshal", "msg", m, "error", err)
			return
		}
		player.ToSend <- b
	}
}

func (r *Room) getWordAt(vocabs *Vocabularies, index int) string {
	if index < 0 || index >= len(r.Config.WordsPerm) {
		return ""
	}

	i := r.Config.WordsPerm[index]
	vocabs.lock.RLock()
	v, ok := vocabs.vocabulary[r.Config.Language]
	defer vocabs.lock.RUnlock()
	if !ok {
		return ""
	}

	primaryLen := len(v.PrimaryWords)
	if i < primaryLen {
		return v.PrimaryWords[i]
	}

	offset := primaryLen
	if r.Config.RudeWords {
		rudeLen := len(v.RudeWords)
		if i < offset+rudeLen {
			return v.RudeWords[i-offset]
		}
		offset += rudeLen
	}

	if i < offset+len(r.Config.AdditionalVocabulary) {
		return r.Config.AdditionalVocabulary[i-offset]
	}

	return ""
}

func (r *Room) CurrentWord(vocabs *Vocabularies) string {
	return r.getWordAt(vocabs, r.CurrentWordIndex)
}

func (r *Room) NextWord(vocabs *Vocabularies) string {
	return r.getWordAt(vocabs, r.CurrentWordIndex+1)
}

func (r *Room) sendState(id uuid.UUID) {
	op := "room.sendState." + r.Id

	m := ServerMessage{
		CurrentState,
		r.ToMap(),
	}
	b, err := json.Marshal(&m)
	if err != nil {
		r.logger.Error(op, "unexpected panic, could not marshal", "msg", m, "error", err)
		return
	}
	r.Players[id].ToSend <- b
}

func (r *Room) ReportUpdate() {
	op := "room.ReportUpdate." + r.Id

	m := ServerMessage{
		MsgType: NewUpdate,
		MsgData: make(map[string]any),
	}
	b, err := json.Marshal(&m)
	if err != nil {
		r.logger.Error(op, "unexpected panic, could not marshal", "msg", m, "error", err)
		return
	}

	for _, player := range r.Players {
		if player.Ready {
			player.ToSend <- b
		}
	}
}

func (r *Room) handleJoin(player *Player) {
	p, ok := r.Players[player.Id]
	if ok {
		p.Ready = true
		p.ToSend = player.ToSend
	} else {
		p = player
		r.Players[p.Id] = p
		p.Ready = true
		r.TurnOrder = append(r.TurnOrder, p.Id)
	}
	r.readyCount++
	r.ReportUpdate()
}

func (r *Room) handleLeave(id uuid.UUID) {
	r.Players[id].Ready = false
	r.readyCount--
	r.ReportUpdate()
}

type PlayerSnapshot struct {
	Id           uuid.UUID
	Name         string
	WordsTried   int
	WordsGuessed int
}

func toClientMessage(v map[string]any) (*ClientMessage, error) {
	u, ok := v["user_id"]
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}
	us, ok := u.(string)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}
	user, err := uuid.Parse(us)
	if err != nil {
		return nil, err
	}

	t, ok := v["type"]
	if !ok {
		return nil, fmt.Errorf("invalid msg type")
	}
	ty, ok := t.(float64)
	if !ok {
		return nil, fmt.Errorf("invalid msg type")
	}
	typ := ClientMessageType(int(ty))

	d, ok := v["data"]
	if !ok {
		return nil, fmt.Errorf("invalid msg data")
	}
	da, ok := d.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid msg data")
	}
	return &ClientMessage{user, typ, da}, nil
}

func (r *Room) RunWriter(ctx context.Context, cancel context.CancelFunc, c *websocket.Conn, player *Player, wsWriteTimeout, pingTimeout time.Duration) {
	op := "room.RunWriter" + r.Id

	ping := time.NewTicker(time.Second * 10)
	for {
		select {
		case msg := <-player.ToSend:
			writeCtx, writeCancel := context.WithTimeout(ctx, wsWriteTimeout)
			err := c.Write(writeCtx, websocket.MessageBinary, msg)
			writeCancel()
			if err != nil {
				r.logger.Error(op, "write error", "playerId", player.Id, "msg", msg, "err", err)
				cancel()
				return // this goroutine does not log errors.
			}
		case <-ping.C:
			pingCtx, pingCancel := context.WithTimeout(ctx, pingTimeout)
			err := c.Ping(pingCtx)
			pingCancel()
			if err != nil {
				r.logger.Error(op, "ping error", "playerId", player.Id, "err", err)
				cancel()
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (r *Room) RunReader(ctx context.Context, cancel context.CancelFunc, c *websocket.Conn, player *Player, maxMessagesPerSecond int) error {
	op := "room.RunReader." + r.Id

	windowStartedAt := time.Now()
	messagesInWindow := 0
	for {
		var v any
		err := wsjson.Read(ctx, c, &v)
		if err != nil {
			cancel()
			closeStatus := websocket.CloseStatus(err)
			if closeStatus == websocket.StatusGoingAway || closeStatus == websocket.StatusNormalClosure {
				r.logger.Info(op, "ws client disconnected normally", "roomId", r.Id)
				return nil
			}
			r.logger.Error(op, "ws read error", "roomId", r.Id, "error", err)
			return err
		}

		if maxMessagesPerSecond > 0 {
			now := time.Now()
			if now.Sub(windowStartedAt) >= time.Second {
				windowStartedAt = now
				messagesInWindow = 0
			}
			messagesInWindow++
			if messagesInWindow > maxMessagesPerSecond {
				_ = c.Close(websocket.StatusPolicyViolation, "too many messages")
				r.logger.Error(op, "websocket message rate limit exceeded", "roomId", r.Id, "playerId", player.Id)
				cancel()
				return fmt.Errorf("websocket message rate limit exceeded")
			}
		}

		m, ok := v.(map[string]any)
		if !ok {
			err := fmt.Errorf("invalid client msg: %v", v)
			cancel()
			return err
		}
		msg, err := toClientMessage(m)
		if err != nil {
			cancel()
			return err
		}
		r.ingest <- msg
	}
}
