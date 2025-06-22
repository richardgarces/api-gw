package cache

import (
	"net/http"
	"sync"
	"time"
)

type cacheEntry struct {
	response []byte
	expires  time.Time
}

var (
	cache = make(map[string]cacheEntry)
	mu    sync.Mutex
)

func CacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.String()
		mu.Lock()
		entry, found := cache[key]
		mu.Unlock()
		if found && entry.expires.After(time.Now()) {
			w.Write(entry.response)
			return
		}
		rw := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(rw, r)
		mu.Lock()
		cache[key] = cacheEntry{response: rw.body, expires: time.Now().Add(10 * time.Second)}
		mu.Unlock()
	})
}

type responseWriter struct {
	http.ResponseWriter
	body []byte
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body = append(rw.body, b...)
	return rw.ResponseWriter.Write(b)
}
