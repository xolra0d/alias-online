package main

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
)

func newMockDatabase() (*Database, redismock.ClientMock) {
	client, mock := redismock.NewClientMock()
	return &Database{client: client}, mock
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func reserveRoomScript() string {
	//	return `
	//local existing = redis.call('GET', KEYS[1])
	//if existing then
	//	return existing
	//end
	//redis.call('SET', KEYS[1], ARGV[1])
	//return ARGV[1]
	//    `
	return "1ac7dac708231c16b2455eb59f223a866dd98b8e"
}

func setWorkerActiveScript() string {
	//	return `
	//local idx = redis.call('LPOS', KEYS[1], ARGV[1])
	//if not idx then
	//  redis.call('RPUSH', KEYS[1], ARGV[1])
	//end
	//redis.call('SET', KEYS[2], "1", "EX", ARGV[2])
	//return 1
	//`
	return "4768da3c5fb4a92fd33c1f144dc69f91d479add2"
}

func TestRetrieveActiveWorkersReturnsOnlyActiveWithCounts(t *testing.T) {
	db, mock := newMockDatabase()
	ctx := context.Background()

	mock.ExpectLRange(WorkersListName, 0, -1).SetVal([]string{"worker-1", "worker-2", "worker-3"})
	mock.ExpectGet(isActiveIdentifier("worker-1")).SetVal("1")
	mock.ExpectSMembers(workerRoomsIdentifier("worker-1")).SetVal([]string{"room-a", "room-b"})
	mock.ExpectGet(isActiveIdentifier("worker-2")).RedisNil()
	mock.ExpectGet(isActiveIdentifier("worker-3")).SetVal("1")
	mock.ExpectSMembers(workerRoomsIdentifier("worker-3")).SetVal([]string{"room-c"})

	workers, err := retrieveActiveWorkers(ctx, db, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("retrieveActiveWorkers() error = %v", err)
	}

	if len(workers) != 2 {
		t.Fatalf("expected 2 active workers, got %d", len(workers))
	}
	if workers["worker-1"] == nil || workers["worker-1"].roomCount != 2 {
		t.Fatalf("expected worker-1 roomCount=2, got %+v", workers["worker-1"])
	}
	if workers["worker-3"] == nil || workers["worker-3"].roomCount != 1 {
		t.Fatalf("expected worker-3 roomCount=1, got %+v", workers["worker-3"])
	}
	if _, ok := workers["worker-2"]; ok {
		t.Fatalf("inactive worker should not be returned")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet redis expectations: %v", err)
	}
}

func TestFindBestWorkerReservesOnLeastBusyWorker(t *testing.T) {
	db, mock := newMockDatabase()

	manager := &RoomManager{
		workers: map[string]*client{
			"worker-1": {IpAddr: "worker-1", roomCount: 3},
			"worker-2": {IpAddr: "worker-2", roomCount: 0},
		},
		db: db,
	}

	mock.ExpectEvalSha(reserveRoomScript(), []string{roomLockIdentifier("room-1")}, "worker-2").SetVal("worker-2")

	worker, err := manager.FindBestWorker(context.Background(), "room-1")
	if err != nil {
		t.Fatalf("FindBestWorker() error = %v", err)
	}
	if worker != "worker-2" {
		t.Fatalf("expected worker-2, got %q", worker)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet redis expectations: %v", err)
	}
}

func TestRegisterRoomIncrementsCountWhenReservationSucceeds(t *testing.T) {
	db, mock := newMockDatabase()

	manager := &RoomManager{
		workers: map[string]*client{
			"worker-1": {IpAddr: "worker-1", roomCount: 1},
		},
		db: db,
	}

	mock.ExpectEvalSha(reserveRoomScript(), []string{roomLockIdentifier("room-1")}, "worker-1").SetVal("worker-1")
	mock.ExpectSAdd(workerRoomsIdentifier("worker-1"), "room-1").SetVal(1)

	worker, err := manager.RegisterRoom(context.Background(), "room-1", "worker-1")
	if err != nil {
		t.Fatalf("RegisterRoom() error = %v", err)
	}
	if worker != "worker-1" {
		t.Fatalf("expected worker-1, got %q", worker)
	}
	if manager.workers["worker-1"].roomCount != 2 {
		t.Fatalf("expected roomCount=2, got %d", manager.workers["worker-1"].roomCount)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet redis expectations: %v", err)
	}
}

func TestRegisterRoomDoesNotIncrementWhenReservedElsewhere(t *testing.T) {
	db, mock := newMockDatabase()

	manager := &RoomManager{
		workers: map[string]*client{
			"worker-1": {IpAddr: "worker-1", roomCount: 1},
		},
		db: db,
	}

	mock.ExpectEvalSha(reserveRoomScript(), []string{roomLockIdentifier("room-1")}, "worker-1").SetVal("worker-2")

	worker, err := manager.RegisterRoom(context.Background(), "room-1", "worker-1")
	if err != nil {
		t.Fatalf("RegisterRoom() error = %v", err)
	}
	if worker != "worker-2" {
		t.Fatalf("expected worker-2, got %q", worker)
	}
	if manager.workers["worker-1"].roomCount != 1 {
		t.Fatalf("expected roomCount unchanged at 1, got %d", manager.workers["worker-1"].roomCount)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet redis expectations: %v", err)
	}
}

func TestReleaseRoomRemovesAndDecrementsCount(t *testing.T) {
	db, mock := newMockDatabase()

	manager := &RoomManager{
		workers: map[string]*client{
			"worker-1": {IpAddr: "worker-1", roomCount: 1},
		},
		db:     db,
		logger: newTestLogger(),
	}

	mock.ExpectSRem(workerRoomsIdentifier("worker-1"), "room-1").SetVal(1)

	if err := manager.ReleaseRoom(context.Background(), "room-1", "worker-1"); err != nil {
		t.Fatalf("ReleaseRoom() error = %v", err)
	}
	if manager.workers["worker-1"].roomCount != 0 {
		t.Fatalf("expected roomCount=0, got %d", manager.workers["worker-1"].roomCount)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet redis expectations: %v", err)
	}
}

func TestSetWorkerActiveUsesConfiguredExpirySeconds(t *testing.T) {
	db, mock := newMockDatabase()
	ctx := context.Background()

	mock.ExpectEvalSha(
		setWorkerActiveScript(),
		[]string{WorkersListName, isActiveIdentifier("worker-1")},
		"worker-1",
		"30",
	).SetVal(1)

	if err := db.SetWorkerActive(ctx, "worker-1", 30*time.Second); err != nil {
		t.Fatalf("SetWorkerActive() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet redis expectations: %v", err)
	}
}
