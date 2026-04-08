package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	mrand "math/rand"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/coder/websocket"
	pbRoomManager "github.com/xolra0d/alias-online/shared/proto/room_manager"
	pbVocabManager "github.com/xolra0d/alias-online/shared/proto/vocab_manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Rooms struct {
	rooms map[string]*Room
	lock  sync.RWMutex

	postgres           *Postgres
	logger             *slog.Logger
	roomManagerClient  pbRoomManager.RoomManagerServiceClient
	vocabManagerClient pbVocabManager.VocabManagerServiceClient

	RunningAddr          string
	WsOriginPatterns     []string
	LoadRoomTimeout      time.Duration
	SaveRoomTimeout      time.Duration
	WsWriteTimeout       time.Duration
	WsPingTimeout        time.Duration
	MaxMessagesPerSecond int
	MaxClockValue        int
	LoadVocabTimeout     time.Duration
}

func NewRooms(
	postgres *Postgres,
	logger *slog.Logger,
	roomManagerClient pbRoomManager.RoomManagerServiceClient,
	vocabManagerClient pbVocabManager.VocabManagerServiceClient,

	runningAddr string,
	WsOriginPatterns []string,
	LoadRoomTimeout time.Duration,
	SaveRoomTimeout time.Duration,
	WsWriteTimeout time.Duration,
	WsPingTimeout time.Duration,
	MaxMessagesPerSecond int,
	maxClockValue int,
	loadVocabTimeout time.Duration,
) *Rooms {
	return &Rooms{
		rooms: make(map[string]*Room),

		postgres:           postgres,
		logger:             logger,
		roomManagerClient:  roomManagerClient,
		vocabManagerClient: vocabManagerClient,

		RunningAddr:          runningAddr,
		WsOriginPatterns:     WsOriginPatterns,
		LoadRoomTimeout:      LoadRoomTimeout,
		SaveRoomTimeout:      SaveRoomTimeout,
		WsWriteTimeout:       WsWriteTimeout,
		WsPingTimeout:        WsPingTimeout,
		MaxMessagesPerSecond: MaxMessagesPerSecond,
		MaxClockValue:        maxClockValue,
		LoadVocabTimeout:     loadVocabTimeout,
	}
}

// GetVocab gets vocab from vocab service.
func (rooms *Rooms) GetVocab(ctx context.Context, name string) (Vocabulary, error) {
	v, err := rooms.vocabManagerClient.GetVocab(ctx, &pbVocabManager.GetVocabRequest{Name: name})
	if err != nil {
		return Vocabulary{}, err
	}
	return Vocabulary{PrimaryWords: v.GetPrimaryWords(), RudeWords: v.GetRudeWords()}, nil

}

func mapToRoomConfig(ctx context.Context, m map[string]any, maxClockValue int, getVocab func(ctx context.Context, name string) (Vocabulary, error)) (*RoomConfig, error) {
	rW, ok := m["rude-words"]
	if !ok {
		return nil, fmt.Errorf("no rude words")
	}
	rudeWords, ok := rW.(bool)
	if !ok {
		return nil, fmt.Errorf("invalid rude words: %v", rW)
	}
	aW, ok := m["additional-vocabulary"]
	if !ok {
		return nil, fmt.Errorf("no additional-vocabulary")
	}
	aV, ok := aW.([]any)
	if !ok {
		return nil, fmt.Errorf("invalid additional-vocabulary: %v", aW)
	}
	additionalVocabulary := make([]string, 0, len(aV))
	for _, w := range aV {
		v, ok := w.(string)
		if !ok {
			return nil, fmt.Errorf("invalid additional-vocabulary: %v", w)
		}
		additionalVocabulary = append(additionalVocabulary, v)
	}

	c, ok := m["clock"]
	if !ok {
		return nil, fmt.Errorf("no clock")
	}
	clock, ok := c.(float64)
	if !ok || clock <= 0 || clock > float64(maxClockValue) {
		return nil, fmt.Errorf("invalid clock: %v", c)
	}
	v, ok := m["language"]
	if !ok {
		return nil, fmt.Errorf("no vocab language")
	}
	vo, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("invalid vocab language: %v", v)
	}

	vocab, err := getVocab(ctx, vo)
	if err != nil {
		return nil, err
	}

	seed := mrand.Int()
	allWords := vocab.PrimaryWords
	if rudeWords {
		allWords = append(allWords, vocab.RudeWords...)
	}
	allWords = append(allWords, additionalVocabulary...)

	r := mrand.New(mrand.NewSource(int64(seed)))
	r.Shuffle(len(allWords), func(i, j int) {
		allWords[i], allWords[j] = allWords[j], allWords[i]
	})

	return &RoomConfig{
		Seed:                 seed,
		AllWords:             allWords,
		Language:             vo,
		RudeWords:            rudeWords,
		AdditionalVocabulary: additionalVocabulary,
		Clock:                int(clock),
	}, nil
}

// RunWS initiates a WS and configuring a room.
func (rooms *Rooms) RunWS(w http.ResponseWriter, r *http.Request, roomId, username, name string) error {
	rooms.lock.Lock()
	room, ok := rooms.rooms[roomId]
	if !ok {
		room = NewPreparingRoom()
		rooms.rooms[roomId] = room
	}
	rooms.lock.Unlock()

	if ok {
		if err := room.WaitUntilOperational(); err != nil {
			return rooms.UpdateToWebsocketsAndRedirect(w, r, err.Error())
		}
	} else {
		ctx, cancel := context.WithTimeout(r.Context(), rooms.LoadRoomTimeout)
		resp, err := rooms.roomManagerClient.RegisterRoom(ctx, &pbRoomManager.RegisterRoomRequest{
			RoomId: roomId,
			Worker: rooms.RunningAddr,
		})
		cancel()
		if err != nil {
			room.SetErrored()
			rooms.logger.Error("failed to register room", "roomId", roomId, "err", err)
			return err
		}
		if resp.Worker != rooms.RunningAddr {
			room.SetErrored()
			// other worker reserved room
			rooms.logger.Info("failed to register room", "roomId", roomId, "otherWorker", resp.Worker)
			return rooms.UpdateToWebsocketsAndRedirect(w, r, resp.GetWorker())
		}
		go room.Run(func(room *Room) {
			ctx, cancel := context.WithTimeout(context.Background(), rooms.SaveRoomTimeout)
			defer cancel()

			err := rooms.postgres.SaveRoom(ctx, room)
			if err != nil {
				rooms.logger.Error("failed to save room", "roomId", roomId, "err", err)
			}
			_, err = rooms.roomManagerClient.ReleaseRoom(ctx, &pbRoomManager.ReleaseRoomRequest{RoomId: roomId, Worker: rooms.RunningAddr})
			if err != nil {
				rooms.logger.Error("failed to release room", "roomId", roomId, "err", err)
			}

			rooms.lock.Lock()
			delete(rooms.rooms, roomId)
			rooms.lock.Unlock()
		})
	}

	acceptOptions := &websocket.AcceptOptions{}
	if len(rooms.WsOriginPatterns) > 0 {
		acceptOptions.OriginPatterns = rooms.WsOriginPatterns
	}
	c, err := websocket.Accept(w, r, acceptOptions)

	if err != nil {
		return err
	}

	defer c.CloseNow()

	mt, mb, err := c.Read(r.Context())
	if mt != websocket.MessageText {
		if !ok {
			room.SetErrored()
		}
		rooms.logger.Error("failed to read message, message is not text", "roomId", roomId, "msg", mb, "err", err)
		c.Close(websocket.StatusUnsupportedData, "unsupported message")
		return nil
	}
	var msg ClientMessage
	err = json.Unmarshal(mb, &msg)
	if err != nil {
		if !ok {
			room.SetErrored()
		}
		rooms.logger.Error("failed to parse message as json", "roomId", roomId, "msg", mb, "err", err)
		c.Close(websocket.StatusUnsupportedData, "invalid message")
		return nil
	}
	switch msg.MsgType {
	case LoadRoom:
		if !ok {
			ctx, cancel := context.WithTimeout(r.Context(), rooms.LoadRoomTimeout)
			newRoom, err := rooms.postgres.LoadRoom(ctx, roomId, rooms.GetVocab)
			cancel()
			if err != nil {
				room.SetErrored()
				rooms.logger.Error("failed to load room", "roomId", roomId, "err", err)
				c.Close(websocket.StatusUnsupportedData, "invalid message")
				return nil
			}
			rooms.lock.Lock()
			room = rooms.rooms[roomId]
			room.UpdateStateFromRoom(newRoom)
			rooms.lock.Unlock()
		}
	case CreateRoom:
		if !ok {
			ctx, cancel := context.WithTimeout(r.Context(), rooms.LoadVocabTimeout)
			cfg, err := mapToRoomConfig(ctx, msg.MsgData, rooms.MaxClockValue, rooms.GetVocab)
			cancel()
			if err != nil {
				room.SetErrored()
				rooms.logger.Error("failed to parse room config", "roomId", roomId, "msg", msg.MsgData, "err", err)
				c.Close(websocket.StatusUnsupportedData, "invalid message")
				return nil
			}
			rooms.lock.Lock()
			room = rooms.rooms[roomId]
			room.UpdateStateFromRoomConfig(roomId, name, username, cfg, rooms.logger)
			rooms.lock.Unlock()
		}
	default:
		if !ok {
			room.SetErrored()
		}
		rooms.logger.Error("failed to read message, message type is invalid", "roomId", roomId, "mt", mt, "msg", mb, "err", err)
		c.Close(websocket.StatusUnsupportedData, "unsupported message")
		return nil
	}

	// I am the one, who initiated a room load - I am the first to play, other shall wait!
	player := NewPlayer(username, name, 0, 0)
	room.Join(player)
	defer room.Leave(username)

	if !ok {
		room.SetOperational()
	}

	// context is shared between reader and writer and is used to tell to stop, when other is done.
	ctx, cancel := context.WithCancel(r.Context()) // todo, REALLY R?

	go room.RunWriter(ctx, cancel, c, player, rooms.WsWriteTimeout, rooms.WsPingTimeout)
	room.RunReader(ctx, cancel, c, player, rooms.MaxMessagesPerSecond)
	return nil
}

func (rooms *Rooms) UpdateToWebsocketsAndRedirect(w http.ResponseWriter, r *http.Request, otherWorker string) error {
	acceptOptions := &websocket.AcceptOptions{}
	if len(rooms.WsOriginPatterns) > 0 {
		acceptOptions.OriginPatterns = rooms.WsOriginPatterns
	}
	c, err := websocket.Accept(w, r, acceptOptions)

	if err != nil {
		return err
	}
	defer c.CloseNow()

	msg, err := json.Marshal(&ServerMessage{
		MsgType: Redirect,
		MsgData: map[string]any{
			"worker": otherWorker,
		},
	})
	if err != nil {
		return err
	}

	writeCtx, writeCancel := context.WithTimeout(context.Background(), rooms.WsWriteTimeout)
	err = c.Write(writeCtx, websocket.MessageBinary, msg)
	writeCancel()
	if err != nil {
		rooms.logger.Error("write error", "playerId", "msg", msg, "err", err)
		return err
	}
	return nil
}

func (rooms *Rooms) ReportLoadedRooms() []string {
	rooms.lock.RLock()
	defer rooms.lock.RUnlock()
	return slices.Collect(maps.Keys(rooms.rooms))
}

func (rooms *Rooms) RunPinger(logger *slog.Logger, PollInterval time.Duration, runningAddr string, shouldStop, done chan struct{}) {
	for {
		ctx, cancel := context.WithCancel(context.Background())

		_, err := rooms.roomManagerClient.PingWorker(ctx, &pbRoomManager.PingWorkerRequest{Worker: runningAddr, LoadedRooms: rooms.ReportLoadedRooms()})
		cancel()
		if err != nil {
			// todo: think about trying to find out other room managers
			logger.Error("failed to ping room manager", "err", err)
			logger.Info("stopped grpc pinger")
			<-shouldStop
			done <- struct{}{}
			return
		}

		timer := time.NewTimer(PollInterval)
		select {
		case <-timer.C:
		case <-shouldStop:
			logger.Info("stopping grpc client")
			done <- struct{}{}
			return
		}
	}
}

func NewRoomManagerClient(roomManagerUrl string, logger *slog.Logger) (pbRoomManager.RoomManagerServiceClient, func() error, error) {
	conn, err := grpc.NewClient(
		roomManagerUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		logger.Error("failed to connect to room manager", "roomManagerUrl", roomManagerUrl, "err", err)
		return nil, nil, err
	}
	//defer conn.Close()
	return pbRoomManager.NewRoomManagerServiceClient(conn), conn.Close, nil
}

func NewVocabManagerClient(vocabManagerUrl string, logger *slog.Logger) (pbVocabManager.VocabManagerServiceClient, func() error, error) {
	conn, err := grpc.NewClient(
		vocabManagerUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		logger.Error("failed to connect to vocab manager", "vocabManagerUrl", vocabManagerUrl, "err", err)
		return nil, nil, err
	}
	//defer conn.Close()
	return pbVocabManager.NewVocabManagerServiceClient(conn), conn.Close, nil
}
