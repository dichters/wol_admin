package antishake

import (
	"testing"
	"time"

	"wol-panel/config"
)

func withAntiShake(enabled bool) func() {
	prev := config.Cfg
	config.Cfg = &config.Config{EnableAntiShake: enabled}
	return func() { config.Cfg = prev }
}

func TestTryLockBlocksWithinTTL(t *testing.T) {
	restore := withAntiShake(true)
	defer restore()

	l := New()
	defer l.Close()

	if !l.TryLock("client-a", "wol") {
		t.Fatal("first TryLock should succeed")
	}
	if l.TryLock("client-a", "wol") {
		t.Fatal("second TryLock within TTL should fail")
	}
	if !l.TryLock("client-b", "wol") {
		t.Fatal("different client should not be blocked")
	}
	if !l.TryLock("client-a", "shutdown") {
		t.Fatal("different operation should not be blocked")
	}
}

func TestTryLockExpiresAfterTTL(t *testing.T) {
	restore := withAntiShake(true)
	defer restore()

	l := New()
	defer l.Close()

	if !l.TryLock("client-a", "wol") {
		t.Fatal("first TryLock should succeed")
	}
	time.Sleep(lockTTL + 500*time.Millisecond)
	if !l.TryLock("client-a", "wol") {
		t.Fatal("TryLock should succeed after TTL expiry")
	}
}

func TestTryLockAlwaysAllowsWhenDisabled(t *testing.T) {
	restore := withAntiShake(false)
	defer restore()

	l := New()
	defer l.Close()

	for i := 0; i < 3; i++ {
		if !l.TryLock("client-a", "wol") {
			t.Fatal("TryLock should always succeed when anti-shake is disabled")
		}
	}
}
