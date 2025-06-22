package security

import (
    "net/http"
    "strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        auth := r.Header.Get("Authorization")
        if !strings.HasPrefix(auth, "Bearer ") || len(auth) < 8 {
            http.Error(w, "No autorizado", http.StatusUnauthorized)
            return
        }
        // Aquí deberías validar el JWT
        next.ServeHTTP(w, r)
    })
}