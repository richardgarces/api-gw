package admin

import (
    "encoding/json"
    "net/http"
)

var openAPISpec = map[string]interface{}{
    "openapi": "3.0.0",
    "info": map[string]interface{}{
        "title":   "API Gateway Admin",
        "version": "1.0.0",
    },
    "paths": map[string]interface{}{
        "/admin/services": map[string]interface{}{
            "get": map[string]interface{}{
                "summary":     "Listar servicios",
                "responses":   map[string]interface{}{"200": map[string]interface{}{"description": "OK"}},
            },
            "post": map[string]interface{}{
                "summary":     "Crear servicio",
                "requestBody": map[string]interface{}{"required": true},
                "responses":   map[string]interface{}{"201": map[string]interface{}{"description": "Creado"}},
            },
        },
        "/admin/routes": map[string]interface{}{
            "get": map[string]interface{}{
                "summary":     "Listar rutas",
                "responses":   map[string]interface{}{"200": map[string]interface{}{"description": "OK"}},
            },
            "post": map[string]interface{}{
                "summary":     "Crear ruta",
                "requestBody": map[string]interface{}{"required": true},
                "responses":   map[string]interface{}{"201": map[string]interface{}{"description": "Creado"}},
            },
        },
    },
}

func OpenAPIHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(openAPISpec)
}