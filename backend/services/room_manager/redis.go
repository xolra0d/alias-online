package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	WorkersListName = "workers" // name for list in redis
)

func isActiveIdentifier(worker string) string {
	return fmt.Sprintf("worker:%s:active", worker)
}

func workerRoomsIdentifier(worker string) string {
	return fmt.Sprintf("worker:%s:rooms", worker)
}

func roomLockIdentifier(roomId string) string {
	return fmt.Sprintf("room:%s:lock", roomId)
}

type Database struct {
	client *redis.Client
}

// NewDatabase creates a new redis client.
func NewDatabase(addr, username, password string, db int) *Database {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})
	return &Database{client: client}
}

// GetAllWorkers returns all ever active workers.
func (d *Database) GetAllWorkers(ctx context.Context) ([]string, error) {
	return d.client.LRange(ctx, WorkersListName, 0, -1).Result()
}

// IsWorkerActive checks if specific worker is currently active.
func (d *Database) IsWorkerActive(ctx context.Context, worker string) (bool, error) {
	active, err := d.client.Get(ctx, isActiveIdentifier(worker)).Bool()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return active, nil
}

// SetWorkerActive sets worker active with timeout.
func (d *Database) SetWorkerActive(ctx context.Context, worker string, exp time.Duration) error {
	registerAndSetActive := redis.NewScript(`
local idx = redis.call('LPOS', KEYS[1], ARGV[1])
if not idx then
  redis.call('RPUSH', KEYS[1], ARGV[1])
end
redis.call('SET', KEYS[2], "1", "EX", ARGV[2])
return 1
`)
	_, err := registerAndSetActive.Run(
		ctx,
		d.client,
		[]string{WorkersListName, isActiveIdentifier(worker)},
		worker,
		strconv.Itoa(int(exp.Seconds())),
	).Result()
	return err
}

// GetWorkerRoomCount returns rooms worker currently holds.
func (d *Database) GetWorkerRoomCount(ctx context.Context, worker string) (int, error) {
	c, err := d.client.SMembers(ctx, workerRoomsIdentifier(worker)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return len(c), err
}

// TryReserveRoom tries to atomically reserve room.
// If succeeded, returns `workerId`. If failed, because other worker already reserved this room, returns their worker id.
func (d *Database) TryReserveRoom(ctx context.Context, roomId, workerId string) (string, error) {
	SetNXOrPrevious := redis.NewScript(`
local existing = redis.call('GET', KEYS[1])
if existing then
	return existing
end
redis.call('SET', KEYS[1], ARGV[1])
return ARGV[1]
    `)
	return SetNXOrPrevious.Run(ctx, d.client, []string{roomLockIdentifier(roomId)}, workerId).Text()
}

// ProlongRoom prolongs lease of room for worker.
func (d *Database) ProlongRoom(ctx context.Context, roomId, workerId string, exp time.Duration) error {
	return d.client.Set(ctx, roomLockIdentifier(roomId), workerId, exp).Err()
}

// AddRoomToWorker registers room under workerIp pool of rooms.
func (d *Database) AddRoomToWorker(ctx context.Context, roomId, workerIp string) error {
	return d.client.SAdd(ctx, workerRoomsIdentifier(workerIp), roomId).Err()
}

// ReleaseRoom removes room from workerIp pool of rooms.
func (d *Database) ReleaseRoom(ctx context.Context, roomId, workerIp string) error {
	return d.client.SRem(ctx, workerRoomsIdentifier(workerIp), roomId).Err()
}
