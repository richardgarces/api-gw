package plugins

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
)

type RateLimitPlugin struct {
	// Mapa: clave -> contador y timestamp
	clients sync.Map
}

type rateInfo struct {
	count     int
	timestamp time.Time
}

func (r *RateLimitPlugin) Wrap(next http.Handler, config map[string]interface{}) http.Handler {
	limit, ok := config["limit"].(float64)
	if !ok || limit <= 0 {
		limit = 60 // Por defecto 60 req/min
	}
	by, _ := config["by"].(string) // "ip", "apikey", "user"
	window := time.Minute

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		key := ""
		switch by {
		case "apikey":
			key = req.Header.Get("X-API-Key")
			if key == "" {
				key = req.URL.Query().Get("api_key")
			}
		case "user":
			// Busca el sub del JWT si existe
			auth := req.Header.Get("Authorization")
			if len(auth) > 7 && auth[:7] == "Bearer " {
				claims, _ := parseJWTClaims(auth[7:])
				if sub, ok := claims["sub"].(string); ok {
					key = sub
				}
			}
		default:
			key = req.RemoteAddr
		}
		if key == "" {
			key = req.RemoteAddr
		}

		now := time.Now()
		val, _ := r.clients.LoadOrStore(key, &rateInfo{count: 0, timestamp: now})
		info := val.(*rateInfo)

		// Reinicia ventana si pasó el tiempo
		if now.Sub(info.timestamp) > window {
			info.count = 0
			info.timestamp = now
		}
		info.count++
		if info.count > int(limit) {
			w.Header().Set("Retry-After", "60")
			http.Error(w, "Rate limit excedido", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, req)
	})
}

// Utilidad para extraer claims de un JWT (sin validar firma)
func parseJWTClaims(tokenString string) (map[string]interface{}, error) {
	// Solo para extraer claims, no valida firma
	parts := strings.Split(tokenString, ".")
	if len(parts) < 2 {
		return nil, nil
	}
	decoded, err := jwtDecodeSegment(parts[1])
	if err != nil {
		return nil, err
	}
	var claims map[string]interface{}
	err = json.Unmarshal(decoded, &claims)
	return claims, err
}

func jwtDecodeSegment(seg string) ([]byte, error) {
	return jwt.DecodeSegment(seg)
}

func init() {
	Register("ratelimit", &RateLimitPlugin{})
}
