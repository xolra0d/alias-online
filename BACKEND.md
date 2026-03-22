# backend

```go
import "github.com/xolra0d/alias-online"
```

## Index

- [func Chain\(h http.Handler, m ...Middleware\) http.Handler](<#Chain>)
- [func InitPool\(postgresUrl string\) \(\*pgxpool.Pool, error\)](<#InitPool>)
- [func ReadUserIP\(r \*http.Request\) string](<#ReadUserIP>)
- [func WriteJSON\(w http.ResponseWriter, status int, data P\) error](<#WriteJSON>)
- [type BaseLogger](<#BaseLogger>)
  - [func NewBaseLogger\(queueLen int\) \*BaseLogger](<#NewBaseLogger>)
  - [func \(l \*BaseLogger\) EndLogging\(\)](<#BaseLogger.EndLogging>)
  - [func \(l \*BaseLogger\) StartLogging\(\)](<#BaseLogger.StartLogging>)
  - [func \(l \*BaseLogger\) WithPrefix\(prefix string\) \*PrefixLogger](<#BaseLogger.WithPrefix>)
- [type ClientMessage](<#ClientMessage>)
- [type ClientMessageType](<#ClientMessageType>)
- [type GameState](<#GameState>)
- [type Handles](<#Handles>)
  - [func \(h \*Handles\) Auth\(\) Middleware](<#Handles.Auth>)
  - [func \(h \*Handles\) AvailableLanguages\(w http.ResponseWriter, \_ \*http.Request\)](<#Handles.AvailableLanguages>)
  - [func \(h \*Handles\) CreateRoom\(w http.ResponseWriter, r \*http.Request\)](<#Handles.CreateRoom>)
  - [func \(h \*Handles\) CreateUser\(w http.ResponseWriter, r \*http.Request\)](<#Handles.CreateUser>)
  - [func \(h \*Handles\) Healthy\(w http.ResponseWriter, \_ \*http.Request\)](<#Handles.Healthy>)
  - [func \(h \*Handles\) InitWS\(w http.ResponseWriter, r \*http.Request\)](<#Handles.InitWS>)
  - [func \(h \*Handles\) IpRateLimiter\(l \*RateLimiter\) Middleware](<#Handles.IpRateLimiter>)
  - [func \(h \*Handles\) UserIdRateLimiter\(l \*RateLimiter\) Middleware](<#Handles.UserIdRateLimiter>)
- [type Middleware](<#Middleware>)
  - [func Logging\(logger \*PrefixLogger\) Middleware](<#Logging>)
- [type P](<#P>)
- [type Player](<#Player>)
- [type Postgres](<#Postgres>)
  - [func \(p \*Postgres\) AddRoom\(ctx context.Context, adminId uuid.UUID, cfg RoomConfig\) \(string, error\)](<#Postgres.AddRoom>)
  - [func \(p \*Postgres\) CreateUser\(ctx context.Context\) \(UserCredentials, error\)](<#Postgres.CreateUser>)
  - [func \(p \*Postgres\) LoadRoom\(ctx context.Context, roomId string, vocabs \*Vocabularies\) \(\*Room, error\)](<#Postgres.LoadRoom>)
  - [func \(p \*Postgres\) LoadVocabs\(ctx context.Context\) \(map\[string\]\*Vocabulary, error\)](<#Postgres.LoadVocabs>)
  - [func \(p \*Postgres\) UpdateRoomState\(ctx context.Context, r \*Room\) error](<#Postgres.UpdateRoomState>)
  - [func \(p \*Postgres\) ValidateUser\(ctx context.Context, credentials UserCredentials\) bool](<#Postgres.ValidateUser>)
- [type PrefixLogger](<#PrefixLogger>)
  - [func \(l \*PrefixLogger\) CopyWithPrefix\(prefix string\) \*PrefixLogger](<#PrefixLogger.CopyWithPrefix>)
  - [func \(l \*PrefixLogger\) Debug\(msg string, args ...any\)](<#PrefixLogger.Debug>)
  - [func \(l \*PrefixLogger\) Error\(msg string, args ...any\)](<#PrefixLogger.Error>)
  - [func \(l \*PrefixLogger\) Info\(msg string, args ...any\)](<#PrefixLogger.Info>)
  - [func \(l \*PrefixLogger\) Warn\(msg string, args ...any\)](<#PrefixLogger.Warn>)
- [type RateLimiter](<#RateLimiter>)
  - [func NewRateLimiter\(limit int, window time.Duration, cleanupEvery int\) \*RateLimiter](<#NewRateLimiter>)
  - [func \(l \*RateLimiter\) Allow\(id string\) bool](<#RateLimiter.Allow>)
- [type Room](<#Room>)
  - [func \(r \*Room\) CurrentWord\(vocabs \*Vocabularies\) string](<#Room.CurrentWord>)
  - [func \(r \*Room\) IncCurrentPlayer\(\)](<#Room.IncCurrentPlayer>)
  - [func \(r \*Room\) NextWord\(vocabs \*Vocabularies\) string](<#Room.NextWord>)
  - [func \(r \*Room\) ReportUpdate\(\)](<#Room.ReportUpdate>)
  - [func \(r \*Room\) Run\(postgres \*Postgres, vocabs \*Vocabularies, rooms \*Rooms\)](<#Room.Run>)
  - [func \(r \*Room\) SaveState\(ctx context.Context, postgres \*Postgres\) error](<#Room.SaveState>)
  - [func \(r \*Room\) ToMap\(\) map\[string\]any](<#Room.ToMap>)
  - [func \(r \*Room\) WordGuessed\(guesser uuid.UUID\)](<#Room.WordGuessed>)
- [type RoomConfig](<#RoomConfig>)
- [type Rooms](<#Rooms>)
  - [func \(rooms \*Rooms\) ServeWS\(w http.ResponseWriter, r \*http.Request, userId uuid.UUID, name, roomId string, postgres \*Postgres, vocabs \*Vocabularies\) error](<#Rooms.ServeWS>)
- [type Secrets](<#Secrets>)
  - [func \(s \*Secrets\) GenerateName\(\) string](<#Secrets.GenerateName>)
  - [func \(s \*Secrets\) GenerateRoomId\(\) string](<#Secrets.GenerateRoomId>)
  - [func \(s \*Secrets\) GenerateSecretBase32\(\) string](<#Secrets.GenerateSecretBase32>)
  - [func \(s \*Secrets\) VerifyPassword\(secret, hash string\) bool](<#Secrets.VerifyPassword>)
- [type ServerConfig](<#ServerConfig>)
- [type ServerMessage](<#ServerMessage>)
- [type ServerMessageType](<#ServerMessageType>)
- [type UserCredentials](<#UserCredentials>)
- [type Vocabularies](<#Vocabularies>)
- [type Vocabulary](<#Vocabulary>)


<a name="Chain"></a>
## func Chain

```go
func Chain(h http.Handler, m ...Middleware) http.Handler
```

Chains multiple m middlewares before h handler

<a name="InitPool"></a>
## func InitPool

```go
func InitPool(postgresUrl string) (*pgxpool.Pool, error)
```

InitPool initializes PostgreSQL pool from \`postgresUrl\`.

<a name="ReadUserIP"></a>
## func ReadUserIP

```go
func ReadUserIP(r *http.Request) string
```

ReadUserIP tries to get user real IP.

<a name="WriteJSON"></a>
## func WriteJSON

```go
func WriteJSON(w http.ResponseWriter, status int, data P) error
```

WriteJSON writes json response \`data\` to \`w\` with \`status\` status code. Returns \`error\` if either marshaling or writing returns an error.

<a name="BaseLogger"></a>
## type BaseLogger

BaseLogger helps to create async loggers, using WithPrefix.

```go
type BaseLogger struct {
    // contains filtered or unexported fields
}
```

<a name="NewBaseLogger"></a>
### func NewBaseLogger

```go
func NewBaseLogger(queueLen int) *BaseLogger
```

NewBaseLogger creates new empty logger.

<a name="BaseLogger.EndLogging"></a>
### func \(\*BaseLogger\) EndLogging

```go
func (l *BaseLogger) EndLogging()
```

EndLogging closes the channel for messages.

<a name="BaseLogger.StartLogging"></a>
### func \(\*BaseLogger\) StartLogging

```go
func (l *BaseLogger) StartLogging()
```

StartLogging starts to retrieve messages received from all PrefixLogger's.

<a name="BaseLogger.WithPrefix"></a>
### func \(\*BaseLogger\) WithPrefix

```go
func (l *BaseLogger) WithPrefix(prefix string) *PrefixLogger
```

WithPrefix creates new async logger.

<a name="ClientMessage"></a>
## type ClientMessage



```go
type ClientMessage struct {
    UserId  uuid.UUID         `json:"user_id"`
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

Handles holds http handles and, realistically, the state of the program.

```go
type Handles struct {
    CreateUserTimeout time.Duration
    CreateRoomTimeout time.Duration
    // contains filtered or unexported fields
}
```

<a name="Handles.Auth"></a>
### func \(\*Handles\) Auth

```go
func (h *Handles) Auth() Middleware
```

Auth reads \`User\-Id\` and \`User\-Secret\` headers and checks if user exists.

<a name="Handles.AvailableLanguages"></a>
### func \(\*Handles\) AvailableLanguages

```go
func (h *Handles) AvailableLanguages(w http.ResponseWriter, _ *http.Request)
```

AvailableLanguages returns loaded vocabs names.

<a name="Handles.CreateRoom"></a>
### func \(\*Handles\) CreateRoom

```go
func (h *Handles) CreateRoom(w http.ResponseWriter, r *http.Request)
```

CreateRoom validates room config from form and inserts it to database, returning roomId.

<a name="Handles.CreateUser"></a>
### func \(\*Handles\) CreateUser

```go
func (h *Handles) CreateUser(w http.ResponseWriter, r *http.Request)
```

CreateUser generates random login, name, secret for user and returns as Credentials.

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

InitWS validates user credentials and tries to update HTTP to Websocket connecton.

<a name="Handles.IpRateLimiter"></a>
### func \(\*Handles\) IpRateLimiter

```go
func (h *Handles) IpRateLimiter(l *RateLimiter) Middleware
```

IpRateLimiter limits resource usage based on IP from ReadUserIP.

<a name="Handles.UserIdRateLimiter"></a>
### func \(\*Handles\) UserIdRateLimiter

```go
func (h *Handles) UserIdRateLimiter(l *RateLimiter) Middleware
```

UserIdRateLimiter limits resource usage based on \`User\-Id\` header.

<a name="Middleware"></a>
## type Middleware



```go
type Middleware func(http.Handler) http.Handler
```

<a name="Logging"></a>
### func Logging

```go
func Logging(logger *PrefixLogger) Middleware
```

Logging Logs each request.

<a name="P"></a>
## type P

P is shortcut for map\[string\]any

```go
type P map[string]any
```

<a name="Player"></a>
## type Player



```go
type Player struct {
    Id   uuid.UUID `json:"id"`
    Name string    `json:"name"`

    Ready        bool `json:"ready"`
    WordsTried   int  `json:"words_tried"`
    WordsGuessed int  `json:"words_guessed"`
    // contains filtered or unexported fields
}
```

<a name="Postgres"></a>
## type Postgres

Postgres helps with postgres\-specific commands.

```go
type Postgres struct {
    // contains filtered or unexported fields
}
```

<a name="Postgres.AddRoom"></a>
### func \(\*Postgres\) AddRoom

```go
func (p *Postgres) AddRoom(ctx context.Context, adminId uuid.UUID, cfg RoomConfig) (string, error)
```

AddRoom generates seeds vocab and saves config to database.

<a name="Postgres.CreateUser"></a>
### func \(\*Postgres\) CreateUser

```go
func (p *Postgres) CreateUser(ctx context.Context) (UserCredentials, error)
```

CreateUser creates random user and inserts into db.

<a name="Postgres.LoadRoom"></a>
### func \(\*Postgres\) LoadRoom

```go
func (p *Postgres) LoadRoom(ctx context.Context, roomId string, vocabs *Vocabularies) (*Room, error)
```

LoadRoom tries to load \`Room\` with its users and turn order. Sets all users to not ready.

<a name="Postgres.LoadVocabs"></a>
### func \(\*Postgres\) LoadVocabs

```go
func (p *Postgres) LoadVocabs(ctx context.Context) (map[string]*Vocabulary, error)
```

LoadVocabs tries to load all vocabs from database.

<a name="Postgres.UpdateRoomState"></a>
### func \(\*Postgres\) UpdateRoomState

```go
func (p *Postgres) UpdateRoomState(ctx context.Context, r *Room) error
```

UpdateRoomState updates room state and state of each player in who were in room.

<a name="Postgres.ValidateUser"></a>
### func \(\*Postgres\) ValidateUser

```go
func (p *Postgres) ValidateUser(ctx context.Context, credentials UserCredentials) bool
```

ValidateUser tries to check if user with \`credentials.Id\` and \`credentials.Secret\` exists.

<a name="PrefixLogger"></a>
## type PrefixLogger

PrefixLogger asyncronously logs logs.

```go
type PrefixLogger struct {
    Prefix string
    // contains filtered or unexported fields
}
```

<a name="PrefixLogger.CopyWithPrefix"></a>
### func \(\*PrefixLogger\) CopyWithPrefix

```go
func (l *PrefixLogger) CopyWithPrefix(prefix string) *PrefixLogger
```

CopyWithPrefix is the same as WithPrefix

<a name="PrefixLogger.Debug"></a>
### func \(\*PrefixLogger\) Debug

```go
func (l *PrefixLogger) Debug(msg string, args ...any)
```

Debug logs data with DEBUG level

<a name="PrefixLogger.Error"></a>
### func \(\*PrefixLogger\) Error

```go
func (l *PrefixLogger) Error(msg string, args ...any)
```

Error logs data with ERROR level

<a name="PrefixLogger.Info"></a>
### func \(\*PrefixLogger\) Info

```go
func (l *PrefixLogger) Info(msg string, args ...any)
```

Info logs data with INFO level

<a name="PrefixLogger.Warn"></a>
### func \(\*PrefixLogger\) Warn

```go
func (l *PrefixLogger) Warn(msg string, args ...any)
```

Warn logs data with WARN level

<a name="RateLimiter"></a>
## type RateLimiter

RateLimiter limits access to specific resource through RateLimiter.Allow func.

```go
type RateLimiter struct {
    // contains filtered or unexported fields
}
```

<a name="NewRateLimiter"></a>
### func NewRateLimiter

```go
func NewRateLimiter(limit int, window time.Duration, cleanupEvery int) *RateLimiter
```

NewRateLimiter creates new rate limiter. \`cleanupEvery\` removes old entries after 1 time per \`cleanupEvery\` requests.

<a name="RateLimiter.Allow"></a>
### func \(\*RateLimiter\) Allow

```go
func (l *RateLimiter) Allow(id string) bool
```

Allow checks if identifier is allowed to access resource.

<a name="Room"></a>
## type Room



```go
type Room struct {
    Id     string      `json:"id"`
    Admin  uuid.UUID   `json:"admin"`
    Config *RoomConfig `json:"config"`

    Players map[uuid.UUID]*Player `json:"players"`

    State GameState `json:"game_state"`

    RemainingTime int `json:"remaining_time"`
    // contains filtered or unexported fields
}
```

<a name="Room.CurrentWord"></a>
### func \(\*Room\) CurrentWord

```go
func (r *Room) CurrentWord(vocabs *Vocabularies) string
```



<a name="Room.IncCurrentPlayer"></a>
### func \(\*Room\) IncCurrentPlayer

```go
func (r *Room) IncCurrentPlayer()
```



<a name="Room.NextWord"></a>
### func \(\*Room\) NextWord

```go
func (r *Room) NextWord(vocabs *Vocabularies) string
```



<a name="Room.ReportUpdate"></a>
### func \(\*Room\) ReportUpdate

```go
func (r *Room) ReportUpdate()
```



<a name="Room.Run"></a>
### func \(\*Room\) Run

```go
func (r *Room) Run(postgres *Postgres, vocabs *Vocabularies, rooms *Rooms)
```



<a name="Room.SaveState"></a>
### func \(\*Room\) SaveState

```go
func (r *Room) SaveState(ctx context.Context, postgres *Postgres) error
```



<a name="Room.ToMap"></a>
### func \(\*Room\) ToMap

```go
func (r *Room) ToMap() map[string]any
```



<a name="Room.WordGuessed"></a>
### func \(\*Room\) WordGuessed

```go
func (r *Room) WordGuessed(guesser uuid.UUID)
```



<a name="RoomConfig"></a>
## type RoomConfig

RoomConfig holds specific room configuration

```go
type RoomConfig struct {
    Language             string   `form:"language" json:"language"`
    RudeWords            bool     `form:"rude-words" json:"rude-words"`
    AdditionalVocabulary []string `form:"additional-vocabulary" json:"additional-vocabulary"`
    Clock                int      `form:"clock" json:"clock"`
    // contains filtered or unexported fields
}
```

<a name="Rooms"></a>
## type Rooms



```go
type Rooms struct {
    MinClock                     int
    MaxClock                     int
    MaxAdditionalVocabularyWords int
    MaxAdditionalWordLength      int

    WSOriginPatterns     []string
    MaxMessagesPerSecond int
    PingTimeout          time.Duration
    WSWriteTimeout       time.Duration
    WSReadTimeout        time.Duration
    LoadRoomTimeout      time.Duration
    SaveRoomTimeout      time.Duration
    // contains filtered or unexported fields
}
```

<a name="Rooms.ServeWS"></a>
### func \(\*Rooms\) ServeWS

```go
func (rooms *Rooms) ServeWS(w http.ResponseWriter, r *http.Request, userId uuid.UUID, name, roomId string, postgres *Postgres, vocabs *Vocabularies) error
```



<a name="Secrets"></a>
## type Secrets



```go
type Secrets struct {
    Argon2idTime    uint32
    Argon2idMemory  uint32
    Argon2idThreads uint8
    Argon2idOutLen  uint32
    // contains filtered or unexported fields
}
```

<a name="Secrets.GenerateName"></a>
### func \(\*Secrets\) GenerateName

```go
func (s *Secrets) GenerateName() string
```

GenerateName creates a new name for account in form \`AdjectiveNoun\(0\-99\)\`.

<a name="Secrets.GenerateRoomId"></a>
### func \(\*Secrets\) GenerateRoomId

```go
func (s *Secrets) GenerateRoomId() string
```

GenerateRoomId creates new 40 bit base32 roomId

<a name="Secrets.GenerateSecretBase32"></a>
### func \(\*Secrets\) GenerateSecretBase32

```go
func (s *Secrets) GenerateSecretBase32() string
```

GenerateSecretBase32 creates secure base32 secret.

<a name="Secrets.VerifyPassword"></a>
### func \(\*Secrets\) VerifyPassword

```go
func (s *Secrets) VerifyPassword(secret, hash string) bool
```

VerifyPassword checks if secret is equal to hash's secret.

<a name="ServerConfig"></a>
## type ServerConfig

ServerConfig holds all runtime configuration loaded from ENV variables.

```go
type ServerConfig struct {
    // APP
    LogMessageMaxQueue int           // Env name: `LOG_MESSAGE_MAX_QUEUE`. Max queue length before client logger needs to wait. Default: 100.
    LoadVocabsTimeout  time.Duration // Env name: `LOAD_VOCABS_TIMEOUT`. Vocabularies load timeout in seconds. Default: 5.

    // DATABASES
    PostgresUrl string // Env name: `POSTGRES_URL`. PostgreSQL connection string. Default: none, will panic, if not set.

    // SECURITY
    Argon2idTime    uint32 // Env name: `ARGON2ID_TIME`. Number of iterations to perform. Default: 2.
    Argon2idMemory  uint32 // Env name: `ARGON2ID_MEMORY`. Amount of memory to use in bytes. Default: 65536.
    Argon2idThreads uint8  // Env name: `ARGON2ID_THREADS`. Degree of parallelism. Default: 1.
    Argon2idOutLen  uint32 // Env name: `ARGON2ID_OUT_LEN`. Desired number of returned bytes. Default: 32.

    // HTTP
    AllowedOrigins           string        // Env name: `ALLOWED_ORIGINS`. Origins to respond to (e.g., http://website.com:12), separetad by comma. Default: none, will exit, if not set.
    RunningAddr              string        // Env name: `RUNNING_ADDR`. Address to run web on. Default: `:8080`.
    ReadTimeout              time.Duration // Env name: `READ_TIMEOUT`. Docs: net/http.Server.ReadTimeout in seconds. Default: 5.
    WriteTimeout             time.Duration // Env name: `WRITE_TIMEOUT`. Docs: net/http.Server.WriteTimeout in seconds. Default: 5.
    IdleTimeout              time.Duration // Env name: `IDLE_TIMEOUT`. Docs: net/http.Server.IdleTimeout in seconds. Default: 30.
    ShutdownTimeout          time.Duration // Env name: `SHUTDOWN_TIMEOUT`. Time for http server to shut down in seconds. Default: 10.
    CreateUserLimitPerWindow int           // Env name: `CREATE_USER_LIMIT_PER_WINDOW`. Number of users allowed to be created for `LimiterWindow` time. Default: 30.
    CreateRoomLimitPerWindow int           // Env name: `CREATE_ROOM_LIMIT_PER_WINDOW`. Number of rooms allowed to be created for `LimiterWindow` time. Default: 30.
    LimiterCleanupEvery      int           // Env name: `LIMITER_CLEANUP_EVERY`. Removes outdated entries after `LimiterCleanupEvery` requests handled. Default: 100.
    LimiterWindow            time.Duration // Env name: `LIMITER_WINDOW`. Defines window for limiters in seconds. Default: 60.
    CreateUserTimeout        time.Duration // Env name: `CREATE_USER_TIMEOUT`. Defines timeout for user creation in seconds. Default: 5.
    CreateRoomTimeout        time.Duration // Env name: `CREATE_ROOM_TIMEOUT`. Defines timeout for room creation in seconds. Default: 5.

    // ROOMS
    MinClock                     int           // Env name: `MIN_CLOCK`. Min number of seconds for clock in round. Default: 1.
    MaxClock                     int           // Env name: `MAX_CLOCK`. Max number of seconds for clock in round. Default: 36000. (10 hours)
    MaxAdditionalVocabularyWords int           // Env name: `MAX_ADDITIONAL_VOCABULARY_WORDS`. Max number of words in additional vocabulary. Default: 1000.
    MaxAdditionalWordLength      int           // Env name: `MAX_ADDITIONAL_WORD_LENGTH`. Max number of runes (UTF-8 chars) in word in additional vocabulary. Default 64.
    LoadRoomTimeout              time.Duration // Env name: `LOAD_ROOM_TIMEOUT`. Max wait time for loading room in seconds. Default: 10.
    SaveRoomTimeout              time.Duration // Env name: `SAVE_ROOM_TIMEOUT`. Max wait time for saving room in seconds. Default: 10.

    // WS
    WSOriginPatterns     string        // Env name: `WS_ORIGIN_PATTERNS`. Defines allowed origins for ws connections, separetad by comma. Default: none, will exit, if not set.
    MaxMessagesPerSecond int           // Env name: `MAX_MESSAGES_PER_SECOND`. Maximum messages sent per second, before connection is closed. Default: 50.
    PingTimeout          time.Duration // Env name: `PING_TIMEOUT`. Max wait time for ping request in seconds. Default: 5.
    WSWriteTimeout       time.Duration // Env name: `WS_WRITE_TIMEOUT`. Max wait time for writing response in seconds. Default: 5.
    WSReadTimeout        time.Duration // Env name: `WS_READ_TIMEOUT`. Max wait time for reading request in seconds. Connection will not be reset after this timeout. Used only to do other things, while no message is present. Default: 5.
}
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
)
```

<a name="UserCredentials"></a>
## type UserCredentials

UserCredentials stores user credentials sent while creating user.

```go
type UserCredentials struct {
    Id     uuid.UUID `json:"id"`
    Name   string    `json:"name"`
    Secret string    `json:"secret"`
}
```

<a name="Vocabularies"></a>
## type Vocabularies



```go
type Vocabularies struct {
    // contains filtered or unexported fields
}
```

<a name="Vocabulary"></a>
## type Vocabulary



```go
type Vocabulary struct {
    PrimaryWords []string
    RudeWords    []string
}
```
