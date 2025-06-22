package plugins

import (
    "net/http"
    "time"
)

type NightBlockPlugin struct{}

func (n *NightBlockPlugin) Wrap(next http.Handler, config map[string]interface{}) http.Handler {
    start, _ := config["start"].(float64) // hora inicio (ej: 22)
    end, _ := config["end"].(float64)     // hora fin (ej: 6)
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        hour := float64(time.Now().Hour())
        if (start < end && hour >= start && hour < end) || (start > end && (hour >= start || hour < end)) {
            http.Error(w, "Acceso bloqueado por horario", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func init() {
    Register("nightblock", &NightBlockPlugin{})
}