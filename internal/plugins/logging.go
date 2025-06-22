package plugins

import (
    "log"
    "net/http"
)

type LoggingPlugin struct{}

func (l *LoggingPlugin) Wrap(next http.Handler, config map[string]interface{}) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("[PLUGIN-LOG] %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

func init() {
    Register("logging", &LoggingPlugin{})
}