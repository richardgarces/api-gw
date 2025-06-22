package plugins

import (
    "net/http"
    "strings"
)

type OAuth2Plugin struct{}

func (o *OAuth2Plugin) Wrap(next http.Handler, config map[string]interface{}) http.Handler {
    expectedToken, _ := config["token"].(string) // Opcional
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        auth := r.Header.Get("Authorization")
        if !strings.HasPrefix(auth, "Bearer ") {
            http.Error(w, "Token OAuth2 requerido", http.StatusUnauthorized)
            return
        }
        token := strings.TrimPrefix(auth, "Bearer ")
        if expectedToken != "" && token != expectedToken {
            http.Error(w, "Token OAuth2 inválido", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func init() {
    Register("oauth2", &OAuth2Plugin{})
}