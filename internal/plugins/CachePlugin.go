package plugins

import (
    "net/http"
    "sync"
    "time"
)

type cacheEntry struct {
    response []byte
    expires  time.Time
}

type CachePlugin struct {
    mu    sync.Mutex
    cache map[string]cacheEntry
}

func (p *CachePlugin) Wrap(next http.Handler, config map[string]interface{}) http.Handler {
    ttl := 10 * time.Second
    if v, ok := config["ttl"].(float64); ok && v > 0 {
        ttl = time.Duration(v) * time.Second
    }
    if p.cache == nil {
        p.cache = make(map[string]cacheEntry)
    }
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        key := r.URL.String()
        p.mu.Lock()
        entry, found := p.cache[key]
        p.mu.Unlock()
        if found && entry.expires.After(time.Now()) {
            w.Write(entry.response)
            return
        }
        rw := &responseWriter{ResponseWriter: w}
        next.ServeHTTP(rw, r)
        p.mu.Lock()
        p.cache[key] = cacheEntry{response: rw.body, expires: time.Now().Add(ttl)}
        p.mu.Unlock()
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

func init() {
    Register("cache", &CachePlugin{})
}