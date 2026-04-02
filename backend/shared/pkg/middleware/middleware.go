package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

// Chain chains multiple m middlewares before h handler
func Chain(h http.Handler, m ...Middleware) http.Handler {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

//// Auth reads `User-Id` and `User-Secret` headers and checks if user exists.
//func Auth(logger *slog.Logger) Middleware {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			id := r.Header.Get("User-Id")
//			secret := r.Header.Get("User-Secret")
//			if id == "" || secret == "" {
//				err := api.WriteJSON(w, http.StatusBadRequest, P{"err": "missing credentials"})
//				if err != nil {
//					logger.Error("could not write response", "err", err)
//				}
//				return
//			}
//			parsedId, err := uuid.Parse(id)
//			if err != nil {
//				err := api.WriteJSON(w, http.StatusBadRequest, P{"err": "invalid credentials"})
//				if err != nil {
//					h.logger.Error("could not write response", "err", err)
//				}
//				return
//			}
//			ok := h.postgres.ValidateUser(r.Context(), database.UserCredentials{Id: parsedId, Secret: secret})
//			if !ok {
//				err := api.WriteJSON(w, http.StatusUnauthorized, P{"err": "invalid credentials"})
//				if err != nil {
//					logger.Error("could not write response", "err", err)
//				}
//				return
//			}
//			next.ServeHTTP(w, r)
//		})
//	}
//}

//// UserIdRateLimiter limits resource usage based on `User-Id` header.
//func (h *Handles) UserIdRateLimiter(l *RateLimiter) Middleware {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			id := r.Header.Get("User-Id")
//			if id == "" {
//				err := WriteJSON(w, http.StatusBadRequest, P{"err": "missing credentials"})
//				if err != nil {
//					h.logger.Error("could not write response", "err", err)
//				}
//				return
//			}
//			if !l.Allow(id) {
//				return
//			}
//			next.ServeHTTP(w, r)
//		})
//	}
//}

//// ReadUserIP tries to get user real IP.
//func ReadUserIP(r *http.Request) string {
//	IPAddress := r.Header.Get("X-Real-Ip")
//	if IPAddress == "" {
//		IPAddress = r.Header.Get("X-Forwarded-For")
//	}
//	if IPAddress == "" {
//		IPAddress = r.RemoteAddr
//	}
//	return IPAddress
//}

//// IpRateLimiter limits resource usage based on IP from ReadUserIP.
//func (h *Handles) IpRateLimiter(l *RateLimiter) Middleware {
//	const op = "transport.ipRateLimiter"
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			ip := ReadUserIP(r)
//			if ip == "" {
//				err := WriteJSON(w, http.StatusBadRequest, P{"err": "missing ip address"})
//				if err != nil {
//					h.logger.Error(op, "could not write response", "err", err)
//				}
//				return
//			}
//			if !l.Allow(ip) {
//				return
//			}
//			next.ServeHTTP(w, r)
//		})
//	}
//}

// Logging Logs each request.
func Logging(logger *slog.Logger) Middleware {
	const op = "transport.Logging"
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("got request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
		})
	}
}
