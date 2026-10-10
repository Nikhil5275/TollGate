package limiter

import (
 "sync"
 "time"
)

type bucket struct { tokens float64; updated time.Time }

// Limiter owns synchronized per-client token buckets.
type Limiter struct {
 mu sync.Mutex
 buckets map[string]bucket
 rate, burst float64
}

func New(rate, burst float64) *Limiter {
 if rate <= 0 || burst < 1 { panic("rate must be positive and burst at least one") }
 return &Limiter{buckets: make(map[string]bucket), rate: rate, burst: burst}
}

// Allow consumes one token or returns the delay until one becomes available.
func (l *Limiter) Allow(key string) (bool, time.Duration) {
 l.mu.Lock()
 defer l.mu.Unlock()
 now := time.Now()
 b, exists := l.buckets[key]
 if !exists { b = bucket{l.burst, now} }
 b.tokens += now.Sub(b.updated).Seconds() * l.rate
 if b.tokens > l.burst { b.tokens = l.burst }
 b.updated = now
 if b.tokens >= 1 {
  b.tokens--
  l.buckets[key] = b
  return true, 0
 }
 l.buckets[key] = b
 return false, time.Duration((1-b.tokens)/l.rate*float64(time.Second))
}
