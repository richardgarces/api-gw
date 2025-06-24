package routes

import (
	"api-gw/internal/balancer"
	"api-gw/internal/models"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Route struct {
	Path   string `bson:"path"`
	Target string `bson:"target"`
}

type RouteManager struct {
	collection *mongo.Collection
	svcColl    *mongo.Collection
	balancers  map[string]*balancer.RoundRobin
	healths    map[string]*balancer.HealthChecker
}

func NewRouteManager(uri string) (*RouteManager, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	db := client.Database("apigw")
	return &RouteManager{
		collection: db.Collection("routes"),
		svcColl:    db.Collection("services"),
		balancers:  make(map[string]*balancer.RoundRobin),
		healths:    make(map[string]*balancer.HealthChecker),
	}, nil
}

func (rm *RouteManager) Resolve(path string) (string, error) {
	var route Route
	err := rm.collection.FindOne(context.TODO(), Route{Path: path}).Decode(&route)
	if err != nil {
		return "", errors.New("ruta no encontrada 2")
	}
	return route.Target, nil
}

// Devuelve la ruta, el servicio y el target
func (rm *RouteManager) ResolveWithService(path string) (*models.Route, *models.Service, string, error) {
	start := time.Now()
	log.Printf("[INFO] Intentando resolver ruta para path: %s", path)

	var route models.Route
	err := rm.collection.FindOne(context.TODO(), bson.M{"path": path}).Decode(&route)
	if err != nil {
		log.Printf("[WARN] Ruta no encontrada para path: %s (duración: %v)", path, time.Since(start))
		return nil, nil, "", errors.New("ruta no encontrada 1")
	}

	//route.Path = "/simple"

	log.Printf("[INFO] Ruta encontrada: %s -> service_id: %s", path, route.ServiceID.Hex())

	var service models.Service
	err = rm.svcColl.FindOne(context.TODO(), bson.M{"_id": route.ServiceID}).Decode(&service)
	if err != nil {
		log.Printf("[ERROR] Servicio no encontrado para route %s (service_id: %v) (duración: %v)", path, route.ServiceID.Hex(), time.Since(start))
		return &route, nil, "", errors.New("servicio no encontrado")
	}
	log.Printf("[INFO] Servicio encontrado: %s (targets: %v)", service.Name, service.Targets)

	if len(service.Targets) == 0 {
		log.Printf("[ERROR] No hay targets configurados para el servicio %s (duración: %v)", service.Name, time.Since(start))
		return &route, &service, "", errors.New("servicio sin targets")
	}
	// Inicializa balanceador y health checker si no existen
	if rm.balancers[service.ID.Hex()] == nil {
		rm.balancers[service.ID.Hex()] = balancer.NewRoundRobin(service.Targets)
		rm.healths[service.ID.Hex()] = balancer.NewHealthChecker(service.Targets, 10*time.Second)
	}
	healthyTargets := rm.healths[service.ID.Hex()].GetHealthyTargets()
	if len(healthyTargets) == 0 {
		log.Printf("[ERROR] No hay targets saludables para el servicio %s (duración: %v)", service.Name, time.Since(start))
		return &route, &service, "", errors.New("no hay targets saludables")
	}
	rr := balancer.NewRoundRobin(healthyTargets)
	target := rr.Next()
	log.Printf("[INFO] Ruta resuelta: %s -> servicio: %s, target: %s (duración: %v)", path, service.Name, target, time.Since(start))
	return &route, &service, target, nil
}

// Métodos para agregar/editar rutas pueden añadirse aquí
