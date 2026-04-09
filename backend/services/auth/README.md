# auth

```go
import "github.com/xolra0d/alias-online/services/auth"
```

## Index

- [func RunGrpcServer\(secrets \*Secrets, postgres \*Postgres, logger \*slog.Logger, addAccountTimeout, findAccountTimeout, jwtCookieTimeout time.Duration, runningAddr string, shutdownTimeout time.Duration\)](<#RunGrpcServer>)
- [func ValidateForLogin\(login, password string\) error](<#ValidateForLogin>)
- [func ValidateForRegister\(name, login, password string\) error](<#ValidateForRegister>)
- [func ValidateLogin\(login string\) error](<#ValidateLogin>)
- [func ValidateName\(name string\) error](<#ValidateName>)
- [func ValidatePassword\(password string\) error](<#ValidatePassword>)
- [type Error](<#Error>)
  - [func NewError\(svcError ServiceError, appError error\) \*Error](<#NewError>)
  - [func \(e \*Error\) Error\(\) string](<#Error.Error>)
- [type Postgres](<#Postgres>)
  - [func NewPostgres\(postgresUrl string, logger \*slog.Logger\) \(\*Postgres, error\)](<#NewPostgres>)
  - [func \(p \*Postgres\) AddAccount\(ctx context.Context, name, login, hash string\) \*Error](<#Postgres.AddAccount>)
  - [func \(p \*Postgres\) Close\(\)](<#Postgres.Close>)
  - [func \(p \*Postgres\) FindAccount\(ctx context.Context, login string\) \(string, \*Error\)](<#Postgres.FindAccount>)
- [type Secrets](<#Secrets>)
  - [func NewSecrets\(logger \*slog.Logger, Argon2idTime uint32, Argon2idMemory uint32, Argon2idThreads uint8, Argon2idOutLen uint32, RsaPrivateKeyFilename string\) \(\*Secrets, error\)](<#NewSecrets>)
  - [func \(s \*Secrets\) GenerateSecretBase32\(\) string](<#Secrets.GenerateSecretBase32>)
  - [func \(s \*Secrets\) NewJWT\(login string, exp time.Time\) \(string, error\)](<#Secrets.NewJWT>)
  - [func \(s \*Secrets\) VerifyPassword\(secret, hash string\) \*Error](<#Secrets.VerifyPassword>)
- [type ServerConfig](<#ServerConfig>)
  - [func LoadServerConfig\(\) \*ServerConfig](<#LoadServerConfig>)
- [type ServiceError](<#ServiceError>)


<a name="RunGrpcServer"></a>
## func RunGrpcServer

```go
func RunGrpcServer(secrets *Secrets, postgres *Postgres, logger *slog.Logger, addAccountTimeout, findAccountTimeout, jwtCookieTimeout time.Duration, runningAddr string, shutdownTimeout time.Duration)
```



<a name="ValidateForLogin"></a>
## func ValidateForLogin

```go
func ValidateForLogin(login, password string) error
```

ValidateForLogin checks if login and password are valid for according fields in db.

<a name="ValidateForRegister"></a>
## func ValidateForRegister

```go
func ValidateForRegister(name, login, password string) error
```

ValidateForRegister checks if name, login and password are valid for according fields in db.

<a name="ValidateLogin"></a>
## func ValidateLogin

```go
func ValidateLogin(login string) error
```

ValidateLogin check if login is valid for login field.

<a name="ValidateName"></a>
## func ValidateName

```go
func ValidateName(name string) error
```

ValidateName check if name is valid for name field.

<a name="ValidatePassword"></a>
## func ValidatePassword

```go
func ValidatePassword(password string) error
```

ValidatePassword check if password is valid for password field.

<a name="Error"></a>
## type Error



```go
type Error struct {
    SvcError ServiceError // Outer error returned to user
    AppError error        // Inner error returned by inner services (db)
}
```

<a name="NewError"></a>
### func NewError

```go
func NewError(svcError ServiceError, appError error) *Error
```



<a name="Error.Error"></a>
### func \(\*Error\) Error

```go
func (e *Error) Error() string
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

NewPostgres creates new instance of postgres client.

<a name="Postgres.AddAccount"></a>
### func \(\*Postgres\) AddAccount

```go
func (p *Postgres) AddAccount(ctx context.Context, name, login, hash string) *Error
```

AddAccount tries to insert a new account into postgres. Error: \- ErrUserAlreadyExists \- if this login already exists. \- ErrInternal \- if there was insertion error.

<a name="Postgres.Close"></a>
### func \(\*Postgres\) Close

```go
func (p *Postgres) Close()
```

Close closes postgres pool.

<a name="Postgres.FindAccount"></a>
### func \(\*Postgres\) FindAccount

```go
func (p *Postgres) FindAccount(ctx context.Context, login string) (string, *Error)
```

FindAccount finds account with specific login and returns hashed password. Error: \- ErrUserNotFound \- if login does not exist. \- ErrInternal \- if there was select error.

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
func NewSecrets(logger *slog.Logger, Argon2idTime uint32, Argon2idMemory uint32, Argon2idThreads uint8, Argon2idOutLen uint32, RsaPrivateKeyFilename string) (*Secrets, error)
```

NewSecrets creates new secrets manager.

<a name="Secrets.GenerateSecretBase32"></a>
### func \(\*Secrets\) GenerateSecretBase32

```go
func (s *Secrets) GenerateSecretBase32() string
```

GenerateSecretBase32 creates secure base32 secret.

<a name="Secrets.NewJWT"></a>
### func \(\*Secrets\) NewJWT

```go
func (s *Secrets) NewJWT(login string, exp time.Time) (string, error)
```

NewJWT Issues new JWT.

<a name="Secrets.VerifyPassword"></a>
### func \(\*Secrets\) VerifyPassword

```go
func (s *Secrets) VerifyPassword(secret, hash string) *Error
```

VerifyPassword checks if secret is equal to hash's secret.

<a name="ServerConfig"></a>
## type ServerConfig



```go
type ServerConfig struct {
    // DATABASES
    PostgresUrl        string        // Env name: `POSTGRES_URL`. PostgreSQL connection string. Default: none, will exit, if not set.
    AddAccountTimeout  time.Duration // Env name: `ADD_ACCOUNT_TIMEOUT`. Max wait time for saving new account in seconds. Default: 10.
    FindAccountTimeout time.Duration // Env name: `FIND_ACCOUNT_TIMEOUT`. Max wait time for finding account in seconds. Default: 10.

    // GRPC SERVER
    RunningAddr     string        // Env name: `RUNNING_ADDR`. Address to run web on. Default: `:8090`.
    ShutdownTimeout time.Duration // Env name: `SHUTDOWN_TIMEOUT`. Time for transport server to shut down in seconds. Default: 10.

    // SECURITY
    Argon2idTime          uint32        // Env name: `ARGON2ID_TIME`. Number of iterations to perform. Default: 2.
    Argon2idMemory        uint32        // Env name: `ARGON2ID_MEMORY`. Amount of memory to use in bytes. Default: 65536.
    Argon2idThreads       uint8         // Env name: `ARGON2ID_THREADS`. Degree of parallelism. Default: 1.
    Argon2idOutLen        uint32        // Env name: `ARGON2ID_OUT_LEN`. Desired number of returned bytes. Default: 32.
    JwtPrivateKeyFilename string        // Env name: `JWT_PRIVATE_KEY_FILENAME`. Path to private key used for creating JWT tokens. Default: none, will exit, if not set.
    JWTCookieTimeout      time.Duration // Env name: `JWT_COOKIE_TIMEOUT`. Time for JWT to expire in seconds. Default: 3600.
}
```

<a name="LoadServerConfig"></a>
### func LoadServerConfig

```go
func LoadServerConfig() *ServerConfig
```



<a name="ServiceError"></a>
## type ServiceError



```go
type ServiceError string
```

<a name="ErrInternal"></a>

```go
const (
    ErrInternal          ServiceError = "internal error"
    ErrWrongCredentials  ServiceError = "wrong credentials"
    ErrUserAlreadyExists ServiceError = "user already exists"
    ErrUserNotFound      ServiceError = "user not found"
)
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
