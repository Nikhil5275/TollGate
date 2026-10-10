package limiter

import (
 "sync"
 "sync/atomic"
 "testing"
)

func TestBurst(t *testing.T) {
 l := New(0.001, 150)
 for i := 0; i < 150; i++ { if ok, _ := l.Allow("client"); !ok { t.Fatalf("request %d rejected", i) } }
 if ok, retry := l.Allow("client"); ok || retry <= 0 { t.Fatal("expected rejection with retry delay") }
 if ok, _ := l.Allow("other"); !ok { t.Fatal("clients must have independent buckets") }
}
func TestConcurrent(t *testing.T) {
 l := New(0.001, 150)
 var accepted atomic.Int64
 var wg sync.WaitGroup
 for i := 0; i < 1000; i++ {
  wg.Add(1)
  go func() { defer wg.Done(); if ok, _ := l.Allow("client"); ok { accepted.Add(1) } }()
 }
 wg.Wait()
 if accepted.Load() != 150 { t.Fatalf("accepted %d, want 150", accepted.Load()) }
}
func BenchmarkAllow(b *testing.B) {
 l := New(100, 150)
 b.RunParallel(func(pb *testing.PB) { for pb.Next() { l.Allow("client") } })
}
