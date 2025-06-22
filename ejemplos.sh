# API Gateway - Ejemplos de uso

## 1. Solicitud a un servicio backend

```sh
curl -k -H "Authorization: Bearer <token>" https://localhost:8080/servicio1
```

## 2. Insertar una nueva ruta (directamente en MongoDB)

```js
db.routes.insertOne({ path: "/nuevo", target: "http://localhost:9003" })
```

## 3. Ejemplo de respuesta cacheada

Haz dos veces la misma petición y la segunda será servida desde caché:

```sh
curl -k -H "Authorization: Bearer <token>" https://localhost:8080/servicio1
```

## 4. Límite de solicitudes

Haz más de 100 peticiones por minuto desde la misma IP y recibirás:

```
HTTP/1.1 429 Too Many Requests
```