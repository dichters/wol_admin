// Package antishake provides a distributed anti-shake lock using Redis.
// When Redis is unavailable it falls back to an in-memory sync.Map lock.
// It is only active when config.EnableAntiShake is true.
package antishake

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"wol_admin/config"
)

const (
	lockTTL = 3 * time.Second
)

// Locker provides anti-shake locking capability.
type Locker struct {
	rdb      *redis.Client
	fallback *sync.Map // in-memory fallback when Redis is down
	useRedis bool
}

// New creates a new Locker. If Redis connection fails, it falls back to in-memory.
func New() *Locker {
	l := &Locker{
		fallback: &sync.Map{},
		useRedis: false,
	}

	if !config.Cfg.EnableAntiShake {
		slog.Info("anti-shake disabled by config")
		return l
	}

	addr := config.Cfg.Redis.IP + ":" + config.Cfg.Redis.Port
	opts := &redis.Options{
		Addr: addr,
	}
	if config.Cfg.Redis.Password != "" {
		opts.Password = config.Cfg.Redis.Password
	}

	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Warn("Redis connection failed, falling back to in-memory anti-shake", "error", err)
		rdb.Close()
		return l
	}

	l.rdb = rdb
	l.useRedis = true
	slog.Info("Redis anti-shake enabled", "addr", addr)
	return l
}

// TryLock attempts to acquire a lock for the given key.
// Returns true if the lock was acquired, false if already locked.
// Lock auto-expires after 3 seconds.
func (l *Locker) TryLock(clientID, operation string) bool {
	if !config.Cfg.EnableAntiShake {
		return true // anti-shake disabled, always allow
	}

	key := fmt.Sprintf("wol_admin:lock:%s:%s", clientID, operation)

	if l.useRedis {
		return l.redisLock(key)
	}
	return l.memoryLock(key)
}

// Close releases Redis connection if active.
func (l *Locker) Close() error {
	if l.rdb != nil {
		return l.rdb.Close()
	}
	return nil
}

func (l *Locker) redisLock(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ok, err := l.rdb.Set(ctx, key, "1", lockTTL).Result()
	if err == redis.Nil {
		return false
	}
	if err != nil {
		slog.Warn("Redis SET failed, falling back to memory lock", "key", key, "error", err)
		return l.memoryLock(key)
	}
	return ok == "1" || ok == "OK"
}

func (l *Locker) memoryLock(key string) bool {
	now := time.Now()
	if v, ok := l.fallback.Load(key); ok {
		if expire, ok2 := v.(time.Time); ok2 && now.Before(expire) {
			return false // still locked
		}
		// expired, remove and allow re-lock
		l.fallback.Delete(key)
	}
	l.fallback.Store(key, now.Add(lockTTL))
	return true
}
