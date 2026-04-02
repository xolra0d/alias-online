package room_manager

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	WorkersListName = "workers"
)

type database struct {
	client *redis.Client
}

func (d *database) GetAllWorkers(ctx context.Context) ([]string, error) {
	return d.client.LRange(ctx, WorkersListName, 0, -1).Result()
}

func (d *database) GetWorkerRoomCount(ctx context.Context, worker string) (int, error) {
	c, err := d.client.HLen(ctx, worker).Result()
	return int(c), err
}

func (d *database) FindMostFreeWorker(ctx context.Context, workers []string) ([]int, error) {
	roomNums := make([]int, 0, len(workers))

	for _, worker := range workers {
		n, err := d.client.HLen(ctx, worker).Result()
		if err != nil {
			return []int{}, err
		}
		roomNums = append(roomNums, int(n))
	}

	return roomNums, nil
}

func (d *database) CheckOrAddRoomID(ctx context.Context, roomId string, workerIP string, workers []string) (string, error) {
	const luaScript = `
local sets = KEYS
local target_set = ARGV[1]
local element = ARGV[2]

for _, set_name in ipairs(sets) do
    if redis.call("SISMEMBER", set_name, element) == 1 then
        return set_name
    end
end

redis.call("SADD", target_set, element)
return target_set
`
	args := []interface{}{workerIP, roomId}

	result, err := d.client.Eval(ctx, luaScript, workers, args...).Result()
	if err != nil {
		return "", fmt.Errorf("lua script error: %w", err)
	}
	return result.(string), nil
}
