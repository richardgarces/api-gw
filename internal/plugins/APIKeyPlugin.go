package plugins

import (
    "net/http"
)

type APIKeyPlugin struct{}

func (a *APIKeyPlugin) Wrap(next http.Handler, config map[string]interface{}) http.Handler {
    expectedKey, ok := config["key"].(string)
    if !ok || expectedKey == "" {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            http.Error(w, "API Key plugin mal configurado", http.StatusInternalServerError)
        })
    }
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        key := r.Header.Get("X-API-Key")
        if key == "" {
			key = r.URL.Query().Get("api_key")
        }
        if key != expectedKey {
            http.Error(w, "API Key inválida", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func init() {
    Register("apikey", &APIKeyPlugin{})
}