package limiter

import (
	"net/http"
	"sync"
	"time"
)

var clients = make(map[string]int)
var mu sync.Mutex

func LimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		mu.Lock()
		clients[ip]++
		count := clients[ip]
		mu.Unlock()
		if count > 100 {
			http.Error(w, "Demasiadas solicitudes", http.StatusTooManyRequests)
			return
		}
		go func() {
			time.Sleep(time.Minute)
			mu.Lock()
			clients[ip]--
			mu.Unlock()
		}()
		next.ServeHTTP(w, r)
	})
}
