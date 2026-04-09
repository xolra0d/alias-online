# room\_worker

```go
import "github.com/xolra0d/services/room_worker"
```

## Index

- [func NewRoomManagerClient\(roomManagerUrl string, logger \*slog.Logger\) \(pbRoomManager.RoomManagerServiceClient, func\(\) error, error\)](<#NewRoomManagerClient>)
- [func NewVocabManagerClient\(vocabManagerUrl string, logger \*slog.Logger\) \(pbVocabManager.VocabManagerServiceClient, func\(\) error, error\)](<#NewVocabManagerClient>)
- [func RunHttpClient\(handles \*Handles, secrets \*Secrets, logger \*slog.Logger, runningAddr string, shutdownTimeout time.Duration, shouldStop, done chan struct\{\}\)](<#RunHttpClient>)
- [type ClientMessage](<#ClientMessage>)
- [type ClientMessageType](<#ClientMessageType>)
- [type GameState](<#GameState>)
- [type Handles](<#Handles>)
  - [func NewHandles\(secrets \*Secrets, logger \*slog.Logger, rooms \*Rooms\) \*Handles](<#NewHandles>)
  - [func \(h \*Handles\) Healthy\(w http.ResponseWriter, \_ \*http.Request\)](<#Handles.Healthy>)
  - [func \(h \*Handles\) InitWS\(w http.ResponseWriter, r \*http.Request\)](<#Handles.InitWS>)
- [type Player](<#Player>)
  - [func NewPlayer\(id, name string, wordsTried, wordsGuessed int\) \*Player](<#NewPlayer>)
- [type Postgres](<#Postgres>)
  - [func NewPostgres\(postgresUrl string, logger \*slog.Logger\) \(\*Postgres, error\)](<#NewPostgres>)
  - [func \(p \*Postgres\) LoadRoom\(ctx context.Context, roomId string, getVocab func\(ctx context.Context, s string\) \(Vocabulary, error\)\) \(\*Room, error\)](<#Postgres.LoadRoom>)
  - [func \(p \*Postgres\) SaveRoom\(ctx context.Context, r \*Room\) \(err error\)](<#Postgres.SaveRoom>)
- [type PrepareState](<#PrepareState>)
  - [func NewPrepareState\(\) \*PrepareState](<#NewPrepareState>)
  - [func \(s \*PrepareState\) SetErrored\(\)](<#PrepareState.SetErrored>)
  - [func \(s \*PrepareState\) SetOperational\(\)](<#PrepareState.SetOperational>)
  - [func \(s \*PrepareState\) WaitUntilOperational\(\) error](<#PrepareState.WaitUntilOperational>)
- [type Room](<#Room>)
  - [func NewPreparingRoom\(\) \*Room](<#NewPreparingRoom>)
  - [func NewRoom\(id string, admin string, cfg \*RoomConfig, players map\[string\]\*Player, turnOrder \[\]string, currentPlayer int, currentWordIndex int, gameState GameState, logger \*slog.Logger\) \*Room](<#NewRoom>)
  - [func \(r \*Room\) CurrentPlayer\(\) string](<#Room.CurrentPlayer>)
  - [func \(r \*Room\) CurrentWord\(\) string](<#Room.CurrentWord>)
  - [func \(r \*Room\) IncCurrentPlayer\(\)](<#Room.IncCurrentPlayer>)
  - [func \(r \*Room\) Ingest\(msg \*ClientMessage\)](<#Room.Ingest>)
  - [func \(r \*Room\) Join\(player \*Player\)](<#Room.Join>)
  - [func \(r \*Room\) Leave\(player string\)](<#Room.Leave>)
  - [func \(r \*Room\) NextWord\(\) string](<#Room.NextWord>)
  - [func \(r \*Room\) ReportUpdate\(\)](<#Room.ReportUpdate>)
  - [func \(r \*Room\) Run\(onEmpty func\(room \*Room\)\)](<#Room.Run>)
  - [func \(r \*Room\) RunReader\(ctx context.Context, cancel context.CancelFunc, c \*websocket.Conn, player \*Player, maxMessagesPerSecond int\)](<#Room.RunReader>)
  - [func \(r \*Room\) RunWriter\(ctx context.Context, cancel context.CancelFunc, c \*websocket.Conn, player \*Player, wsWriteTimeout, pingTimeout time.Duration\)](<#Room.RunWriter>)
  - [func \(r \*Room\) SetErrored\(\)](<#Room.SetErrored>)
  - [func \(r \*Room\) SetOperational\(\)](<#Room.SetOperational>)
  - [func \(r \*Room\) ToMap\(\) map\[string\]any](<#Room.ToMap>)
  - [func \(r \*Room\) UpdateStateFromRoom\(newRoom \*Room\)](<#Room.UpdateStateFromRoom>)
  - [func \(r \*Room\) UpdateStateFromRoomConfig\(roomId, name, admin string, cfg \*RoomConfig, logger \*slog.Logger\)](<#Room.UpdateStateFromRoomConfig>)
  - [func \(r \*Room\) WaitUntilOperational\(\) error](<#Room.WaitUntilOperational>)
  - [func \(r \*Room\) WordGuessed\(guesser string\)](<#Room.WordGuessed>)
- [type RoomConfig](<#RoomConfig>)
- [type Rooms](<#Rooms>)
  - [func NewRooms\(postgres \*Postgres, logger \*slog.Logger, roomManagerClient pbRoomManager.RoomManagerServiceClient, vocabManagerClient pbVocabManager.VocabManagerServiceClient, runningAddr string, WsOriginPatterns \[\]string, LoadRoomTimeout time.Duration, SaveRoomTimeout time.Duration, WsWriteTimeout time.Duration, WsPingTimeout time.Duration, MaxMessagesPerSecond int, maxClockValue int, loadVocabTimeout time.Duration\) \*Rooms](<#NewRooms>)
  - [func \(rooms \*Rooms\) GetVocab\(ctx context.Context, name string\) \(Vocabulary, error\)](<#Rooms.GetVocab>)
  - [func \(rooms \*Rooms\) ReportLoadedRooms\(\) \[\]string](<#Rooms.ReportLoadedRooms>)
  - [func \(rooms \*Rooms\) RunPinger\(logger \*slog.Logger, PollInterval time.Duration, runningAddr string, shouldStop, done chan struct\{\}\)](<#Rooms.RunPinger>)
  - [func \(rooms \*Rooms\) RunWS\(w http.ResponseWriter, r \*http.Request, roomId, username, name string\) error](<#Rooms.RunWS>)
  - [func \(rooms \*Rooms\) UpdateToWebsocketsAndRedirect\(w http.ResponseWriter, r \*http.Request, otherWorker string\) error](<#Rooms.UpdateToWebsocketsAndRedirect>)
- [type Secrets](<#Secrets>)
  - [func NewSecrets\(jwtPublicTokenPath string, logger \*slog.Logger\) \(\*Secrets, error\)](<#NewSecrets>)
  - [func \(s \*Secrets\) CheckJwt\(tokenString string\) \(string, error\)](<#Secrets.CheckJwt>)
- [type ServerConfig](<#ServerConfig>)
  - [func LoadServerConfig\(\) \*ServerConfig](<#LoadServerConfig>)
- [type ServerMessage](<#ServerMessage>)
- [type ServerMessageType](<#ServerMessageType>)
- [type Vocabulary](<#Vocabulary>)


<a name="NewRoomManagerClient"></a>
## func NewRoomManagerClient

```go
func NewRoomManagerClient(roomManagerUrl string, logger *slog.Logger) (pbRoomManager.RoomManagerServiceClient, func() error, error)
```



<a name="NewVocabManagerClient"></a>
## func NewVocabManagerClient

```go
func NewVocabManagerClient(vocabManagerUrl string, logger *slog.Logger) (pbVocabManager.VocabManagerServiceClient, func() error, error)
```



<a name="RunHttpClient"></a>
## func RunHttpClient

```go
func RunHttpClient(handles *Handles, secrets *Secrets, logger *slog.Logger, runningAddr string, shutdownTimeout time.Duration, shouldStop, done chan struct{})
```



<a name="ClientMessage"></a>
## type ClientMessage



```go
type ClientMessage struct {
    UserId  string            `json:"user_id"`
    MsgType ClientMessageType `json:"type"`
    MsgData map[string]any    `json:"data"`
}
```

<a name="ClientMessageType"></a>
## type ClientMessageType



```go
type ClientMessageType int
```

<a name="GetState"></a>

```go
const (
    GetState ClientMessageType = iota
    StartRound
    GetWord
    TryGuess
    FinishGame
    GetNewWord
    CreateRoom
    LoadRoom
)
```

<a name="GameState"></a>
## type GameState



```go
type GameState int
```

<a name="RoundOver"></a>

```go
const (
    RoundOver GameState = iota
    Explaining
    Finished
)
```

<a name="Handles"></a>
## type Handles



```go
type Handles struct {
    // contains filtered or unexported fields
}
```

<a name="NewHandles"></a>
### func NewHandles

```go
func NewHandles(secrets *Secrets, logger *slog.Logger, rooms *Rooms) *Handles
```



<a name="Handles.Healthy"></a>
### func \(\*Handles\) Healthy

```go
func (h *Handles) Healthy(w http.ResponseWriter, _ *http.Request)
```

Healthy handles /ok requests

<a name="Handles.InitWS"></a>
### func \(\*Handles\) InitWS

```go
func (h *Handles) InitWS(w http.ResponseWriter, r *http.Request)
```

InitWS validates user credentials and tries to update HTTP to Websocket connection.

<a name="Player"></a>
## type Player



```go
type Player struct {
    Id     string      `json:"id"`
    Name   string      `json:"name"`
    ToSend chan []byte `json:"-"`

    Ready        bool `json:"ready"`
    WordsTried   int  `json:"words_tried"`
    WordsGuessed int  `json:"words_guessed"`
}
```

<a name="NewPlayer"></a>
### func NewPlayer

```go
func NewPlayer(id, name string, wordsTried, wordsGuessed int) *Player
```



<a name="Postgres"></a>
## type Postgres



```go
type Postgres struct {
    // contains filtered or unexported fields
}
```

<a name="NewPostgres"></a>
### func NewPostgres

```go
func NewPostgres(postgresUrl string, logger *slog.Logger) (*Postgres, error)
```



<a name="Postgres.LoadRoom"></a>
### func \(\*Postgres\) LoadRoom

```go
func (p *Postgres) LoadRoom(ctx context.Context, roomId string, getVocab func(ctx context.Context, s string) (Vocabulary, error)) (*Room, error)
```

LoadRoom tries to load \`Room\` with its users and turn order. Sets all users to not ready.

<a name="Postgres.SaveRoom"></a>
### func \(\*Postgres\) SaveRoom

```go
func (p *Postgres) SaveRoom(ctx context.Context, r *Room) (err error)
```

SaveRoom tries to save \`Room\` with its users and turn order.

<a name="PrepareState"></a>
## type PrepareState



```go
type PrepareState struct {
    // contains filtered or unexported fields
}
```

<a name="NewPrepareState"></a>
### func NewPrepareState

```go
func NewPrepareState() *PrepareState
```



<a name="PrepareState.SetErrored"></a>
### func \(\*PrepareState\) SetErrored

```go
func (s *PrepareState) SetErrored()
```



<a name="PrepareState.SetOperational"></a>
### func \(\*PrepareState\) SetOperational

```go
func (s *PrepareState) SetOperational()
```



<a name="PrepareState.WaitUntilOperational"></a>
### func \(\*PrepareState\) WaitUntilOperational

```go
func (s *PrepareState) WaitUntilOperational() error
```



<a name="Room"></a>
## type Room



```go
type Room struct {
    Id     string      `json:"id"`
    Admin  string      `json:"admin"`
    Config *RoomConfig `json:"config"`

    Players map[string]*Player `json:"players"`

    CurrentWordIndex int

    TurnOrder []string // circular queue

    State GameState `json:"game_state"`

    RemainingTime int `json:"remaining_time"`
    // contains filtered or unexported fields
}
```

<a name="NewPreparingRoom"></a>
### func NewPreparingRoom

```go
func NewPreparingRoom() *Room
```



<a name="NewRoom"></a>
### func NewRoom

```go
func NewRoom(id string, admin string, cfg *RoomConfig, players map[string]*Player, turnOrder []string, currentPlayer int, currentWordIndex int, gameState GameState, logger *slog.Logger) *Room
```



<a name="Room.CurrentPlayer"></a>
### func \(\*Room\) CurrentPlayer

```go
func (r *Room) CurrentPlayer() string
```



<a name="Room.CurrentWord"></a>
### func \(\*Room\) CurrentWord

```go
func (r *Room) CurrentWord() string
```



<a name="Room.IncCurrentPlayer"></a>
### func \(\*Room\) IncCurrentPlayer

```go
func (r *Room) IncCurrentPlayer()
```



<a name="Room.Ingest"></a>
### func \(\*Room\) Ingest

```go
func (r *Room) Ingest(msg *ClientMessage)
```



<a name="Room.Join"></a>
### func \(\*Room\) Join

```go
func (r *Room) Join(player *Player)
```



<a name="Room.Leave"></a>
### func \(\*Room\) Leave

```go
func (r *Room) Leave(player string)
```



<a name="Room.NextWord"></a>
### func \(\*Room\) NextWord

```go
func (r *Room) NextWord() string
```



<a name="Room.ReportUpdate"></a>
### func \(\*Room\) ReportUpdate

```go
func (r *Room) ReportUpdate()
```



<a name="Room.Run"></a>
### func \(\*Room\) Run

```go
func (r *Room) Run(onEmpty func(room *Room))
```



<a name="Room.RunReader"></a>
### func \(\*Room\) RunReader

```go
func (r *Room) RunReader(ctx context.Context, cancel context.CancelFunc, c *websocket.Conn, player *Player, maxMessagesPerSecond int)
```



<a name="Room.RunWriter"></a>
### func \(\*Room\) RunWriter

```go
func (r *Room) RunWriter(ctx context.Context, cancel context.CancelFunc, c *websocket.Conn, player *Player, wsWriteTimeout, pingTimeout time.Duration)
```



<a name="Room.SetErrored"></a>
### func \(\*Room\) SetErrored

```go
func (r *Room) SetErrored()
```



<a name="Room.SetOperational"></a>
### func \(\*Room\) SetOperational

```go
func (r *Room) SetOperational()
```



<a name="Room.ToMap"></a>
### func \(\*Room\) ToMap

```go
func (r *Room) ToMap() map[string]any
```



<a name="Room.UpdateStateFromRoom"></a>
### func \(\*Room\) UpdateStateFromRoom

```go
func (r *Room) UpdateStateFromRoom(newRoom *Room)
```



<a name="Room.UpdateStateFromRoomConfig"></a>
### func \(\*Room\) UpdateStateFromRoomConfig

```go
func (r *Room) UpdateStateFromRoomConfig(roomId, name, admin string, cfg *RoomConfig, logger *slog.Logger)
```



<a name="Room.WaitUntilOperational"></a>
### func \(\*Room\) WaitUntilOperational

```go
func (r *Room) WaitUntilOperational() error
```



<a name="Room.WordGuessed"></a>
### func \(\*Room\) WordGuessed

```go
func (r *Room) WordGuessed(guesser string)
```



<a name="RoomConfig"></a>
## type RoomConfig

RoomConfig holds specific room\_worker configuration

```go
type RoomConfig struct {
    Seed                 int      `form:"-" json:"-"`
    AllWords             []string `form:"-" json:"-"` // words permutation, unique for every room, dependent on Seed
    Language             string   `form:"language" json:"language"`
    RudeWords            bool     `form:"rude-words" json:"rude-words"`
    AdditionalVocabulary []string `form:"additional-vocabulary" json:"additional-vocabulary"`
    Clock                int      `form:"clock" json:"clock"`
}
```

<a name="Rooms"></a>
## type Rooms



```go
type Rooms struct {
    RunningAddr          string
    WsOriginPatterns     []string
    LoadRoomTimeout      time.Duration
    SaveRoomTimeout      time.Duration
    WsWriteTimeout       time.Duration
    WsPingTimeout        time.Duration
    MaxMessagesPerSecond int
    MaxClockValue        int
    LoadVocabTimeout     time.Duration
    // contains filtered or unexported fields
}
```

<a name="NewRooms"></a>
### func NewRooms

```go
func NewRooms(postgres *Postgres, logger *slog.Logger, roomManagerClient pbRoomManager.RoomManagerServiceClient, vocabManagerClient pbVocabManager.VocabManagerServiceClient, runningAddr string, WsOriginPatterns []string, LoadRoomTimeout time.Duration, SaveRoomTimeout time.Duration, WsWriteTimeout time.Duration, WsPingTimeout time.Duration, MaxMessagesPerSecond int, maxClockValue int, loadVocabTimeout time.Duration) *Rooms
```



<a name="Rooms.GetVocab"></a>
### func \(\*Rooms\) GetVocab

```go
func (rooms *Rooms) GetVocab(ctx context.Context, name string) (Vocabulary, error)
```

GetVocab gets vocab from vocab service.

<a name="Rooms.ReportLoadedRooms"></a>
### func \(\*Rooms\) ReportLoadedRooms

```go
func (rooms *Rooms) ReportLoadedRooms() []string
```



<a name="Rooms.RunPinger"></a>
### func \(\*Rooms\) RunPinger

```go
func (rooms *Rooms) RunPinger(logger *slog.Logger, PollInterval time.Duration, runningAddr string, shouldStop, done chan struct{})
```



<a name="Rooms.RunWS"></a>
### func \(\*Rooms\) RunWS

```go
func (rooms *Rooms) RunWS(w http.ResponseWriter, r *http.Request, roomId, username, name string) error
```

RunWS initiates a WS and configuring a room.

<a name="Rooms.UpdateToWebsocketsAndRedirect"></a>
### func \(\*Rooms\) UpdateToWebsocketsAndRedirect

```go
func (rooms *Rooms) UpdateToWebsocketsAndRedirect(w http.ResponseWriter, r *http.Request, otherWorker string) error
```



<a name="Secrets"></a>
## type Secrets



```go
type Secrets struct {
    // contains filtered or unexported fields
}
```

<a name="NewSecrets"></a>
### func NewSecrets

```go
func NewSecrets(jwtPublicTokenPath string, logger *slog.Logger) (*Secrets, error)
```

NewSecrets creates new secrets manager.

<a name="Secrets.CheckJwt"></a>
### func \(\*Secrets\) CheckJwt

```go
func (s *Secrets) CheckJwt(tokenString string) (string, error)
```

CheckJwt validates jwt to be valid. returns username or error.

<a name="ServerConfig"></a>
## type ServerConfig



```go
type ServerConfig struct {
    // DATABASES
    PostgresUrl              string        // Env name: `POSTGRES_URL`. PostgreSQL connection string. Default: none, will panic, if not set.
    LoadVocabTimeout         time.Duration // Env name: `LOAD_VOCAB_TIMEOUT`. Max wait time for loading vocab in seconds. Default: 10.
    ClosePostgresConnTimeout time.Duration // Env name: `CLOSE_POSTGRES_CONN_TIMEOUT`. Max wait time for closing postgres connection in seconds. Default: 10.

    // HTTP
    AllowedOrigins   string        // Env name: `ALLOWED_ORIGINS`. Origins to respond to (e.g., http://website.com:12), separated by comma. Default: none, will exit, if not set.
    RunningAddr      string        // Env name: `RUNNING_ADDR`. Address to bind HTTP server to. Default: `:8050`.
    WorkerPublicAddr string        // Env name: `WORKER_PUBLIC_ADDR`. Public worker address sent to room_manager and clients. Default: value of `RUNNING_ADDR`.
    ShutdownTimeout  time.Duration // Env name: `SHUTDOWN_TIMEOUT`. Time for transport server to shut down in seconds. Default: 10.

    // WORKERS
    WorkerPollInterval time.Duration // Env name: `WORKER_POLL_INTERVAL`. Wait time for worker pings in seconds. Default: 10.

    WsOriginPatterns     string        // Env name: `WS_ORIGIN_PATTERNS`. Defines allowed origins for ws connections, separated by comma. Default: none, will exit, if not set.
    MaxMessagesPerSecond int           // Env name: `MAX_MESSAGES_PER_SECOND`. Maximum messages sent per second, before connection is closed. Default: 50.
    WsPingTimeout        time.Duration // Env name: `PING_TIMEOUT`. Max wait time for ping request in seconds. Default: 5.
    WsWriteTimeout       time.Duration // Env name: `WS_WRITE_TIMEOUT`. Max wait time for writing response in seconds. Default: 5.
    LoadRoomTimeout      time.Duration // Env name: `LOAD_ROOM_TIMEOUT`. Max wait time for loading room_worker in seconds. Default: 10.
    SaveRoomTimeout      time.Duration // Env name: `SAVE_ROOM_TIMEOUT`. Max wait time for saving room_worker in seconds. Default: 10.
    MaxClockValue        int           // Env name: `MAX_CLOCK_VALUE`. Max clock value used for room state. Default 36000.
    JwtPublicKeyPath     string        // Env name: `JWT_PUBLIC_KEY_PATH`. Path to the JWT public key file. Default: none, will exit, if not set.
    RoomManagerUrl       string        // Env name: `ROOM_MANAGER_URL`. Address of the room manager service. Default: `localhost:8060`.
    VocabManagerUrl      string        // Env name: `VOCAB_MANAGER_URL`. Address of the vocab manager service. Default: `localhost:8070`.
}
```

<a name="LoadServerConfig"></a>
### func LoadServerConfig

```go
func LoadServerConfig() *ServerConfig
```



<a name="ServerMessage"></a>
## type ServerMessage



```go
type ServerMessage struct {
    MsgType ServerMessageType `json:"msg_type"`
    MsgData map[string]any    `json:"msg_data"`
}
```

<a name="ServerMessageType"></a>
## type ServerMessageType



```go
type ServerMessageType int
```

<a name="NewUpdate"></a>

```go
const (
    NewUpdate ServerMessageType = iota
    CurrentState
    YourWord
    WordGuessed
    RightGuess
    WrongGuess
    Redirect // if other worker reserved room
)
```

<a name="Vocabulary"></a>
## type Vocabulary



```go
type Vocabulary struct {
    PrimaryWords []string
    RudeWords    []string
}
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
