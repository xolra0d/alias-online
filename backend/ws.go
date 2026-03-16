package main

import (
	"container/ring"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
)

var pingTimeout = sync.OnceValue(func() time.Duration {
	pingT, err := strconv.ParseUint(os.Getenv("PING_TIMEOUT"), 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Duration(pingT) * time.Second
})

type Vocabulary struct {
	PrimaryWords []string
	RudeWords    []string
}

type Vocabularies struct {
	vocabulary map[string]Vocabulary
	lock       sync.RWMutex
}

type ClientMessageType int

const (
	GetState ClientMessageType = iota
	StartRound
	GetWord
	TryGuess
	FinishGame
	SkipWord
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
	toSend chan []byte

	Ready        bool `json:"ready"`
	WordsTried   int  `json:"words_tried"`
	WordsGuessed int  `json:"words_guessed"`
}

type GameState int

const (
	RoundOver GameState = iota
	Explaining
	Finished
)

type Room struct {
	Id     string      `json:"id"`
	Admin  uuid.UUID   `json:"admin"`
	Config *RoomConfig `json:"config"`

	Players map[uuid.UUID]*Player `json:"players"`
	ingest  chan *ClientMessage
	count   uint
	join    chan *Player
	leave   chan uuid.UUID

	currentWord   int
	currentPlayer *ring.Ring // ring of UUID TODO: report current player
	State         GameState  `json:"game_state"`
	ticker        *time.Ticker
	RemainingTime int `json:"remaining_time"`
}

func (r *Room) ToMap() map[string]any {
	var currentPlayer uuid.UUID
	if r.currentPlayer != nil {
		id, ok := r.currentPlayer.Value.(uuid.UUID)
		if ok {
			currentPlayer = id
		}
	}
	return map[string]any{
		"admin":          r.Admin,
		"config":         r.Config,
		"players":        r.Players,
		"game_state":     r.State,
		"remaining_time": r.RemainingTime,
		"current_player": currentPlayer,
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

			r.ticker.Stop()
			r.State = RoundOver
			r.RemainingTime = r.Config.Clock

			for range r.Players {
				r.currentPlayer = r.currentPlayer.Next()
				id, ok := r.currentPlayer.Value.(uuid.UUID)
				if !ok {
					panic("invalid player type")
				}
				if r.Players[id].Ready {
					break
				}
			}
			r.currentWord += 1
			r.ReportUpdate()
		case msg := <-r.ingest:
			r.handleMessage(msg, vocabs)
		case player := <-r.join:
			r.handleJoin(player)

			for range r.Players {
				if !r.Players[r.currentPlayer.Value.(uuid.UUID)].Ready {
					r.currentPlayer = r.currentPlayer.Next()
				}
			}
		case id := <-r.leave:
			r.handleLeave(id, postgres, rooms)
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
			id, ok := r.currentPlayer.Value.(uuid.UUID)
			if !ok {
				panic("invalid player type in ring")
			}
			if msg.UserId != id {
				return
			}
			r.State = Explaining
			r.ticker = time.NewTicker(time.Second) // update Room.RemainingTime every second.
			r.ReportUpdate()
		case FinishGame:
			if msg.UserId != r.Admin {
				return
			}
			r.State = Finished

			r.ReportUpdate()
		}
	case Explaining:
		switch msg.MsgType {
		case GetState:
			r.sendState(msg.UserId)
		case GetWord:
			id, ok := r.currentPlayer.Value.(uuid.UUID)
			if !ok {
				panic("invalid player type in ring")
			}
			if msg.UserId != id {
				return
			}
			word := r.CurrentWord(vocabs)
			fmt.Println(word)
			b, err := json.Marshal(ServerMessage{
				MsgType: YourWord,
				MsgData: map[string]any{
					"word": word,
				},
			})
			if err != nil {
				panic(err)
			}
			r.Players[msg.UserId].toSend <- b
		case SkipWord:
			id, ok := r.currentPlayer.Value.(uuid.UUID)
			if !ok {
				panic("invalid player type in ring")
			}
			if msg.UserId != id {
				return
			}
			3
		case TryGuess:
			g, ok := msg.MsgData["guess"]
			if !ok {
				return
			}
			guess, ok := g.(string)
			if !ok {
				return
			}
			if guess == r.CurrentWord(vocabs) {
				id, ok := r.currentPlayer.Value.(uuid.UUID)
				if !ok {
					panic(r.currentPlayer)
				}
				r.Players[id].WordsGuessed++

				r.currentWord++
				r.Players[id].WordsTried++

				r.WordGuessed(id, msg.UserId)
				r.ReportUpdate()
			} else {
				b, err := json.Marshal(&ServerMessage{
					MsgType: WrongGuess,
					MsgData: map[string]any{},
				})
				if err != nil {
					panic(err)
				}
				r.Players[msg.UserId].toSend <- b
			}
		}
	case Finished:
		if msg.MsgType == GetState {
			r.sendState(msg.UserId)
		}
	}
}

func (r *Room) WordGuessed(explainer, guesser uuid.UUID) {
	for _, player := range r.Players {
		if !player.Ready {
			continue
		}

		if player.Id == explainer {
			b, err := json.Marshal(&ServerMessage{
				MsgType: WordGuessed,
				MsgData: map[string]any{
					"guesser": guesser,
				},
			})
			if err != nil {
				panic(err)
			}
			player.toSend <- b
		} else if player.Id == guesser {
			b, err := json.Marshal(&ServerMessage{
				MsgType: RightGuess,
				MsgData: map[string]any{
					"guesser": guesser,
				},
			})
			if err != nil {
				panic(err)
			}
			player.toSend <- b
		} else {
			b, err := json.Marshal(&ServerMessage{
				MsgType: WordGuessed,
				MsgData: map[string]any{
					"guesser": guesser,
				},
			})
			if err != nil {
				panic(err)
			}
			player.toSend <- b
		}
	}
}

func (r *Room) getWordAt(vocabs *Vocabularies, index int) string {
	if index < 0 || index >= len(r.Config.words) {
		return ""
	}

	i := r.Config.words[index]
	vocabs.lock.RLock()
	defer vocabs.lock.RUnlock()
	v := vocabs.vocabulary[r.Config.Language]

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
	b, err := json.Marshal(ServerMessage{
		CurrentState,
		r.ToMap(),
	})
	if err != nil {
		panic(err)
	}
	r.Players[id].toSend <- b
}

func (r *Room) ReportUpdate() {
	for _, player := range r.Players {
		if player.Ready {
			player.toSend <- NewUpdateMsg
		}
	}
}

var NewUpdateMsg, _ = json.Marshal(ServerMessage{
	MsgType: NewUpdate,
	MsgData: make(map[string]any),
})

func (r *Room) handleJoin(player *Player) { // TODO: check ring impl
	p, ok := r.Players[player.Id]
	if ok {
		p.Ready = true
		p.toSend = player.toSend
	} else {
		p = player
		r.Players[p.Id] = p
		p.Ready = true

		node := ring.New(1)
		node.Value = p.Id
		if r.currentPlayer == nil {
			r.currentPlayer = node
		} else {
			r.currentPlayer.Link(node)
		}
	}
	r.count++
	r.ReportUpdate()
}

func (r *Room) SaveState(postgres *Postgres) {
	err := postgres.UpdateRoomState(context.Background(), r)
	if err != nil {
		panic(err)
	}
}

func (r *Room) handleLeave(id uuid.UUID, postgres *Postgres, rooms *Rooms) {
	r.Players[id].Ready = false
	r.count--
	if r.count == 0 {
		r.SaveState(postgres)
		rooms.lock.RLock()
		defer rooms.lock.RUnlock()
		delete(rooms.rooms, r.Id)
		return
	}

	for _, player := range r.Players {
		if player.Ready {
			player.toSend <- NewUpdateMsg
		}
	}
}

type Rooms struct {
	rooms map[string]*Room
	lock  sync.RWMutex
}

func (r *Rooms) ServeWS(writer http.ResponseWriter, reader *http.Request, userId uuid.UUID, roomId string, postgres *Postgres, vocabs *Vocabularies) error {
	c, err := websocket.Accept(writer, reader, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // todo:
	})
	if err != nil {
		return err
	}
	defer c.CloseNow()

	r.lock.RLock()
	room, ok := r.rooms[roomId]
	r.lock.RUnlock()
	if !ok {
		// nobody is in room => room is not loaded
		newRoom, err := postgres.LoadRoom(context.Background(), roomId, vocabs)
		if err != nil {
			return err
		}
		room = newRoom
		r.lock.Lock()
		r.rooms[roomId] = room
		r.lock.Unlock()
		go room.Run(postgres, vocabs, r)
	}
	toSend := make(chan []byte, 20)
	player := &Player{Id: userId, toSend: toSend}

	room.join <- player
	defer func() { room.leave <- userId }()

	errc := make(chan error, 1)
	done := make(chan struct{}, 1)
	go func() {
		ping := time.NewTicker(pingTimeout())
		for {
			select {
			case msg := <-toSend:
				err = c.Write(context.Background(), websocket.MessageBinary, msg)
				if err != nil {
					errc <- err
					return
				}
			case <-ping.C:
				err := c.Ping(context.Background())
				if err != nil {
					errc <- err
					return
				}
			case <-done:
				return
			}
		}
	}()

	var v any
	for {
		err := wsjson.Read(context.Background(), c, &v)
		if err != nil {
			if !errors.Is(err, context.DeadlineExceeded) {
				done <- struct{}{}
				return err
			}
			continue
		}
		if len(errc) != 0 {
			return <-errc
		}
		m, ok := v.(map[string]any)
		if !ok {
			done <- struct{}{}
			return errors.New(fmt.Sprintf("invalid client msg: %s", v))
		}
		msg, err := toClientMessage(m)
		if err != nil {
			return err
		}
		room.ingest <- msg
	}
}

func toClientMessage(v map[string]any) (*ClientMessage, error) {
	u, ok := v["user_id"]
	if !ok {
		return nil, errors.New("invalid user id")
	}
	us, ok := u.(string)
	if !ok {
		return nil, errors.New("invalid user id")
	}
	user, err := uuid.Parse(us)
	if err != nil {
		return nil, err
	}

	t, ok := v["type"]
	if !ok {
		return nil, errors.New("invalid msg type")
	}
	ty, ok := t.(float64)
	if !ok {
		return nil, errors.New("invalid msg type")
	}
	typ := ClientMessageType(int(ty))

	d, ok := v["data"]
	if !ok {
		return nil, errors.New("invalid msg data")
	}
	da, ok := d.(map[string]any)
	if !ok {
		return nil, errors.New("invalid msg data")
	}
	return &ClientMessage{user, typ, da}, nil
}
