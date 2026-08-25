package transport

import (
	"sync"
	"time"
)

// windowData holds window data for each user.
type windowData struct {
	start    time.Time
	count    int
	lastSeen time.Time
}

// RateLimiter limits access to specific resource through RateLimiter.Allow func.
type RateLimiter struct {
	limit          int
	window         time.Duration
	users          map[string]*windowData
	cleanupEvery   int
	cleanupCounter int
	mu             sync.Mutex
}

// NewRateLimiter creates new rate limiter. `cleanupEvery` removes old entries after 1 time per `cleanupEvery` requests.
func NewRateLimiter(limit int, window time.Duration, cleanupEvery int) *RateLimiter {
	return &RateLimiter{
		limit:        limit,
		window:       window,
		cleanupEvery: cleanupEvery,
		users:        map[string]*windowData{},
	}
}

// Allow checks if identifier is allowed to access resource.
func (l *RateLimiter) Allow(id string) bool {
	now := time.Now()

	if l.limit <= 0 {
		return true
	}
	if id == "" {
		id = "unknown"
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	l.cleanupCounter++
	if l.cleanupCounter >= l.cleanupEvery {
		for k, w := range l.users {
			if now.Sub(w.lastSeen) > 2*l.window { // 2*window, because we want to remember those, who spam.
				delete(l.users, k)
			}
		}
	}

	w, ok := l.users[id]
	if !ok {
		l.users[id] = &windowData{
			start:    now,
			count:    1,
			lastSeen: now,
		}
		return true
	}

	if now.Sub(w.lastSeen) >= l.window {
		w.start = now
		w.count = 1
		w.lastSeen = now
		return true
	}

	if w.count >= l.limit {
		w.lastSeen = now
		return false
	}

	w.count++
	w.lastSeen = now
	return true
}
