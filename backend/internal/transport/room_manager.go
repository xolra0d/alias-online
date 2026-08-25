package transport

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"github.com/xolra0d/alias-online/internal/config"
	"github.com/xolra0d/alias-online/internal/database"
	"github.com/xolra0d/alias-online/internal/room"
)

type Rooms struct {
	rooms  map[string]*room.Room
	logger *config.Logger
	lock   sync.RWMutex

	minClock                     int
	maxClock                     int
	maxAdditionalVocabularyWords int
	maxAdditionalWordLength      int

	wsOriginPatterns     []string
	maxMessagesPerSecond int
	pingTimeout          time.Duration
	wsWriteTimeout       time.Duration
	loadRoomTimeout      time.Duration
	saveRoomTimeout      time.Duration
}

func NewRooms(
	rooms map[string]*room.Room,
	logger *config.Logger,
	minClock int,
	maxClock int,
	maxAdditionalVocabularyWords int,
	maxAdditionalWordLength int,
	wsOriginPatterns []string,
	maxMessagesPerSecond int,
	pingTimeout time.Duration,
	wsWriteTimeout time.Duration,
	loadRoomTimeout time.Duration,
	saveRoomTimeout time.Duration,
) *Rooms {
	return &Rooms{
		rooms:                        rooms,
		logger:                       logger,
		minClock:                     minClock,
		maxClock:                     maxClock,
		maxAdditionalVocabularyWords: maxAdditionalVocabularyWords,
		maxAdditionalWordLength:      maxAdditionalWordLength,
		wsOriginPatterns:             wsOriginPatterns,
		maxMessagesPerSecond:         maxMessagesPerSecond,
		pingTimeout:                  pingTimeout,
		wsWriteTimeout:               wsWriteTimeout,
		loadRoomTimeout:              loadRoomTimeout,
		saveRoomTimeout:              saveRoomTimeout,
	}
}

func (rooms *Rooms) ServeWS(w http.ResponseWriter, r *http.Request, userId uuid.UUID, name, roomId string, postgres *database.Postgres, vocabs *room.Vocabularies) error {
	const op = "transport.ServeWS"

	rooms.lock.RLock()
	loaded, ok := rooms.rooms[roomId]
	rooms.lock.RUnlock()

	if !ok {
		ctx, cancel := context.WithTimeout(context.Background(), rooms.loadRoomTimeout)
		newRoom, err := postgres.LoadRoom(ctx, roomId, vocabs)
		cancel()
		if err != nil {
			return err
		}
		loaded = newRoom
	}

	acceptOptions := &websocket.AcceptOptions{}
	if len(rooms.wsOriginPatterns) > 0 {
		acceptOptions.OriginPatterns = rooms.wsOriginPatterns
	}
	c, err := websocket.Accept(w, r, acceptOptions)

	if err != nil {
		return err
	}
	defer func(c *websocket.Conn) {
		err := c.CloseNow()
		if err != nil {
			rooms.logger.Error(op, "error closing websocket connection", "error", err)
		}
	}(c)

	if !ok {
		rooms.lock.Lock()
		existingRoom, exists := rooms.rooms[roomId]
		if exists {
			loaded = existingRoom
		} else {
			rooms.rooms[roomId] = loaded
			go loaded.Run(vocabs, func(room *room.Room) error {
				ctx, cancel := context.WithTimeout(context.Background(), rooms.saveRoomTimeout)
				defer cancel()

				err := postgres.SaveRoomSnapshot(ctx, room)

				rooms.lock.Lock()
				delete(rooms.rooms, roomId)
				rooms.lock.Unlock()

				return err
			})
		}
		rooms.lock.Unlock()
	}

	player := room.NewPlayer(userId, name, 0, 0)

	loaded.Join(player)
	defer loaded.Leave(userId)

	// context is shared between reader and writer and is used to tell to stop via ctx.Done(), when other is done.
	ctx, cancel := context.WithCancel(r.Context())

	go loaded.RunWriter(ctx, cancel, c, player, rooms.wsWriteTimeout, rooms.pingTimeout)
	return loaded.RunReader(ctx, cancel, c, player, rooms.maxMessagesPerSecond)
}
