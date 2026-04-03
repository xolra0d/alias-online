package main

import (
	"time"

	"github.com/xolra0d/alias-online/shared/pkg/config"
)

type ServerConfig struct {
	// DATABASES
	PostgresUrl        string        // Env name: `POSTGRES_URL`. PostgreSQL connection string. Default: none, will panic, if not set.
	AddAccountTimeout  time.Duration // Env name: `ADD_ACCOUNT_TIMEOUT`. Max wait time for saving new account in seconds. Default: 10.
	FindAccountTimeout time.Duration // Env name: `FIND_ACCOUNT_TIMEOUT`. Max wait time for finding account in seconds. Default: 10.

	// HTTP
	AllowedOrigins  string        // Env name: `ALLOWED_ORIGINS`. Origins to respond to (e.g., http://website.com:12), separated by comma. Default: none, will exit, if not set.
	RunningAddr     string        // Env name: `RUNNING_ADDR`. Address to run web on. Default: `:8080`.
	ShutdownTimeout time.Duration // Env name: `SHUTDOWN_TIMEOUT`. Time for transport server to shut down in seconds. Default: 10.

	// SECURITY
	JWTCookieTimeout time.Duration // Env name: `JWT_COOKIE_TIMEOUT`. Time for JWT to expire in seconds. Default: 3600.
	//JWTCookiePath         string        // Env name: `JWT_COOKIE_PATH`. Path param for JWT cookie. Default: "/".
	//JWTCookieSecure       bool          // Env name: `JWT_COOKIE_SECURE`. Whether cookie should be stored as SECURE. Set false for dev (with http it will not set cookie), true for prod (https). Default: None, will exit, if not set.
	//JWTCookieHTTPOnly     bool          // Env name: `JWT_COOKIE_HTTP_ONLY`. Cookies cannot be accessed by JavaScript. Default: true.
	//JWTCookieDomain       string        // Env name: `JWT_COOKIE_DOMAIN`. Domain for cookie to be accessible. Set "localhost" for dev, and e.g., "xolra0d.com". Default: None, will exit, if not set.
}

func LoadServerConfig() *ServerConfig {
	postgresUrl := config.GetEnvOrExit("POSTGRES_URL")
	addAccountTimeout := config.StringToSeconds("ADD_ACCOUNT_TIMEOUT", config.GetEnvOrFallback("ADD_ACCOUNT_TIMEOUT", "10"))
	findAccountTimeout := config.StringToSeconds("FIND_ACCOUNT_TIMEOUT", config.GetEnvOrFallback("FIND_ACCOUNT_TIMEOUT", "10"))

	allowedOrigins := config.GetEnvOrExit("ALLOWED_ORIGINS")
	runningAddr := config.GetEnvOrFallback("RUNNING_ADDR", ":8080")
	shutdownTimeout := config.StringToSeconds("SHUTDOWN_TIMEOUT", config.GetEnvOrFallback("SHUTDOWN_TIMEOUT", "10"))

	jwtCookieTimeout := config.StringToSeconds("JWT_COOKIE_TIMEOUT", config.GetEnvOrFallback("JWT_COOKIE_TIMEOUT", "3600"))
	//jwtCookiePath := config.GetEnvOrFallback("JWT_COOKIE_PATH", "/")
	//jwtCookieSecure := config.StringToBool("JWT_COOKIE_SECURE", config.GetEnvOrExit("JWT_COOKIE_SECURE"))
	//jwtCookieHTTPOnly := config.StringToBool("JWT_COOKIE_HTTP_ONLY", config.GetEnvOrFallback("JWT_COOKIE_HTTP_ONLY", "true"))
	//jwtCookieDomain := config.GetEnvOrExit("JWT_COOKIE_DOMAIN")

	return &ServerConfig{
		PostgresUrl:        postgresUrl,
		AddAccountTimeout:  addAccountTimeout,
		FindAccountTimeout: findAccountTimeout,

		AllowedOrigins:  allowedOrigins,
		RunningAddr:     runningAddr,
		ShutdownTimeout: shutdownTimeout,

		JWTCookieTimeout: jwtCookieTimeout,
		//JWTCookiePath:         jwtCookiePath,
		//JWTCookieSecure:       jwtCookieSecure,
		//JWTCookieHTTPOnly:     jwtCookieHTTPOnly,
		//JWTCookieDomain:       jwtCookieDomain,
	}
}
