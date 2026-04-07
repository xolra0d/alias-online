package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	WorkersListName = "workers"
)

func isActiveIdentifier(worker string) string {
	return fmt.Sprintf("worker:%s:active", worker)
}

func WorkerRoomsIdentifier(worker string) string {
	return fmt.Sprintf("worker:%s:rooms", worker)
}

func RoomLockIdentifier(roomId string) string {
	return fmt.Sprintf("room:%s:lock", roomId)
}

type Database struct {
	client *redis.Client
}

func NewDatabase(addr, username, password string, db int) *Database {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})
	return &Database{client: client}
}

func (d *Database) GetAllWorkers(ctx context.Context) ([]string, error) {
	return d.client.LRange(ctx, WorkersListName, 0, -1).Result()
}

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

func (d *Database) SetWorkerActive(ctx context.Context, worker string, exp time.Duration) error {
	return d.client.Set(ctx, isActiveIdentifier(worker), "1", exp).Err() // any value is fine
}

func (d *Database) GetWorkerRoomCount(ctx context.Context, worker string) (int, error) {
	c, err := d.client.HLen(ctx, WorkerRoomsIdentifier(worker)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}
		return 0, err
	}
	return int(c), err
}

func (d *Database) TryReserveRoom(ctx context.Context, roomId, workerId string) (string, error) {
	SetNXOrPrevious := redis.NewScript(`
local existing = redis.call('GET', KEYS[1])
if existing then
	return existing
end
redis.call('SET', KEYS[1], ARGV[1])
return ARGV[1]
    `)
	return SetNXOrPrevious.Run(ctx, d.client, []string{RoomLockIdentifier(roomId)}, workerId).Text()
}

func (d *Database) ProlongRoom(ctx context.Context, roomId, workerId string, exp time.Duration) error {
	return d.client.Set(ctx, RoomLockIdentifier(roomId), workerId, exp).Err()
}

func (d *Database) AddRoomToWorker(ctx context.Context, roomId, workerIp string) error {
	return d.client.SAdd(ctx, WorkerRoomsIdentifier(workerIp), roomId).Err()
}

func (d *Database) ReleaseRoom(ctx context.Context, roomId, workerIp string) error {
	return d.client.SRem(ctx, WorkerRoomsIdentifier(workerIp), roomId).Err()
}
