package main

import (
	"container/ring"
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

type WSClientMessageType int

const (
	GetState WSClientMessageType = iota
	TryGuess
)

type WSServerMessageType int

const (
	NewUpdate WSServerMessageType = iota
	GameState
)

type ClientMessage struct {
	UserId  uuid.UUID
	msgType WSClientMessageType
	msgData map[string]any
}

type ServerMessage struct {
	MsgType WSServerMessageType `json:"msg_type"`
	MsgData map[string]any      `json:"msg_data"`
}

type Player struct {
	Id     uuid.UUID
	ToSend chan *ServerMessage

	ready bool
	score uint
}

type Room struct {
	Admin   uuid.UUID
	Ingest  chan *ClientMessage
	Players map[uuid.UUID]*Player
	Count   uint
	Join    chan *Player
	Leave   chan uuid.UUID

	Config    *RoomConfig
	UsedWords map[int]struct{}

	CurrentPlayer *ring.Ring // ring of UUID
	CurrentWord   int
	Clock         time.Ticker
}

func (r *Room) Run() {
	for {
		select {
		case <-r.Clock.C:
			// todo: no word was guessed, move to next player

		case msg := <-r.Ingest:
			switch msg.msgType {
			case GetState:
				r.Players[msg.UserId].ToSend <- &ServerMessage{} // todo: snapshot state
			case TryGuess:
				// check if word

			}
		case player := <-r.Join:
			p, ok := r.Players[player.Id]
			if ok {
				p.ready = true
				p.ToSend = player.ToSend
			} else {
				p = player
				r.Players[p.Id] = p
				p.ready = true

				node := ring.New(1)
				node.Value = p.Id
				r.CurrentPlayer = node.Link(node)
			}
			for _, player := range r.Players {
				if player.ready {
					player.ToSend <- &ServerMessage{
						MsgType: NewUpdate,
					}
				}
			}
			r.Count++
		case id := <-r.Leave:
			r.Players[id].ready = false
			r.Count--
			if r.Count == 0 {
				// todo: save state to db

				return
			}

			for _, player := range r.Players {
				if player.ready {
					player.ToSend <- &ServerMessage{
						MsgType: NewUpdate,
					}
				}
			}
		}
	}
}

type Rooms struct {
	rooms map[string]*Room
	lock  sync.RWMutex
}

func (r *Rooms) ServeWS(writer http.ResponseWriter, reader *http.Request, userId uuid.UUID, roomId string, postgres *Postgres) error {
	c, err := upgrader.Upgrade(writer, reader, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	var rLock, wLock sync.Mutex

	ping := make(chan struct{})
	done := make(chan struct{})

	c.SetPingHandler(func(appData string) error {
		ping <- struct{}{}
		wLock.Lock()
		defer wLock.Unlock()
		return c.WriteMessage(websocket.PongMessage, []byte{})
	})

	go func() {
		t := time.NewTicker(pingTimeout())
		for {
			select {
			case <-t.C:
				done <- struct{}{}
				c.Close()
				return
			case <-ping:
				t.Reset(pingTimeout())
			}
		}
	}()

	r.lock.RLock()
	room, ok := r.rooms[roomId]
	r.lock.RUnlock()
	if !ok {
		// nobody is in room => room is not loaded
		newRoom, err := postgres.LoadRoom(context.Background(), roomId)
		if err != nil {
			return err
		}
		room = newRoom
		r.lock.Lock()
		r.rooms[roomId] = room
		r.lock.Unlock()
		go room.Run()
	}
	toSend := make(chan *ServerMessage, 20)
	player := &Player{Id: userId, ToSend: toSend}

	room.Join <- player
	defer func() { room.Leave <- userId }()

	for {
		rLock.Lock()
		msg, ok, err := tryReadMessage(c)
		rLock.Unlock()
		if err != nil {
			select {
			case <-done:
				return nil
			default:
				return err
			}
		}
		if !ok {
			for {
				if len(toSend) == 0 {
					break
				}
				wLock.Lock()
				msg := <-toSend
				m, err := json.Marshal(msg)
				if err != nil {
					// todo:
				}
				err = c.WriteMessage(websocket.BinaryMessage, m)
				if err != nil {
					// todo:
				}
				wLock.Unlock()
			}
			continue
		}

		var jsonMap map[string]interface{}
		if err := json.Unmarshal(msg, &jsonMap); err != nil {
			return err
		}

		msgType, ok := jsonMap["msg_type"].(int)
		if !ok {
			// todo:
		}
		msgData, ok := jsonMap["msg_data"].(map[string]any)
		if !ok {
			// todo:
		}

		room.Ingest <- &ClientMessage{
			userId,
			WSClientMessageType(msgType),
			msgData,
		}
	}
}
func tryReadMessage(conn *websocket.Conn) ([]byte, bool, error) {
	conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	defer conn.SetReadDeadline(time.Time{})

	mt, msg, err := conn.ReadMessage()
	if err != nil {
		if netErr, ok := errors.AsType[net.Error](err); ok && netErr.Timeout() {
			return nil, false, nil
		}
		return nil, false, err
	}
	if mt != websocket.BinaryMessage {
		return nil, false, errors.New("unexpected message type")
	}

	return msg, true, nil
}
