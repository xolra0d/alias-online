package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"maps"
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
	}
}

func (rooms *Rooms) GetVocab(ctx context.Context, name string) (Vocabulary, error) {
	v, err := rooms.vocabManagerClient.GetVocab(ctx, &pbVocabManager.GetVocabRequest{Name: name})
	if err != nil {
		return Vocabulary{}, err
	}
	return Vocabulary{PrimaryWords: v.GetPrimaryWords(), RudeWords: v.GetRudeWords()}, nil

}

func (rooms *Rooms) RunWS(w http.ResponseWriter, r *http.Request, roomId, username, name string) error {
	rooms.lock.RLock()
	loaded, ok := rooms.rooms[roomId]
	rooms.lock.RUnlock()

	if !ok {
		ctx, cancel := context.WithTimeout(r.Context(), rooms.LoadRoomTimeout)
		newRoom, err := rooms.postgres.LoadRoom(ctx, roomId, rooms.GetVocab)
		cancel()
		if err != nil {
			rooms.logger.Error("failed to load room", "roomId", roomId, "err", err)
			return err
		}

		resp, err := rooms.roomManagerClient.RegisterRoom(ctx, &pbRoomManager.RegisterRoomRequest{
			RoomId: roomId,
			Worker: rooms.RunningAddr,
		})
		if err != nil {
			rooms.logger.Error("failed to register room", "roomId", roomId, "err", err)
			return err
		}
		if resp.Worker != rooms.RunningAddr {
			// other worker reserved room
			rooms.logger.Info("failed to register room", "roomId", roomId, "otherWorker", resp.Worker)
			return rooms.UpdateToWebsocketsAndErr(w, r, resp.GetWorker())
		}
		loaded = newRoom
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

	if !ok {
		rooms.lock.Lock()
		existingRoom, exists := rooms.rooms[roomId]
		if exists {
			loaded = existingRoom
			rooms.lock.Unlock()
		} else {
			rooms.rooms[roomId] = loaded
			rooms.lock.Unlock()
			go loaded.Run(func(room *Room) {
				ctx, cancel := context.WithTimeout(context.Background(), rooms.SaveRoomTimeout)
				defer cancel()

				err := rooms.postgres.SaveRoom(ctx, room)
				if err != nil {
					rooms.logger.Error("failed to save room", "roomId", roomId, "err", err)
				}
				_, err = rooms.roomManagerClient.ReleaseRoom(ctx, &pbRoomManager.ReleaseRoomRequest{RoomId: roomId})
				if err != nil {
					rooms.logger.Error("failed to release room", "roomId", roomId, "err", err)
				}

				rooms.lock.Lock()
				delete(rooms.rooms, roomId)
				rooms.lock.Unlock()

			})
		}
	}

	player := NewPlayer(username, name, 0, 0)

	loaded.Join(player)
	defer loaded.Leave(username)

	// context is shared between reader and writer and is used to tell to stop, when other is done.
	ctx, cancel := context.WithCancel(r.Context())

	go loaded.RunWriter(ctx, cancel, c, player, rooms.WsWriteTimeout, rooms.WsPingTimeout)
	return loaded.RunReader(ctx, cancel, c, player, rooms.MaxMessagesPerSecond)
}

func (rooms *Rooms) UpdateToWebsocketsAndErr(w http.ResponseWriter, r *http.Request, otherWorker string) error {
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
	conn, err := grpc.NewClient(roomManagerUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to connect to room manager", "roomManagerUrl", roomManagerUrl, "err", err)
		return nil, nil, err
	}
	//defer conn.Close()
	return pbRoomManager.NewRoomManagerServiceClient(conn), conn.Close, nil
}

func NewVocabManagerClient(vocabManagerUrl string, logger *slog.Logger) (pbVocabManager.VocabManagerServiceClient, func() error, error) {
	conn, err := grpc.NewClient(vocabManagerUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("failed to connect to vocab manager", "vocabManagerUrl", vocabManagerUrl, "err", err)
		return nil, nil, err
	}
	//defer conn.Close()
	return pbVocabManager.NewVocabManagerServiceClient(conn), conn.Close, nil
}
