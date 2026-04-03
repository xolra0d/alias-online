package room_manager

import (
	"context"
	"maps"
	"slices"
	"sync"
	"time"
)

type RoomManager struct {
	workers map[string]int // map of addrPort: roomCount
	lock    sync.RWMutex

	database *database
}

func NewManager(
	database *database,
) *RoomManager {
	return &RoomManager{
		workers:  map[string]int{},
		database: database,
	}
}

func (m *RoomManager) FindRoomIP(ctx context.Context, roomId string) (string, error) {
	const op = "room_manager.FindRoomIP"

	bestIP := ""
	lowestCount := 0
	m.lock.RLock()
	for ip, count := range m.workers {
		if count < lowestCount {
			bestIP = ip
			lowestCount = count
		}
	}
	m.lock.RUnlock()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	m.lock.RLock()
	workerIP, err := m.database.CheckOrAddRoomID(ctx, roomId, bestIP, slices.Collect(maps.Keys(m.workers)))
	m.lock.RUnlock()
	cancel()

	if err != nil {
		return "", err
	}
	return workerIP, nil
}

func (m *RoomManager) SearchForNewRoomWorkers() {
	const op = "room_manager.SearchForNewRoomWorkers"

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		workers, err := m.database.GetAllWorkers(ctx)
		cancel()
		if err != nil {
			// todo!!
		}

		lookup := make(map[string]int, len(workers))
		for _, worker := range workers {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			c, err := m.database.GetWorkerRoomCount(ctx, worker)
			cancel()
			if err != nil {
				// todo!!!
			}
			lookup[worker] = c
		}

		eq := true
		m.lock.RLock()
		for w := range m.workers {
			if _, ok := lookup[w]; !ok {
				eq = false
				break
			}
		}
		if len(lookup) != len(m.workers) {
			eq = false
		}
		m.lock.RUnlock()

		if !eq {

			m.lock.Lock()
			for worker := range m.workers {
				if _, ok := lookup[worker]; !ok {
					delete(m.workers, worker)
				}
			}
			for workerIP := range lookup {
				if _, ok := m.workers[workerIP]; !ok {
					m.workers[workerIP] = 0
				}
			}
			m.lock.Unlock()
		}

		time.Sleep(time.Second * 10)
	}
}

//room_worker.lock.RLock()
//loaded, ok := room_worker.room_worker[roomId]
//room_worker.lock.RUnlock()
//
//if !ok {
//ctx, cancel := context.WithTimeout(context.Background(), room_worker.loadRoomTimeout)
//newRoom, err := postgres.LoadRoom(ctx, roomId, vocabs)
//cancel()
//if err != nil {
//return err
//}
//loaded = newRoom
//}
//
//acceptOptions := &websocket.AcceptOptions{}
//if len(room_worker.wsOriginPatterns) > 0 {
//acceptOptions.OriginPatterns = room_worker.wsOriginPatterns
//}
//c, err := websocket.Accept(w, r, acceptOptions)
//
//if err != nil {
//return err
//}
//defer func(c *websocket.Conn) {
//	err := c.CloseNow()
//	if err != nil {
//		room_worker.logger.Error(op, "error closing websocket connection", "error", err)
//	}
//}(c)
//
//if !ok {
//room_worker.lock.Lock()
//existingRoom, exists := room_worker.room_worker[roomId]
//if exists {
//loaded = existingRoom
//} else {
//room_worker.room_worker[roomId] = loaded
//go loaded.Run(vocabs, func(room_worker *room_worker.Room) error {
//ctx, cancel := context.WithTimeout(context.Background(), room_worker.saveRoomTimeout)
//defer cancel()
//
//err := postgres.SaveRoomSnapshot(ctx, room_worker)
//
//room_worker.lock.Lock()
//delete(room_worker.room_worker, roomId)
//room_worker.lock.Unlock()
//
//return err
//})
//}
//room_worker.lock.Unlock()
//}
//
//player := room_worker.NewPlayer(userId, name, 0, 0)
//
//loaded.Join(player)
//defer loaded.Leave(userId)
//
//// context is shared between reader and writer and is used to tell to stop via ctx.Done(), when other is done.
//ctx, cancel := context.WithCancel(r.Context())
//
//go loaded.RunWriter(ctx, cancel, c, player, room_worker.wsWriteTimeout, room_worker.pingTimeout)
//return loaded.RunReader(ctx, cancel, c, player, room_worker.maxMessagesPerSecond)
