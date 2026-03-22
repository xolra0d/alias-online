package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
)

type Vocabulary struct {
	PrimaryWords []string
	RudeWords    []string
}

type Vocabularies struct {
	vocabulary map[string]*Vocabulary
	lock       sync.RWMutex
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
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	toSend chan []byte

	Ready        bool `json:"ready"`
	WordsTried   int  `json:"words_tried"`
	WordsGuessed int  `json:"words_guessed"`
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

	currentWord   int
	wordShown     bool
	turnOrder     []uuid.UUID // circular queue
	currentPlayer int         // points into turnOrder. When RoundOver, points to next player
	State         GameState   `json:"game_state"`
	ticker        *time.Ticker
	RemainingTime int `json:"remaining_time"`

	logger *PrefixLogger
}

type Rooms struct {
	rooms  map[string]*Room
	logger *PrefixLogger
	lock   sync.RWMutex

	MinClock                     int
	MaxClock                     int
	MaxAdditionalVocabularyWords int
	MaxAdditionalWordLength      int

	WSOriginPatterns     []string
	MaxMessagesPerSecond int
	PingTimeout          time.Duration
	WSWriteTimeout       time.Duration
	LoadRoomTimeout      time.Duration
	SaveRoomTimeout      time.Duration
}

func (r *Room) IncCurrentPlayer() {
	if r.currentPlayer == len(r.turnOrder)-1 {
		r.currentPlayer = 0
	} else {
		r.currentPlayer++
	}
}

func (r *Room) ToMap() map[string]any {
	return map[string]any{
		"admin":          r.Admin,
		"config":         r.Config,
		"players":        r.Players,
		"game_state":     r.State,
		"remaining_time": r.RemainingTime,
		"current_player": r.turnOrder[r.currentPlayer],
	}
}

func (r *Room) Run(postgres *Postgres, vocabs *Vocabularies, rooms *Rooms) {
	for {
		select {
		case <-r.ticker.C:
			if r.RemainingTime != 0 {
				r.RemainingTime--
				continue
			}

			r.logger.Info("round ended", "currentPlayer", r.turnOrder[r.currentPlayer])

			r.ticker.Stop()
			r.State = RoundOver
			r.RemainingTime = r.Config.Clock

			for range r.turnOrder {
				r.IncCurrentPlayer()
				if r.Players[r.turnOrder[r.currentPlayer]].Ready {
					break
				}
			}
			r.currentWord += 1
			r.wordShown = false
			r.ReportUpdate()
		case msg := <-r.ingest:
			r.handleMessage(msg, vocabs)
		case player := <-r.join:
			r.handleJoin(player)

			for range r.turnOrder {
				if r.Players[r.turnOrder[r.currentPlayer]].Ready {
					break
				}
				r.IncCurrentPlayer()
			}
			r.logger.Info("player joined", "playerId", player.Id)
		case id := <-r.leave:
			r.handleLeave(id)
			r.logger.Info("player left", "playerId", id)
			if r.readyCount == 0 {
				ctx, cancel := context.WithTimeout(context.Background(), rooms.SaveRoomTimeout)
				err := r.SaveState(ctx, postgres)
				cancel()
				if err != nil {
					r.logger.Error("failed to save state", "error", err)
				}
				rooms.lock.Lock()
				delete(rooms.rooms, r.Id)
				rooms.lock.Unlock()

				r.logger.Info("room is removed")
				return
			}
		}
	}
}

func (r *Room) handleMessage(msg *ClientMessage, vocabs *Vocabularies) {
	switch r.State {
	case RoundOver:
		switch msg.MsgType {
		case GetState:
			r.sendState(msg.UserId)
		case StartRound:
			if msg.UserId != r.turnOrder[r.currentPlayer] {
				return
			}
			r.State = Explaining
			r.ticker = time.NewTicker(time.Second) // update Room.RemainingTime every second.
			r.wordShown = false
			r.ReportUpdate()
			r.logger.Info("round started", "currentPlayer", r.turnOrder[r.currentPlayer])
		case FinishGame:
			if msg.UserId != r.Admin {
				return
			}
			r.State = Finished
			r.ReportUpdate()
			r.logger.Info("game is finished")
		default:
			r.logger.Warn("unknown msg type", "msg", msg)
		}
	case Explaining:
		switch msg.MsgType {
		case GetState:
			r.sendState(msg.UserId)
		case GetWord:
			if msg.UserId != r.turnOrder[r.currentPlayer] {
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
				r.logger.Error("unexpected panic, could not marshal", "msg", m, "error", err)
				return
			}
			r.Players[msg.UserId].toSend <- b

			if !r.wordShown {
				r.wordShown = true
				r.Players[msg.UserId].WordsTried++
				r.ReportUpdate()
			}
		case GetNewWord:
			if msg.UserId != r.turnOrder[r.currentPlayer] {
				return
			}
			r.currentWord++
			r.wordShown = true
			r.Players[msg.UserId].WordsTried++

			word := r.CurrentWord(vocabs)
			m := ServerMessage{
				MsgType: YourWord,
				MsgData: map[string]any{
					"word": word,
				},
			}
			b, err := json.Marshal(&m)
			if err != nil {
				r.logger.Error("unexpected panic, could not marshal", "msg", m, "error", err)
				return
			}
			r.Players[msg.UserId].toSend <- b
			r.ReportUpdate()
		case TryGuess:
			if msg.UserId == r.turnOrder[r.currentPlayer] {
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
				id := r.turnOrder[r.currentPlayer]
				r.Players[id].WordsGuessed++

				r.logger.Info("Word guessed", "word", word)

				r.currentWord++
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
					r.logger.Error("unexpected panic, could not marshal", "msg", m, "error", err)
					return
				}
				r.Players[msg.UserId].toSend <- b
			}
		default:
			r.logger.Warn("unknown msg type", "msg", msg)
		}
	case Finished:
		if msg.MsgType == GetState {
			r.sendState(msg.UserId)
		}
	}
}

func (r *Room) WordGuessed(guesser uuid.UUID) {
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
			r.logger.Error("unexpected panic, could not marshal", "msg", m, "error", err)
			return
		}
		player.toSend <- b
	}
}

func (r *Room) getWordAt(vocabs *Vocabularies, index int) string {
	if index < 0 || index >= len(r.Config.wordsPerm) {
		return ""
	}

	i := r.Config.wordsPerm[index]
	vocabs.lock.RLock()
	v := vocabs.vocabulary[r.Config.Language]
	vocabs.lock.RUnlock()

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
	return r.getWordAt(vocabs, r.currentWord)
}

func (r *Room) NextWord(vocabs *Vocabularies) string {
	return r.getWordAt(vocabs, r.currentWord+1)
}

func (r *Room) sendState(id uuid.UUID) {
	m := ServerMessage{
		CurrentState,
		r.ToMap(),
	}
	b, err := json.Marshal(&m)
	if err != nil {
		r.logger.Error("unexpected panic, could not marshal", "msg", m, "error", err)
		return
	}
	r.Players[id].toSend <- b
}

func (r *Room) ReportUpdate() {
	m := ServerMessage{
		MsgType: NewUpdate,
		MsgData: make(map[string]any),
	}
	b, err := json.Marshal(&m)
	if err != nil {
		r.logger.Error("unexpected panic, could not marshal", "msg", m, "error", err)
		return
	}

	for _, player := range r.Players {
		if player.Ready {
			player.toSend <- b
		}
	}
}

func (r *Room) handleJoin(player *Player) {
	p, ok := r.Players[player.Id]
	if ok {
		p.Ready = true
		p.toSend = player.toSend
	} else {
		p = player
		r.Players[p.Id] = p
		p.Ready = true
		r.turnOrder = append(r.turnOrder, p.Id)
	}
	r.readyCount++
	r.ReportUpdate()
}

func (r *Room) SaveState(ctx context.Context, postgres *Postgres) error {
	err := postgres.UpdateRoomState(ctx, r)
	if err != nil {
		r.logger.Error("could not save room state", "error", err)
	}
	return err
}

func (r *Room) handleLeave(id uuid.UUID) {
	r.Players[id].Ready = false
	r.readyCount--
	r.ReportUpdate()
}

func (rooms *Rooms) ServeWS(w http.ResponseWriter, r *http.Request, userId uuid.UUID, name, roomId string, postgres *Postgres, vocabs *Vocabularies) error {
	rooms.lock.RLock()
	room, ok := rooms.rooms[roomId]
	rooms.lock.RUnlock()

	var loadedRoom *Room
	if !ok {
		ctx, cancel := context.WithTimeout(context.Background(), rooms.LoadRoomTimeout)
		newRoom, err := postgres.LoadRoom(ctx, roomId, vocabs)
		cancel()
		if err != nil {
			return err
		}
		loadedRoom = newRoom
	}

	acceptOptions := &websocket.AcceptOptions{}
	if len(rooms.WSOriginPatterns) > 0 {
		acceptOptions.OriginPatterns = rooms.WSOriginPatterns
	}
	c, err := websocket.Accept(w, r, acceptOptions)
	if err != nil {
		return err
	}
	defer func(c *websocket.Conn) {
		err := c.CloseNow()
		if err != nil {
			rooms.logger.Error("error closing websocket connection", "error", err)
		}
	}(c)

	if !ok {
		rooms.lock.Lock()
		existingRoom, exists := rooms.rooms[roomId]
		if exists {
			room = existingRoom
		} else {
			room = loadedRoom
			rooms.rooms[roomId] = room
			go room.Run(postgres, vocabs, rooms)
		}
		rooms.lock.Unlock()
	}

	toSend := make(chan []byte, 20)
	player := &Player{Id: userId, Name: name, toSend: toSend}

	room.join <- player
	defer func() { room.leave <- userId }()

	ctx, cancel := context.WithCancel(r.Context())

	go func() {
		ping := time.NewTicker(time.Second * 10)
		for {
			select {
			case msg := <-toSend:
				writeCtx, writeCancel := context.WithTimeout(ctx, rooms.WSWriteTimeout)
				err := c.Write(writeCtx, websocket.MessageBinary, msg)
				writeCancel()
				if err != nil {
					room.logger.Error("write error", "playerId", player.Id, "msg", msg, "err", err)
					cancel()
					return // this goroutine does not log errors.
				}
			case <-ping.C:
				pingCtx, pingCancel := context.WithTimeout(ctx, rooms.PingTimeout)
				err := c.Ping(pingCtx)
				pingCancel()
				if err != nil {
					room.logger.Error("ping error", "playerId", player.Id, "err", err)
					cancel()
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	windowStartedAt := time.Now()
	messagesInWindow := 0
	for {
		var v any
		err := wsjson.Read(ctx, c, &v)
		if err != nil {
			cancel()
			rooms.logger.Error("ws read error", "roomId", roomId, "error", err)
			return err
		}

		if rooms.MaxMessagesPerSecond > 0 {
			now := time.Now()
			if now.Sub(windowStartedAt) >= time.Second {
				windowStartedAt = now
				messagesInWindow = 0
			}
			messagesInWindow++
			if messagesInWindow > rooms.MaxMessagesPerSecond {
				_ = c.Close(websocket.StatusPolicyViolation, "too many messages")
				rooms.logger.Error("websocket message rate limit exceeded", "roomId", roomId, "playerId", userId)
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
		room.ingest <- msg
	}
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
