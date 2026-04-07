package main

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"
)

type client struct {
	IpAddr    string
	roomCount int
	Exp       time.Time
}

type RoomManager struct {
	workers map[string]*client
	lock    sync.RWMutex

	db     *Database
	logger *slog.Logger

	done       chan struct{}
	doneCtx    context.Context
	doneCancel context.CancelFunc

	PollInterval                 time.Duration
	WorkerExpiry                 time.Duration
	RetrieveActiveWorkersTimeout time.Duration
}

func NewManager(
	database *Database,
	logger *slog.Logger,

	PollInterval time.Duration,
	WorkerExpiry time.Duration,
	RetrieveActiveWorkersTimeout time.Duration,
) *RoomManager {
	ctx, cancel := context.WithTimeout(context.Background(), RetrieveActiveWorkersTimeout)
	w, err := retrieveActiveWorkers(ctx, database, time.Now().Add(WorkerExpiry))
	cancel()
	if err != nil {
		logger.Error("Failed to retrieve active workers during manager startup", "error", err)
		os.Exit(1)
	}

	ctx, cancel = context.WithCancel(context.Background())
	return &RoomManager{
		workers: w,

		logger: logger,
		db:     database,

		done:       make(chan struct{}),
		doneCtx:    ctx,
		doneCancel: cancel,

		PollInterval:                 PollInterval,
		WorkerExpiry:                 WorkerExpiry,
		RetrieveActiveWorkersTimeout: RetrieveActiveWorkersTimeout,
	}
}

func retrieveActiveWorkers(ctx context.Context, db *Database, WorkerExpiry time.Time) (map[string]*client, error) {
	workers, err := db.GetAllWorkers(ctx)
	if err != nil {
		return nil, err
	}

	activeWorkers := map[string]*client{}

	for _, w := range workers {
		active, err := db.IsWorkerActive(ctx, w)
		if err != nil {
			return nil, err
		}
		if !active {
			continue
		}
		roomCount, err := db.GetWorkerRoomCount(ctx, w)
		if err != nil {
			return nil, err
		}
		activeWorkers[w] = &client{w, roomCount, WorkerExpiry}
	}

	return activeWorkers, nil
}

func (m *RoomManager) ScanForNewWorkers() {
	for {
		newExp := time.Now().Add(m.WorkerExpiry)

		ctx, cancel := context.WithTimeout(m.doneCtx, m.RetrieveActiveWorkersTimeout)
		activeWorkers, err := retrieveActiveWorkers(ctx, m.db, newExp)
		cancel()
		if err != nil {
			if ctx.Err() == nil {
				m.logger.Error("could not retrieve active workers", "err", err)
			}
			m.done <- struct{}{}
			return
		}

		workersAdded := []string{}
		workersDeleted := []string{}

		m.lock.Lock()
		for id, w := range activeWorkers {
			if c, ok := m.workers[id]; ok {
				c.Exp = w.Exp
				c.roomCount = w.roomCount
			} else {
				m.workers[id] = &client{w.IpAddr, w.roomCount, newExp}
				workersAdded = append(workersAdded, id)
			}
		}
		for id, w := range m.workers {
			if w.Exp.Before(time.Now()) {
				delete(m.workers, id)
				workersDeleted = append(workersDeleted, id)
			}
		}
		m.lock.Unlock()

		for _, w := range workersAdded {
			m.logger.Info("Added worker", "worker", w)
		}
		for _, w := range workersDeleted {
			m.logger.Info("Removed worker", "worker", w)
		}

		timer := time.NewTimer(m.PollInterval)
		select {
		case <-m.doneCtx.Done():
			timer.Stop()
			m.done <- struct{}{}
			return
		case <-timer.C:
		}
	}
}

func (m *RoomManager) StopScanForNewWorkers() {
	m.logger.Info("stopping observation loop")
	m.doneCancel()
	<-m.done
	m.logger.Info("stopped observation loop")
}

func (m *RoomManager) SetWorkerActive(ctx context.Context, worker string) error {
	return m.db.SetWorkerActive(ctx, worker, m.WorkerExpiry)
}

func (m *RoomManager) FindMostFreeWorker() string {
	bestIp := ""
	lowestCount := 1 << 30

	m.lock.RLock()
	for ip, c := range m.workers {
		if c.roomCount < lowestCount {
			bestIp = ip
			lowestCount = c.roomCount
		}
	}
	m.lock.RUnlock()

	return bestIp
}

func (m *RoomManager) FindBestWorker(ctx context.Context, roomId string) (string, error) {
	optimal := m.FindMostFreeWorker()
	worker, err := m.db.TryReserveRoom(ctx, roomId, optimal)
	if err != nil {
		return "", err
	}
	return worker, nil
}

func (m *RoomManager) RegisterRoom(ctx context.Context, roomId, workerIp string) (string, error) {
	worker, err := m.db.TryReserveRoom(ctx, roomId, workerIp)
	if err != nil {
		return "", err
	}
	if worker != workerIp {
		return worker, nil
	}

	// it's us, who reserved the room
	err = m.db.AddRoomToWorker(ctx, roomId, workerIp)
	if err != nil {
		return "", err
	}

	m.lock.Lock()
	m.workers[workerIp].roomCount++
	m.lock.Unlock()

	return worker, nil
}

func (m *RoomManager) ProlongRoom(ctx context.Context, roomId, workerIp string) error {
	return m.db.ProlongRoom(ctx, roomId, workerIp, time.Second*30)
}

func (m *RoomManager) ReleaseRoom(ctx context.Context, roomId, workerIp string) error {
	err := m.db.ReleaseRoom(ctx, roomId, workerIp)
	if err != nil {
		return err
	}

	m.lock.Lock()
	m.workers[workerIp].roomCount--
	m.lock.Unlock()

	return nil
}
