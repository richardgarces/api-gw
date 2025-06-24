package main

import (
	"api-gw/internal/admin"
	"api-gw/internal/config"
	"api-gw/internal/db"
	"api-gw/internal/limiter"
	"api-gw/internal/monitor"
	"api-gw/internal/proxy"
	"api-gw/internal/routes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	cfg := config.Load()
	os.Setenv("MONGO_DATABASE", cfg.MongoDatabase)

	client, err := db.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Error conectando a MongoDB: %v", err)
	}
	defer client.Disconnect(nil)

	//routeManager, err := routes.NewRouteManager("mongodb://localhost:27017")
	routeManager, err := routes.NewRouteManager(cfg.MongoURI)

	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// Handler combinado para "/" y proxy general
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Api - Gateway OK"))
			return
		}
		monitor.PrometheusMiddleware(
			limiter.LimitMiddleware(
				monitor.LogMiddleware(
					proxy.NewReverseProxy(routeManager),
				),
			),
		).ServeHTTP(w, r)
	})

	mux.HandleFunc("/admin/services", admin.ServicesHandler)
	mux.HandleFunc("/admin/routes", admin.RoutesHandler)
	mux.HandleFunc("/admin/openapi", admin.OpenAPIHandler)
	mux.Handle("/metrics", monitor.PrometheusHandler())
	mux.Handle("/admin/", http.StripPrefix("/admin/", http.FileServer(http.Dir("webadmin"))))

	// Configuración de JWT para el panel admin
	jwtSecret := os.Getenv("ADMIN_JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "miclaveultrasecreta"
	}
	jwtDurationStr := os.Getenv("ADMIN_JWT_DURATION")
	jwtDuration := time.Hour * 24 * 365 // 1 año por defecto
	if jwtDurationStr != "" {
		if dur, err := strconv.Atoi(jwtDurationStr); err == nil {
			jwtDuration = time.Duration(dur) * time.Second
		}
	}

	// Endpoint para obtener el token JWT para el panel admin
	mux.HandleFunc("/admin/token", func(w http.ResponseWriter, r *http.Request) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": "admin",
			"exp": time.Now().Add(jwtDuration).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(jwtSecret))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	})

	log.Printf("API Gateway escuchando en :%s", cfg.HTTPPort)
	log.Fatal(http.ListenAndServe(":"+cfg.HTTPPort, mux))
}
