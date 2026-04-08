package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type Vocabulary struct {
	PrimaryWords []string
	RudeWords    []string
}

type ClientMessageType int

const (
	GetState ClientMessageType = iota
	StartRound
	GetWord
	TryGuess
	FinishGame
	GetNewWord
	CreateRoom
	LoadRoom
)

type ServerMessageType int

const (
	NewUpdate ServerMessageType = iota
	CurrentState
	YourWord
	WordGuessed
	RightGuess
	WrongGuess
	Redirect // if other worker reserved room
)

type GameState int

const (
	RoundOver GameState = iota
	Explaining
	Finished
)

type ClientMessage struct {
	UserId  string            `json:"user_id"`
	MsgType ClientMessageType `json:"type"`
	MsgData map[string]any    `json:"data"`
}

type ServerMessage struct {
	MsgType ServerMessageType `json:"msg_type"`
	MsgData map[string]any    `json:"msg_data"`
}

type Player struct {
	Id     string      `json:"id"`
	Name   string      `json:"name"`
	ToSend chan []byte `json:"-"`

	Ready        bool `json:"ready"`
	WordsTried   int  `json:"words_tried"`
	WordsGuessed int  `json:"words_guessed"`
}

func NewPlayer(id, name string, wordsTried, wordsGuessed int) *Player {
	return &Player{
		id,
		name,
		make(chan []byte, 10),
		false,
		wordsTried,
		wordsGuessed,
	}
}

// RoomConfig holds specific room_worker configuration
type RoomConfig struct {
	Seed                 int      `form:"-" json:"-"`
	AllWords             []string `form:"-" json:"-"` // words permutation, unique for every room, dependent on Seed
	Language             string   `form:"language" json:"language"`
	RudeWords            bool     `form:"rude-words" json:"rude-words"`
	AdditionalVocabulary []string `form:"additional-vocabulary" json:"additional-vocabulary"`
	Clock                int      `form:"clock" json:"clock"`
}

type Room struct {
	Id     string      `json:"id"`
	Admin  string      `json:"admin"`
	Config *RoomConfig `json:"config"`

	Players    map[string]*Player `json:"players"`
	ingest     chan *ClientMessage
	readyCount uint
	join       chan *Player
	leave      chan string

	CurrentWordIndex int
	wordShown        bool
	TurnOrder        []string  // circular queue
	currentPlayer    int       // points into TurnOrder. When RoundOver, points to next player
	State            GameState `json:"game_state"`
	ticker           *time.Ticker
	RemainingTime    int `json:"remaining_time"`

	logger       *slog.Logger
	prepareState *PrepareState
}

func (r *Room) WaitUntilOperational() error {
	return r.prepareState.WaitUntilOperational()
}

func (r *Room) SetOperational() {
	r.prepareState.SetOperational()
}

func (r *Room) SetErrored() {
	r.prepareState.SetErrored()
}

func (r *Room) UpdateStateFromRoom(newRoom *Room) {
	r.Id = newRoom.Id
	r.Admin = newRoom.Admin
	r.Config = newRoom.Config
	r.Players = newRoom.Players
	r.CurrentWordIndex = newRoom.CurrentWordIndex
	r.TurnOrder = newRoom.TurnOrder
	r.currentPlayer = newRoom.currentPlayer
	r.State = newRoom.State
	r.logger = newRoom.logger
}

func (r *Room) UpdateStateFromRoomConfig(roomId, name, admin string, cfg *RoomConfig, logger *slog.Logger) {
	r.Id = roomId
	r.Admin = admin
	r.Config = cfg
	r.Players = map[string]*Player{admin: {Id: admin, Name: name}}
	r.TurnOrder = []string{admin}
	r.RemainingTime = cfg.Clock
	r.logger = logger
}

func NewPreparingRoom() *Room {
	return &Room{
		ingest:        make(chan *ClientMessage, 50),
		readyCount:    0,
		join:          make(chan *Player, 5),
		leave:         make(chan string, 5),
		prepareState:  NewPrepareState(),
		ticker:        &time.Ticker{},
		RemainingTime: 100,
	}
}

func NewRoom(
	id string,
	admin string,
	cfg *RoomConfig,
	players map[string]*Player,
	turnOrder []string,
	currentPlayer int,
	currentWordIndex int,
	gameState GameState,
	logger *slog.Logger,
) *Room {
	return &Room{
		Id:     id,
		Admin:  admin,
		Config: cfg,

		Players:    players,
		ingest:     make(chan *ClientMessage, 50),
		readyCount: 0,
		join:       make(chan *Player, 5),
		leave:      make(chan string, 5),

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

func (r *Room) Leave(player string) {
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

func (r *Room) CurrentPlayer() string {
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

func (r *Room) Run(onEmpty func(room *Room)) {
	for {
		select {
		case <-r.ticker.C:
			if r.RemainingTime != 0 {
				r.RemainingTime--
				continue
			}

			r.logger.Info("round ended", "roomId", r.Id, "currentPlayer", r.TurnOrder[r.currentPlayer])

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
			r.handleMessage(msg)
		case player := <-r.join:
			r.handleJoin(player)

			for range r.TurnOrder {
				if r.Players[r.TurnOrder[r.currentPlayer]].Ready {
					break
				}
				r.IncCurrentPlayer()
			}
			r.logger.Info("player joined", "roomId", r.Id, "playerId", player.Id)
		case id := <-r.leave:
			r.handleLeave(id)
			r.logger.Info("player left", "roomId", r.Id, "playerId", id)
			if r.readyCount == 0 {
				onEmpty(r)
				r.logger.Info("room is removed", "roomId", r.Id)
				return
			}
		}
	}
}

func (r *Room) handleMessage(msg *ClientMessage) {
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
			r.logger.Info("round started", "roomId", r.Id, "currentPlayer", r.TurnOrder[r.currentPlayer])
		case FinishGame:
			if msg.UserId != r.Admin {
				return
			}
			r.State = Finished
			r.ReportUpdate()
			r.logger.Info("game is finished", "roomId", r.Id)
		default:
			r.logger.Warn("unknown msg type", "roomId", r.Id, "msgType", msg.MsgType, "userId", msg.UserId)
		}
	case Explaining:
		switch msg.MsgType {
		case GetState:
			r.sendState(msg.UserId)
		case GetWord:
			if msg.UserId != r.TurnOrder[r.currentPlayer] {
				return
			}
			word := r.CurrentWord()
			m := ServerMessage{
				MsgType: YourWord,
				MsgData: map[string]any{
					"word": word,
				},
			}
			b, err := json.Marshal(&m)
			if err != nil {
				r.logger.Error("could not marshal YourWord", "roomId", r.Id, "userId", msg.UserId, "error", err)
				return
			}
			r.Players[msg.UserId].ToSend <- b

			if !r.wordShown {
				r.wordShown = true
				r.logger.Info("new word shown", "roomId", r.Id, "word", word, "currentPlayer", msg.UserId)
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

			word := r.CurrentWord()
			m := ServerMessage{
				MsgType: YourWord,
				MsgData: map[string]any{
					"word": word,
				},
			}

			r.logger.Info("new word", "roomId", r.Id, "word", word, "currentPlayer", msg.UserId)

			b, err := json.Marshal(&m)
			if err != nil {
				r.logger.Error("could not marshal YourWord", "roomId", r.Id, "userId", msg.UserId, "error", err)
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

			word := r.CurrentWord()
			if guess == word {
				id := r.TurnOrder[r.currentPlayer]
				r.Players[id].WordsGuessed++

				r.logger.Info("word guessed", "roomId", r.Id, "word", word, "guesser", msg.UserId, "currentPlayer", id)

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
					r.logger.Error("could not marshal WrongGuess", "roomId", r.Id, "userId", msg.UserId, "error", err)
					return
				}
				r.Players[msg.UserId].ToSend <- b
			}
		default:
			r.logger.Warn("unknown msg type", "roomId", r.Id, "msgType", msg.MsgType, "userId", msg.UserId)
		}
	case Finished:
		if msg.MsgType == GetState {
			r.sendState(msg.UserId)
		}
	}
}

func (r *Room) WordGuessed(guesser string) {
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
			r.logger.Error("could not marshal WordGuessed", "roomId", r.Id, "guesser", guesser, "playerId", player.Id, "error", err)
			return
		}
		player.ToSend <- b
	}
}

func (r *Room) getWordAt(index int) string {
	if index < 0 || index >= len(r.Config.AllWords) {
		return ""
	}

	return r.Config.AllWords[index]
}

func (r *Room) CurrentWord() string {
	return r.getWordAt(r.CurrentWordIndex)
}

func (r *Room) NextWord() string {
	return r.getWordAt(r.CurrentWordIndex + 1)
}

func (r *Room) sendState(id string) {
	m := ServerMessage{
		CurrentState,
		r.ToMap(),
	}
	b, err := json.Marshal(&m)
	if err != nil {
		r.logger.Error("could not marshal CurrentState", "roomId", r.Id, "userId", id, "error", err)
		return
	}
	r.Players[id].ToSend <- b
}

func (r *Room) ReportUpdate() {
	m := ServerMessage{
		MsgType: NewUpdate,
		MsgData: make(map[string]any),
	}
	b, err := json.Marshal(&m)
	if err != nil {
		r.logger.Error("could not marshal NewUpdate", "roomId", r.Id, "error", err)
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

func (r *Room) handleLeave(id string) {
	r.Players[id].Ready = false
	r.readyCount--
	r.ReportUpdate()
}

func toClientMessage(v map[string]any) (*ClientMessage, error) {
	u, ok := v["user_id"]
	if !ok {
		return nil, fmt.Errorf("invalid user id")
	}
	user, ok := u.(string)
	if !ok {
		return nil, fmt.Errorf("invalid user id")
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
	ping := time.NewTicker(time.Second * 10)
	for {
		select {
		case msg := <-player.ToSend:
			writeCtx, writeCancel := context.WithTimeout(ctx, wsWriteTimeout)
			err := c.Write(writeCtx, websocket.MessageBinary, msg)
			writeCancel()
			if err != nil {
				r.logger.Error("write error", "roomId", r.Id, "playerId", player.Id, "error", err)
				cancel()
				return
			}
		case <-ping.C:
			pingCtx, pingCancel := context.WithTimeout(ctx, pingTimeout)
			err := c.Ping(pingCtx)
			pingCancel()
			if err != nil {
				r.logger.Error("ping error", "roomId", r.Id, "playerId", player.Id, "error", err)
				cancel()
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (r *Room) RunReader(ctx context.Context, cancel context.CancelFunc, c *websocket.Conn, player *Player, maxMessagesPerSecond int) {
	windowStartedAt := time.Now()
	messagesInWindow := 0
	for {
		var v any
		err := wsjson.Read(ctx, c, &v)
		if err != nil {
			cancel()
			closeStatus := websocket.CloseStatus(err)
			if closeStatus == websocket.StatusGoingAway || closeStatus == websocket.StatusNormalClosure {
				r.logger.Info("ws client disconnected normally", "roomId", r.Id, "playerId", player.Id)
				return
			}
			r.logger.Error("ws read error", "roomId", r.Id, "playerId", player.Id, "error", err)
			return
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
				r.logger.Error("websocket message rate limit exceeded", "roomId", r.Id, "playerId", player.Id)
				cancel()
				return
			}
		}

		m, ok := v.(map[string]any)
		if !ok {
			r.logger.Error("invalid client msg type assertion", "roomId", r.Id, "playerId", player.Id)
			cancel()
			return
		}
		msg, err := toClientMessage(m)
		if err != nil {
			r.logger.Error("invalid client msg", "roomId", r.Id, "playerId", player.Id, "error", err)
			cancel()
			return
		}
		r.ingest <- msg
	}
}
