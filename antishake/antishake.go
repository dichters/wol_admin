// Package antishake provides an in-memory anti-shake lock backed by a TTL cache.
// It is only active when config.EnableAntiShake is true.
package antishake

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"wol-panel/config"
)

const (
	lockTTL = 3 * time.Second
)

// Locker provides anti-shake locking capability backed by a TTL cache.
type Locker struct {
	mu    sync.Mutex
	cache *ttlcache.Cache[string, struct{}]
}

// New creates a new Locker. When anti-shake is disabled a no-op Locker is returned.
func New() *Locker {
	l := &Locker{}

	if !config.Cfg.EnableAntiShake {
		slog.Info("anti-shake disabled by config")
		return l
	}

	l.cache = ttlcache.New[string, struct{}](
		ttlcache.WithTTL[string, struct{}](lockTTL),
		ttlcache.WithDisableTouchOnHit[string, struct{}](),
	)
	go l.cache.Start()

	slog.Info("in-memory anti-shake enabled", "ttl", lockTTL.String())
	return l
}

// TryLock attempts to acquire a lock for the given key.
// Returns true if the lock was acquired, false if already locked.
// Lock auto-expires after lockTTL.
func (l *Locker) TryLock(clientID, operation string) bool {
	if !config.Cfg.EnableAntiShake || l.cache == nil {
		return true // anti-shake disabled, always allow
	}

	key := fmt.Sprintf("wol-panel:lock:%s:%s", clientID, operation)

	l.mu.Lock()
	defer l.mu.Unlock()

	if item := l.cache.Get(key); item != nil {
		return false // still locked
	}

	l.cache.Set(key, struct{}{}, lockTTL)
	return true
}

// Close stops the TTL cache janitor.
func (l *Locker) Close() error {
	if l.cache != nil {
		l.cache.Stop()
	}
	return nil
}
