# api

```go
import "github.com/xolra0d/alias-online/shared/pkg/api"
```

## Index

- [func WriteJSON\(w http.ResponseWriter, status int, data map\[string\]any\) error](<#WriteJSON>)


<a name="WriteJSON"></a>
## func WriteJSON

```go
func WriteJSON(w http.ResponseWriter, status int, data map[string]any) error
```

WriteJSON writes json response \`data\` to \`w\` with \`status\` status code. Returns \`error\` if either marshaling or writing returns an error.

# config

```go
import "github.com/xolra0d/alias-online/shared/pkg/config"
```

## Index

- [func GetEnvOrExit\(name string\) string](<#GetEnvOrExit>)
- [func GetEnvOrFallback\(name, fallback string\) string](<#GetEnvOrFallback>)
- [func StringToBool\(name, val string\) bool](<#StringToBool>)
- [func StringToSeconds\(name, val string\) time.Duration](<#StringToSeconds>)
- [func StringToUInt\(name, val string\) uint64](<#StringToUInt>)


<a name="GetEnvOrExit"></a>
## func GetEnvOrExit

```go
func GetEnvOrExit(name string) string
```

GetEnvOrExit tries to get \`name\` env var. If not set \- exits.

<a name="GetEnvOrFallback"></a>
## func GetEnvOrFallback

```go
func GetEnvOrFallback(name, fallback string) string
```

GetEnvOrFallback tries to get \`name\` env var and return it. If not set \- returns \`fallback\`.

<a name="StringToBool"></a>
## func StringToBool

```go
func StringToBool(name, val string) bool
```

StringToBool parses \`val\` into bool. Uses \`name\` for panic.

<a name="StringToSeconds"></a>
## func StringToSeconds

```go
func StringToSeconds(name, val string) time.Duration
```

StringToSeconds parses \`val\` \(unsigned num\) into num seconds. Uses \`name\` for panic.

<a name="StringToUInt"></a>
## func StringToUInt

```go
func StringToUInt(name, val string) uint64
```

StringToUInt parses \`val\` \(unsigned num\) into uint64. Uses \`name\` for panic.

# logger

```go
import "github.com/xolra0d/alias-online/shared/pkg/logger"
```

## Index

- [type Handler](<#Handler>)
  - [func NewHandler\(opts \*slog.HandlerOptions\) \*Handler](<#NewHandler>)
  - [func \(h \*Handler\) Enabled\(ctx context.Context, level slog.Level\) bool](<#Handler.Enabled>)
  - [func \(h \*Handler\) Handle\(ctx context.Context, r slog.Record\) error](<#Handler.Handle>)
  - [func \(h \*Handler\) WithAttrs\(attrs \[\]slog.Attr\) slog.Handler](<#Handler.WithAttrs>)
  - [func \(h \*Handler\) WithGroup\(name string\) slog.Handler](<#Handler.WithGroup>)


<a name="Handler"></a>
## type Handler



```go
type Handler struct {
    // contains filtered or unexported fields
}
```

<a name="NewHandler"></a>
### func NewHandler

```go
func NewHandler(opts *slog.HandlerOptions) *Handler
```



<a name="Handler.Enabled"></a>
### func \(\*Handler\) Enabled

```go
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool
```



<a name="Handler.Handle"></a>
### func \(\*Handler\) Handle

```go
func (h *Handler) Handle(ctx context.Context, r slog.Record) error
```



<a name="Handler.WithAttrs"></a>
### func \(\*Handler\) WithAttrs

```go
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler
```



<a name="Handler.WithGroup"></a>
### func \(\*Handler\) WithGroup

```go
func (h *Handler) WithGroup(name string) slog.Handler
```



# middleware

```go
import "github.com/xolra0d/alias-online/shared/pkg/middleware"
```

## Index

- [Constants](<#constants>)
- [func Chain\(h http.Handler, m ...Middleware\) http.Handler](<#Chain>)
- [func LoggingStreamInterceptor\(logger \*slog.Logger\) grpc.StreamServerInterceptor](<#LoggingStreamInterceptor>)
- [func LoggingUnaryInterceptor\(logger \*slog.Logger\) grpc.UnaryServerInterceptor](<#LoggingUnaryInterceptor>)
- [type Middleware](<#Middleware>)
  - [func AuthJWT\(logger \*slog.Logger, validate func\(tokenString string\) \(username string, err error\)\) Middleware](<#AuthJWT>)
  - [func Logging\(logger \*slog.Logger\) Middleware](<#Logging>)
  - [func NewCSRF\(allowedOrigins \[\]string\) Middleware](<#NewCSRF>)
  - [func NewCors\(allowedOrigins, allowedMethods, allowedHeaders \[\]string, allowCredentials bool\) Middleware](<#NewCors>)
  - [func RequestRateLimiter\(l \*RateLimiter, getId func\(r \*http.Request\) string, logger \*slog.Logger\) Middleware](<#RequestRateLimiter>)
- [type RateLimiter](<#RateLimiter>)
  - [func NewRateLimiter\(limit int, window time.Duration, cleanupEvery int\) \*RateLimiter](<#NewRateLimiter>)
  - [func \(l \*RateLimiter\) Allow\(id string\) bool](<#RateLimiter.Allow>)


## Constants

<a name="LoginCookieName"></a>

```go
const (
    LoginCookieName = "login_token"
    LoginContextKey = "login"
)
```

<a name="Chain"></a>
## func Chain

```go
func Chain(h http.Handler, m ...Middleware) http.Handler
```

Chain chains multiple m middlewares before h handler

<a name="LoggingStreamInterceptor"></a>
## func LoggingStreamInterceptor

```go
func LoggingStreamInterceptor(logger *slog.Logger) grpc.StreamServerInterceptor
```



<a name="LoggingUnaryInterceptor"></a>
## func LoggingUnaryInterceptor

```go
func LoggingUnaryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor
```



<a name="Middleware"></a>
## type Middleware



```go
type Middleware func(http.Handler) http.Handler
```

<a name="AuthJWT"></a>
### func AuthJWT

```go
func AuthJWT(logger *slog.Logger, validate func(tokenString string) (username string, err error)) Middleware
```

AuthJWT reads \`Authorization\` cookie and checks if it is valid.

<a name="Logging"></a>
### func Logging

```go
func Logging(logger *slog.Logger) Middleware
```

Logging Logs http each request.

<a name="NewCSRF"></a>
### func NewCSRF

```go
func NewCSRF(allowedOrigins []string) Middleware
```



<a name="NewCors"></a>
### func NewCors

```go
func NewCors(allowedOrigins, allowedMethods, allowedHeaders []string, allowCredentials bool) Middleware
```



<a name="RequestRateLimiter"></a>
### func RequestRateLimiter

```go
func RequestRateLimiter(l *RateLimiter, getId func(r *http.Request) string, logger *slog.Logger) Middleware
```

RequestRateLimiter limits resource usage based on ???.

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

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
