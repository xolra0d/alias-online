# main

```go
import "github.com/xolra0d/alias-online/services/main"
```

## Index

- [func IPRateLimiter\(limit int, window time.Duration, cleanupEvery int, logger \*slog.Logger\) middleware.Middleware](<#IPRateLimiter>)
- [func NewAuthClient\(authUrl string, logger \*slog.Logger\) \(pbAuth.AuthServiceClient, func\(\) error, error\)](<#NewAuthClient>)
- [func NewRoomManagerClient\(roomManagerUrl string, logger \*slog.Logger\) \(pbRoomManager.RoomManagerServiceClient, func\(\) error, error\)](<#NewRoomManagerClient>)
- [func NewVocabManagerClient\(vocabManagerUrl string, logger \*slog.Logger\) \(pbVocabManager.VocabManagerServiceClient, func\(\) error, error\)](<#NewVocabManagerClient>)
- [func RunServer\(mux \*http.ServeMux, csrf, cors middleware.Middleware, logger \*slog.Logger, runningAddr string, shutdownTimeout time.Duration\)](<#RunServer>)
- [type Handles](<#Handles>)
  - [func NewHTTPHandles\(authClient pbAuth.AuthServiceClient, vocabManagerClient pbVocabManager.VocabManagerServiceClient, roomManagerClient pbRoomManager.RoomManagerServiceClient, logger \*slog.Logger, addAccountTimeout, findAccountTimeout, JWTCookieTimeout time.Duration, JWTCookiePath string, JWTCookieSecure bool, JWTCookieHTTPOnly bool, JWTCookieDomain string\) \*Handles](<#NewHTTPHandles>)
  - [func \(h \*Handles\) AvailableVocabs\(w http.ResponseWriter, r \*http.Request\)](<#Handles.AvailableVocabs>)
  - [func \(h \*Handles\) Healthy\(w http.ResponseWriter, \_ \*http.Request\)](<#Handles.Healthy>)
  - [func \(h \*Handles\) Login\(w http.ResponseWriter, r \*http.Request\)](<#Handles.Login>)
  - [func \(h \*Handles\) Play\(w http.ResponseWriter, r \*http.Request\)](<#Handles.Play>)
  - [func \(h \*Handles\) Register\(w http.ResponseWriter, r \*http.Request\)](<#Handles.Register>)
- [type Secrets](<#Secrets>)
  - [func NewSecrets\(jwtPublicTokenPath string, logger \*slog.Logger\) \(\*Secrets, error\)](<#NewSecrets>)
  - [func \(s \*Secrets\) CheckJwt\(tokenString string\) \(string, error\)](<#Secrets.CheckJwt>)
- [type ServerConfig](<#ServerConfig>)
  - [func LoadServerConfig\(\) \*ServerConfig](<#LoadServerConfig>)


<a name="IPRateLimiter"></a>
## func IPRateLimiter

```go
func IPRateLimiter(limit int, window time.Duration, cleanupEvery int, logger *slog.Logger) middleware.Middleware
```



<a name="NewAuthClient"></a>
## func NewAuthClient

```go
func NewAuthClient(authUrl string, logger *slog.Logger) (pbAuth.AuthServiceClient, func() error, error)
```



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



<a name="RunServer"></a>
## func RunServer

```go
func RunServer(mux *http.ServeMux, csrf, cors middleware.Middleware, logger *slog.Logger, runningAddr string, shutdownTimeout time.Duration)
```



<a name="Handles"></a>
## type Handles



```go
type Handles struct {
    AddAccountTimeout  time.Duration
    FindAccountTimeout time.Duration
    JWTCookieTimeout   time.Duration
    JWTCookiePath      string
    JWTCookieSecure    bool
    JWTCookieHTTPOnly  bool
    JWTCookieDomain    string
    // contains filtered or unexported fields
}
```

<a name="NewHTTPHandles"></a>
### func NewHTTPHandles

```go
func NewHTTPHandles(authClient pbAuth.AuthServiceClient, vocabManagerClient pbVocabManager.VocabManagerServiceClient, roomManagerClient pbRoomManager.RoomManagerServiceClient, logger *slog.Logger, addAccountTimeout, findAccountTimeout, JWTCookieTimeout time.Duration, JWTCookiePath string, JWTCookieSecure bool, JWTCookieHTTPOnly bool, JWTCookieDomain string) *Handles
```



<a name="Handles.AvailableVocabs"></a>
### func \(\*Handles\) AvailableVocabs

```go
func (h *Handles) AvailableVocabs(w http.ResponseWriter, r *http.Request)
```

AvailableVocabs handles /ok requests

<a name="Handles.Healthy"></a>
### func \(\*Handles\) Healthy

```go
func (h *Handles) Healthy(w http.ResponseWriter, _ *http.Request)
```

Healthy handles /ok requests

<a name="Handles.Login"></a>
### func \(\*Handles\) Login

```go
func (h *Handles) Login(w http.ResponseWriter, r *http.Request)
```



<a name="Handles.Play"></a>
### func \(\*Handles\) Play

```go
func (h *Handles) Play(w http.ResponseWriter, r *http.Request)
```



<a name="Handles.Register"></a>
### func \(\*Handles\) Register

```go
func (h *Handles) Register(w http.ResponseWriter, r *http.Request)
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
    AddAccountTimeout  time.Duration // Env name: `ADD_ACCOUNT_TIMEOUT`. Max wait time for saving new account in seconds. Default: 10.
    FindAccountTimeout time.Duration // Env name: `FIND_ACCOUNT_TIMEOUT`. Max wait time for finding account in seconds. Default: 10.

    // HTTP
    VocabsURLGateway      string        // Env name: `VOCABS_URL_GATEWAY`. Gateway for vocab_manager service. Default: none, will exit, if not set.
    AuthURLGateway        string        // Env name: `AUTH_URL_GATEWAY`. Gateway for auth service. Default: none, will exit, if not set.
    RoomManagerURLGateway string        // Env name: `ROOM_MANAGER_URL_GATEWAY`. Gateway for room_manager service. Default: none, will exit, if not set.
    AllowedOrigins        string        // Env name: `ALLOWED_ORIGINS`. Origins to respond to (e.g., http://website.com:12), separated by comma. Default: none, will exit, if not set.
    RunningAddr           string        // Env name: `RUNNING_ADDR`. Address to run web on. Default: `:8080`.
    ShutdownTimeout       time.Duration // Env name: `SHUTDOWN_TIMEOUT`. Time for transport server to shut down in seconds. Default: 10.

    // SECURITY
    JWTCookieTimeout  time.Duration // Env name: `JWT_COOKIE_TIMEOUT`. Time for JWT to expire in seconds. Default: 3600.
    JWTCookiePath     string        // Env name: `JWT_COOKIE_PATH`. Path param for JWT cookie. Default: "/".
    JWTCookieSecure   bool          // Env name: `JWT_COOKIE_SECURE`. Whether cookie should be stored as SECURE. Set false for dev (with http it will not set cookie), true for prod (https). Default: None, will exit, if not set.
    JWTCookieHTTPOnly bool          // Env name: `JWT_COOKIE_HTTP_ONLY`. Cookies cannot be accessed by JavaScript. Default: true.
    JWTCookieDomain   string        // Env name: `JWT_COOKIE_DOMAIN`. Domain for cookie to be accessible. Set "localhost" for dev, and e.g., "xolra0d.com". Default: None, will exit, if not set.
    JwtPublicKeyPath  string        // Env name: `JWT_PUBLIC_KEY_PATH`. Path to the JWT public key file. Default: none, will exit, if not set.
}
```

<a name="LoadServerConfig"></a>
### func LoadServerConfig

```go
func LoadServerConfig() *ServerConfig
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
