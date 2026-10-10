package main

import (
 "fmt"
 "log"
 "math"
 "net"
 "net/http"
 "net/http/httputil"
 "net/url"
 "os"
 "strings"
 "time"

 "github.com/Nikhil5275/TollGate/internal/limiter"
 "github.com/Nikhil5275/TollGate/internal/ring"
)

func env(key, fallback string) string { if v := os.Getenv(key); v != "" { return v }; return fallback }

func main() {
 self := env("NODE_URL", "http://localhost:8080")
 nodes := strings.Split(env("PEERS", self), ",")
 owners := ring.New(nodes)
 buckets := limiter.New(100, 150)
 proxies := make(map[string]*httputil.ReverseProxy)
 for _, peer := range nodes {
  u, err := url.Parse(peer)
  if err != nil || u.Host == "" { log.Fatalf("invalid peer: %q", peer) }
  proxies[peer] = httputil.NewSingleHostReverseProxy(u)
 }
 mux := http.NewServeMux()
 mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "ok") })
 mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  key := r.Header.Get("X-Client-ID")
  if key == "" { key, _, _ = net.SplitHostPort(r.RemoteAddr) }
  owner := owners.Owner(key)
  if owner != self {
   // Carry the resolved key so the owner does not hash the proxy's address.
   r.Header.Set("X-Client-ID", key)
   proxies[owner].ServeHTTP(w, r)
   return
  }
  w.Header().Set("X-TollGate-Node", self)
  if allowed, retry := buckets.Allow(key); !allowed {
   w.Header().Set("Retry-After", fmt.Sprint(int(math.Ceil(retry.Seconds()))))
   http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
   return
  }
  fmt.Fprintln(w, "request accepted")
 })
 server := &http.Server{Addr: ":"+env("PORT", "8080"), Handler: mux, ReadHeaderTimeout: 5*time.Second, WriteTimeout: 10*time.Second, IdleTimeout: 60*time.Second}
 log.Printf("TollGate listening on %s; node=%s", server.Addr, self)
 log.Fatal(server.ListenAndServe())
}
