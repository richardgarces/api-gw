# API Gateway

## Ejemplo de uso

```sh
curl -k -H "Authorization: Bearer <token>" https://localhost:8080/servicio1/ruta
```



// Ejemplo para POST /admin/services
{
  "name": "servicio1",
  "targets": ["http://localhost:9001"],
  "plugins": [
    {"type": "logging", "config": {}}
  ]
}


{
  "name": "servicio2",
  "targets": ["http://localhost:9002"],
  "plugins": [
    {
      "type": "jwt",
      "config": {
        "secret": "miclaveultrasecreta"
      }
    }
  ]
}


curl -H "Authorization: Bearer <TOKEN_JWT>" http://localhost:8080/servicio2


{
  "type": "apikey",
  "config": {
    "key": "miclaveapi"
  }
}

{
  "type": "oauth2",
  "config": {
    "token": "token_oauth2_de_prueba"
  }
}

"plugins": [
  {"type": "apikey", "config": {"key": "miclaveapi"}},
  {"type": "jwt", "config": {"secret": "miclaveultrasecreta"}},
  {"type": "oauth2", "config": {"token": "token_oauth2_de_prueba"}}
]

curl -H "X-API-Key: miclaveapi" http://localhost:8080/servicio1

curl "http://localhost:8080/servicio1?api_key=miclaveapi"

curl -H "Authorization: Bearer token_oauth2_de_prueba" http://localhost:8080/servicio1


{
  "type": "ratelimit",
  "config": {
    "limit": 10,
    "by": "ip"
  }
}

"by": "ip" limita por dirección IP.
"by": "apikey" limita por API Key.
"by": "user" limita por el campo sub del JWT.

{
  "type": "request_headers",
  "config": {
    "set": {
      "X-Request-Source": "api-gateway",
      "X-User": "demo"
    }
  }
}

{
  "type": "response_headers",
  "config": {
    "set": {
      "X-Powered-By": "Go-API-Gateway"
    }
  }
}


{
  "type": "nightblock",
  "config": {
    "start": 22,
    "end": 6
  }
}

openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=localhost"



Crear un servicio con plugins:

curl -X POST http://localhost:8080/admin/services \
  -H "Content-Type: application/json" \
  -d '{
    "name":"servicio1",
    "targets":["http://localhost:9001"],
    "plugins":[
      {"type":"jwt","config":{"secret":"miclaveultrasecreta"}},
      {"type":"ratelimit","config":{"limit":5,"by":"ip"}}
    ]
  }'


Crear una ruta:

curl -X POST http://localhost:8080/admin/routes \
  -H "Content-Type: application/json" \
  -d '{
    "path":"/servicio1",
    "service_id":"<ID_DEL_SERVICIO>",
    "plugins":[{"type":"request_headers","config":{"set":{"X-Request-Source":"api-gateway"}}}]
  }'


Acceso a panel web:
http://localhost:8080/admin/

Acceso a OpenAPI:
http://localhost:8080/admin/openapi

Acceso a métricas Prometheus:
http://localhost:8080/metrics


8. Pasos para levantar todo
Levanta MongoDB y el gateway:
Accede al panel web y administra tus servicios/rutas/plugins.
Haz peticiones a través del gateway usando JWT, API Key, OAuth2, etc.
9. Extiende con tus propios plugins
Crea un archivo en plugins, implementa la interfaz y regístralo en init().
Documenta el plugin en /admin/openapi o /admin/plugins.
10. Pruebas automáticas
Ejecuta:

make test
para correr los tests de tus plugins y middlewares.
