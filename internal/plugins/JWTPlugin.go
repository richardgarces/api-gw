package plugins

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type JWTPlugin struct{}

func (j *JWTPlugin) Wrap(next http.Handler, config map[string]interface{}) http.Handler {
	secret, ok := config["secret"].(string)
	if !ok || secret == "" {
		// Si no hay secreto, rechaza todas las peticiones
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "JWT plugin mal configurado", http.StatusInternalServerError)
		})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Token JWT requerido", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(auth, "Bearer ")
		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Solo soporta HS256 por simplicidad
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Método de firma inválido")
			}
			return []byte(secret), nil
		})
		if err != nil {
			http.Error(w, "Token JWT inválido: "+err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func init() {
	Register("jwt", &JWTPlugin{})
}