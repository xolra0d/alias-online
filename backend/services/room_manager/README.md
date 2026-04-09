# room\_manager

```go
import "github.com/xolra0d/alias-online/services/room_manager"
```

## Index

- [Constants](<#constants>)
- [func RunGrpcServer\(roomManager \*RoomManager, logger \*slog.Logger, runningAddr string, shutdownTimeout time.Duration\)](<#RunGrpcServer>)
- [type Database](<#Database>)
  - [func NewDatabase\(addr, username, password string, db int\) \*Database](<#NewDatabase>)
  - [func \(d \*Database\) AddRoomToWorker\(ctx context.Context, roomId, workerIp string\) error](<#Database.AddRoomToWorker>)
  - [func \(d \*Database\) GetAllWorkers\(ctx context.Context\) \(\[\]string, error\)](<#Database.GetAllWorkers>)
  - [func \(d \*Database\) GetWorkerRoomCount\(ctx context.Context, worker string\) \(int, error\)](<#Database.GetWorkerRoomCount>)
  - [func \(d \*Database\) IsWorkerActive\(ctx context.Context, worker string\) \(bool, error\)](<#Database.IsWorkerActive>)
  - [func \(d \*Database\) ProlongRoom\(ctx context.Context, roomId, workerId string, exp time.Duration\) error](<#Database.ProlongRoom>)
  - [func \(d \*Database\) ReleaseRoom\(ctx context.Context, roomId, workerIp string\) error](<#Database.ReleaseRoom>)
  - [func \(d \*Database\) SetWorkerActive\(ctx context.Context, worker string, exp time.Duration\) error](<#Database.SetWorkerActive>)
  - [func \(d \*Database\) TryReserveRoom\(ctx context.Context, roomId, workerId string\) \(string, error\)](<#Database.TryReserveRoom>)
- [type RoomManager](<#RoomManager>)
  - [func NewManager\(database \*Database, logger \*slog.Logger, PollInterval time.Duration, WorkerExpiry time.Duration, RetrieveActiveWorkersTimeout time.Duration\) \*RoomManager](<#NewManager>)
  - [func \(m \*RoomManager\) FindBestWorker\(ctx context.Context, roomId string\) \(string, error\)](<#RoomManager.FindBestWorker>)
  - [func \(m \*RoomManager\) FindMostFreeWorker\(\) string](<#RoomManager.FindMostFreeWorker>)
  - [func \(m \*RoomManager\) ProlongRoom\(ctx context.Context, roomId, workerIp string\) error](<#RoomManager.ProlongRoom>)
  - [func \(m \*RoomManager\) RegisterRoom\(ctx context.Context, roomId, workerIp string\) \(string, error\)](<#RoomManager.RegisterRoom>)
  - [func \(m \*RoomManager\) ReleaseRoom\(ctx context.Context, roomId, workerIp string\) error](<#RoomManager.ReleaseRoom>)
  - [func \(m \*RoomManager\) ScanForNewWorkers\(\)](<#RoomManager.ScanForNewWorkers>)
  - [func \(m \*RoomManager\) SetWorkerActive\(ctx context.Context, worker string\) error](<#RoomManager.SetWorkerActive>)
  - [func \(m \*RoomManager\) StopScanForNewWorkers\(\)](<#RoomManager.StopScanForNewWorkers>)
- [type ServerConfig](<#ServerConfig>)
  - [func LoadServerConfig\(\) \*ServerConfig](<#LoadServerConfig>)


## Constants

<a name="WorkersListName"></a>

```go
const (
    WorkersListName = "workers" // name for list in redis
)
```

<a name="RunGrpcServer"></a>
## func RunGrpcServer

```go
func RunGrpcServer(roomManager *RoomManager, logger *slog.Logger, runningAddr string, shutdownTimeout time.Duration)
```



<a name="Database"></a>
## type Database



```go
type Database struct {
    // contains filtered or unexported fields
}
```

<a name="NewDatabase"></a>
### func NewDatabase

```go
func NewDatabase(addr, username, password string, db int) *Database
```

NewDatabase creates a new redis client.

<a name="Database.AddRoomToWorker"></a>
### func \(\*Database\) AddRoomToWorker

```go
func (d *Database) AddRoomToWorker(ctx context.Context, roomId, workerIp string) error
```

AddRoomToWorker registers room under workerIp pool of rooms.

<a name="Database.GetAllWorkers"></a>
### func \(\*Database\) GetAllWorkers

```go
func (d *Database) GetAllWorkers(ctx context.Context) ([]string, error)
```

GetAllWorkers returns all ever active workers.

<a name="Database.GetWorkerRoomCount"></a>
### func \(\*Database\) GetWorkerRoomCount

```go
func (d *Database) GetWorkerRoomCount(ctx context.Context, worker string) (int, error)
```

GetWorkerRoomCount returns rooms worker currently holds.

<a name="Database.IsWorkerActive"></a>
### func \(\*Database\) IsWorkerActive

```go
func (d *Database) IsWorkerActive(ctx context.Context, worker string) (bool, error)
```

IsWorkerActive checks if specific worker is currently active.

<a name="Database.ProlongRoom"></a>
### func \(\*Database\) ProlongRoom

```go
func (d *Database) ProlongRoom(ctx context.Context, roomId, workerId string, exp time.Duration) error
```

ProlongRoom prolongs lease of room for worker.

<a name="Database.ReleaseRoom"></a>
### func \(\*Database\) ReleaseRoom

```go
func (d *Database) ReleaseRoom(ctx context.Context, roomId, workerIp string) error
```

ReleaseRoom removes room from workerIp pool of rooms.

<a name="Database.SetWorkerActive"></a>
### func \(\*Database\) SetWorkerActive

```go
func (d *Database) SetWorkerActive(ctx context.Context, worker string, exp time.Duration) error
```

SetWorkerActive sets worker active with timeout.

<a name="Database.TryReserveRoom"></a>
### func \(\*Database\) TryReserveRoom

```go
func (d *Database) TryReserveRoom(ctx context.Context, roomId, workerId string) (string, error)
```

TryReserveRoom tries to atomically reserve room. If succeeded, returns \`workerId\`. If failed, because other worker already reserved this room, returns their worker id.

<a name="RoomManager"></a>
## type RoomManager



```go
type RoomManager struct {
    PollInterval                 time.Duration
    WorkerExpiry                 time.Duration
    RetrieveActiveWorkersTimeout time.Duration
    // contains filtered or unexported fields
}
```

<a name="NewManager"></a>
### func NewManager

```go
func NewManager(database *Database, logger *slog.Logger, PollInterval time.Duration, WorkerExpiry time.Duration, RetrieveActiveWorkersTimeout time.Duration) *RoomManager
```

NewManager creates new room manager.

<a name="RoomManager.FindBestWorker"></a>
### func \(\*RoomManager\) FindBestWorker

```go
func (m *RoomManager) FindBestWorker(ctx context.Context, roomId string) (string, error)
```

FindBestWorker tries to reserve room with the optimal worker from loaded ones.

<a name="RoomManager.FindMostFreeWorker"></a>
### func \(\*RoomManager\) FindMostFreeWorker

```go
func (m *RoomManager) FindMostFreeWorker() string
```

FindMostFreeWorker returns worker with the least rooms loaded.

<a name="RoomManager.ProlongRoom"></a>
### func \(\*RoomManager\) ProlongRoom

```go
func (m *RoomManager) ProlongRoom(ctx context.Context, roomId, workerIp string) error
```

ProlongRoom prolongs lease of room for worker.

<a name="RoomManager.RegisterRoom"></a>
### func \(\*RoomManager\) RegisterRoom

```go
func (m *RoomManager) RegisterRoom(ctx context.Context, roomId, workerIp string) (string, error)
```

RegisterRoom reserves room, adds it to rooms of worker and returns worker. If room is reserved by other worker \- returns that worker.

<a name="RoomManager.ReleaseRoom"></a>
### func \(\*RoomManager\) ReleaseRoom

```go
func (m *RoomManager) ReleaseRoom(ctx context.Context, roomId, workerIp string) error
```

ReleaseRoom removes room from workerIp pool of rooms.

<a name="RoomManager.ScanForNewWorkers"></a>
### func \(\*RoomManager\) ScanForNewWorkers

```go
func (m *RoomManager) ScanForNewWorkers()
```

ScanForNewWorkers retrieves active workers, removes loaded unactive ones and refreshes timeouts.

<a name="RoomManager.SetWorkerActive"></a>
### func \(\*RoomManager\) SetWorkerActive

```go
func (m *RoomManager) SetWorkerActive(ctx context.Context, worker string) error
```

SetWorkerActive sets worker active with timeout.

<a name="RoomManager.StopScanForNewWorkers"></a>
### func \(\*RoomManager\) StopScanForNewWorkers

```go
func (m *RoomManager) StopScanForNewWorkers()
```

StopScanForNewWorkers stops scan loop.

<a name="ServerConfig"></a>
## type ServerConfig



```go
type ServerConfig struct {
    // HTTP
    RunningAddr     string        // Env name: `RUNNING_ADDR`. Address to run web on. Default: `:8060`.
    ShutdownTimeout time.Duration // Env name: `SHUTDOWN_TIMEOUT`. Time for transport server to shut down in seconds. Default: 10.

    // Workers
    PollInterval                 time.Duration // Env name: `WORKERS_POLL_INTERVAL`. Wait time between searches for new workers in seconds. Default: 10.
    WorkerExpiry                 time.Duration // Env name: `WORKER_EXPIRY`. Wait time before worker expires in seconds. Default: 30.
    RetrieveActiveWorkersTimeout time.Duration

    // DATABASE
    RedisAddr     string // Env name: `REDIS_ADDR`. Redis server address. Default: `localhost:6379`.
    RedisUsername string // Env name: `REDIS_USERNAME`. Redis auth username. Default: "".
    RedisPassword string // Env name: `REDIS_PASSWORD`. Redis auth password. Default: ".
    RedisDB       int    // Env name: `REDIS_DB`. Redis database index. Default: 0.
}
```

<a name="LoadServerConfig"></a>
### func LoadServerConfig

```go
func LoadServerConfig() *ServerConfig
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
