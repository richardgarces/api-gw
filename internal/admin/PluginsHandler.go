package admin

import (
    "api-gw/internal/plugins"
    "encoding/json"
    "net/http"
)

func PluginsHandler(w http.ResponseWriter, r *http.Request) {
    var names []string
    for name := range plugins.Registry() {
        names = append(names, name)
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(names)
}