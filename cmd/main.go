package main

import (
	"api-gw/internal/admin"
	"api-gw/internal/cache"
	"api-gw/internal/config"
	"api-gw/internal/db"
	"api-gw/internal/limiter"
	"api-gw/internal/monitor"
	"api-gw/internal/proxy"
	"api-gw/internal/routes"
	"api-gw/internal/security"
	"fmt"
	"log"
	"net/http"
	"os"
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

	routeManager, err := routes.NewRouteManager("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/admin/services", admin.ServicesHandler)
	mux.HandleFunc("/admin/routes", admin.RoutesHandler)
	mux.HandleFunc("/admin/openapi", admin.OpenAPIHandler)
	mux.Handle("/metrics", monitor.PrometheusHandler())
	mux.Handle("/", monitor.PrometheusMiddleware(
		security.AuthMiddleware(
			limiter.LimitMiddleware(
				cache.CacheMiddleware(
					monitor.LogMiddleware(
						proxy.NewReverseProxy(routeManager),
					),
				),
			),
		),
	))
	mux.Handle("/admin/", http.StripPrefix("/admin/", http.FileServer(http.Dir("webadmin"))))

	log.Printf("API Gateway escuchando en :%s", cfg.HTTPPort)
	log.Fatal(http.ListenAndServe(":"+cfg.HTTPPort, mux))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "usuario1",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("miclaveultrasecreta"))
	fmt.Println(tokenString)
}
